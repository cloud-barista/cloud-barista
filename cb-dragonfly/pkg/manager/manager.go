package manager

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cloud-barista/cb-dragonfly/pkg/collector"
	"github.com/cloud-barista/cb-dragonfly/pkg/metricstore"
	"github.com/cloud-barista/cb-dragonfly/pkg/metricstore/influxdbv1"
	"github.com/cloud-barista/cb-dragonfly/pkg/realtimestore"
	"github.com/cloud-barista/cb-dragonfly/pkg/realtimestore/etcd"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"net"
	"os"
	"sync"
	"time"
)

// TODO: implements
// TODO: 1. API Server
// TODO: 2. Scheduling Collector...
// TODO: 3. Configuring Policy...

type CollectManager struct {
	Config        Config
	InfluxdDB     metricstore.Storage
	Etcd          realtimestore.Storage
	Aggregator    collector.Aggregator
	WaitGroup     *sync.WaitGroup
	UdpCOnn       *net.UDPConn
	CollectorList map[string]*collector.MetricCollector
	CollectorChan map[string]*chan string
	HostInfo      collector.HostInfo
	HostCnt       int
}

// 콜렉터 매니저 초기화
func NewCollectorManager() (*CollectManager, error) {
	manager := CollectManager{}
	manager.LoadConfiguration()

	influxConfig := influxdbv1.Config{
		ClientOptions: []influxdbv1.ClientOptions{
			{
				URL:      manager.Config.InfluxDB.EndpointUrl,
				Username: manager.Config.InfluxDB.UserName,
				Password: manager.Config.InfluxDB.Password,
			},
		},
		Database: manager.Config.InfluxDB.Database,
	}

	// InfluxDB 연결
	influx, err := metricstore.NewStorage(metricstore.InfluxDBV1Type, influxConfig)
	if err != nil {
		logrus.Error("Failed to initialize influxDB")
		return nil, err
	}
	manager.InfluxdDB = influx

	etcdConfig := etcd.Config{
		ClientOptions: etcd.ClientOptions{
			Endpoints: manager.Config.Etcd.EndpointUrl,
		},
	}

	// etcd 연결
	etcd, err := realtimestore.NewStorage(realtimestore.ETCDV2Type, etcdConfig)
	if err != nil {
		logrus.Error("Failed to initialize etcd")
		return nil, err
	}
	manager.Etcd = etcd

	// 호스트 리스트 초기화
	hostList := collector.HostInfo{
		HostMap: &map[string]string{},
		L:       &sync.RWMutex{},
	}
	manager.HostInfo = hostList

	return &manager, nil
}

// config 파일 로드
func (manager *CollectManager) LoadConfiguration() error {
	configPath := os.Getenv("CBMON_PATH") + "/conf/config.yaml"

	bytes, err := ioutil.ReadFile(configPath)
	if err != nil {
		logrus.Error("Failed to read configuration file in: ", configPath)
		return err
	}

	err = yaml.Unmarshal(bytes, &manager.Config)
	if err != nil {
		logrus.Error("Failed to unmarshal configuration file")
		return err
	}

	return nil
}

// 기존의 실시간 모니터링 데이터 삭제
func (manager *CollectManager) FlushMonitoringData() error {
	// 모니터링 콜렉터 태그 정보 삭제
	manager.Etcd.DeleteMetric("/host-list")
	manager.Etcd.DeleteMetric("/collector")
	/*if err := manager.Etcd.DeleteMetric("/host-list"); err != nil {
		return err
	}
	if err := manager.Etcd.DeleteMetric("/collector"); err != nil {
		return err
	}*/
	// 실시간 모니터링 정보 삭제
	manager.Etcd.DeleteMetric("/host")
	/*if err := manager.Etcd.DeleteMetric("/host"); err != nil {
		return err
	}*/
	return nil
}

// TODO: 모니터링 정책 설정
func (manager *CollectManager) SetConfiguration() error {
	return nil
}

//func (manager *CollectManager) StartCollector(wg *sync.WaitGroup, aggregateChan chan string) error {
func (manager *CollectManager) StartCollector(wg *sync.WaitGroup) error {
	// UDP 서버 설정
	udpAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("0.0.0.0:%d", manager.Config.CollectManager.CollectorPort))
	if err != nil {
		logrus.Error("Failed to resolve UDP server address: ", err)
		return err
	}
	udpConn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		logrus.Error("Failed to listen UDP server address: ", err)
		return err
	}
	//defer udpConn.Close()

	// TODO: UDP 멀티 소켓 처리
	/*listenConfig := net.ListenConfig{Control: func(network, address string, c syscall.RawConn) (err error) {
		return c.Control(func(fd uintptr) {
			err = windows.SetsockoptInt(windows.Handle(fd), windows.SOL_SOCKET, windows.SO_REUSEADDR, 1)
			if err != nil {
				panic(err)
			}
		})
	}}*/
	/*udpConn, err := listenConfig.ListenPacket(context.Background(), "udp", fmt.Sprintf(":%d", manager.Config.CollectManager.CollectorPort))
	if err != nil {
		panic(err)
	}*/

	manager.WaitGroup = wg
	manager.UdpCOnn = udpConn

	manager.CollectorList = map[string]*collector.MetricCollector{}
	manager.CollectorChan = map[string]*chan string{}

	for i := 0; i < manager.Config.CollectManager.CollectorCnt; i++ {
		err := manager.CreateCollector()
		if err != nil {
			logrus.Error("failed to create collector", err)
			continue
		}
	}

	return nil
}

func (manager *CollectManager) CreateCollector() error {
	// 실시간 데이터 저장을 위한 collector 고루틴 실행
	mc := collector.NewMetricCollector(manager.Config.CollectManager.CollectorPort, manager.Config.Monitoring.CollectorInterval, &manager.Etcd, &manager.InfluxdDB, collector.AVG, &manager.HostInfo, manager.CollectorChan)
	manager.WaitGroup.Add(1)
	go mc.StartCollector(manager.UdpCOnn, manager.WaitGroup)

	// 실시간 데이터 처리를 위한 Aggregator 고루틴 실행
	aggregateChan := make(chan string)
	go mc.StartAggregator(manager.WaitGroup, &aggregateChan)

	manager.CollectorList[mc.UUID] = &mc
	manager.CollectorChan[mc.UUID] = &aggregateChan

	return nil
}

func (manager *CollectManager) StopCollector(uuid string) error {
	if _, ok := manager.CollectorList[uuid]; ok {
		// 실행 중인 콜렉터 고루틴 종료 (콜렉터 활성화 플래그 변경)
		manager.CollectorList[uuid].Active = false
		delete(manager.CollectorList, uuid)
		fmt.Println(fmt.Sprintf("###############################################"))
		fmt.Println(fmt.Sprintf("[%s] stop collector", uuid))
		for key := range manager.CollectorList {
			fmt.Println(key)
		}
		fmt.Println(fmt.Sprintf("###############################################"))
		return nil
	} else {
		return errors.New(fmt.Sprintf("failed to get collector by id, uuid: %s", uuid))
	}
}

func (manager *CollectManager) GetConfigInfo() (MonConfig, error) {
	// etcd 저장소 조회
	configNode, err := manager.Etcd.ReadMetric("/mon/config")
	if err != nil {
		return MonConfig{}, err
	}
	// MonConfig 매핑
	var config MonConfig
	err = json.Unmarshal([]byte(configNode.Value), &config)
	if err != nil {
		return MonConfig{}, err
	}
	return config, nil
}

func (manager *CollectManager) StartAggregateScheduler(wg *sync.WaitGroup, c *map[string]*chan string) {
	defer wg.Done()
	for {
		/*select {
		case <-ctx.Done():
			logrus.Debug("Stop scheduling for aggregate metric")
			return
		default:
		}*/
		// aggregate 주기 정보 조회
		monConfig, err := manager.GetConfigInfo()
		if err != nil {
			logrus.Error("failed to get monitoring config info", err)
		}

		time.Sleep(time.Duration(monConfig.CollectorInterval) * time.Second)

		manager.HostCnt = len(*manager.HostInfo.HostMap)

		for _, channel := range *c {
			*channel <- "aggregate"
		}
	}
}

// 콜렉터 스케일 인/아웃 관리 스케줄러
func (manager *CollectManager) StartScaleScheduler(wg *sync.WaitGroup) {
	defer wg.Done()
	cs := NewCollectorScheduler(manager)
	for {
		// 스케줄링 주기 정보 조회
		monConfig, err := manager.GetConfigInfo()
		if err != nil {
			logrus.Error("failed to get monitoring config info", err)
		}

		time.Sleep(time.Duration(monConfig.SchedulingInterval) * time.Second)

		// Check Scale-In/Out Logic (호스트 수 기준 Scaling In/Out)
		err = cs.CheckScaleCondition()
		if err != nil {
			logrus.Error("failed to check scale in/out condition", err)
		}
	}
}
