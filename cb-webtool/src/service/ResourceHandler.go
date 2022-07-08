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
	tbmcir "github.com/cloud-barista/cb-webtool/src/model/tumblebug/mcir"
	tbmcis "github.com/cloud-barista/cb-webtool/src/model/tumblebug/mcis"

	util "github.com/cloud-barista/cb-webtool/src/util"

	"github.com/labstack/echo"
)

func RegFirewallRules(nameSpaceID string, securityGroupID string, firewallRuleReq *tbmcir.TbFirewallRulesWrapper) (*tbmcir.TbSecurityGroupInfo, model.WebStatus) {
	var originalUrl = "/ns/{nsId}/resources/securityGroup/{securityGroupId}/rules"
	var paramMapper = make(map[string]string)
	paramMapper["{nsId}"] = nameSpaceID
	paramMapper["{securityGroupId}"] = securityGroupID
	urlParam := util.MappingUrlParameter(originalUrl, paramMapper)
	url := util.TUMBLEBUG + urlParam

	pbytes, _ := json.Marshal(firewallRuleReq)
	fmt.Println(string(pbytes))
	resp, err := util.CommonHttp(url, pbytes, http.MethodPost)
	securityGroupInfo := tbmcir.TbSecurityGroupInfo{}
	if err != nil {
		fmt.Println(err)
		return &securityGroupInfo, model.WebStatus{StatusCode: 500, Message: err.Error()}
	}

	respBody := resp.Body
	respStatus := resp.StatusCode
	log.Println("respStatusCode = ", resp.StatusCode)
	log.Println("respStatus = ", resp.Status)
	if respStatus != 200 && respStatus != 201 {
		errorInfo := model.ErrorInfo{}
		json.NewDecoder(respBody).Decode(&errorInfo)
		fmt.Println(errorInfo)
		return &securityGroupInfo, model.WebStatus{StatusCode: 500, Message: errorInfo.Message}
	}

	// 응답에 생성한 객체값이 옴
	json.NewDecoder(respBody).Decode(&securityGroupInfo)
	fmt.Println(securityGroupInfo)
	// return respBody, respStatusCode
	return &securityGroupInfo, model.WebStatus{StatusCode: respStatus}
}

func DelFirewallRules(nameSpaceID string, securityGroupID string, firewallRuleReq *tbmcir.TbFirewallRulesWrapper) (*tbmcir.TbSecurityGroupInfo, model.WebStatus) {
	var originalUrl = "/ns/{nsId}/resources/securityGroup/{securityGroupId}/rules"
	var paramMapper = make(map[string]string)
	paramMapper["{nsId}"] = nameSpaceID
	paramMapper["{securityGroupId}"] = securityGroupID
	urlParam := util.MappingUrlParameter(originalUrl, paramMapper)
	url := util.TUMBLEBUG + urlParam

	pbytes, _ := json.Marshal(firewallRuleReq)
	fmt.Println(string(pbytes))
	resp, err := util.CommonHttp(url, pbytes, http.MethodDelete)

	securityGroupInfo := tbmcir.TbSecurityGroupInfo{}
	if err != nil {
		fmt.Println(err)
		return &securityGroupInfo, model.WebStatus{StatusCode: 500, Message: err.Error()}
	}

	respBody := resp.Body
	respStatus := resp.StatusCode
	resultInfo := model.ResultInfo{}

	json.NewDecoder(respBody).Decode(&resultInfo)
	log.Println(resultInfo)
	log.Println("ResultMessage : " + resultInfo.Message)

	if respStatus != 200 && respStatus != 201 {
		return &securityGroupInfo, model.WebStatus{StatusCode: respStatus, Message: resultInfo.Message}
	}
	json.NewDecoder(respBody).Decode(&securityGroupInfo)
	fmt.Println(securityGroupInfo)
	return &securityGroupInfo, model.WebStatus{StatusCode: respStatus}
}

// 해당 namespace의 vpc 목록 조회
//func GetVnetList(nameSpaceID string) (io.ReadCloser, error) {
func GetVnetList(nameSpaceID string) ([]tbmcir.TbVNetInfo, model.WebStatus) {
	fmt.Println("GetVnetList ************ : ")
	var originalUrl = "/ns/{nsId}/resources/vNet"
	var paramMapper = make(map[string]string)
	paramMapper["{nsId}"] = nameSpaceID
	urlParam := util.MappingUrlParameter(originalUrl, paramMapper)

	url := util.TUMBLEBUG + urlParam
	// url := util.TUMBLEBUG + "/ns/" + nameSpaceID + "/resources/vNet"

	pbytes, _ := json.Marshal(nameSpaceID)
	// body, err := util.CommonHttpGet(url)
	resp, err := util.CommonHttp(url, pbytes, http.MethodGet)

	if err != nil {
		fmt.Println(err)
		return nil, model.WebStatus{StatusCode: 500, Message: err.Error()}
	}
	// defer body.Close()
	respBody := resp.Body
	respStatus := resp.StatusCode

	// return respBody, respStatus
	log.Println(respBody)
	vNetInfoList := map[string][]tbmcir.TbVNetInfo{}
	json.NewDecoder(respBody).Decode(&vNetInfoList)
	//spew.Dump(body)
	fmt.Println(vNetInfoList["vNet"])

	return vNetInfoList["vNet"], model.WebStatus{StatusCode: respStatus}
}

// ID목록만 조회
func GetVnetListByID(nameSpaceID string, filterKeyParam string, filterValParam string) ([]string, model.WebStatus) {
	fmt.Println("GetVnetList ************ : ")
	var originalUrl = "/ns/{nsId}/resources/vNet"
	var paramMapper = make(map[string]string)
	paramMapper["{nsId}"] = nameSpaceID
	urlParam := util.MappingUrlParameter(originalUrl, paramMapper)

	//if optionParam != ""{
	//	urlParam += "?option=" + optionParam
	//}
	urlParam += "?option=id"
	if filterKeyParam != "" {
		urlParam += "&filterKey=" + filterKeyParam
		urlParam += "&filterVal=" + filterValParam
	}
	url := util.TUMBLEBUG + urlParam
	// url := util.TUMBLEBUG + "/ns/" + nameSpaceID + "/resources/vNet"

	//pbytes, _ := json.Marshal(nameSpaceID)
	// body, err := util.CommonHttpGet(url)
	//resp, err := util.CommonHttp(url, pbytes, http.MethodGet)
	resp, err := util.CommonHttpWithoutParam(url, http.MethodGet)

	if err != nil {
		fmt.Println(err)
		return nil, model.WebStatus{StatusCode: 500, Message: err.Error()}
	}
	// defer body.Close()
	respBody := resp.Body
	respStatus := resp.StatusCode

	// return respBody, respStatus
	log.Println(respBody)
	//vNetInfoList := map[string][]string{}
	vNetInfoList := tbcommon.TbIdList{}
	json.NewDecoder(respBody).Decode(&vNetInfoList)
	//spew.Dump(body)
	//fmt.Println(vNetInfoList["idList"])

	//return vNetInfoList["idList"], model.WebStatus{StatusCode: respStatus}
	return vNetInfoList.IDList, model.WebStatus{StatusCode: respStatus}
}

// List 조회시 optionParam 추가
func GetVnetListByOption(nameSpaceID string, optionParam string, filterKeyParam string, filterValParam string) ([]tbmcir.TbVNetInfo, model.WebStatus) {
	fmt.Println("GetVnetListByOption ************ : ")
	var originalUrl = "/ns/{nsId}/resources/vNet"
	var paramMapper = make(map[string]string)
	paramMapper["{nsId}"] = nameSpaceID
	urlParam := util.MappingUrlParameter(originalUrl, paramMapper)

	if optionParam != "" {
		urlParam += "?option=" + optionParam
	}
	if filterKeyParam != "" {
		urlParam += "&filterKey=" + filterKeyParam
		urlParam += "&filterVal=" + filterValParam
	}
	url := util.TUMBLEBUG + urlParam
	// url := util.TUMBLEBUG + "/ns/" + nameSpaceID + "/resources/vNet"

	//pbytes, _ := json.Marshal(nameSpaceID)
	// body, err := util.CommonHttpGet(url)
	//resp, err := util.CommonHttp(url, pbytes, http.MethodGet)
	resp, err := util.CommonHttpWithoutParam(url, http.MethodGet)

	if err != nil {
		fmt.Println(err)
		return nil, model.WebStatus{StatusCode: 500, Message: err.Error()}
	}
	// defer body.Close()
	respBody := resp.Body
	respStatus := resp.StatusCode

	// return respBody, respStatus
	log.Println(respBody)
	vNetInfoList := map[string][]tbmcir.TbVNetInfo{}
	json.NewDecoder(respBody).Decode(&vNetInfoList)
	//spew.Dump(body)
	fmt.Println(vNetInfoList["vNet"])

	return vNetInfoList["vNet"], model.WebStatus{StatusCode: respStatus}
}

// vpc 상세 조회-> ResourceHandler로 이동
func GetVpcData(nameSpaceID string, vNetID string) (*tbmcir.TbVNetInfo, model.WebStatus) {
	var originalUrl = "/ns/{nsId}/resources/vNet/{vNetId}"
	var paramMapper = make(map[string]string)
	paramMapper["{nsId}"] = nameSpaceID
	paramMapper["{vNetId}"] = vNetID
	urlParam := util.MappingUrlParameter(originalUrl, paramMapper)
	url := util.TUMBLEBUG + urlParam
	// url := util.TUMBLEBUG + "/ns/" + nameSpaceID + "/resources/vNet/" + vNetID

	resp, err := util.CommonHttpWithoutParam(url, http.MethodGet)
	vNetInfo := tbmcir.TbVNetInfo{}
	if err != nil {
		fmt.Println(err)
		return &vNetInfo, model.WebStatus{StatusCode: 500, Message: err.Error()}
	}

	respBody := resp.Body
	respStatus := resp.StatusCode

	// json.NewDecoder(body).Decode(&vNetInfo)
	json.NewDecoder(respBody).Decode(&vNetInfo)
	fmt.Println(vNetInfo)

	// return vNetInfo, err
	return &vNetInfo, model.WebStatus{StatusCode: respStatus}
}

// vpc 등록
// option=register 항목은 TB에서 자동으로 넣을 때 사용하는 param으로 webtool에서 사용하지 않음.
func RegVpc(nameSpaceID string, vnetRegInfo *tbmcir.TbVNetReq) (*tbmcir.TbVNetInfo, model.WebStatus) {
	var originalUrl = "/ns/{nsId}/resources/vNet"
	var paramMapper = make(map[string]string)
	paramMapper["{nsId}"] = nameSpaceID
	urlParam := util.MappingUrlParameter(originalUrl, paramMapper)

	url := util.TUMBLEBUG + urlParam
	// url := util.TUMBLEBUG + "/ns/" + nameSpaceID + "/resources/vNet"

	fmt.Println("vnetRegInfo : ", vnetRegInfo)

	pbytes, _ := json.Marshal(vnetRegInfo)
	fmt.Println(string(pbytes))
	resp, err := util.CommonHttp(url, pbytes, http.MethodPost)
	vNetInfo := tbmcir.TbVNetInfo{}
	if err != nil {
		fmt.Println(err)
		return &vNetInfo, model.WebStatus{StatusCode: 500, Message: err.Error()}
	}

	respBody := resp.Body
	respStatus := resp.StatusCode
	fmt.Println("respStatus ", respStatus)

	if respStatus == 500 {
		webStatus := model.WebStatus{}
		json.NewDecoder(respBody).Decode(&webStatus)
		fmt.Println(webStatus)
		webStatus.StatusCode = respStatus
		return &vNetInfo, webStatus
	}
	// 응답에 생성한 객체값이 옴
	json.NewDecoder(respBody).Decode(&vNetInfo)
	fmt.Println(vNetInfo)

	return &vNetInfo, model.WebStatus{StatusCode: respStatus}
}

// vpc 삭제
func DelVpc(nameSpaceID string, vNetID string) (model.WebStatus, model.WebStatus) {
	webStatus := model.WebStatus{}
	// if ValidateString(vNetID) != nil {
	if len(vNetID) == 0 {
		log.Println("vNetID 가 없으면 해당 namespace의 모든 vpc가 삭제되므로 처리할 수 없습니다.")
		return webStatus, model.WebStatus{StatusCode: 4040, Message: "vNetID 가 없으면 해당 namespace의 모든 vpc가 삭제되므로 처리할 수 없습니다."}
	}
	var originalUrl = "/ns/{nsId}/resources/vNet/{vNetId}"
	var paramMapper = make(map[string]string)
	paramMapper["{nsId}"] = nameSpaceID
	paramMapper["{vNetId}"] = vNetID
	urlParam := util.MappingUrlParameter(originalUrl, paramMapper)
	url := util.TUMBLEBUG + urlParam
	// url := util.TUMBLEBUG + "/ns/" + nameSpaceID + "/resources/vNet/" + vNetID

	fmt.Println("vNetID : ", vNetID)

	pbytes, _ := json.Marshal(vNetID)
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
}

// 전체 vpc 삭제
func DelAllVpc(nameSpaceID string) (tbcommon.TbSimpleMsg, model.WebStatus) {
	var originalUrl = "/ns/{nsId}/resources/vNet"
	var paramMapper = make(map[string]string)
	paramMapper["{nsId}"] = nameSpaceID

	urlParam := util.MappingUrlParameter(originalUrl, paramMapper)

	url := util.TUMBLEBUG + urlParam

	resp, err := util.CommonHttpWithoutParam(url, http.MethodDelete)

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

// 해당 namespace의 SecurityGroup 목록 조회
func GetSecurityGroupList(nameSpaceID string) ([]tbmcir.TbSecurityGroupInfo, model.WebStatus) {
	fmt.Println("GetSecurityGroupList ************ : ")
	var originalUrl = "/ns/{nsId}/resources/securityGroup"
	var paramMapper = make(map[string]string)
	paramMapper["{nsId}"] = nameSpaceID
	urlParam := util.MappingUrlParameter(originalUrl, paramMapper)

	url := util.TUMBLEBUG + urlParam
	// url := util.TUMBLEBUG + "/ns/" + nameSpaceID + "/resources/securityGroup"

	pbytes, _ := json.Marshal(nameSpaceID)
	resp, err := util.CommonHttp(url, pbytes, http.MethodGet)

	if err != nil {
		fmt.Println(err)
		return nil, model.WebStatus{StatusCode: 500, Message: err.Error()}
	}
	// defer body.Close()
	respBody := resp.Body
	respStatus := resp.StatusCode

	// return respBody, respStatus
	log.Println(respBody)
	securityGroupList := map[string][]tbmcir.TbSecurityGroupInfo{}

	json.NewDecoder(respBody).Decode(&securityGroupList)
	//spew.Dump(body)
	fmt.Println(securityGroupList["securityGroup"])

	return securityGroupList["securityGroup"], model.WebStatus{StatusCode: respStatus}
}

// ID만 조회
func GetSecurityGroupListByOptionID(nameSpaceID string, optionParam string, filterKeyParam string, filterValParam string) ([]string, model.WebStatus) {
	fmt.Println("GetSecurityGroupListByOptionID ************ : ")
	var originalUrl = "/ns/{nsId}/resources/securityGroup"
	var paramMapper = make(map[string]string)
	paramMapper["{nsId}"] = nameSpaceID
	urlParam := util.MappingUrlParameter(originalUrl, paramMapper)

	if optionParam != "" {
		urlParam += "?option=" + optionParam
	} else {
		urlParam += "?option="
	}
	if filterKeyParam != "" {
		urlParam += "&filterKey=" + filterKeyParam
		urlParam += "&filterVal=" + filterValParam
	}

	url := util.TUMBLEBUG + urlParam
	// url := util.TUMBLEBUG + "/ns/" + nameSpaceID + "/resources/securityGroup"
	resp, err := util.CommonHttpWithoutParam(url, http.MethodGet)
	//resp, err := util.CommonHttp(url, pbytes, http.MethodGet)

	if err != nil {
		// 	// Tumblebug 접속 확인하라고
		// fmt.Println(err)
		// panic(err)
		return nil, model.WebStatus{StatusCode: 500, Message: err.Error()}
	}

	respBody := resp.Body
	respStatus := resp.StatusCode

	//securityGroupList := map[string][]string{}
	securityGroupList := tbcommon.TbIdList{}
	// defer body.Close()
	json.NewDecoder(respBody).Decode(&securityGroupList)
	//spew.Dump(body)
	//fmt.Println(securityGroupList["idList"])
	fmt.Println(securityGroupList.IDList)

	//return securityGroupList["idList"], model.WebStatus{StatusCode: respStatus}
	return securityGroupList.IDList, model.WebStatus{StatusCode: respStatus}
}

// SecurityGroupList 조회 시 Option에 해당하는 값만 조회. GetSecurityGroupList와 TB 호출은 동일하나 option 사용으로 받아오는 param이 다름
func GetSecurityGroupListByOption(nameSpaceID string, optionParam string, filterKeyParam string, filterValParam string) ([]tbmcir.TbSecurityGroupInfo, model.WebStatus) {
	fmt.Println("GetSecurityGroupListByOption ************ : ")
	var originalUrl = "/ns/{nsId}/resources/securityGroup"
	var paramMapper = make(map[string]string)
	paramMapper["{nsId}"] = nameSpaceID
	urlParam := util.MappingUrlParameter(originalUrl, paramMapper)

	if optionParam != "" {
		urlParam += "?option=" + optionParam
	} else {
		urlParam += "?option="
	}
	if filterKeyParam != "" {
		urlParam += "&filterKey=" + filterKeyParam
		urlParam += "&filterVal=" + filterValParam
	}

	url := util.TUMBLEBUG + urlParam
	// url := util.TUMBLEBUG + "/ns/" + nameSpaceID + "/resources/securityGroup"
	resp, err := util.CommonHttpWithoutParam(url, http.MethodGet)
	//resp, err := util.CommonHttp(url, pbytes, http.MethodGet)

	if err != nil {
		// 	// Tumblebug 접속 확인하라고
		// fmt.Println(err)
		// panic(err)
		return nil, model.WebStatus{StatusCode: 500, Message: err.Error()}
	}

	respBody := resp.Body
	respStatus := resp.StatusCode

	securityGroupList := map[string][]tbmcir.TbSecurityGroupInfo{}

	json.NewDecoder(respBody).Decode(&securityGroupList)
	//spew.Dump(body)
	fmt.Println(securityGroupList["securityGroup"])

	return securityGroupList["securityGroup"], model.WebStatus{StatusCode: respStatus}
}

// SecurityGroup 상세 조회
func GetSecurityGroupData(nameSpaceID string, securityGroupID string) (*tbmcir.TbSecurityGroupInfo, model.WebStatus) {
	var originalUrl = "/ns/{nsId}/resources/securityGroup/{securityGroupId}"
	var paramMapper = make(map[string]string)
	paramMapper["{nsId}"] = nameSpaceID
	paramMapper["{securityGroupId}"] = securityGroupID
	urlParam := util.MappingUrlParameter(originalUrl, paramMapper)
	url := util.TUMBLEBUG + urlParam
	// url := util.TUMBLEBUG + "/ns/" + nameSpaceID + "/resources/securityGroup/" + securityGroupID

	fmt.Println("nameSpaceID : ", nameSpaceID)

	// pbytes, _ := json.Marshal(nameSpaceID)
	// body, err := util.CommonHttpGet(url)
	resp, err := util.CommonHttpWithoutParam(url, http.MethodGet)
	securityGroupInfo := tbmcir.TbSecurityGroupInfo{}
	if err != nil {
		fmt.Println(err)
		return &securityGroupInfo, model.WebStatus{StatusCode: 500, Message: err.Error()}
	}

	respBody := resp.Body
	respStatus := resp.StatusCode

	json.NewDecoder(respBody).Decode(&securityGroupInfo)
	fmt.Println(securityGroupInfo)

	return &securityGroupInfo, model.WebStatus{StatusCode: respStatus}
}

// SecurityGroup 등록
// option=register 항목은 TB에서 자동으로 넣을 때 사용하는 param으로 webtool에서 사용하지 않음.
func RegSecurityGroup(nameSpaceID string, securityGroupRegInfo *tbmcir.TbSecurityGroupReq) (*tbmcir.TbSecurityGroupInfo, model.WebStatus) {
	fmt.Println("RegSecurityGroup : ")

	var originalUrl = "/ns/{nsId}/resources/securityGroup"
	var paramMapper = make(map[string]string)
	paramMapper["{nsId}"] = nameSpaceID
	urlParam := util.MappingUrlParameter(originalUrl, paramMapper)

	url := util.TUMBLEBUG + urlParam
	// url := util.TUMBLEBUG + "/ns/" + nameSpaceID + "/resources/securityGroup"

	pbytes, _ := json.Marshal(securityGroupRegInfo)
	fmt.Println(string(pbytes))
	resp, err := util.CommonHttp(url, pbytes, http.MethodPost)
	securityGroupInfo := tbmcir.TbSecurityGroupInfo{}
	if err != nil {
		log.Println("-----")
		fmt.Println(err)
		log.Println("-----1111")
		fmt.Println(err.Error())
		log.Println("-----222")
		return &securityGroupInfo, model.WebStatus{StatusCode: 500, Message: err.Error()}
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
	json.NewDecoder(respBody).Decode(&securityGroupInfo)
	fmt.Println(securityGroupInfo)
	// return respBody, respStatusCode
	return &securityGroupInfo, model.WebStatus{StatusCode: respStatus}
}

// 해당 Namespace의 모든 SecurityGroup 삭제
func DelAllSecurityGroup(nameSpaceID string) (model.WebStatus, model.WebStatus) {
	var originalUrl = "/ns/{nsId}/resources/securityGroup"
	var paramMapper = make(map[string]string)
	paramMapper["{nsId}"] = nameSpaceID
	urlParam := util.MappingUrlParameter(originalUrl, paramMapper)
	url := util.TUMBLEBUG + urlParam
	// url := util.TUMBLEBUG + "/ns/" + nameSpaceID + "/resources/securityGroup/"

	resp, err := util.CommonHttp(url, nil, http.MethodDelete)
	webStatus := model.WebStatus{}
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

// SecurityGroup 삭제
func DelSecurityGroup(nameSpaceID string, securityGroupID string) (model.WebStatus, model.WebStatus) {
	webStatus := model.WebStatus{}
	// if ValidateString(vNetID) != nil {
	if len(securityGroupID) == 0 {
		log.Println("securityGroupID 가 없으면 해당 namespace의 모든 securityGroup이 삭제되므로 처리할 수 없습니다.")
		return webStatus, model.WebStatus{StatusCode: 4040, Message: "securityGroupID 가 없으면 해당 namespace의 모든 securityGroup이 삭제되므로 처리할 수 없습니다."}
	}

	var originalUrl = "/ns/{nsId}/resources/securityGroup/{securityGroupId}"
	var paramMapper = make(map[string]string)
	paramMapper["{nsId}"] = nameSpaceID
	paramMapper["{securityGroupId}"] = securityGroupID
	urlParam := util.MappingUrlParameter(originalUrl, paramMapper)
	url := util.TUMBLEBUG + urlParam
	// url := util.TUMBLEBUG + "/ns/" + nameSpaceID + "/resources/securityGroup/" + securityGroupID

	fmt.Println("securityGroupID : ", securityGroupID)

	pbytes, _ := json.Marshal(securityGroupID)
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
	// respStatus := resp.Status
	// log.Println("respStatusCode = ", respStatusCode)
	// log.Println("respStatus = ", respStatus)

	// return respBody, model.WebStatus{StatusCode: respStatus}
}

// SSHKey 목록 조회 : /ns/{nsId}/resources/sshKey
func GetSshKeyInfoList(nameSpaceID string) ([]tbmcir.TbSshKeyInfo, model.WebStatus) {
	fmt.Println("GetSshKeyInfoList ************ : ")
	var originalUrl = "/ns/{nsId}/resources/sshKey"
	var paramMapper = make(map[string]string)
	paramMapper["{nsId}"] = nameSpaceID
	urlParam := util.MappingUrlParameter(originalUrl, paramMapper)

	url := util.TUMBLEBUG + urlParam
	// url := util.TUMBLEBUG + "/ns/" + nameSpaceID + "/resources/sshKey"

	pbytes, _ := json.Marshal(nameSpaceID)
	resp, err := util.CommonHttp(url, pbytes, http.MethodGet)

	if err != nil {
		fmt.Println(err)
		return nil, model.WebStatus{StatusCode: 500, Message: err.Error()}
	}
	// defer body.Close()
	respBody := resp.Body
	respStatus := resp.StatusCode

	// return respBody, respStatus
	log.Println(respBody)
	sshKeyList := map[string][]tbmcir.TbSshKeyInfo{}

	json.NewDecoder(respBody).Decode(&sshKeyList)
	//spew.Dump(body)
	fmt.Println(sshKeyList["sshKey"])

	return sshKeyList["sshKey"], model.WebStatus{StatusCode: respStatus}

}

func GetSshKeyInfoListByID(nameSpaceID string, filterKeyParam string, filterValParam string) ([]string, model.WebStatus) {
	fmt.Println("GetSshKeyInfoListByID ************ : ")
	var originalUrl = "/ns/{nsId}/resources/sshKey"
	var paramMapper = make(map[string]string)
	paramMapper["{nsId}"] = nameSpaceID
	urlParam := util.MappingUrlParameter(originalUrl, paramMapper)
	//url := util.TUMBLEBUG + urlParam
	urlParam += "?option=id"
	if filterKeyParam != "" {
		urlParam += "&filterKey=" + filterKeyParam
		urlParam += "&filterVal=" + filterValParam
	}
	url := util.TUMBLEBUG + urlParam
	// url := util.TUMBLEBUG + "/ns/" + nameSpaceID + "/resources/sshKey"

	//resp, err := util.CommonHttp(url, pbytes, http.MethodGet)
	resp, err := util.CommonHttpWithoutParam(url, http.MethodGet)

	if err != nil {
		fmt.Println(err)
		return nil, model.WebStatus{StatusCode: 500, Message: err.Error()}
	}
	// defer body.Close()
	respBody := resp.Body
	respStatus := resp.StatusCode

	// return respBody, respStatus
	log.Println(respBody)

	//sshKeyList := map[string][]string{}
	sshKeyList := tbcommon.TbIdList{}
	// defer body.Close()
	json.NewDecoder(respBody).Decode(&sshKeyList)
	//spew.Dump(body)
	//log.Println(sshKeyList["idList"])
	log.Println(sshKeyList.IDList)

	//return sshKeyList["idList"], model.WebStatus{StatusCode: respStatus}
	return sshKeyList.IDList, model.WebStatus{StatusCode: respStatus}
}

func GetSshKeyInfoListByOption(nameSpaceID string, optionParam string, filterKeyParam string, filterValParam string) ([]tbmcir.TbSshKeyInfo, model.WebStatus) {
	fmt.Println("GetSshKeyInfoList ************ : ")
	var originalUrl = "/ns/{nsId}/resources/sshKey"
	var paramMapper = make(map[string]string)
	paramMapper["{nsId}"] = nameSpaceID
	urlParam := util.MappingUrlParameter(originalUrl, paramMapper)

	if optionParam != "" {
		urlParam += "?option=" + optionParam
	}
	if filterKeyParam != "" {
		urlParam += "&filterKey=" + filterKeyParam
		urlParam += "&filterVal=" + filterValParam
	}
	url := util.TUMBLEBUG + urlParam
	// url := util.TUMBLEBUG + "/ns/" + nameSpaceID + "/resources/sshKey"

	//resp, err := util.CommonHttp(url, pbytes, http.MethodGet)
	resp, err := util.CommonHttpWithoutParam(url, http.MethodGet)

	if err != nil {
		fmt.Println(err)
		return nil, model.WebStatus{StatusCode: 500, Message: err.Error()}
	}
	// defer body.Close()
	respBody := resp.Body
	respStatus := resp.StatusCode

	// return respBody, respStatus
	log.Println(respBody)
	sshKeyList := map[string][]tbmcir.TbSshKeyInfo{}

	json.NewDecoder(respBody).Decode(&sshKeyList)
	//spew.Dump(body)
	fmt.Println(sshKeyList["sshKey"])

	return sshKeyList["sshKey"], model.WebStatus{StatusCode: respStatus}
}

// sshKey 상세 조회
func GetSshKeyData(nameSpaceID string, sshKeyID string) (*tbmcir.TbSshKeyInfo, model.WebStatus) {
	var originalUrl = "/ns/{nsId}/resources/sshKey/{sshKeyId}"
	var paramMapper = make(map[string]string)
	paramMapper["{nsId}"] = nameSpaceID
	paramMapper["{sshKeyId}"] = sshKeyID
	urlParam := util.MappingUrlParameter(originalUrl, paramMapper)
	url := util.TUMBLEBUG + urlParam
	// url := util.TUMBLEBUG + "/ns/" + nameSpaceID + "/resources/sshKey/" + sshKeyID

	fmt.Println("nameSpaceID : ", nameSpaceID)

	// pbytes, _ := json.Marshal(nameSpaceID)
	// body, err := util.CommonHttpGet(url)
	resp, err := util.CommonHttpWithoutParam(url, http.MethodGet)
	sshKeyInfo := tbmcir.TbSshKeyInfo{}
	if err != nil {
		fmt.Println(err)
		return &sshKeyInfo, model.WebStatus{StatusCode: 500, Message: err.Error()}
	}

	respBody := resp.Body
	respStatus := resp.StatusCode

	json.NewDecoder(respBody).Decode(&sshKeyInfo)
	fmt.Println(sshKeyInfo)

	return &sshKeyInfo, model.WebStatus{StatusCode: respStatus}
}

// sshKey 등록
// option=register 항목은 TB에서 자동으로 넣을 때 사용하는 param으로 webtool에서 사용하지 않음.
func RegSshKey(nameSpaceID string, sshKeyRegInfo *tbmcir.TbSshKeyReq) (*tbmcir.TbSshKeyInfo, model.WebStatus) {
	var originalUrl = "/ns/{nsId}/resources/sshKey"
	var paramMapper = make(map[string]string)
	paramMapper["{nsId}"] = nameSpaceID
	urlParam := util.MappingUrlParameter(originalUrl, paramMapper)

	url := util.TUMBLEBUG + urlParam
	// url := util.TUMBLEBUG + "/ns/" + nameSpaceID + "/resources/sshKey"

	// fmt.Println("vnetInfo : ", vnetInfo)

	pbytes, _ := json.Marshal(sshKeyRegInfo)
	fmt.Println(string(pbytes))
	resp, err := util.CommonHttp(url, pbytes, http.MethodPost)
	sshKeyInfo := tbmcir.TbSshKeyInfo{}
	if err != nil {
		fmt.Println(err)
		return &sshKeyInfo, model.WebStatus{StatusCode: 500, Message: err.Error()}
	}
	log.Println("resp = ", resp)
	respBody := resp.Body
	respStatus := resp.StatusCode
	log.Println("respBody = ", respBody)
	// respStatus := resp.Status
	// log.Println("respStatusCode = ", respStatusCode)
	// log.Println("respStatus = ", respStatus)

	// 응답에 생성한 객체값이 옴

	json.NewDecoder(respBody).Decode(&sshKeyInfo)
	fmt.Println(sshKeyInfo)
	// return respBody, respStatusCode
	return &sshKeyInfo, model.WebStatus{StatusCode: respStatus}
}

func UpdateSshKey(nameSpaceID string, sshKeyId string, sshKeyInfo *tbmcir.TbSshKeyInfo) (*tbmcir.TbSshKeyInfo, model.WebStatus) {
	var originalUrl = "/ns/{nsId}/resources/sshKey/{sshKeyId}"
	var paramMapper = make(map[string]string)
	paramMapper["{nsId}"] = nameSpaceID
	paramMapper["{sshKeyId}"] = sshKeyId
	urlParam := util.MappingUrlParameter(originalUrl, paramMapper)
	url := util.TUMBLEBUG + urlParam
	// url := util.TUMBLEBUG + "/ns/" + nameSpaceID + "/resources/sshKey"

	// fmt.Println("vnetInfo : ", vnetInfo)

	pbytes, _ := json.Marshal(sshKeyInfo)
	fmt.Println(string(pbytes))
	resp, err := util.CommonHttp(url, pbytes, http.MethodPut)
	sshKeyInfoResponse := tbmcir.TbSshKeyInfo{}
	if err != nil {
		fmt.Println(err)
		return &sshKeyInfoResponse, model.WebStatus{StatusCode: 500, Message: err.Error()}
	}
	log.Println("resp = ", resp)
	respBody := resp.Body
	respStatus := resp.StatusCode
	log.Println("respBody = ", respBody)

	json.NewDecoder(respBody).Decode(&sshKeyInfoResponse)
	fmt.Println(sshKeyInfoResponse)

	return &sshKeyInfoResponse, model.WebStatus{StatusCode: respStatus}
}

// sshKey 삭제
func DelSshKey(nameSpaceID string, sshKeyID string) (model.WebStatus, model.WebStatus) {
	webStatus := model.WebStatus{}
	// if ValidateString(sshKeyID) != nil {
	if len(sshKeyID) == 0 {
		log.Println("sshKeyID 가 없으면 해당 namespace의 모든 sshKeyID 삭제되므로 처리할 수 없습니다.")
		return webStatus, model.WebStatus{StatusCode: 4040, Message: "sshKeyID 가 없으면 해당 namespace의 모든 sshKeyID 삭제되므로 처리할 수 없습니다."}
	}

	var originalUrl = "/ns/{nsId}/resources/sshKey/{sshKeyId}"
	var paramMapper = make(map[string]string)
	paramMapper["{nsId}"] = nameSpaceID
	paramMapper["{sshKeyId}"] = sshKeyID
	urlParam := util.MappingUrlParameter(originalUrl, paramMapper)
	url := util.TUMBLEBUG + urlParam
	// url := util.TUMBLEBUG + "/ns/" + nameSpaceID + "/resources/sshKey/" + sshKeyID

	fmt.Println("sshKeyID : ", sshKeyID)

	pbytes, _ := json.Marshal(sshKeyID)
	resp, err := util.CommonHttp(url, pbytes, http.MethodDelete)

	if err != nil {
		fmt.Println(err)
		return webStatus, model.WebStatus{StatusCode: 500, Message: err.Error()}
	}
	fmt.Println("resp : ", resp)

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
	// return respBody, model.WebStatus{StatusCode: respStatus}
}

// 전체 sshKey 삭제
func DelAllSshKey(nameSpaceID string) (tbcommon.TbSimpleMsg, model.WebStatus) {
	var originalUrl = "/ns/{nsId}/resources/sshKey"
	var paramMapper = make(map[string]string)
	paramMapper["{nsId}"] = nameSpaceID

	urlParam := util.MappingUrlParameter(originalUrl, paramMapper)

	url := util.TUMBLEBUG + urlParam

	resp, err := util.CommonHttpWithoutParam(url, http.MethodDelete)

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

// VirtualMachineImage 목록 조회
func GetVirtualMachineImageInfoList(nameSpaceID string) ([]tbmcir.TbImageInfo, model.WebStatus) {
	fmt.Println("GetVirtualMachineImageInfoList ************ : ")
	// var originalUrl = "/ns/{nsId}/resources/image"
	var originalUrl = "/ns/{nsId}/resources/image"
	var paramMapper = make(map[string]string)
	paramMapper["{nsId}"] = nameSpaceID
	urlParam := util.MappingUrlParameter(originalUrl, paramMapper)
	url := util.TUMBLEBUG + urlParam
	// url := util.TUMBLEBUG + "/ns/" + nameSpaceID + "/resources/image"

	pbytes, _ := json.Marshal(nameSpaceID)
	resp, err := util.CommonHttp(url, pbytes, http.MethodGet)

	if err != nil {
		fmt.Println(err)
		return nil, model.WebStatus{StatusCode: 500, Message: err.Error()}
	}
	// TODO : defer를 넣어줘야 할 듯. defer body.Close()
	respBody := resp.Body
	respStatus := resp.StatusCode

	virtualMachineImageList := map[string][]tbmcir.TbImageInfo{}

	json.NewDecoder(respBody).Decode(&virtualMachineImageList)
	fmt.Println(virtualMachineImageList["image"])

	// robots, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Printf("%s", robots)

	return virtualMachineImageList["image"], model.WebStatus{StatusCode: respStatus}
}

// VirtualMachineImage 목록에서 Option으로 ID 목록만 가져오는 function
func GetVirtualMachineImageInfoListByID(nameSpaceID string, filterKeyParam string, filterValParam string) ([]string, model.WebStatus) {
	fmt.Println("GetVirtualMachineImageInfoListByID ************ : ")
	// var originalUrl = "/ns/{nsId}/resources/image"
	var originalUrl = "/ns/{nsId}/resources/image"
	var paramMapper = make(map[string]string)
	paramMapper["{nsId}"] = nameSpaceID

	// url := util.TUMBLEBUG + urlParam
	urlParam := util.MappingUrlParameter(originalUrl, paramMapper)

	//if optionParam != ""{
	//	urlParam += "?option=" + optionParam
	//}
	urlParam += "?option=id"
	if filterKeyParam != "" {
		urlParam += "&filterKey=" + filterKeyParam
		urlParam += "&filterVal=" + filterValParam
	}
	url := util.TUMBLEBUG + urlParam
	//url := util.TUMBLEBUG + urlParam + optionParamVal
	// url := util.TUMBLEBUG + "/ns/" + nameSpaceID + "/resources/image"

	pbytes, _ := json.Marshal(nameSpaceID)
	resp, err := util.CommonHttp(url, pbytes, http.MethodGet)

	if err != nil {
		fmt.Println(err)
		return nil, model.WebStatus{StatusCode: 500, Message: err.Error()}
	}
	// TODO : defer를 넣어줘야 할 듯. defer body.Close()
	respBody := resp.Body
	respStatus := resp.StatusCode

	// virtualMachineImageIdList := map[string][]tbcommon.TbIdList{}
	virtualMachineImageIdList := tbcommon.TbIdList{}

	json.NewDecoder(respBody).Decode(&virtualMachineImageIdList)
	//fmt.Println(virtualMachineImageIdList.IDList)

	// robots, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Printf("%s", robots)

	return virtualMachineImageIdList.IDList, model.WebStatus{StatusCode: respStatus}
}

func GetVirtualMachineImageInfoListByOption(nameSpaceID string, optionParam string, filterKeyParam string, filterValParam string) ([]tbmcir.TbImageInfo, model.WebStatus) {
	fmt.Println("GetVirtualMachineImageInfoListByOption ************ : ")
	// var originalUrl = "/ns/{nsId}/resources/image"
	var originalUrl = "/ns/{nsId}/resources/image"
	var paramMapper = make(map[string]string)
	paramMapper["{nsId}"] = nameSpaceID

	// url := util.TUMBLEBUG + urlParam
	urlParam := util.MappingUrlParameter(originalUrl, paramMapper)

	if optionParam != "" {
		urlParam += "?option=" + optionParam
	} else {
		urlParam += "?option="
	}
	if filterKeyParam != "" {
		urlParam += "&filterKey=" + filterKeyParam
		urlParam += "&filterVal=" + filterValParam
	}
	url := util.TUMBLEBUG + urlParam
	//url := util.TUMBLEBUG + urlParam + optionParamVal
	// url := util.TUMBLEBUG + "/ns/" + nameSpaceID + "/resources/image"

	pbytes, _ := json.Marshal(nameSpaceID)
	resp, err := util.CommonHttp(url, pbytes, http.MethodGet)

	if err != nil {
		fmt.Println(err)
		return nil, model.WebStatus{StatusCode: 500, Message: err.Error()}
	}
	// TODO : defer를 넣어줘야 할 듯. defer body.Close()
	respBody := resp.Body
	respStatus := resp.StatusCode

	virtualMachineImageList := map[string][]tbmcir.TbImageInfo{}

	json.NewDecoder(respBody).Decode(&virtualMachineImageList)
	fmt.Println(virtualMachineImageList["image"])

	// robots, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Printf("%s", robots)

	return virtualMachineImageList["image"], model.WebStatus{StatusCode: respStatus}
}

// VirtualMachineImage 상세 조회
func GetVirtualMachineImageData(nameSpaceID string, virtualMachineImageID string) (*tbmcir.TbImageInfo, model.WebStatus) {
	fmt.Println("GetVirtualMachineImageData ************ : ")
	var originalUrl = "/ns/{nsId}/resources/image/{imageId}"
	var paramMapper = make(map[string]string)
	paramMapper["{nsId}"] = nameSpaceID
	paramMapper["{imageId}"] = virtualMachineImageID
	urlParam := util.MappingUrlParameter(originalUrl, paramMapper)
	url := util.TUMBLEBUG + urlParam
	// url := util.TUMBLEBUG + "/ns/" + nameSpaceID + "/resources/image/" + virtualMachineImageID

	resp, err := util.CommonHttpWithoutParam(url, http.MethodGet)
	virtualMachineImageInfo := tbmcir.TbImageInfo{}
	if err != nil {
		fmt.Println(err)
		return &virtualMachineImageInfo, model.WebStatus{StatusCode: 500, Message: err.Error()}
	}

	respBody := resp.Body
	respStatus := resp.StatusCode

	json.NewDecoder(respBody).Decode(&virtualMachineImageInfo)
	fmt.Println(virtualMachineImageInfo)

	return &virtualMachineImageInfo, model.WebStatus{StatusCode: respStatus}
}

func UpdateVirtualMachineImage(nameSpaceID string, virtualMachineImageID string, imageInfo *tbmcir.TbImageInfo) (*tbmcir.TbImageInfo, model.WebStatus) {
	fmt.Println("UpdateVirtualMachineImageData ************ : ")
	var originalUrl = "/ns/{nsId}/resources/image/{imageId}"
	var paramMapper = make(map[string]string)
	paramMapper["{nsId}"] = nameSpaceID
	paramMapper["{imageId}"] = virtualMachineImageID
	urlParam := util.MappingUrlParameter(originalUrl, paramMapper)
	url := util.TUMBLEBUG + urlParam
	// url := util.TUMBLEBUG + "/ns/" + nameSpaceID + "/resources/image/" + virtualMachineImageID

	pbytes, _ := json.Marshal(imageInfo) // action=registerWithInfo, registerWithId param이 regInfo안에 모두 있으므로 별도로 나누어 호출하지않고 그냥 사용
	fmt.Println(string(pbytes))
	resp, err := util.CommonHttp(url, pbytes, http.MethodPut)
	virtualMachineImageInfo := tbmcir.TbImageInfo{}
	if err != nil {
		fmt.Println(err)
		return &virtualMachineImageInfo, model.WebStatus{StatusCode: 500, Message: err.Error()}
	}

	respBody := resp.Body
	respStatus := resp.StatusCode

	json.NewDecoder(respBody).Decode(&virtualMachineImageInfo)
	fmt.Println(virtualMachineImageInfo)

	return &virtualMachineImageInfo, model.WebStatus{StatusCode: respStatus}
}

// VirtualMachineImage 등록 registeringMethod = imageID
func RegVirtualMachineImage(nameSpaceID string, registType string, virtualMachineImageRegInfo *tbmcir.TbImageReq) (*tbmcir.TbImageInfo, model.WebStatus) {
	fmt.Println("RegVirtualMachineImage ************ : ")
	if registType == "" {
		registType = "registerWithId" // registerWithId 또는 registerWithInfo
	}

	// API에는 registeringMethod로 표현되어있으나 실제로는 action임.
	var originalUrl = "/ns/{nsId}/resources/image?action=registerWithId"
	var paramMapper = make(map[string]string)
	paramMapper["{nsId}"] = nameSpaceID
	paramMapper["{registType}"] = registType
	urlParam := util.MappingUrlParameter(originalUrl, paramMapper)
	url := util.TUMBLEBUG + urlParam
	// url := util.TUMBLEBUG + "/ns/" + nameSpaceID + "/resources/image?action=registerWithInfo" //
	// url := util.TUMBLEBUG + "/ns/" + nameSpaceID + "/resources/image?action=registerWithId"//

	pbytes, _ := json.Marshal(virtualMachineImageRegInfo) // action=registerWithInfo, registerWithId param이 regInfo안에 모두 있으므로 별도로 나누어 호출하지않고 그냥 사용
	fmt.Println(string(pbytes))
	resp, err := util.CommonHttp(url, pbytes, http.MethodPost)
	virtualMachineImageInfo := tbmcir.TbImageInfo{}
	if err != nil {
		fmt.Println(err)
		return &virtualMachineImageInfo, model.WebStatus{StatusCode: 500, Message: err.Error()}
	}

	// data, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Printf("%s\n", string(data))

	respBody := resp.Body
	respStatus := resp.StatusCode
	// respStatus := resp.Status
	// log.Println("respStatusCode = ", respStatusCode)
	// log.Println("respStatus = ", respStatus)

	// 응답에 생성한 객체값이 옴

	json.NewDecoder(respBody).Decode(&virtualMachineImageInfo)
	fmt.Println(virtualMachineImageInfo)
	// return respBody, respStatusCode
	return &virtualMachineImageInfo, model.WebStatus{StatusCode: respStatus}
}

// VirtualMachineImage 등록 registeringMethod = imageID
// 생성시 action에 regist
func RegVirtualMachineImageWithInfo(nameSpaceID string, registType string, virtualMachineImageRegInfo *tbmcir.TbImageInfo) (*tbmcir.TbImageInfo, model.WebStatus) {
	fmt.Println("RegVirtualMachineImage ************ : ")
	if registType == "" {
		registType = "registerWithId" // registerWithId 또는 registerWithInfo
	}

	var originalUrl = "/ns/{nsId}/resources/image?registeringMethod=registerWithInfo"
	var paramMapper = make(map[string]string)
	paramMapper["{nsId}"] = nameSpaceID
	paramMapper["{registType}"] = registType
	urlParam := util.MappingUrlParameter(originalUrl, paramMapper)
	url := util.TUMBLEBUG + urlParam
	// url := util.TUMBLEBUG + "/ns/" + nameSpaceID + "/resources/image?action=registerWithInfo" //
	// url := util.TUMBLEBUG + "/ns/" + nameSpaceID + "/resources/image?action=registerWithId"//

	pbytes, _ := json.Marshal(virtualMachineImageRegInfo) // action=registerWithInfo, registerWithId param이 regInfo안에 모두 있으므로 별도로 나누어 호출하지않고 그냥 사용
	fmt.Println(string(pbytes))
	resp, err := util.CommonHttp(url, pbytes, http.MethodPost)
	virtualMachineImageInfo := tbmcir.TbImageInfo{}
	if err != nil {
		fmt.Println(err)
		return &virtualMachineImageInfo, model.WebStatus{StatusCode: 500, Message: err.Error()}
	}

	// data, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Printf("%s\n", string(data))

	respBody := resp.Body
	respStatus := resp.StatusCode
	// respStatus := resp.Status
	// log.Println("respStatusCode = ", respStatusCode)
	// log.Println("respStatus = ", respStatus)

	// 응답에 생성한 객체값이 옴

	json.NewDecoder(respBody).Decode(&virtualMachineImageInfo)
	fmt.Println(virtualMachineImageInfo)
	// return respBody, respStatusCode
	return &virtualMachineImageInfo, model.WebStatus{StatusCode: respStatus}
}

// 해당 namespace의 모든 VirtualMachineImage 삭제
func DelAllVirtualMachineImage(nameSpaceID string) (model.WebStatus, model.WebStatus) {
	// if ValidateString(VirtualMachineImageID) != nil {
	var originalUrl = "/ns/{nsId}/resources/image"
	var paramMapper = make(map[string]string)
	paramMapper["{nsId}"] = nameSpaceID
	urlParam := util.MappingUrlParameter(originalUrl, paramMapper)
	url := util.TUMBLEBUG + urlParam
	// url := util.TUMBLEBUG + "/ns/" + nameSpaceID + "/resources/image"

	// resp, err := util.CommonHttp(url, pbytes, http.MethodDelete)
	resp, err := util.CommonHttpWithoutParam(url, http.MethodDelete)
	webStatus := model.WebStatus{}
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
	// return respBody, model.WebStatus{StatusCode: respStatus}
}

// 해당 namespace의 특정 VirtualMachineImage 삭제
func DelVirtualMachineImage(nameSpaceID string, virtualMachineImageID string) (model.WebStatus, model.WebStatus) {
	webStatus := model.WebStatus{}
	// if ValidateString(VirtualMachineImageID) != nil {
	if len(virtualMachineImageID) == 0 {
		log.Println("ImageID 가 없으면 해당 namespace의 모든 image가 삭제되므로 처리할 수 없습니다.")
		return webStatus, model.WebStatus{StatusCode: 4040, Message: "ImageID 가 없으면 해당 namespace의 모든 image가 삭제되므로 처리할 수 없습니다."}
	}

	var originalUrl = "/ns/{nsId}/resources/image/{imageId}"
	var paramMapper = make(map[string]string)
	paramMapper["{nsId}"] = nameSpaceID
	paramMapper["{imageId}"] = virtualMachineImageID
	urlParam := util.MappingUrlParameter(originalUrl, paramMapper)
	url := util.TUMBLEBUG + urlParam
	// url := util.TUMBLEBUG + "/ns/" + nameSpaceID + "/resources/image/" + virtualMachineImageID

	pbytes, _ := json.Marshal(virtualMachineImageID)
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
	// return respBody, model.WebStatus{StatusCode: respStatus}
}

// 자신의 provider에 등록된 resource 조회
func GetInspectResourceList(inspectResource *tbcommon.RestInspectResourcesRequest) (*tbmcis.InspectResource, model.WebStatus) {
	fmt.Println("GetInspectResourceList ************ : ")
	//https://www.javaer101.com/ko/article/5704925.html 참조 : 값이 있는 것만 넘기기
	var originalUrl = "/inspectResources"
	urlParam := util.MappingUrlParameter(originalUrl, nil)
	url := util.TUMBLEBUG + urlParam

	pbytes, _ := json.Marshal(inspectResource)
	resp, err := util.CommonHttp(url, pbytes, http.MethodPost)

	respBody := resp.Body
	respStatus := resp.StatusCode

	inspectResourcesResponse := tbmcis.InspectResource{}
	if err != nil {
		failResultInfo := tbcommon.TbSimpleMsg{}
		json.NewDecoder(respBody).Decode(&failResultInfo)
		log.Println(failResultInfo)
		//return &inspectResourcesResponse, model.WebStatus{StatusCode: 500, Message: err.Error()}
		return &inspectResourcesResponse, model.WebStatus{StatusCode: 500, Message: failResultInfo.Message}
	}

	json.NewDecoder(respBody).Decode(&inspectResourcesResponse)
	fmt.Println(inspectResourcesResponse)

	return &inspectResourcesResponse, model.WebStatus{StatusCode: respStatus}

}

/*
	CSP와 Tumblebug에 등록된 모든 리소스 비교
	전체이므로 별도의 parameter 없음.
*/
func GetInspectResourcesOverview() (*tbmcis.InspectResourceAllResult, model.WebStatus) {
	fmt.Println("Inspect Resources Overview (vNet, securityGroup, sshKey, vm) registered in CB-Tumblebug and CSP for all connections")
	var originalUrl = "/inspectResourcesOverview"
	urlParam := util.MappingUrlParameter(originalUrl, nil)

	url := util.TUMBLEBUG + urlParam
	// url := util.TUMBLEBUG + "/ns"

	resp, err := util.CommonHttp(url, nil, http.MethodGet)

	respBody := resp.Body
	respStatus := resp.StatusCode

	inspectResourceAllResult := tbmcis.InspectResourceAllResult{}
	if err != nil {
		failResultInfo := tbcommon.TbSimpleMsg{}
		json.NewDecoder(respBody).Decode(&failResultInfo)
		log.Println(failResultInfo)
		//return &inspectResourcesResponse, model.WebStatus{StatusCode: 500, Message: err.Error()}
		return &inspectResourceAllResult, model.WebStatus{StatusCode: 500, Message: failResultInfo.Message}
	}

	json.NewDecoder(respBody).Decode(&inspectResourceAllResult)
	fmt.Println(inspectResourceAllResult)

	return &inspectResourceAllResult, model.WebStatus{StatusCode: respStatus}
}

func RegCspResources(resourcesRequest *tbcommon.RestRegisterCspNativeResourcesRequest, optionParam string) (*tbmcis.RegisterResourceResult, model.WebStatus) {
	var originalUrl = "/registerCspResources"
	urlParam := util.MappingUrlParameter(originalUrl, nil)
	if optionParam != "" {
		urlParam += "?option=" + optionParam
	}
	url := util.TUMBLEBUG + urlParam

	pbytes, _ := json.Marshal(resourcesRequest)
	resp, err := util.CommonHttp(url, pbytes, http.MethodPost)

	respBody := resp.Body
	respStatus := resp.StatusCode

	registerResourceResult := tbmcis.RegisterResourceResult{}
	if err != nil {
		failResultInfo := tbcommon.TbSimpleMsg{}
		json.NewDecoder(respBody).Decode(&failResultInfo)
		log.Println(failResultInfo)
		//return &inspectResourcesResponse, model.WebStatus{StatusCode: 500, Message: err.Error()}
		return &registerResourceResult, model.WebStatus{StatusCode: 500, Message: failResultInfo.Message}
	}

	json.NewDecoder(respBody).Decode(registerResourceResult)
	fmt.Println(registerResourceResult)

	return &registerResourceResult, model.WebStatus{StatusCode: respStatus}

}

func RegCspResourcesAll(resourcesRequest *tbcommon.RestRegisterCspNativeResourcesRequestAll, optionParam string) (*tbmcis.RegisterResourceAllResult, model.WebStatus) {
	var originalUrl = "/registerCspResourcesAll"
	urlParam := util.MappingUrlParameter(originalUrl, nil)
	if optionParam != "" {
		urlParam += "?option=" + optionParam
	}
	url := util.TUMBLEBUG + urlParam

	pbytes, _ := json.Marshal(resourcesRequest)
	resp, err := util.CommonHttp(url, pbytes, http.MethodPost)

	respBody := resp.Body
	respStatus := resp.StatusCode

	registerResourceResult := tbmcis.RegisterResourceAllResult{}
	if err != nil {
		failResultInfo := tbcommon.TbSimpleMsg{}
		json.NewDecoder(respBody).Decode(&failResultInfo)
		log.Println(failResultInfo)
		//return &inspectResourcesResponse, model.WebStatus{StatusCode: 500, Message: err.Error()}
		return &registerResourceResult, model.WebStatus{StatusCode: 500, Message: failResultInfo.Message}
	}

	json.NewDecoder(respBody).Decode(registerResourceResult)
	fmt.Println(registerResourceResult)

	return &registerResourceResult, model.WebStatus{StatusCode: respStatus}

}

func GetLoadCommonResource() (tbcommon.TbSimpleMsg, model.WebStatus) {
	fmt.Println("Load Common Resources from internal asset files (Spec, Image)")
	var originalUrl = "/loadCommonResource"
	urlParam := util.MappingUrlParameter(originalUrl, nil)

	url := util.TUMBLEBUG + urlParam
	// url := util.TUMBLEBUG + "/ns"

	resp, err := util.CommonHttp(url, nil, http.MethodGet)

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

// VM Image 조회
func LookupVirtualMachineImageList(connectionName string) (tbmcir.SpiderImageInfos, model.WebStatus) {
	fmt.Println("LookupVirtualMachineImageList ************ : ", connectionName)
	var originalUrl = "/lookupImages"
	urlParam := util.MappingUrlParameter(originalUrl, nil)
	url := util.TUMBLEBUG + urlParam
	// url := util.TUMBLEBUG + "/lookupImage"

	// body, err := util.CommonHttpGet(url)
	// resp, err := util.CommonHttpWithoutParam(url, http.MethodGet)
	//common.TbConnectionName
	paramMap := map[string]string{"connectionName": connectionName}

	pbytes, _ := json.Marshal(paramMap)
	log.Println(string(pbytes))
	resp, err := util.CommonHttp(url, pbytes, http.MethodPost)
	log.Println("LookupVirtualMachineImageList called 1 ")
	if err != nil {
		fmt.Println(err)
		return nil, model.WebStatus{StatusCode: 500, Message: err.Error()}
	}
	// defer body.Close()
	respBody := resp.Body
	respStatus := resp.StatusCode
	log.Println("LookupVirtualMachineImageList called 2 ", respStatus)
	// return respBody, respStatus
	// log.Println(respBody)
	lookupImageList := map[string]tbmcir.SpiderImageInfos{}

	json.NewDecoder(respBody).Decode(&lookupImageList)
	log.Println("LookupVirtualMachineImageList called 3 ")

	return lookupImageList["image"], model.WebStatus{StatusCode: respStatus}
}

// deprecated :  imageId 받는 것에서 connection 받는 것 까지로 변경 됨
// func LookupVirtualMachineImageData(virtualMachineImageID string) (*tbmcir.TbImageInfo, model.WebStatus) {
// 	var originalUrl = "/lookupImage/{imageId}"
// 	var paramMapper = make(map[string]string)
// 	paramMapper["{imageId}"] = virtualMachineImageID
// 	urlParam := util.MappingUrlParameter(originalUrl, paramMapper)
// 	url := util.TUMBLEBUG + urlParam
// 	// url := util.TUMBLEBUG + "/lookupImage/" + virtualMachineImageID

// 	// pbytes, _ := json.Marshal(nameSpaceID)
// 	resp, err := util.CommonHttpWithoutParam(url, http.MethodGet)
// 	virtualMachineImageInfo := tbmcir.TbImageInfo{}
// 	if err != nil {
// 		fmt.Println(err)
// 		return &virtualMachineImageInfo, model.WebStatus{StatusCode: 500, Message: err.Error()}
// 	}

// 	respBody := resp.Body
// 	respStatus := resp.StatusCode

// 	json.NewDecoder(respBody).Decode(virtualMachineImageInfo)
// 	fmt.Println(virtualMachineImageInfo)

// 	return &virtualMachineImageInfo, model.WebStatus{StatusCode: respStatus}
// }

// 특정 provider의 특정 image정보 조회
func LookupVirtualMachineImageData(restLookupImageRequest *tbmcir.RestLookupImageRequest) (*tbmcir.SpiderImageInfo, model.WebStatus) {
	var originalUrl = "/lookupImage"
	var paramMapper = make(map[string]string)
	urlParam := util.MappingUrlParameter(originalUrl, paramMapper)
	url := util.TUMBLEBUG + urlParam

	pbytes, _ := json.Marshal(restLookupImageRequest)
	//resp, err := util.CommonHttp(url, pbytes, http.MethodGet)
	resp, err := util.CommonHttp(url, pbytes, http.MethodPost)

	virtualMachineImageInfo := tbmcir.SpiderImageInfo{}
	if err != nil {
		fmt.Println(err)
		return &virtualMachineImageInfo, model.WebStatus{StatusCode: 500, Message: err.Error()}
	}

	respBody := resp.Body
	respStatus := resp.StatusCode

	json.NewDecoder(respBody).Decode(virtualMachineImageInfo)
	fmt.Println(virtualMachineImageInfo)

	return &virtualMachineImageInfo, model.WebStatus{StatusCode: respStatus}
}

// csp에 등록된 정보조회.
func FetchVirtualMachineImageList(nameSpaceID string) ([]tbcommon.TbSimpleMsg, model.WebStatus) {
	fmt.Println("FetchVirtualMachineImageList ************ : ", nameSpaceID)
	var originalUrl = "/ns/{nsId}/resources/fetchImages"
	var paramMapper = make(map[string]string)
	paramMapper["{nsId}"] = nameSpaceID
	urlParam := util.MappingUrlParameter(originalUrl, paramMapper)
	url := util.TUMBLEBUG + urlParam
	// url := util.TUMBLEBUG + "/ns/" + nameSpaceID + "/resources/fetchImages"

	resp, err := util.CommonHttp(url, nil, http.MethodPost)
	// resp, err := util.CommonHttpWithoutParam(url, http.MethodPost)
	if err != nil {
		fmt.Println(err)
		return nil, model.WebStatus{StatusCode: 500, Message: err.Error()}
	}
	// defer body.Close()
	respBody := resp.Body
	respStatus := resp.StatusCode
	fetchImageList := map[string]tbmcir.SpiderImageInfos{}

	json.NewDecoder(respBody).Decode(&fetchImageList)
	log.Println("FetchVirtualMachineImageList called ")

	// return fetchImageList["image"], model.WebStatus{StatusCode: respStatus}
	return nil, model.WebStatus{StatusCode: respStatus} // TODO : simpleMsg return 하도록 변경할 것
}

// VirtualMachineImage 상세 조회
func SearchVirtualMachineImageList(nameSpaceID string, restSearchImageRequest *tbmcir.RestSearchImageRequest) ([]tbmcir.TbImageInfo, model.WebStatus) {
	var originalUrl = "/ns/{nsId}/resources/searchImage"
	var paramMapper = make(map[string]string)
	paramMapper["{nsId}"] = nameSpaceID
	urlParam := util.MappingUrlParameter(originalUrl, paramMapper)
	url := util.TUMBLEBUG + urlParam
	// url := util.TUMBLEBUG + "/ns/" + nameSpaceID + "/resources/searchImage/"

	// pbytes, _ := json.Marshal(nameSpaceID)
	pbytes, _ := json.Marshal(restSearchImageRequest)
	// resp, err := util.CommonHttpWithoutParam(url, http.MethodGet)
	resp, err := util.CommonHttp(url, pbytes, http.MethodPost)
	virtualMachineImageInfo := tbmcir.RestGetAllImageResponse{}
	if err != nil {
		fmt.Println(err)
		return virtualMachineImageInfo.Image, model.WebStatus{StatusCode: 500, Message: err.Error()}
	}

	respBody := resp.Body
	respStatus := resp.StatusCode

	json.NewDecoder(respBody).Decode(&virtualMachineImageInfo)
	fmt.Println(virtualMachineImageInfo)

	return virtualMachineImageInfo.Image, model.WebStatus{StatusCode: respStatus}
}

// VMSpec 목록 조회
func GetVmSpecInfoList(nameSpaceID string) ([]tbmcir.TbSpecInfo, model.WebStatus) {
	fmt.Println("GetVMSpecInfoList ************ : ")
	var originalUrl = "/ns/{nsId}/resources/spec"
	var paramMapper = make(map[string]string)
	paramMapper["{nsId}"] = nameSpaceID
	urlParam := util.MappingUrlParameter(originalUrl, paramMapper)
	url := util.TUMBLEBUG + urlParam
	// url := util.TUMBLEBUG + "/ns/" + nameSpaceID + "/resources/spec"

	// pbytes, _ := json.Marshal(nameSpaceID)
	// resp, err := util.CommonHttp(url, pbytes, http.MethodGet)
	resp, err := util.CommonHttpWithoutParam(url, http.MethodGet)
	if err != nil {
		fmt.Println(err)
		return nil, model.WebStatus{StatusCode: 500, Message: err.Error()}
	}
	// TODO : defer를 넣어줘야 할 듯. defer body.Close()
	respBody := resp.Body
	respStatus := resp.StatusCode

	// return respBody, respStatus
	log.Println(respBody)
	vmSpecList := tbmcir.RestGetAllSpecResponse{}

	json.NewDecoder(respBody).Decode(&vmSpecList)
	//spew.Dump(body)
	fmt.Println(vmSpecList.Spec)

	return vmSpecList.Spec, model.WebStatus{StatusCode: respStatus}
}

func GetVmSpecInfoListByID(nameSpaceID string, filterKeyParam string, filterValParam string) ([]string, model.WebStatus) {
	fmt.Println("GetVMSpecInfoList ************ : ")
	var originalUrl = "/ns/{nsId}/resources/spec"
	var paramMapper = make(map[string]string)
	paramMapper["{nsId}"] = nameSpaceID

	urlParam := util.MappingUrlParameter(originalUrl, paramMapper)
	// url := util.TUMBLEBUG + urlParam

	//if optionParam != ""{
	//	urlParam += "?option=" + optionParam
	//}
	urlParam += "?option=id"
	if filterKeyParam != "" {
		urlParam += "&filterKey=" + filterKeyParam
		urlParam += "&filterVal=" + filterValParam
	}
	url := util.TUMBLEBUG + urlParam
	// url := util.TUMBLEBUG + "/ns/" + nameSpaceID + "/resources/spec"

	// pbytes, _ := json.Marshal(nameSpaceID)
	// resp, err := util.CommonHttp(url, pbytes, http.MethodGet)
	resp, err := util.CommonHttpWithoutParam(url, http.MethodGet)
	if err != nil {
		fmt.Println(err)
		return nil, model.WebStatus{StatusCode: 500, Message: err.Error()}
	}
	// TODO : defer를 넣어줘야 할 듯. defer body.Close()
	respBody := resp.Body
	respStatus := resp.StatusCode

	// return respBody, respStatus
	log.Println(respBody)
	vmSpecList := tbcommon.TbIdList{}

	json.NewDecoder(respBody).Decode(&vmSpecList)
	//spew.Dump(body)
	fmt.Println(vmSpecList.IDList)

	return vmSpecList.IDList, model.WebStatus{StatusCode: respStatus}
}

func GetVmSpecInfoListByOption(nameSpaceID string, optionParam string, filterKeyParam string, filterValParam string) ([]tbmcir.TbSpecInfo, model.WebStatus) {
	fmt.Println("GetVMSpecInfoList ************ : ")
	var originalUrl = "/ns/{nsId}/resources/spec"
	var paramMapper = make(map[string]string)
	paramMapper["{nsId}"] = nameSpaceID

	urlParam := util.MappingUrlParameter(originalUrl, paramMapper)

	if optionParam != "" {
		urlParam += "?option=" + optionParam
	} else {
		urlParam += "?option="
	}
	if filterKeyParam != "" {
		urlParam += "&filterKey=" + filterKeyParam
		urlParam += "&filterVal=" + filterValParam
	}
	url := util.TUMBLEBUG + urlParam
	// url := util.TUMBLEBUG + "/ns/" + nameSpaceID + "/resources/spec"

	// pbytes, _ := json.Marshal(nameSpaceID)
	// resp, err := util.CommonHttp(url, pbytes, http.MethodGet)
	resp, err := util.CommonHttpWithoutParam(url, http.MethodGet)
	if err != nil {
		fmt.Println(err)
		return nil, model.WebStatus{StatusCode: 500, Message: err.Error()}
	}
	// TODO : defer를 넣어줘야 할 듯. defer body.Close()
	respBody := resp.Body
	respStatus := resp.StatusCode

	// return respBody, respStatus
	log.Println(respBody)
	vmSpecList := tbmcir.RestGetAllSpecResponse{}

	json.NewDecoder(respBody).Decode(&vmSpecList)
	//spew.Dump(body)
	fmt.Println(vmSpecList.Spec)

	return vmSpecList.Spec, model.WebStatus{StatusCode: respStatus}
}

// VMSpec 상세 조회
func GetVmSpecInfoData(nameSpaceID string, vmSpecID string) (*tbmcir.TbSpecInfo, model.WebStatus) {
	fmt.Println("GetVMSpecInfoData ************ : ")
	var originalUrl = "/ns/{nsId}/resources/spec/{specId}"
	var paramMapper = make(map[string]string)
	paramMapper["{nsId}"] = nameSpaceID
	paramMapper["{specId}"] = vmSpecID
	urlParam := util.MappingUrlParameter(originalUrl, paramMapper)
	url := util.TUMBLEBUG + urlParam
	// url := util.TUMBLEBUG + "/ns/" + nameSpaceID + "/resources/spec/" + vmSpecID

	// pbytes, _ := json.Marshal(nameSpaceID)
	resp, err := util.CommonHttpWithoutParam(url, http.MethodGet)
	vmSpecInfo := tbmcir.TbSpecInfo{}
	if err != nil {
		fmt.Println(err)
		return &vmSpecInfo, model.WebStatus{StatusCode: 500, Message: err.Error()}
	}

	respBody := resp.Body
	respStatus := resp.StatusCode

	json.NewDecoder(respBody).Decode(&vmSpecInfo)
	fmt.Println(vmSpecInfo)

	return &vmSpecInfo, model.WebStatus{StatusCode: respStatus}
}

// VMSpecInfo 등록
func RegVmSpec(nameSpaceID string, specregisteringMethod string, vmSpecRegInfo *tbmcir.TbSpecReq) (*tbmcir.TbSpecInfo, model.WebStatus) {
	fmt.Println("RegVMSpec ************ : ")
	if specregisteringMethod == "" {
		specregisteringMethod = "registerWithInfo" // registerWithInfo or Else 이므로 registerWithInfo 를 넣거나 아니거나.
	}

	// else인 경우에는 4개의 parameter만 있음{
	// 	"connectionName": "string",
	// 	"cspSpecName": "string",
	// 	"description": "string",
	// 	"name": "string"
	//   }
	var originalUrl = "/ns/{nsId}/resources/spec?registeringMethod={specregisteringMethod}"
	var paramMapper = make(map[string]string)
	paramMapper["{nsId}"] = nameSpaceID
	paramMapper["{specregisteringMethod}"] = specregisteringMethod
	urlParam := util.MappingUrlParameter(originalUrl, paramMapper)
	url := util.TUMBLEBUG + urlParam
	// "https://localhost:1323/tumblebug/ns/ns01/resources/spec?registeringMethod=registerWithInfo"
	// url := util.TUMBLEBUG + "/ns/" + nameSpaceID + "/resources/spec?action=registerWithInfo"// parameter를 모두 받지않기 때문에 param의 data type이 틀려 오류남.
	// url := util.TUMBLEBUG + "/ns/" + nameSpaceID + "/resources/spec" // 그래서 action 인자없이 전송

	pbytes, _ := json.Marshal(vmSpecRegInfo)
	resp, err := util.CommonHttp(url, pbytes, http.MethodPost)

	vmSpecInfo := tbmcir.TbSpecInfo{}
	if err != nil {
		fmt.Println(err)
		return &vmSpecInfo, model.WebStatus{StatusCode: 500, Message: err.Error()}
	}

	respBody := resp.Body
	respStatus := resp.StatusCode

	// 응답에 생성한 객체값이 옴
	returnStatus := model.WebStatus{}
	if respStatus != 200 && respStatus != 201 { // 호출은 정상이나, 가져온 결과값이 200, 201아닌 경우 message에 담겨있는 것을 WebStatus에 set
		errorInfo := model.ErrorInfo{}
		json.NewDecoder(respBody).Decode(&errorInfo)
		fmt.Println("respStatus != 200 reason ", errorInfo)
		returnStatus.Message = errorInfo.Message
	} else {
		json.NewDecoder(respBody).Decode(&vmSpecInfo)
		fmt.Println(vmSpecInfo)
	}
	returnStatus.StatusCode = respStatus

	// return respBody, respStatusCode
	return &vmSpecInfo, returnStatus
}

// specRegisteringMethod에 따라 requestMethod가 다르므로 function 분리 함
func RegVmSpecWithInfo(nameSpaceID string, specregisteringMethod string, vmSpecRegInfo *tbmcir.TbSpecInfo) (*tbmcir.TbSpecInfo, model.WebStatus) {
	fmt.Println("RegVMSpec ************ : ")
	if specregisteringMethod == "" {
		specregisteringMethod = "registerWithInfo" // registerWithInfo or Else 이므로 registerWithInfo 를 넣거나 아니거나.
	}

	// else인 경우에는 4개의 parameter만 있음{
	// 	"connectionName": "string",
	// 	"cspSpecName": "string",
	// 	"description": "string",
	// 	"name": "string"
	//   }
	//var originalUrl = "/ns/{nsId}/resources/spec?registeringMethod={specregisteringMethod}"
	var originalUrl = "/ns/{nsId}/resources/spec?registeringMethod={specregisteringMethod}"

	var paramMapper = make(map[string]string)
	paramMapper["{nsId}"] = nameSpaceID
	paramMapper["{specregisteringMethod}"] = specregisteringMethod
	urlParam := util.MappingUrlParameter(originalUrl, paramMapper)
	url := util.TUMBLEBUG + urlParam
	// "https://localhost:1323/tumblebug/ns/ns01/resources/spec?registeringMethod=registerWithInfo"
	// url := util.TUMBLEBUG + "/ns/" + nameSpaceID + "/resources/spec?action=registerWithInfo"// parameter를 모두 받지않기 때문에 param의 data type이 틀려 오류남.
	// url := util.TUMBLEBUG + "/ns/" + nameSpaceID + "/resources/spec" // 그래서 action 인자없이 전송

	pbytes, _ := json.Marshal(vmSpecRegInfo)
	resp, err := util.CommonHttp(url, pbytes, http.MethodPost)

	vmSpecInfo := tbmcir.TbSpecInfo{}
	if err != nil {
		fmt.Println(err)
		return &vmSpecInfo, model.WebStatus{StatusCode: 500, Message: err.Error()}
	}

	respBody := resp.Body
	respStatus := resp.StatusCode

	// 응답에 생성한 객체값이 옴
	returnStatus := model.WebStatus{}
	if respStatus != 200 && respStatus != 201 { // 호출은 정상이나, 가져온 결과값이 200, 201아닌 경우 message에 담겨있는 것을 WebStatus에 set
		errorInfo := model.ErrorInfo{}
		json.NewDecoder(respBody).Decode(&errorInfo)
		fmt.Println("respStatus != 200 reason ", errorInfo)
		returnStatus.Message = errorInfo.Message
	} else {
		json.NewDecoder(respBody).Decode(&vmSpecInfo)
		fmt.Println(vmSpecInfo)
	}
	returnStatus.StatusCode = respStatus

	// return respBody, respStatusCode
	return &vmSpecInfo, returnStatus
}

func UpdateVMSpec(nameSpaceID string, vmSpecRegInfo *tbmcir.TbSpecInfo) (*tbmcir.TbSpecInfo, model.WebStatus) {
	fmt.Println("UpdateVMSpec ************ : ")
	vmSpecID := vmSpecRegInfo.ID
	var originalUrl = "/ns/{nsId}/resources/spec/{specId}"
	var paramMapper = make(map[string]string)
	paramMapper["{nsId}"] = nameSpaceID
	paramMapper["{specId}"] = vmSpecID
	urlParam := util.MappingUrlParameter(originalUrl, paramMapper)
	url := util.TUMBLEBUG + urlParam
	// url := util.TUMBLEBUG + "/ns/" + nameSpaceID + "/resources/spec"

	pbytes, _ := json.Marshal(vmSpecRegInfo)
	resp, err := util.CommonHttp(url, pbytes, http.MethodPut)
	vmSpecInfo := tbmcir.TbSpecInfo{}
	if err != nil {
		fmt.Println(err)
		return &vmSpecInfo, model.WebStatus{StatusCode: 500, Message: err.Error()}
	}

	respBody := resp.Body
	respStatus := resp.StatusCode

	json.NewDecoder(respBody).Decode(&vmSpecInfo)
	fmt.Println(vmSpecInfo)

	return &vmSpecInfo, model.WebStatus{StatusCode: respStatus}
}

// 해당 namespace의 모든 VMSpec 삭제 : TODO : 로그인 유저의 동일 namespace일 때만 삭제가능하도록
func DelAllVMSpec(nameSpaceID string) (model.WebStatus, model.WebStatus) {
	fmt.Println("DelAllVMSpec ************ : ")
	var originalUrl = "/ns/{nsId}/resources/spec"
	var paramMapper = make(map[string]string)
	paramMapper["{nsId}"] = nameSpaceID
	urlParam := util.MappingUrlParameter(originalUrl, paramMapper)
	url := util.TUMBLEBUG + urlParam
	// url := util.TUMBLEBUG + "/ns/" + nameSpaceID + "/resources/spec"

	resp, err := util.CommonHttp(url, nil, http.MethodDelete)
	webStatus := model.WebStatus{}
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

	// return respBody, model.WebStatus{StatusCode: respStatus}
}

// VMSpec 삭제
func DelVMSpec(nameSpaceID string, vmSpecID string) (model.WebStatus, model.WebStatus) {
	webStatus := model.WebStatus{}
	// if ValidateString(VMSpecID) != nil {
	if len(vmSpecID) == 0 {
		log.Println("specID 가 없으면 해당 namespace의 모든 image가 삭제되므로 처리할 수 없습니다.")
		return webStatus, model.WebStatus{StatusCode: 4040, Message: "specID 가 없으면 해당 namespace의 모든 image가 삭제되므로 처리할 수 없습니다."}
	}

	var originalUrl = "/ns/{nsId}/resources/spec/{specId}"
	var paramMapper = make(map[string]string)
	paramMapper["{nsId}"] = nameSpaceID
	paramMapper["{specId}"] = vmSpecID
	urlParam := util.MappingUrlParameter(originalUrl, paramMapper)
	url := util.TUMBLEBUG + urlParam
	// url := util.TUMBLEBUG + "/ns/" + nameSpaceID + "/resources/spec/" + vmSpecID

	pbytes, _ := json.Marshal(vmSpecID)
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
	// return respBody, model.WebStatus{StatusCode: respStatus}
}

func LookupVmSpecInfoList(connectionName *tbcommon.TbConnectionName) (tbmcir.SpiderSpecInfos, model.WebStatus) {
	fmt.Println("LookupVmSpecInfoList ************ : ")
	var originalUrl = "/lookupSpecs"
	urlParam := util.MappingUrlParameter(originalUrl, nil)
	url := util.TUMBLEBUG + urlParam
	// url := util.TUMBLEBUG + "/lookupSpec"

	pbytes, _ := json.Marshal(connectionName)
	// fmt.Println(string(pbytes))
	resp, err := util.CommonHttp(url, pbytes, http.MethodPost)
	// resp, err := util.CommonHttpWithoutParam(url, http.MethodGet)
	if err != nil {
		fmt.Println(err)
		return nil, model.WebStatus{StatusCode: 500, Message: err.Error()}
	}
	// defer body.Close()
	respBody := resp.Body
	respStatus := resp.StatusCode

	// return respBody, respStatus
	// log.Println(respBody)
	vmSpecList := map[string]tbmcir.SpiderSpecInfos{}

	json.NewDecoder(respBody).Decode(&vmSpecList)
	// fmt.Println(vmSpecList["vmspec"])

	return vmSpecList["vmspec"], model.WebStatus{StatusCode: respStatus}

}

// deprecated 호출경로 변경
// func LookupVmSpecInfoData(vmSpecName string) (*tbmcir.RestGetAllSpecResponse, model.WebStatus) {
// 	var originalUrl = "/lookupSpec/{specName}"
// 	var paramMapper = make(map[string]string)
// 	paramMapper["{specName}"] = vmSpecName
// 	urlParam := util.MappingUrlParameter(originalUrl, paramMapper)
// 	url := util.TUMBLEBUG + urlParam
// 	// url := util.TUMBLEBUG + "/lookupSpec/" + vmSpecName

// 	resp, err := util.CommonHttpWithoutParam(url, http.MethodGet)
// 	vmSpecInfo := tbmcir.RestGetAllSpecResponse{}
// 	if err != nil {
// 		fmt.Println(err)
// 		return &vmSpecInfo, model.WebStatus{StatusCode: 500, Message: err.Error()}
// 	}

// 	respBody := resp.Body
// 	respStatus := resp.StatusCode

// 	json.NewDecoder(respBody).Decode(&vmSpecInfo)
// 	fmt.Println(vmSpecInfo)

// 	return &vmSpecInfo, model.WebStatus{StatusCode: respStatus}
// }

func LookupVmSpecInfoData(restLookupSpecRequest *tbmcir.RestLookupSpecRequest) (*tbmcir.SpiderSpecInfo, model.WebStatus) {
	var originalUrl = "/lookupSpec"

	urlParam := util.MappingUrlParameter(originalUrl, nil)
	url := util.TUMBLEBUG + urlParam

	pbytes, _ := json.Marshal(restLookupSpecRequest)
	resp, err := util.CommonHttp(url, pbytes, http.MethodPost)

	vmSpecInfo := tbmcir.SpiderSpecInfo{}
	if err != nil {
		fmt.Println(err)
		return &vmSpecInfo, model.WebStatus{StatusCode: 500, Message: err.Error()}
	}

	respBody := resp.Body
	respStatus := resp.StatusCode

	json.NewDecoder(respBody).Decode(&vmSpecInfo)
	fmt.Println(vmSpecInfo)

	return &vmSpecInfo, model.WebStatus{StatusCode: respStatus}
}

// Fetch는 결과만 return
func FetchVmSpecInfoList(nameSpaceID string) (*tbcommon.TbSimpleMsg, model.WebStatus) {
	fmt.Println("FetchVmSpecInfoList ************ : ", nameSpaceID)
	var originalUrl = "/ns/{nsId}/resources/fetchSpecs"
	var paramMapper = make(map[string]string)
	paramMapper["{nsId}"] = nameSpaceID
	urlParam := util.MappingUrlParameter(originalUrl, paramMapper)
	url := util.TUMBLEBUG + urlParam
	// url := util.TUMBLEBUG + "/ns/" + nameSpaceID + "/resources/fetchSpecs"

	resultInfo := tbcommon.TbSimpleMsg{}
	// resp, err := util.CommonHttp(url, nil, http.MethodGet)
	resp, err := util.CommonHttpWithoutParam(url, http.MethodPost)
	if err != nil {
		fmt.Println(err)
		return &resultInfo, model.WebStatus{StatusCode: 500, Message: err.Error()}
	}
	// defer body.Close()
	respBody := resp.Body
	respStatus := resp.StatusCode
	// fetchSpecList := map[string][]tbmcir.RestGetAllSpecResponse{}

	json.NewDecoder(respBody).Decode(&resultInfo)
	log.Println("FetchVmSpecList called ")

	return &resultInfo, model.WebStatus{StatusCode: respStatus}
}

// 오래걸리므로 비동기로 처리
func FetchVmSpecInfoListByAsync(nameSpaceID string, c echo.Context) {
	log.Println("FetchVmSpecInfoListByAsync ************ : ", nameSpaceID)
	var originalUrl = "/ns/{nsId}/resources/fetchSpecs"
	var paramMapper = make(map[string]string)
	paramMapper["{nsId}"] = nameSpaceID
	urlParam := util.MappingUrlParameter(originalUrl, paramMapper)
	url := util.TUMBLEBUG + urlParam
	// url := util.TUMBLEBUG + "/ns/" + nameSpaceID + "/resources/fetchSpecs"

	// resultInfo := tbcommon.TbSimpleMsg{}

	taskKey := nameSpaceID + "||" + "VMSpec" + "||" + "Fetch"

	// resp, err := util.CommonHttp(url, nil, http.MethodGet)
	// resp, err := util.CommonHttpWithoutParam(url, http.MethodPost)
	_, err := util.CommonHttpWithoutParam(url, http.MethodPost)
	if err != nil {
		fmt.Println(err)
		StoreWebsocketMessage(util.TASK_TYPE_VMSPEC, taskKey, util.VMSPEC_LIFECYCLE_CREATE, util.TASK_STATUS_FAIL, c) // session에 작업내용 저장
	}
	// defer body.Close()
	// respBody := resp.Body
	// respStatus := resp.StatusCode
	// // fetchSpecList := map[string][]tbmcir.RestGetAllSpecResponse{}

	// json.NewDecoder(respBody).Decode(&resultInfo)

	log.Println("FetchVmSpecList called ")
	StoreWebsocketMessage(util.TASK_TYPE_VMSPEC, taskKey, util.VMSPEC_LIFECYCLE_CREATE, util.TASK_STATUS_COMPLETE, c) // session에 작업내용 저장
}

// resourcesGroup.PUT("/vmspec/put/:specID", controller.VmSpecPutProc)	// RegProc _ SshKey 같이 앞으로 넘길까
// resourcesGroup.POST("/vmspec/filterspecs", controller.FilterVmSpecList)

// spec들을 filterling
func FilterVmSpecInfoList(nameSpaceID string, vmSpecRegInfo *tbmcir.TbSpecInfo) ([]tbmcir.TbSpecInfo, model.WebStatus) {
	fmt.Println("FilterVmSpecInfoList ************ : ", nameSpaceID)
	var originalUrl = "/ns/{nsId}/resources/filterSpecs"
	var paramMapper = make(map[string]string)
	paramMapper["{nsId}"] = nameSpaceID
	urlParam := util.MappingUrlParameter(originalUrl, paramMapper)
	url := util.TUMBLEBUG + urlParam
	// url := util.TUMBLEBUG + "/ns/" + nameSpaceID + "/resources/filterSpecs"
	// /ns/{nsId}/resources/filterSpecs
	// resp, err := util.CommonHttp(url, nil, http.MethodGet)
	// resp, err := util.CommonHttpWithoutParam(url, http.MethodPost)

	pbytes, _ := json.Marshal(vmSpecRegInfo)
	resp, err := util.CommonHttp(url, pbytes, http.MethodGet)

	if err != nil {
		fmt.Println(err)
		return nil, model.WebStatus{StatusCode: 500, Message: err.Error()}
	}
	// defer body.Close()
	respBody := resp.Body
	respStatus := resp.StatusCode
	filterSpecList := tbmcir.RestFilterSpecsResponse{}

	json.NewDecoder(respBody).Decode(&filterSpecList)
	log.Println("FilterVmSpecInfoList called ")

	return filterSpecList.SpaceInfo, model.WebStatus{StatusCode: respStatus}
}

// resourcesGroup.POST("/vmspec/filterspecsbyrange", controller.FilterVmSpecListByRange)
func FilterVmSpecInfoListByRange(nameSpaceID string, vmSpecRangeMinMax *tbmcir.FilterSpecsByRangeRequest) ([]tbmcir.TbSpecInfo, model.WebStatus) {
	webStatus := model.WebStatus{}
	resultInfo := tbmcir.RestFilterSpecsResponse{}

	fmt.Println("FilterVmSpecInfoListByRange ************ : ", nameSpaceID)
	var originalUrl = "/ns/{nsId}/resources/filterSpecsByRange"
	var paramMapper = make(map[string]string)
	paramMapper["{nsId}"] = nameSpaceID
	urlParam := util.MappingUrlParameter(originalUrl, paramMapper)
	url := util.TUMBLEBUG + urlParam
	// url := util.TUMBLEBUG + "/ns/" + nameSpaceID + "/resources/filterSpecsByRange"
	// /ns/{nsId}/resources/filterSpecsByRange

	pbytes, _ := json.Marshal(vmSpecRangeMinMax)
	resp, err := util.CommonHttp(url, pbytes, http.MethodPost)

	if err != nil {
		fmt.Println(err)
		return resultInfo.SpaceInfo, model.WebStatus{StatusCode: 500, Message: err.Error()}
	}
	// defer body.Close()
	respBody := resp.Body
	respStatus := resp.StatusCode

	if respStatus != 200 && respStatus != 201 {
		failResultInfo := tbcommon.TbSimpleMsg{}
		json.NewDecoder(respBody).Decode(&failResultInfo)
		log.Println(resultInfo)
		log.Println("ResultMessage : " + failResultInfo.Message)
		return resultInfo.SpaceInfo, model.WebStatus{StatusCode: respStatus, Message: failResultInfo.Message}
	}

	json.NewDecoder(respBody).Decode(&resultInfo)
	log.Println(resultInfo)
	webStatus.StatusCode = respStatus

	return resultInfo.SpaceInfo, model.WebStatus{StatusCode: respStatus}
}

func DelDefaultResources(nameSpaceID string) (*tbcommon.TbIdList, model.WebStatus) {
	var originalUrl = "/ns/{nsId}/defaultResources"
	var paramMapper = make(map[string]string)
	paramMapper["{nsId}"] = nameSpaceID
	urlParam := util.MappingUrlParameter(originalUrl, paramMapper)
	url := util.TUMBLEBUG + urlParam

	resp, err := util.CommonHttp(url, nil, http.MethodDelete)
	idList := tbcommon.TbIdList{}

	if err != nil {
		fmt.Println(err)
		return &idList, model.WebStatus{StatusCode: 500, Message: err.Error()}
	}
	// return body, err
	respBody := resp.Body
	respStatus := resp.StatusCode

	if respStatus != 200 && respStatus != 201 {
		failResultInfo := tbcommon.TbSimpleMsg{}
		json.NewDecoder(respBody).Decode(&failResultInfo)
		log.Println("ResultMessage : " + failResultInfo.Message)
		return &idList, model.WebStatus{StatusCode: respStatus, Message: failResultInfo.Message}
	}

	json.NewDecoder(respBody).Decode(&idList)
	fmt.Println(idList)

	return &idList, model.WebStatus{StatusCode: respStatus}
}

func LoadDefaultResources(nameSpaceID string, optionParam string, connectionName string) (model.WebStatus, model.WebStatus) {
	var originalUrl = "/ns/{nsId}/loadDefaultResource"
	var paramMapper = make(map[string]string)
	paramMapper["{nsId}"] = nameSpaceID
	urlParam := util.MappingUrlParameter(originalUrl, paramMapper)
	urlParam += "?option=" + optionParam
	if connectionName != "" {
		urlParam += "&connectionName=" + connectionName
	}
	url := util.TUMBLEBUG + urlParam

	resp, err := util.CommonHttp(url, nil, http.MethodGet)
	webStatus := model.WebStatus{}
	if err != nil {
		fmt.Println(err)
		return webStatus, model.WebStatus{StatusCode: 500, Message: err.Error()}
	}

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
}
