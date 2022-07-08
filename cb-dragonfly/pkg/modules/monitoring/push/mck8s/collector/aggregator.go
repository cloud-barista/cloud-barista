package collector

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	agentmetadata "github.com/cloud-barista/cb-dragonfly/pkg/api/core/agent/common"
	v1 "github.com/cloud-barista/cb-dragonfly/pkg/storage/metricstore/influxdb/v1"
	"github.com/cloud-barista/cb-dragonfly/pkg/types"
	"github.com/cloud-barista/cb-dragonfly/pkg/util"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/thoas/go-funk"
)

type TelegrafMetric struct {
	Name      string                 `json:"name"`
	Tags      map[string]string      `json:"tags"`
	Fields    map[string]interface{} `json:"fields"`
	Timestamp int64                  `json:"timestamp"`
}

type Aggregator struct {
	CreateOrder   int
	AggregateType types.AggregateType
}

func (a *Aggregator) AggregateMetric(kafkaAdminClient *kafka.AdminClient, kafkaConsumerConn *kafka.Consumer, topic string) {
	curTime := time.Now().Unix()
	reconnectTry := 0
	var topicMsgBytes [][]byte

	// 토픽 메세지 조회
	for {

		reconnectTry++

		// 토픽 조회 재시도 횟수 제한 체크
		if reconnectTry >= types.ReadConnectionTimeout {
			//errMsg := fmt.Sprintf("exceed maximum kafka connection try, try count=%d", reconnectTry)
			//util.GetLogger().Info(errMsg)

			// 토픽 조회 재시도 횟수 제한 초과 시 1초 sleep 처리
			time.Sleep(1 * time.Second)

			break
		}

		topicMsg, err := kafkaConsumerConn.ReadMessage(5 * time.Second)
		if err != nil {
			//errMsg := fmt.Sprintf("fail to read topic message with topic %s, error=%s", topic, err)
			//util.GetLogger().Info(errMsg)

			// 메세지 큐 조회 에러 시 1초 sleep 처리
			time.Sleep(1 * time.Second)

			continue
		}

		if topicMsg != nil {
			// 토픽 메세지 저장
			topicMsgBytes = append(topicMsgBytes, topicMsg.Value)
			//fmt.Println(fmt.Sprintf("#### Group_%d MCK8S collector - add topic %d ####", a.CreateOrder, len(topicMsgBytes)))
			// 토픽 생성 시간 체크
			if topicMsg.Timestamp.Unix() > curTime {
				break
			}
			// 토픽 메세지 및 타임아웃 정보 초기화
			reconnectTry = 0
			topicMsg = nil
		}
	}

	// 에이전트 메타데이터 정보 조회
	agentInfo, err := agentmetadata.GetAgentByUUID(topic)
	if err != nil {
		//errMsg := fmt.Sprintf("failed to get agent metadata with UUID %s, error=%s", topic, err.Error())
		//util.GetLogger().Error(errMsg)
	}

	// 토픽에 모니터링 데이더가 특정 횟수 이상 쌓이지 않을 경우 에이전트 Unhealthy 처리
	if len(topicMsgBytes) == 0 {
		//util.GetLogger().Info("failed to get monitoring data from kafka, data bytes is zero")

		// 이미 에이전트 헬스체크 상태가 Unhealthy 경우에는 메타데이터 업데이트 스킵
		if agentInfo == nil {
			return
		}
		if agentmetadata.AgentHealth(agentInfo.AgentHealth) == agentmetadata.Unhealthy {
			return
		}

		// 에이전트 메타데이터 헬스체크 상태 설정 변경
		updatedAgentInfo := agentmetadata.AgentInstallInfo{
			ServiceType: agentInfo.ServiceType,
			NsId:        agentInfo.NsId,
			Mck8sId:     agentInfo.Mck8sId,
		}
		agentInfo.AgentUnhealthyRespCnt += 1
		curAgentHealthy := agentmetadata.Healthy
		if agentInfo.AgentUnhealthyRespCnt > 5 {
			curAgentHealthy = agentmetadata.Unhealthy
		}

		// 에이전트 메타데이터 헬스체크 상태 변경
		_, _, err = agentmetadata.PutAgent(updatedAgentInfo, agentInfo.AgentUnhealthyRespCnt, agentmetadata.Enable, curAgentHealthy)
		if err != nil {
			util.GetLogger().Error(err)
		}
		fmt.Printf("[%s] <MCK8S> Group_%d collector - update AgentUnhealthyRespCnt %d\n", time.Now().Format(time.RFC3339), a.CreateOrder, agentInfo.AgentUnhealthyRespCnt)

		if curAgentHealthy == agentmetadata.Unhealthy {
			fmt.Printf("[%s] <MCK8S> Group_%d collector - delete Topic %s\n", time.Now().Format(time.RFC3339), a.CreateOrder, topic)
			_, err := kafkaAdminClient.DeleteTopics(context.Background(), []string{topic})
			if err != nil {
				errMsg := fmt.Sprintf("failed to delete topic %s, error=%s", topic, err.Error())
				util.GetLogger().Error(errMsg)
			}
		}

		return
	}

	// 토픽 데이터 처리 시 에이전트 메타데이터 헬스상태 변경
	if agentInfo != nil {
		if agentmetadata.AgentHealth(agentInfo.AgentHealth) == agentmetadata.Unhealthy {
			// 에이전트 메타데이터 헬스체크 상태 변경
			updatedAgentInfo := agentmetadata.AgentInstallInfo{
				ServiceType: agentInfo.ServiceType,
				NsId:        agentInfo.NsId,
				Mck8sId:     agentInfo.Mck8sId,
			}
			_, _, err = agentmetadata.PutAgent(updatedAgentInfo, 0, agentmetadata.Enable, agentmetadata.Healthy)
			if err != nil {
				util.GetLogger().Error(err)
			}
			fmt.Printf("[%s] <MCK8S> Group_%d collector - update AgentStatus %s\n", time.Now().Format(time.RFC3339), a.CreateOrder, agentmetadata.Healthy)
		}
	}

	// 토픽 메세지 파싱
	metrics := make([]TelegrafMetric, len(topicMsgBytes))
	for idx, topicMsgByte := range topicMsgBytes {
		var metric TelegrafMetric
		err := json.Unmarshal(topicMsgByte, &metric)
		if err == nil {
			metrics[idx] = metric
		}
	}

	fmt.Printf("[%s] <MCK8S> Collected Topic %s : %d\n", time.Now().Format(time.RFC3339), topic, len(topicMsgBytes))
	// 모니터링 데이터 처리
	a.Aggregate(metrics)
}

// Aggregate 쿠버네티스 노드, 파드 메트릭 처리 및 저장
func (a *Aggregator) Aggregate(metrics []TelegrafMetric) {

	fmt.Printf("[%s] <MCK8S> EXECUTE Group_%d aggregator\n", time.Now().Format(time.RFC3339), a.CreateOrder)

	// 쿠버네티스 노드 메트릭 처리 및 저장
	a.aggregateNodeMetric(metrics)

	// 쿠버네티스 파드 메트릭 처리 및 저장
	a.aggregatePodMetric(metrics, "kubernetes_pod_container")

	// 쿠버네티스 파드 네트워크 메트릭 처리 및 저장
	a.aggregatePodMetric(metrics, "kubernetes_pod_network")
}

// aggregateNodeMetric 쿠버네티스 노드 메트릭 처리 및 저장
func (a *Aggregator) aggregateNodeMetric(metrics []TelegrafMetric) {
	metricName := "kubernetes_node"

	// 1. 토픽 메세지에서 노드 메트릭 메세지를 필터링
	nodeMetricFilter := funk.Filter(metrics, func(metric TelegrafMetric) bool {
		return metric.Name == "kubernetes_node"
	})
	nodeMetricArr := nodeMetricFilter.([]TelegrafMetric)
	if len(nodeMetricArr) == 0 {
		return
	}

	// 2. 전체 클러스터에서 노드 목록 추출
	nodeNameFilter := funk.Uniq(funk.Get(nodeMetricArr, "Tags.node_name"))
	nodeNameArr := nodeNameFilter.([]string)

	// 개별 노드에 대한 모니터링 메트릭 처리 및 저장
	for _, nodeName := range nodeNameArr {
		currentNodeMetricArr := funk.Filter(nodeMetricArr, func(metric TelegrafMetric) bool {
			return metric.Tags["node_name"] == nodeName
		})
		nodeMetric := aggregateMetric(metricName, currentNodeMetricArr.([]TelegrafMetric), string(a.AggregateType))
		err := v1.GetInstance().WriteOnDemandMetric(v1.DefaultDatabase, nodeMetric.Name, nodeMetric.Tags, nodeMetric.Fields)
		if err != nil {
			util.GetLogger().Error(fmt.Sprintf("failed to write metric, error=%s", err.Error()))
			continue
		}
	}
}

// aggregateMetric 쿠버네티스 파드 메트릭 처리 및 저장
func (a *Aggregator) aggregatePodMetric(metrics []TelegrafMetric, metricName string) {
	// 데이터가 없을 경우
	if len(metrics) == 0 {
		return
	}

	// 1. 토픽 메세지에서 파드 메트릭 메세지를 필터링
	podMetricFilter := funk.Filter(metrics, func(metric TelegrafMetric) bool {
		return metric.Name == metricName
	})
	podMetricArr := podMetricFilter.([]TelegrafMetric)
	if len(podMetricArr) == 0 {
		return
	}

	// 2. 파드 별 메트릭 처리
	podNameFilter := funk.Uniq(funk.Get(podMetricArr, "Tags.pod_name"))
	podNameArr := podNameFilter.([]string)

	// 3. 단일 파드에 대한 파드 메트릭 처리 및 저장
	for _, podName := range podNameArr {
		currentPodMetricArr := funk.Filter(podMetricArr, func(metric TelegrafMetric) bool {
			return metric.Tags["pod_name"] == podName
		})
		podMetric := aggregateMetric(metricName, currentPodMetricArr.([]TelegrafMetric), string(a.AggregateType))
		err := v1.GetInstance().WriteOnDemandMetric(v1.DefaultDatabase, podMetric.Name, podMetric.Tags, podMetric.Fields)
		if err != nil {
			util.GetLogger().Error(fmt.Sprintf("failed to write metric, error=%s", err.Error()))
			continue
		}
	}
}

// aggregateMetric 쿠버네티스 메트릭 처리 및 저장
func aggregateMetric(metricType string, metrics []TelegrafMetric, criteria string) TelegrafMetric {
	aggregatedMetric := TelegrafMetric{
		Name:   metricType,
		Tags:   map[string]string{},
		Fields: map[string]interface{}{},
	}

	// 모니터링 메트릭 태그 정보 설정
	metricTags := make(map[string]string)
	for k, v := range metrics[0].Tags {
		metricTags[k] = v
	}
	aggregatedMetric.Tags = metricTags

	// last
	if criteria == "last" {
		// 가장 최신의 타임스탬프 값 가져오기
		timestampArr := funk.Get(metrics, "Timestamp")
		latestTimestamp := funk.MaxInt64(timestampArr.([]int64))
		// 최신 타임스탬프 값의 모니터링 메트릭 조회
		latestMetric := funk.Filter(metrics, func(metric TelegrafMetric) bool {
			return metric.Timestamp == latestTimestamp
		})
		return latestMetric.(TelegrafMetric)
	}

	// min, max, avg
	fieldList := funk.Keys(metrics[0].Fields)
	for _, fieldKey := range fieldList.([]string) {
		// 필드 별 데이터만 추출
		fieldDataArr := funk.Get(metrics, fmt.Sprintf("Fields.%s", fieldKey))
		// min, max, avg, last 처리
		var aggregatedVal float64
		if criteria == "min" {
			aggregatedVal = funk.MinFloat64(convertToFloat64Arr(fieldDataArr))
		} else if criteria == "max" {
			aggregatedVal = funk.MaxFloat64(convertToFloat64Arr(fieldDataArr))
		} else if criteria == "avg" {
			aggregatedVal = funk.SumFloat64(convertToFloat64Arr(fieldDataArr))
			aggregatedVal = aggregatedVal / float64(len(fieldDataArr.([]interface{})))
		}
		aggregatedMetric.Fields[fieldKey] = aggregatedVal
	}
	return aggregatedMetric
}

func convertToFloat64Arr(interfaceArr interface{}) []float64 {
	var float64Arr []float64
	for _, elem := range interfaceArr.([]interface{}) {
		float64Arr = append(float64Arr, elem.(float64))
	}
	return float64Arr
}
