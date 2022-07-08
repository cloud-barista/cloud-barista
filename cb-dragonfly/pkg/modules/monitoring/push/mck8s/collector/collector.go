package collector

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/cloud-barista/cb-dragonfly/pkg/config"
	"github.com/cloud-barista/cb-dragonfly/pkg/types"
	"github.com/cloud-barista/cb-dragonfly/pkg/util"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type MetricCollector struct {
	KafkaAdminClient  *kafka.AdminClient
	KafkaConsumerConn *kafka.Consumer
	CreateOrder       int
	Aggregator        Aggregator
	Ch                chan string
}

func (mc *MetricCollector) DoCollect(wg *sync.WaitGroup) error {
	defer wg.Done()
	for {
		select {
		case chanData := <-mc.Ch:
			if len(chanData) != 0 {

				// 콜렉터 삭제 요청일 경우 토픽 구독 취소 및 삭제 처리
				if chanData == "close" {
					close(mc.Ch)

					mc.KafkaAdminClient.Close()
					_ = mc.KafkaConsumerConn.Unsubscribe()
					err := mc.KafkaConsumerConn.Close()
					if err != nil {
						errMsg := fmt.Sprintf("fail to delete mck8s metic collector, kafka close connection failed with error=%s", err)
						util.GetLogger().Error(errMsg)
						return errors.New(errMsg)
					}

					fmt.Printf("[%s] <MCK8S> DELETE Group_%d collector\n", time.Now().Format(time.RFC3339), mc.CreateOrder)
					return nil
				}

				deliveredTopic := chanData
				fmt.Printf("[%s] <MCK8S> EXECUTE Group_%d collector - topic: %s\n", time.Now().Format(time.RFC3339), mc.CreateOrder, deliveredTopic)

				// 토픽 데이터 구독
				err := mc.KafkaConsumerConn.SubscribeTopics([]string{deliveredTopic}, nil)
				if err != nil {
					errMsg := fmt.Sprintf("fail to subscribe topic with topic %s, error=%s", deliveredTopic, err)
					util.GetLogger().Error(errMsg)
					return errors.New(errMsg)
				}

				// 토픽 데이터 처리
				mc.Aggregator.AggregateMetric(mc.KafkaAdminClient, mc.KafkaConsumerConn, deliveredTopic)
			}
			break
		}
	}
}

func NewMetricCollector(aggregateType types.AggregateType, createOrder int) (MetricCollector, error) {

	kafkaConfig := &kafka.ConfigMap{
		"bootstrap.servers":  fmt.Sprintf("%s", config.GetDefaultConfig().Kafka.EndpointUrl),
		"group.id":           fmt.Sprintf("%d", createOrder),
		"enable.auto.commit": true,
		"auto.offset.reset":  "earliest",
	}

	// kafka 관리자 커넥션 설정
	kafkaAdminClient, err := kafka.NewAdminClient(kafkaConfig)
	if err != nil {
		errMsg := fmt.Sprintf("fail to create mcks metic collector, kafka admin connection failed with error=%s", err)
		util.GetLogger().Error(errMsg)
		return MetricCollector{}, errors.New(errMsg)
	}

	// kafka 컨슈머 커넥션 설정
	kafkaConsumerConn, err := kafka.NewConsumer(kafkaConfig)
	if err != nil {
		errMsg := fmt.Sprintf("fail to create mcks metic collector, kafka consumer connection failed with error=%s", err)
		util.GetLogger().Error(errMsg)
		return MetricCollector{}, errors.New(errMsg)
	}

	ch := make(chan string)
	mc := MetricCollector{
		KafkaAdminClient:  kafkaAdminClient,
		KafkaConsumerConn: kafkaConsumerConn,
		Aggregator: Aggregator{
			CreateOrder:   createOrder,
			AggregateType: aggregateType,
		},
		Ch:          ch,
		CreateOrder: createOrder,
	}
	fmt.Printf("[%s] <MCK8S> CREATE collector\n", time.Now().Format(time.RFC3339))
	return mc, nil
}
