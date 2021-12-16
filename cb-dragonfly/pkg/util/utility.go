package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/cloud-barista/cb-dragonfly/pkg/types"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"math"
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
	return strings.Split(topicsStrings, "^")[1:]
}

func MergeTopicsToOneString(topicsSlice []string) string {
	var combinedTopicString string
	for _, topic := range topicsSlice {
		combinedTopicString = fmt.Sprintf("%s^%s", combinedTopicString, topic)
	}
	return combinedTopicString
}

func CalculateNumberOfCollector(topicCount int, maxHostCount int) int {
	collectorCount := topicCount / maxHostCount
	if topicCount%maxHostCount != 0 {
		collectorCount += 1
	}
	return collectorCount
}

func GetAllTopicBySort(topicsSlice []string) []string {
	if len(topicsSlice) == 0 {
		return []string{}
	}
	sort.Slice(topicsSlice, func(i, j int) bool {
		return topicsSlice[i] < topicsSlice[j]
	})
	return topicsSlice
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

func GetCspCollectorIdx(topic string) (collectorIdx int) {
	topicSplit := strings.Split(topic, "_")
	cspType := strings.ToUpper(topicSplit[len(topicSplit)-1])
	switch cspType {
	case types.Alibaba:
		collectorIdx = 0
		break
	case types.Aws:
		collectorIdx = 1
		break
	case types.Azure:
		collectorIdx = 2
		break
	case types.Cloudit:
		collectorIdx = 3
		break
	case types.Cloudtwin:
		collectorIdx = 4
		break
	case types.Docker:
		collectorIdx = 5
		break
	case types.Gcp:
		collectorIdx = 6
		break
	case types.Openstack:
		collectorIdx = 7
		break
	}
	return
}

func Unique(intSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range intSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	sort.Strings(list)
	return list
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

func Int32Ptr(i int32) *int32 { return &i }

func Int64Ptr(i int64) *int64 { return &i }

func DeploymentTemplate(collectorCreateOrder int, collectorUUID string, env []apiv1.EnvVar) *appsv1.Deployment {
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:   fmt.Sprintf("%s%d-%s", types.DeploymentName, collectorCreateOrder, collectorUUID),
			Labels: map[string]string{types.LabelKey: collectorUUID},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: Int32Ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					types.LabelKey: collectorUUID,
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						types.LabelKey: collectorUUID,
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  fmt.Sprintf("%s%d-%s", types.DeploymentName, collectorCreateOrder, collectorUUID),
							Image: types.CollectorImage,
							Ports: []apiv1.ContainerPort{},
							Env:   env,
							VolumeMounts: []apiv1.VolumeMount{
								{
									Name:      "config-volume",
									MountPath: "/go/src/github.com/cloud-barista/cb-dragonfly/conf",
								},
							},
							SecurityContext: &apiv1.SecurityContext{
								RunAsUser: Int64Ptr(0),
							},
						},
					},
					Volumes: []apiv1.Volume{
						{
							Name: "config-volume",
							VolumeSource: apiv1.VolumeSource{
								ConfigMap: &apiv1.ConfigMapVolumeSource{
									LocalObjectReference: apiv1.LocalObjectReference{
										Name: "cb-dragonfly-config",
									},
								},
							},
						},
					},
					ServiceAccountName: "cb-dragonfly",
				},
			},
		},
	}
}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func ToFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}

//func SliceContainsItem(slice []string, item string) bool {
//	set := make(map[string]struct{}, len(slice))
//	for _, s := range slice {
//		set[s] = struct{}{}
//	}
//
//	_, ok := set[item]
//	return ok
//}
