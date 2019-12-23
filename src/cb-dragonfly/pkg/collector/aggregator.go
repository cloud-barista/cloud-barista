package collector

import (
	"encoding/json"
	"fmt"
	"github.com/cloud-barista/cb-dragonfly/pkg/metricstore"
	"github.com/cloud-barista/cb-dragonfly/pkg/realtimestore"
	"github.com/davecgh/go-spew/spew"
	"github.com/sirupsen/logrus"
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
	var vmList []string     // 콜렉터에 태그된 VM 목록
	var metricList []string // 모니터링 메트릭 목록

	/* 1. 콜렉터에 태그된 VM 목록, 모니터링 메트릭 목록 가져오기 */

	// VM 목록 가져오기
	node, err := a.Etcd.ReadMetric(fmt.Sprintf("/collector/%s/host", collectorId))
	if err != nil {
		logrus.Error("Failed to get tagging vm list", err)
		return err
	}

	if node == nil {
		return nil
	}

	for idx, vm := range node.Nodes {
		vmId := strings.Split(vm.Key, "/")[4]
		vmList = append(vmList, vmId)

		// 메트릭 목록 가져오기
		if idx == 0 {
			metricNode, err := a.Etcd.ReadMetric(fmt.Sprintf("/host/%s/metric", vmId))
			if err != nil {
				logrus.Error("Failed to get vm metric list", err)
				return err
			}
			for _, metric := range metricNode.Nodes {
				metricName := strings.Split(metric.Key, "/")[4]
				metricList = append(metricList, metricName)
			}
		}
	}

	/* 2. 상세 모니터링 데이터 목록 가져오기 */
	monMap := map[string]map[string]map[string]interface{}{}

	for _, vmId := range vmList {
		metricMap := map[string]map[string]interface{}{}
		for _, metricName := range metricList {

			if metricName == "disk" || metricName == "diskio" {
				err = a.AggregateDiskMetric(vmId, metricName)
				if err != nil {
					logrus.Error("Failed to aggregate disk metric", err)
				}
				continue
			}

			// 모니터링 데이터 조회
			metricDataKey := fmt.Sprintf("/host/%s/metric/%s", vmId, metricName)
			metricDataNode, err := a.Etcd.ReadMetric(metricDataKey)
			if err != nil {
				return err
			}

			// 모니터링 데이터 파싱 (string to json)
			metricDetailMap := make(map[string]interface{})
			for _, data := range metricDataNode.Nodes {
				timestamp := strings.Split(data.Key, "/")[5]
				metricData := make(map[string]interface{})
				err := json.Unmarshal([]byte(data.Value), &metricData)
				if err != nil {
					logrus.Error("Failed to convert json string to map", err)
					return err
				}
				metricDetailMap[timestamp] = metricData
			}

			metricMap[metricName] = metricDetailMap
		}
		monMap[vmId] = metricMap
	}

	/* 3. 모니터링 데이터 Aggregate */
	aggregateMap := map[string]interface{}{}

	for hostId, host := range monMap {
		metricMap := map[string]interface{}{}
		for metricName, metric := range host {
			aggregateMetric, err := a.CalculateMetric(metricName, metric, a.AggregateType.toString())
			if err != nil {
				logrus.Error("Failed to aggregate data", err)
				return err
			}
			metricMap[metricName] = aggregateMetric
		}
		aggregateMap[hostId] = metricMap
	}

	/* 4. 모니터링 데이터 저장 (InfluxDB) */
	err = a.InfluxDB.WriteMetric(aggregateMap)
	if err != nil {
		return err
	}

	/* 5. 모니터링 데이터 초기화 (etcd) */
	err = a.FlushMetric(vmList)
	if err != nil {
		return err
	}

	return nil
}

func (a *Aggregator) AggregateDiskMetric(vmId string, metricName string) error {

	/* 1. device 정보 가져오기 */
	deviceDataKey := fmt.Sprintf("/host/%s/metric/%s", vmId, metricName)
	deviceNode, err := a.Etcd.ReadMetric(deviceDataKey)
	if err != nil {
		logrus.Error("Failed to get device list", err)
		return err
	}

	/* 2. device 별 모니터링 메트릭 데이터 계산 */
	deviceMap := map[string]interface{}{}
	for _, device := range deviceNode.Nodes {

		// device 정보 가져오기
		deviceName := strings.Split(device.Key, "/")[5]
		/*if strings.Contains(deviceName, "%") {
			deviceName = strings.ReplaceAll(deviceName, "%", "/")
		}*/

		// 모니터링 메트릭 정보 조회
		metricDataKey := fmt.Sprintf("/host/%s/metric/%s/%s", vmId, metricName, deviceName)
		metricDataNode, err := a.Etcd.ReadMetric(metricDataKey)
		if err != nil {
			logrus.Error("Failed to get disk, diskio metric", err)
			return err
		}

		// 모니터링 데이터 파싱 (string to json)
		metricDetailMap := make(map[string]interface{})
		for _, data := range metricDataNode.Nodes {
			timestamp := strings.Split(data.Key, "/")[5]
			metricData := make(map[string]interface{})
			err := json.Unmarshal([]byte(data.Value), &metricData)
			if err != nil {
				logrus.Error("Failed to convert json string to map", err)
				return err
			}
			metricDetailMap[timestamp] = metricData
		}

		deviceMap[deviceName], err = a.CalculateMetric(metricName, metricDetailMap, a.AggregateType.toString())
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
			//resultMap["used_percent"] = resultMap["used"].(float64) / resultMap["total"].(float64)
			deviceCnt := len(deviceMap)
			if deviceCnt == 0 {
				resultMap["used_percent"] = 0
			} else {
				resultMap["used_percent"] = resultMap["used_percent"].(float64) / float64(deviceCnt)
			}
		}
	}
	/*
		else if metricName == "diskio" {
			deviceCnt := len(deviceMap)
			if deviceCnt == 0 {
				resultMap["read_bytes"] = 0
				resultMap["write_bytes"] = 0
			} else {
				//resultMap["read_bytes"] = resultMap["read_bytes "].(float64) / float64(deviceCnt)
				//resultMap["write_bytes"] = resultMap["write_bytes"].(float64) / float64(deviceCnt)
				//resultMap["iops_read"] = resultMap["iops_read"].(float64) / float64(deviceCnt)
				//resultMap["iops_write"] = resultMap["iops_write"].(float64) / float64(deviceCnt)
			}
		}
	*/

	//spew.Dump(deviceMap)
	spew.Dump(resultMap)

	aggregateMap := make(map[string]interface{})
	metricMap := make(map[string]interface{})

	metricMap[metricName] = resultMap
	aggregateMap[vmId] = metricMap

	/* 4. 모니터링 데이터 저장 (InfluxDB) */
	err = a.InfluxDB.WriteMetric(aggregateMap)
	if err != nil {
		logrus.Error("Failed to save write disk, diskio metric", err)
		return err
	}

	return nil
}

// 실시간 모니터링 데이터 통계 값 계산 (MIN, MAX, AVG, LAST)
func (a *Aggregator) CalculateMetric(metricName string, metric map[string]interface{}, aggregateType string) (interface{}, error) {

	/* 1. 실시간 모니터링 데이터 계산을 위한 맵 생성 */
	metricArr := map[string][]float64{}
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
	//spew.Dump(metricArr)

	/* 2. 실시간 모니터링 데이터 통계 로직 적용 */
	resultMap := map[string]interface{}{}
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
}

// etcd 저장소에 저장된 모든 모니터링 데이터 삭제 (초기화)
func (a *Aggregator) FlushMetric(vmList []string) error {
	for _, vmId := range vmList {
		err := a.Etcd.DeleteMetric(fmt.Sprintf("/host/%s", vmId))
		if err != nil {
			return err
		}
	}
	return nil
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
	aggregateMetric, err := a.CalculateMetric(metricName, metricDetailMap, aggregateType)
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

		deviceMap[deviceName], err = a.CalculateMetric(metricName, metricDetailMap, aggregateType)
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
			//resultMap["used_percent"] = resultMap["used"].(float64) / resultMap["total"].(float64)
			deviceCnt := len(deviceMap)
			if deviceCnt == 0 {
				resultMap["used_percent"] = 0
			} else {
				resultMap["used_percent"] = resultMap["used_percent"].(float64) / float64(deviceCnt)
			}
		}
	}

	spew.Dump(resultMap)
	return resultMap, nil
}
