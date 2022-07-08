package mcis

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/cloud-barista/cb-dragonfly/pkg/api/core/agent/common"
	"github.com/cloud-barista/cb-dragonfly/pkg/modules/monitoring/push/mcis/collector"
	"github.com/cloud-barista/cb-dragonfly/pkg/util"

	"github.com/cloud-barista/cb-dragonfly/pkg/types"
)

const (
	AgentPort    = 8888
	AgentTimeout = 30
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
	resultMetric, err := collector.ConvertMonMetric(metric, metricData[metric.ToAgentMetricKey()])
	if err != nil {
		return nil, err
	}
	return resultMetric, nil
}

type PacketsInfo struct {
	DestinationIp    string
	PacketCnt        int
	TotalPacketBytes int
	Msg              string
}

type NetworkPacketsResult struct {
	WatchTime    string
	PacketsInfos map[int]PacketsInfo
}

func GetMCISOnDemandPacketInfo(nsId string, mcisId string, vmId string, watchTime string) (NetworkPacketsResult, int, error) {
	agentList, err := common.ListAgent()
	if err != nil {
		fmt.Println("Fail to Get AgentList From CB-Store")
		return NetworkPacketsResult{}, http.StatusInternalServerError, err
	}
	var sourceAgentIP string
	var targetAgentInfo []common.AgentInfo

	for _, agentMetadata := range agentList {
		if agentMetadata.McisId == mcisId && agentMetadata.NsId == nsId {
			if agentMetadata.VmId == vmId {
				sourceAgentIP = agentMetadata.PublicIp
			} else {
				targetAgentInfo = append(targetAgentInfo, agentMetadata)
			}
		}
	}

	if sourceAgentIP == "" || len(targetAgentInfo) == 0 {
		return NetworkPacketsResult{}, http.StatusOK, nil
	}

	client := http.Client{
		Timeout: AgentTimeout * time.Second,
	}

	wg := sync.WaitGroup{}
	wg.Add(len(targetAgentInfo))

	result := NetworkPacketsResult{
		WatchTime:    watchTime,
		PacketsInfos: map[int]PacketsInfo{},
	}

	for idx, targetAgent := range targetAgentInfo {
		agentUrl := fmt.Sprintf("http://%s:%d/cb-dragonfly/mcis/dstip/%s/watchtime/%s", sourceAgentIP, AgentPort, targetAgent.PublicIp, watchTime)
		idx := idx
		targetAgent := targetAgent
		go func() {
			defer wg.Done()
			packetsInfo := PacketsInfo{}
			resp, err := client.Get(agentUrl)
			if err != nil {
				fmt.Println("err: "+targetAgent.PublicIp+", msg: ", err)
				return
			}
			body, err2 := ioutil.ReadAll(resp.Body)
			defer resp.Body.Close()
			if err2 != nil {
				fmt.Println("err: "+targetAgent.PublicIp+", msg: ", err2)
				return
			}
			_ = json.Unmarshal(body, &packetsInfo)
			result.PacketsInfos[idx] = packetsInfo
		}()
	}
	wg.Wait()
	return result, http.StatusOK, err
}

type ProcessUsage struct {
	Pid      string
	CpuUsage string
	MemUsage string
	Command  string
}

func GetMCISOnDemandProcessInfo(publicIp string) (map[string][]ProcessUsage, int, error) {

	userProcess := map[string][]ProcessUsage{}
	client := http.Client{
		Timeout: AgentTimeout * time.Second,
	}
	agentUrl := fmt.Sprintf("http://%s:%d/cb-dragonfly/mcis/process", publicIp, AgentPort)
	resp, err := client.Get(agentUrl)
	if err != nil {
		return map[string][]ProcessUsage{}, http.StatusInternalServerError, err
	}
	body, err2 := ioutil.ReadAll(resp.Body)
	if err2 != nil {
		return map[string][]ProcessUsage{}, http.StatusInternalServerError, err
	}
	_ = json.Unmarshal(body, &userProcess)

	return userProcess, http.StatusOK, nil

}

type vmSpec struct {
	Namespace      string
	Id             string
	Name           string
	ConnectionName string
	CspSpecName    string
	NumvCPU        int
	MemGiB         int
	CostPerHour    float64
}

type McisVMSpecs struct {
	AvgNumvCpu     float64
	AvgMemGiB      float64
	AvgCostPerHour float64
	VmSpec         []vmSpec
}

func GetMCISSpecInfo(nsId string, mcisId string, auth string) (McisVMSpecs, int, error) {
	mcisVMSpecs := McisVMSpecs{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/ns/%s/mcis/%s", types.TBRestAPIURL, nsId, mcisId), nil)
	req.Header.Add("Authorization", auth)
	client := &http.Client{}
	resp, _ := client.Do(req)
	defer resp.Body.Close()
	bytes, _ := ioutil.ReadAll(resp.Body)
	respStr := strings.Split(string(bytes), `"specId":`)[1:]
	if len(respStr) == 0 {
		return mcisVMSpecs, http.StatusOK, nil
	}
	//nsId = "common"
	for _, specIdStr := range respStr {
		specId := strings.ReplaceAll(strings.Split(specIdStr, ",")[0], `"`, "")
		req, err = http.NewRequest("GET", fmt.Sprintf("%s/ns/%s/resources/spec/%s", types.TBRestAPIURL, nsId, specId), nil)
		req.Header.Add("Authorization", auth)
		resp, _ = client.Do(req)
		bytes, _ = ioutil.ReadAll(resp.Body)
		vmSpec := vmSpec{}
		err = json.Unmarshal(bytes, &vmSpec)
		if err != nil {
			fmt.Println(err)
		}
		mcisVMSpecs.AvgNumvCpu += float64(vmSpec.NumvCPU)
		mcisVMSpecs.AvgMemGiB += float64(vmSpec.MemGiB)
		mcisVMSpecs.AvgCostPerHour += vmSpec.CostPerHour
		mcisVMSpecs.VmSpec = append(mcisVMSpecs.VmSpec, vmSpec)
	}
	vmSpecCnt := float64(len(mcisVMSpecs.VmSpec))
	mcisVMSpecs.AvgNumvCpu /= vmSpecCnt
	mcisVMSpecs.AvgMemGiB /= vmSpecCnt
	mcisVMSpecs.AvgCostPerHour /= vmSpecCnt
	mcisVMSpecs.AvgCostPerHour = util.ToFixed(mcisVMSpecs.AvgCostPerHour, 5)

	return mcisVMSpecs, http.StatusOK, nil
}
