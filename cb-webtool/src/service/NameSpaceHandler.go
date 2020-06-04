package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"reflect"
)

// var NameSpaceUrl = "http://15.165.16.67:1323"
var NameSpaceUrl = os.Getenv("TUMBLE_URL")

type NSInfo struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func GetNS(nsID string) NSInfo {
	url := NameSpaceUrl + "ns" + nsID

	resp, err := http.Get(url)

	if err != nil {
		fmt.Println("request URL : ", url)
	}

	defer resp.Body.Close()
	nsInfo := NSInfo{}
	json.NewDecoder(resp.Body).Decode(&nsInfo)
	fmt.Println("nsInfo : ", nsInfo.ID)
	return nsInfo

}

func GetNSList() []NSInfo {
	url := NameSpaceUrl + "/ns"
	fmt.Println("============= NameSpace URL =============", url)
	resp, err := http.Get(url)
	fmt.Println("=============result GetNSList =============", resp)
	//spew.Dump(resp)
	if err != nil {
		fmt.Println("========= GetNSList Error : ", err)
		fmt.Println("request URL : ", url)
		return nil
	}

	defer resp.Body.Close()
	nsInfo := map[string][]NSInfo{}

	json.NewDecoder(resp.Body).Decode(&nsInfo)
	//fmt.Println("nsInfo : ", nsInfo["ns"][0].ID)
	return nsInfo["ns"]

}

// func RegNS() error {

// }

func RequestGet(url string) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("request URL : ", url)
	}

	defer resp.Body.Close()
	nsInfo := map[string][]NSInfo{}
	fmt.Println("nsInfo type : ", reflect.TypeOf(nsInfo))
	json.NewDecoder(resp.Body).Decode(&nsInfo)
	fmt.Println("nsInfo : ", nsInfo["ns"][0].ID)

	// data, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	fmt.Println("Get Data Error")
	// }
	// fmt.Println("GetData : ", string(data))

}
