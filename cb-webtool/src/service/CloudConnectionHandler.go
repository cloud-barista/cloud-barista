package service

import (
	"encoding/json"
	"fmt"
	"io"
	"log"

	// "strings"
	"bytes"

	// "math"
	"net/http"
	// "strconv"
	// "sync"
	//"io/ioutil"
	//"github.com/davecgh/go-spew/spew"

	"github.com/cloud-barista/cb-webtool/src/model"
	// spider "github.com/cloud-barista/cb-webtool/src/model/spider"
	"github.com/cloud-barista/cb-webtool/src/model/spider"

	tbcommon "github.com/cloud-barista/cb-webtool/src/model/tumblebug/common"

	util "github.com/cloud-barista/cb-webtool/src/util"
)

//var CloudConnectionUrl = "http://15.165.16.67:1024"
// var CloudConnectionUrl = os.Getenv("SPIDER_URL")	// Const.go로 이동
// var TumbleUrl = os.Getenv("TUMBLE_URL") // Const.go로 이동

// type KeyValueInfo struct {
// 	Key   string `json:"Key"`
// 	Value string `json:"Value"`
// }
// type RegionInfo struct {
// 	RegionName       string `json:"RegionName"`
// 	ProviderName     string `json:"ProviderName"`
// 	KeyValueInfoList []KeyValueInfo
// }

// 뭐에쓰는 거지?
type RESP struct {
	Region []struct {
		RegionName       string                    `json:"RegionName"`
		ProviderName     string                    `json:"ProviderName"`
		KeyValueInfoList []spider.KeyValueInfoList `json:"KeyValueInfoList"`
	} `json:"region"`
}

// 뭐에쓰는 거지?
type ImageRESP struct {
	Image []struct {
		id               string                    `json:"id"`
		name             string                    `json:"name"`
		connectionName   string                    `json:"connectionName"`
		cspImageId       string                    `json:"cspImageId"`
		cspImageName     string                    `json:"cspImageName"`
		description      string                    `json:"description"`
		guestOS          string                    `json:"guestOS"`
		status           string                    `json:"status"`
		KeyValueInfoList []spider.KeyValueInfoList `json:"KeyValueList"`
	} `json:"image"`
}
type Image struct {
	id               string                    `json:"id"`
	name             string                    `json:"name"`
	connectionName   string                    `json:"connectionName"`
	cspImageId       string                    `json:"cspImageId"`
	cspImageName     string                    `json:"cspImageName"`
	description      string                    `json:"description"`
	guestOS          string                    `json:"guestOS"`
	status           string                    `json:"status"`
	KeyValueInfoList []spider.KeyValueInfoList `json:"KeyValueList"`
}
type IPStackInfo struct {
	IP          string  `json:"ip"`
	Lat         float64 `json:"latitude"`
	Long        float64 `json:"longitude"`
	CountryCode string  `json:"country_code"`
	VMName      string
	VMID        string
	Status      string
}

// 목록 : ListData
// 1개 : Data
// 등록 : Reg
// 삭제 : Del

// Cloud Provider 목록
func GetCloudOSList() ([]string, model.WebStatus) {

	var originalUrl = "/cloudos"
	urlParam := util.MappingUrlParameter(originalUrl, nil)
	url := util.SPIDER + urlParam

	resp, err := util.CommonHttp(url, nil, http.MethodGet)
	// defer resp.Close()
	if err != nil {
		fmt.Println(err)
		return nil, model.WebStatus{StatusCode: 500, Message: err.Error()}
	}

	respBody := resp.Body
	respStatus := resp.StatusCode

	cloudOs := map[string][]string{}
	json.NewDecoder(respBody).Decode(&cloudOs)
	fmt.Println(cloudOs["cloudos"])
	// TODO : mock을 추가할 것
	return cloudOs["cloudos"], model.WebStatus{StatusCode: respStatus}
}

// provider 별 connection count, connection 있는 provider 수
func GetCloudConnectionCountMap(cloudConnectionConfigInfoList []spider.CloudConnectionConfigInfo) (map[string]int, int) {
	connectionConfigCountMap := make(map[string]int)
	for _, connectionInfo := range cloudConnectionConfigInfoList {
		count := 0
		val, exists := connectionConfigCountMap[util.GetProviderName(connectionInfo.ProviderName)]
		if !exists {
			count = 1
		} else {
			count = val + 1
		}
		connectionConfigCountMap[util.GetProviderName(connectionInfo.ProviderName)] = count
	}

	providerCount := 0
	for i, _ := range connectionConfigCountMap {
		if i == "" {
		}
		providerCount++
	}
	return connectionConfigCountMap, providerCount
}

// 현재 설정된 connection 목록 GetConnectionConfigListData -> GetCloudConnectionConfigList로 변경
func GetCloudConnectionConfigList() ([]spider.CloudConnectionConfigInfo, model.WebStatus) {
	var originalUrl = "/connectionconfig"
	urlParam := util.MappingUrlParameter(originalUrl, nil)
	url := util.SPIDER + urlParam

	resp, err := util.CommonHttp(url, nil, http.MethodGet)
	// defer body.Close()

	if err != nil {
		fmt.Println(err)
		return nil, model.WebStatus{StatusCode: 500, Message: err.Error()}
	}

	respBody := resp.Body
	respStatus := resp.StatusCode

	cloudConnectionConfigInfo := map[string][]spider.CloudConnectionConfigInfo{}
	json.NewDecoder(respBody).Decode(&cloudConnectionConfigInfo)
	fmt.Println(cloudConnectionConfigInfo["connectionconfig"])

	return cloudConnectionConfigInfo["connectionconfig"], model.WebStatus{StatusCode: respStatus}
}

// Connection 상세
func GetCloudConnectionConfigData(configName string) (spider.CloudConnectionConfigInfo, model.WebStatus) {
	var originalUrl = "/connectionconfig/{{config_name}}"

	var paramMapper = make(map[string]string)
	paramMapper["{{config_name}}"] = configName
	urlParam := util.MappingUrlParameter(originalUrl, paramMapper)

	url := util.SPIDER + urlParam

	// url := util.SPIDER + "/connectionconfig/" + configName
	fmt.Println("=========== GetCloudConnectionConfigData : ", configName)
	cloudConnectionConfigInfo := spider.CloudConnectionConfigInfo{}

	resp, err := util.CommonHttp(url, nil, http.MethodGet)
	// defer body.Close()

	if err != nil {
		fmt.Println(err)
		return cloudConnectionConfigInfo, model.WebStatus{StatusCode: 500, Message: err.Error()}
	}

	respBody := resp.Body
	respStatus := resp.StatusCode

	json.NewDecoder(respBody).Decode(&cloudConnectionConfigInfo)
	fmt.Println(cloudConnectionConfigInfo)
	return cloudConnectionConfigInfo, model.WebStatus{StatusCode: respStatus}
}

// CloudConnectionConfigInfo 등록
// func RegCloudConnectionConfig(cloudConnectionConfigInfo *model.CloudConnectionConfigInfo) (io.ReadCloser, model.WebStatus) {
func RegCloudConnectionConfig(cloudConnectionConfigInfo *spider.CloudConnectionConfigInfo) (*spider.CloudConnectionConfigInfo, model.WebStatus) {
	var originalUrl = "/connectionconfig"
	urlParam := util.MappingUrlParameter(originalUrl, nil)

	url := util.SPIDER + urlParam
	// buff := bytes.NewBuffer(pbytes)
	// url := util.SPIDER + "/connectionconfig"

	fmt.Println("cloudConnectionConfigInfo : ", cloudConnectionConfigInfo)

	// body, err := util.CommonHttpPost(url, regionInfo)
	pbytes, _ := json.Marshal(cloudConnectionConfigInfo)
	// body, err := util.CommonHttpPost(url, pbytes)
	resp, err := util.CommonHttp(url, pbytes, http.MethodPost)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return nil, model.WebStatus{StatusCode: 500, Message: err.Error()}
	// }

	// respBody := resp.Body
	// respStatus := resp.StatusCode
	//cloudConnectionConfigInfo
	// return respBody, model.WebStatus{StatusCode: respStatus}

	returnCloudConnectionConfigInfo := spider.CloudConnectionConfigInfo{}
	if err != nil {
		fmt.Println(err)
		return &returnCloudConnectionConfigInfo, model.WebStatus{StatusCode: 500, Message: err.Error()}
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
		json.NewDecoder(respBody).Decode(&returnCloudConnectionConfigInfo)
		fmt.Println(returnCloudConnectionConfigInfo)
	}
	returnStatus.StatusCode = respStatus

	return &returnCloudConnectionConfigInfo, returnStatus
}

// CloudConnectionConfigInfo 삭제
func DelCloudConnectionConfig(configName string) (io.ReadCloser, model.WebStatus) {
	var originalUrl = "/connectionconfig/{{config_name}}"

	var paramMapper = make(map[string]string)
	paramMapper["{{config_name}}"] = configName
	urlParam := util.MappingUrlParameter(originalUrl, paramMapper)

	url := util.SPIDER + urlParam

	// buff := bytes.NewBuffer(pbytes)
	// url := util.SPIDER + "/connectionconfig/" + configName

	fmt.Println("DelCloudConnectionConfig : ", configName)

	// body, err := util.CommonHttpPost(url, regionInfo)

	pbytes, _ := json.Marshal(configName)
	// body, err := util.CommonHttpDelete(url, pbytes)
	resp, err := util.CommonHttp(url, pbytes, http.MethodDelete)

	if err != nil {
		fmt.Println(err)
		return nil, model.WebStatus{StatusCode: 500, Message: err.Error()}
	}
	// return body, err
	respBody := resp.Body
	respStatus := resp.StatusCode

	return respBody, model.WebStatus{StatusCode: respStatus}
}

// 현재 설정된 region 목록
func GetRegionList() ([]spider.RegionInfo, model.WebStatus) {
	var originalUrl = "/region"
	urlParam := util.MappingUrlParameter(originalUrl, nil)

	url := util.SPIDER + urlParam

	// url := util.SPIDER + "/region"
	// fmt.Println("=========== GetRegionListData : ", url)

	resp, err := util.CommonHttp(url, nil, http.MethodGet)
	// defer body.Close()

	if err != nil {
		fmt.Println(err)
		return nil, model.WebStatus{StatusCode: 500, Message: err.Error()}
	}

	respBody := resp.Body
	respStatus := resp.StatusCode

	regionList := map[string][]spider.RegionInfo{}
	json.NewDecoder(respBody).Decode(&regionList)
	fmt.Println(regionList["region"])

	return regionList["region"], model.WebStatus{StatusCode: respStatus}
}

func GetRegionData(regionName string) (*tbcommon.TbRegion, model.WebStatus) {
	var originalUrl = "/region/{regionName}"

	var paramMapper = make(map[string]string)
	paramMapper["{regionName}"] = regionName
	urlParam := util.MappingUrlParameter(originalUrl, paramMapper)

	url := util.SPIDER + urlParam

	// url := util.SPIDER + "/region/" + regionName
	fmt.Println("=========== GetRegionData : ", regionName)

	resp, err := util.CommonHttp(url, nil, http.MethodGet)
	// defer body.Close()
	// regionInfo := spider.RegionInfo{}
	regionInfo := tbcommon.TbRegion{}
	if err != nil {
		fmt.Println(err)
		return &regionInfo, model.WebStatus{StatusCode: 500, Message: err.Error()}
	}

	respBody := resp.Body
	respStatus := resp.StatusCode

	// regionList := map[string][]string{}
	// // regionList := map[string][]spider.RegionInfo{}
	// json.NewDecoder(body).Decode(&regionList)
	// fmt.Println(regionList)	// map[KeyValueInfoList:[] ProviderName:[] RegionName:[]]
	// // fmt.Println(regionList["connectionconfig"])

	json.NewDecoder(respBody).Decode(&regionInfo)
	// fmt.Println(regionInfo)
	// fmt.Println(regionInfo.KeyValueInfoList)
	return &regionInfo, model.WebStatus{StatusCode: respStatus}
}

// Region 등록
// func RegRegion(regionInfo *spider.RegionInfo) (io.ReadCloser, model.WebStatus) {
func RegRegion(regionInfo *spider.RegionInfo) (*spider.RegionInfo, model.WebStatus) {
	var originalUrl = "/region"
	urlParam := util.MappingUrlParameter(originalUrl, nil)

	url := util.SPIDER + urlParam

	// buff := bytes.NewBuffer(pbytes)
	// url := util.SPIDER + "/region"

	fmt.Println("RegRegion : ", regionInfo)

	// body, err := util.CommonHttpPost(url, regionInfo)
	pbytes, _ := json.Marshal(regionInfo)
	// body, err := util.CommonHttpPost(url, pbytes)
	resp, err := util.CommonHttp(url, pbytes, http.MethodPost)

	// if err != nil {
	// 	fmt.Println(err)
	// 	return nil, model.WebStatus{StatusCode: 500, Message: err.Error()}
	// }

	// respBody := resp.Body
	// respStatus := resp.StatusCode

	// return respBody, model.WebStatus{StatusCode: respStatus}
	respBody := resp.Body
	respStatus := resp.StatusCode

	returnRegionInfo := spider.RegionInfo{}
	returnStatus := model.WebStatus{}

	if err != nil {
		fmt.Println(err)
		return &returnRegionInfo, model.WebStatus{StatusCode: 500, Message: err.Error()}
	}

	if respStatus != 200 && respStatus != 201 { // 호출은 정상이나, 가져온 결과값이 200, 201아닌 경우 message에 담겨있는 것을 WebStatus에 set
		errorInfo := model.ErrorInfo{}
		json.NewDecoder(respBody).Decode(&errorInfo)
		fmt.Println("respStatus != 200 reason ", errorInfo)
		returnStatus.Message = errorInfo.Message
	} else {
		json.NewDecoder(respBody).Decode(&returnRegionInfo)
		fmt.Println(returnRegionInfo)
	}
	returnStatus.StatusCode = respStatus

	return &returnRegionInfo, returnStatus
}

// Region 삭제
func DelRegion(regionName string) (io.ReadCloser, model.WebStatus) {
	var originalUrl = "/region/{{region_name}}"

	var paramMapper = make(map[string]string)
	paramMapper["{{region_name}}"] = regionName
	urlParam := util.MappingUrlParameter(originalUrl, paramMapper)

	url := util.SPIDER + urlParam

	// buff := bytes.NewBuffer(pbytes)
	// url := util.SPIDER + "/region/" + regionName

	fmt.Println("DelRegion : ", regionName)

	// body, err := util.CommonHttpPost(url, regionInfo)

	pbytes, _ := json.Marshal(regionName)
	// body, err := util.CommonHttpDelete(url, pbytes)
	resp, err := util.CommonHttp(url, pbytes, http.MethodDelete)

	if err != nil {
		fmt.Println(err)
		return nil, model.WebStatus{StatusCode: 500, Message: err.Error()}
	}
	// return body, err
	respBody := resp.Body
	respStatus := resp.StatusCode

	return respBody, model.WebStatus{StatusCode: respStatus}
}

// 현재 설정된 credential 목록 : 목록에서는 key의 value는 ...으로 표시
func GetCredentialList() ([]spider.CredentialInfo, model.WebStatus) {
	var originalUrl = "/credential"
	urlParam := util.MappingUrlParameter(originalUrl, nil)

	url := util.SPIDER + urlParam

	// SPIDER == SPIDER
	// url := util.SPIDER + "/credential"
	// fmt.Println("=========== GetRegionListData : ", url)

	resp, err := util.CommonHttp(url, nil, http.MethodGet)
	// defer body.Close()

	if err != nil {
		fmt.Println(err)
		return nil, model.WebStatus{StatusCode: 500, Message: err.Error()}
	}

	respBody := resp.Body
	respStatus := resp.StatusCode

	credentialList := map[string][]spider.CredentialInfo{}
	json.NewDecoder(respBody).Decode(&credentialList)
	fmt.Println(credentialList["credential"])
	// TODO : key의 value에 ...표시
	for _, credentialInfo := range credentialList["credential"] {
		fmt.Println("credentialInfo : ", credentialInfo)
		keyValueInfoList := credentialInfo.KeyValueInfoList
		fmt.Println("before keyValueInfoList : ", keyValueInfoList)
		for _, keyValueInfo := range keyValueInfoList {
			keyValueInfo.Value = "..."
		}
		// fmt.Println("after keyValueInfoList : ", keyValueInfoList)
	}

	return credentialList["credential"], model.WebStatus{StatusCode: respStatus}
}

// Credential 상세조회
func GetCredentialData(credentialName string) (*spider.CredentialInfo, model.WebStatus) {
	var originalUrl = "/credential/{{credential_name}}"

	var paramMapper = make(map[string]string)
	paramMapper["{{credential_name}}"] = credentialName
	urlParam := util.MappingUrlParameter(originalUrl, paramMapper)

	url := util.SPIDER + urlParam

	// url := util.SPIDER + "/credential/" + credentialName
	fmt.Println("=========== GetCredentialData : ", credentialName)

	resp, err := util.CommonHttp(url, nil, http.MethodGet)
	// defer body.Close()
	credentialInfo := spider.CredentialInfo{}
	if err != nil {
		fmt.Println(err)
		return &credentialInfo, model.WebStatus{StatusCode: 500, Message: err.Error()}
	}

	respBody := resp.Body
	respStatus := resp.StatusCode

	json.NewDecoder(respBody).Decode(&credentialInfo)
	// fmt.Println(credentialInfo)
	// fmt.Println(credentialInfo.KeyValueInfoList)
	return &credentialInfo, model.WebStatus{StatusCode: respStatus}
}

// Credential 등록
func RegCredential(credentialInfo *spider.CredentialInfo) (*spider.CredentialInfo, model.WebStatus) {
	var originalUrl = "/credential"
	urlParam := util.MappingUrlParameter(originalUrl, nil)

	url := util.SPIDER + urlParam

	fmt.Println("RegCredential : ", credentialInfo)
	returnCredentialInfo := spider.CredentialInfo{}
	returnStatus := model.WebStatus{}

	// GCP의 경우 value에 \n 이 포함되어 있기 때문에 이것이 넘어올 때는 \\n 형태로 넘어온다.  이것을 \으로 replace 해야
	//
	var credentialBuffer bytes.Buffer
	if credentialInfo.ProviderName == "GCP" {

		credentialBuffer.WriteString(`{`)
		credentialBuffer.WriteString(`"CredentialName":"` + credentialInfo.CredentialName + `"`)
		credentialBuffer.WriteString(`,"ProviderName":"` + credentialInfo.ProviderName + `"`)
		credentialBuffer.WriteString(`,"KeyValueInfoList":[`)

		for mapIndex, keyValueInfo := range credentialInfo.KeyValueInfoList {
			// for mapIndex, _ := range credentialInfo.KeyValueInfoList {
			gcpKey := keyValueInfo.Key
			gcpValue := keyValueInfo.Value
			// replacedValue := gcpValue
			if mapIndex > 0 {
				credentialBuffer.WriteString(`,`)
			}
			credentialBuffer.WriteString(`{"Key":"` + gcpKey + `","Value":"` + gcpValue + `"}`)
			// if gcpKey == "private_key" {
			// // 	fmt.Println(gcpValue)
			// // 	// fmt.Println(keyValueInfo)
			// // 	fmt.Println("--------- before / after -----------------")
			// 	replacedValue = strings.Replace(gcpValue, "\\n", "\n", -1)
			// 	// replacedValue = strings.Replace(gcpValue, "\n", "\\n", -1)
			// 	// replacedValue = "`" + gcpValue + "`"
			// 	keyValueInfo.Value = replacedValue
			// 	fmt.Println(replacedValue)
			// 	// fmt.Println(keyValueInfo)

			// }
		}
		// fmt.Println("GCP RegCredential : ", credentialInfo)

		credentialBuffer.WriteString(`]`)
		credentialBuffer.WriteString(`}`)
		// fmt.Println("******* ")

		// fmt.Println(credentialBuffer.String())
		resp, err := util.CommonHttpBytes(url, &credentialBuffer, http.MethodPost)

		// pbytes, marshalErr := json.Marshal(credentialInfo)
		// if marshalErr != nil {
		// 	fmt.Println(" ------------------ ")
		// 	fmt.Println(marshalErr)
		// }
		// fmt.Println(string(pbytes))

		// resp, err := util.CommonHttp(url, pbytes, http.MethodPost)
		if err != nil {
			fmt.Println(err)
			return &returnCredentialInfo, model.WebStatus{StatusCode: 500, Message: err.Error()}
		}

		respBody := resp.Body
		respStatus := resp.StatusCode

		if respStatus != 200 && respStatus != 201 { // 호출은 정상이나, 가져온 결과값이 200, 201아닌 경우 message에 담겨있는 것을 WebStatus에 set
			errorInfo := model.ErrorInfo{}
			json.NewDecoder(respBody).Decode(&errorInfo)
			fmt.Println("respStatus != 200 reason ", errorInfo)
			returnStatus.Message = errorInfo.Message
		} else {
			json.NewDecoder(respBody).Decode(&returnCredentialInfo)
			fmt.Println(returnCredentialInfo)
		}
		returnStatus.StatusCode = respStatus
	} else if credentialInfo.ProviderName == "OPENSTACK" {

		credentialBuffer.WriteString(`{`)
		credentialBuffer.WriteString(`"CredentialName":"` + credentialInfo.CredentialName + `"`)
		credentialBuffer.WriteString(`,"ProviderName":"` + credentialInfo.ProviderName + `"`)
		credentialBuffer.WriteString(`,"KeyValueInfoList":[`)

		for mapIndex, keyValueInfo := range credentialInfo.KeyValueInfoList {
			// for mapIndex, _ := range credentialInfo.KeyValueInfoList {
			openstackKey := keyValueInfo.Key
			openstackValue := keyValueInfo.Value
			// replacedValue := gcpValue
			if mapIndex > 0 {
				credentialBuffer.WriteString(`,`)
			}
			credentialBuffer.WriteString(`{"Key":"` + openstackKey + `","Value":"` + openstackValue + `"}`)
			// if gcpKey == "private_key" {
			// // 	fmt.Println(gcpValue)
			// // 	// fmt.Println(keyValueInfo)
			// // 	fmt.Println("--------- before / after -----------------")
			// 	replacedValue = strings.Replace(gcpValue, "\\n", "\n", -1)
			// 	// replacedValue = strings.Replace(gcpValue, "\n", "\\n", -1)
			// 	// replacedValue = "`" + gcpValue + "`"
			// 	keyValueInfo.Value = replacedValue
			// 	fmt.Println(replacedValue)
			// 	// fmt.Println(keyValueInfo)

			// }
		}
		// fmt.Println("GCP RegCredential : ", credentialInfo)

		credentialBuffer.WriteString(`]`)
		credentialBuffer.WriteString(`}`)
		// fmt.Println("******* ")

		// fmt.Println(credentialBuffer.String())
		resp, err := util.CommonHttpBytes(url, &credentialBuffer, http.MethodPost)

		// pbytes, marshalErr := json.Marshal(credentialInfo)
		// if marshalErr != nil {
		// 	fmt.Println(" ------------------ ")
		// 	fmt.Println(marshalErr)
		// }
		// fmt.Println(string(pbytes))

		// resp, err := util.CommonHttp(url, pbytes, http.MethodPost)
		if err != nil {
			fmt.Println(err)
			return &returnCredentialInfo, model.WebStatus{StatusCode: 500, Message: err.Error()}
		}

		respBody := resp.Body
		respStatus := resp.StatusCode

		if respStatus != 200 && respStatus != 201 { // 호출은 정상이나, 가져온 결과값이 200, 201아닌 경우 message에 담겨있는 것을 WebStatus에 set
			errorInfo := model.ErrorInfo{}
			json.NewDecoder(respBody).Decode(&errorInfo)
			fmt.Println("respStatus != 200 reason ", errorInfo)
			returnStatus.Message = errorInfo.Message
		} else {
			json.NewDecoder(respBody).Decode(&returnCredentialInfo)
			fmt.Println(returnCredentialInfo)
		}
		returnStatus.StatusCode = respStatus
	} else {
		// body, err := util.CommonHttpPost(url, regionInfo)

		pbytes, marshalErr := json.Marshal(credentialInfo)
		// pbytes, marshalErr := json.Marshal(credentialInfo)
		if marshalErr != nil {
			fmt.Println(" ------------------ ")
			fmt.Println(marshalErr)
		}
		fmt.Println(string(pbytes))
		resp, err := util.CommonHttp(url, pbytes, http.MethodPost)

		if err != nil {
			fmt.Println(err)
			return &returnCredentialInfo, model.WebStatus{StatusCode: 500, Message: err.Error()}
		}

		respBody := resp.Body
		respStatus := resp.StatusCode

		if respStatus != 200 && respStatus != 201 { // 호출은 정상이나, 가져온 결과값이 200, 201아닌 경우 message에 담겨있는 것을 WebStatus에 set
			errorInfo := model.ErrorInfo{}
			json.NewDecoder(respBody).Decode(&errorInfo)
			fmt.Println("respStatus != 200 reason ", errorInfo)
			returnStatus.Message = errorInfo.Message
		} else {
			json.NewDecoder(respBody).Decode(&returnCredentialInfo)
			fmt.Println(returnCredentialInfo)
		}
		returnStatus.StatusCode = respStatus
	}

	// if err != nil {
	// 	fmt.Println(err)
	// 	return nil, model.WebStatus{StatusCode: 500, Message: err.Error()}
	// }

	// respBody := resp.Body
	// respStatus := resp.StatusCode

	// // util.DisplayResponse(resp)
	// return respBody, model.WebStatus{StatusCode: respStatus}

	// respBody := resp.Body
	// respStatus := resp.StatusCode

	// if err != nil {
	// 	fmt.Println(err)
	// 	return &returnCredentialInfo, model.WebStatus{StatusCode: 500, Message: err.Error()}
	// }

	// if respStatus != 200 && respStatus != 201 { // 호출은 정상이나, 가져온 결과값이 200, 201아닌 경우 message에 담겨있는 것을 WebStatus에 set
	// 	errorInfo := model.ErrorInfo{}
	// 	json.NewDecoder(respBody).Decode(&errorInfo)
	// 	fmt.Println("respStatus != 200 reason ", errorInfo)
	// 	returnStatus.Message = errorInfo.Message
	// } else {
	// 	json.NewDecoder(respBody).Decode(&returnCredentialInfo)
	// 	fmt.Println(returnCredentialInfo)
	// }
	// returnStatus.StatusCode = respStatus

	return &returnCredentialInfo, returnStatus
}

// Credential 삭제
func DelCredential(credentialName string) (model.WebStatus, model.WebStatus) {
	webStatus := model.WebStatus{}

	var originalUrl = "/credential/{{credential_name}}"

	var paramMapper = make(map[string]string)
	paramMapper["{{credential_name}}"] = credentialName
	urlParam := util.MappingUrlParameter(originalUrl, paramMapper)

	url := util.SPIDER + urlParam
	// buff := bytes.NewBuffer(pbytes)
	// url := util.SPIDER + "/credential/" + credentialName

	fmt.Println("DelCredential : ", credentialName)

	pbytes, _ := json.Marshal(credentialName)
	// body, err := util.CommonHttpDelete(url, pbytes)
	resp, err := util.CommonHttp(url, pbytes, http.MethodDelete)

	if err != nil {
		fmt.Println(err)
		return webStatus, model.WebStatus{StatusCode: 500, Message: err.Error()}
	}
	// return body, err
	respBody := resp.Body
	respStatus := resp.StatusCode
	resultInfo := model.ResultInfo{}

	json.NewDecoder(respBody).Decode(&resultInfo)
	log.Println(resultInfo)
	log.Println("ResultMessage : " + resultInfo.Message)

	if respStatus != 200 && respStatus != 201 {
		return model.WebStatus{}, model.WebStatus{StatusCode: respStatus, Message: resultInfo.Message}
	}
	webStatus.StatusCode = respStatus
	webStatus.Message = resultInfo.Message
	return webStatus, model.WebStatus{StatusCode: respStatus}
	//return respBody, model.WebStatus{StatusCode: respStatus}
}

// 현재 설정된 Driver 목록
func GetDriverList() ([]spider.DriverInfo, model.WebStatus) {
	var originalUrl = "/driver"
	urlParam := util.MappingUrlParameter(originalUrl, nil)

	url := util.SPIDER + urlParam

	// url := util.SPIDER + "/driver"
	fmt.Println("=========== GetDriverListData : ", url)

	resp, err := util.CommonHttp(url, nil, http.MethodGet)
	// defer body.Close()

	if err != nil {
		fmt.Println(err)
		return nil, model.WebStatus{StatusCode: 500, Message: err.Error()}
	}

	respBody := resp.Body
	respStatus := resp.StatusCode

	driverList := map[string][]spider.DriverInfo{}
	json.NewDecoder(respBody).Decode(&driverList)
	fmt.Println(driverList["driver"])

	return driverList["driver"], model.WebStatus{StatusCode: respStatus}
}

// Driver 상세조회
func GetDriverData(driverlName string) (*spider.DriverInfo, model.WebStatus) {
	var originalUrl = "/driver/{{driver_name}}"

	var paramMapper = make(map[string]string)
	paramMapper["{{driver_name}}"] = driverlName
	urlParam := util.MappingUrlParameter(originalUrl, paramMapper)

	url := util.SPIDER + urlParam

	//url := util.SPIDER + "/driver/" + driverlName
	fmt.Println("=========== GetDriverData : ", url)

	resp, err := util.CommonHttp(url, nil, http.MethodGet)
	// defer body.Close()
	driverInfo := spider.DriverInfo{}
	if err != nil {
		fmt.Println(err)
		return &driverInfo, model.WebStatus{StatusCode: 500, Message: err.Error()}
	}

	respBody := resp.Body
	respStatus := resp.StatusCode

	json.NewDecoder(respBody).Decode(&driverInfo)
	fmt.Println(driverInfo)
	return &driverInfo, model.WebStatus{StatusCode: respStatus}
}

// Driver 등록
func RegDriver(driverInfo *spider.DriverInfo) (*spider.DriverInfo, model.WebStatus) {
	var originalUrl = "/driver"
	urlParam := util.MappingUrlParameter(originalUrl, nil)

	url := util.SPIDER + urlParam
	// buff := bytes.NewBuffer(pbytes)
	// url := util.SPIDER + "/driver"

	fmt.Println("driverInfo : ", driverInfo)

	// body, err := util.CommonHttpPost(url, regionInfo)
	pbytes, _ := json.Marshal(driverInfo)
	resp, err := util.CommonHttp(url, pbytes, http.MethodPost)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return nil, model.WebStatus{StatusCode: 500, Message: err.Error()}
	// }

	// respBody := resp.Body
	// respStatus := resp.StatusCode
	// return respBody, model.WebStatus{StatusCode: respStatus}

	respBody := resp.Body
	respStatus := resp.StatusCode

	returnDriverInfo := spider.DriverInfo{}
	returnStatus := model.WebStatus{}

	if err != nil {
		fmt.Println(err)
		return &returnDriverInfo, model.WebStatus{StatusCode: 500, Message: err.Error()}
	}

	if respStatus != 200 && respStatus != 201 { // 호출은 정상이나, 가져온 결과값이 200, 201아닌 경우 message에 담겨있는 것을 WebStatus에 set
		errorInfo := model.ErrorInfo{}
		json.NewDecoder(respBody).Decode(&errorInfo)
		fmt.Println("respStatus != 200 reason ", errorInfo)
		returnStatus.Message = errorInfo.Message
	} else {
		json.NewDecoder(respBody).Decode(&returnDriverInfo)
		fmt.Println(returnDriverInfo)
	}
	returnStatus.StatusCode = respStatus

	return &returnDriverInfo, returnStatus
}

// Driver 삭제
func DelDriver(driverName string) (io.ReadCloser, model.WebStatus) {
	var originalUrl = "/driver/{{driver_name}}"

	var paramMapper = make(map[string]string)
	paramMapper["{{driver_name}}"] = driverName
	urlParam := util.MappingUrlParameter(originalUrl, paramMapper)

	url := util.SPIDER + urlParam
	// buff := bytes.NewBuffer(pbytes)
	// url := util.SPIDER + "/driver/" + driverName

	fmt.Println("driverName : ", driverName)

	pbytes, _ := json.Marshal(driverName)
	// body, err := util.CommonHttpDelete(url, pbytes)
	resp, err := util.CommonHttp(url, pbytes, http.MethodDelete)

	if err != nil {
		fmt.Println(err)
		return nil, model.WebStatus{StatusCode: 500, Message: err.Error()}
	}
	// return body, err
	respBody := resp.Body
	respStatus := resp.StatusCode
	return respBody, model.WebStatus{StatusCode: respStatus}
}

///////////// Config
// 현재 설정된 Config 목록  TODO :Spider에서 /config 가 없는 것 같은데.... 나중에 확인해서 안쓰면 제거할 것
func GetConfigList() ([]spider.ConfigInfo, model.WebStatus) {
	url := util.SPIDER + "/config"
	fmt.Println("=========== GetConfigListData : ", url)

	resp, err := util.CommonHttp(url, nil, http.MethodGet)
	// defer body.Close()

	if err != nil {
		fmt.Println(err)
		return nil, model.WebStatus{StatusCode: 500, Message: err.Error()}
	}

	respBody := resp.Body
	respStatus := resp.StatusCode

	configList := map[string][]spider.ConfigInfo{}
	json.NewDecoder(respBody).Decode(&configList)
	fmt.Println(configList["config"])

	return configList["config"], model.WebStatus{StatusCode: respStatus}
}

// Config 상세조회
func GetConfigData(configID string) (*spider.ConfigInfo, model.WebStatus) {
	url := util.SPIDER + "/config/" + configID
	fmt.Println("=========== GetConfigData : ", url)

	resp, err := util.CommonHttp(url, nil, http.MethodGet)
	// defer body.Close()
	configInfo := spider.ConfigInfo{}
	if err != nil {
		fmt.Println(err)
		return &configInfo, model.WebStatus{StatusCode: 500, Message: err.Error()}
	}

	respBody := resp.Body
	respStatus := resp.StatusCode

	json.NewDecoder(respBody).Decode(&configInfo)
	fmt.Println(configInfo)
	return &configInfo, model.WebStatus{StatusCode: respStatus}
}

// Driver 등록
func RegConfig(configInfo *spider.ConfigInfo) (*spider.ConfigInfo, model.WebStatus) {
	// buff := bytes.NewBuffer(pbytes)
	url := util.SPIDER + "/config"

	fmt.Println("configInfo : ", configInfo)

	// body, err := util.CommonHttpPost(url, regionInfo)
	pbytes, _ := json.Marshal(configInfo)
	resp, err := util.CommonHttp(url, pbytes, http.MethodPost)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return nil, model.WebStatus{StatusCode: 500, Message: err.Error()}
	// }

	// respBody := resp.Body
	// respStatus := resp.StatusCode
	// return respBody, model.WebStatus{StatusCode: respStatus}
	respBody := resp.Body
	respStatus := resp.StatusCode

	returnConfigInfo := spider.ConfigInfo{}
	returnStatus := model.WebStatus{}

	if err != nil {
		fmt.Println(err)
		return &returnConfigInfo, model.WebStatus{StatusCode: 500, Message: err.Error()}
	}

	if respStatus != 200 && respStatus != 201 { // 호출은 정상이나, 가져온 결과값이 200, 201아닌 경우 message에 담겨있는 것을 WebStatus에 set
		errorInfo := model.ErrorInfo{}
		json.NewDecoder(respBody).Decode(&errorInfo)
		fmt.Println("respStatus != 200 reason ", errorInfo)
		returnStatus.Message = errorInfo.Message
	} else {
		json.NewDecoder(respBody).Decode(&returnConfigInfo)
		fmt.Println(returnConfigInfo)
	}
	returnStatus.StatusCode = respStatus

	return &returnConfigInfo, returnStatus
}

// Driver 삭제
func DelConfig(configID string) (io.ReadCloser, model.WebStatus) {

	// buff := bytes.NewBuffer(pbytes)
	url := util.SPIDER + "/config/" + configID

	fmt.Println("configID : ", configID)

	pbytes, _ := json.Marshal(configID)
	// body, err := util.CommonHttpDelete(url, pbytes)
	resp, err := util.CommonHttp(url, pbytes, http.MethodDelete)

	if err != nil {
		fmt.Println(err)
		return nil, model.WebStatus{StatusCode: 500, Message: err.Error()}
	}
	// return body, err
	respBody := resp.Body
	respStatus := resp.StatusCode
	return respBody, model.WebStatus{StatusCode: respStatus}
}

/////////////
// 해당 namespace의 vpc 목록 조회 -> ResourceHandler로 이동
// func GetVnetList(nameSpaceID string) (io.ReadCloser, error) {
// url := TumbleUrl + "ns/" + nameSpaceID + "/resources/vNet"

// fmt.Println("nameSpaceID : ", nameSpaceID)

// pbytes, _ := json.Marshal(nameSpaceID)
// body, err := util.CommonHttp(url, pbytes, http.MethodGet)

// if err != nil {
// 	fmt.Println(err)
// }
// return body, err
// }

// // vpc 상세 조회-> ResourceHandler로 이동
// func GetVpcData(nameSpaceID string, vNetID string) (io.ReadCloser, error) {
// 	url := TumbleUrl + "ns/" + nameSpaceID + "/resources/vNet"

// 	fmt.Println("nameSpaceID : ", nameSpaceID)

// 	pbytes, _ := json.Marshal(nameSpaceID)
// 	body, err := util.CommonHttp(url, pbytes, http.MethodGet)

// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	return body, err
// }

// vpc 등록 -> ResourceHandler로 이동
// func RegVpc(nameSpaceID string, vnetInfo *model.VNetInfo) (io.ReadCloser, error) {
// 	url := TumbleUrl + "ns/" + nameSpaceID + "/resources/vNet"

// 	fmt.Println("nameSpaceID : ", nameSpaceID)

// 	pbytes, _ := json.Marshal(vnetInfo)
// 	body, err := util.CommonHttp(url, pbytes, http.MethodPost)

// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	return body, err
// }

// vpc 삭제 -> ResourceHandler로 이동
// func DelVpc(nameSpaceID string, vNetID string) (io.ReadCloser, error) {
// 	url := TumbleUrl + "ns/" + nameSpaceID + "/resources/vNet" + vNetID

// 	fmt.Println("nameSpaceID : ", nameSpaceID)

// 	pbytes, _ := json.Marshal(vNetID)
// 	body, err := util.CommonHttp(url, pbytes, http.MethodDelete)

// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	return body, err
// }

// func GetConnectionconfig(drivername string) CloudConnectionInfo {
// 	url := NameSpaceUrl + "/driver/" + drivername

// 	// resp, err := http.Get(url)

// 	// if err != nil {
// 	fmt.Println("request URL : ", url)
// 	// }

// 	// defer resp.Body.Close()
// 	body := HttpGetHandler(url)
// 	nsInfo := CloudConnectionInfo{}

// 	json.NewDecoder(body).Decode(&nsInfo)
// 	fmt.Println("nsInfo : ", nsInfo.ID)
// 	return nsInfo

// }
// func GetImageList() []Image {
// 	url := CloudConnectionUrl + "/connectionconfig"
// 	// resp, err := http.Get(url)
// 	// if err != nil {
// 	// 	fmt.Println("request URL : ", url)
// 	// }

// 	// defer resp.Body.Close()
// 	body := HttpGetHandler(url)
// 	defer body.Close()

// 	nsInfo := ImageRESP{}
// 	json.NewDecoder(body).Decode(&nsInfo)
// 	fmt.Println("nsInfo : ", nsInfo.Image[0].id)
// 	var info []Image
// 	for _, item := range nsInfo.Image {
// 		reg := Image{
// 			id:             item.id,
// 			name:           item.name,
// 			connectionName: item.connectionName,
// 			cspImageId:     item.cspImageId,
// 			cspImageName:   item.cspImageName,
// 			description:    item.description,
// 			guestOS:        item.guestOS,
// 			status:         item.status,
// 		}
// 		info = append(info, reg)
// 	}
// 	return info

// }

// func GetConnectionList() []CloudConnectionInfo {
// 	url := CloudConnectionUrl + "/connectionconfig"
// 	// resp, err := http.Get(url)
// 	// if err != nil {
// 	// 	fmt.Println("request URL : ", url)
// 	// }

// 	// defer resp.Body.Close()
// 	body := HttpGetHandler(url)
// 	defer body.Close()

// 	nsInfo := map[string][]CloudConnectionInfo{}
// 	json.NewDecoder(body).Decode(&nsInfo)
// 	//fmt.Println("nsInfo : ", nsInfo["connectionconfig"][0].ID)
// 	return nsInfo["connectionconfig"]

// }

// func GetDriverReg() []CloudConnectionInfo {
// 	url := NameSpaceUrl + "/driver"

// 	body := HttpGetHandler(url)
// 	defer body.Close()
// 	nsInfo := map[string][]CloudConnectionInfo{}
// 	json.NewDecoder(body).Decode(&nsInfo)

// 	return nsInfo["driver"]

// }

// func GetCredentialList() []CloudConnectionInfo {
// 	url := CloudConnectionUrl + "/credential"
// 	// resp, err := http.Get(url)
// 	// if err != nil {
// 	// 	fmt.Println("request URL : ", url)
// 	// }

// 	// defer resp.Body.Close()
// 	body := HttpGetHandler(url)
// 	defer body.Close()
// 	nsInfo := map[string][]CloudConnectionInfo{}
// 	json.NewDecoder(body).Decode(&nsInfo)
// 	// fmt.Println("nsInfo : ", nsInfo["credential"][0].ID)
// 	return nsInfo["credential"]

// }

// func GetRegionList() []RegionInfo {
// 	url := CloudConnectionUrl + "/region"
// 	fmt.Println("=========== Get Start Region List : ", url)

// 	body := HttpGetHandler(url)
// 	defer body.Close()

// 	nsInfo := RESP{}

// 	json.NewDecoder(body).Decode(&nsInfo)
// 	var info []RegionInfo

// 	for _, item := range nsInfo.Region {

// 		reg := RegionInfo{
// 			RegionName:   item.RegionName,
// 			ProviderName: item.ProviderName,
// 		}
// 		info = append(info, reg)
// 	}
// 	fmt.Println("info region list : ", info)
// 	return info

// }

// func GetCredentialReg() []CloudConnectionInfo {
// 	url := CloudConnectionUrl + "/credential"

// 	body := HttpGetHandler(url)
// 	defer body.Close()

// 	nsInfo := map[string][]CloudConnectionInfo{}
// 	json.NewDecoder(body).Decode(&nsInfo)
// 	fmt.Println("nsInfo : ", nsInfo["credential"][0].ID)
// 	return nsInfo["credential"]

// }

// func GetGeoMetryInfo(wg *sync.WaitGroup, ip_address string, status string, vm_id string, vm_name string, returnResult *[]IPStackInfo) {
// 	defer wg.Done() //goroutin sync done

// 	apiUrl := "http://api.ipstack.com/"
// 	access_key := "86c895286435070c0369a53d2d0b03d1"
// 	url := apiUrl + ip_address + "?access_key=" + access_key
// 	resp, err := http.Get(url)
// 	fmt.Println("GetGeoMetryInfo request URL : ", url)
// 	if err != nil {
// 		fmt.Println("GetGeoMetryInfo request URL : ", url)
// 	}
// 	defer resp.Body.Close()

// 	//그냥 스트링으로 반환해서 프론트에서 JSON.parse로 처리 하는 방법도 괜찮네
// 	//spew.Dump(resp.Body)
// 	// bytes, _ := ioutil.ReadAll(resp.Body)
// 	// str := string(bytes)
// 	// fmt.Println(str)
// 	// *returnStr = append(*returnStr, str)

// 	ipStackInfo := IPStackInfo{
// 		VMID:   vm_id,
// 		Status: status,
// 		VMName: vm_name,
// 	}

// 	json.NewDecoder(resp.Body).Decode(&ipStackInfo)
// 	fmt.Println("Get GeoMetry INFO :", ipStackInfo)

// 	*returnResult = append(*returnResult, ipStackInfo)
// }

// func RegNS() error {

// }

// func RequestGet(url string) {
// 	resp, err := http.Get(url)
// 	if err != nil {
// 		fmt.Println("request URL : ", url)
// 	}

// 	defer resp.Body.Close()
// 	nsInfo := map[string][]NSInfo{}
// 	fmt.Println("nsInfo type : ", reflect.TypeOf(nsInfo))
// 	json.NewDecoder(resp.Body).Decode(&nsInfo)
// 	fmt.Println("nsInfo : ", nsInfo["ns"][0].ID)

// 	// data, err := ioutil.ReadAll(resp.Body)
// 	// if err != nil {
// 	// 	fmt.Println("Get Data Error")
// 	// }
// 	// fmt.Println("GetData : ", string(data))

// }
