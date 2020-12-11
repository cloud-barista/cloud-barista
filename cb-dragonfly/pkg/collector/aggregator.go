package collector

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"

	"github.com/cloud-barista/cb-dragonfly/pkg/metricstore/influxdb/influxdbv1"
	"github.com/cloud-barista/cb-dragonfly/pkg/types"
)

type Aggregator struct {
	AggregateType AggregateType
}

func (a *Aggregator) AggregateMetric(kafkaConn *kafka.Consumer, topics []string) ([]string, error) {

	currentTime := time.Now().Unix()
	stayConnCount := 0

	var msgSlice [][]byte
	var msgTopic []string
	for {
		stayConnCount += 1
		msg, err := kafkaConn.ReadMessage(1 * time.Second)
		if err != nil {
			logrus.Debug(err)
		}
		if msg != nil {
			msgTime := msg.Timestamp.Unix()
			msgSlice = append(msgSlice, msg.Value)
			msgTopic = append(msgTopic, *msg.TopicPartition.Topic)
			if msgTime > currentTime {
				break
			}
			stayConnCount = 0
			msg = nil
		}
		if stayConnCount == READ_CONNECTION_TIMEOUT {
			break
		}
	}
	fmt.Println(fmt.Sprintf("%v : %d\n", topics, len(msgSlice)))

	tagInfo := make(map[string]map[string]string)
	if len(msgSlice) != 0 {
		uniqueResponseSlice := make(map[string]map[string]map[string][]float64)
		for idx, value := range msgSlice {
			response := TelegrafMetric{}
			_ = json.Unmarshal(value, &response)

			vmTopic := msgTopic[idx]
			if _, ok := tagInfo[vmTopic]; ok {
				for key, tag := range response.Tags {
					if key == types.NSID || key == types.MCISID || key == types.VMID || key == types.OSTYPE || key == types.CSPTYPE {
						tagInfo[vmTopic][key] = tag.(string)
					}
				}
			} else {
				tagInfo[vmTopic] = make(map[string]string)
				for key, tag := range response.Tags {
					if key == types.NSID || key == types.MCISID || key == types.VMID || key == types.OSTYPE || key == types.CSPTYPE {
						tagInfo[vmTopic][key] = tag.(string)
					}
				}
			}

			if _, ok := uniqueResponseSlice[vmTopic]; ok {
				if _, ok := uniqueResponseSlice[vmTopic][response.Name]; ok {
					for metricName, val := range response.Fields {
						uniqueResponseSlice[vmTopic][response.Name][metricName] = append(uniqueResponseSlice[vmTopic][response.Name][metricName], val.(float64))
					}
				} else {
					uniqueResponseSlice[vmTopic][response.Name] = make(map[string][]float64)
					for metricName, val := range response.Fields {
						uniqueResponseSlice[vmTopic][response.Name][metricName] = append(uniqueResponseSlice[vmTopic][response.Name][metricName], val.(float64))
					}
				}
			} else {
				uniqueResponseSlice[vmTopic] = make(map[string]map[string][]float64)
				uniqueResponseSlice[vmTopic][response.Name] = make(map[string][]float64)
				for metricName, val := range response.Fields {
					uniqueResponseSlice[vmTopic][response.Name][metricName] = append(uniqueResponseSlice[vmTopic][response.Name][metricName], val.(float64))
				}
			}
		}
		result, err := a.CalculateMetric(uniqueResponseSlice, tagInfo, a.AggregateType.toString())
		if err != nil {
			logrus.Debug(err)
		}
		err = influxdbv1.GetInstance().WriteMetric(result)
		if err != nil {
			return []string{}, err
		}
	}

	currentTopics := unique(msgTopic)
	for _, topic := range currentTopics {
		GetInstance().StoreDelete(DELTOPICS + topic)
	}
	delTopics := []string{}
	needCheckTopics := ReturnDiffTopicList(topics, currentTopics)
	if len(needCheckTopics) != 0 {
		for _, topic := range needCheckTopics {
			if GetInstance().StoreGet(DELTOPICS+topic) == "" {
				GetInstance().StorePut(DELTOPICS+topic, "0")
			} else {
				count, _ := strconv.Atoi(GetInstance().StoreGet(DELTOPICS + topic))
				count++
				GetInstance().StorePut(DELTOPICS+topic, strconv.Itoa(count))
			}
			checkNum, _ := strconv.Atoi(GetInstance().StoreGet(DELTOPICS + topic))
			if checkNum >= 2 {
				delTopics = append(delTopics, topic)
			}
		}
	}

	return ReturnDiffTopicList(topics, delTopics), nil
}

func (a *Aggregator) CalculateMetric(responseMap map[string]map[string]map[string][]float64, tagMap map[string]map[string]string, aggregateType string) (map[string]interface{}, error) {

	resultMap := map[string]interface{}{}

	for vmTopic, metric := range responseMap {
		resultMap[vmTopic] = make(map[string]interface{})
		for metricName, metricSlice := range metric {
			metric := map[string]interface{}{}
			for key, slice := range metricSlice {
				switch aggregateType {
				case MINIMUM:
					sort.Sort(sort.Float64Slice(slice))
					metric[key] = slice[0]
				case MAXIMUM:
					sort.Sort(sort.Reverse(sort.Float64Slice(slice)))
					metric[key] = slice[0]
				case AVERAGE:
					var sum float64
					for _, v := range slice {
						sum += v
					}
					metric[key] = sum / float64(len(slice))
				case LATEST:
					metric[key] = slice[len(slice)-1]
				}
				resultMap[vmTopic].(map[string]interface{})[metricName] = metric
			}
			resultMap[vmTopic].(map[string]interface{})["tagInfo"] = tagMap[vmTopic]
		}
	}
	return resultMap, nil
}
