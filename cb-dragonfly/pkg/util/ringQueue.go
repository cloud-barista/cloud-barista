package util

import (
	"encoding/json"
	"fmt"
	que "github.com/Workiva/go-datastructures/queue"
	"sync"
)

type TopicStructure struct {
	Policy string
	Topic string
}

var ringQueueOnce sync.Once
var ringQueue *que.Queue

func GetRingQueue() *que.Queue {
	ringQueueOnce.Do(func() {
		ringQueue = que.New(100)
	})
	return ringQueue
}

func RingQueuePut(topicManagePolicy string, topic string) error {
	var topicBytes []byte
	var err error
	topicStructure := TopicStructure{
		Policy: topicManagePolicy,
		Topic: topic,
	}
	if topicBytes, err = json.Marshal(topicStructure); err != nil {
		fmt.Println("error?")
		return err
	}
	if err = GetRingQueue().Put(topicBytes); err != nil {
		return err
	}
	return nil
}
