package influxdb

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

func MappingOnDemandMonMetric(metricName string, metricVal map[string]interface{}) (map[string]interface{}, error) {

	// Metric 구조체 Map 변환
	//var metricKeyArr []string
	metricCols := map[string]interface{}{}

	switch metricName {
	case "cpu":
		//metricKeyArr = metricstore.Cpu{}.GetField()
		metricCols["cpu_utilization"] = metricVal["usage_utilization"]
		metricCols["cpu_system"] = metricVal["usage_system"]
		metricCols["cpu_idle"] = metricVal["usage_idle"]
		metricCols["cpu_iowait"] = metricVal["usage_iowait"]
		metricCols["cpu_hintr"] = metricVal["usage_irq"]
		metricCols["cpu_sintr"] = metricVal["usage_softirq"]
	case "cpufreq":
		metricCols["cpu_speed"] = metricVal["cur_freq"]
	case "mem":
		//metricKeyArr = metricstore.Memory{}.GetField()
		metricCols["mem_utilization"] = metricVal["used_percent"]
		metricCols["mem_total"] = metricVal["total"]
		metricCols["mem_used"] = metricVal["used"]
		metricCols["mem_free"] = metricVal["free"]
		metricCols["mem_shared"] = metricVal["shared"]
		metricCols["mem_buffers"] = metricVal["buffered"]
		metricCols["mem_cached"] = metricVal["cached"]
	case "disk":
		//metricKeyArr = metricstore.Disk{}.GetField()
		metricCols["disk_utilization"] = metricVal["used_percent"]
		metricCols["disk_total"] = metricVal["total"]
		metricCols["disk_used"] = metricVal["used"]
		metricCols["disk_free"] = metricVal["free"]
	case "diskio":
		//metricKeyArr = metricstore.DiskIO{}.GetField()
		metricCols["kb_read"] = metricVal["read_bytes"]
		metricCols["kb_written"] = metricVal["write_bytes"]
		metricCols["ops_read"] = metricVal["iops_read"]
		metricCols["ops_write"] = metricVal["iops_write"]
	case "net":
		//metricKeyArr = metricstore.Network{}.GetField()
		metricCols["bytes_in"] = metricVal["bytes_recv"]
		metricCols["bytes_out"] = metricVal["bytes_sent"]
		metricCols["pkts_in"] = metricVal["packets_recv"]
		metricCols["pkts_out"] = metricVal["packets_sent"]
	default:
		err := errors.New("not found metric")
		return nil, err
	}

	return metricCols, nil
}
