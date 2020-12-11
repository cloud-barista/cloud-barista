package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/cloud-barista/cb-dragonfly/pkg/types"
	"net/url"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

func StructToMap(i interface{}) (values url.Values) {
	values = url.Values{}
	iVal := reflect.ValueOf(i).Elem()
	typ := iVal.Type()
	for i := 0; i < iVal.NumField(); i++ {
		f := iVal.Field(i)
		// You ca use tags here...
		// tag := typ.Field(i).Tag.Get("tagname")
		// Convert each type into a string for the url.Values string map
		var v string
		switch f.Interface().(type) {
		case int, int8, int16, int32, int64:
			v = strconv.FormatInt(f.Int(), 10)
		case uint, uint8, uint16, uint32, uint64:
			v = strconv.FormatUint(f.Uint(), 10)
		case float32:
			v = strconv.FormatFloat(f.Float(), 'f', 4, 32)
		case float64:
			v = strconv.FormatFloat(f.Float(), 'f', 4, 64)
		case []byte:
			v = string(f.Bytes())
		case string:
			v = f.String()
		}
		values.Set(typ.Field(i).Name, v)
	}
	return
}

func ToMap(val interface{}) (map[string]interface{}, error) {

	// Convert struct to bytes
	bytes := new(bytes.Buffer)
	if err := json.NewEncoder(bytes).Encode(val); err != nil {
		return nil, err
	}

	// Convert bytes to map
	byteData := bytes.Bytes()
	resultMap := map[string]interface{}{}
	if err := json.Unmarshal(byteData, &resultMap); err != nil {
		return nil, err
	}

	return resultMap, nil
}

func GetFields(val reflect.Value) []string {
	var fieldArr []string
	t := val.Type()
	for i := 0; i < t.NumField(); i++ {
		fieldArr = append(fieldArr, t.Field(i).Tag.Get("json"))
	}
	return fieldArr
}

func SplitOneStringToTopicsSlice(topicsStrings string) []string {
	return strings.Split(topicsStrings, "&")[1:]
}

func MergeTopicsToOneString(topicsSlice []string) string {
	var combinedTopicString string
	for _, topic := range topicsSlice {
		combinedTopicString = fmt.Sprintf("%s&%s", combinedTopicString, topic)
	}
	return combinedTopicString
}

func CalculateNumberOfCollector(topicCount int, maxHostCount int) int {
	collectorCount := topicCount / maxHostCount
	if topicCount%maxHostCount != 0 || topicCount == 0 {
		collectorCount += 1
	}
	return collectorCount
}

func ReturnDiffTopicList(a, b []string) (diff []string) {
	m := make(map[string]bool)
	for _, item := range b {
		m[item] = true
	}
	for _, item := range a {
		if _, ok := m[item]; !ok {
			diff = append(diff, item)
		}
	}
	return
}

func GetAllTopicBySort(topicsSlice []string) []string {
	if len(topicsSlice) == 0 {
		return []string{}
	}
	sort.Slice(topicsSlice, func(i, j int) bool {
		return topicsSlice[i] < topicsSlice[j]
	})
	return topicsSlice[1:]
}

func MakeCollectorTopicMap(allTopics []string, maxHostCount int) (map[int][]string, []int) {

	if len(allTopics) == 0 {
		return map[int][]string{}, []int{}
	}

	collectorTopicMap := map[int][]string{}
	collectorTopicCnt := []int{}
	allTopicsLen := len(allTopics)
	startIdx := 0
	endIdx := 0

	collectorCount := CalculateNumberOfCollector(allTopicsLen, maxHostCount)

	for collectorCountIdx := 0; collectorCountIdx < collectorCount; collectorCountIdx++ {
		if allTopicsLen < maxHostCount {
			endIdx = len(allTopics)
		} else {
			endIdx = (collectorCountIdx + 1) * maxHostCount
		}
		aTopics := allTopics[startIdx:endIdx]
		collectorTopicMap[collectorCountIdx] = aTopics

		collectorTopicCnt = append(collectorTopicCnt, len(aTopics))

		startIdx = endIdx
		allTopicsLen -= maxHostCount
	}
	return collectorTopicMap, collectorTopicCnt
}

func MakeCollectorTopicMapBasedCSP(allTopics []string) map[int][]string {

	if len(allTopics) == 0 {
		return map[int][]string{}
	}

	collectorTopicMap := map[int][]string{
		0: []string{},
		1: []string{},
		2: []string{},
		3: []string{},
		4: []string{},
		5: []string{},
	}
	for _, topic := range allTopics {
		splitTopic := strings.Split(topic, "_")
		cspType := splitTopic[len(splitTopic)-1]
		switch cspType {
		case types.ALIBABA:
			collectorTopicMap[0] = append(collectorTopicMap[0], topic)
			break
		case types.AWS:
			collectorTopicMap[1] = append(collectorTopicMap[1], topic)
			break
		case types.AZURE:
			collectorTopicMap[2] = append(collectorTopicMap[2], topic)
			break
		case types.CLOUDIT:
			collectorTopicMap[3] = append(collectorTopicMap[3], topic)
			break
		case types.CLOUDTWIN:
			collectorTopicMap[4] = append(collectorTopicMap[4], topic)
			break
		case types.DOCKER:
			collectorTopicMap[5] = append(collectorTopicMap[5], topic)
			break
		case types.GCP:
			collectorTopicMap[6] = append(collectorTopicMap[6], topic)
			break
		case types.OPENSTACK:
			collectorTopicMap[7] = append(collectorTopicMap[7], topic)
			break
		}
	}
	return collectorTopicMap
}
