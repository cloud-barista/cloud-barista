package util

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/Workiva/go-datastructures/queue"
)

type TopicStructure struct {
	Policy string
	Topic  string
}

// MCIS 큐
var ringQueueOnce sync.Once
var ringQueue *queue.Queue

func GetRingQueue() *queue.Queue {
	ringQueueOnce.Do(func() {
		ringQueue = queue.New(1000)
	})
	return ringQueue
}

// MCK8S 큐

var mck8sRingQueueOnce sync.Once
var mck8sRingQueue *queue.Queue

func GetMCK8SRingQueue() *queue.Queue {
	mck8sRingQueueOnce.Do(func() {
		mck8sRingQueue = queue.New(1000)
	})
	return mck8sRingQueue
}

func RingQueuePut(topicManagePolicy string, topic string) error {
	var topicBytes []byte
	var err error
	topicStructure := TopicStructure{
		Policy: topicManagePolicy,
		Topic:  topic,
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

func PutMCK8SRingQueue(topicManagePolicy string, topic string) error {
	var topicBytes []byte
	var err error
	topicStructure := TopicStructure{
		Policy: topicManagePolicy,
		Topic:  topic,
	}
	if topicBytes, err = json.Marshal(topicStructure); err != nil {
		fmt.Println("error?")
		return err
	}
	if err = GetMCK8SRingQueue().Put(topicBytes); err != nil {
		return err
	}
	return nil
}
