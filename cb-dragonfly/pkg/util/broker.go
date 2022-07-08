package util

import (
	"fmt"
	"sync"

	"github.com/cloud-barista/cb-dragonfly/pkg/config"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type KafkaStruct struct {
	AdminClient    *kafka.AdminClient
	ConsumerClient *kafka.Consumer
}

var onceBroker sync.Once
var kafkaConn KafkaStruct

// TODO: Broker connection close logic
//defer func() {
//	if err := kafka.Close(); err != nil {
//		panic(err)
//	}
//}()
func Initialize() error {
	kafkaConfig := &kafka.ConfigMap{
		"bootstrap.servers":  fmt.Sprintf("%s", config.GetDefaultConfig().Kafka.EndpointUrl),
		"group.id":           "myGroup",
		"enable.auto.commit": true,
		"auto.offset.reset":  "earliest",
	}
	admin, err := kafka.NewAdminClient(kafkaConfig)
	consumer, err := kafka.NewConsumer(kafkaConfig)

	if err != nil {
		GetLogger().Error(err)
		GetLogger().Error("failed to  load-balancer kafka connection")
		return err
	} else {
		kafkaConn.AdminClient = admin
		kafkaConn.ConsumerClient = consumer
	}
	return nil
}

func GetInstance() (*KafkaStruct, error) {
	var err error
	onceBroker.Do(func() {
		err = Initialize()
	})
	if err != nil {
		return &KafkaStruct{}, err
	}
	return &kafkaConn, err
}

//func (k *KafkaStruct) GetAllTopics() []string {
//	topics := []string{}
//	getTopics, err := k.AdminClient.GetMetadata(nil, true, -1)
//	if err != nil {
//		GetLogger().Error(err)
//		GetLogger().Error("failed to  get all topics list")
//		return nil
//	} else {
//		for topic, _ := range getTopics.Topics {
//			topics = append(topics, topic)
//		}
//		return topics
//	}
//}

//func (k *KafkaStruct) DeleteTopics(topics []string) error {
//	_, err := k.AdminClient.DeleteTopics(context.Background(), topics)
//	if err != nil {
//		GetLogger().Error(err)
//		GetLogger().Error("failed to  delete topic list from broker")
//		return err
//	}
//	return nil
//}
