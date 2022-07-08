package metric

import (
	"errors"

	"github.com/influxdata/influxdb1-client/models"
)

func MappingMonMetric(metricName string, metricVal *interface{}) (interface{}, error) {

	var mappingMetric models.Row
	var ok bool
	var err error

	if mappingMetric, ok = (*metricVal).(models.Row); !ok {
		return nil, errors.New("failed to mapping monitoring metric")
	}

	// Metric 구조체 Map 변환
	var metricKeyArr []string
	switch metricName {
	case "cpu":
		metricKeyArr = Cpu{}.GetField()
	case "cpufreq":
		metricKeyArr = Cpufreq{}.GetField()
	case "memory":
		metricKeyArr = Memory{}.GetField()
	case "disk":
		metricKeyArr = Disk{}.GetField()
	case "diskio":
		metricKeyArr = DiskIO{}.GetField()
	case "network":
		metricKeyArr = Network{}.GetField()
	case "net":
		metricKeyArr = Network{}.GetField()
	case "node":
		metricKeyArr = MCK8SNode{}.GetField()
	case "pod":
		metricKeyArr = MCK8SPod{}.GetField()
	default:
		err = errors.New("not found metric")
	}
	if err != nil {
		return nil, err
	}

	// 메트릭 정보 설정
	metricCols := make([]string, len(metricKeyArr)+1)
	for idx, metricKey := range metricKeyArr {
		if idx == 0 {
			metricCols[0] = "time"
		}
		metricCols[idx+1] = metricKey
	}

	mappingMetric.Columns = metricCols

	resultMap := map[string]interface{}{}
	resultMap["name"] = metricName
	resultMap["tags"] = mappingMetric.Tags
	resultMap["values"] = ConvertMetricValFormat(metricCols, mappingMetric.Values)

	return resultMap, nil
}

func ConvertMetricValFormat(metricKeyArr []string, metricVal [][]interface{}) []interface{} {
	convertedMetricVal := make([]interface{}, len(metricVal))
	for i, metricVal := range metricVal {
		newMetricVal := map[string]interface{}{}
		for j, key := range metricKeyArr {
			newMetricVal[key] = metricVal[j]
		}
		convertedMetricVal[i] = newMetricVal
	}
	return convertedMetricVal
}
