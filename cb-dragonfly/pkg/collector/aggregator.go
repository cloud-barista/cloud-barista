package collector

import (
	"encoding/json"
	"fmt"
	"github.com/cloud-barista/cb-dragonfly/pkg/metricstore"
	"github.com/cloud-barista/cb-dragonfly/pkg/realtimestore"
	"github.com/sirupsen/logrus"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type AggregateType string

const (
	MIN  AggregateType = "min"
	MAX  AggregateType = "max"
	AVG  AggregateType = "avg"
	LAST AggregateType = "last"
)

func (a AggregateType) toString() string {
	switch a {
	case MIN:
		return "min"
	case MAX:
		return "max"
	case AVG:
		return "avg"
	case LAST:
		return "last"
	default:
		return ""
	}
}

type Aggregator struct {
	Etcd          realtimestore.Storage
	InfluxDB      metricstore.Storage
	AggregateType AggregateType
}

// 실시간 모니터링 데이터 Aggregate
func (a *Aggregator) AggregateMetric(collectorId string) error {

	/* Monitoring metric data tree : aggregatedMap
		Depth : vmId / parentMetricName / childMetricName / value
	 	= map[string] interface | map[string] interface | map[string] interface {}
		= map[string] interface { map[string] interface { map[string] interface {} } }
	*/

	/*1. Get VM List from ETCD */
	aggregatedMap := map[string]interface{}{}
	getVmList, err := a.Etcd.ReadMetric(fmt.Sprintf("/collector/%s/host", collectorId))

	if err != nil {
		if err.Error()[0:3] != "100" {
			logrus.Error("Failed to get vm list from ETCD : ", err)
			return err
		} else {
			logrus.Error("It is empty ETCD. Failed to get vm list from ETCD : ", err)
			return nil
		}
	} else if getVmList == nil {
		return nil
	}

	var vmList []string

	for _, vm := range getVmList.Nodes {

		vmId := strings.Split(vm.Key, "/")[4]
		vmList = append(vmList, vmId)
		vmIdNode, err := a.Etcd.ReadMetric(fmt.Sprintf("/host/%s/metric", vmId))

		if err != nil {
			logrus.Error("Failed to get vm metric list", err)
			return err
		}
		/* 2. Get metric List from ETCD */
		parentMetric := map[string]interface{}{}

		for _, metric := range vmIdNode.Nodes {

			parentMetricName := strings.Split(metric.Key, "/")[4]
			metricDataNode, err := a.Etcd.ReadMetric(metric.Key)

			if err != nil {
				logrus.Error("Failed to get vm metric list", err)
				return err
			}
			/*3. Aggregating metric data*/
			if hasGrandch, _ := regexp.MatchString("[a-zA-Z]", strings.Split(metricDataNode.Nodes[0].Key, "/")[5]); hasGrandch {
				// if metric name is disk or diskio which has grandchildren, if TRUE logic execution.
				childMetric := make(map[string]interface{})

				for _, chidMetricNode := range metricDataNode.Nodes {

					chidMetricData, err := a.Etcd.ReadMetric(chidMetricNode.Key)

					if err != nil {
						logrus.Error("Failed to get child metric list", err)
						return err
					}

					childMetricName := strings.Split(chidMetricData.Nodes[0].Key, "/")[5]
					timestamp := strings.Split(chidMetricData.Nodes[0].Key, "/")[6]

					grandChildata, err := a.Etcd.ReadMetric(chidMetricData.Nodes[0].Key)

					if err != nil {
						logrus.Error("Failed to get grand child metric list", err)
						return err
					}

					grandChildMetricString := make(map[string]interface{})
					err = json.Unmarshal([]byte(grandChildata.Value), &grandChildMetricString)

					if err != nil {
						logrus.Error("Failed to convert json string to map", err)
						return err
					}

					grandChildDetailedMetric := map[string]interface{}{timestamp: grandChildMetricString}
					grandchildMetric, err := a.CalculateMetric(childMetricName, grandChildDetailedMetric, a.AggregateType.toString(), !hasGrandch)

					if err != nil {
						logrus.Error("Failed to aggregate grandchildMetric data", err)
						return err
					}

					childMetric[childMetricName] = grandchildMetric
				}

				parentMetric[parentMetricName], err = a.CalculateMetric(parentMetricName, childMetric, a.AggregateType.toString(), hasGrandch)

				if err != nil {
					logrus.Error("Failed to aggregate parentMetric data", err)
					return err
				}

			} else { // cpu, mem, swap... etc (disk, diskio metric excludes)

				childDetailMetric := make(map[string]interface{})
				timestamp := strings.Split(metricDataNode.Nodes[0].Key, "/")[5]
				err = json.Unmarshal([]byte(metricDataNode.Nodes[0].Value), &childDetailMetric)

				if err != nil {
					logrus.Error("Failed to convert json string to map", err)
					return err
				}

				childMetricString := map[string]interface{}{timestamp: childDetailMetric}
				// metric aggregate start
				childMetric, err := a.CalculateMetric(parentMetricName, childMetricString, a.AggregateType.toString(), hasGrandch)

				if err != nil {
					logrus.Error("Failed to aggregate childMetric data", err)
					return err
				}
				// metric aggregate end
				parentMetric[parentMetricName] = childMetric
			}
		}
		// tagging data processing start
		metricTagMap := make(map[string]interface{})
		metricTagData, err := a.Etcd.ReadMetric(fmt.Sprintf("/host/%s/tag", vmId))

		if err != nil {
			logrus.Error("Failed to get tag data from ETCD", err)
			return err
		}

		// Add Vm tag(mcisId, hostId, osType) to metricMap
		err = json.Unmarshal([]byte(metricTagData.Value), &metricTagMap)
		if err != nil {
			logrus.Error("Failed to convert tag json string to map", err)
			return err
		}

		parentMetric["tag"] = metricTagMap
		// tagging data processing end

		aggregatedMap[vmId] = parentMetric
	}
	/* 4. 모니터링 데이터 저장 (InfluxDB) */

	err = a.InfluxDB.WriteMetric(aggregatedMap)
	if err != nil {
		return err
	}

	return nil
}

// 실시간 모니터링 데이터 통계 값 계산 (MIN, MAX, AVG, LAST)
func (a *Aggregator) CalculateMetric(metricName string, metric map[string]interface{}, aggregateType string, deterProgress bool) (interface{}, error) {

	metricArr := map[string][]float64{}
	resultMap := map[string]interface{}{}

	if deterProgress {
		goto diskProgress
	}
	/* 1. 실시간 모니터링 데이터 계산을 위한 맵 생성 */
	//if a.AggregateType == MIN || a.AggregateType == MAX || a.AggregateType == AVG {
	if aggregateType == "min" || aggregateType == "max" || aggregateType == "avg" {
		for key := range metric {
			for metric, val := range metric[key].(map[string]interface{}) {
				if metricArr[metric] == nil {
					metricArr[metric] = []float64{}
				}
				metricArr[metric] = append(metricArr[metric], val.(float64))
			}
		}
	}

	/* 2. 실시간 모니터링 데이터 통계 로직 적용 */
	switch aggregateType {
	// 조회된 실시간 데이터 목록 기준 MIN, MAX, AVG, LAST 값 계산
	//case MIN:
	case "min":
		// 조회된 실시간 데이터 목록 기준 MIN 값 계산
		for key := range metricArr {
			sort.Sort(sort.Float64Slice(metricArr[key]))
			resultMap[key] = metricArr[key][0]
		}
	//case MAX:
	case "max":
		// 조회된 실시간 데이터 목록 기준 MAX 값 계산
		for key := range metricArr {
			sort.Sort(sort.Reverse(sort.Float64Slice(metricArr[key])))
			resultMap[key] = metricArr[key][0]
		}
	//case AVG:
	case "avg":
		// 조회된 실시간 데이터 목록 기준 AVG 값 계산
		for key := range metricArr {
			var sum float64
			if metricArr[key] == nil {
				continue
			}
			for _, v := range metricArr[key] {
				sum += v
			}
			resultMap[key] = sum / float64(len(metricArr[key]))
		}
	//case LAST:
	case "last":
		// 마지막 타임스탬프 기준 실시간 모니터링 데이터 가져오기
		var timestampArr []int
		for key := range metric {
			timestamp, err := strconv.Atoi(key)
			if err != nil {
				logrus.Error("Failed to convert string to int")
			}
			timestampArr = append(timestampArr, timestamp)
		}

		sort.Sort(sort.Reverse(sort.IntSlice(timestampArr)))

		lastTimestamp := timestampArr[0]
		resultMap = metric[strconv.Itoa(lastTimestamp)].(map[string]interface{})
	}

	return resultMap, nil

diskProgress:

	deviceMap := metric

	for _, metricMap := range deviceMap {
		for key, val := range metricMap.(map[string]interface{}) {
			if key == "io_time" || key == "iops_in_progress" {
				continue
			}
			if resultMap[key] == nil {
				resultMap[key] = 0.0
			}
			resultMap[key] = resultMap[key].(float64) + val.(float64)
		}
	}

	// disk 메트릭의 경우 usage_utilization 메트릭 항목 재계산
	if metricName == "disk" {
		if resultMap["total"].(float64) == 0 {
			resultMap["used_percent"] = 0
		} else {
			//resultMap["used_percent"] = resultMap["used"].(float64) / resultMap["total"].(float64)
			deviceCnt := len(deviceMap)
			if deviceCnt == 0 {
				resultMap["used_percent"] = 0
			} else {
				resultMap["used_percent"] = resultMap["used_percent"].(float64) / float64(deviceCnt)
			}
		}
	}

	return resultMap, nil
}

// 실시간 모니터링 데이터 조회
func (a *Aggregator) GetAggregateMetric(vmId string, metricName string, aggregateType string) (map[string]interface{}, error) {

	metricDetailMap := make(map[string]interface{})

	// 모니터링 데이터 조회
	metricDataKey := fmt.Sprintf("/host/%s/metric/%s", vmId, metricName)
	metricDataNode, err := a.Etcd.ReadMetric(metricDataKey)
	if err != nil {
		logrus.Error("Failed to get metric", err)
		return nil, err
	}

	// 모니터링 데이터 파싱 (string to json)
	for _, data := range metricDataNode.Nodes {
		timestamp := strings.Split(data.Key, "/")[5]
		metricData := make(map[string]interface{})
		err := json.Unmarshal([]byte(data.Value), &metricData)
		if err != nil {
			logrus.Error("Failed to convert json string to map", err)
			return nil, err
		}
		metricDetailMap[timestamp] = metricData
	}

	// 모니터링 데이터 Aggregate
	aggregateMetric, err := a.CalculateMetric(metricName, metricDetailMap, aggregateType, false)
	if err != nil {
		logrus.Error("Failed to aggregate data", err)
		return nil, err
	}

	return aggregateMetric.(map[string]interface{}), nil
}

// 실시간 모니터링 데이터 조회 (disk, diskio)
func (a *Aggregator) GetAggregateDiskMetric(vmId string, metricName string, aggregateType string) (map[string]interface{}, error) {

	/* 1. device 정보 가져오기 */
	deviceDataKey := fmt.Sprintf("/host/%s/metric/%s", vmId, metricName)
	deviceNode, err := a.Etcd.ReadMetric(deviceDataKey)
	if err != nil {
		logrus.Error("Failed to get device list", err)
		return nil, err
	}

	/* 2. device 별 모니터링 메트릭 데이터 계산 */
	deviceMap := map[string]interface{}{}
	for _, device := range deviceNode.Nodes {

		// device 정보 가져오기
		deviceName := strings.Split(device.Key, "/")[5]

		// 모니터링 메트릭 정보 조회
		metricDataKey := fmt.Sprintf("/host/%s/metric/%s/%s", vmId, metricName, deviceName)
		metricDataNode, err := a.Etcd.ReadMetric(metricDataKey)
		if err != nil {
			logrus.Error("Failed to get disk, diskio metric", err)
			return nil, err
		}

		// 모니터링 데이터 파싱 (string to json)
		metricDetailMap := make(map[string]interface{})
		for _, data := range metricDataNode.Nodes {
			timestamp := strings.Split(data.Key, "/")[5]
			metricData := make(map[string]interface{})
			err := json.Unmarshal([]byte(data.Value), &metricData)
			if err != nil {
				logrus.Error("Failed to convert json string to map", err)
				return nil, err
			}
			metricDetailMap[timestamp] = metricData
		}

		deviceMap[deviceName], err = a.CalculateMetric(metricName, metricDetailMap, aggregateType, false)
		if err != nil {
			logrus.Error("Failed to calculate disk, diskio metric")
		}
	}

	/* 3. device 별 모니터링 메트릭 데이터 합산 */
	resultMap := map[string]interface{}{}
	for _, metricMap := range deviceMap {
		for key, val := range metricMap.(map[string]interface{}) {
			if key == "io_time" || key == "iops_in_progress" {
				continue
			}
			if resultMap[key] == nil {
				resultMap[key] = 0.0
			}
			resultMap[key] = resultMap[key].(float64) + val.(float64)
		}
	}

	// disk 메트릭의 경우 usage_utilization 메트릭 항목 재계산
	if metricName == "disk" {
		if resultMap["total"].(float64) == 0 {
			resultMap["used_percent"] = 0
		} else {
			deviceCnt := len(deviceMap)
			if deviceCnt == 0 {
				resultMap["used_percent"] = 0
			} else {
				resultMap["used_percent"] = resultMap["used_percent"].(float64) / float64(deviceCnt)
			}
		}
	}

	return resultMap, nil
}
