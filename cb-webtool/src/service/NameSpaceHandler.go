package service

import (
	"encoding/json"
	// "errors"
	"fmt"
	"io"
	"log"
	"net/http"

	// "os"
	// "bytes"
	// "reflect"
	//"github.com/davecgh/go-spew/spew"
	model "github.com/cloud-barista/cb-webtool/src/model"
	// "github.com/cloud-barista/cb-webtool/src/model/spider"
	"github.com/cloud-barista/cb-webtool/src/model/tumblebug"

	util "github.com/cloud-barista/cb-webtool/src/util"
)

// var NameSpaceUrl = "http://15.165.16.67:1323"
// var NameSpaceUrl = os.Getenv("TUMBLE_URL")

// type NSInfo struct {
// 	ID          string `json:"id"`
// 	Name        string `json:"name"`
// 	Description string `json:"description"`
// }

// 저장된 namespace가 없을 때 최초 1개 생성하고 해당 namespace 정보를 return  : 검증 필요(TODO : 이미 namespace가 있어서 확인 못함)
func CreateDefaultNamespace() (*tumblebug.NameSpaceInfo, model.WebStatus) {
	// nsInfo := new(model.NSInfo)
	nameSpaceInfo := tumblebug.NameSpaceInfo{}

	// 사용자의 namespace 목록조회
	nsList, nsStatus := GetNameSpaceList()
	if nsStatus.StatusCode == 500 {
		log.Println(" nsStatus  ", nsStatus)
		return nil, nsStatus
	}

	if len(nsList) > 0 {
		nsStatus.StatusCode = 101
		nsStatus.Message = "Namespace already exists"
		//return &nameSpaceInfo, errors.New(101, "Namespace already exists. size="+len(nsList))
		return &nameSpaceInfo, nsStatus
	}

	// create default namespace
	nameSpaceInfo.Name = "NS-01" // default namespace name
	//nameSpaceInfo.ID = "NS-01"
	nameSpaceInfo.Description = "default name space name"
	respBody, respStatus := RegNameSpace(&nameSpaceInfo)
	log.Println(" respBody  ", respBody) // respBody에 namespace Id가 있으면 할당해서 사용할 것
	if respStatus.StatusCode != 200 && respStatus.StatusCode != 201 {
		log.Println(" nsCreateErr  ", respStatus)
		return &nameSpaceInfo, respStatus
	}
	// respBody := resp.Body
	// respStatus := resp.StatusCode

	return &nameSpaceInfo, respStatus
}

// 사용자의 namespace 목록 조회
func GetNameSpaceList() ([]tumblebug.NameSpaceInfo, model.WebStatus) {
	fmt.Println("GetNameSpaceList start")
	var originalUrl = "/ns"
	urlParam := util.MappingUrlParameter(originalUrl, nil)

	url := util.TUMBLEBUG + urlParam
	// url := util.TUMBLEBUG + "/ns"

	resp, err := util.CommonHttp(url, nil, http.MethodGet)
	//body := HttpGetHandler(url)

	if err != nil {
		// 	// Tumblebug 접속 확인하라고
		// fmt.Println(err)
		// panic(err)
		return nil, model.WebStatus{StatusCode: 500, Message: err.Error()}
	}

	respBody := resp.Body
	respStatus := resp.StatusCode

	nameSpaceInfoList := map[string][]tumblebug.NameSpaceInfo{}
	// defer body.Close()
	json.NewDecoder(respBody).Decode(&nameSpaceInfoList)
	//spew.Dump(body)
	fmt.Println(nameSpaceInfoList["ns"])

	return nameSpaceInfoList["ns"], model.WebStatus{StatusCode: respStatus}
}

// Get namespace
func GetNameSpaceData(nameSpaceID string) (tumblebug.NameSpaceInfo, model.WebStatus) {
	fmt.Println("GetNameSpaceData start")
	var originalUrl = "/ns/{nsId}"
	var paramMapper = make(map[string]string)
	paramMapper["{nsId}"] = nameSpaceID
	urlParam := util.MappingUrlParameter(originalUrl, paramMapper)

	url := util.TUMBLEBUG + urlParam
	// url := util.TUMBLEBUG + "/ns/" + nameSpaceID

	resp, err := util.CommonHttp(url, nil, http.MethodGet)

	nameSpaceInfo := tumblebug.NameSpaceInfo{}
	if err != nil {
		return nameSpaceInfo, model.WebStatus{StatusCode: 500, Message: err.Error()}
	}

	respBody := resp.Body
	respStatus := resp.StatusCode

	// defer body.Close()
	json.NewDecoder(respBody).Decode(&nameSpaceInfo)
	fmt.Println(nameSpaceInfo)

	return nameSpaceInfo, model.WebStatus{StatusCode: respStatus}
}

// NameSpace 등록
func RegNameSpace(nameSpaceInfo *tumblebug.NameSpaceInfo) (io.ReadCloser, model.WebStatus) {
	// buff := bytes.NewBuffer(pbytes)
	var originalUrl = "/ns"
	urlParam := util.MappingUrlParameter(originalUrl, nil)
	url := util.TUMBLEBUG + urlParam
	// url := util.TUMBLEBUG + "/ns"

	//body, err := util.CommonHttpPost(url, nameSpaceInfo)
	pbytes, _ := json.Marshal(nameSpaceInfo)
	resp, err := util.CommonHttp(url, pbytes, http.MethodPost)
	if err != nil {
		fmt.Println(err)
		return nil, model.WebStatus{StatusCode: 500, Message: err.Error()}
	}
	// return body, err
	respBody := resp.Body
	respStatus := resp.StatusCode
	return respBody, model.WebStatus{StatusCode: respStatus}
}

// NameSpace 수정 : namespace 없데이트 기능 없음
func UpdateNameSpace(nameSpaceID string, nameSpaceInfo *tumblebug.NameSpaceInfo) (io.ReadCloser, model.WebStatus) {
	var originalUrl = "/ns/{nsId}"
	var paramMapper = make(map[string]string)
	paramMapper["{nsId}"] = nameSpaceID
	urlParam := util.MappingUrlParameter(originalUrl, paramMapper)
	url := util.TUMBLEBUG + urlParam
	// url := util.TUMBLEBUG + "/ns"

	pbytes, _ := json.Marshal(nameSpaceInfo)
	resp, err := util.CommonHttp(url, pbytes, http.MethodPut)

	if err != nil {
		fmt.Println(err)
		return nil, model.WebStatus{StatusCode: 500, Message: err.Error()}
	}
	// return body, err
	respBody := resp.Body
	respStatus := resp.StatusCode
	return respBody, model.WebStatus{StatusCode: respStatus}
}

// NameSpace 삭제
func DelNameSpace(nameSpaceID string) (io.ReadCloser, model.WebStatus) {
	var originalUrl = "/ns/{nsId}"
	var paramMapper = make(map[string]string)
	paramMapper["{nsId}"] = nameSpaceID
	urlParam := util.MappingUrlParameter(originalUrl, paramMapper)
	url := util.TUMBLEBUG + urlParam
	// url := util.TUMBLEBUG + "/ns/" + nameSpaceID

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
