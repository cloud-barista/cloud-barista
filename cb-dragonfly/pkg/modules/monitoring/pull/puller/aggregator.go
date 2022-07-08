package puller

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/cloud-barista/cb-dragonfly/pkg/api/core/agent/common"

	"github.com/cloud-barista/cb-dragonfly/pkg/config"
	"github.com/cloud-barista/cb-dragonfly/pkg/storage/cbstore"
	"github.com/cloud-barista/cb-dragonfly/pkg/storage/metricstore/influxdb/metric"
	v1 "github.com/cloud-barista/cb-dragonfly/pkg/storage/metricstore/influxdb/v1"
	"github.com/cloud-barista/cb-dragonfly/pkg/types"
	"github.com/sirupsen/logrus"
)

type PullAggregator struct {
	Storage   v1.Storage
	CBStore   cbstore.CBStore
	AgentList map[string]common.AgentInfo
}

func NewPullAggregator() (*PullAggregator, error) {
	pullAggregator := PullAggregator{
		Storage: *v1.GetInstance(),
		CBStore: *cbstore.GetInstance(),
	}
	return &pullAggregator, nil
}

func (pa *PullAggregator) StartAggregate() error {
	metricArr := []types.Metric{types.Cpu, types.CpuFrequency, types.Memory, types.Disk, types.Network, types.DiskIO}
	aggregateInterval := time.Duration(config.GetInstance().Monitoring.PullerAggregateInterval)
	for {
		time.Sleep(aggregateInterval * time.Second)

		err := pa.syncAgentList()
		if err != nil {
			fmt.Println(err)
			return err
		}

		if len(pa.AgentList) == 0 {
			time.Sleep(aggregateInterval * time.Second)
			continue
		}

		go pa.AggregateMetric(pa.AgentList, metricArr, config.GetInstance().Monitoring.AggregateType)

	}
}

func (pa *PullAggregator) AggregateMetric(agentList map[string]common.AgentInfo, metricArr []types.Metric, aggregateType string) {
	for _, targetAgent := range agentList {
		if targetAgent.AgentType == types.PushPolicy {
			continue
		}
		for _, metricKind := range metricArr {
			var receivedMetric interface{}
			var err error
			var mappedMetric interface{}
			monRequestInfo := types.DBMetricRequestInfo{
				MonitoringMechanism: config.GetInstance().Monitoring.DefaultPolicy == types.PushPolicy,
				NsID:                targetAgent.NsId,
				ServiceType:         types.MCIS,
				ServiceID:           targetAgent.McisId,
				VMID:                targetAgent.VmId,
				MetricName:          metricKind.ToAgentMetricKey(),
				Period:              "m",
				AggegateType:        aggregateType,
				Duration:            "5m",
			}
			receivedMetric, err = pa.Storage.ReadMetric(monRequestInfo)
			if err != nil {
				logrus.Println(err)
			}
			if receivedMetric == nil {
				continue
			}
			mappedMetric, err = metric.MappingMonMetric(metricKind.ToString(), &receivedMetric)
			var metricName string
			var valueLength float64
			tagArr := map[string]string{}
			reqValue := map[string]interface{}{}
			if metricKind.ToString() == types.Network.ToString() || metricKind.ToString() == types.DiskIO.ToString() {
				for k, v := range mappedMetric.(map[string]interface{}) {
					if k == "values" {
						for _, vv := range v.([]interface{}) {
							for metricKey, metricValue := range vv.(map[string]interface{}) {
								valueLength += 1
								if metricKey == "time" {
									continue
								}
								compare, _ := metricValue.(json.Number).Float64()
								if reqValue[metricKey] == nil || aggregateType == types.LAST.ToString() {
									reqValue[metricKey] = compare
								} else {
									origin, _ := reqValue[metricKey].(float64)
									var vSum float64
									if aggregateType == types.MAX.ToString() {
										if origin < compare {
											reqValue[metricKey] = compare
										}
									}
									if aggregateType == types.MIN.ToString() {
										if origin > compare {
											reqValue[metricKey] = compare
										}
									}
									if aggregateType == types.AVG.ToString() {
										vSum += compare
										reqValue[metricKey] = vSum / valueLength

									}
								}
							}
						}
					}
					if k == "name" {
						metricName = v.(string)
					}
					if k == "tags" {
						for tKey, tValue := range v.(map[string]string) {
							tagArr[tKey] = tValue
						}
					}
				}
			} else {
				convertedMetric := mappedMetric.(map[string]interface{})
				metricName = convertedMetric["name"].(string)
				for k, v := range convertedMetric["tags"].(map[string]string) {
					tagArr[k] = v
				}
				for _, value := range convertedMetric["values"].([]interface{}) {
					for k, v := range value.(map[string]interface{}) {
						if k == "time" {
							continue
						}
						if v == nil {
							v = json.Number("0")
						}
						inputData, _ := v.(json.Number).Float64()
						reqValue[k] = inputData

					}
				}

			}
			err = pa.Storage.WriteOnDemandMetric(v1.DefaultDatabase, metricName, tagArr, reqValue)
			if err != nil {
				logrus.Println(err)
			}
			err = pa.Storage.DeleteMetric(v1.PullDatabase, metricName, "5m")
			if err != nil {
				logrus.Println(err)
			}
		}
	}
}

func (pa *PullAggregator) CalculateMetric() (map[string]interface{}, error) {
	return nil, nil
}

func (pa *PullAggregator) syncAgentList() error {
	syncedAgentList, err := common.ListAgent()
	if err != nil {
		fmt.Println(err)
		return err
	}
	pa.AgentList = syncedAgentList
	return nil
}
