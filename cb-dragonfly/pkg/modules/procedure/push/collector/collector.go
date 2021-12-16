package collector

import (
	"fmt"
	"github.com/cloud-barista/cb-dragonfly/pkg/config"
	"github.com/sirupsen/logrus"
	"sort"
	"sync"
	"time"

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

func NewMetricCollector(aggregateType types.AggregateType, createOrder int) (MetricCollector, error) {

	KafkaConfig = &kafka.ConfigMap{
		"bootstrap.servers":  fmt.Sprintf("%s", config.GetDefaultConfig().Kafka.EndpointUrl),
		"group.id":           fmt.Sprintf("%d", createOrder),
		"enable.auto.commit": true,
		//"max.poll.interval.ms": 300000,
		//"session.timeout.ms": 15000,
		//"max.poll.records": 1000,
		//"max.poll.interval": 6,
		"auto.offset.reset":  "earliest",
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

func (mc *MetricCollector) Collector(wg *sync.WaitGroup) error {

	deadOrAliveCnt := map[string] int{}

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

				err := mc.ConsumerKafkaConn.SubscribeTopics(DeliveredTopicList, nil)
				if err != nil {
					fmt.Println(err)
				}
				start := time.Now()
				aliveTopics, _ := mc.Aggregator.AggregateMetric(mc.ConsumerKafkaConn, DeliveredTopicList)
				elapsed := time.Since(start)
				sort.Strings(aliveTopics)
				fmt.Println("Aggregate Time: ", elapsed)
				for _, aliveTopic := range aliveTopics {
					if _, ok := deadOrAliveCnt[aliveTopic]; ok {
						delete(deadOrAliveCnt, aliveTopic)
					}
				}
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
