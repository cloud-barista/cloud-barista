package service

import (
	"encoding/json"
	"fmt"
	"net/http"
)

var CloudConnectionUrl = "http://15.165.16.67:1024"

type CloudConnectionInfo struct {
	ID             string `json:"id"`
	ConfigName     string `json:"ConfigName"`
	ProviderName   string `json:"ProviderName"`
	DriverName     string `json:"DriverName"`
	CredentialName string `json:"CredentialName"`
	RegionName     string `json:"RegionName"`
	Description    string `json:"description"`
}

func GetConnectionconfig(drivername string) CloudConnectionInfo {
	url := NameSpaceUrl + "/driver/" + drivername

	resp, err := http.Get(url)

	if err != nil {
		fmt.Println("request URL : ", url)
	}

	defer resp.Body.Close()
	nsInfo := CloudConnectionInfo{}
	json.NewDecoder(resp.Body).Decode(&nsInfo)
	fmt.Println("nsInfo : ", nsInfo.ID)
	return nsInfo

}

func GetConnectionList() []CloudConnectionInfo {
	url := CloudConnectionUrl + "/connectionconfig"
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("request URL : ", url)
	}

	defer resp.Body.Close()
	nsInfo := map[string][]CloudConnectionInfo{}
	json.NewDecoder(resp.Body).Decode(&nsInfo)
	fmt.Println("nsInfo : ", nsInfo["ns"][0].ID)
	return nsInfo["ns"]

}

func GetDriverReg() []CloudConnectionInfo {
	url := NameSpaceUrl + "/driver"
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("request URL : ", url)
	}

	defer resp.Body.Close()
	nsInfo := map[string][]CloudConnectionInfo{}
	json.NewDecoder(resp.Body).Decode(&nsInfo)
	fmt.Println("nsInfo : ", nsInfo["ns"][0].ID)
	return nsInfo["ns"]

}

func GetCredentialList() []CloudConnectionInfo {
	url := NameSpaceUrl + "/driver"
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("request URL : ", url)
	}

	defer resp.Body.Close()
	nsInfo := map[string][]CloudConnectionInfo{}
	json.NewDecoder(resp.Body).Decode(&nsInfo)
	fmt.Println("nsInfo : ", nsInfo["ns"][0].ID)
	return nsInfo["ns"]

}

func GetCredentialReg() []CloudConnectionInfo {
	url := NameSpaceUrl + "/driver"
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("request URL : ", url)
	}

	defer resp.Body.Close()
	nsInfo := map[string][]CloudConnectionInfo{}
	json.NewDecoder(resp.Body).Decode(&nsInfo)
	fmt.Println("nsInfo : ", nsInfo["ns"][0].ID)
	return nsInfo["ns"]

}

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
