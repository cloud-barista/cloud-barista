package mcis

import (

	//"encoding/json"
	"github.com/tidwall/gjson"

	"fmt"
	"io/ioutil"

	//"log"

	//"strings"
	"strconv"

	"bytes"
	"mime/multipart"

	// REST API (echo)
	"net/http"
	

	"sync"

	"github.com/cloud-barista/cb-tumblebug/src/core/common"
)

const MonMetricAll string = "all"
const MonMetricCpu string = "cpu"
const MonMetricCpufreq string = "cpufreq"
const MonMetricMem string = "mem"
const MonMetricNet string = "net"
const MonMetricSwap string = "swap"
const MonMetricDisk string = "disk"
const MonMetricDiskio string = "diskio"


type MonAgentInstallReq struct {
	Ns_id   string `json:"ns_id,omitempty"`
	Mcis_id   string `json:"mcis_id,omitempty"`
	Vm_id     string `json:"vm_id,omitempty"`
	Public_ip string `json:"public_ip,omitempty"`
	User_name string `json:"user_name,omitempty"`
	Ssh_key   string `json:"ssh_key,omitempty"`
	Csp_type   string `json:"cspType,omitempty"`
}
/*
type DfTelegrafMetric struct {
	Name      string                 `json:"name"`
	Tags      map[string]interface{} `json:"tags"`
	Fields    map[string]interface{} `json:"fields"`
	Timestamp int64                  `json:"timestamp"`
	TagInfo   map[string]interface{} `json:"tagInfo"`
}
*/

type MonResultSimple struct {
	Metric   string `json:"metric"`
	VmId   string `json:"vmId"`
	Value   string `json:"value"`
	Err   string `json:"err"`
}

type MonResultSimpleResponse struct {
	NsId   string `json:"nsId"`
	McisId   string `json:"mcisId"`
	McisMonitoring []MonResultSimple `json:"mcisMonitoring"`
}

// Module for checking CB-Dragonfly endpoint (call get config)
func CheckDragonflyEndpoint() error {
	cmd := "/config"

	url := common.DRAGONFLY_REST_URL + cmd
	method := "GET"
  
	client := &http.Client {
	}
	req, err := http.NewRequest(method, url, nil)
  
	if err != nil {
	  fmt.Println(err)
	  return err
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return err
	}
  
	fmt.Println(string(body))
	return nil
}

func CallMonitoringAsync(wg *sync.WaitGroup, nsID string, mcisID string, vmID string, vmIP string, userName string, privateKey string, method string, cmd string, returnResult *[]SshCmdResult) {

	defer wg.Done() //goroutin sync done

	url := common.DRAGONFLY_REST_URL + cmd
	fmt.Println("\n\n[Calling DRAGONFLY] START")
	fmt.Println("url: " + url + " method: " + method)

	tempReq := MonAgentInstallReq{
		Ns_id:   nsID,
		Mcis_id:   mcisID,
		Vm_id:     vmID,
		Public_ip: vmIP,
		User_name: userName,
		Ssh_key:   privateKey,
	}
	fmt.Printf("\n[Request body to CB-DRAGONFLY]\n")
	common.PrintJsonPretty(tempReq)

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	_ = writer.WriteField("ns_id", nsID)
	_ = writer.WriteField("mcis_id", mcisID)
	_ = writer.WriteField("vm_id", vmID)
	_ = writer.WriteField("public_ip", vmIP)
	_ = writer.WriteField("user_name", userName)
	_ = writer.WriteField("ssh_key", privateKey)
	_ = writer.WriteField("cspType", "test")
	err := writer.Close()

	errStr := ""
	if err != nil {
		common.CBLog.Error(err)
		errStr = err.Error()
	}

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	req, err := http.NewRequest(method, url, payload)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	res, err := client.Do(req)

	fmt.Println("Called CB-DRAGONFLY API")
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		common.CBLog.Error(err)
		errStr = err.Error()
	}

	fmt.Println("HTTP Status code " + strconv.Itoa(res.StatusCode))
	switch {
	case res.StatusCode >= 400 || res.StatusCode < 200:
		err := fmt.Errorf(string(body))
		common.CBLog.Error(err)
		errStr = err.Error()
	}

	result := string(body)

	//wg.Done() //goroutin sync done

	vmInfoTmp, _ := GetVmObject(nsID, mcisID, vmID)

	sshResultTmp := SshCmdResult{}
	sshResultTmp.Mcis_id = mcisID
	sshResultTmp.Vm_id = vmID
	sshResultTmp.Vm_ip = vmIP

	if err != nil {
		sshResultTmp.Result = errStr
		sshResultTmp.Err = err
		*returnResult = append(*returnResult, sshResultTmp)
		vmInfoTmp.MonAgentStatus = "failed"
	} else {
		fmt.Println("result " + result)
		sshResultTmp.Result = result
		sshResultTmp.Err = nil
		*returnResult = append(*returnResult, sshResultTmp)
		vmInfoTmp.MonAgentStatus = "installed"
	}
	
	UpdateVmInfo(nsID, mcisID, vmInfoTmp)

}

func InstallMonitorAgentToMcis(nsId string, mcisId string, req *McisCmdReq) (AgentInstallContentWrapper, error) {

	nsId = common.GenId(nsId)
	check, lowerizedName, _ := LowerizeAndCheckMcis(nsId, mcisId)
	mcisId = lowerizedName

	if check == false {
		temp := AgentInstallContentWrapper{}
		err := fmt.Errorf("The mcis " + mcisId + " does not exist.")
		return temp, err
	}

	content := AgentInstallContentWrapper{}

	//install script
	cmd := "/agent/install"

	vmList, err := ListVmId(nsId, mcisId)
	if err != nil {
		common.CBLog.Error(err)
		return content, err
	}

	fmt.Println("[Install agent for each VM]")

	//goroutin sync wg
	var wg sync.WaitGroup

	var resultArray []SshCmdResult

	method := "POST"

	for _, v := range vmList {
		vmObjTmp, _ := GetVmObject(nsId, mcisId, v)
		fmt.Println("MonAgentStatus : " + vmObjTmp.MonAgentStatus)
		
		if(vmObjTmp.MonAgentStatus != "installed"){

			vmId := v
			vmIp := GetVmIp(nsId, mcisId, vmId)
	
			// find vaild username
			userName, sshKey := GetVmSshKey(nsId, mcisId, vmId)
			userNames := []string{SshDefaultUserName01, SshDefaultUserName02, SshDefaultUserName03, SshDefaultUserName04, userName, req.User_name}
			userName = VerifySshUserName(vmIp, userNames, sshKey)
	
			fmt.Println("[CallMonitoringAsync] " + mcisId + "/" + vmId + "(" + vmIp + ")" + "with userName:" + userName)

			wg.Add(1)
			go CallMonitoringAsync(&wg, nsId, mcisId, vmId, vmIp, userName, sshKey, method, cmd, &resultArray)
		}
	}
	wg.Wait() //goroutin sync wg

	for _, v := range resultArray {

		resultTmp := AgentInstallContent{}
		resultTmp.Mcis_id = mcisId
		resultTmp.Vm_id = v.Vm_id
		resultTmp.Vm_ip = v.Vm_ip
		resultTmp.Result = v.Result
		content.Result_array = append(content.Result_array, resultTmp)
		//fmt.Println("result from goroutin " + v)
	}

	//fmt.Printf("%+v\n", content)
	common.PrintJsonPretty(content)

	return content, nil

}

func GetMonitoringData(nsId string, mcisId string, metric string) (MonResultSimpleResponse, error) {

	nsId = common.GenId(nsId)
	check, lowerizedName, _ := LowerizeAndCheckMcis(nsId, mcisId)
	mcisId = lowerizedName

	if check == false {
		temp := MonResultSimpleResponse{}
		err := fmt.Errorf("The mcis " + mcisId + " does not exist.")
		return temp, err
	}

	content := MonResultSimpleResponse{}

	vmList, err := ListVmId(nsId, mcisId)
	if err != nil {
		//common.CBLog.Error(err)
		return content, err
	}

	//goroutin sync wg
	var wg sync.WaitGroup

	var resultArray []MonResultSimple

	method := "GET"

	for _, v := range vmList {
		wg.Add(1)

		vmId := v
		vmIp := GetVmIp(nsId, mcisId, vmId)

		// DF: Get vm on-demand monitoring metric info
		// Path Para: /ns/:ns_id/mcis/:mcis_id/vm/:vm_id/agent_ip/:agent_ip/metric/:metric_name/ondemand-monitoring-info
		cmd := "/ns/" + nsId + "/mcis/" + mcisId + "/vm/" + vmId + "/agent_ip/" + vmIp + "/metric/" + metric + "/ondemand-monitoring-info"
		fmt.Println("[CMD] " + cmd)

		go CallGetMonitoringAsync(&wg, nsId, mcisId, vmId, vmIp, method, metric, cmd, &resultArray)

	}
	wg.Wait() //goroutin sync wg

	content.NsId = nsId
	content.McisId = mcisId
	for _, v := range resultArray {
		content.McisMonitoring = append(content.McisMonitoring, v)
		//fmt.Println("result from goroutin " + v)
	}

	//fmt.Printf("%+v\n", content)
	common.PrintJsonPretty(content)

	return content, nil

}

func CallGetMonitoringAsync(wg *sync.WaitGroup, nsID string, mcisID string, vmID string, vmIP string, method string, metric string, cmd string, returnResult *[]MonResultSimple) {

	defer wg.Done() //goroutin sync done

	url := common.DRAGONFLY_REST_URL + cmd
	fmt.Println("\n\n[Calling DRAGONFLY] START")
	fmt.Println("url: " + url + " method: " + method)

	tempReq := MonAgentInstallReq{
		Mcis_id: mcisID,
		Vm_id:   vmID,
	}
	fmt.Printf("\n[Request body to CB-DRAGONFLY]\n")
	common.PrintJsonPretty(tempReq)

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	req, err := http.NewRequest(method, url, nil)
	errStr := ""
	if err != nil {
		common.CBLog.Error(err)
		errStr = err.Error()
	}

	res, err := client.Do(req)

	fmt.Println("Called CB-DRAGONFLY API")
	if err != nil {
		common.CBLog.Error(err)
		errStr = err.Error()
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		common.CBLog.Error(err)
		errStr = err.Error()
	}

	fmt.Println("HTTP Status code " + strconv.Itoa(res.StatusCode))
	switch {
	case res.StatusCode >= 400 || res.StatusCode < 200:
		err := fmt.Errorf(string(body))
		common.CBLog.Error(err)
		errStr = err.Error()
	}

	result := string(body)

	switch {
		case metric == MonMetricCpu:
			value := gjson.Get(string(body), "values.cpu_utilization")
			result = value.String()
		case metric == MonMetricMem:
			value := gjson.Get(string(body), "values.mem_utilization")
			result = value.String()
		case metric == MonMetricDisk:
			value := gjson.Get(string(body), "values.disk_utilization")
			result = value.String()
		case metric == MonMetricNet:
			value := gjson.Get(string(body), "values.bytes_out")
			result = value.String()
		default:
	}

	//wg.Done() //goroutin sync done

	ResultTmp := MonResultSimple{}
	ResultTmp.VmId = vmID
	ResultTmp.Metric = metric

	if err != nil {
		ResultTmp.Value = errStr
		ResultTmp.Err = err.Error()
		*returnResult = append(*returnResult, ResultTmp)
	} else {
		fmt.Println("result " + result)
		ResultTmp.Value = result
		*returnResult = append(*returnResult, ResultTmp)
	}

}
