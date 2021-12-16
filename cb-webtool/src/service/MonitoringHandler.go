package service

import (
	"encoding/json"
	"fmt"
	"io"

	// "io/ioutil"
	"log"
	"net/http"

	// "strconv"
	// "os"
	model "github.com/cloud-barista/cb-webtool/src/model"
	// "github.com/cloud-barista/cb-webtool/src/model/spider"
	"github.com/cloud-barista/cb-webtool/src/model/dragonfly"
	// "github.com/cloud-barista/cb-webtool/src/model/tumblebug"
	// tbcommon "github.com/cloud-barista/cb-webtool/src/model/tumblebug/common"
	// tbmcir "github.com/cloud-barista/cb-webtool/src/model/tumblebug/mcir"
	tbmcis "github.com/cloud-barista/cb-webtool/src/model/tumblebug/mcis"

	util "github.com/cloud-barista/cb-webtool/src/util"
)

// VM 에 모니터링 Agent 설치
///ns/{nsId}/monitoring/install/mcis/{mcisId}
func RegBenchmarkAgentInVm(nameSpaceID string, mcisID string, vmMonitoringAgentReg *tbmcis.McisCmdReq) (*tbmcis.AgentInstallContentWrapper, model.WebStatus) {
	fmt.Println("RegBenchmarkAgentInVm ************ : ")
	var originalUrl = "/ns/{nsId}/monitoring/install/mcis/{mcisId}"

	var paramMapper = make(map[string]string)
	paramMapper["{nsId}"] = nameSpaceID
	paramMapper["{mcisId}"] = mcisID
	urlParam := util.MappingUrlParameter(originalUrl, paramMapper)

	url := util.TUMBLEBUG + urlParam
	// url := util.TUMBLEBUG + "/ns/" + nameSpaceID + "/monitoring/install/mcis/"

	pbytes, _ := json.Marshal(vmMonitoringAgentReg)
	fmt.Println(string(pbytes))
	resp, err := util.CommonHttp(url, pbytes, http.MethodPost)

	vmMonitoringAgentInfo := tbmcis.AgentInstallContentWrapper{}
	if err != nil {
		fmt.Println(err)
		return &vmMonitoringAgentInfo, model.WebStatus{StatusCode: 500, Message: err.Error()}
	}

	respBody := resp.Body
	respStatus := resp.StatusCode

	returnStatus := model.WebStatus{}
	if respStatus != 200 && respStatus != 201 { // 호출은 정상이나, 가져온 결과값이 200, 201아닌 경우 message에 담겨있는 것을 WebStatus에 set
		errorInfo := model.ErrorInfo{}
		json.NewDecoder(respBody).Decode(&errorInfo)
		fmt.Println("respStatus != 200 reason ", errorInfo)
		returnStatus.Message = errorInfo.Message
	} else {
		json.NewDecoder(respBody).Decode(&vmMonitoringAgentInfo)
		fmt.Println(vmMonitoringAgentInfo)
	}
	returnStatus.StatusCode = respStatus

	return &vmMonitoringAgentInfo, returnStatus
}

func RegMonitoringAgentInVm(nameSpaceID string, mcisID string, vmMonitoringAgentReg *dragonfly.VmMonitoringInstallReg) (*model.WebStatus, model.WebStatus) {
	fmt.Println("RegMonitoringAgentInVm ************ : ")
	//var originalUrl = "/agent/install"
	var originalUrl = "/agent"
	urlParam := util.MappingUrlParameter(originalUrl, nil)

	// Command  string `json:"command"`
	// PublicIp string `json:"ip"`
	// McisID   string `json:"mcis_id"`
	// SshKey   string `json:"ssh_key"`
	// UserName string `json:"user_name"`
	// VmID     string `json:"vm_id"`
	url := util.DRAGONFLY + urlParam
	// url := util.TUMBLEBUG + "/ns/" + nameSpaceID + "/monitoring/install/mcis/"
	// _ = writer.WriteField("mcis_id", "7e3130a0-a811-47b8-a82c-b155267edef5")
	// _ = writer.WriteField("vm_id", "c9f01668-4db4-4521-8d54-14501e31d2c7")
	// _ = writer.WriteField("public_ip", "0.0.0.0")
	// _ = writer.WriteField("user_name", "cbuser")
	// _ = writer.WriteField("ssh_key", "-----BEGIN RSA PRIVATE KEY-----\n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\n-----END RSA PRIVATE KEY-----\n")
	// _ = writer.WriteField("cspType", "test")

	// pbytes, _ := json.Marshal(vmMonitoringAgentReg)
	// fmt.Println(string(pbytes))
	// resp, err := util.CommonHttp(url, pbytes, http.MethodPost)

	fmt.Println("RegMonitoringAgentInVm : ", url)

	fmt.Println(vmMonitoringAgentReg)

	urlValues, convertErr := util.StructToMapByJson(vmMonitoringAgentReg)
	if convertErr != nil {
		log.Println(convertErr)
	}
	//vmMonitoringInstallReg
	fmt.Println(urlValues)
	//resp, err := util.CommonHttpFormData(url, urlValues, http.MethodPost)
	pbytes, _ := json.Marshal(vmMonitoringAgentReg)
	resp, err := util.CommonHttp(url, pbytes, http.MethodPost)

	webStatus := model.WebStatus{}
	if err != nil {
		fmt.Println(err)
		return &webStatus, model.WebStatus{StatusCode: 500, Message: err.Error()}
	}
	// return body, err
	respBody := resp.Body
	respStatus := resp.StatusCode
	resultInfo := model.ResultInfo{}

	json.NewDecoder(respBody).Decode(&resultInfo)
	log.Println(resultInfo)
	log.Println("ResultMessage : " + resultInfo.Message)

	if respStatus != 200 && respStatus != 201 {
		return &model.WebStatus{}, model.WebStatus{StatusCode: respStatus, Message: resultInfo.Message}
	}
	webStatus.StatusCode = respStatus
	webStatus.Message = resultInfo.Message
	return &webStatus, model.WebStatus{StatusCode: respStatus}
}

// Get Monitoring Data
func GetVmMonitoringInfoData(nameSpaceID string, mcisID string, metric string) (*tbmcis.MonResultSimpleResponse, model.WebStatus) {
	var originalUrl = "/ns/{nsId}/monitoring/mcis/{mcisId}/metric/{metric}"

	var paramMapper = make(map[string]string)
	paramMapper["{nsId}"] = nameSpaceID
	paramMapper["{mcisId}"] = mcisID
	paramMapper["{metric}"] = metric
	urlParam := util.MappingUrlParameter(originalUrl, paramMapper)

	url := util.TUMBLEBUG + urlParam
	// url := util.TUMBLEBUG + "/ns/" + nameSpaceID + "/monitoring/mcis/" + mcisID + "/metric/" + metric

	resp, err := util.CommonHttpWithoutParam(url, http.MethodGet)

	// defer body.Close()
	vmMonitoringResult := tbmcis.MonResultSimpleResponse{}
	if err != nil {
		fmt.Println(err)
		return &vmMonitoringResult, model.WebStatus{StatusCode: 500, Message: err.Error()}
	}

	respBody := resp.Body
	respStatus := resp.StatusCode

	json.NewDecoder(respBody).Decode(&vmMonitoringResult)
	fmt.Println(vmMonitoringResult)

	return &vmMonitoringResult, model.WebStatus{StatusCode: respStatus}
}

// VM monitoring
// Get vm monitoring info
// 멀티 클라우드 인프라 VM 모니터링 정보 조회
func GetVmMonitoring(vmMonitoring *dragonfly.VmMonitoring) (map[string]interface{}, model.WebStatus) {
	//func GetVmMonitoring(vmMonitoring *dragonfly.VmMonitoring) (*dragonfly.VmMonitoringInfo, model.WebStatus) {
	nameSpaceID := vmMonitoring.NameSpaceID
	mcisID := vmMonitoring.McisID
	vmID := vmMonitoring.VmID
	metric := vmMonitoring.Metric
	periodType := vmMonitoring.PeriodType
	statisticsCriteria := vmMonitoring.StatisticsCriteria
	duration := vmMonitoring.Duration

	var originalUrl = "/ns/:ns_id/mcis/:mcis_id/vm/:vm_id/metric/:metric_name/info?periodType={periodType}&statisticsCriteria={statisticsCriteria}&duration={duration}"
	//{{ip}}:{{port}}/dragonfly/ns/:ns_id/mcis/:mcis_id/vm/:vm_id/metric/:metric_name/info?periodType=m&statisticsCriteria=last&duration=5m
	var paramMapper = make(map[string]string)
	paramMapper[":ns_id"] = nameSpaceID
	paramMapper[":mcis_id"] = mcisID
	paramMapper[":vm_id"] = vmID
	paramMapper[":metric_name"] = metric
	paramMapper["{periodType}"] = periodType
	paramMapper["{statisticsCriteria}"] = statisticsCriteria
	paramMapper["{duration}"] = duration
	urlParam := util.MappingUrlParameter(originalUrl, paramMapper)

	url := util.DRAGONFLY + urlParam

	//"/mcis/"+mcis_id+"/vm/"+vm_id+"/metric/"+metric+"/info?periodType="+periodType+"&statisticsCriteria="+statisticsCriteria+"&duration="+duration;
	// urlParam := "periodType=" + vmMonitoring.PeriodType + "&statisticsCriteria=" + vmMonitoring.StatisticsCriteria + "&duration=" + vmMonitoring.Duration
	// url := util.DRAGONFLY + "/ns/" + vmMonitoring.NameSpaceID + "/mcis/" + vmMonitoring.McisID + "/vm/" + vmMonitoring.VmID + "/metric/" + vmMonitoring.Metric + "/info?" + urlParam

	// resp, err := util.CommonHttp(url, nil, http.MethodGet)
	resp, err := util.CommonHttpWithoutParam(url, http.MethodGet)

	// defer body.Close()
	//vmMonitoringInfo := dragonfly.VmMonitoringInfo{}
	vmMonitoringInfo := make(map[string]interface{})
	if err != nil {
		fmt.Println(err)
		return vmMonitoringInfo, model.WebStatus{StatusCode: 500, Message: err.Error()}
	}
	// util.DisplayResponse(resp) // 수신내용 확인

	respBody := resp.Body
	respStatus := resp.StatusCode

	if metric == "cpu" {
		vmMonitoringInfoByCpu := dragonfly.VmMonitoringInfoByCpu{}
		json.NewDecoder(respBody).Decode(&vmMonitoringInfoByCpu)

		//vmMonitoringInfo.ValuesByCpu = vmMonitoringInfoByCpu
		vmMonitoringInfo[metric] = vmMonitoringInfoByCpu
	} else if metric == "memory" {
		vmMonitoringInfoByMemory := dragonfly.VmMonitoringInfoByMemory{}
		json.NewDecoder(respBody).Decode(&vmMonitoringInfoByMemory)

		//vmMonitoringInfo.ValuesByMemory = vmMonitoringInfoByMemory
		vmMonitoringInfo[metric] = vmMonitoringInfoByMemory
	} else if metric == "disk" {
		vmMonitoringInfoByDisk := dragonfly.VmMonitoringInfoByDisk{}
		json.NewDecoder(respBody).Decode(&vmMonitoringInfoByDisk)

		//vmMonitoringInfo.ValuesByDisk = vmMonitoringInfoByDisk
		vmMonitoringInfo[metric] = vmMonitoringInfoByDisk
	} else if metric == "network" {
		vmMonitoringInfoByNetwork := dragonfly.VmMonitoringInfoByNetwork{}
		json.NewDecoder(respBody).Decode(&vmMonitoringInfoByNetwork)

		//vmMonitoringInfo.ValuesByNetwork = vmMonitoringInfoByNetwork
		vmMonitoringInfo[metric] = vmMonitoringInfoByNetwork
	}

	//json.NewDecoder(respBody).Decode(&vmMonitoringInfo)
	//fmt.Println(vmMonitoringInfo)

	//return &vmMonitoringInfo, model.WebStatus{StatusCode: respStatus}
	return vmMonitoringInfo, model.WebStatus{StatusCode: respStatus}
}

// 멀티 클라우드 인프라 VM 온디맨드 모니터링 정보 조회
// Get MCIS on-demand monitoring metric info
func GetMcisOnDemandMonitoringMetricInfo(agentIp string, metricName string, vmMonitoring *dragonfly.VmMonitoring) (*dragonfly.McisMonitoringOnDemandInfo, model.WebStatus) {
	nameSpaceID := vmMonitoring.NameSpaceID
	mcisID := vmMonitoring.McisID
	vmID := vmMonitoring.VmID
	// agentIp := vmMonitoring.AgentIp
	// metricName := vmMonitoring.MetricName

	var originalUrl = "/ns/:ns_id/mcis/:mcis_id/vm/:vm_id/agent_ip/:agent_ip/mcis_metric/:metric_name/mcis-monitoring-info"
	//{{ip}}:{{port}}/dragonfly/ns/:ns_id/mcis/:mcis_id/vm/:vm_id/agent_ip/:agent_ip/mcis_metric/:metric_name/mcis-monitoring-info
	var paramMapper = make(map[string]string)
	paramMapper[":ns_id"] = nameSpaceID
	paramMapper[":mcis_id"] = mcisID
	paramMapper[":vm_id"] = vmID
	paramMapper[":agent_ip"] = agentIp       // 에이전트 아이피
	paramMapper[":metric_name"] = metricName // 메트릭 정보 ( "InitDB" | "ResetDB" | "CpuM" | "CpuS" | "MemR" | "MemW" | "FioW" | "FioR" | "DBW" | DBR" | "Rtt" | "Mrtt" )
	urlParam := util.MappingUrlParameter(originalUrl, paramMapper)

	url := util.DRAGONFLY + urlParam

	// url := util.DRAGONFLY + "/ns/" + vmMonitoring.NameSpaceID + "/mcis/" + vmMonitoring.McisID + "/vm/" + vmMonitoring.VmID + "/agent_ip/" + vmMonitoring.AgentIP + "/mcis_metric/" + vmMonitoring.MetricName + "/mcis-monitoring-info"
	// url := util.DRAGONFLY + "/ns/" + vmMonitoring.NameSpaceID + "/mcis/" + vmMonitoring.McisID + "/vm/" + vmMonitoring.VmID // TODO : 객체에 parameter추가해야 함

	resp, err := util.CommonHttpWithoutParam(url, http.MethodGet)

	// defer body.Close()
	mcisMonitoringInfo := dragonfly.McisMonitoringOnDemandInfo{}
	if err != nil {
		fmt.Println(err)
		return &mcisMonitoringInfo, model.WebStatus{StatusCode: 500, Message: err.Error()}
	}
	// util.DisplayResponse(resp) // 수신내용 확인

	respBody := resp.Body
	respStatus := resp.StatusCode

	json.NewDecoder(respBody).Decode(&mcisMonitoringInfo)
	fmt.Println(mcisMonitoringInfo)

	return &mcisMonitoringInfo, model.WebStatus{StatusCode: respStatus}
}

// 멀티 클라우드 인프라 VM 온디맨드 모니터링 정보 조회
// Get vm on-demand monitoring metric info
func GetVmOnDemandMonitoringMetricInfo(agentIp string, metricName string, vmMonitoring *dragonfly.VmMonitoring) (*dragonfly.VmMonitoringOnDemandInfo, model.WebStatus) {
	nameSpaceID := vmMonitoring.NameSpaceID
	mcisID := vmMonitoring.McisID
	vmID := vmMonitoring.VmID
	// agentIp := vmMonitoring.AgentIp
	// metricName := vmMonitoring.MetricName

	var originalUrl = "/ns/:ns_id/mcis/:mcis_id/vm/:vm_id/agent_ip/:agent_ip/metric/:metric_name/ondemand-monitoring-info"
	// {{ip}}:{{port}}/dragonfly/ns/:ns_id/mcis/:mcis_id/vm/:vm_id/agent_ip/:agent_ip/metric/:metric_name/ondemand-monitoring-info
	var paramMapper = make(map[string]string)
	paramMapper[":ns_id"] = nameSpaceID
	paramMapper[":mcis_id"] = mcisID
	paramMapper[":vm_id"] = vmID
	paramMapper[":agent_ip"] = agentIp
	paramMapper[":metric_name"] = metricName
	urlParam := util.MappingUrlParameter(originalUrl, paramMapper)

	url := util.DRAGONFLY + urlParam

	// url := util.DRAGONFLY + "/ns/" + vmMonitoring.NameSpaceID + "/mcis/" + vmMonitoring.McisID + "/vm/" + vmMonitoring.VmID + "/agent_ip/" + vmMonitoring.AgentIP + "/mcis_metric/" + vmMonitoring.MetricName + "/mcis-monitoring-info"
	// url := util.DRAGONFLY + "/ns/" + vmMonitoring.NameSpaceID + "/mcis/" + vmMonitoring.McisID + "/vm/" + vmMonitoring.VmID // TODO : 객체에 parameter추가해야 함

	resp, err := util.CommonHttpWithoutParam(url, http.MethodGet)

	// defer body.Close()
	vmMonitoringInfo := dragonfly.VmMonitoringOnDemandInfo{}
	if err != nil {
		fmt.Println(err)
		return &vmMonitoringInfo, model.WebStatus{StatusCode: 500, Message: err.Error()}
	}
	// util.DisplayResponse(resp) // 수신내용 확인

	respBody := resp.Body
	respStatus := resp.StatusCode

	json.NewDecoder(respBody).Decode(&vmMonitoringInfo)
	fmt.Println(vmMonitoringInfo)

	return &vmMonitoringInfo, model.WebStatus{StatusCode: respStatus}
}

// 모니터링 정책 조회
// Get monitoring config
func GetMonitoringConfig() (*dragonfly.MonitoringConfig, model.WebStatus) {
	var originalUrl = "/config"
	//{{ip}}:{{port}}/dragonfly/config
	urlParam := util.MappingUrlParameter(originalUrl, nil)

	url := util.DRAGONFLY + urlParam
	// url := util.DRAGONFLY + "/config"
	resp, err := util.CommonHttpWithoutParam(url, http.MethodGet)

	// defer body.Close()
	monitoringConfig := dragonfly.MonitoringConfig{}

	if err != nil {
		log.Println(err)
		return &monitoringConfig, model.WebStatus{StatusCode: 500, Message: err.Error()}
	}
	// util.DisplayResponse(resp) // 수신내용 확인
	respBody := resp.Body
	respStatus := resp.StatusCode

	json.NewDecoder(respBody).Decode(&monitoringConfig)
	log.Println(monitoringConfig)

	return &monitoringConfig, model.WebStatus{StatusCode: respStatus}
}

// 모니터링 정책 설정
func PutMonigoringConfig(monitoringConfigReg *dragonfly.MonitoringConfigReg) (*dragonfly.MonitoringConfig, model.WebStatus) {
	var originalUrl = "/config"
	//{{ip}}:{{port}}/dragonfly/config
	urlParam := util.MappingUrlParameter(originalUrl, nil)

	url := util.DRAGONFLY + urlParam
	// url := util.DRAGONFLY + "/config"

	fmt.Println("Update MonigoringConfigReg : ", url)

	//fmt.Println(monitoringConfigReg)
	//
	urlValues, convertErr := util.StructToMapByJson(monitoringConfigReg)
	if convertErr != nil {
		log.Println(convertErr)
	}

	fmt.Println(urlValues)
	//resp, err := util.CommonHttpFormData(url, urlValues, http.MethodPut)
	pbytes, _ := json.Marshal(monitoringConfigReg)
	resp, err := util.CommonHttp(url, pbytes, http.MethodPut)
	resultMonitoringConfig := dragonfly.MonitoringConfig{}

	if err != nil {
		log.Println("-----")
		fmt.Println(err)
		log.Println("-----1111")
		fmt.Println(err.Error())
		log.Println("-----222")
		return &resultMonitoringConfig, model.WebStatus{StatusCode: 500, Message: err.Error()}
	}

	respBody := resp.Body
	respStatus := resp.StatusCode
	log.Println("respStatusCode = ", resp.StatusCode)
	log.Println("respStatus = ", resp.Status)
	if respStatus != 200 && respStatus != 201 {
		errorInfo := model.ErrorInfo{}
		json.NewDecoder(respBody).Decode(&errorInfo)
		fmt.Println(errorInfo)
		return nil, model.WebStatus{StatusCode: 500, Message: errorInfo.Message}
	}

	// 응답에 생성한 객체값이 옴
	json.NewDecoder(respBody).Decode(&resultMonitoringConfig)
	fmt.Println(resultMonitoringConfig)
	// return respBody, respStatusCode
	return &resultMonitoringConfig, model.WebStatus{StatusCode: respStatus}
}

// 모니터링 정책 초기화
func ResetMonigoringConfig(monitoringConfig *dragonfly.MonitoringConfig) (*dragonfly.MonitoringConfig, model.WebStatus) {
	var originalUrl = "/config/reset"
	//{{ip}}:{{port}}/dragonfly/config/reset
	urlParam := util.MappingUrlParameter(originalUrl, nil)

	url := util.DRAGONFLY + urlParam
	// url := util.DRAGONFLY + "/config/reset"

	resp, err := util.CommonHttp(url, nil, http.MethodPut)
	resultMonitoringConfig := dragonfly.MonitoringConfig{}
	if err != nil {
		log.Println("-----")
		fmt.Println(err)
		log.Println("-----1111")
		fmt.Println(err.Error())
		log.Println("-----222")
		return &resultMonitoringConfig, model.WebStatus{StatusCode: 500, Message: err.Error()}
	}

	respBody := resp.Body
	respStatus := resp.StatusCode
	log.Println("respStatusCode = ", resp.StatusCode)
	log.Println("respStatus = ", resp.Status)
	if respStatus != 200 && respStatus != 201 {
		errorInfo := model.ErrorInfo{}
		json.NewDecoder(respBody).Decode(&errorInfo)
		fmt.Println(errorInfo)
		return nil, model.WebStatus{StatusCode: 500, Message: errorInfo.Message}
	}

	// 응답에 생성한 객체값이 옴
	json.NewDecoder(respBody).Decode(&resultMonitoringConfig)
	fmt.Println(resultMonitoringConfig)
	// return respBody, respStatusCode
	return &resultMonitoringConfig, model.WebStatus{StatusCode: respStatus}
}

// Install agent to vm
// 모니터링 에이전트 설치 : 위에 RegMonitoringAgentInVm 와 뭐가 다른거지?
func InstallAgentToVm(nameSpaceID string, vmMonitoringInstallReg *dragonfly.VmMonitoringInstallReg) (*dragonfly.VmMonitoringInstallReg, model.WebStatus) {
	//var originalUrl = "/agent/install"
	var originalUrl = "/agent/install"
	//{{ip}}:{{port}}/dragonfly/agent/install
	urlParam := util.MappingUrlParameter(originalUrl, nil)

	url := util.DRAGONFLY + urlParam
	// url := util.DRAGONFLY + "/agent/install/"

	pbytes, _ := json.Marshal(vmMonitoringInstallReg)
	resp, err := util.CommonHttp(url, pbytes, http.MethodPost)

	returnVmMonitoringInstallReg := dragonfly.VmMonitoringInstallReg{}
	returnStatus := model.WebStatus{}

	respBody := resp.Body
	respStatus := resp.StatusCode

	if err != nil {
		fmt.Println(err)
		return &returnVmMonitoringInstallReg, model.WebStatus{StatusCode: 500, Message: err.Error()}
	}

	if respStatus != 200 && respStatus != 201 { // 호출은 정상이나, 가져온 결과값이 200, 201아닌 경우 message에 담겨있는 것을 WebStatus에 set
		errorInfo := model.ErrorInfo{}
		json.NewDecoder(respBody).Decode(&errorInfo)
		fmt.Println("respStatus != 200 reason ", errorInfo)
		returnStatus.Message = errorInfo.Message
	} else {
		json.NewDecoder(respBody).Decode(&returnVmMonitoringInstallReg)
		fmt.Println(returnVmMonitoringInstallReg)
	}
	returnStatus.StatusCode = respStatus

	return &returnVmMonitoringInstallReg, returnStatus
}

// 모니터링 에이전트 제거
// Uninstall agent to vm
func UnInstallAgentToVm(nameSpaceID string, vmMonitoringInstallReg *dragonfly.VmMonitoringInstallReg) (*dragonfly.VmMonitoringInstallReg, model.WebStatus) {
	//var originalUrl = "/agent/uninstall"
	var originalUrl = "/agent"
	//{{ip}}:{{port}}/dragonfly/agent/uninstall
	urlParam := util.MappingUrlParameter(originalUrl, nil)

	url := util.DRAGONFLY + urlParam
	// url := util.DRAGONFLY + "/agent/uninstall/"

	pbytes, _ := json.Marshal(vmMonitoringInstallReg)
	resp, err := util.CommonHttp(url, pbytes, http.MethodPost)

	returnVmMonitoringInstallReg := dragonfly.VmMonitoringInstallReg{}
	returnStatus := model.WebStatus{}

	respBody := resp.Body
	respStatus := resp.StatusCode

	if err != nil {
		fmt.Println(err)
		return &returnVmMonitoringInstallReg, model.WebStatus{StatusCode: 500, Message: err.Error()}
	}

	if respStatus != 200 && respStatus != 201 { // 호출은 정상이나, 가져온 결과값이 200, 201아닌 경우 message에 담겨있는 것을 WebStatus에 set
		errorInfo := model.ErrorInfo{}
		json.NewDecoder(respBody).Decode(&errorInfo)
		fmt.Println("respStatus != 200 reason ", errorInfo)
		returnStatus.Message = errorInfo.Message
	} else {
		json.NewDecoder(respBody).Decode(&returnVmMonitoringInstallReg)
		fmt.Println(returnVmMonitoringInstallReg)
	}
	returnStatus.StatusCode = respStatus

	return &returnVmMonitoringInstallReg, returnStatus
}

// 알람 목록 조회
// List monitoring alert
func GetMonitoringAlertList() ([]dragonfly.VmMonitoringAlertInfo, model.WebStatus) {
	fmt.Print("#########GetMonitoringAlertList############")
	var originalUrl = "/alert/tasks"
	// {{ip}}:{{port}}/dragonfly/alert/tasks
	urlParam := util.MappingUrlParameter(originalUrl, nil)

	url := util.DRAGONFLY + urlParam
	// url := util.DRAGONFLY + "/alert/tasks"
	resp, err := util.CommonHttp(url, nil, http.MethodGet)

	// vmMonitoringAlertInfoList := dragonfly.VmMonitoringAlertInfo{}
	vmMonitoringAlertInfoList := []dragonfly.VmMonitoringAlertInfo{}
	if err != nil {
		fmt.Println(err)
		return vmMonitoringAlertInfoList, model.WebStatus{StatusCode: 500, Message: err.Error()}
	}

	respBody := resp.Body
	respStatus := resp.StatusCode

	//vmMonitoringAlertInfoList := []dragonfly.VmMonitoringAlertInfo{}
	json.NewDecoder(respBody).Decode(&vmMonitoringAlertInfoList)

	// robots, _ := ioutil.ReadAll(resp.Body)
	// log.Println(fmt.Print(string(robots)))

	// json.NewDecoder(respBody).Decode(&vmMonitoringAlertInfoList)
	// fmt.Println(vmMonitoringAlertInfoList)

	return vmMonitoringAlertInfoList, model.WebStatus{StatusCode: respStatus}
}

// 알람  조회
// monitoring alert
func GetMonitoringAlertData(taskName string) (dragonfly.VmMonitoringAlertInfo, model.WebStatus) {
	var originalUrl = "/alert/task/:task_name"
	// {{ip}}:{{port}}/dragonfly/alert/task/:task_name
	var paramMapper = make(map[string]string)
	paramMapper[":task_name"] = taskName
	urlParam := util.MappingUrlParameter(originalUrl, paramMapper)

	url := util.DRAGONFLY + urlParam
	// url := util.DRAGONFLY + "/alert/task/" + taskName
	resp, err := util.CommonHttp(url, nil, http.MethodGet)

	vmMonitoringAlertInfo := dragonfly.VmMonitoringAlertInfo{}
	if err != nil {
		fmt.Println(err)
		return vmMonitoringAlertInfo, model.WebStatus{StatusCode: 500, Message: err.Error()}
	}

	respBody := resp.Body
	respStatus := resp.StatusCode

	json.NewDecoder(respBody).Decode(&vmMonitoringAlertInfo)

	return vmMonitoringAlertInfo, model.WebStatus{StatusCode: respStatus}
}

// 알람 생성
// Create Monitoring Alert
func RegMonitoringAlert(vmMonitoringAlertInfo *dragonfly.VmMonitoringAlertInfo) (*dragonfly.VmMonitoringAlertInfo, model.WebStatus) {
	fmt.Println("RegMonitoringAlert ************ : ")
	var originalUrl = "/alert/task"
	// {{ip}}:{{port}}/dragonfly/alert/task
	urlParam := util.MappingUrlParameter(originalUrl, nil)

	url := util.DRAGONFLY + urlParam
	// url := util.DRAGONFLY + "/alert/task"

	pbytes, _ := json.Marshal(vmMonitoringAlertInfo)
	fmt.Println(string(pbytes))
	resp, err := util.CommonHttp(url, pbytes, http.MethodPost)

	//urlValues, convertErr := util.StructToMapByJson(vmMonitoringAlertInfo)
	//if convertErr != nil {
	//	log.Println(convertErr)
	//}
	//
	//fmt.Println(urlValues)
	//resp, err := util.CommonHttpFormData(url, urlValues, http.MethodPost)

	resultVmMonitoringAlertInfo := dragonfly.VmMonitoringAlertInfo{}
	if err != nil {
		fmt.Println(err)
		return &resultVmMonitoringAlertInfo, model.WebStatus{StatusCode: 500, Message: err.Error()}
	}

	respBody := resp.Body
	respStatus := resp.StatusCode

	returnStatus := model.WebStatus{}
	if respStatus != 200 && respStatus != 201 { // 호출은 정상이나, 가져온 결과값이 200, 201아닌 경우 message에 담겨있는 것을 WebStatus에 set
		errorInfo := model.ErrorInfo{}
		json.NewDecoder(respBody).Decode(&errorInfo)
		fmt.Println("respStatus != 200 reason ", errorInfo)
		returnStatus.Message = errorInfo.Message
	} else {
		json.NewDecoder(respBody).Decode(&resultVmMonitoringAlertInfo)
		fmt.Println(resultVmMonitoringAlertInfo)
	}
	returnStatus.StatusCode = respStatus

	return &resultVmMonitoringAlertInfo, returnStatus
}

// 알람 수정
// Update Monitoring Alert
func PutMonitoringAlert(taskName string, vmMonitoringAlertInfo *dragonfly.VmMonitoringAlertInfo) (*dragonfly.VmMonitoringAlertInfo, model.WebStatus) {
	fmt.Println("PutMonitoringAlert ************ : ")
	var originalUrl = "/alert/task"
	// {{ip}}:{{port}}/dragonfly/alert/task
	urlParam := util.MappingUrlParameter(originalUrl, nil)

	url := util.DRAGONFLY + urlParam
	// url := util.DRAGONFLY + "/alert/task/" + taskName

	pbytes, _ := json.Marshal(vmMonitoringAlertInfo)
	fmt.Println(string(pbytes))
	resp, err := util.CommonHttp(url, pbytes, http.MethodPut)

	resultVmMonitoringAlertInfo := dragonfly.VmMonitoringAlertInfo{}
	if err != nil {
		fmt.Println(err)
		return &resultVmMonitoringAlertInfo, model.WebStatus{StatusCode: 500, Message: err.Error()}
	}

	respBody := resp.Body
	respStatus := resp.StatusCode

	returnStatus := model.WebStatus{}
	if respStatus != 200 && respStatus != 201 { // 호출은 정상이나, 가져온 결과값이 200, 201아닌 경우 message에 담겨있는 것을 WebStatus에 set
		errorInfo := model.ErrorInfo{}
		json.NewDecoder(respBody).Decode(&errorInfo)
		fmt.Println("respStatus != 200 reason ", errorInfo)
		returnStatus.Message = errorInfo.Message
	} else {
		json.NewDecoder(respBody).Decode(&resultVmMonitoringAlertInfo)
		fmt.Println(resultVmMonitoringAlertInfo)
	}
	returnStatus.StatusCode = respStatus

	return &resultVmMonitoringAlertInfo, returnStatus
}

// 알람 제거
// Delete Monitoring Alert
func DelMonitoringAlert(taskName string) (io.ReadCloser, model.WebStatus) {
	var originalUrl = "/alert/task/:task_name"
	// {{ip}}:{{port}}/dragonfly/alert/task/:task_name
	var paramMapper = make(map[string]string)
	paramMapper[":task_name"] = taskName
	urlParam := util.MappingUrlParameter(originalUrl, paramMapper)

	url := util.DRAGONFLY + urlParam
	// url := util.DRAGONFLY + "/alert/task/" + taskName

	if taskName == "" {
		return nil, model.WebStatus{StatusCode: 500, Message: "TaskName is required"}
	}

	// 경로안에 parameter가 있어 추가 param없이 호출 함.
	resp, err := util.CommonHttp(url, nil, http.MethodDelete)
	if err != nil {
		fmt.Println(err)
		return nil, model.WebStatus{StatusCode: 500, Message: err.Error()}
	}
	// return body, err
	respBody := resp.Body
	respStatus := resp.StatusCode
	return respBody, model.WebStatus{StatusCode: respStatus}
}

// 알람 이벤트 핸들러 조회
// Get monitoring alert event-handler
// type : 이벤트 핸들러 유형 ( "slack" | "smtp" )
// name : slackHandler(EventHandlerName)
func GetMonitoringAlertEventHandlerData(eventHandlerType string, eventName string) (dragonfly.VmMonitoringAlertInfo, model.WebStatus) {
	var originalUrl = "/alert/eventhandler/type/:type/event/:name"
	//{{ip}}:{{port}}/dragonfly/alert/eventhandler/type/:type/event/:name
	var paramMapper = make(map[string]string)
	paramMapper[":type"] = eventHandlerType
	paramMapper[":name"] = eventName
	urlParam := util.MappingUrlParameter(originalUrl, paramMapper)

	url := util.DRAGONFLY + urlParam
	// url := util.DRAGONFLY + "/alert/eventhandler/type/" + eventHandlerType + "/event/" + eventName

	resp, err := util.CommonHttp(url, nil, http.MethodGet)

	vmMonitoringAlertInfo := dragonfly.VmMonitoringAlertInfo{}
	if err != nil {
		fmt.Println(err)
		return vmMonitoringAlertInfo, model.WebStatus{StatusCode: 500, Message: err.Error()}
	}

	respBody := resp.Body
	respStatus := resp.StatusCode

	json.NewDecoder(respBody).Decode(&vmMonitoringAlertInfo)

	return vmMonitoringAlertInfo, model.WebStatus{StatusCode: respStatus}
}

// 알람 이벤트 핸들러 목록 조회
// List monitoring alert event handler
func GetMonitoringAlertEventHandlerList() ([]dragonfly.VmMonitoringAlertEventHandlerInfo, model.WebStatus) {
	fmt.Print("#########GetMonitoringAlertEventHandlerList############")
	var originalUrl = "/alert/eventhandlers"
	//var originalUrl = "/alert/eventhandlers?eventType=slack"
	// {{ip}}:{{port}}/dragonfly/alert/eventhandlers?eventType=smtp
	urlParam := util.MappingUrlParameter(originalUrl, nil)

	url := util.DRAGONFLY + urlParam
	// url := util.DRAGONFLY + "/alert/eventhandlers" + "?eventType=smtp"

	resp, err := util.CommonHttp(url, nil, http.MethodGet)
	vmMonitoringAlertEventHandlerInfoList := []dragonfly.VmMonitoringAlertEventHandlerInfo{}
	if err != nil {
		fmt.Println(err)
		return vmMonitoringAlertEventHandlerInfoList, model.WebStatus{StatusCode: 500, Message: err.Error()}
	}

	respBody := resp.Body
	respStatus := resp.StatusCode

	json.NewDecoder(respBody).Decode(&vmMonitoringAlertEventHandlerInfoList)

	return vmMonitoringAlertEventHandlerInfoList, model.WebStatus{StatusCode: respStatus}
}

// 알람 이벤트 핸들러 생성
// Create monitoring alert event-handler
func RegMonitoringAlertEventHandler(vmMonitoringAlertEventHandlerInfoReg *dragonfly.VmMonitoringAlertEventHandlerInfoReg) (*dragonfly.VmMonitoringAlertEventHandlerInfoReg, model.WebStatus) {
	fmt.Println("RegMonitoringAlertEventHandler ************ : ")
	var originalUrl = "/alert/eventhandler"
	// {{ip}}:{{port}}/dragonfly/alert/eventhandler
	urlParam := util.MappingUrlParameter(originalUrl, nil)

	url := util.DRAGONFLY + urlParam
	// url := util.DRAGONFLY + "/alert/eventhandler"

	pbytes, _ := json.Marshal(vmMonitoringAlertEventHandlerInfoReg)
	fmt.Println(string(pbytes))
	resp, err := util.CommonHttp(url, pbytes, http.MethodPost)

	//urlValues, convertErr := util.StructToMapByJson(vmMonitoringAlertEventHandlerInfoReg)
	//if convertErr != nil {
	//	log.Println(convertErr)
	//}
	//
	//fmt.Println(urlValues)
	//resp, err := util.CommonHttpFormData(url, urlValues, http.MethodPost)

	resultVmMonitoringAlertEventHandlerInfoReg := dragonfly.VmMonitoringAlertEventHandlerInfoReg{}
	if err != nil {
		fmt.Println(err)
		return &resultVmMonitoringAlertEventHandlerInfoReg, model.WebStatus{StatusCode: 500, Message: err.Error()}
	}

	respBody := resp.Body
	respStatus := resp.StatusCode

	returnStatus := model.WebStatus{}
	if respStatus != 200 && respStatus != 201 { // 호출은 정상이나, 가져온 결과값이 200, 201아닌 경우 message에 담겨있는 것을 WebStatus에 set
		errorInfo := model.ErrorInfo{}
		json.NewDecoder(respBody).Decode(&errorInfo)
		fmt.Println("respStatus != 200 reason ", errorInfo)
		returnStatus.Message = errorInfo.Message
	} else {
		json.NewDecoder(respBody).Decode(&resultVmMonitoringAlertEventHandlerInfoReg)
		fmt.Println(resultVmMonitoringAlertEventHandlerInfoReg)
	}
	returnStatus.StatusCode = respStatus

	return &resultVmMonitoringAlertEventHandlerInfoReg, returnStatus
}

// 알람 이벤트 핸들러 수정( handlerType=slack)
func PutMonitoringAlertEventHandlerSlack(eventHandlerType string, eventName string, vmMonitoringAlertEventHandlerSlackInfo *dragonfly.EventHandlerOptionSlack) (*dragonfly.VmMonitoringAlertEventHandlerSlackInfo, model.WebStatus) {
	fmt.Println("PutMonitoringAlertEventHandler ************ : ")
	var originalUrl = "/alert/eventhandler/type/:type/event/:name"
	// {{ip}}:{{port}}/dragonfly/alert/eventhandler/type/:type/event/:name
	var paramMapper = make(map[string]string)
	paramMapper[":type"] = eventHandlerType
	paramMapper[":name"] = eventName
	urlParam := util.MappingUrlParameter(originalUrl, paramMapper)

	url := util.DRAGONFLY + urlParam
	// url := util.DRAGONFLY + "/alert/eventhandler/type/" + eventHandlerType + "/event/" + eventName

	pbytes, _ := json.Marshal(vmMonitoringAlertEventHandlerSlackInfo)
	fmt.Println(string(pbytes))
	resp, err := util.CommonHttp(url, pbytes, http.MethodPut)

	resultVmMonitoringAlertEventHandlerInfo := dragonfly.VmMonitoringAlertEventHandlerSlackInfo{}
	if err != nil {
		fmt.Println(err)
		return &resultVmMonitoringAlertEventHandlerInfo, model.WebStatus{StatusCode: 500, Message: err.Error()}
	}

	respBody := resp.Body
	respStatus := resp.StatusCode

	returnStatus := model.WebStatus{}
	if respStatus != 200 && respStatus != 201 { // 호출은 정상이나, 가져온 결과값이 200, 201아닌 경우 message에 담겨있는 것을 WebStatus에 set
		errorInfo := model.ErrorInfo{}
		json.NewDecoder(respBody).Decode(&errorInfo)
		fmt.Println("respStatus != 200 reason ", errorInfo)
		returnStatus.Message = errorInfo.Message
	} else {
		json.NewDecoder(respBody).Decode(&resultVmMonitoringAlertEventHandlerInfo)
		fmt.Println(resultVmMonitoringAlertEventHandlerInfo)
	}
	returnStatus.StatusCode = respStatus

	return &resultVmMonitoringAlertEventHandlerInfo, returnStatus
}

// 알람 이벤트 핸들러 수정( handlerType=smtp)
func PutMonitoringAlertEventHandlerSmtp(eventHandlerType string, eventName string, vmMonitoringAlertEventHandlerInfo *dragonfly.EventHandlerOptionSmtp) (*dragonfly.VmMonitoringAlertEventHandlerSmtpInfo, model.WebStatus) {
	fmt.Println("PutMonitoringAlertEventHandlerSmtp ************ : ")
	var originalUrl = "/alert/eventhandler/type/:type/event/:name"
	// {{ip}}:{{port}}/dragonfly/alert/eventhandler/type/:type/event/:name
	var paramMapper = make(map[string]string)
	paramMapper[":type"] = eventHandlerType
	paramMapper[":name"] = eventName
	urlParam := util.MappingUrlParameter(originalUrl, paramMapper)

	url := util.DRAGONFLY + urlParam
	// url := util.DRAGONFLY + "/alert/eventhandler/type/" + eventHandlerType + "/event/" + eventName
	// {{ip}}:{{port}}/dragonfly/alert/eventhandler/type/:type/event/:name

	pbytes, _ := json.Marshal(vmMonitoringAlertEventHandlerInfo)
	fmt.Println(string(pbytes))
	resp, err := util.CommonHttp(url, pbytes, http.MethodPut)

	resultVmMonitoringAlertEventHandlerInfo := dragonfly.VmMonitoringAlertEventHandlerSmtpInfo{}
	if err != nil {
		fmt.Println(err)
		return &resultVmMonitoringAlertEventHandlerInfo, model.WebStatus{StatusCode: 500, Message: err.Error()}
	}

	respBody := resp.Body
	respStatus := resp.StatusCode

	returnStatus := model.WebStatus{}
	if respStatus != 200 && respStatus != 201 { // 호출은 정상이나, 가져온 결과값이 200, 201아닌 경우 message에 담겨있는 것을 WebStatus에 set
		errorInfo := model.ErrorInfo{}
		json.NewDecoder(respBody).Decode(&errorInfo)
		fmt.Println("respStatus != 200 reason ", errorInfo)
		returnStatus.Message = errorInfo.Message
	} else {
		json.NewDecoder(respBody).Decode(&resultVmMonitoringAlertEventHandlerInfo)
		fmt.Println(resultVmMonitoringAlertEventHandlerInfo)
	}
	returnStatus.StatusCode = respStatus

	return &resultVmMonitoringAlertEventHandlerInfo, returnStatus
}

// 알람 제거
// Delete monitoring alert event-handler
func DelMonitoringAlertEventHandler(eventHandlerType string, eventName string) (io.ReadCloser, model.WebStatus) {
	var originalUrl = "/alert/eventhandler/type/:type/event/:name"
	// {{ip}}:{{port}}/dragonfly/alert/eventhandler/type/:type/event/:name
	var paramMapper = make(map[string]string)
	paramMapper[":type"] = eventHandlerType
	paramMapper[":name"] = eventName
	urlParam := util.MappingUrlParameter(originalUrl, paramMapper)

	url := util.DRAGONFLY + urlParam
	// url := util.DRAGONFLY + "/alert/eventhandler/type/" + eventHandlerType + "/event/" + eventName

	if eventHandlerType == "" {
		return nil, model.WebStatus{StatusCode: 500, Message: "eventHandlerType is required"}
	}
	if eventName == "" {
		return nil, model.WebStatus{StatusCode: 500, Message: "eventName is required"}
	}

	// 경로안에 parameter가 있어 추가 param없이 호출 함.
	resp, err := util.CommonHttp(url, nil, http.MethodDelete)
	if err != nil {
		fmt.Println(err)
		return nil, model.WebStatus{StatusCode: 500, Message: err.Error()}
	}
	// return body, err
	respBody := resp.Body
	respStatus := resp.StatusCode
	return respBody, model.WebStatus{StatusCode: respStatus}
}

// 알람 로그 정보 목록 조회
// List monitoring alert event
func GetMonitoringAlertLogList(taskName string, logLevel string) ([]dragonfly.VmMonitoringAlertLog, model.WebStatus) {
	if logLevel == "" {
		logLevel = "warning"
	}
	var originalUrl = "/alert/task/:task_name/events?level={logLevel}"
	// {{ip}}:{{port}}/dragonfly/alert/task/:task_name/events?level=warning
	var paramMapper = make(map[string]string)
	paramMapper[":task_name"] = taskName
	paramMapper["{logLevel}"] = logLevel
	urlParam := util.MappingUrlParameter(originalUrl, paramMapper)
	//
	url := util.DRAGONFLY + urlParam

	//resp, err := util.CommonHttpFormData(url, nil, http.MethodGet)
	resp, err := util.CommonHttp(url, nil, http.MethodGet)

	if err != nil {
		fmt.Println(err)
		return nil, model.WebStatus{StatusCode: 500, Message: err.Error()}
	}

	respBody := resp.Body
	respStatus := resp.StatusCode

	vmMonitoringAlertLogList := []dragonfly.VmMonitoringAlertLog{}
	json.NewDecoder(respBody).Decode(&vmMonitoringAlertLogList)

	// fmt.Println("check")
	// fmt.Println(vmMonitoringAlertLogList)

	return vmMonitoringAlertLogList, model.WebStatus{StatusCode: respStatus}
}
