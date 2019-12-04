package metricstore

import (
	"errors"
	"github.com/influxdata/influxdb1-client/models"
)

func MappingMonMetric(metricName string, metricVal *interface{}) (*interface{}, error) {

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
	case "mem":
		metricKeyArr = Memory{}.GetField()
	case "disk":
		metricKeyArr = Disk{}.GetField()
	case "diskio":
		metricKeyArr = DiskIO{}.GetField()
	case "net":
		metricKeyArr = Network{}.GetField()
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

	var resultMap interface{}
	resultMap = mappingMetric
	return &resultMap, nil
}
