package kafka

import (
	"context"
	"fmt"
	"sync"

	"github.com/cloud-barista/cb-dragonfly/pkg/config"
	"github.com/cloud-barista/cb-dragonfly/pkg/util"

	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

type KafkaStruct struct {
	AdminClient *kafka.AdminClient
	L           *sync.RWMutex
}

var once sync.Once
var adminKafka KafkaStruct

// TODO: Broker connection close logic
//defer func() {
//	if err := kafka.Close(); err != nil {
//		panic(err)
//	}
//}()
func Initialize() error {
	admin, err := kafka.NewAdminClient(&kafka.ConfigMap{
		"bootstrap.servers": fmt.Sprintf("%s", config.GetDefaultConfig().GetKafkaConfig().GetKafkaEndpointUrl()),
		//"bootstrap.servers":  "192.168.130.7",
		"group.id":          "myGroup",
		"auto.offset.reset": "earliest",
	})

	if err != nil {
		util.GetLogger().Error(err)
		util.GetLogger().Error("failed to  load-balancer kafka connection")
		return err
	} else {
		adminKafka.AdminClient = admin
	}
	return nil
}

func GetInstance() (*KafkaStruct, error) {
	var err error
	once.Do(func() {
		err = Initialize()
	})
	if err != nil {
		return &KafkaStruct{}, err
	}
	return &adminKafka, err
}

func (k *KafkaStruct) GetAllTopics() []string {
	topics := []string{}
	getTopics, err := k.AdminClient.GetMetadata(nil, true, -1)
	if err != nil {
		util.GetLogger().Error(err)
		util.GetLogger().Error("failed to  get all topics list")
		return nil
	} else {
		for topic, _ := range getTopics.Topics {
			topics = append(topics, topic)
		}
		return topics
	}
}

func (k *KafkaStruct) DeleteTopics(topics []string) error {
	_, err := k.AdminClient.DeleteTopics(context.Background(), topics)
	if err != nil {
		util.GetLogger().Error(err)
		util.GetLogger().Error("failed to  delete topic list from broker")
		return err
	}
	return nil
}
