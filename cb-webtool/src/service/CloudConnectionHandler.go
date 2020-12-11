package service

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"os"
	"strconv"
	"sync"
	//"io/ioutil"
	//"github.com/davecgh/go-spew/spew"
)

//var CloudConnectionUrl = "http://15.165.16.67:1024"
var CloudConnectionUrl = os.Getenv("SPIDER_URL")
var TumbleUrl = os.Getenv("TUMBLE_URL")

type CloudConnectionInfo struct {
	ID             string `json:"id"`
	ConfigName     string `json:"ConfigName"`
	ProviderName   string `json:"ProviderName"`
	DriverName     string `json:"DriverName"`
	CredentialName string `json:"CredentialName"`
	RegionName     string `json:"RegionName"`
	Description    string `json:"description"`
}
type KeyValueInfo struct {
	Key   string `json:"Key"`
	Value string `json:"Value"`
}
type RegionInfo struct {
	RegionName       string `json:"RegionName"`
	ProviderName     string `json:"ProviderName"`
	KeyValueInfoList []KeyValueInfo
}
type RESP struct {
	Region []struct {
		RegionName       string         `json:"RegionName"`
		ProviderName     string         `json:"ProviderName"`
		KeyValueInfoList []KeyValueInfo `json:"KeyValueInfoList"`
	} `json:"region"`
}

type ImageRESP struct {
	Image []struct {
		id             string         `json:"id"`
		name           string         `json:"name"`
		connectionName string         `json:"connectionName"`
		cspImageId     string         `json:"cspImageId"`
		cspImageName   string         `json:"cspImageName"`
		description    string         `json:"description"`
		guestOS        string         `json:"guestOS"`
		status         string         `json:"status"`
		KeyValueList   []KeyValueInfo `json:"KeyValueList"`
	} `json:"image"`
}
type Image struct {
	id             string         `json:"id"`
	name           string         `json:"name"`
	connectionName string         `json:"connectionName"`
	cspImageId     string         `json:"cspImageId"`
	cspImageName   string         `json:"cspImageName"`
	description    string         `json:"description"`
	guestOS        string         `json:"guestOS"`
	status         string         `json:"status"`
	KeyValueList   []KeyValueInfo `json:"KeyValueList"`
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

type KeepZero float64

func (f KeepZero) MarshalJSON() ([]byte, error) {
	if float64(f) == float64(int(f)) {
		return []byte(strconv.FormatFloat(float64(f), 'f', 1, 32)), nil
	}
	return []byte(strconv.FormatFloat(float64(f), 'f', -1, 32)), nil
}

type myFloat64 float64

func (mf myFloat64) MarshalJSON() ([]byte, error) {
	const ε = 1e-12
	v := float64(mf)
	w, f := math.Modf(v)
	if f < ε {
		return []byte(fmt.Sprintf(`%v.0`, math.Trunc(w))), nil
	}
	return json.Marshal(v)
}

func GetConnectionconfig(drivername string) CloudConnectionInfo {
	url := NameSpaceUrl + "/driver/" + drivername

	// resp, err := http.Get(url)

	// if err != nil {
	// 	fmt.Println("request URL : ", url)
	// }

	// defer resp.Body.Close()
	body := HttpGetHandler(url)
	nsInfo := CloudConnectionInfo{}

	json.NewDecoder(body).Decode(&nsInfo)
	fmt.Println("nsInfo : ", nsInfo.ID)
	return nsInfo

}
func GetImageList() []Image {
	url := CloudConnectionUrl + "/connectionconfig"
	// resp, err := http.Get(url)
	// if err != nil {
	// 	fmt.Println("request URL : ", url)
	// }

	// defer resp.Body.Close()
	body := HttpGetHandler(url)
	defer body.Close()

	nsInfo := ImageRESP{}
	json.NewDecoder(body).Decode(&nsInfo)
	fmt.Println("nsInfo : ", nsInfo.Image[0].id)
	var info []Image
	for _, item := range nsInfo.Image {
		reg := Image{
			id:             item.id,
			name:           item.name,
			connectionName: item.connectionName,
			cspImageId:     item.cspImageId,
			cspImageName:   item.cspImageName,
			description:    item.description,
			guestOS:        item.guestOS,
			status:         item.status,
		}
		info = append(info, reg)
	}
	return info

}

func GetConnectionList() []CloudConnectionInfo {
	url := CloudConnectionUrl + "/connectionconfig"
	// resp, err := http.Get(url)
	// if err != nil {
	// 	fmt.Println("request URL : ", url)
	// }

	// defer resp.Body.Close()
	body := HttpGetHandler(url)
	defer body.Close()

	nsInfo := map[string][]CloudConnectionInfo{}
	json.NewDecoder(body).Decode(&nsInfo)
	//fmt.Println("nsInfo : ", nsInfo["connectionconfig"][0].ID)
	return nsInfo["connectionconfig"]

}

func GetDriverReg() []CloudConnectionInfo {
	url := NameSpaceUrl + "/driver"

	body := HttpGetHandler(url)
	defer body.Close()
	nsInfo := map[string][]CloudConnectionInfo{}
	json.NewDecoder(body).Decode(&nsInfo)

	return nsInfo["driver"]

}

func GetCredentialList() []CloudConnectionInfo {
	url := CloudConnectionUrl + "/credential"
	// resp, err := http.Get(url)
	// if err != nil {
	// 	fmt.Println("request URL : ", url)
	// }

	// defer resp.Body.Close()
	body := HttpGetHandler(url)
	defer body.Close()
	nsInfo := map[string][]CloudConnectionInfo{}
	json.NewDecoder(body).Decode(&nsInfo)
	// fmt.Println("nsInfo : ", nsInfo["credential"][0].ID)
	return nsInfo["credential"]

}
func GetRegionList() []RegionInfo {
	url := CloudConnectionUrl + "/region"
	fmt.Println("=========== Get Start Region List : ", url)

	body := HttpGetHandler(url)
	defer body.Close()

	nsInfo := RESP{}

	json.NewDecoder(body).Decode(&nsInfo)
	var info []RegionInfo

	for _, item := range nsInfo.Region {

		reg := RegionInfo{
			RegionName:   item.RegionName,
			ProviderName: item.ProviderName,
		}
		info = append(info, reg)
	}
	fmt.Println("info region list : ", info)
	return info

}

func GetCredentialReg() []CloudConnectionInfo {
	url := CloudConnectionUrl + "/credential"

	body := HttpGetHandler(url)
	defer body.Close()

	nsInfo := map[string][]CloudConnectionInfo{}
	json.NewDecoder(body).Decode(&nsInfo)
	fmt.Println("nsInfo : ", nsInfo["credential"][0].ID)
	return nsInfo["credential"]

}

func GetGeoMetryInfo(wg *sync.WaitGroup, ip_address string, status string, vm_id string, vm_name string, returnResult *[]IPStackInfo) {
	defer wg.Done() //goroutin sync done

	apiUrl := "http://api.ipstack.com/"
	access_key := "86c895286435070c0369a53d2d0b03d1"
	url := apiUrl + ip_address + "?access_key=" + access_key
	resp, err := http.Get(url)
	fmt.Println("GetGeoMetryInfo request URL : ", url)
	if err != nil {
		fmt.Println("GetGeoMetryInfo request URL : ", url)
	}
	defer resp.Body.Close()

	//그냥 스트링으로 반환해서 프론트에서 JSON.parse로 처리 하는 방법도 괜찮네
	//spew.Dump(resp.Body)
	// bytes, _ := ioutil.ReadAll(resp.Body)
	// str := string(bytes)
	// fmt.Println(str)
	// *returnStr = append(*returnStr, str)

	ipStackInfo := IPStackInfo{
		VMID:   vm_id,
		Status: status,
		VMName: vm_name,
	}

	json.NewDecoder(resp.Body).Decode(&ipStackInfo)
	fmt.Println("Get GeoMetry INFO :", ipStackInfo)

	*returnResult = append(*returnResult, ipStackInfo)

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
