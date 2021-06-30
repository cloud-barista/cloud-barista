package collector

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"sort"
	"strings"
	"sync"

	"github.com/cloud-barista/cb-dragonfly/pkg/config"
	"github.com/cloud-barista/cb-dragonfly/pkg/types"
	"github.com/cloud-barista/cb-dragonfly/pkg/util"

	"github.com/google/go-cmp/cmp"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

type MetricCollector struct {
	CreateOrder       int
	ConsumerKafkaConn *kafka.Consumer
	AdminKafkaConn    *kafka.AdminClient
	Aggregator        Aggregator
	Ch                chan string
}

type TelegrafMetric struct {
	Name      string                 `json:"name"`
	Tags      map[string]interface{} `json:"tags"`
	Fields    map[string]interface{} `json:"fields"`
	Timestamp int64                  `json:"timestamp"`
	TagInfo   map[string]interface{} `json:"tagInfo"`
}

var KafkaConfig *kafka.ConfigMap

func NewMetricCollector(aggregateType types.AggregateType, createOrder int) (MetricCollector, error) {

	KafkaConfig = &kafka.ConfigMap{
		"bootstrap.servers": fmt.Sprintf("%s", config.GetDefaultConfig().GetKafkaConfig().GetKafkaEndpointUrl()),
		//"bootstrap.servers":  "192.168.130.7",
		"group.id":           fmt.Sprintf("%d", createOrder),
		"enable.auto.commit": true,
		"auto.offset.reset":  "earliest",
	}

	consumerKafkaConn, err := kafka.NewConsumer(KafkaConfig)
	if err != nil {
		util.GetLogger().Error("Fail to create collector kafka consumer", err)
		util.GetLogger().Error(err)
		return MetricCollector{}, err
	}
	adminKafkaConn, err := kafka.NewAdminClient(KafkaConfig)
	if err != nil {
		util.GetLogger().Error("Fail to create collector kafka consumer", err)
		util.GetLogger().Error(err)
		return MetricCollector{}, err
	}
	ch := make(chan string)
	mc := MetricCollector{
		ConsumerKafkaConn: consumerKafkaConn,
		AdminKafkaConn:    adminKafkaConn,
		CreateOrder:       createOrder,
		Aggregator: Aggregator{
			AggregateType: aggregateType,
		},
		Ch: ch,
	}
	fmt.Println(fmt.Sprintf("#### Group_%d collector Create ####", createOrder))
	return mc, nil
}

const (
	DEFAULT                = "default"
	SUBSCRIBE_ALIVE_TOPICS = "subscribeAliveTopics"
	UNSUBSCRIBE            = "unsubscribe"
	CLOSE                  = "close"
	TOPIC_HEARTBEAT        = 6
)

func (mc *MetricCollector) Collector(wg *sync.WaitGroup) error {

	defer wg.Done()
	aliveTopics := []string{}
	isDeadTopics := []string{}
	topicsPolicy := DEFAULT
	topicsHeartBeat := TOPIC_HEARTBEAT
	for {
		select {
		case processDecision := <-mc.Ch:
			if processDecision == CLOSE {

				close(mc.Ch)
				err := mc.ConsumerKafkaConn.Close()
				if err != nil {
					logrus.Debug("Fail to collector kafka connection close")
				}
				mc.AdminKafkaConn.Close()
				fmt.Println(fmt.Sprintf("#### Group_%d collector Delete ####", mc.CreateOrder))
				return nil

			} else if len(processDecision) != 0 {

				DeliveredTopicList := unique(strings.Split(processDecision, "&")[1:])
				fmt.Println(fmt.Sprintf("Group_%d collector Delivered : %s", mc.CreateOrder, DeliveredTopicList))
				sort.Strings(aliveTopics)

				switch topicsPolicy {
				case SUBSCRIBE_ALIVE_TOPICS:
					_ = mc.ConsumerKafkaConn.SubscribeTopics(aliveTopics, nil)
					topicsPolicy = DEFAULT
					break
				case UNSUBSCRIBE:
					_ = mc.ConsumerKafkaConn.Unsubscribe()
					topicsPolicy = DEFAULT
					break
				default:
					_ = mc.ConsumerKafkaConn.SubscribeTopics(DeliveredTopicList, nil)
					break
				}
				aliveTopics, _ = mc.Aggregator.AggregateMetric(mc.ConsumerKafkaConn, DeliveredTopicList)
				if !cmp.Equal(DeliveredTopicList, aliveTopics) {
					if len(DeliveredTopicList) != 0 && len(aliveTopics) == 0 {
						topicsPolicy = UNSUBSCRIBE
					} else {
						topicsPolicy = SUBSCRIBE_ALIVE_TOPICS
					}
					diffTopics := ReturnDiffTopicList(DeliveredTopicList, aliveTopics)
					if cmp.Equal(isDeadTopics, diffTopics) {
						topicsHeartBeat--
					}
					isDeadTopics = diffTopics
					if topicsHeartBeat == 0 {
						isDeadTopics = []string{}
						topicsHeartBeat = TOPIC_HEARTBEAT
						topicsPolicy = DEFAULT
						break
					}
					_, _ = mc.AdminKafkaConn.DeleteTopics(context.Background(), diffTopics)
				}
			} else {
				isDeadTopics = []string{}
				topicsHeartBeat = TOPIC_HEARTBEAT
				topicsPolicy = DEFAULT
			}
			break
		}
	}
}
