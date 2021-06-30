package metric

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/cloud-barista/cb-dragonfly/pkg/collector"
	"github.com/cloud-barista/cb-dragonfly/pkg/types"
)

const (
	AgentPort    = 8888
	AgentTimeout = 10
)

func GetVMOnDemandMonInfo(metricName string, publicIP string) (interface{}, int, error) {
	metric := types.Metric(metricName)

	// 메트릭 타입 유효성 체크
	if metric == types.None {
		return nil, http.StatusInternalServerError, errors.New(fmt.Sprintf("not found metric : %s", metricName))
	}

	// disk, diskio 메트릭 조회
	if metric == types.Disk {
		diskMetric, err := getVMOnDemandMonInfo(types.Disk, publicIP)
		if err != nil {
			return nil, http.StatusInternalServerError, err
		}
		diskioMetric, err := getVMOnDemandMonInfo(types.DiskIO, publicIP)
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

	// cpu, cpufreq, memory, network 메트릭 조회
	resultMetric, err := getVMOnDemandMonInfo(metric, publicIP)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	return resultMetric, http.StatusOK, nil
}

func getVMOnDemandMonInfo(metric types.Metric, publicIP string) (map[string]interface{}, error) {
	client := http.Client{
		Timeout: AgentTimeout * time.Second,
	}
	agentUrl := fmt.Sprintf("http://%s:%d/cb-dragonfly/metric/%s", publicIP, AgentPort, metric.ToAgentMetricKey())
	resp, err := client.Get(agentUrl)
	if err != nil {
		return nil, err
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
	resultMetric, err := convertMonMetric(metric, metricData[metric.ToAgentMetricKey()])
	if err != nil {
		return nil, err
	}
	return resultMetric, nil
}

func convertMonMetric(metric types.Metric, metricVal collector.TelegrafMetric) (map[string]interface{}, error) {
	metricMap := map[string]interface{}{}
	metricMap["name"] = metricVal.Name
	tagMap := map[string]interface{}{
		"nsId":    metricVal.Tags["nsId"],
		"mcisId":  metricVal.Tags["mcisId"],
		"vmId":    metricVal.Tags["vmId"],
		"osType":  metricVal.Tags["osType"],
		"cspType": metricVal.Tags["cspType"],
	}
	metricMap["tags"] = tagMap

	metricCols, err := mappingOnDemandMetric(metric, metricVal.Fields)
	if err != nil {
		return nil, err
	}
	metricMap["values"] = metricCols
	metricMap["time"] = time.Now().UTC()
	return metricMap, nil
}

func mappingOnDemandMetric(metric types.Metric, metricVal map[string]interface{}) (map[string]interface{}, error) {
	// Metric 구조체 Map 변환
	metricCols := map[string]interface{}{}

	switch metric {
	case types.Cpu:
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
	case types.CpuFrequency:
		metricCols["cpu_speed"] = metricVal["cur_freq"]
	case types.Memory:
		metricCols["mem_utilization"] = metricVal["used_percent"]
		metricCols["mem_total"] = metricVal["total"]
		metricCols["mem_used"] = metricVal["used"]
		metricCols["mem_free"] = metricVal["free"]
		metricCols["mem_shared"] = metricVal["shared"]
		metricCols["mem_buffers"] = metricVal["buffered"]
		metricCols["mem_cached"] = metricVal["cached"]
	case types.Disk:
		metricCols["disk_utilization"] = metricVal["used_percent"]
		metricCols["disk_total"] = metricVal["total"]
		metricCols["disk_used"] = metricVal["used"]
		metricCols["disk_free"] = metricVal["free"]
	case types.DiskIO:
		metricCols["kb_read"] = metricVal["read_bytes"]
		metricCols["kb_written"] = metricVal["write_bytes"]
		metricCols["ops_read"] = metricVal["iops_read"]
		metricCols["ops_write"] = metricVal["iops_write"]
	case types.Network:
		metricCols["bytes_in"] = metricVal["bytes_recv"]
		metricCols["bytes_out"] = metricVal["bytes_sent"]
		metricCols["pkts_in"] = metricVal["packets_recv"]
		metricCols["pkts_out"] = metricVal["packets_sent"]
	case types.None:
	default:
		err := errors.New("not found metric")
		return nil, err
	}

	return metricCols, nil
}
