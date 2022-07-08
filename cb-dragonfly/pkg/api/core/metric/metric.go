package metric

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cloud-barista/cb-dragonfly/pkg/util"
	"net/http"
	"sort"
	"strings"

	influxdbmetric "github.com/cloud-barista/cb-dragonfly/pkg/storage/metricstore/influxdb/metric"
	v1 "github.com/cloud-barista/cb-dragonfly/pkg/storage/metricstore/influxdb/v1"

	"github.com/cloud-barista/cb-dragonfly/pkg/types"
	"github.com/influxdata/influxdb1-client/models"
)

// GetMonInfo 모니터링 메트릭 조회
func GetMonInfo(info types.DBMetricRequestInfo) (interface{}, int, error) {
	metric := types.Metric(info.MetricName)

	// 메트릭 타입 유효성 체크
	if metric == types.None {
		return nil, http.StatusInternalServerError, errors.New(fmt.Sprintf("not found metric : %s", info.MetricName))
	}

	switch metric {

	case types.Cpu, types.CpuFrequency, types.Memory, types.Network:
		if !util.CheckMCISType(info.ServiceType) {
			return nil, http.StatusBadRequest, errors.New(fmt.Sprintf("not supported metric data for %s, metric=%s", info.ServiceType, info.MetricName))
		}
		// cpu, cpufreq, memory, network 메트릭 조회
		cpuMetric, err := v1.GetInstance().ReadMetric(info)
		if err != nil {
			return nil, http.StatusInternalServerError, err
		}
		if cpuMetric == nil {
			return nil, http.StatusNotFound, errors.New(fmt.Sprintf("not found metric data, metric=%s", info.MetricName))
		}
		resultMetric, err := influxdbmetric.MappingMonMetric(metric.ToString(), &cpuMetric)
		if err != nil {
			return nil, http.StatusInternalServerError, err
		}
		return resultMetric, http.StatusOK, nil

	case types.Disk:
		if !util.CheckMCISType(info.ServiceType) {
			return nil, http.StatusBadRequest, errors.New(fmt.Sprintf("not supported metric data for %s, metric=%s", info.ServiceType, info.MetricName))
		}
		// disk, diskio 메트릭 조회
		diskMetric, err := v1.GetInstance().ReadMetric(info)
		if err != nil {
			return nil, http.StatusInternalServerError, err
		}
		diskIoMetric, err := v1.GetInstance().ReadMetric(info)
		if err != nil {
			return nil, http.StatusInternalServerError, err
		}
		if diskMetric == nil && diskIoMetric == nil {
			return nil, http.StatusNotFound, errors.New(fmt.Sprintf("not found metric data, metric=%s", info.MetricName))
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
		resultMap["name"] = info.MetricName
		resultMap["tags"] = resultRow.Tags
		resultMap["values"] = influxdbmetric.ConvertMetricValFormat(resultRow.Columns, resultRow.Values)
		return resultMap, http.StatusOK, nil

	case types.MCK8S_NODE:
		if !util.CheckMCK8SType(info.ServiceType) {
			return nil, http.StatusBadRequest, errors.New(fmt.Sprintf("not supported metric data for %s, metric=%s", info.ServiceType, info.MetricName))
		}

		info.MetricName = "kubernetes_node"
		nodeMetric, err := v1.GetInstance().ReadMetric(info)
		if err != nil {
			return nil, http.StatusInternalServerError, err
		}
		if nodeMetric == nil {
			return nil, http.StatusNotFound, errors.New(fmt.Sprintf("not found metric data, metric=%s", types.Node))
		}

		var resultData types.DBData
		byteData, err := json.Marshal(nodeMetric)
		if err != nil {
			return nil, http.StatusInternalServerError, errors.New("internal server error with parsing metric data")
		}
		if err = json.Unmarshal(byteData, &resultData); err != nil {
			return nil, http.StatusInternalServerError, errors.New("internal server error with parsing metric data")
		}

		resultData.Name = string(types.MCK8S_NODE)
		return resultData, http.StatusOK, nil

	case types.MCK8S_POD:
		if !util.CheckMCK8SType(info.ServiceType) {
			return nil, http.StatusBadRequest, errors.New(fmt.Sprintf("not supported metric data for %s, metric=%s", info.ServiceType, info.MetricName))
		}
		resultData := types.DBData{}
		mck8sMeasurements := []string{"kubernetes_pod_container", "kubernetes_pod_network"}

		for i, measurement := range mck8sMeasurements {
			info.MetricName = measurement
			podMetric, err := v1.GetInstance().ReadMetric(info)
			if err != nil {
				return nil, http.StatusInternalServerError, err
			}
			if podMetric == nil {
				return nil, http.StatusNotFound, errors.New(fmt.Sprintf("not found metric data, metric=%s", types.MCK8S_POD))
			}

			var tmpData types.DBData
			byteData, err := json.Marshal(podMetric)
			if err != nil {
				return nil, http.StatusInternalServerError, errors.New("internal server error with parsing metric data")
			}
			if err = json.Unmarshal(byteData, &tmpData); err != nil {
				return nil, http.StatusInternalServerError, errors.New("internal server error with parsing metric data")
			}

			for _, column := range tmpData.Columns {
				resultData.Columns = append(resultData.Columns, column)
			}

			if i == 0 {
				for _, tmp := range tmpData.Values {
					resultData.Values = append(resultData.Values, tmp)
				}
			} else {
				for _, tmp := range tmpData.Values {
					var copiedData []interface{}
					for _, tmpElement := range tmp {
						copiedData = append(copiedData, tmpElement)
					}
					for resultIndex, appendedValue := range resultData.Values {
						if strings.EqualFold(tmp[0].(string), appendedValue[0].(string)) {

							excludeTimeData := append(copiedData[:0], copiedData[1:]...)
							for _, extractedData := range excludeTimeData {
								resultData.Values[resultIndex] = append(resultData.Values[resultIndex], extractedData)
							}
						}
						continue
					}
				}
			}
			resultData.Tags = tmpData.Tags
		}

		resultData.Name = string(types.MCK8S_POD)
		resultData.Columns = util.Unique(resultData.Columns, false)
		return resultData, http.StatusOK, nil
	default:
		return nil, http.StatusInternalServerError, errors.New(fmt.Sprintf("NOT FOUND METRIC : %s", info.MetricName))
	}
}
