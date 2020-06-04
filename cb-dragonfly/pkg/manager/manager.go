package manager

import (
	"bytes"
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
	Config            Config
	InfluxdDB         metricstore.Storage
	Etcd              realtimestore.Storage
	Aggregator        collector.Aggregator
	WaitGroup         *sync.WaitGroup
	UdpCOnn           *net.UDPConn
	metricL           *sync.RWMutex
	CollectorIdx      []string
	CollectorUUIDAddr map[string]*collector.MetricCollector
	AggregatingChan   map[string]*chan string
	TransmitDataChan  map[string]*chan collector.TelegrafMetric
	AgentQueueTTL     map[string]time.Time
	AgentQueueColN    map[string]int
}

// 콜렉터 매니저 초기화
func NewCollectorManager() (*CollectManager, error) {
	manager := CollectManager{}
	err := manager.LoadConfiguration()
	if err != nil {
		return nil, err
	}

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

	manager.metricL = &sync.RWMutex{}

	manager.Etcd = etcd

	manager.AgentQueueTTL = map[string]time.Time{}
	manager.AgentQueueColN = map[string]int{}

	return &manager, nil
}

// 기존의 실시간 모니터링 데이터 삭제
func (manager *CollectManager) FlushMonitoringData() error {
	// 모니터링 콜렉터 태그 정보 삭제
	//manager.Etcd.DeleteMetric("/host-list")
	manager.Etcd.DeleteMetric("/collector")

	// 실시간 모니터링 정보 삭제
	manager.Etcd.DeleteMetric("/host")
	manager.Etcd.DeleteMetric("/mon")

	manager.SetConfigurationToETCD()

	return nil
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

// TODO: 모니터링 정책 설정
func (manager *CollectManager) SetConfigurationToETCD() error {
	monConfig := MonConfig{
		AgentInterval:      manager.Config.Monitoring.AgentInterval,
		CollectorInterval:  manager.Config.Monitoring.CollectorInterval,
		SchedulingInterval: manager.Config.Monitoring.ScheduleInterval,
		MaxHostCount:       manager.Config.Monitoring.MaxHostCount,
		AgentTtl:           manager.Config.Monitoring.AgentTtl,
	}

	// TODO: 구조체 map[string]interface{} 타입으로 Unmarshal
	// TODO: 추후에 별도의 map 변환 함수 (toMap() 개발)
	reqBodyBytes := new(bytes.Buffer)
	if err := json.NewEncoder(reqBodyBytes).Encode(monConfig); err != nil {
		return err
	}
	byteData := reqBodyBytes.Bytes()

	jsonMap := map[string]interface{}{}
	if err := json.Unmarshal(byteData, &jsonMap); err != nil {
		return err
	}

	// etcd 저장소에 모니터링 정책 저장 후 결과 값 반환
	err := manager.Etcd.WriteMetric("/mon/config", jsonMap)
	if err != nil {
		return err
	}

	return nil
}

func (manager *CollectManager) CreateLoadBalancer(wg *sync.WaitGroup) error {

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
	//manager.WaitGroup = wg
	manager.UdpCOnn = udpConn

	manager.WaitGroup.Add(1)
	go manager.ManageAgentTtl(manager.WaitGroup)

	manager.WaitGroup.Add(1)
	go manager.StartLoadBalancer(manager.UdpCOnn, manager.WaitGroup)

	return nil
}

func (manager *CollectManager) StartLoadBalancer(udpConn net.PacketConn, wg *sync.WaitGroup) {

	defer wg.Done()

	for {
		metric := collector.TelegrafMetric{}
		buf := make([]byte, 1024*10)

		n, _, err := udpConn.ReadFrom(buf)

		if err != nil {
			logrus.Error("UDPLoadBalancer : failed to read bytes: ", err)
		}
		manager.metricL.Lock()
		if err := json.Unmarshal(buf[0:n], &metric); err != nil {
			logrus.Error("Failed to decode json to buf: ", string(buf[0:n]))
			continue
		}
		manager.metricL.Unlock()
		hostId := metric.Tags["hostID"].(string)

		manager.AgentQueueTTL[hostId] = time.Now()

		_, alreadyRegistered := manager.AgentQueueColN[hostId]

		if !alreadyRegistered {
			manager.metricL.Lock()
			manager.AgentQueueColN[hostId] = -1
			manager.metricL.Unlock()
		}

		err = manager.ManageAgentQueue(hostId, manager.AgentQueueColN, metric)
		if err != nil {
			logrus.Error("ManageAgentQueue Error", err)
		}
	}
}

func (manager *CollectManager) ManageAgentQueue(hostId string, AgentQueueColN map[string]int, metric collector.TelegrafMetric) error {

	colN := AgentQueueColN[hostId]
	colUUID := ""
	var cashingCAddr *collector.MetricCollector
	// Case : new Data which is not allocated at collector
	for idx, cUUID := range manager.CollectorIdx {

		cAddr := manager.CollectorUUIDAddr[cUUID]

		if cAddr != nil {
			if _, alreadyRegistered := (*cAddr).MarkingAgent[hostId]; alreadyRegistered {
				if idx != 0 {
					cashingCAddr = cAddr
					break
				}
				colUUID = manager.CollectorIdx[colN]
				*(manager.TransmitDataChan[colUUID]) <- metric
				return nil
			}
		}
	}

	for idx, cUUID := range manager.CollectorIdx {

		cAddr := manager.CollectorUUIDAddr[cUUID]

		if len((*cAddr).MarkingAgent) < manager.Config.Monitoring.MaxHostCount {

			if cashingCAddr != nil {
				delete((*cashingCAddr).MarkingAgent, hostId)
			}
			manager.metricL.Lock()
			(*cAddr).MarkingAgent[hostId] = hostId
			manager.metricL.Unlock()
			AgentQueueColN[hostId] = idx
			colN = AgentQueueColN[hostId]
			colUUID = manager.CollectorIdx[colN]
			*(manager.TransmitDataChan[colUUID]) <- metric
			return nil
		}
	}

	return nil
}

func (manager *CollectManager) ManageAgentTtl(wg *sync.WaitGroup) {

	defer wg.Done()

	//monConfig, err := manager.GetConfigInfo()
	agentTtl := manager.Config.Monitoring.AgentTtl

	for {
		currentTime := time.Now()
		if len(manager.AgentQueueTTL) != 0 {
			manager.metricL.RLock()
			for hostId, arrivedTime := range manager.AgentQueueTTL {

				if currentTime.Sub(arrivedTime) > time.Duration(agentTtl)*time.Second {
					if _, ok := manager.AgentQueueTTL[hostId]; ok {
						//manager.metricL.RLock()
						delete(manager.AgentQueueTTL, hostId)
						//manager.metricL.RUnlock()
					}
					colN := manager.AgentQueueColN[hostId]
					cUUID := ""
					if colN >= 0 && colN < len(manager.CollectorIdx) {
						cUUID = manager.CollectorIdx[colN]
					} else {
						continue
					}
					c := manager.CollectorUUIDAddr[cUUID]
					if _, ok := manager.AgentQueueColN[hostId]; ok {
						delete(manager.AgentQueueColN, hostId)
					}
					if _, ok := (*c).MarkingAgent[hostId]; ok {
						delete((*c).MarkingAgent, hostId)
					}
					err := manager.Etcd.DeleteMetric(fmt.Sprintf("/collector/%s/host/%s", cUUID, hostId))
					if err != nil {
						logrus.Error("Fail to delete hostInfo ETCD data")
					}
					err = manager.Etcd.DeleteMetric(fmt.Sprintf("/host/%s", hostId))
					if err != nil {
						logrus.Error("Fail to delete expired ETCD data")
					}
				}
			}
			manager.metricL.RUnlock()
		}
		time.Sleep(1 * time.Second)
	}
}

//func (manager *CollectManager) StartCollector(wg *sync.WaitGroup, aggregateChan chan string) error {
func (manager *CollectManager) StartCollector(wg *sync.WaitGroup) error {
	manager.WaitGroup = wg

	manager.CollectorIdx = []string{}
	manager.CollectorUUIDAddr = map[string]*collector.MetricCollector{}
	manager.AggregatingChan = map[string]*chan string{}
	manager.TransmitDataChan = map[string]*chan collector.TelegrafMetric{}

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
	mc := collector.NewMetricCollector(map[string]string{}, manager.metricL, manager.Config.Monitoring.CollectorInterval, &manager.Etcd, &manager.InfluxdDB, collector.AVG, manager.AggregatingChan, manager.TransmitDataChan)
	manager.metricL.Lock()
	manager.CollectorIdx = append(manager.CollectorIdx, mc.UUID)
	manager.metricL.Unlock()
	transmitDataChan := make(chan collector.TelegrafMetric)
	manager.WaitGroup.Add(1)
	go mc.StartCollector(manager.UdpCOnn, manager.WaitGroup, &transmitDataChan)

	manager.WaitGroup.Add(1)
	aggregateChan := make(chan string)
	go mc.StartAggregator(manager.WaitGroup, &aggregateChan)

	manager.TransmitDataChan[mc.UUID] = &transmitDataChan
	manager.CollectorUUIDAddr[mc.UUID] = &mc
	manager.AggregatingChan[mc.UUID] = &aggregateChan

	return nil
}

func (manager *CollectManager) StopCollector(uuid string) error {
	if _, ok := manager.CollectorUUIDAddr[uuid]; ok {
		// 실행 중인 콜렉터 고루틴 종료 (콜렉터 활성화 플래그 변경)
		manager.CollectorUUIDAddr[uuid].Active = false
		delete(manager.CollectorUUIDAddr, uuid)
		//*(manager.TransmitDataChan[uuid]) <- collector.TelegrafMetric{}
		//manager.CollectorCnt -= 1
		return nil
	} else {
		return errors.New(fmt.Sprintf("failed to get collector by id, uuid: %s", uuid))
	}
}

//func (manager *CollectManager) GetConfigInfo() (MonConfig, error) {
//	// etcd 저장소 조회
//	configNode, err := manager.Etcd.ReadMetric("/mon/config")
//	if err != nil {
//		return MonConfig{}, err
//	}
//	// MonConfig 매핑
//	var config MonConfig
//	err = json.Unmarshal([]byte(configNode.Value), &config)
//	if err != nil {
//		return MonConfig{}, err
//	}
//	return config, nil
//}

func (manager *CollectManager) StartAggregateScheduler(wg *sync.WaitGroup, c *map[string]*chan string) {
	defer wg.Done()
	for {
		// aggregate 주기 정보 조회
		collectorInterval := manager.Config.Monitoring.CollectorInterval

		//// Print Session Start /////
		//fmt.Print("\nTTL queue List : ")
		//sortedAgentQueueTTL := make([] int, 0)
		//for key, _ := range manager.AgentQueueTTL{
		//	value, _ := strconv.Atoi(strings.Split(key,"-")[2])
		//	sortedAgentQueueTTL = append(sortedAgentQueueTTL, value)
		//}
		//sort.Slice(sortedAgentQueueTTL, func(i, j int) bool {
		//	return sortedAgentQueueTTL[i] < sortedAgentQueueTTL[j]
		//})
		//for _, value := range sortedAgentQueueTTL {
		//	fmt.Print(value, ", ")
		//}
		//fmt.Print(fmt.Sprintf(" / Total : %d", len(sortedAgentQueueTTL)))
		//fmt.Print("\n")
		//fmt.Println("The number of collector : ", len(manager.CollectorIdx))
		//// Print Session End /////

		time.Sleep(time.Duration(collectorInterval) * time.Second)

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
		schedulingInterval := manager.Config.Monitoring.ScheduleInterval

		time.Sleep(time.Duration(schedulingInterval) * time.Second)

		// Check Scale-In/Out Logic ( len(AgentTTLQueue) 기준 Scaling In/Out)
		err := cs.CheckScaleCondition()
		if err != nil {
			logrus.Error("failed to check scale in/out condition", err)
		}
	}
}
