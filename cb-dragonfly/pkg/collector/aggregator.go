package collector

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/cloud-barista/cb-dragonfly/pkg/cbstore"
	"github.com/cloud-barista/cb-dragonfly/pkg/metricstore/influxdb/v1"
	"github.com/cloud-barista/cb-dragonfly/pkg/util"
	"github.com/sirupsen/logrus"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"

	"github.com/cloud-barista/cb-dragonfly/pkg/types"
)

type Aggregator struct {
	AggregateType types.AggregateType
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
			logrus.Debug("From AggregateMetric, pre-topics conn based kafkaConn bring about above err : ", err)
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
		if stayConnCount == types.ReadConnectionTimeout {
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
					if key == types.NsId || key == types.McisId || key == types.VmId || key == types.OsType || key == types.CspType {
						tagInfo[vmTopic][key] = tag.(string)
					}
				}
			} else {
				tagInfo[vmTopic] = make(map[string]string)
				for key, tag := range response.Tags {
					if key == types.NsId || key == types.McisId || key == types.VmId || key == types.OsType || key == types.CspType {
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
		result, err := a.CalculateMetric(uniqueResponseSlice, tagInfo, a.AggregateType.ToString())
		if err != nil {
			util.GetLogger().Error(err)
		}
		err = v1.GetInstance().WriteMetric(v1.DefaultDatabase, result)
		if err != nil {
			return []string{}, err
		}
	}

	currentTopics := unique(msgTopic)
	for _, topic := range currentTopics {
		cbstore.GetInstance().StoreDelete(types.DELTOPICS + topic)
	}
	delTopics := []string{}
	needCheckTopics := ReturnDiffTopicList(topics, currentTopics)
	if len(needCheckTopics) != 0 {
		for _, topic := range needCheckTopics {
			if cbstore.GetInstance().StoreGet(types.DELTOPICS+topic) == "" {
				cbstore.GetInstance().StorePut(types.DELTOPICS+topic, "0")
			} else {
				count, _ := strconv.Atoi(cbstore.GetInstance().StoreGet(types.DELTOPICS + topic))
				count++
				cbstore.GetInstance().StorePut(types.DELTOPICS+topic, strconv.Itoa(count))
			}
			checkNum, _ := strconv.Atoi(cbstore.GetInstance().StoreGet(types.DELTOPICS + topic))
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
				switch types.AggregateType(aggregateType) {
				case types.MIN:
					sort.Sort(sort.Float64Slice(slice))
					metric[key] = slice[0]
				case types.MAX:
					sort.Sort(sort.Reverse(sort.Float64Slice(slice)))
					metric[key] = slice[0]
				case types.AVG:
					var sum float64
					for _, v := range slice {
						sum += v
					}
					metric[key] = sum / float64(len(slice))
				case types.LAST:
					metric[key] = slice[len(slice)-1]
				}
				resultMap[vmTopic].(map[string]interface{})[metricName] = metric
			}
			resultMap[vmTopic].(map[string]interface{})["tagInfo"] = tagMap[vmTopic]
		}
	}
	return resultMap, nil
}
