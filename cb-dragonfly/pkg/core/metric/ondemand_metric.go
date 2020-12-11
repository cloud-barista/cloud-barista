package metric

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/cloud-barista/cb-dragonfly/pkg/collector"
)

func GetVMOnDemandMonInfo(nsId string, mcisId string, vmId string, metricName string, publicIP string) (interface{}, int, error) {
	// disk, diskio 메트릭 조회
	if metricName == Disk {
		diskMetric, err := getVMOnDemandMonInfo(Disk, publicIP)
		if err != nil {
			return nil, http.StatusInternalServerError, err
		}
		diskioMetric, err := getVMOnDemandMonInfo(DiskIO, publicIP)
		if err != nil {
			return nil, http.StatusInternalServerError, err
		}

		diskMetricMap := diskMetric["values"].(map[string]interface{})
		diskioMetricMap := diskioMetric["values"].(map[string]interface{})
		for k, v := range diskioMetricMap {
			diskMetricMap[k] = v
		}

		return diskMetric, http.StatusOK, nil
	}

	var metricKey string
	switch metricName {
	case Cpu:
		metricKey = "cpu"
	case CpuFreqency:
		metricKey = "cpufreq"
	case Memory:
		metricKey = "mem"
	case Disk:
		metricKey = "disk"
	case Network:
		metricKey = "net"
	default:
		return nil, http.StatusInternalServerError, errors.New(fmt.Sprintf("not found metric : %s", metricName))
	}

	// cpu, cpufreq, memory, network 메트릭 조회
	resultMetric, err := getVMOnDemandMonInfo(metricKey, publicIP)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	return resultMetric, http.StatusOK, nil
}

func getVMOnDemandMonInfo(metricName string, publicIP string) (map[string]interface{}, error) {
	resp, err := http.Get(fmt.Sprintf("http://%s:8080/cb-dragonfly/metric/%s", publicIP, metricName))
	if err != nil {
		return nil, errors.New("agent server is closed")
	}
	defer resp.Body.Close()

	var metricData = map[string]collector.TelegrafMetric{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &metricData)
	if err != nil {
		return nil, err
	}
	resultMetric, err := convertMonMetric(metricName, metricData[metricName])
	if err != nil {
		return nil, err
	}
	return resultMetric, nil
}

func convertMonMetric(metricKey string, metricVal collector.TelegrafMetric) (map[string]interface{}, error) {
	metricMap := map[string]interface{}{}
	metricMap["name"] = metricVal.Name
	tagMap := map[string]interface{}{
		"nsId":   metricVal.Tags["nsId"],
		"mcisId": metricVal.Tags["mcisId"],
		"vmId":   metricVal.Tags["vmId"],
	}
	metricMap["tags"] = tagMap

	metricCols, err := mappingOnDemandMetric(metricKey, metricVal.Fields)
	if err != nil {
		return nil, err
	}
	metricMap["values"] = metricCols
	metricMap["time"] = time.Now().UTC()
	return metricMap, nil
}

func mappingOnDemandMetric(metricName string, metricVal map[string]interface{}) (map[string]interface{}, error) {
	// Metric 구조체 Map 변환
	metricCols := map[string]interface{}{}

	switch metricName {
	case "cpu":
		metricCols["cpu_utilization"] = metricVal["usage_utilization"]
		metricCols["cpu_system"] = metricVal["usage_system"]
		metricCols["cpu_idle"] = metricVal["usage_idle"]
		metricCols["cpu_iowait"] = metricVal["usage_iowait"]
		metricCols["cpu_hintr"] = metricVal["usage_irq"]
		metricCols["cpu_sintr"] = metricVal["usage_softirq"]
		metricCols["cpu_user"] = metricVal["usage_user"]
		metricCols["cpu_nice"] = metricVal["usage_nice"]
		metricCols["cpu_steal"] = metricVal["usage_steal"]
		metricCols["cpu_guest"] = metricVal["usage_guest"]
		metricCols["cpu_guest_nice"] = metricVal["usage_guest_nice"]
	case "cpufreq":
		metricCols["cpu_speed"] = metricVal["cur_freq"]
	case "mem":
		metricCols["mem_utilization"] = metricVal["used_percent"]
		metricCols["mem_total"] = metricVal["total"]
		metricCols["mem_used"] = metricVal["used"]
		metricCols["mem_free"] = metricVal["free"]
		metricCols["mem_shared"] = metricVal["shared"]
		metricCols["mem_buffers"] = metricVal["buffered"]
		metricCols["mem_cached"] = metricVal["cached"]
	case "disk":
		metricCols["disk_utilization"] = metricVal["used_percent"]
		metricCols["disk_total"] = metricVal["total"]
		metricCols["disk_used"] = metricVal["used"]
		metricCols["disk_free"] = metricVal["free"]
	case "diskio":
		metricCols["kb_read"] = metricVal["read_bytes"]
		metricCols["kb_written"] = metricVal["write_bytes"]
		metricCols["ops_read"] = metricVal["iops_read"]
		metricCols["ops_write"] = metricVal["iops_write"]
	case "net":
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
