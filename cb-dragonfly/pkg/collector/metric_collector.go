package collector

import (
	"encoding/json"
	"fmt"
	"github.com/cloud-barista/cb-dragonfly/pkg/metricstore"
	"github.com/cloud-barista/cb-dragonfly/pkg/realtimestore"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"go.etcd.io/etcd/client"
	"net"
	"strings"
	"sync"
	"time"
)

type MetricCollector struct {
	UUID              string
	UDPPort           int
	AggregateInterval int
	InfluxDB          metricstore.Storage
	Etcd              realtimestore.Storage
	Aggregator        Aggregator
	HostInfo          *HostInfo
	CollectorChan     map[string]*chan string
	Active            bool
	//UdpConn         net.UDPConn
}

type TelegrafMetric struct {
	Name      string                 `json:"name"`
	Tags      map[string]interface{} `json:"tags"`
	Fields    map[string]interface{} `json:"fields"`
	Timestamp int64                  `json:"timestamp"`
}

type DeviceInfo struct {
	HostID     string `json:"host_id"`
	MetricName string `json:"host_id"`
	DeviceName string `json:"device_name"`
}

// 메트릭 콜렉터 초기화
func NewMetricCollector(udpPort int, interval int, etcd *realtimestore.Storage, influxDB *metricstore.Storage, aggregateType AggregateType, hostList *HostInfo, collectorChan map[string]*chan string) MetricCollector {

	// UUID 생성
	uuid := uuid.New().String()

	// 모니터링 메트릭 Collector 초기화
	mc := MetricCollector{
		UUID:              uuid,
		UDPPort:           udpPort,
		AggregateInterval: interval,
		Etcd:              *etcd,
		Aggregator: Aggregator{
			Etcd:          *etcd,
			InfluxDB:      *influxDB,
			AggregateType: aggregateType,
		},
		HostInfo:      hostList,
		CollectorChan: collectorChan,
		Active:        true,
	}

	return mc
}

//func (mc *MetricCollector) Start(listenConfig net.ListenConfig, wg *sync.WaitGroup) {
func (mc *MetricCollector) StartCollector(udpConn net.PacketConn, wg *sync.WaitGroup) {
	// TODO: UDP 멀티 소켓 처리
	/*udpConn, err := listenConfig.ListenPacket(context.Background(), "udp", fmt.Sprintf(":%d", mc.UDPPort))
	if err != nil {
		panic(err)
	}*/

	// Telegraf 에이전트에서 보내는 모니터링 메트릭 수집
	defer wg.Done()

	for {

		if !mc.Active {
			fmt.Println(fmt.Sprintf("###############################################"))
			fmt.Println(fmt.Sprintf("[%s] stop collector", mc.UUID))
			fmt.Println(fmt.Sprintf("###############################################"))
			break
		}

		//logrus.Debug("[" + mc.UUID + "] Waiting...")

		buf := make([]byte, 1024*10)
		n, _, err := udpConn.ReadFrom(buf)
		if err != nil {
			logrus.Error("[+"+mc.UUID+"] Failed to read bytes: ", err)
			continue
		}
		metric := TelegrafMetric{}
		if err := json.Unmarshal(buf[0:n], &metric); err != nil {
			logrus.Error("[+"+mc.UUID+"] Failed to decode json to buf: ", string(buf[0:n]))
			continue
		}

		hostId := metric.Tags["hostID"].(string)

		// Tagging host
		if mc.Active {
			err = mc.TagHost(hostId)
			if err != nil {
				logrus.Error("[+"+mc.UUID+"] Failed to tagging host: ", err)
				continue
			}
		}

		curTimestamp := time.Now().Unix()
		var diskName string
		var metricKey string

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

		if err := mc.Etcd.WriteMetric(metricKey, metric.Fields); err != nil {
			logrus.Error(err)
		}

		/*
			host := metric.Tags["host"].(string)
			logrus.Debug("======================================================================")
			logrus.Debugf("UUID: %s", mc.UUID)
			logrus.Debugf("From %s", addr)
			logrus.Debugf("Metric: %v", metric)
			logrus.Debugf("Name: %s", metric.Name)
			logrus.Debugf("Tags: %s", metric.Tags)
			logrus.Debugf("Fields: %s", metric.Fields) // TODO: 수집 시 파싱 (실시간 데이터 처리 위해서)
			logrus.Debugf("Host: %s", host)
			logrus.Debug("======================================================================")
		*/
	}
}

func (mc *MetricCollector) Stop() {
	// TODO: 모니터링 메트릭 수집 고루틴 종료
}

func (mc *MetricCollector) StartAggregator(wg *sync.WaitGroup, c *chan string) {
	defer wg.Done()
	for {
		select {
		// TODO: 모니터링 메트릭 Aggregate 고루틴 종료
		case <-*c:
			logrus.Debug("======================================================================")
			logrus.Debug("[" + mc.UUID + "]Start Aggregate!!")
			err := mc.Aggregator.AggregateMetric(mc.UUID)
			if err != nil {
				logrus.Error("["+mc.UUID+"]Failed to aggregate meric", err)
			}
			err = mc.UntagHost()
			if err != nil {
				logrus.Error("["+mc.UUID+"]Failed to untag host", err)
			}
			logrus.Debug("======================================================================")

			// 콜렉터 비활성화 시 aggregate 채널 삭제
			if !mc.Active {
				fmt.Println("###############################################")
				fmt.Println(fmt.Sprintf("[%s] stop aggregator", mc.UUID))
				fmt.Println("###############################################")
				// aggregate 채널 삭제
				close(*c)
				delete(mc.CollectorChan, mc.UUID)
				return
			}
		}
	}
}

func (mc *MetricCollector) TagHost(hostId string) error {

	// 전체 호스트 목록에 등록되어 있는 지 체크
	isTagged := true
	hostKey := fmt.Sprintf("/host-list/%s", hostId)
	_, err := mc.Etcd.ReadMetric(hostKey)
	if err != nil {
		etcdErr := err.(client.Error)
		if etcdErr.Code == 100 { // ErrorCode 100 = Key Not Found Error
			isTagged = false
		}
	}

	/*fmt.Println("===================================================================")
	fmt.Println("[" + mc.UUID + "] Get Tag info")
	fmt.Println(mc.HostInfo.HostMap)
	fmt.Println("isTagged=" + strconv.FormatBool(isTagged))
	fmt.Println("===================================================================")*/

	// 등록되어 있지 않은 호스트라면 호스트 목록에 등록 후 현재 콜렉터 기준으로 태깅
	if !isTagged && mc.HostInfo.GetHostById(hostId) == "" {
		// 전체 호스트 목록에 등록
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

		/*fmt.Println("===================================================================")
		fmt.Println("[" + mc.UUID + "] Add Tag info")
		fmt.Println(*mc.HostInfo)
		fmt.Println("===================================================================")*/
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
