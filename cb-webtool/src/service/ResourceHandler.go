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
	"github.com/cloud-barista/cb-webtool/src/model/tumblebug"

	util "github.com/cloud-barista/cb-webtool/src/util"
)

// 해당 namespace의 vpc 목록 조회
//func GetVnetList(nameSpaceID string) (io.ReadCloser, error) {
func GetVnetList(nameSpaceID string) ([]tumblebug.VNetInfo, model.WebStatus) {
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
	vNetInfoList := map[string][]tumblebug.VNetInfo{}
	json.NewDecoder(respBody).Decode(&vNetInfoList)
	//spew.Dump(body)
	fmt.Println(vNetInfoList["vNet"])

	return vNetInfoList["vNet"], model.WebStatus{StatusCode: respStatus}

}

// vpc 상세 조회-> ResourceHandler로 이동
func GetVpcData(nameSpaceID string, vNetID string) (*tumblebug.VNetInfo, model.WebStatus) {
	var originalUrl = "/ns/{nsId}/resources/vNet/{vNetId}"
	var paramMapper = make(map[string]string)
	paramMapper["{nsId}"] = nameSpaceID
	paramMapper["{vNetId}"] = vNetID
	urlParam := util.MappingUrlParameter(originalUrl, paramMapper)
	url := util.TUMBLEBUG + urlParam
	// url := util.TUMBLEBUG + "/ns/" + nameSpaceID + "/resources/vNet/" + vNetID

	resp, err := util.CommonHttpWithoutParam(url, http.MethodGet)
	vNetInfo := tumblebug.VNetInfo{}
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
// func RegVpc(nameSpaceID string, vnetRegInfo *tumblebug.VNetRegInfo) (io.ReadCloser, int) {
func RegVpc(nameSpaceID string, vnetRegInfo *tumblebug.VNetRegInfo) (*tumblebug.VNetInfo, model.WebStatus) {
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
	vNetInfo := tumblebug.VNetInfo{}
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

// 해당 namespace의 SecurityGroup 목록 조회
func GetSecurityGroupList(nameSpaceID string) ([]tumblebug.SecurityGroupInfo, model.WebStatus) {
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
	securityGroupList := map[string][]tumblebug.SecurityGroupInfo{}

	json.NewDecoder(respBody).Decode(&securityGroupList)
	//spew.Dump(body)
	fmt.Println(securityGroupList["securityGroup"])

	return securityGroupList["securityGroup"], model.WebStatus{StatusCode: respStatus}

}

// SecurityGroup 상세 조회
func GetSecurityGroupData(nameSpaceID string, securityGroupID string) (*tumblebug.SecurityGroupInfo, model.WebStatus) {
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
	securityGroupInfo := tumblebug.SecurityGroupInfo{}
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
func RegSecurityGroup(nameSpaceID string, securityGroupRegInfo *tumblebug.SecurityGroupRegInfo) (*tumblebug.SecurityGroupInfo, model.WebStatus) {
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
	securityGroupInfo := tumblebug.SecurityGroupInfo{}
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
func GetSshKeyInfoList(nameSpaceID string) ([]tumblebug.SshKeyInfo, model.WebStatus) {
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
	sshKeyList := map[string][]tumblebug.SshKeyInfo{}

	json.NewDecoder(respBody).Decode(&sshKeyList)
	//spew.Dump(body)
	fmt.Println(sshKeyList["sshKey"])

	return sshKeyList["sshKey"], model.WebStatus{StatusCode: respStatus}

}

// sshKey 상세 조회
func GetSshKeyData(nameSpaceID string, sshKeyID string) (*tumblebug.SshKeyInfo, model.WebStatus) {
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
	sshKeyInfo := tumblebug.SshKeyInfo{}
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
func RegSshKey(nameSpaceID string, sshKeyRegInfo *tumblebug.SshKeyRegInfo) (*tumblebug.SshKeyInfo, model.WebStatus) {
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
	sshKeyInfo := tumblebug.SshKeyInfo{}
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

// sshKey 삭제
func DelSshKey(nameSpaceID string, sshKeyID string) (model.WebStatus, model.WebStatus) {
	webStatus := model.WebStatus{}
	// if ValidateString(sshKeyID) != nil {
	if len(sshKeyID) == 0 {
		log.Println("securityGroupID 가 없으면 해당 namespace의 모든 securityGroup이 삭제되므로 처리할 수 없습니다.")
		return webStatus, model.WebStatus{StatusCode: 4040, Message: "securityGroupID 가 없으면 해당 namespace의 모든 securityGroup이 삭제되므로 처리할 수 없습니다."}
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

// VirtualMachineImage 목록 조회
func GetVirtualMachineImageInfoList(nameSpaceID string) ([]tumblebug.VirtualMachineImageInfo, model.WebStatus) {
	fmt.Println("GetVirtualMachineImageInfoList ************ : ")
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

	virtualMachineImageList := map[string][]tumblebug.VirtualMachineImageInfo{}

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
func GetVirtualMachineImageData(nameSpaceID string, virtualMachineImageID string) (*tumblebug.VirtualMachineImageInfo, model.WebStatus) {
	fmt.Println("GetVirtualMachineImageData ************ : ")
	var originalUrl = "/ns/{nsId}/resources/image/{imageId}"
	var paramMapper = make(map[string]string)
	paramMapper["{nsId}"] = nameSpaceID
	paramMapper["{imageId}"] = virtualMachineImageID
	urlParam := util.MappingUrlParameter(originalUrl, paramMapper)
	url := util.TUMBLEBUG + urlParam
	// url := util.TUMBLEBUG + "/ns/" + nameSpaceID + "/resources/image/" + virtualMachineImageID

	resp, err := util.CommonHttpWithoutParam(url, http.MethodGet)
	virtualMachineImageInfo := tumblebug.VirtualMachineImageInfo{}
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

// VirtualMachineImage 등록
func RegVirtualMachineImage(nameSpaceID string, registType string, virtualMachineImageRegInfo *tumblebug.VirtualMachineImageRegInfo) (*tumblebug.VirtualMachineImageInfo, model.WebStatus) {
	fmt.Println("RegVirtualMachineImage ************ : ")
	if registType == "" {
		registType = "registerWithId" // registerWithId 또는 registerWithInfo
	}

	var originalUrl = "/ns/{nsId}/resources/image?action={registType}"
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
	virtualMachineImageInfo := tumblebug.VirtualMachineImageInfo{}
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
func GetInspectResourceList(inspectResource *tumblebug.InspectResourcesRequest) (*tumblebug.InspectResourcesResponse, model.WebStatus) {
	fmt.Println("GetInspectResourceList ************ : ")
	//https://www.javaer101.com/ko/article/5704925.html 참조 : 값이 있는 것만 넘기기
	var originalUrl = "/inspectResources"
	urlParam := util.MappingUrlParameter(originalUrl, nil)
	url := util.TUMBLEBUG + urlParam

	pbytes, _ := json.Marshal(inspectResource)
	resp, err := util.CommonHttp(url, pbytes, http.MethodPost)

	inspectResourcesResponse := tumblebug.InspectResourcesResponse{}
	if err != nil {
		fmt.Println(err)
		return &inspectResourcesResponse, model.WebStatus{StatusCode: 500, Message: err.Error()}
	}

	respBody := resp.Body
	respStatus := resp.StatusCode

	json.NewDecoder(respBody).Decode(inspectResourcesResponse)
	fmt.Println(inspectResourcesResponse)

	return &inspectResourcesResponse, model.WebStatus{StatusCode: respStatus}

}

// VM Image 조회
func LookupVirtualMachineImageList(connectionName string) ([]tumblebug.VirtualMachineLookupImageInfo, model.WebStatus) {
	fmt.Println("LookupVirtualMachineImageList ************ : ", connectionName)
	var originalUrl = "/lookupImages"
	urlParam := util.MappingUrlParameter(originalUrl, nil)
	url := util.TUMBLEBUG + urlParam
	// url := util.TUMBLEBUG + "/lookupImage"

	// body, err := util.CommonHttpGet(url)
	// resp, err := util.CommonHttpWithoutParam(url, http.MethodGet)

	paramMap := map[string]string{"connectionName": connectionName}

	pbytes, _ := json.Marshal(paramMap)
	log.Println(string(pbytes))
	resp, err := util.CommonHttp(url, pbytes, http.MethodGet)
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
	lookupImageList := map[string][]tumblebug.VirtualMachineLookupImageInfo{}

	json.NewDecoder(respBody).Decode(&lookupImageList)
	log.Println("LookupVirtualMachineImageList called 3 ")

	return lookupImageList["image"], model.WebStatus{StatusCode: respStatus}
}

func LookupVirtualMachineImageData(virtualMachineImageID string) (*tumblebug.VirtualMachineImageInfo, model.WebStatus) {
	var originalUrl = "/lookupImage/{imageId}"
	var paramMapper = make(map[string]string)
	paramMapper["{imageId}"] = virtualMachineImageID
	urlParam := util.MappingUrlParameter(originalUrl, paramMapper)
	url := util.TUMBLEBUG + urlParam
	// url := util.TUMBLEBUG + "/lookupImage/" + virtualMachineImageID

	// pbytes, _ := json.Marshal(nameSpaceID)
	resp, err := util.CommonHttpWithoutParam(url, http.MethodGet)
	virtualMachineImageInfo := tumblebug.VirtualMachineImageInfo{}
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
func FetchVirtualMachineImageList(nameSpaceID string) ([]tumblebug.VirtualMachineLookupImageInfo, model.WebStatus) {
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
	fetchImageList := map[string][]tumblebug.VirtualMachineLookupImageInfo{}

	json.NewDecoder(respBody).Decode(&fetchImageList)
	log.Println("FetchVirtualMachineImageList called ")

	log.Println("FetchVirtualMachineImageList called 5 ")
	return fetchImageList["image"], model.WebStatus{StatusCode: respStatus}
}

// VirtualMachineImage 상세 조회
func SearchVirtualMachineImageList(nameSpaceID string, virtualMachineImageID string) (*tumblebug.VirtualMachineImageInfo, model.WebStatus) {
	var originalUrl = "/ns/{nsId}/resources/searchImage"
	var paramMapper = make(map[string]string)
	paramMapper["{nsId}"] = nameSpaceID
	urlParam := util.MappingUrlParameter(originalUrl, paramMapper)
	url := util.TUMBLEBUG + urlParam
	// url := util.TUMBLEBUG + "/ns/" + nameSpaceID + "/resources/searchImage/"

	// pbytes, _ := json.Marshal(nameSpaceID)
	resp, err := util.CommonHttpWithoutParam(url, http.MethodGet)
	virtualMachineImageInfo := tumblebug.VirtualMachineImageInfo{}
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

// VMSpec 목록 조회
func GetVmSpecInfoList(nameSpaceID string) ([]tumblebug.VmSpecInfo, model.WebStatus) {
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
	vmSpecList := map[string][]tumblebug.VmSpecInfo{}

	json.NewDecoder(respBody).Decode(&vmSpecList)
	//spew.Dump(body)
	fmt.Println(vmSpecList["spec"])

	return vmSpecList["spec"], model.WebStatus{StatusCode: respStatus}

}

// VMSpec 상세 조회
func GetVmSpecInfoData(nameSpaceID string, vmSpecID string) (*tumblebug.VmSpecInfo, model.WebStatus) {
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
	vmSpecInfo := tumblebug.VmSpecInfo{}
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
func RegVmSpec(nameSpaceID string, specregisteringMethod string, vmSpecRegInfo *tumblebug.VmSpecRegInfo) (*tumblebug.VmSpecInfo, model.WebStatus) {
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

	vmSpecInfo := tumblebug.VmSpecInfo{}
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

func UpdateVMSpec(nameSpaceID string, vmSpecRegInfo *tumblebug.VmSpecRegInfo) (*tumblebug.VmSpecInfo, model.WebStatus) {
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
	vmSpecInfo := tumblebug.VmSpecInfo{}
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

func LookupVmSpecInfoList(connectionName *tumblebug.TbConnectionName) ([]tumblebug.SpiderSpecInfo, model.WebStatus) {
	fmt.Println("LookupVmSpecInfoList ************ : ")
	var originalUrl = "/lookupSpecs"
	urlParam := util.MappingUrlParameter(originalUrl, nil)
	url := util.TUMBLEBUG + urlParam
	// url := util.TUMBLEBUG + "/lookupSpec"

	pbytes, _ := json.Marshal(connectionName)
	// fmt.Println(string(pbytes))
	resp, err := util.CommonHttp(url, pbytes, http.MethodGet)
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
	vmSpecList := map[string][]tumblebug.SpiderSpecInfo{}

	json.NewDecoder(respBody).Decode(&vmSpecList)
	// fmt.Println(vmSpecList["vmspec"])

	return vmSpecList["vmspec"], model.WebStatus{StatusCode: respStatus}

}

func LookupVmSpecInfoData(vmSpecName string) (*tumblebug.VmSpecInfo, model.WebStatus) {
	var originalUrl = "/lookupSpec/{specName}"
	var paramMapper = make(map[string]string)
	paramMapper["{specName}"] = vmSpecName
	urlParam := util.MappingUrlParameter(originalUrl, paramMapper)
	url := util.TUMBLEBUG + urlParam
	// url := util.TUMBLEBUG + "/lookupSpec/" + vmSpecName

	resp, err := util.CommonHttpWithoutParam(url, http.MethodGet)
	vmSpecInfo := tumblebug.VmSpecInfo{}
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

func FetchVmSpecInfoList(nameSpaceID string) ([]tumblebug.VmSpecInfo, model.WebStatus) {
	fmt.Println("FetchVmSpecInfoList ************ : ", nameSpaceID)
	var originalUrl = "/ns/{nsId}/resources/fetchSpecs"
	var paramMapper = make(map[string]string)
	paramMapper["{nsId}"] = nameSpaceID
	urlParam := util.MappingUrlParameter(originalUrl, paramMapper)
	url := util.TUMBLEBUG + urlParam
	// url := util.TUMBLEBUG + "/ns/" + nameSpaceID + "/resources/fetchSpecs"

	// resp, err := util.CommonHttp(url, nil, http.MethodGet)
	resp, err := util.CommonHttpWithoutParam(url, http.MethodPost)
	if err != nil {
		fmt.Println(err)
		return nil, model.WebStatus{StatusCode: 500, Message: err.Error()}
	}
	// defer body.Close()
	respBody := resp.Body
	respStatus := resp.StatusCode
	fetchSpecList := map[string][]tumblebug.VmSpecInfo{}

	json.NewDecoder(respBody).Decode(&fetchSpecList)
	log.Println("FetchVmSpecList called ")

	return fetchSpecList["spec"], model.WebStatus{StatusCode: respStatus}
}

// resourcesGroup.PUT("/vmspec/put/:specID", controller.VmSpecPutProc)	// RegProc _ SshKey 같이 앞으로 넘길까
// resourcesGroup.POST("/vmspec/filterspecs", controller.FilterVmSpecList)

// spec들을 filterling
func FilterVmSpecInfoList(nameSpaceID string, vmSpecRegInfo *tumblebug.VmSpecRegInfo) ([]tumblebug.VmSpecInfo, model.WebStatus) {
	fmt.Println("FilterVmSpecInfoList ************ : ", nameSpaceID)
	var originalUrl = "/ns/{nsId}/resources/filterSpecs"
	var paramMapper = make(map[string]string)
	paramMapper["{nsId}"] = nameSpaceID
	urlParam := util.MappingUrlParameter(originalUrl, paramMapper)
	url := util.TUMBLEBUG + urlParam
	// url := util.TUMBLEBUG + "/ns/" + nameSpaceID + "/resources/filterSpecs"
	// /ns/{nsId}/resources/filterSpecs
	// resp, err := util.CommonHttp(url, nil, http.MethodGet)
	resp, err := util.CommonHttpWithoutParam(url, http.MethodPost)
	if err != nil {
		fmt.Println(err)
		return nil, model.WebStatus{StatusCode: 500, Message: err.Error()}
	}
	// defer body.Close()
	respBody := resp.Body
	respStatus := resp.StatusCode
	fetchSpecList := map[string][]tumblebug.VmSpecInfo{}

	json.NewDecoder(respBody).Decode(&fetchSpecList)
	log.Println("FilterVmSpecInfoList called ")

	return fetchSpecList["spec"], model.WebStatus{StatusCode: respStatus}
}

// resourcesGroup.POST("/vmspec/filterspecsbyrange", controller.FilterVmSpecListByRange)
func FilterVmSpecInfoListByRange(nameSpaceID string, vmSpecRangeMinMax *tumblebug.VmSpecRangeReqInfo) ([]tumblebug.TbSpecInfo, model.WebStatus) {
	webStatus := model.WebStatus{}
	// vmSpecInfo := tumblebug.VmSpecRangeInfo{}
	vmSpecInfo := map[string][]tumblebug.TbSpecInfo{}
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
		return vmSpecInfo["spec"], model.WebStatus{StatusCode: 500, Message: err.Error()}
	}
	// defer body.Close()
	respBody := resp.Body
	respStatus := resp.StatusCode

	if respStatus != 200 && respStatus != 201 {
		resultInfo := model.ResultInfo{}

		json.NewDecoder(respBody).Decode(&resultInfo)
		log.Println(resultInfo)
		log.Println("ResultMessage : " + resultInfo.Message)
		return vmSpecInfo["spec"], model.WebStatus{StatusCode: respStatus, Message: resultInfo.Message}
	}

	json.NewDecoder(respBody).Decode(&vmSpecInfo)
	log.Println(vmSpecInfo)
	webStatus.StatusCode = respStatus

	return vmSpecInfo["spec"], model.WebStatus{StatusCode: respStatus}
}
