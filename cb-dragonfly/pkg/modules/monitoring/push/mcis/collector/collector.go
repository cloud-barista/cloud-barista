package collector

import (
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/cloud-barista/cb-dragonfly/pkg/config"
	"github.com/sirupsen/logrus"

	"github.com/cloud-barista/cb-dragonfly/pkg/types"
	"github.com/cloud-barista/cb-dragonfly/pkg/util"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/google/go-cmp/cmp"
)

type MetricCollector struct {
	CreateOrder       int
	ConsumerKafkaConn *kafka.Consumer
	Aggregator        Aggregator
	Ch                chan []string
}

var KafkaConfig *kafka.ConfigMap

// NewMetricCollector
//   - Go-routine 기반 collector 입니다.
//   - go channel 및 Kafka 기반으로 topic 분배받기 & topic 구독 및 가져오기를 수행합니다.
func NewMetricCollector(aggregateType types.AggregateType, createOrder int) (MetricCollector, error) {

	KafkaConfig = &kafka.ConfigMap{
		"bootstrap.servers":  fmt.Sprintf("%s", config.GetDefaultConfig().Kafka.EndpointUrl),
		"group.id":           fmt.Sprintf("%d", createOrder),
		"enable.auto.commit": true,
		//"max.poll.interval.ms": 300000,
		//"session.timeout.ms": 15000,
		//"max.poll.records": 1000,
		//"max.poll.interval": 6,
		"auto.offset.reset": "earliest",
	}

	consumerKafkaConn, err := kafka.NewConsumer(KafkaConfig)
	if err != nil {
		util.GetLogger().Error("Fail to create collector kafka consumer, Kafka Connection Fail", err)
		util.GetLogger().Error(err)
		return MetricCollector{}, err
	}
	ch := make(chan []string)
	mc := MetricCollector{
		ConsumerKafkaConn: consumerKafkaConn,
		CreateOrder:       createOrder,
		Aggregator: Aggregator{
			AggregateType: aggregateType,
		},
		Ch: ch,
	}
	fmt.Println(fmt.Sprintf("#### Group_%d collector Create ####", createOrder))
	return mc, nil
}

// Collector
//   - 콜렉터 매니저로부터 "close" 채널 값을 받으면 종료합니다. (고루틴 채널 중지 => 삭제)
//   - 콜렉터 매니저로부터 topic 리스트 값을 받으면 kafka 에 해당 topic 을 기준으로 데이터를 가져옵니다.
//   - kafka 에 요청한 topic 리스트 들 중 데이터가 3회 이상 넘어오지 않는 topic 의 경우 < 스케줄러가 활용하는 topic Queue > 에 삭제할 topic 으로 등록합니다.
func (mc *MetricCollector) Collector(wg *sync.WaitGroup) error {

	deadOrAliveCnt := map[string]int{}

	defer wg.Done()
	for {
		select {
		case processDecision := <-mc.Ch:
			if len(processDecision) != 0 {
				if processDecision[0] == "close" {
					close(mc.Ch)
					_ = mc.ConsumerKafkaConn.Unsubscribe()
					err := mc.ConsumerKafkaConn.Close()
					if err != nil {
						logrus.Debug("Fail to collector kafka connection close")
					}
					fmt.Println(fmt.Sprintf("#### Group_%d collector Delete ####", mc.CreateOrder))
					return nil
				}
				DeliveredTopicList := processDecision
				sort.Strings(DeliveredTopicList)
				fmt.Println(fmt.Sprintf("Group_%d collector Delivered : %s", mc.CreateOrder, DeliveredTopicList))

				// 분배받은 kafka topic 구독
				err := mc.ConsumerKafkaConn.SubscribeTopics(DeliveredTopicList, nil)
				if err != nil {
					fmt.Println(err)
				}
				start := time.Now()
				// 분배 받은 topic 들을 기준으로 데이터 수집, 가공, DB 저장
				// 이 후에 실제로 데이터를 가져와 위 로직을 수행한 topic 리스트 추출 => aliveTopics
				aliveTopics, _ := mc.Aggregator.AggregateMetric(mc.ConsumerKafkaConn, DeliveredTopicList)
				elapsed := time.Since(start)
				sort.Strings(aliveTopics)
				fmt.Println("Aggregate Time: ", elapsed)
				for _, aliveTopic := range aliveTopics {
					if _, ok := deadOrAliveCnt[aliveTopic]; ok {
						delete(deadOrAliveCnt, aliveTopic)
					}
				}
				// 분배받은 topic 리스트와 aliveTopics 를 비교하여 deadTopics 추출
				// 3회 이상 데이터를 주지 않는 topic의 경우 topic 삭제 queue 에 등록
				if !cmp.Equal(DeliveredTopicList, aliveTopics) {
					_ = mc.ConsumerKafkaConn.Unsubscribe()
					deadTopics := util.ReturnDiffTopicList(DeliveredTopicList, aliveTopics)
					for _, delTopic := range deadTopics {
						if _, ok := deadOrAliveCnt[delTopic]; !ok {
							deadOrAliveCnt[delTopic] = 0
						} else if ok {
							if deadOrAliveCnt[delTopic] == 2 {
								if err := util.RingQueuePut(types.TopicDel, delTopic); err != nil {
									logrus.Debug(fmt.Sprintf("failed to put topics to ring queue, error=%s", err))
								}
								delete(deadOrAliveCnt, delTopic)
							}
							deadOrAliveCnt[delTopic] += 1
						}
					}
				}
			}
			break
		}
	}
}
