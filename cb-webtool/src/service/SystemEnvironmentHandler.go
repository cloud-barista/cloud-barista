package service

import (
	"encoding/json"
	"fmt"
	"log"

	// "math"
	"net/http"
	// "strconv"
	// "sync"

	//"github.com/davecgh/go-spew/spew"
	model "github.com/cloud-barista/cb-webtool/src/model"
	// "github.com/cloud-barista/cb-webtool/src/model/spider"
	// "github.com/cloud-barista/cb-webtool/src/model/tumblebug"
	tbcommon "github.com/cloud-barista/cb-webtool/src/model/tumblebug/common"
	// tbmcir "github.com/cloud-barista/cb-webtool/src/model/tumblebug/mcir"
	// tbmcis "github.com/cloud-barista/cb-webtool/src/model/tumblebug/mcis"

	util "github.com/cloud-barista/cb-webtool/src/util"

	// "github.com/labstack/echo"
)


func DelAllTbConfig() (tbcommon.TbSimpleMsg, model.WebStatus) {
	fmt.Println("Delete all configs start")
	var originalUrl = "/config"
	urlParam := util.MappingUrlParameter(originalUrl, nil)

	url := util.TUMBLEBUG + urlParam
	// url := util.TUMBLEBUG + "/ns"

	resp, err := util.CommonHttp(url, nil, http.MethodDelete)
	//body := HttpGetHandler(url)

	resultInfo := tbcommon.TbSimpleMsg{}

	if err != nil {
		return resultInfo, model.WebStatus{StatusCode: 500, Message: err.Error()}
	}

	respBody := resp.Body
	respStatus := resp.StatusCode
	

	json.NewDecoder(respBody).Decode(&resultInfo)
	log.Println(resultInfo)
	log.Println("ResultMessage : " + resultInfo.Message)

	if respStatus != 200 && respStatus != 201 {
		return resultInfo, model.WebStatus{StatusCode: respStatus, Message: resultInfo.Message}
	}

	return resultInfo, model.WebStatus{StatusCode: respStatus}
}


func GetAllTbConfig() (tbcommon.TbRestGetAllConfigResponse, model.WebStatus) {
	fmt.Println("Get all configs list start")
	var originalUrl = "/config"
	urlParam := util.MappingUrlParameter(originalUrl, nil)

	url := util.TUMBLEBUG + urlParam
	// url := util.TUMBLEBUG + "/ns"

	resp, err := util.CommonHttp(url, nil, http.MethodGet)
	//body := HttpGetHandler(url)

	resultAllConfigInfo := tbcommon.TbRestGetAllConfigResponse{}

	if err != nil {
		return resultAllConfigInfo, model.WebStatus{StatusCode: 500, Message: err.Error()}
	}

	respBody := resp.Body
	respStatus := resp.StatusCode
	
	// log.Println("ResultMessage : " + resultInfo.Message)

	if respStatus != 200 && respStatus != 201 {
		failResultInfo := tbcommon.TbSimpleMsg{}
        json.NewDecoder(respBody).Decode(&failResultInfo)
        return resultAllConfigInfo, model.WebStatus{StatusCode: respStatus, Message: failResultInfo.Message}
	}

	json.NewDecoder(respBody).Decode(&resultAllConfigInfo)
	log.Println(resultAllConfigInfo)

	return resultAllConfigInfo, model.WebStatus{StatusCode: respStatus}
}


func RegOrUpdateConfig(configReg tbcommon.TbConfigReq) (tbcommon.TbConfigInfo, model.WebStatus) {
	fmt.Println("Create or Update config ")
	//https://www.javaer101.com/ko/article/5704925.html 참조 : 값이 있는 것만 넘기기
	var originalUrl = "​/config"
	urlParam := util.MappingUrlParameter(originalUrl, nil)
	url := util.TUMBLEBUG + urlParam

	pbytes, _ := json.Marshal(configReg)
	resp, err := util.CommonHttp(url, pbytes, http.MethodPost)

	resultConfigInfo := tbcommon.TbConfigInfo{}

	if err != nil {
		return resultConfigInfo, model.WebStatus{StatusCode: 500, Message: err.Error()}
	}

	respBody := resp.Body
	respStatus := resp.StatusCode
	
	if respStatus != 200 && respStatus != 201 {
		failResultInfo := tbcommon.TbSimpleMsg{}
        json.NewDecoder(respBody).Decode(&failResultInfo)
        return resultConfigInfo, model.WebStatus{StatusCode: respStatus, Message: failResultInfo.Message}
	}

	json.NewDecoder(respBody).Decode(&resultConfigInfo)
	fmt.Println(resultConfigInfo)

	return resultConfigInfo, model.WebStatus{StatusCode: respStatus}

}

func DelTbConfig(configID string) (tbcommon.TbConfigInfo, model.WebStatus) {
	var originalUrl = "/config/{configId}"

	var paramMapper = make(map[string]string)
	paramMapper["{configId}"] = configID
	urlParam := util.MappingUrlParameter(originalUrl, paramMapper)

	url := util.TUMBLEBUG + urlParam
	resp, err := util.CommonHttpWithoutParam(url, http.MethodDelete)

	resultConfigInfo := tbcommon.TbConfigInfo{}

	if err != nil {
		return resultConfigInfo, model.WebStatus{StatusCode: 500, Message: err.Error()}
	}

	respBody := resp.Body
	respStatus := resp.StatusCode
	

	if respStatus != 200 && respStatus != 201 {
		failResultInfo := tbcommon.TbSimpleMsg{}
        json.NewDecoder(respBody).Decode(&failResultInfo)
        return resultConfigInfo, model.WebStatus{StatusCode: respStatus, Message: failResultInfo.Message}
	}

	json.NewDecoder(respBody).Decode(&resultConfigInfo)
	fmt.Println(resultConfigInfo)

	return resultConfigInfo, model.WebStatus{StatusCode: respStatus}
}


func GetTbConfig(configID string) (tbcommon.TbConfigInfo, model.WebStatus) {
	var originalUrl = "/config/{configId}"

	var paramMapper = make(map[string]string)
	paramMapper["{configId}"] = configID
	urlParam := util.MappingUrlParameter(originalUrl, paramMapper)

	url := util.TUMBLEBUG + urlParam
	resp, err := util.CommonHttpWithoutParam(url, http.MethodGet)

	resultConfigInfo := tbcommon.TbConfigInfo{}

	if err != nil {
		return resultConfigInfo, model.WebStatus{StatusCode: 500, Message: err.Error()}
	}

	respBody := resp.Body
	respStatus := resp.StatusCode
	

	if respStatus != 200 && respStatus != 201 {
		failResultInfo := tbcommon.TbSimpleMsg{}
        json.NewDecoder(respBody).Decode(&failResultInfo)
        return resultConfigInfo, model.WebStatus{StatusCode: respStatus, Message: failResultInfo.Message}
	}

	json.NewDecoder(respBody).Decode(&resultConfigInfo)
	fmt.Println(resultConfigInfo)

	return resultConfigInfo, model.WebStatus{StatusCode: respStatus}
}


