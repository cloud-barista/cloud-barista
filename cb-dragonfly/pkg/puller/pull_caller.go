package puller

import (
	"fmt"
	"net/http"
	"time"

	agentmetadata "github.com/cloud-barista/cb-dragonfly/pkg/core/agent"
	"github.com/cloud-barista/cb-dragonfly/pkg/core/metric"
	"github.com/cloud-barista/cb-dragonfly/pkg/metricstore/influxdb/v1"
	"github.com/cloud-barista/cb-dragonfly/pkg/types"
)

const (
	AgentUnhealthyCnt = 5
)

type PullCaller struct {
	AgentList map[string]agentmetadata.AgentInfo
}

func NewPullCaller(agentList map[string]agentmetadata.AgentInfo) (PullCaller, error) {
	return PullCaller{AgentList: agentList}, nil
}

func (pc PullCaller) StartPull() {
	for uuid, agent := range pc.AgentList {
		// Check agent status
		if agent.AgentState == string(agentmetadata.Disable) {
			continue
		}
		// Check agent health
		if agent.AgentHealth == string(agentmetadata.Unhealthy) {
			// Call healthcheck API
			err := pc.healthcheck(uuid, agent)
			if err != nil {
				fmt.Println(err)
			}
			continue
		}
		go pc.pullMetric(uuid, agent)
	}
	fmt.Println(fmt.Sprintf("[%s] finished pulling loop", time.Now().Local().String()))
}

func (pc PullCaller) healthcheck(uuid string, agent agentmetadata.AgentInfo) error {
	client := http.Client{
		Timeout: metric.AgentTimeout * time.Second,
	}
	agentUrl := fmt.Sprintf("http://%s:%d/cb-dragonfly/healthcheck", agent.PublicIp, metric.AgentPort)
	resp, _ := client.Get(agentUrl)
	if resp != nil {
		if resp.StatusCode == http.StatusNoContent {
			agent.AgentHealth = string(agentmetadata.Healthy)
			err := agentmetadata.PutAgentMetadataToStore(uuid, agent)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (pc PullCaller) pullMetric(uuid string, agent agentmetadata.AgentInfo) {

	pullerIdx := time.Now().Unix()
	metricArr := []types.Metric{types.Cpu, types.CpuFrequency, types.Memory, types.Disk, types.DiskIO, types.Network}
	for _, pullMetric := range metricArr {

		if agent.AgentState == string(agentmetadata.Disable) || agent.AgentHealth == string(agentmetadata.Unhealthy) {
			// TODO: Call healthcheck API
			continue
		}

		fmt.Printf("[%d][%s][%s] CALL API: http://%s:%d/cb-dragonfly/metric/%s\n", pullerIdx, time.Now().Local().String(), uuid, agent.PublicIp, metric.AgentPort, pullMetric.ToAgentMetricKey())

		// Pulling agent
		result, statusCode, err := metric.GetVMOnDemandMonInfo(pullMetric.ToString(), agent.PublicIp)

		// Update Agent Health
		updated := false
		if statusCode == http.StatusOK && agent.AgentHealth == string(agentmetadata.Unhealthy) {
			updated = true
			agent.AgentHealth = string(agentmetadata.Healthy)
		}
		if statusCode != http.StatusOK {
			updated = true
			agent.AgentUnhealthyRespCnt += 1
			if agent.AgentUnhealthyRespCnt > AgentUnhealthyCnt {
				agent.AgentHealth = string(agentmetadata.Unhealthy)
			}
		}

		if updated {
			err := agentmetadata.PutAgentMetadataToStore(uuid, agent)
			if err != nil {
				continue
			}
		}

		if result == nil {
			continue
		}

		// 메트릭 정보 파싱
		metricData := result.(map[string]interface{})
		metricName := metricData["name"].(string)
		if metricName == "" {
			continue
		}
		tagArr := map[string]string{}
		for k, v := range metricData["tags"].(map[string]interface{}) {
			tagArr[k] = v.(string)
		}
		metricVal := metricData["values"].(map[string]interface{})

		// 메트릭 정보 InfluxDB 저장
		err = v1.GetInstance().WriteOnDemandMetric(v1.PullDatabase, metricName, tagArr, metricVal)
		if err != nil {
			fmt.Println(err)
		}
	}
}
