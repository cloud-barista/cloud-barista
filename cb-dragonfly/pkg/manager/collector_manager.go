package manager

import (
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/cloud-barista/cb-dragonfly/pkg/cbstore"
	"github.com/cloud-barista/cb-dragonfly/pkg/collector"
	"github.com/cloud-barista/cb-dragonfly/pkg/config"
	"github.com/cloud-barista/cb-dragonfly/pkg/types"
	"github.com/cloud-barista/cb-dragonfly/pkg/util"
	"github.com/mitchellh/mapstructure"
)

const (
	KafkaConnectionRetryCnt      = 5
	KafkaConnectionRetryInterval = 5
)

// TODO: VM OR CONTAINER BASED COLLECTOR SCALE OUT => CHANNEL TO API

type CollectManager struct {
	CollectorGroupManageMap map[int][]*collector.MetricCollector
	WaitGroup               *sync.WaitGroup
	collectorPolicy         string
}

func NewCollectorManager() (*CollectManager, error) {
	manager := CollectManager{}

	retryCnt := KafkaConnectionRetryCnt
	waitInterval := KafkaConnectionRetryInterval * time.Second
	for i := 0; i <= retryCnt; i++ {
		_, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", config.GetInstance().GetKafkaConfig().GetKafkaEndpointUrl(), config.GetInstance().GetKafkaConfig().InternalPort), waitInterval)
		if err != nil {
			if i == retryCnt {
				util.GetLogger().Error("kafka is not responding %s", "kafka is not responding ", err.Error())
				return nil, err
			} else {
				util.GetLogger().Error(fmt.Sprintf("\nRetry Attempt : %d, Now ReConn to Kafka... ERR MSG: %s ", i+1, err.Error()))
			}
		} else {
			util.GetLogger().Info("kafka is responding")
			break
		}
		time.Sleep(KafkaConnectionRetryInterval * time.Second)
	}

	manager.collectorPolicy = strings.ToUpper(config.GetInstance().Monitoring.MonitoringPolicy)
	manager.CollectorGroupManageMap = map[int][]*collector.MetricCollector{}

	return &manager, nil
}

func (manager *CollectManager) FlushMonitoringData() {
	err := os.RemoveAll("./meta_db")
	if err != nil {
		util.GetLogger().Error(fmt.Sprintf("failed to flush monitoring data error=%s", err.Error()))
	}
	manager.SetConfigurationToMemoryDB()
}

func (manager *CollectManager) SetConfigurationToMemoryDB() {
	monConfigMap := map[string]interface{}{}
	err := mapstructure.Decode(config.GetInstance().Monitoring, &monConfigMap)
	if err != nil {
		util.GetLogger().Error(fmt.Sprintf("failed to decode monConfigMap, error=%s", err))
	}
	for key, val := range monConfigMap {
		err := cbstore.GetInstance().StorePut(types.MoNConfig+"/"+key, fmt.Sprintf("%v", val))
		if err != nil {
			util.GetLogger().Error(fmt.Sprintf("failed to put monitoring configuration info, error=%s", err))
		}
	}
}

func (manager *CollectManager) StartCollectorGroup(wg *sync.WaitGroup) error {
	manager.WaitGroup = wg
	if manager.collectorPolicy == types.AgentCnt {
		startCollectorGroupCnt := config.GetInstance().CollectManager.CollectorGroupCnt
		for i := 0; i < startCollectorGroupCnt; i++ {
			err := manager.CreateCollectorGroup()
			if err != nil {
				util.GetLogger().Error("failed to create collector group", err)
				return err
			}
		}
	}
	if manager.collectorPolicy == types.CSP {
		for i := 0; i < types.TotalCspCnt; i++ {
			err := manager.CreateCollectorGroup()
			if err != nil {
				util.GetLogger().Error("failed to create collector group", err)
				return err
			}
		}
	}
	return nil
}

func (manager *CollectManager) CreateCollectorGroup() error {

	manager.WaitGroup.Add(1)
	collectorCreateOrder := len(manager.CollectorGroupManageMap)
	var collectorList []*collector.MetricCollector
	mc, err := collector.NewMetricCollector(
		types.AVG,
		collectorCreateOrder,
	)
	if err != nil {
		return err
	}
	collectorList = append(collectorList, &mc)
	go func() {
		err := mc.Collector(manager.WaitGroup)
		if err != nil {
			util.GetLogger().Error("failed to  create Collector")
		}
	}()
	manager.CollectorGroupManageMap[collectorCreateOrder] = collectorList
	return nil
}

func (manager *CollectManager) StopCollectorGroup() error {
	collectorIdx := len(manager.CollectorGroupManageMap) - 1
	if collectorIdx == 0 {
		return nil
	} else {
		for _, cAddr := range manager.CollectorGroupManageMap[collectorIdx] {
			(*cAddr).Ch <- "close"
		}
		delete(manager.CollectorGroupManageMap, collectorIdx)
	}
	return nil
}

func (manager *CollectManager) StartScheduler(wg *sync.WaitGroup) error {
	defer wg.Done()
	scheduler, erro := NewCollectorScheduler(manager)
	if erro != nil {
		util.GetLogger().Error("Failed to initialize influxDB")
		return erro
	}
	go func() {
		err := scheduler.Scheduler()
		if err != nil {
			erro = err
		}
	}()
	if erro != nil {
		util.GetLogger().Error("Failed to make scheduler")
		return erro
	}
	return nil
}
