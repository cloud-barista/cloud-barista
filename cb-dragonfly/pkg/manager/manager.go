package manager

import (
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"

	"github.com/cloud-barista/cb-dragonfly/pkg/collector"
	"github.com/cloud-barista/cb-dragonfly/pkg/config"
	"github.com/cloud-barista/cb-dragonfly/pkg/localstore"
	"github.com/cloud-barista/cb-dragonfly/pkg/metricstore/influxdb/influxdbv1"
	"github.com/cloud-barista/cb-dragonfly/pkg/types"
)

// TODO: implements
// TODO: VM OR CONTAINER BASED COLLECTOR SCALE OUT => CHANNEL TO API

type CollectManager struct {
	CollectorGroupManageMap map[int][]*collector.MetricCollector
	WaitGroup               *sync.WaitGroup
	collectorPolicy         string
}

func NewCollectorManager() (*CollectManager, error) {
	manager := CollectManager{}

	influxConfig := influxdbv1.Config{
		ClientOptions: []influxdbv1.ClientOptions{
			{
				URL: fmt.Sprintf("%s:%d", config.GetInstance().GetInfluxDBConfig().EndpointUrl, config.GetInstance().GetInfluxDBConfig().InternalPort),
				//URL: "http://192.168.130.7:28086",
				Username: config.GetInstance().GetInfluxDBConfig().UserName,
				Password: config.GetInstance().GetInfluxDBConfig().Password,
			},
		},
		Database: config.GetInstance().GetInfluxDBConfig().Database,
	}

	err := influxdbv1.Initialize(influxConfig)
	if err != nil {
		logrus.Error("Failed to initialize influxDB")
		return nil, err
	}

	timeout := time.Duration(1 * time.Second)
	_, err = net.DialTimeout("tcp", fmt.Sprintf("%s:%d", config.GetInstance().GetKafkaConfig().GetKafkaEndpointUrl(), config.GetInstance().GetKafkaConfig().InternalPort), timeout)
	if err != nil {
		fmt.Printf("%s %s \n", "kafka is not responding", err.Error())
		logrus.Error(err)
		return nil, err
	} else {
		fmt.Printf("kafka is responding")
	}

	manager.collectorPolicy = strings.ToUpper(config.GetInstance().Monitoring.MonitoringPolicy)
	manager.CollectorGroupManageMap = map[int][]*collector.MetricCollector{}

	return &manager, nil
}

func (manager *CollectManager) FlushMonitoringData() {
	err := os.RemoveAll("./meta_db")
	if err != nil {
		fmt.Println(err)
	}
	manager.SetConfigurationToMemoryDB()
}

func (manager *CollectManager) SetConfigurationToMemoryDB() {
	monConfigMap := map[string]interface{}{}
	mapstructure.Decode(config.GetInstance().Monitoring, &monConfigMap)
	for key, val := range monConfigMap {
		err := localstore.GetInstance().StorePut(types.MONCONFIG+"/"+key, fmt.Sprintf("%v", val))
		if err != nil {
			logrus.Debug(err)
		}
	}
}

func (manager *CollectManager) StartCollectorGroup(wg *sync.WaitGroup) error {
	manager.WaitGroup = wg
	if manager.collectorPolicy == types.AGENTCOUNT {
		startCollectorGroupCnt := config.GetInstance().CollectManager.CollectorGroupCnt
		for i := 0; i < startCollectorGroupCnt; i++ {
			err := manager.CreateCollectorGroup()
			if err != nil {
				logrus.Error("failed to create collector group", err)
				return err
			}
		}
	}
	if manager.collectorPolicy == types.CSP {
		for i := 0; i < types.TOTALCSPCNT; i++ {
			err := manager.CreateCollectorGroup()
			if err != nil {
				logrus.Error("failed to create collector group", err)
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
	//for i := 0; i < config.GetInstance().CollectManager.GroupPerCollectCnt; i++ {
	mc, err := collector.NewMetricCollector(
		collector.AVG,
		collectorCreateOrder,
	)
	if err != nil {
		return err
	}
	collectorList = append(collectorList, &mc)
	go func() {
		err := mc.Collector(manager.WaitGroup)
		if err != nil {
			logrus.Debug("Fail to create Collector")
		}
	}()
	//}
	manager.CollectorGroupManageMap[collectorCreateOrder] = collectorList
	return nil
}

func (manager *CollectManager) StopCollectorGroup() error {
	collectorIdx := len(manager.CollectorGroupManageMap) - 1
	if collectorIdx == 0 {
		return nil
	} else {
		for _, cAddr := range manager.CollectorGroupManageMap[collectorIdx] {
			(*cAddr).Active = false
		}
		delete(manager.CollectorGroupManageMap, collectorIdx)
	}
	return nil
}

func (manager *CollectManager) StartScheduler(wg *sync.WaitGroup) error {
	defer wg.Done()
	scheduler, erro := NewCollectorScheduler(manager)
	if erro != nil {
		logrus.Error("Failed to initialize influxDB")
		return erro
	}
	go func() {
		err := scheduler.Scheduler()
		if err != nil {
			erro = err
		}
	}()
	if erro != nil {
		logrus.Error("Failed to make scheduler")
		return erro
	}
	return nil
}
