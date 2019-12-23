package realtimestore

import (
	"errors"
)

func MappingMonMetric(metricName string, metricVal map[string]interface{}) (map[string]interface{}, error) {

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
