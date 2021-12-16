package metric

import (
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strings"

	"github.com/cloud-barista/cb-dragonfly/pkg/config"
	influxdbmetric "github.com/cloud-barista/cb-dragonfly/pkg/storage/metricstore/influxdb/metric"
	v1 "github.com/cloud-barista/cb-dragonfly/pkg/storage/metricstore/influxdb/v1"

	"github.com/cloud-barista/cb-dragonfly/pkg/types"
	"github.com/influxdata/influxdb1-client/models"
)

// CBMCISMetric 단일 MCIS Milkyway 메트릭
type CBMCISMetric struct {
	Result  string `json:"result"`
	Unit    string `json:"unit"`
	Desc    string `json:"desc"`
	Elapsed string `json:"elapsed"`
	SpecId  string `json:"specid"`
}

// MCBMCISMetric 멀티 MCIS Milkyway 메트릭
type MCBMCISMetric struct {
	ResultArray []CBMCISMetric `json:"resultarray"`
}

// Request GET Request 단일 Body 정보
type Request struct {
	Host string `json:"host"`
	Spec string `json:"spec"`
}

// Mrequest GET Request 멀티 Body 정보
type Mrequest struct {
	MultiHost []Request `json:"multihost"`
}

type Parameter struct {
	agent_ip    string
	mcis_metric string
}

// GetVMMonInfo 가상머신 모니터링 메트릭 조회
func GetVMMonInfo(nsId string, mcisId string, vmId string, metricName string, period string, aggregateType string, duration string) (interface{}, int, error) {
	metric := types.Metric(metricName)

	// 메트릭 타입 유효성 체크
	if metric == types.None {
		return nil, http.StatusInternalServerError, errors.New(fmt.Sprintf("not found metric : %s", metricName))
	}

	switch metric {

	case types.Cpu, types.CpuFrequency, types.Memory, types.Network:

		// cpu, cpufreq, memory, network 메트릭 조회
		cpuMetric, err := v1.GetInstance().ReadMetric(config.GetInstance().Monitoring.DefaultPolicy == types.PushPolicy, nsId, mcisId, vmId, metric.ToAgentMetricKey(), period, aggregateType, duration)
		if err != nil {
			return nil, http.StatusInternalServerError, err
		}
		if cpuMetric == nil {
			return nil, http.StatusNotFound, errors.New(fmt.Sprintf("not found metric data, metric=%s", metricName))
		}
		resultMetric, err := influxdbmetric.MappingMonMetric(metric.ToString(), &cpuMetric)
		if err != nil {
			return nil, http.StatusInternalServerError, err
		}
		return resultMetric, http.StatusOK, nil

	case types.Disk:

		// disk, diskio 메트릭 조회
		diskMetric, err := v1.GetInstance().ReadMetric(config.GetInstance().Monitoring.DefaultPolicy == types.PushPolicy, nsId, mcisId, vmId, types.Disk.ToString(), period, aggregateType, duration)
		if err != nil {
			return nil, http.StatusInternalServerError, err
		}
		diskIoMetric, err := v1.GetInstance().ReadMetric(config.GetInstance().Monitoring.DefaultPolicy == types.PushPolicy, nsId, mcisId, vmId, types.DiskIO.ToString(), period, aggregateType, duration)
		if err != nil {
			return nil, http.StatusInternalServerError, err
		}
		if diskMetric == nil && diskIoMetric == nil {
			return nil, http.StatusNotFound, errors.New(fmt.Sprintf("not found metric data, metric=%s", metricName))
		}

		diskRow := diskMetric.(models.Row)
		diskIoRow := diskIoMetric.(models.Row)

		// Aggregate Metric
		var resultRow models.Row
		resultRow.Name = types.Disk.ToString()
		resultRow.Tags = diskRow.Tags
		resultRow.Columns = append(resultRow.Columns, diskRow.Columns[0:]...)
		resultRow.Columns = append(resultRow.Columns, diskIoRow.Columns[1:]...)

		// TimePoint 맵 생성 (disk, diskio 메트릭)
		timePointMap := make(map[string]string, len(diskRow.Values))
		for _, val := range diskRow.Values {
			timePoint := val[0].(string)
			timePointMap[timePoint] = timePoint
		}
		for _, val := range diskIoRow.Values {
			timePoint := val[0].(string)
			if _, exist := timePointMap[timePoint]; !exist {
				timePointMap[timePoint] = timePoint
			}
		}

		// TimePoint 배열 생성
		idx := 0
		timePointArr := make([]string, len(timePointMap))
		for _, timePoint := range timePointMap {
			timePointArr[idx] = timePoint
			idx++
		}
		sort.Strings(timePointArr)

		// TimePoint 배열 기준 모니터링 메트릭 Aggregate
		for _, tp := range timePointArr {

			metricVal := make([]interface{}, 1)
			metricVal[0] = tp

			// disk 메트릭 aggregate
			diskMetricAdded := false
			for idx, val := range diskRow.Values {
				t := val[0].(string)
				if strings.EqualFold(t, tp) {
					metricVal = append(metricVal, val[1:]...)
					diskMetricAdded = true
					break
				}
				// 해당 TimePoint에 해당하는 disk 메트릭이 없을 경우 0으로 값 초기화
				if !diskMetricAdded && (idx == len(diskRow.Values)-1) {
					initVal := make([]interface{}, len(val)-1)
					for i := range initVal {
						initVal[i] = 0
					}
					metricVal = append(metricVal, initVal...)
				}
			}

			// diskio 메트릭 aggregate
			diskIoMetricAdded := false
			for idx, val := range diskIoRow.Values {
				t := val[0].(string)
				if strings.EqualFold(t, tp) {
					metricVal = append(metricVal, val[1:]...)
					diskIoMetricAdded = true
					break
				}
				// 해당 TimePoint에 해당하는 disk 메트릭이 없을 경우 0으로 값 초기화
				if !diskIoMetricAdded && (idx == len(diskIoRow.Values)-1) {
					initVal := make([]interface{}, len(val)-1)
					for i := range initVal {
						initVal[i] = 0
					}
					metricVal = append(metricVal, initVal...)
				}
			}

			resultRow.Values = append(resultRow.Values, metricVal)
		}

		resultMap := map[string]interface{}{}
		resultMap["name"] = metricName
		resultMap["tags"] = resultRow.Tags
		resultMap["values"] = influxdbmetric.ConvertMetricValFormat(resultRow.Columns, resultRow.Values)
		return resultMap, http.StatusOK, nil

	default:
		return nil, http.StatusInternalServerError, errors.New(fmt.Sprintf("NOT FOUND METRIC : %s", metricName))
	}

	return nil, http.StatusInternalServerError, errors.New(fmt.Sprintf("not found metric : %s", metricName))
}
