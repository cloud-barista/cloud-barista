package collector

import (
	"encoding/json"
	"fmt"
	"github.com/cloud-barista/cb-dragonfly/pkg/metricstore"
	"github.com/cloud-barista/cb-dragonfly/pkg/realtimestore"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"net"
	"strings"
	"sync"
	"time"
)

type MetricCollector struct {
	MarkingAgent      map[string]string
	UUID              string
	AggregateInterval int
	InfluxDB          metricstore.Storage
	metricL           *sync.RWMutex
	Etcd              realtimestore.Storage
	Aggregator        Aggregator
	//HostInfo          *HostInfo
	AggregatingChan  map[string]*chan string
	TransmitDataChan map[string]*chan TelegrafMetric
	Active           bool
	//UdpConn         net.UDPConn
}

type TelegrafMetric struct {
	Name      string                 `json:"name"`
	Tags      map[string]interface{} `json:"tags"`
	Fields    map[string]interface{} `json:"fields"`
	Timestamp int64                  `json:"timestamp"`
	TagInfo   map[string]interface{} `json:"tagInfo"`
}

type TagMetric struct {
	Tags map[string]interface{} `json:"tags"`
}

type DeviceInfo struct {
	HostID     string `json:"host_id"`
	MetricName string `json:"host_id"`
	DeviceName string `json:"device_name"`
}

// 메트릭 콜렉터 초기화
func NewMetricCollector(markingAgent map[string]string, mutexLock *sync.RWMutex, interval int, etcd *realtimestore.Storage, influxDB *metricstore.Storage, aggregateType AggregateType /*hostList *HostInfo, */, aggregatingChan map[string]*chan string, transmitDataChan map[string]*chan TelegrafMetric) MetricCollector {

	// UUID 생성
	uuid := uuid.New().String()

	// 모니터링 메트릭 Collector 초기화
	mc := MetricCollector{
		MarkingAgent:      markingAgent,
		UUID:              uuid,
		AggregateInterval: interval,
		Etcd:              *etcd,
		metricL:           mutexLock,
		Aggregator: Aggregator{
			Etcd:          *etcd,
			InfluxDB:      *influxDB,
			AggregateType: aggregateType,
		},
		//HostInfo:      hostList,
		AggregatingChan:  aggregatingChan,
		TransmitDataChan: transmitDataChan,
		Active:           true,
	}

	return mc
}

//func (mc *MetricCollector) Start(listenConfig net.ListenConfig, wg *sync.WaitGroup) {
func (mc *MetricCollector) StartCollector(udpConn net.PacketConn, wg *sync.WaitGroup, ch *chan TelegrafMetric) error {

	// Telegraf 에이전트에서 보내는 모니터링 메트릭 수집
	defer wg.Done()
	for {

		metric := TelegrafMetric{}
		select {
		case metric = <-*ch:
			if !mc.Active {
				// tagging 채널 삭제
				close(*ch)
				delete(mc.TransmitDataChan, mc.UUID)
				break
			}

			goto Start
		}

	Start:

		hostId, ok := metric.Tags["hostID"].(string)

		if !ok {
			continue
		}
		mc.metricL.RLock()
		if _, ok := mc.MarkingAgent[hostId]; !ok {
			continue
		}
		mc.metricL.RUnlock()
		collectorInfo := fmt.Sprintf("/collector/%s/host/%s", mc.UUID, hostId)
		err := mc.Etcd.WriteMetric(collectorInfo, "")

		if err != nil {
			return err
		}

		curTimestamp := time.Now().Unix()
		var diskName string
		var metricKey string
		var osTypeKey string

		mc.metricL.RLock()
		switch strings.ToLower(metric.Name) {
		case "disk":
			diskName = metric.Tags["device"].(string)
			diskName = strings.ReplaceAll(diskName, "/", "%")
			metricKey = fmt.Sprintf("/host/%s/metric/%s/%s/%d", hostId, metric.Name, diskName, curTimestamp)
		case "diskio":
			diskName := metric.Tags["name"].(string)
			diskName = strings.ReplaceAll(diskName, "/", "%")
			metricKey = fmt.Sprintf("/host/%s/metric/%s/%s/%d", hostId, metric.Name, diskName, curTimestamp)
		default:
			metricKey = fmt.Sprintf("/host/%s/metric/%s/%d", hostId, metric.Name, curTimestamp)
		}
		mc.metricL.RUnlock()

		if err := mc.Etcd.WriteMetric(metricKey, metric.Fields); err != nil {
			logrus.Error(err)
		}

		metric.TagInfo = map[string]interface{}{}
		metric.TagInfo["mcisId"] = hostId
		metric.TagInfo["hostId"] = hostId
		metric.TagInfo["osType"] = metric.Tags["osType"].(string)

		osTypeKey = fmt.Sprintf("/host/%s/tag", hostId)

		if err := mc.Etcd.WriteMetric(osTypeKey, metric.TagInfo); err != nil {
			logrus.Error(err)
		}
	}
}

func (mc *MetricCollector) StartAggregator(wg *sync.WaitGroup, c *chan string) {
	defer wg.Done()
	for {
		select {
		// check aggregating model
		case <-*c:
			logrus.Debug("======================================================================")
			logrus.Debug("[" + mc.UUID + "]Start Aggregate!!")
			fmt.Println("["+mc.UUID+"] Start Aggregate!!", time.Now())
			//fmt.Println("mc.UUID : ", mc.UUID)
			err := mc.Aggregator.AggregateMetric(mc.UUID)
			if err != nil {
				logrus.Error("["+mc.UUID+"]Failed to aggregate meric", err)
			}
			logrus.Debug("======================================================================")

			//// Print Session Start /////
			//fmt.Print("mc.MarkingAgent : ")
			//sortedMarkingAgent := make([] int, 0)
			//for key, _ := range mc.MarkingAgent {
			//	value, _ := strconv.Atoi(strings.Split(key,"-")[2])
			//	sortedMarkingAgent = append(sortedMarkingAgent, value)
			//}
			//sort.Slice(sortedMarkingAgent, func(i, j int) bool {
			//	return sortedMarkingAgent[i] < sortedMarkingAgent[j]
			//})
			//for _, value := range sortedMarkingAgent {
			//	fmt.Print(value, ", ")
			//}
			//fmt.Print(fmt.Sprintf(" / Total : %d", len(sortedMarkingAgent)))
			//fmt.Print("\n")
			//// Print Session End /////

			// 콜렉터 비활성화 시 aggregate 채널 삭제
			if !mc.Active {
				// aggregate 채널 삭제
				fmt.Println("Deleting aggregate channel!")
				close(*c)
				delete(mc.AggregatingChan, mc.UUID)
				return
			}
		}
	}
}

func (mc *MetricCollector) MyMarshal(metric interface{}) (string, error) {
	var metricVal string

	_, ok := metric.(map[string]interface{})
	if ok {
		mc.metricL.Lock()
		bytes, err := json.Marshal(metric)
		mc.metricL.Unlock()
		if err != nil {
			logrus.Error("Failed to marshaling realtime monitoring data to JSON: ", err)
			return "", err
		}
		metricVal = fmt.Sprintf("%s", bytes)
	} else {
		metricVal = metric.(string)
	}

	return metricVal, nil
}

/*
func (mc *MetricCollector) UntagHost() error {

	// 현재 콜렉터에 태그된 호스트 목록 가져오기
	var hostArr []string
	tagKey := fmt.Sprintf("/collector/%s/host", mc.UUID)
	resp, err := mc.Etcd.ReadMetric(tagKey)
	if err != nil {
		return err
	}

	// 호스트 목록 슬라이스 생성
	for _, vm := range resp.Nodes {
		hostId := strings.Split(vm.Key, "/")[4]
		hostArr = append(hostArr, hostId)
	}

	// 전체 호스트 목록에서 VM 태그 목록 삭제
	for _, hostId := range hostArr {
		hostKey := fmt.Sprintf("/host-list/%s", hostId)
		err := mc.Etcd.DeleteMetric(hostKey)
		if err != nil {
			logrus.Error("Failed to untag VM", err)
			return err
		}
	}
	mc.HostInfo.DeleteHost(hostArr)

	// 콜렉터에 등록된 VM 태그 목록 삭제
	tagKey = fmt.Sprintf("/collector/%s", mc.UUID)
	err = mc.Etcd.DeleteMetric(tagKey)
	if err != nil {
		logrus.Error("Failed to untag VM", err)
		return err
	}

	return nil
}

func (mc *MetricCollector) TagHost(hostId string) error {

	// 호스트 목록에 등록되어 있는 지 체크
	isTagged := true
	hostKey := fmt.Sprintf("/host-list/%s", hostId)
	_, err := mc.Etcd.ReadMetric(hostKey)

	//fmt.Println(hostId)
	if err != nil {
		if v, ok := err.(client.Error); ok {
			if v.Code == 100 { // ErrorCode 100 = Key Not Found Error
				isTagged = false
			}
		} else {
			logrus.Error("Failed to get host-list", err)
			return err
		}
	}

	// TODO: 추후 로컬 변수가 아니라 etcd 기준으로 Mutex 처리
	// 호스트 목록에 등록되지 않았지만 내부 로컬 변수에 남아있는 데이터 삭제 처리
	if !isTagged && mc.HostInfo.GetHostById(hostId) != "" {
		//fmt.Println(hostId)
		mc.HostInfo.DeleteHost([]string{hostId})
	}

	// 등록되어 있지 않은 호스트라면 호스트 목록에 등록 후 현재 콜렉터 기준으로 태깅
	if !isTagged && mc.HostInfo.GetHostById(hostId) == "" {
		// 호스트 목록에 등록
		//fmt.Println(hostId)
		mc.HostInfo.AddHost(hostId)
		err := mc.Etcd.WriteMetric(hostKey, hostKey)
		if err != nil {
			return err
		}
		// 현재 콜렉터 기준 태깅
		tagKey := fmt.Sprintf("/collector/%s/host/%s", mc.UUID, hostId)
		err = mc.Etcd.WriteMetric(tagKey, "")
		if err != nil {
			return err
		}
	}

	return nil
}

func (mc *MetricCollector) UntagHost() error {

	// 현재 콜렉터에 태그된 호스트 목록 가져오기
	var hostArr []string
	tagKey := fmt.Sprintf("/collector/%s/host", mc.UUID)
	resp, err := mc.Etcd.ReadMetric(tagKey)
	if err != nil {
		return err
	}

	// 호스트 목록 슬라이스 생성
	for _, vm := range resp.Nodes {
		hostId := strings.Split(vm.Key, "/")[4]
		hostArr = append(hostArr, hostId)
	}

	// 전체 호스트 목록에서 VM 태그 목록 삭제
	for _, hostId := range hostArr {
		hostKey := fmt.Sprintf("/host-list/%s", hostId)
		err := mc.Etcd.DeleteMetric(hostKey)
		if err != nil {
			logrus.Error("Failed to untag VM", err)
			return err
		}
	}
	mc.HostInfo.DeleteHost(hostArr)

	// 콜렉터에 등록된 VM 태그 목록 삭제
	tagKey = fmt.Sprintf("/collector/%s", mc.UUID)
	err = mc.Etcd.DeleteMetric(tagKey)
	if err != nil {
		logrus.Error("Failed to untag VM", err)
		return err
	}

	return nil
}
*/
