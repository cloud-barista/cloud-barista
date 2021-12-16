package service

import (
	"encoding/json"
	"fmt"
	tbmcir "github.com/cloud-barista/cb-webtool/src/model/tumblebug/mcir"

	//tbmcir "github.com/cloud-barista/cb-webtool/src/model/tumblebug/mcir"
	"io"
	"log"
	"net/http"

	// "github.com/davecgh/go-spew/spew"

	// "os"
	// model "github.com/cloud-barista/cb-webtool/src/model"
	"github.com/cloud-barista/cb-webtool/src/model"
	// spider "github.com/cloud-barista/cb-webtool/src/model/spider"
	// "github.com/cloud-barista/cb-webtool/src/model/tumblebug"

	tbcommon "github.com/cloud-barista/cb-webtool/src/model/tumblebug/common"
	// tbmcir "github.com/cloud-barista/cb-webtool/src/model/tumblebug/mcir"
	tbmcis "github.com/cloud-barista/cb-webtool/src/model/tumblebug/mcis"

	webtool "github.com/cloud-barista/cb-webtool/src/model/webtool"

	// "github.com/go-session/echo-session"

	// echosession "github.com/go-session/echo-session"
	"github.com/labstack/echo"

	util "github.com/cloud-barista/cb-webtool/src/util"
)

//var MCISUrl = "http://15.165.16.67:1323"
//var SPiderUrl = "http://15.165.16.67:1024"

// var SpiderUrl = os.Getenv("SPIDER_URL")// util.SPIDER
// var MCISUrl = os.Getenv("TUMBLE_URL")// util.TUMBLEBUG

// MCIS 목록 조회   : option (id, simple, status) 추가할 것.
func GetMcisList(nameSpaceID string, optionParam string) ([]tbmcis.TbMcisInfo, model.WebStatus) {
	var originalUrl = "/ns/{nsId}/mcis"

	var paramMapper = make(map[string]string)
	paramMapper["{nsId}"] = nameSpaceID
	urlParam := util.MappingUrlParameter(originalUrl, paramMapper)

	optionParamVal := ""

	if optionParam == "" {
		optionParam = "status"
	}
	// install, init, cpus, cpum, memR, memW, fioR, fioW, dbR, dbW, rtt, mrtt, clean
	//if optionParam != "" {
	//	optionParamVal = "?option=" + optionParam
	//}
	if optionParam == "all" {
		optionParamVal = "" // all 은 optionParam값이 없는 경우임.
	} else {
		optionParamVal = "?option=" + optionParam
	}

	// url := util.TUMBLEBUG + urlParam
	url := util.TUMBLEBUG + urlParam + optionParamVal
	// url := util.TUMBLEBUG + "/ns/" + nameSpaceID + "/mcis"
	resp, err := util.CommonHttp(url, nil, http.MethodGet)
	// resp, err := util.CommonHttpWithoutParam(url, http.MethodGet)

	if err != nil {
		fmt.Println(err)
		return nil, model.WebStatus{StatusCode: 500, Message: err.Error()}
	}

	respBody := resp.Body
	respStatus := resp.StatusCode

	mcisList := map[string][]tbmcis.TbMcisInfo{}
	returnStatus := model.WebStatus{}

	if respStatus != 200 && respStatus != 201 { // 호출은 정상이나, 가져온 결과값이 200, 201아닌 경우 message에 담겨있는 것을 WebStatus에 set
		//errorInfo := model.ErrorInfo{}
		//json.NewDecoder(respBody).Decode(&errorInfo)
		//fmt.Println("respStatus != 200 reason ", errorInfo)
		//returnStatus.Message = errorInfo.Message

		failResultInfo := tbcommon.TbSimpleMsg{}
		json.NewDecoder(respBody).Decode(&failResultInfo)
		return nil, model.WebStatus{StatusCode: respStatus, Message: failResultInfo.Message}
	}
	json.NewDecoder(respBody).Decode(&mcisList)
	fmt.Println(mcisList["mcis"])

	returnStatus.StatusCode = respStatus
	log.Println(respBody)
	// util.DisplayResponse(resp) // 수신내용 확인

	return mcisList["mcis"], returnStatus
}

func GetMcisListByID(nameSpaceID string) ([]string, model.WebStatus) {
	var originalUrl = "/ns/{nsId}/mcis"

	var paramMapper = make(map[string]string)
	paramMapper["{nsId}"] = nameSpaceID
	urlParam := util.MappingUrlParameter(originalUrl, paramMapper)

	//if optionParam != ""{
	//	urlParam += "?option=" + optionParam
	//}
	urlParam += "?option=id"
	url := util.TUMBLEBUG + urlParam
	// url := util.TUMBLEBUG + "/ns/" + nameSpaceID + "/mcis"
	resp, err := util.CommonHttp(url, nil, http.MethodGet)
	// resp, err := util.CommonHttpWithoutParam(url, http.MethodGet)

	if err != nil {
		fmt.Println(err)
		return nil, model.WebStatus{StatusCode: 500, Message: err.Error()}
	}

	respBody := resp.Body
	respStatus := resp.StatusCode

	mcisList := tbcommon.TbIdList{}
	if respStatus != 200 && respStatus != 201 { // 호출은 정상이나, 가져온 결과값이 200, 201아닌 경우 message에 담겨있는 것을 WebStatus에 set
		//errorInfo := model.ErrorInfo{}
		//json.NewDecoder(respBody).Decode(&errorInfo)
		//fmt.Println("respStatus != 200 reason ", errorInfo)
		//returnStatus.Message = errorInfo.Message

		failResultInfo := tbcommon.TbSimpleMsg{}
		json.NewDecoder(respBody).Decode(&failResultInfo)
		return nil, model.WebStatus{StatusCode: respStatus, Message: failResultInfo.Message}
	}

	json.NewDecoder(respBody).Decode(&mcisList)
	//spew.Dump(body)
	fmt.Println(mcisList.IDList)

	return mcisList.IDList, model.WebStatus{StatusCode: respStatus}
}

func GetMcisListByOption(nameSpaceID string, optionParam string) ([]tbmcis.TbMcisInfo, model.WebStatus) {
	var originalUrl = "/ns/{nsId}/mcis"

	var paramMapper = make(map[string]string)
	paramMapper["{nsId}"] = nameSpaceID
	urlParam := util.MappingUrlParameter(originalUrl, paramMapper)

	if optionParam != "" {
		urlParam += "?option=" + optionParam
	}

	url := util.TUMBLEBUG + urlParam
	// url := util.TUMBLEBUG + "/ns/" + nameSpaceID + "/mcis"
	resp, err := util.CommonHttp(url, nil, http.MethodGet)
	// resp, err := util.CommonHttpWithoutParam(url, http.MethodGet)

	if err != nil {
		fmt.Println(err)
		return nil, model.WebStatus{StatusCode: 500, Message: err.Error()}
	}

	respBody := resp.Body
	respStatus := resp.StatusCode

	mcisList := map[string][]tbmcis.TbMcisInfo{}
	returnStatus := model.WebStatus{}

	if respStatus != 200 && respStatus != 201 { // 호출은 정상이나, 가져온 결과값이 200, 201아닌 경우 message에 담겨있는 것을 WebStatus에 set
		failResultInfo := tbcommon.TbSimpleMsg{}
		json.NewDecoder(respBody).Decode(&failResultInfo)
		return nil, model.WebStatus{StatusCode: respStatus, Message: failResultInfo.Message}
	}
	json.NewDecoder(respBody).Decode(&mcisList)
	fmt.Println(mcisList["mcis"])

	returnStatus.StatusCode = respStatus
	log.Println(respBody)
	// util.DisplayResponse(resp) // 수신내용 확인

	return mcisList["mcis"], returnStatus
}

// 특정 MCIS 조회
// action : status, suspend, resume, reboot, terminate, refine
// option : id, -
// [CONTROL] : common.SimpleMsg
// [DEFAULT] : mcis.TbMcisInfo
// [ID] : common.IdList
// [STATUS] : mcis.McisStatusInfo : status는 swagger에 정의되어 있지 않음. slack에 물어봐야 하나
func GetMcisData(nameSpaceID string, mcisID string) (*tbmcis.TbMcisInfo, model.WebStatus) {
	var originalUrl = "/ns/{nsId}/mcis/{mcisId}"

	var paramMapper = make(map[string]string)
	paramMapper["{nsId}"] = nameSpaceID
	paramMapper["{mcisId}"] = mcisID

	optionParamVal := "?option=default"

	urlParam := util.MappingUrlParameter(originalUrl, paramMapper)

	url := util.TUMBLEBUG + urlParam + optionParamVal
	resp, err := util.CommonHttpWithoutParam(url, http.MethodGet)

	// defer body.Close()
	mcisInfo := tbmcis.TbMcisInfo{}
	if err != nil {
		fmt.Println(err)
		return &mcisInfo, model.WebStatus{StatusCode: 500, Message: err.Error()}
	}
	// util.DisplayResponse(resp) // 수신내용 확인

	respBody := resp.Body
	respStatus := resp.StatusCode

	if respStatus != 200 && respStatus != 201 { // 호출은 정상이나, 가져온 결과값이 200, 201아닌 경우 message에 담겨있는 것을 WebStatus에 set
		failResultInfo := tbcommon.TbSimpleMsg{}
		json.NewDecoder(respBody).Decode(&failResultInfo)
		return &mcisInfo, model.WebStatus{StatusCode: respStatus, Message: failResultInfo.Message}
	}

	json.NewDecoder(respBody).Decode(&mcisInfo)
	fmt.Println(mcisInfo)

	// resultBody, err := ioutil.ReadAll(respBody)
	// if err == nil {
	// 	str := string(resultBody)
	// 	println(str)
	// }
	// pbytes, _ := json.Marshal(respBody)
	// fmt.Println(string(pbytes))

	return &mcisInfo, model.WebStatus{StatusCode: respStatus}
}

func GetMcisDataByStatus(nameSpaceID string, mcisID string, optionParam string) (tbmcis.McisStatusInfo, model.WebStatus) {
	var originalUrl = "/ns/{nsId}/mcis/{mcisId}"

	var paramMapper = make(map[string]string)
	paramMapper["{nsId}"] = nameSpaceID
	paramMapper["{mcisId}"] = mcisID

	optionParamVal := "?option=default"
	// install, init, cpus, cpum, memR, memW, fioR, fioW, dbR, dbW, rtt, mrtt, clean
	if optionParam != "" {
		optionParamVal = "?option=" + optionParam
	}

	urlParam := util.MappingUrlParameter(originalUrl, paramMapper)

	url := util.TUMBLEBUG + urlParam + optionParamVal
	resp, err := util.CommonHttpWithoutParam(url, http.MethodGet)

	// defer body.Close()
	//mcisStatusInfo := tbmcis.McisStatusInfo{}
	//mcisStatusInfo := tbmcis.McisStatusInfo{}
	mcisStatusInfo := map[string]tbmcis.McisStatusInfo{}
	if err != nil {
		fmt.Println(err)
		failStatusInfo := tbmcis.McisStatusInfo{}
		return failStatusInfo, model.WebStatus{StatusCode: 500, Message: err.Error()}
	}
	// util.DisplayResponse(resp) // 수신내용 확인

	respBody := resp.Body
	respStatus := resp.StatusCode

	if respStatus != 200 && respStatus != 201 { // 호출은 정상이나, 가져온 결과값이 200, 201아닌 경우 message에 담겨있는 것을 WebStatus에 set
		failResultInfo := tbcommon.TbSimpleMsg{}
		failStatusInfo := tbmcis.McisStatusInfo{}
		json.NewDecoder(respBody).Decode(&failResultInfo)
		return failStatusInfo, model.WebStatus{StatusCode: respStatus, Message: failResultInfo.Message}
	}

	json.NewDecoder(respBody).Decode(&mcisStatusInfo)
	fmt.Println(mcisStatusInfo)

	// resultBody, err := ioutil.ReadAll(respBody)
	// if err == nil {
	// 	str := string(resultBody)
	// 	println(str)
	// }
	// pbytes, _ := json.Marshal(respBody)
	// fmt.Println(string(pbytes))

	return mcisStatusInfo["status"], model.WebStatus{StatusCode: respStatus}
}

func GetMcisDataByID(nameSpaceID string, mcisID string) (*tbcommon.TbIdList, model.WebStatus) {
	var originalUrl = "/ns/{nsId}/mcis/{mcisId}"

	var paramMapper = make(map[string]string)
	paramMapper["{nsId}"] = nameSpaceID
	paramMapper["{mcisId}"] = mcisID

	optionParamVal := "?option=id"

	urlParam := util.MappingUrlParameter(originalUrl, paramMapper)

	url := util.TUMBLEBUG + urlParam + optionParamVal
	resp, err := util.CommonHttpWithoutParam(url, http.MethodGet)

	// defer body.Close()
	vmIDList := tbcommon.TbIdList{}
	if err != nil {
		fmt.Println(err)
		return &vmIDList, model.WebStatus{StatusCode: 500, Message: err.Error()}
	}
	// util.DisplayResponse(resp) // 수신내용 확인

	respBody := resp.Body
	respStatus := resp.StatusCode

	if respStatus != 200 && respStatus != 201 { // 호출은 정상이나, 가져온 결과값이 200, 201아닌 경우 message에 담겨있는 것을 WebStatus에 set
		failResultInfo := tbcommon.TbSimpleMsg{}
		json.NewDecoder(respBody).Decode(&failResultInfo)
		return &vmIDList, model.WebStatus{StatusCode: respStatus, Message: failResultInfo.Message}
	}

	//json.NewDecoder(respBody).Decode(&mcisStatusInfo)
	//fmt.Println(mcisStatusInfo)

	return &vmIDList, model.WebStatus{StatusCode: respStatus}
}

// MCIS 등록. VM도 함께 등록
func RegMcis(nameSpaceID string, mcisInfo *tbmcis.TbMcisReq) (*tbmcis.TbMcisInfo, model.WebStatus) {
	var originalUrl = "/ns/{nsId}/mcis"

	var paramMapper = make(map[string]string)
	paramMapper["{nsId}"] = nameSpaceID
	urlParam := util.MappingUrlParameter(originalUrl, paramMapper)

	url := util.TUMBLEBUG + urlParam
	// url := util.TUMBLEBUG + "/ns/" + nameSpaceID + "/mcis"

	pbytes, _ := json.Marshal(mcisInfo)
	resp, err := util.CommonHttp(url, pbytes, http.MethodPost)

	returnMcisInfo := tbmcis.TbMcisInfo{}
	returnStatus := model.WebStatus{}

	respBody := resp.Body
	respStatus := resp.StatusCode

	if err != nil {
		fmt.Println(err)
		return &returnMcisInfo, model.WebStatus{StatusCode: 500, Message: err.Error()}
	}

	if respStatus != 200 && respStatus != 201 { // 호출은 정상이나, 가져온 결과값이 200, 201아닌 경우 message에 담겨있는 것을 WebStatus에 set
		//errorInfo := model.ErrorInfo{}
		//json.NewDecoder(respBody).Decode(&errorInfo)
		//fmt.Println("respStatus != 200 reason ", errorInfo)
		//returnStatus.Message = errorInfo.Message
		failResultInfo := tbcommon.TbSimpleMsg{}
		json.NewDecoder(respBody).Decode(&failResultInfo)
		return &returnMcisInfo, model.WebStatus{StatusCode: respStatus, Message: failResultInfo.Message}
	}

	json.NewDecoder(respBody).Decode(&returnMcisInfo)
	fmt.Println(returnMcisInfo)

	returnStatus.StatusCode = respStatus

	// return respBody, respStatusCode
	return &returnMcisInfo, returnStatus
}

// MCIS 등록, channel 이용 : thread이긴 하나 ch를 통해 결과를 받은 후 처리되므로 다를 바가 없음.
// go routine으로 호출하므로 호출결과를 echo-session에 저장  -> web socket으로 front-end 에 전달
func RegMcisByAsync(nameSpaceID string, mcisInfo *tbmcis.TbMcisReq, c echo.Context) {
	var originalUrl = "/ns/{nsId}/mcis"

	var paramMapper = make(map[string]string)
	paramMapper["{nsId}"] = nameSpaceID
	urlParam := util.MappingUrlParameter(originalUrl, paramMapper)

	url := util.TUMBLEBUG + urlParam
	// url := util.TUMBLEBUG + "/ns/" + nameSpaceID + "/mcis"

	pbytes, _ := json.Marshal(mcisInfo)
	resp, err := util.CommonHttp(url, pbytes, http.MethodPost)

	//returnMcisInfo := tbmcis.TbMcisInfo{}
	//returnStatus := model.WebStatus{}

	respBody := resp.Body
	respStatus := resp.StatusCode

	taskKey := nameSpaceID + "||" + "mcis" + "||" + mcisInfo.Name // TODO : 공통 function으로 뺄 것.

	if err != nil {
		fmt.Println(err)
		log.Println("RegMcisByAsync ", err)
		// websocketMessage := websocket.WebSocketMessage{}
		// websocketMessage.Status = "fail"
		// websocketMessage.ProcessTime = time.Now()
		// return &returnMcisInfo, model.WebStatus{StatusCode: 500, Message: err.Error()}

		// websocket으로 전달할 data set
		StoreWebsocketMessage(util.TASK_TYPE_MCIS, taskKey, util.MCIS_LIFECYCLE_CREATE, util.TASK_STATUS_FAIL, c) // session에 작업내용 저장

		// socketDataStore, storedOk := store.Get("socketdata") // 만들고 들어오기 때문에 저장된 게 없을 수 없음..
		// if !storedOk {                                       // 저장된 거 가져올 때 오류면
		// 	log.Println("get stored socketData error ")
		// } else {
		// 	socketDataMap := socketDataStore.(map[int64]websocket.WebSocketMessage)
		// 	// socketDataMap := map[string]websocket.WebSocketMessage{} // 신규로 만들 필요 없음. 이미 있음(by taskKey)

		// 	// socketDataMap[taskKey] = websocketMessage// 해당 객체 갱신.
		// 	// websocketMessage := socketDataMap[taskKey]
		// 	websocketMessage := websocket.WebSocketMessage{}
		// 	websocketMessage.TaskType = "mcis"
		// 	websocketMessage.TaskKey = taskKey
		// 	websocketMessage.Status = "fail"
		// 	websocketMessage.ProcessTime = time.Now()

		// 	socketDataMap[time.Now().UnixNano()] = websocketMessage
		// 	store.Set("socketdata", socketDataMap)

		// 	store.Save()
		// 	log.Println("get stored socketData fail ")
		// }
	}

	if respStatus != 200 && respStatus != 201 { // 호출은 정상이나, 가져온 결과값이 200, 201아닌 경우 message에 담겨있는 것을 WebStatus에 set
		//util.DisplayResponse(resp) // 결과 확인용

		//errorInfo := model.ErrorInfo{}
		//json.NewDecoder(respBody).Decode(&errorInfo)
		//fmt.Println("respStatus != 200 reason ", errorInfo)
		//returnStatus.Message = errorInfo.Message
		failResultInfo := tbcommon.TbSimpleMsg{}
		json.NewDecoder(respBody).Decode(&failResultInfo)
		log.Println("RegMcisByAsync ", failResultInfo)
		StoreWebsocketMessage(util.TASK_TYPE_MCIS, taskKey, util.MCIS_LIFECYCLE_CREATE, util.TASK_STATUS_FAIL, c) // session에 작업내용 저장
	} else {

		//if respStatus != 200 && respStatus != 201 { // 호출은 정상이나, 가져온 결과값이 200, 201아닌 경우 message에 담겨있는 것을 WebStatus에 set
		//	errorInfo := model.ErrorInfo{}
		//	json.NewDecoder(respBody).Decode(&errorInfo)
		//	fmt.Println("respStatus != 200 reason ", errorInfo)
		//	returnStatus.Message = errorInfo.Message
		//} else {
		//	json.NewDecoder(respBody).Decode(&returnMcisInfo)
		//	fmt.Println(returnMcisInfo)
		//}
		//returnStatus.StatusCode = respStatus

		// // TODO : 결과값 설정 확인할 것.
		// socketDataStore, storedOk := store.Get("socketdata")
		// if !storedOk { // 저장된 거 가져올 때 오류면
		// 	log.Println("get stored socketData error 2 ")
		// }
		// socketDataMap := socketDataStore.(map[int64]websocket.WebSocketMessage)
		// websocketMessage := websocket.WebSocketMessage{}
		// // websocketMessage := socketDataMap[time.Now().UnixNano()]
		// // websocketMessage := socketDataMap[taskKey]
		// websocketMessage.TaskType = "mcis"
		// websocketMessage.TaskKey = taskKey
		// websocketMessage.Status = "complete"
		// websocketMessage.ProcessTime = time.Now()
		// store.Set("socketdata", socketDataMap)
		// store.Save()
		// log.Println("set stored socketData complete ", respStatus)
		// websocket으로 전달할 data set

		StoreWebsocketMessage(util.TASK_TYPE_MCIS, taskKey, util.MCIS_LIFECYCLE_CREATE, util.TASK_STATUS_COMPLETE, c) // session에 작업내용 저장
	}
}

// MCIS에 VM 추가 등록
func RegVm(nameSpaceID string, mcisID string, vmInfo *tbmcis.TbVmReq) (*tbmcis.TbVmInfo, model.WebStatus) {
	// var mcisInfoID = mcisInfo.ID // path의 mcisID와 전송되는 parameter의 mcisID 비교용
	// var vmList = mcisInfo.Vms

	// 전송은 vm -> 수신 vm
	returnVmInfo := tbmcis.TbVmInfo{}
	returnStatus := model.WebStatus{}
	fmt.Println("111")

	var originalUrl = "/ns/{nsId}/mcis/{mcisId}/vm" // 1개만 추가할 때
	// if len(vmList) == 0 {
	// 	return nil, model.WebStatus{StatusCode: 500, Message: "There no Vm info"}
	// }
	// fmt.Println("222")
	// // mcisID 변조 체크
	// if mcisID != mcisInfoID {
	// 	return nil, model.WebStatus{StatusCode: 500, Message: "MCIS Info not valid"}
	// }
	fmt.Println("333")
	var paramMapper = make(map[string]string)
	paramMapper["{nsId}"] = nameSpaceID
	paramMapper["{mcisId}"] = mcisID
	urlParam := util.MappingUrlParameter(originalUrl, paramMapper)
	fmt.Println(originalUrl)
	url := util.TUMBLEBUG + urlParam
	// fmt.Println(vmList)
	// pbytes, _ := json.Marshal(vmList)
	fmt.Println(vmInfo)
	pbytes, _ := json.Marshal(vmInfo)
	resp, err := util.CommonHttp(url, pbytes, http.MethodPost)
	fmt.Println("result ", err)
	respBody := resp.Body
	respStatus := resp.StatusCode

	if err != nil {
		fmt.Println(err)
		log.Println("RegVm ", err)
		return &returnVmInfo, model.WebStatus{StatusCode: 500, Message: err.Error()}
	}

	if respStatus != 200 && respStatus != 201 { // 호출은 정상이나, 가져온 결과값이 200, 201아닌 경우 message에 담겨있는 것을 WebStatus에 set
		//util.DisplayResponse(resp) // 결과 확인용

		//errorInfo := model.ErrorInfo{}
		//json.NewDecoder(respBody).Decode(&errorInfo)
		//fmt.Println("respStatus != 200 reason ", errorInfo)
		//returnStatus.Message = errorInfo.Message
		failResultInfo := tbcommon.TbSimpleMsg{}
		json.NewDecoder(respBody).Decode(&failResultInfo)
		log.Println("RegVm ", failResultInfo)
		return &returnVmInfo, model.WebStatus{StatusCode: respStatus, Message: failResultInfo.Message}
	}

	json.NewDecoder(respBody).Decode(&returnVmInfo)
	//fmt.Println(returnVmInfo)

	returnStatus.StatusCode = respStatus
	//fmt.Println(respBody)
	//fmt.Println(respStatus)

	return &returnVmInfo, returnStatus
}

// VM 등록
func AsyncRegVm(nameSpaceID string, mcisID string, vmInfo *tbmcis.TbVmReq, c echo.Context) {
	var originalUrl = "/ns/{nsId}/mcis/{mcisId}/vm" // 1개만 추가할 때

	var paramMapper = make(map[string]string)
	paramMapper["{nsId}"] = nameSpaceID
	paramMapper["{mcisId}"] = mcisID
	urlParam := util.MappingUrlParameter(originalUrl, paramMapper)
	fmt.Println(originalUrl)
	url := util.TUMBLEBUG + urlParam
	// fmt.Println(vmList)
	// pbytes, _ := json.Marshal(vmList)
	fmt.Println(vmInfo)
	pbytes, _ := json.Marshal(vmInfo)
	resp, err := util.CommonHttp(url, pbytes, http.MethodPost)

	respBody := resp.Body
	respStatus := resp.StatusCode
	fmt.Println("AsyncRegVm result ", err)
	log.Println(resp)

	// // TODO : 결과값 설정 확인할 것.
	// store := echosession.FromContext(c)
	// socketDataStore, storedOk := store.Get("socketdata")
	// if !storedOk { // 저장된 거 가져올 때 오류면
	// 	log.Println("get stored socketData error 2 ")
	// }
	// socketDataMap := socketDataStore.(map[string]websocket.WebSocketMessage)
	// websocketMessage := socketDataMap[taskKey]
	// // websocketMessage := socketDataMap[taskKey]
	// websocketMessage.Status = "complete"
	// websocketMessage.ProcessTime = time.Now()
	// store.Set("socketdata", socketDataMap)
	// store.Save()
	// log.Println("set stored socketData complete ", respStatus)

	returnVmInfo := tbmcis.TbVmInfo{}
	taskKey := nameSpaceID + "||" + "vm" + "||" + mcisID + "||" + vmInfo.Name
	if respStatus != 200 && respStatus != 201 { // 호출은 정상이나, 가져온 결과값이 200, 201아닌 경우 message에 담겨있는 것을 WebStatus에 set
		//errorInfo := model.ErrorInfo{}
		//json.NewDecoder(respBody).Decode(&errorInfo)
		//fmt.Println("respStatus != 200 reason ", errorInfo)

		failResultInfo := tbcommon.TbSimpleMsg{}
		json.NewDecoder(respBody).Decode(&failResultInfo)
		log.Println("RegVm ", failResultInfo)
		StoreWebsocketMessage(util.TASK_TYPE_MCIS, taskKey, util.VM_LIFECYCLE_CREATE, util.TASK_STATUS_FAIL, c) // session에 작업내용 저장
	} else { // return이 없으므로 else에서 처리
		json.NewDecoder(respBody).Decode(&returnVmInfo)
		fmt.Println(returnVmInfo)
		StoreWebsocketMessage(util.TASK_TYPE_MCIS, taskKey, util.VM_LIFECYCLE_CREATE, util.TASK_STATUS_COMPLETE, c) // session에 작업내용 저장
	}
}

// MCIS에 VM 추가 등록
func RegVmGroup(nameSpaceID string, mcisID string, vmGroupInfo *tbmcis.TbVmReq) (*tbmcis.TbMcisInfo, model.WebStatus) {
	// var mcisInfoID = mcisInfo.ID
	// var vmList = mcisInfo.Vm

	var originalUrl = "/ns/{nsId}/mcis/{mcisId}/vmgroup" // 여러개 추가할 때
	// if len(vmList) == 0 {
	// 	return nil, model.WebStatus{StatusCode: 500, Message: "There no Vm info"}
	// }

	// mcisID 변조 체크
	// if mcisID != mcisInfoID {
	// 	return nil, model.WebStatus{StatusCode: 500, Message: "MCIS Info not valid"}
	// }

	var paramMapper = make(map[string]string)
	paramMapper["{nsId}"] = nameSpaceID
	paramMapper["{mcisId}"] = mcisID
	urlParam := util.MappingUrlParameter(originalUrl, paramMapper)

	url := util.TUMBLEBUG + urlParam

	pbytes, _ := json.Marshal(vmGroupInfo)
	resp, err := util.CommonHttp(url, pbytes, http.MethodPost)

	// 전송은 vm이나 수신은 mcisInfo (mcis안에 vm목록이 있음.)
	returnMcisInfo := tbmcis.TbMcisInfo{}
	returnStatus := model.WebStatus{}

	if err != nil {
		fmt.Println(err)
		return &returnMcisInfo, model.WebStatus{StatusCode: 500, Message: err.Error()}
	}

	respBody := resp.Body
	respStatus := resp.StatusCode
	returnStatus.StatusCode = respStatus

	if respStatus != 200 && respStatus != 201 { // 호출은 정상이나, 가져온 결과값이 200, 201아닌 경우 message에 담겨있는 것을 WebStatus에 set
		//errorInfo := model.ErrorInfo{}
		//json.NewDecoder(respBody).Decode(&errorInfo)
		//fmt.Println("respStatus != 200 reason ", errorInfo)
		//returnStatus.Message = errorInfo.Message
		failResultInfo := tbcommon.TbSimpleMsg{}
		json.NewDecoder(respBody).Decode(&failResultInfo)
		log.Println("RegVmGroup ", failResultInfo)
		return &returnMcisInfo, model.WebStatus{StatusCode: respStatus, Message: failResultInfo.Message}
	}

	json.NewDecoder(respBody).Decode(&returnMcisInfo)
	fmt.Println(returnMcisInfo)

	// return respBody, respStatusCode
	return &returnMcisInfo, returnStatus
}

// Create MCIS Dynamically from common spec and image
// async 로 만들 지
func RegMcisDynamic(nameSpaceID string, mcisDynamicReq *tbmcis.TbMcisDynamicReq) (*tbmcis.TbMcisInfo, model.WebStatus) {
	var originalUrl = "/ns/{nsId}/mcisDynamic"

	var paramMapper = make(map[string]string)
	paramMapper["{nsId}"] = nameSpaceID
	urlParam := util.MappingUrlParameter(originalUrl, paramMapper)

	url := util.TUMBLEBUG + urlParam

	pbytes, _ := json.Marshal(mcisDynamicReq)
	resp, err := util.CommonHttp(url, pbytes, http.MethodPost)

	returnMcisInfo := tbmcis.TbMcisInfo{}
	returnStatus := model.WebStatus{}

	if err != nil {
		fmt.Println(err)
		return &returnMcisInfo, model.WebStatus{StatusCode: 500, Message: err.Error()}
	}

	respBody := resp.Body
	respStatus := resp.StatusCode
	returnStatus.StatusCode = respStatus

	if respStatus != 200 && respStatus != 201 { // 호출은 정상이나, 가져온 결과값이 200, 201아닌 경우 message에 담겨있는 것을 WebStatus에 set
		failResultInfo := tbcommon.TbSimpleMsg{}
		json.NewDecoder(respBody).Decode(&failResultInfo)
		log.Println("RegMcisDynamic ", failResultInfo)
		return &returnMcisInfo, model.WebStatus{StatusCode: respStatus, Message: failResultInfo.Message}
	}

	json.NewDecoder(respBody).Decode(&returnMcisInfo)
	fmt.Println(returnMcisInfo)

	return &returnMcisInfo, returnStatus
}

// Recommend MCIS plan (filter and priority)
// 실제로는 추천 image 목록
// async 로 만들 지
func RegMcisRecommendVm(nameSpaceID string, mcisDeploymentPlan *tbmcis.DeploymentPlan) ([]tbmcir.TbSpecInfo, model.WebStatus) {
	var originalUrl = "/ns/{nsId}/mcisRecommendVm"

	var paramMapper = make(map[string]string)
	paramMapper["{nsId}"] = nameSpaceID // default는 common
	urlParam := util.MappingUrlParameter(originalUrl, paramMapper)

	url := util.TUMBLEBUG + urlParam

	pbytes, _ := json.Marshal(mcisDeploymentPlan)
	resp, err := util.CommonHttp(url, pbytes, http.MethodPost)

	returnVmSpecs := []tbmcir.TbSpecInfo{}
	returnStatus := model.WebStatus{}

	if err != nil {
		fmt.Println(err)
		return nil, model.WebStatus{StatusCode: 500, Message: err.Error()}
	}

	respBody := resp.Body
	respStatus := resp.StatusCode
	returnStatus.StatusCode = respStatus

	if respStatus != 200 && respStatus != 201 { // 호출은 정상이나, 가져온 결과값이 200, 201아닌 경우 message에 담겨있는 것을 WebStatus에 set
		failResultInfo := tbcommon.TbSimpleMsg{}
		json.NewDecoder(respBody).Decode(&failResultInfo)
		log.Println("RegMcisDynamic ", failResultInfo)
		return returnVmSpecs, model.WebStatus{StatusCode: respStatus, Message: failResultInfo.Message}
	}

	json.NewDecoder(respBody).Decode(&returnVmSpecs)
	fmt.Println(returnVmSpecs)

	return returnVmSpecs, returnStatus
}

////////////////
// func GetVMStatus(vm_name string, connectionConfig string) string {
// 	url := SpiderUrl + "/vmstatus/" + vm_name + "?connection_name=" + connectionConfig
// 	// resp, err := http.Get(url)
// 	// if err != nil {
// 	// 	fmt.Println("request URL : ", url)
// 	// }

// 	// defer resp.Body.Close()
// 	body := HttpGetHandler(url)
// 	defer body.Close()
// 	info := map[string]McisInfo{}
// 	json.NewDecoder(body).Decode(&info)
// 	fmt.Println("VM Status : ", info["status"].Status)
// 	return info["status"].Status

// }

// MCIS 목록에서 mcis 상태별 count map반환
func GetMcisStatusCountMap(mcisInfo tbmcis.TbMcisInfo) map[string]int {
	mcisStatusRunning := 0
	mcisStatusStopped := 0
	mcisStatusTerminated := 0

	// log.Println(" mcisInfo  ", index, mcisInfo)
	mcisStatus := util.GetMcisStatus(mcisInfo.Status)
	if mcisStatus == util.MCIS_STATUS_RUNNING {
		mcisStatusRunning++
	} else if mcisStatus == util.MCIS_STATUS_TERMINATED {
		mcisStatusTerminated++
	} else {
		mcisStatusStopped++
	}
	mcisStatusMap := make(map[string]int)
	mcisStatusMap["RUNNING"] = mcisStatusRunning
	mcisStatusMap["STOPPED"] = mcisStatusStopped
	mcisStatusMap["TERMINATED"] = mcisStatusTerminated
	mcisStatusMap["TOTAL"] = mcisStatusRunning + mcisStatusStopped + mcisStatusTerminated
	// mcisStatusTotalMap[mcisInfo.ID] = mcisStatusMap

	return mcisStatusMap
}

// MCIS의 vm별 statun와 vm 상태별 count
// key는 vmID + vmName, value는 vmStatus
func GetSimpleVmWithStatusCountMap(mcisInfo tbmcis.TbMcisInfo) ([]webtool.VmSimpleInfo, map[string]int) {
	// log.Println(" mcisInfo  ", index, mcisInfo)
	// vmStatusMap := make(map[string]int)
	// vmStatusMap := map[string]string{} // vmName : vmStatus

	vmStatusCountMap := map[string]int{}
	totalVmStatusCount := 0
	vmList := mcisInfo.Vm
	var vmSimpleList []webtool.VmSimpleInfo
	for vmIndex, vmInfo := range vmList {
		// log.Println(" vmInfo  ", vmIndex, vmInfo)
		vmStatus := util.GetVmStatus(vmInfo.Status) // lowercase로 변환

		locationInfo := vmInfo.Location
		vmLatitude := locationInfo.Latitude
		vmLongitude := locationInfo.Longitude

		log.Println(locationInfo)
		//
		vmSimpleObj := webtool.VmSimpleInfo{
			VmIndex:   vmIndex + 1,
			VmID:      vmInfo.ID,
			VmName:    vmInfo.Name,
			VmStatus:  vmStatus,
			Latitude:  vmLatitude,
			Longitude: vmLongitude,

			// export용 추가
			VmConnectionName: vmInfo.ConnectionName,
			VmDescription:    vmInfo.Description,
			VmImageId:        vmInfo.ImageID,
			VmLabel:          vmInfo.Label,
			// TODO : securityGroupId는 []string 이므로 변환작업 필요
			VmSecurityGroupIds: vmInfo.SecurityGroupIDs, //"securityGroupIIds": [		{		  "nameId": "string",		  "systemId": "string"		}	  ],
			VmSpecId:           vmInfo.SpecID,
			VmSshKeyId:         vmInfo.SshKeyID,
			VmSubnetId:         vmInfo.SubnetID,
			VmVnetId:           vmInfo.VNetID,
			VmGroupSize:        1, //? 는 없는데.. vmGroupId만 있는데... 1로 기본 setting
			VmUserAccount:      vmInfo.VmUserAccount,
			VmUserPassword:     vmInfo.VmUserPassword,
		}

		vmSimpleList = append(vmSimpleList, vmSimpleObj)

		log.Println("vmStatus " + vmStatus + ", Status " + vmInfo.Status)
		log.Println(vmInfo.SecurityGroupIDs)
		vmStatusCount := 0
		val, exists := vmStatusCountMap[vmStatus]
		if exists {
			vmStatusCount = val + 1
			totalVmStatusCount += 1
		} else {
			vmStatusCount = 1
			totalVmStatusCount += 1
		}
		vmStatusCountMap[vmStatus] = vmStatusCount
	}

	vmStatusCountMap["TOTAL"] = totalVmStatusCount
	log.Println(vmStatusCountMap)
	// vmStatusCountMap := make(map[string]int)
	// UI에서 사칙연산이 되지 않아 controller에서 계산한 뒤 넘겨 줌. 아나면 function을 정의해서 넘겨야 함
	// vmStatusCountMap[util.VM_STATUS_RUNNING] = vmStatusRunning
	// vmStatusCountMap[util.VM_STATUS_RESUMING] = vmStatusResuming
	// vmStatusCountMap[util.VM_STATUS_INCLUDE] = vmStatusInclude
	// vmStatusCountMap[util.VM_STATUS_SUSPENDED] = vmStatusSuspended
	// vmStatusCountMap[util.VM_STATUS_TERMINATED] = vmStatusTerminated
	// vmStatusCountMap[util.VM_STATUS_UNDEFINED] = vmStatusUndefined
	// vmStatusCountMap[util.VM_STATUS_PARTIAL] = vmStatusPartial
	// vmStatusCountMap[util.VM_STATUS_ETC] = vmStatusEtc
	// log.Println("mcisInfo.ID  ", mcisInfo.ID)
	// mcisIdArr[mcisIndex] = mcisInfo.ID	// 바로 넣으면 Runtime Error구만..
	// vmStatusArr[mcisIndex] = vmStatusCountMap

	// UI에서는 3가지로 통합하여 봄
	// vmStatusCountMap["RUNNING"] = vmStatusRunning
	// vmStatusCountMap["STOPPED"] = vmStatusInclude + vmStatusSuspended + vmStatusUndefined + vmStatusPartial + vmStatusEtc
	// vmStatusCountMap["TERMINATED"] = vmStatusTerminated
	// vmStatusTotalMap[mcisInfo.ID] = vmStatusCountMap
	// vmIdArr = append(vmIdArr, vmInfo.ID)
	// vmStatusArr = append(vmStatusArr, vmStatusCountMap)

	// log.Println("mcisIndex  ", mcisIndex)

	// vmStatusCountMap["RUNNING"] = vmStatusRunning
	// vmStatusCountMap["STOPPED"] = vmStatusInclude + vmStatusSuspended + vmStatusUndefined + vmStatusPartial + vmStatusEtc
	// vmStatusCountMap["TERMINATED"] = vmStatusTerminated
	// vmStatusCountMap["TOTAL"] = vmStatusCountMap["RUNNING"] + vmStatusCountMap["STOPPED"] + vmStatusCountMap["TERMINATED"]

	return vmSimpleList, vmStatusCountMap

}

// MCIS별 connection count
func GetVmConnectionCountMap(mcisInfo tbmcis.TbMcisInfo) map[string]int {
	// connectionCountTotal := 0
	// connectionCountByMcis := 0
	// vmCountTotal := 0
	// vmRunningCountByMcis := 0
	// vmStoppedCountByMcis := 0
	// vmTerminatedCountByMcis := 0
	// vmStatusUndefined := 0
	// vmStatusPartial := 0
	// vmStatusEtc := 0
	// vmStatusTerminated := 0

	// log.Println(" mcisInfo  ", index, mcisInfo)
	// vmList := mcisInfo.VMs
	// for vmIndex, vmInfo := range vmList {
	// 	// log.Println(" vmInfo  ", vmIndex, vmInfo)
	// 	vmConnection := util.GetVmConnectionName(vmInfo.ConnectionName)

	// }
	vmStatusMap := make(map[string]int)
	// UI에서 사칙연산이 되지 않아 controller에서 계산한 뒤 넘겨 줌.
	// vmStatusMap[util.VM_STATUS_RUNNING] = vmStatusRunning
	// vmStatusMap[util.VM_STATUS_RESUMING] = vmStatusResuming
	// vmStatusMap[util.VM_STATUS_INCLUDE] = vmStatusInclude
	// vmStatusMap[util.VM_STATUS_SUSPENDED] = vmStatusSuspended
	// vmStatusMap[util.VM_STATUS_TERMINATED] = vmStatusTerminated
	// vmStatusMap[util.VM_STATUS_UNDEFINED] = vmStatusUndefined
	// vmStatusMap[util.VM_STATUS_PARTIAL] = vmStatusPartial
	// vmStatusMap[util.VM_STATUS_ETC] = vmStatusEtc
	// log.Println("mcisInfo.ID  ", mcisInfo.ID)
	// mcisIdArr[mcisIndex] = mcisInfo.ID	// 바로 넣으면 Runtime Error구만..
	// vmStatusArr[mcisIndex] = vmStatusMap

	// UI에서는 3가지로 통합하여 봄
	// vmStatusMap["RUNNING"] = vmStatusRunning
	// vmStatusMap["STOPPED"] = vmStatusInclude + vmStatusSuspended + vmStatusUndefined + vmStatusPartial + vmStatusEtc
	// vmStatusMap["TERMINATED"] = vmStatusTerminated
	// vmStatusTotalMap[mcisInfo.ID] = vmStatusMap
	// vmIdArr = append(vmIdArr, vmInfo.ID)
	// vmStatusArr = append(vmStatusArr, vmStatusMap)

	// log.Println("mcisIndex  ", mcisIndex)

	// vmStatusMap := make(map[string]int)
	// vmStatusMap["RUNNING"] = vmStatusRunning
	// vmStatusMap["STOPPED"] = vmStatusInclude + vmStatusSuspended + vmStatusUndefined + vmStatusPartial + vmStatusEtc
	// vmStatusMap["TERMINATED"] = vmStatusTerminated

	return vmStatusMap

}

// 해당 MCIS의 VM 연결 수
func GetVmConnectionCountByMcis(mcisInfo tbmcis.TbMcisInfo) map[string]int {
	// log.Println(" mcisInfo  ", index, mcisInfo)
	vmList := mcisInfo.Vm
	// mcisConnectionCountMap := make(map[string]int)
	mcisConnectionCountMap := map[string]int{}

	totalConnectionCount := 0
	log.Println("GetVMConnectionCountByMcis map length ", len(mcisConnectionCountMap))
	for _, vmInfo := range vmList {
		// log.Println(" vmInfo  ", vmIndex, vmInfo)
		locationInfo := vmInfo.Location
		// cloudType := locationInfo.CloudType // CloudConnection
		providerCount := 0
		val, exists := mcisConnectionCountMap[util.GetProviderName(locationInfo.CloudType)]
		if exists {
			providerCount = val + 1
			// totalConnectionCount += 1 // 이미 있는 경우에는 count추가필요없음
		} else {
			providerCount = 1
			totalConnectionCount += 1
		}
		// log.Println("GetProviderName ", locationInfo.CloudType)
		mcisConnectionCountMap[util.GetProviderName(locationInfo.CloudType)] = providerCount
	}
	// log.Println("GetVMConnectionCountByMcis map length ", len(mcisConnectionCountMap))
	// log.Println("GetVMConnectionCountByMcis map ", mcisConnectionCountMap)
	return mcisConnectionCountMap
}

// MCIS의 특정 VM 조회
// action : status, suspend, resume, reboot, terminate
func GetVMofMcisData(nameSpaceID string, mcisID string, vmID string) (*tbmcis.TbVmInfo, model.WebStatus) {
	var originalUrl = "/ns/{nsId}/mcis/{mcisId}/vm/{vmId}"

	var paramMapper = make(map[string]string)
	paramMapper["{nsId}"] = nameSpaceID
	paramMapper["{mcisId}"] = mcisID
	paramMapper["{vmId}"] = vmID
	urlParam := util.MappingUrlParameter(originalUrl, paramMapper)

	url := util.TUMBLEBUG + urlParam
	// url := util.TUMBLEBUG + "/ns/" + nameSpaceID + "/mcis/" + mcisID + "/vm/" + vmID

	// resp, err := util.CommonHttp(url, nil, http.MethodGet)
	resp, err := util.CommonHttpWithoutParam(url, http.MethodGet)

	// defer body.Close()
	vmInfo := tbmcis.TbVmInfo{}
	if err != nil {
		fmt.Println(err)
		return &vmInfo, model.WebStatus{StatusCode: 500, Message: err.Error()}
	}

	respBody := resp.Body
	respStatus := resp.StatusCode

	if respStatus != 200 && respStatus != 201 { // 호출은 정상이나, 가져온 결과값이 200, 201아닌 경우 message에 담겨있는 것을 WebStatus에 set
		failResultInfo := tbcommon.TbSimpleMsg{}
		json.NewDecoder(respBody).Decode(&failResultInfo)
		log.Println("GetVMofMcisData ", failResultInfo)
		return &vmInfo, model.WebStatus{StatusCode: respStatus, Message: failResultInfo.Message}
	}

	json.NewDecoder(respBody).Decode(&vmInfo)
	fmt.Println("respStatus = ", respStatus)
	fmt.Println(vmInfo)

	return &vmInfo, model.WebStatus{StatusCode: respStatus}
}

// MCIS의 Status변경
// LifeCycle 의 경우 요청에 대한 응답이 바로 오므로 asyncMethod를 따로 만들지 않음. 응답시간이 오래걸리는 경우 syncXXX 를 만들고 echo 를 같이 넘겨 결과 처리하도록 해야 함.
func McisLifeCycle(mcisLifeCycle *webtool.McisLifeCycle) (*webtool.McisLifeCycle, model.WebStatus) {
	nameSpaceID := mcisLifeCycle.NameSpaceID
	mcisID := mcisLifeCycle.McisID
	lifeCycleType := mcisLifeCycle.LifeCycleType

	var originalUrl = "/ns/{nsId}/control/mcis/{mcisId}?action={type}"

	var paramMapper = make(map[string]string)
	paramMapper["{nsId}"] = nameSpaceID
	paramMapper["{mcisId}"] = mcisID
	paramMapper["{type}"] = lifeCycleType
	urlParam := util.MappingUrlParameter(originalUrl, paramMapper)

	url := util.TUMBLEBUG + urlParam
	// url := util.TUMBLEBUG + "/ns/" + mcisLifeCycle.NameSpaceID + "/mcis/" + mcisLifeCycle.McisID + "?action=" + mcisLifeCycle.LifeCycleType
	//// var url = CommonURL+"/ns/"+nameSpace+"/mcis/"+mcis_id+"?action="+type
	//pbytes, _ := json.Marshal(mcisLifeCycle)
	//resp, err := util.CommonHttp(url, pbytes, http.MethodGet) // POST로 받기는 했으나 실제로는 Get으로 날아감.
	resp, err := util.CommonHttpWithoutParam(url, http.MethodGet)
	resultMcisLifeCycle := webtool.McisLifeCycle{}
	if err != nil {
		fmt.Println("McisLifeCycle err")
		fmt.Println(err)
		return &resultMcisLifeCycle, model.WebStatus{StatusCode: 500, Message: err.Error()}
	}

	respBody := resp.Body
	respStatus := resp.StatusCode
	// respStatus := resp.Status
	// log.Println("respStatusCode = ", respStatusCode)
	// log.Println("respStatus = ", respStatus)

	// 응답에 생성한 객체값이 옴

	if respStatus != 200 && respStatus != 201 {
		// statusInfo := model.WebStatus{}
		// fmt.Println("McisLifeCycle respStatus ", respStatus)
		// fmt.Println(respBody)
		// json.NewDecoder(respBody).Decode(statusInfo)
		// fmt.Println(statusInfo)
		// fmt.Println(statusInfo.Message)

		//errorInfo := model.ErrorInfo{}
		//json.NewDecoder(respBody).Decode(&errorInfo)
		//fmt.Println("respStatus != 200 reason ", errorInfo)

		failResultInfo := tbcommon.TbSimpleMsg{}
		json.NewDecoder(respBody).Decode(&failResultInfo)
		log.Println("McisLifeCycle ", failResultInfo)
		return &resultMcisLifeCycle, model.WebStatus{StatusCode: respStatus, Message: failResultInfo.Message}

		//return &resultMcisLifeCycle, model.WebStatus{StatusCode: respStatus, Message: errorInfo.Message}
	}
	// return body, err
	// respBody := resp.Body
	// respStatus := resp.StatusCode
	// return respBody, respStatus

	json.NewDecoder(respBody).Decode(resultMcisLifeCycle)
	fmt.Println(resultMcisLifeCycle)

	return &resultMcisLifeCycle, model.WebStatus{StatusCode: respStatus}

}
func McisLifeCycleByAsync(mcisLifeCycle *webtool.McisLifeCycle, c echo.Context) {
	nameSpaceID := mcisLifeCycle.NameSpaceID
	mcisID := mcisLifeCycle.McisID
	lifeCycleType := mcisLifeCycle.LifeCycleType

	var originalUrl = "/ns/{nsId}/control/mcis/{mcisId}?action={type}"

	var paramMapper = make(map[string]string)
	paramMapper["{nsId}"] = nameSpaceID
	paramMapper["{mcisId}"] = mcisID
	paramMapper["{type}"] = lifeCycleType
	urlParam := util.MappingUrlParameter(originalUrl, paramMapper)

	url := util.TUMBLEBUG + urlParam
	resp, err := util.CommonHttpWithoutParam(url, http.MethodGet)

	taskKey := nameSpaceID + "||" + "mcis" + "||" + mcisLifeCycle.McisID
	////////////
	if err != nil {
		fmt.Println("McisLifeCycle err")
		fmt.Println(err)

		// websocket으로 전달할 data set
		StoreWebsocketMessage(util.TASK_TYPE_MCIS, taskKey, mcisLifeCycle.LifeCycleType, util.TASK_STATUS_FAIL, c) // session에 작업내용 저장
	}

	respBody := resp.Body
	respStatus := resp.StatusCode

	if respStatus != 200 && respStatus != 201 {
		failResultInfo := tbcommon.TbSimpleMsg{}
		json.NewDecoder(respBody).Decode(&failResultInfo)
		log.Println("McisLifeCycle ", failResultInfo)
		StoreWebsocketMessage(util.TASK_TYPE_MCIS, taskKey, mcisLifeCycle.LifeCycleType, util.TASK_STATUS_FAIL, c) // session에 작업내용 저장
	} else {
		resultMcisLifeCycle := webtool.McisLifeCycle{}
		json.NewDecoder(respBody).Decode(resultMcisLifeCycle)
		fmt.Println(resultMcisLifeCycle)
		StoreWebsocketMessage(util.TASK_TYPE_MCIS, taskKey, mcisLifeCycle.LifeCycleType, util.TASK_STATUS_COMPLETE, c) // session에 작업내용 저장
	}
}

// MCIS의 VM Status변경 : 요청에 대한 응답이 바로 오므로 async 만들지 않음
func McisVmLifeCycle(vmLifeCycle *webtool.VmLifeCycle) (*webtool.VmLifeCycle, model.WebStatus) {
	var originalUrl = "/ns/{nsId}/control/mcis/{mcisId}/vm/{vmId}?action={type}"

	var paramMapper = make(map[string]string)
	paramMapper["{nsId}"] = vmLifeCycle.NameSpaceID
	paramMapper["{mcisId}"] = vmLifeCycle.McisID
	paramMapper["{vmId}"] = vmLifeCycle.VmID
	paramMapper["{type}"] = vmLifeCycle.LifeCycleType
	urlParam := util.MappingUrlParameter(originalUrl, paramMapper)

	url := util.TUMBLEBUG + urlParam
	// url := util.TUMBLEBUG + "/ns/" + vmLifeCycle.NameSpaceID + "/mcis/" + vmLifeCycle.McisID + "/vm/" + vmLifeCycle.VmID + "?action=" + vmLifeCycle.LifeCycleType
	///url = CommonURL+"/ns/"+nameSpace+"/mcis/"+mcis_id+"/vm/"+vm_id+"?action="+type
	//pbytes, _ := json.Marshal(vmLifeCycle)
	//resp, err := util.CommonHttp(url, pbytes, http.MethodGet) // POST로 받기는 했으나 실제로는 Get으로 날아감.
	resp, err := util.CommonHttpWithoutParam(url, http.MethodGet)
	resultVmLifeCycle := webtool.VmLifeCycle{}
	if err != nil {
		fmt.Println(err)
		return &resultVmLifeCycle, model.WebStatus{StatusCode: 500, Message: err.Error()}
	}

	respBody := resp.Body
	respStatus := resp.StatusCode

	if respStatus != 200 && respStatus != 201 {
		failResultInfo := tbcommon.TbSimpleMsg{}
		json.NewDecoder(respBody).Decode(&failResultInfo)
		log.Println("McisVmLifeCycle ", failResultInfo)
		return &resultVmLifeCycle, model.WebStatus{StatusCode: respStatus, Message: failResultInfo.Message}
	}

	// 응답에 생성한 객체값이 옴
	json.NewDecoder(respBody).Decode(resultVmLifeCycle)
	fmt.Println(resultVmLifeCycle)

	return &resultVmLifeCycle, model.WebStatus{StatusCode: respStatus}
}

func McisVmLifeCycleByAsync(vmLifeCycle *webtool.VmLifeCycle, c echo.Context) {
	var originalUrl = "/ns/{nsId}/control/mcis/{mcisId}/vm/{vmId}?action={type}"

	var paramMapper = make(map[string]string)
	paramMapper["{nsId}"] = vmLifeCycle.NameSpaceID
	paramMapper["{mcisId}"] = vmLifeCycle.McisID
	paramMapper["{vmId}"] = vmLifeCycle.VmID
	paramMapper["{type}"] = vmLifeCycle.LifeCycleType
	urlParam := util.MappingUrlParameter(originalUrl, paramMapper)

	url := util.TUMBLEBUG + urlParam
	resp, err := util.CommonHttpWithoutParam(url, http.MethodGet)
	resultVmLifeCycle := webtool.VmLifeCycle{}

	taskKey := vmLifeCycle.NameSpaceID + "||" + "vm" + "||" + vmLifeCycle.McisID + "||" + vmLifeCycle.VmID

	if err != nil {
		fmt.Println(err)
		StoreWebsocketMessage(util.TASK_TYPE_VM, taskKey, vmLifeCycle.LifeCycleType, util.TASK_STATUS_FAIL, c) // session에 작업내용 저장
	}

	respBody := resp.Body
	respStatus := resp.StatusCode

	if respStatus != 200 && respStatus != 201 {
		failResultInfo := tbcommon.TbSimpleMsg{}
		json.NewDecoder(respBody).Decode(&failResultInfo)
		log.Println("McisVmLifeCycle ", failResultInfo)
		StoreWebsocketMessage(util.TASK_TYPE_VM, taskKey, vmLifeCycle.LifeCycleType, util.TASK_STATUS_FAIL, c) // session에 작업내용 저장
	}

	// 응답에 생성한 객체값이 옴
	json.NewDecoder(respBody).Decode(resultVmLifeCycle)
	fmt.Println(resultVmLifeCycle)
	StoreWebsocketMessage(util.TASK_TYPE_VM, taskKey, vmLifeCycle.LifeCycleType, util.TASK_STATUS_COMPLETE, c) // session에 작업내용 저장
	//return &resultVmLifeCycle, model.WebStatus{StatusCode: respStatus}
}

// 벤치마크?? MCIS 조회. 근데 왜 결과는 resultarray지?
// TODO : 여러개 return되면 method이름을 xxxData -> xxxList 로 바꿀 것
func GetBenchmarkMcisData(nameSpaceID string, mcisID string, hostIp string, optionParam string) ([]tbmcis.BenchmarkInfo, model.WebStatus) {
	var originalUrl = "/ns/{nsId}/benchmark/mcis/{mcisId}"

	var paramMapper = make(map[string]string)
	paramMapper["{nsId}"] = nameSpaceID
	paramMapper["{mcisId}"] = mcisID
	urlParam := util.MappingUrlParameter(originalUrl, paramMapper)

	optionParamVal := ""
	// install, init, cpus, cpum, memR, memW, fioR, fioW, dbR, dbW, rtt, mrtt, clean
	if optionParam != "" {
		optionParamVal = "?option=" + optionParam
	}

	// url := util.TUMBLEBUG + urlParam + "?option=" + optionParam
	url := util.TUMBLEBUG + urlParam + optionParamVal
	// url := util.TUMBLEBUG + "/ns/" + nameSpaceID + "/benchmark/mcis/" + mcisID
	// /ns/{nsId}/benchmark/mcis/{mcisId}
	pbytes, _ := json.Marshal(hostIp)
	resp, err := util.CommonHttp(url, pbytes, http.MethodPost) // 조회이나 POST로 호출

	respBody := resp.Body
	respStatus := resp.StatusCode
	// defer body.Close()
	//resultBenchmarkInfos := []tbmcis.BenchmarkInfo{}
	resultBenchmarkInfos := map[string][]tbmcis.BenchmarkInfo{}
	if err != nil {
		fmt.Println(err)
		return nil, model.WebStatus{StatusCode: 500, Message: err.Error()}
	}
	// util.DisplayResponse(resp) // 수신내용 확인

	if respStatus != 200 && respStatus != 201 {
		failResultInfo := tbcommon.TbSimpleMsg{}
		json.NewDecoder(respBody).Decode(&failResultInfo)
		log.Println("GetBenchmarkMcisData ", failResultInfo)
		return nil, model.WebStatus{StatusCode: respStatus, Message: failResultInfo.Message}
	}

	json.NewDecoder(respBody).Decode(&resultBenchmarkInfos)
	fmt.Println(resultBenchmarkInfos)
	//resultarray
	return resultBenchmarkInfos["resultarray"], model.WebStatus{StatusCode: respStatus}
}

// List all MCISs
func GetBenchmarkAllMcisList(nameSpaceID string, mcisID string, hostIp string) ([]tbmcis.BenchmarkInfo, model.WebStatus) {
	var originalUrl = "/ns/{nsId}/benchmarkAll/mcis/{mcisId}"

	var paramMapper = make(map[string]string)
	paramMapper["{nsId}"] = nameSpaceID
	paramMapper["{mcisId}"] = mcisID
	urlParam := util.MappingUrlParameter(originalUrl, paramMapper)

	url := util.TUMBLEBUG + urlParam
	// url := util.TUMBLEBUG + "/ns/" + nameSpaceID + "/benchmark/mcis/" + mcisID

	pbytes, _ := json.Marshal(hostIp)
	resp, err := util.CommonHttp(url, pbytes, http.MethodPost) // 조회이나 POST로 호출

	respBody := resp.Body
	respStatus := resp.StatusCode

	// defer body.Close()
	//resultBenchmarkInfos := []tbmcis.BenchmarkInfo{}
	//if err != nil {
	//	fmt.Println(err)
	//	return &resultBenchmarkInfos, model.WebStatus{StatusCode: 500, Message: err.Error()}
	//}
	resultBenchmarkInfos := map[string][]tbmcis.BenchmarkInfo{}
	if err != nil {
		fmt.Println(err)
		return nil, model.WebStatus{StatusCode: 500, Message: err.Error()}
	}

	if respStatus != 200 && respStatus != 201 {
		failResultInfo := tbcommon.TbSimpleMsg{}
		json.NewDecoder(respBody).Decode(&failResultInfo)
		log.Println("GetBenchmarkAllMcisList ", failResultInfo)
		return nil, model.WebStatus{StatusCode: respStatus, Message: failResultInfo.Message}
	}

	json.NewDecoder(respBody).Decode(&resultBenchmarkInfos)
	fmt.Println(resultBenchmarkInfos)

	return resultBenchmarkInfos["resultarray"], model.WebStatus{StatusCode: respStatus}
}

// MCIS에 명령 내리기
func CommandMcis(nameSpaceID string, mcisID string, mcisCommandInfo *tbmcis.McisCmdReq) (tbmcis.RestPostCmdMcisResponseWrapper, model.WebStatus) {
	// webStatus := model.WebStatus{}
	resultInfo := tbmcis.RestPostCmdMcisResponseWrapper{}

	var originalUrl = "/ns/{nsId}/cmd/mcis/{mcisId}"

	var paramMapper = make(map[string]string)
	paramMapper["{nsId}"] = nameSpaceID
	paramMapper["{mcisId}"] = mcisID
	urlParam := util.MappingUrlParameter(originalUrl, paramMapper)

	url := util.TUMBLEBUG + urlParam

	pbytes, _ := json.Marshal(mcisCommandInfo)
	resp, err := util.CommonHttp(url, pbytes, http.MethodPost)

	if err != nil {
		fmt.Println(err)
		return resultInfo, model.WebStatus{StatusCode: 500, Message: err.Error()}
	}
	fmt.Println("resp : ", resp)

	// return body, err
	respBody := resp.Body
	respStatus := resp.StatusCode
	//resultInfo := model.ResultInfo{}
	log.Println("ResultStatusCode : ", respStatus)

	// 실패시 Message에 성공시 Result에 string으로 담겨 온다.
	if respStatus != 200 && respStatus != 201 {
		failResult := tbcommon.TbSimpleMsg{}
		json.NewDecoder(respBody).Decode(&failResult)
		log.Println("ResultMessage : " + failResult.Message)

		// return model.WebStatus{}, model.WebStatus{StatusCode: respStatus, Message: failResult.Message}
		return resultInfo, model.WebStatus{StatusCode: respStatus, Message: failResult.Message}
	}

	log.Println(respBody)
	// spew.Dump(respBody)

	json.NewDecoder(respBody).Decode(&resultInfo)
	log.Println(resultInfo)

	// webStatus.StatusCode = respStatus
	return resultInfo, model.WebStatus{StatusCode: respStatus}
}

// 특정 VM에 명령내리기
func CommandVmOfMcis(nameSpaceID string, mcisID string, vmID string, mcisCommandInfo *tbmcis.McisCmdReq) (tbmcis.RestPostCmdMcisVmResponse, model.WebStatus) {
	// webStatus := model.WebStatus{}
	resultInfo := tbmcis.RestPostCmdMcisVmResponse{}

	var originalUrl = "/ns/{nsId}/cmd/mcis/{mcisId}/vm/{vmId}"

	var paramMapper = make(map[string]string)
	paramMapper["{nsId}"] = nameSpaceID
	paramMapper["{mcisId}"] = mcisID
	paramMapper["{vmId}"] = vmID
	urlParam := util.MappingUrlParameter(originalUrl, paramMapper)

	url := util.TUMBLEBUG + urlParam
	// url := util.TUMBLEBUG + "/ns/" + nameSpaceID + "/mcis/" + mcisID + "/vm/" + vmID

	pbytes, _ := json.Marshal(mcisCommandInfo)
	resp, err := util.CommonHttp(url, pbytes, http.MethodPost)

	if err != nil {
		fmt.Println(err)
		return resultInfo, model.WebStatus{StatusCode: 500, Message: err.Error()}
	}
	fmt.Println("resp : ", resp)

	// return body, err
	respBody := resp.Body
	respStatus := resp.StatusCode
	// resultInfo := model.ResultInfo{}

	log.Println(respBody)
	// spew.Dump(respBody)

	log.Println("ResultStatusCode : ", respStatus)

	// 실패시 Message에 성공시 Result에 string으로 담겨 온다.
	if respStatus != 200 && respStatus != 201 {
		// log.Println("ResultMessage : " + resultInfo.Message)
		// return model.WebStatus{}, model.WebStatus{StatusCode: respStatus, Message: resultInfo.Message}
		failResult := tbcommon.TbSimpleMsg{}
		json.NewDecoder(respBody).Decode(&failResult)
		return resultInfo, model.WebStatus{StatusCode: respStatus, Message: failResult.Message}

	}
	log.Println("ResultMessage : " + resultInfo.Result)

	json.NewDecoder(respBody).Decode(&resultInfo)
	log.Println(resultInfo)

	// webStatus.StatusCode = respStatus
	// webStatus.Message = resultInfo.Result
	return resultInfo, model.WebStatus{StatusCode: respStatus}
}

//Install the benchmark agent to specified MCIS
func InstallBenchmarkAgentToMcis(nameSpaceID string, mcisID string, mcisCommandInfo *tbmcis.McisCmdReq) (*tbmcis.RestPostCmdMcisResponseWrapper, model.WebStatus) {
	var originalUrl = "/ns/{nsId}/installBenchmarkAgent/mcis/{mcisId}"

	var paramMapper = make(map[string]string)
	paramMapper["{nsId}"] = nameSpaceID
	paramMapper["{mcisId}"] = mcisID
	urlParam := util.MappingUrlParameter(originalUrl, paramMapper)

	url := util.TUMBLEBUG + urlParam
	// url := util.TUMBLEBUG + "/ns/" + nameSpaceID + "/mcis/" + mcisID
	// /ns/{nsId}/install/mcis/{mcisId}
	pbytes, _ := json.Marshal(mcisCommandInfo)
	resp, err := util.CommonHttp(url, pbytes, http.MethodPost)

	returnMcisCommandResult := tbmcis.RestPostCmdMcisResponseWrapper{}
	returnStatus := model.WebStatus{}

	respBody := resp.Body
	respStatus := resp.StatusCode
	returnStatus.StatusCode = respStatus

	if err != nil {
		fmt.Println(err)
		return &returnMcisCommandResult, model.WebStatus{StatusCode: 500, Message: err.Error()}
	}

	if respStatus != 200 && respStatus != 201 {
		failResultInfo := tbcommon.TbSimpleMsg{}
		json.NewDecoder(respBody).Decode(&failResultInfo)
		log.Println("InstallBenchmarkAgentToMcis ", failResultInfo)
		return &returnMcisCommandResult, model.WebStatus{StatusCode: respStatus, Message: failResultInfo.Message}
	}

	json.NewDecoder(respBody).Decode(&returnMcisCommandResult)
	fmt.Println(returnMcisCommandResult)

	// return respBody, respStatusCode
	return &returnMcisCommandResult, returnStatus

	// resultMcisCommandResult := tbmcis.AgentInstallContentWrapper{}
	// if err != nil {
	// 	fmt.Println(err)
	// 	return &resultMcisCommandResult, model.WebStatus{StatusCode: 500, Message: err.Error()}
	// }

	// respBody := resp.Body
	// respStatus := resp.StatusCode

	// // TODO : result는 resultArray인데....
	// json.NewDecoder(respBody).Decode(resultMcisCommandResult)
	// fmt.Println(resultMcisCommandResult)
	// return &resultMcisCommandResult, model.WebStatus{StatusCode: respStatus}
	// if err != nil {
	// 	fmt.Println(err)
	// 	return nil, model.WebStatus{StatusCode: 500, Message: err.Error()}
	// }

	// respBody := resp.Body
	// respStatus := resp.StatusCode

	// return respBody, model.WebStatus{StatusCode: respStatus}
}

// Delete All MCISs
func DelAllMcis(nameSpaceID string, optionParam string) (tbcommon.TbSimpleMsg, model.WebStatus) {
	var originalUrl = "/ns/{nsId}/mcis"

	var paramMapper = make(map[string]string)
	paramMapper["{nsId}"] = nameSpaceID
	urlParam := util.MappingUrlParameter(originalUrl, paramMapper)

	optionParamVal := ""

	if optionParam != "" {
		optionParamVal = "?option=" + optionParam
	}

	// url := util.TUMBLEBUG + urlParam
	url := util.TUMBLEBUG + urlParam + optionParamVal

	// 경로안에 parameter가 있어 추가 param없이 호출 함.
	resp, err := util.CommonHttpWithoutParam(url, http.MethodDelete)

	resultInfo := tbcommon.TbSimpleMsg{}

	if err != nil {
		fmt.Println(err)
		return resultInfo, model.WebStatus{StatusCode: 500, Message: err.Error()}
	}

	respBody := resp.Body
	respStatus := resp.StatusCode

	if respStatus != 200 && respStatus != 201 {
		failResultInfo := tbcommon.TbSimpleMsg{}
		json.NewDecoder(respBody).Decode(&failResultInfo)
		log.Println("DelAllMcis ", failResultInfo)
		return resultInfo, model.WebStatus{StatusCode: respStatus, Message: failResultInfo.Message}
	}

	return resultInfo, model.WebStatus{StatusCode: respStatus}
}

// MCIS 삭제. TODO : 해당 namespace의 MCIS만 삭제 가능... 창 두개에서 1개는 MCIS삭제, 1개는 namespace 변경이 있을 수 있으므로 UI에서 namespace도 넘겨서 비교할 것.
// optionParam은 없거나 force 가 있음.
func DelMcis(nameSpaceID string, mcisID string, optionParam string) (io.ReadCloser, model.WebStatus) {
	var originalUrl = "/ns/{nsId}/mcis/{mcisId}"

	var paramMapper = make(map[string]string)
	paramMapper["{nsId}"] = nameSpaceID
	paramMapper["{mcisId}"] = mcisID

	urlParam := util.MappingUrlParameter(originalUrl, paramMapper)

	optionParamVal := ""
	// install, init, cpus, cpum, memR, memW, fioR, fioW, dbR, dbW, rtt, mrtt, clean
	if optionParam != "" {
		optionParamVal = "?option=" + optionParam
	}

	// url := util.TUMBLEBUG + urlParam
	url := util.TUMBLEBUG + urlParam + optionParamVal
	// url := util.TUMBLEBUG + "/ns/" + nameSpaceID + "/mcis/" + mcisID

	if mcisID == "" {
		return nil, model.WebStatus{StatusCode: 500, Message: "MCIS ID is required"}
	}

	// 경로안에 parameter가 있어 추가 param없이 호출 함.
	resp, err := util.CommonHttp(url, nil, http.MethodDelete)
	if err != nil {
		fmt.Println(err)
		return nil, model.WebStatus{StatusCode: 500, Message: err.Error()}
	}
	respBody := resp.Body
	respStatus := resp.StatusCode

	if respStatus != 200 && respStatus != 201 {
		failResultInfo := tbcommon.TbSimpleMsg{}
		json.NewDecoder(respBody).Decode(&failResultInfo)
		log.Println("DelMcis ", failResultInfo)
		return nil, model.WebStatus{StatusCode: respStatus, Message: failResultInfo.Message}
	}

	return respBody, model.WebStatus{StatusCode: respStatus}
}

// MCIS에 VM 생성. path에 mcisID가 있음. VMInfo에는 mcisID가 없음.
func RegVM(nameSpaceID string, mcisID string, vmInfo *tbmcis.TbVmInfo) (*tbmcis.TbVmInfo, model.WebStatus) {
	var originalUrl = "/ns/{nsId}/mcis/{mcisId}/vm"

	var paramMapper = make(map[string]string)
	paramMapper["{nsId}"] = nameSpaceID
	paramMapper["{mcisId}"] = mcisID
	urlParam := util.MappingUrlParameter(originalUrl, paramMapper)

	url := util.TUMBLEBUG + urlParam
	// url := util.TUMBLEBUG + "/ns/" + nameSpaceID + "/mcis/" + mcisID + "/vm"

	pbytes, _ := json.Marshal(vmInfo)
	resp, err := util.CommonHttp(url, pbytes, http.MethodPost)

	returnVmInfo := tbmcis.TbVmInfo{}
	returnStatus := model.WebStatus{}

	respBody := resp.Body
	respStatus := resp.StatusCode

	if err != nil {
		fmt.Println(err)
		return &returnVmInfo, model.WebStatus{StatusCode: 500, Message: err.Error()}
	}

	if respStatus != 200 && respStatus != 201 { // 호출은 정상이나, 가져온 결과값이 200, 201아닌 경우 message에 담겨있는 것을 WebStatus에 set
		errorInfo := model.ErrorInfo{}
		json.NewDecoder(respBody).Decode(&errorInfo)
		fmt.Println("respStatus != 200 reason ", errorInfo)
		returnStatus.Message = errorInfo.Message
	} else {
		json.NewDecoder(respBody).Decode(&returnVmInfo)
		fmt.Println(returnVmInfo)
	}
	returnStatus.StatusCode = respStatus

	// return respBody, respStatusCode
	return &returnVmInfo, returnStatus

	// resultVmResult := tumblebug.VmInfo{}
	// if err != nil {
	// 	fmt.Println(err)
	// 	return &resultVmResult, model.WebStatus{StatusCode: 500, Message: err.Error()}
	// }

	// respBody := resp.Body
	// respStatus := resp.StatusCode

	// // TODO : result는 resultArray인데....
	// json.NewDecoder(respBody).Decode(resultVmResult)
	// fmt.Println(resultVmResult)
	// return &resultVmResult, model.WebStatus{StatusCode: respStatus}
	// if err != nil {
	// 	fmt.Println(err)
	// 	return nil, model.WebStatus{StatusCode: 500, Message: err.Error()}
	// }

	// respBody := resp.Body
	// respStatus := resp.StatusCode

	// return respBody, model.WebStatus{StatusCode: respStatus}
}

func DelVM(nameSpaceID string, mcisID string, vmID string) (io.ReadCloser, model.WebStatus) {
	var originalUrl = "/ns/{nsId}/mcis/{mcisId}/vm/{vmId}"

	var paramMapper = make(map[string]string)
	paramMapper["{nsId}"] = nameSpaceID
	paramMapper["{mcisId}"] = mcisID
	paramMapper["{vmId}"] = vmID
	urlParam := util.MappingUrlParameter(originalUrl, paramMapper)

	url := util.TUMBLEBUG + urlParam
	// url := util.TUMBLEBUG + "/ns/" + nameSpaceID + "/mcis/" + mcisID + "/vm"
	// /ns/{nsId}/mcis/{mcisId}/vm/{vmId}

	if vmID == "" {
		return nil, model.WebStatus{StatusCode: 500, Message: "vmID ID is required"}
	}

	// 경로안에 parameter가 있어 추가 param없이 호출 함.
	resp, err := util.CommonHttp(url, nil, http.MethodDelete)
	// body, err := util.CommonHttpDelete(url, pbytes)
	if err != nil {
		fmt.Println(err)
		return nil, model.WebStatus{StatusCode: 500, Message: err.Error()}
	}
	// return body, err
	respBody := resp.Body
	respStatus := resp.StatusCode
	return respBody, model.WebStatus{StatusCode: respStatus}
}

// 특정 VM 조회
func GetVmData(nameSpaceID string, mcisID string, vmID string) (*tbmcis.TbVmInfo, model.WebStatus) {
	var originalUrl = "/ns/{nsId}/mcis/{mcisId}/vm/{vmId}"

	var paramMapper = make(map[string]string)
	paramMapper["{nsId}"] = nameSpaceID
	paramMapper["{mcisId}"] = mcisID
	paramMapper["{vmId}"] = vmID
	urlParam := util.MappingUrlParameter(originalUrl, paramMapper)

	url := util.TUMBLEBUG + urlParam
	// url := util.TUMBLEBUG + "/ns/" + nameSpaceID + "/mcis/" + mcisID + "/vm/" + vmID

	resp, err := util.CommonHttpWithoutParam(url, http.MethodGet)

	// defer body.Close()
	vmInfo := tbmcis.TbVmInfo{}
	if err != nil {
		fmt.Println(err)
		return &vmInfo, model.WebStatus{StatusCode: 500, Message: err.Error()}
	}
	// util.DisplayResponse(resp) // 수신내용 확인

	respBody := resp.Body
	respStatus := resp.StatusCode

	if respStatus != 200 && respStatus != 201 {
		failResultInfo := tbcommon.TbSimpleMsg{}
		json.NewDecoder(respBody).Decode(&failResultInfo)
		log.Println("GetVmData ", failResultInfo)
		return &vmInfo, model.WebStatus{StatusCode: respStatus, Message: failResultInfo.Message}
	}

	json.NewDecoder(respBody).Decode(&vmInfo)
	fmt.Println(vmInfo)

	return &vmInfo, model.WebStatus{StatusCode: respStatus}
}

// Get MCIS recommendation
// Deprecated at 0.4.5
// func GetMcisRecommand(nameSpaceID string, mcisID string, mcisRecommandReq *tumblebug.McisRecommendReq) (*tumblebug.McisRecommendInfo, model.WebStatus) {
// 	var originalUrl = "/ns/{nsId}/mcis/recommend"

// 	var paramMapper = make(map[string]string)
// 	paramMapper["{nsId}"] = nameSpaceID
// 	urlParam := util.MappingUrlParameter(originalUrl, paramMapper)

// 	url := util.TUMBLEBUG + urlParam
// 	// url := util.TUMBLEBUG + "/ns/" + nameSpaceID + "/mcis/recommend"

// 	pbytes, _ := json.Marshal(mcisRecommandReq)
// 	resp, err := util.CommonHttp(url, pbytes, http.MethodPost)

// 	returnMcisRecommendInfo := tumblebug.McisRecommendInfo{}
// 	returnStatus := model.WebStatus{}

// 	respBody := resp.Body
// 	respStatus := resp.StatusCode

// 	if err != nil {
// 		fmt.Println(err)
// 		return &returnMcisRecommendInfo, model.WebStatus{StatusCode: 500, Message: err.Error()}
// 	}
// 	log.Println(respBody)
// 	spew.Dump(respBody)
// 	if respStatus != 200 && respStatus != 201 { // 호출은 정상이나, 가져온 결과값이 200, 201아닌 경우 message에 담겨있는 것을 WebStatus에 set
// 		errorInfo := model.ErrorInfo{}
// 		json.NewDecoder(respBody).Decode(&errorInfo)
// 		fmt.Println("respStatus != 200 reason ", errorInfo)
// 		returnStatus.Message = errorInfo.Message
// 	} else {
// 		json.NewDecoder(respBody).Decode(&returnMcisRecommendInfo)
// 		fmt.Println(returnMcisRecommendInfo)
// 	}
// 	returnStatus.StatusCode = respStatus

// 	return &returnMcisRecommendInfo, returnStatus

// 	// mcisRecommandesult := tumblebug.McisRecommendInfo{}
// 	// if err != nil {
// 	// 	fmt.Println(err)
// 	// 	return &mcisRecommandesult, model.WebStatus{StatusCode: 500, Message: err.Error()}
// 	// }

// 	// respBody := resp.Body
// 	// respStatus := resp.StatusCode

// 	// // TODO : result는 resultArray인데....
// 	// json.NewDecoder(respBody).Decode(mcisRecommandesult)
// 	// fmt.Println(mcisRecommandesult)
// 	// return &mcisRecommandesult, model.WebStatus{StatusCode: respStatus}
// 	// if err != nil {
// 	// 	fmt.Println(err)
// 	// 	return nil, model.WebStatus{StatusCode: 500, Message: err.Error()}
// 	// }

// 	// respBody := resp.Body
// 	// respStatus := resp.StatusCode

// 	// return respBody, model.WebStatus{StatusCode: respStatus}
// }
