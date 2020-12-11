package service

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"reflect"
	//"github.com/davecgh/go-spew/spew"
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

	body := HttpGetHandler(url)
	defer body.Close()
	nsInfo := NSInfo{}
	json.NewDecoder(body).Decode(&nsInfo)
	fmt.Println("nsInfo : ", nsInfo.ID)
	return nsInfo

}

func GetNSList() []NSInfo {
	url := NameSpaceUrl + "/ns"
	fmt.Println("============= NameSpace URL =============", url)
	// authInfo := controller.AuthenticationHandler()
	// req, err := http.NewRequest("GET", url, nil)
	// if err != nil {

	// }
	// req.Header.Add("Authorization", authInfo)
	// client := &http.Client{}
	// resp, err := client.Do(req)
	// fmt.Println("=============result GetNSList =============", resp)
	// //spew.Dump(resp)
	// if err != nil {
	// 	fmt.Println("========= GetNSList Error : ", err)
	// 	fmt.Println("request URL : ", url)
	// 	return nil
	// }

	// defer resp.Body.Close()
	body := HttpGetHandler(url)
	nsInfo := map[string][]NSInfo{}
	defer body.Close()
	json.NewDecoder(body).Decode(&nsInfo)
	//spew.Dump(body)
	return nsInfo["ns"]

}

func GetNSCnt() int {
	url := NameSpaceUrl + "/ns"
	fmt.Println("============= NameSpace URL =============", url)

	// defer resp.Body.Close()
	body := HttpGetHandler(url)
	nsInfo := map[string][]NSInfo{}
	defer body.Close()
	json.NewDecoder(body).Decode(&nsInfo)
	//spew.Dump(body)
	if nsInfo["ns"] == nil {
		return 0
	} else {
		return len(nsInfo["ns"])

	}

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

func HttpGetHandler(url string) io.ReadCloser {
	authInfo := AuthenticationHandler()

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", authInfo)

	client := &http.Client{}
	resp, _ := client.Do(req)

	//defer resp.Body.Close()

	return resp.Body
}
func AuthenticationHandler() string {

	api_username := os.Getenv("API_USERNAME")
	api_password := os.Getenv("API_PASSWORD")

	//The header "KEY: VAL" is "Authorization: Basic {base64 encoded $USERNAME:$PASSWORD}".
	apiUserInfo := api_username + ":" + api_password
	encA := base64.StdEncoding.EncodeToString([]byte(apiUserInfo))
	//req.Header.Add("Authorization", "Basic"+encA)
	return "Basic " + encA

}
