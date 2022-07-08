package util

// 공통으로 사용할 function 정의
import (
	// "encoding/base64"
	// "fmt"
	// "io"
	// "io/ioutil"
	"log"
	// "net/http"
	// "net/url"
	"os"
	// "reflect"
	// "strconv"
	"strings"
	// "time"
	// "bytes"
	"encoding/json"
	// "math"
	// "io/ioutil"
	echosession "github.com/go-session/echo-session"
	"github.com/labstack/echo"
	// "github.com/cloud-barista/cb-webtool/src/model"
)

func GetUserInfo(c echo.Context, userID string) (map[string]string, bool) {
	user := os.Getenv("LoginUser")
	store := echosession.FromContext(c) // store내 param은 모두 소문자.
	result, ok := store.Get(userID)
	if !ok && userID == user { // setup.env에서 받아온 userID와 비교
		setAdminUser(c)
		result, ok = store.Get(userID)
	}
	return result.(map[string]string), ok
}

// 관리자 정보 set.
func setAdminUser(c echo.Context) {

	user := os.Getenv("LoginUser")
	email := os.Getenv("LoginEmail")
	pass := os.Getenv("LoginPassword")

	store := echosession.FromContext(c) // store내 param은 모두 소문자.
	obj := map[string]string{
		"userid":           user,
		"username":         user,
		"email":            email,
		"password":         pass,
		"defaultnamespage": "",
		"accesstoken":      "",
		"refreshtoken":     "",
	}
	store.Set(user, obj)
	store.Save() // 사용자정보를 따로 저장하지 않으므로 설정파일에 유저를 set.
	log.Println("Set Admin")
}

func SetStore(c echo.Context, storeKeyName string, storeKeyValue interface{}) {
	store := echosession.FromContext(c) // store내 param은 모두 소문자.
	store.Set(storeKeyName, storeKeyValue)
	store.Save()
}

// providerName 소문자로
func GetProviderName(provider string) string {
	return strings.ToLower(provider)
}

// MCIS 상태값의 앞부분만 사용. 소문자로
func GetMcisStatus(mcisStatus string) string {
	statusArr := strings.Split(mcisStatus, "-")
	returnStatus := strings.ToLower(statusArr[0])

	if returnStatus == MCIS_STATUS_RUNNING {
		returnStatus = "running"
	} else if returnStatus == MCIS_STATUS_INCLUDE {
		returnStatus = "stop"
	} else if returnStatus == MCIS_STATUS_SUSPENDED {
		returnStatus = "stop"
	} else if returnStatus == MCIS_STATUS_TERMINATED {
		returnStatus = "terminate"
	} else if returnStatus == MCIS_STATUS_PARTIAL {
		returnStatus = "stop"
	} else if returnStatus == MCIS_STATUS_ETC {
		returnStatus = "stop"
	} else {
		returnStatus = "stop"
	}
	return returnStatus
}

func GetMcksStatus(mcksStatus string) string {
	statusArr := strings.Split(mcksStatus, "-")
	returnStatus := strings.ToLower(statusArr[0])

	if returnStatus == MCKS_STATUS_RUNNING {
		returnStatus = "running"
	} else if returnStatus == MCKS_STATUS_INCLUDE {
		returnStatus = "stop"
	} else if returnStatus == MCKS_STATUS_SUSPENDED {
		returnStatus = "stop"
	} else if returnStatus == MCKS_STATUS_TERMINATED {
		returnStatus = "terminate"
	} else if returnStatus == MCKS_STATUS_PARTIAL {
		returnStatus = "stop"
	} else if returnStatus == MCKS_STATUS_ETC {
		returnStatus = "stop"
	} else {
		returnStatus = "stop"
	}
	return returnStatus
}

// VM 상태를 UI에서 표현하는 방식으로 변경
func GetVmStatus(vmStatus string) string {
	returnVmStatus := strings.ToLower(vmStatus) // 소문자로 변환

	if returnVmStatus == VM_STATUS_RUNNING {
		returnVmStatus = VM_STATUS_RUNNING
		// }else if vmStatus == util.VM_STATUS_RESUMING {
		// 	vmStatusResuming++
	} else if returnVmStatus == VM_STATUS_INCLUDE {
		returnVmStatus = VM_STATUS_INCLUDE
	} else if returnVmStatus == VM_STATUS_SUSPENDED || returnVmStatus == VM_STATUS_STOPPED {
		returnVmStatus = VM_STATUS_SUSPENDED
	} else if returnVmStatus == VM_STATUS_TERMINATED {
		returnVmStatus = VM_STATUS_TERMINATED
		// }else if returnVmStatus == util.VM_STATUS_UNDEFINED {
		// 	vmStatusUndefined++
		// }else if returnVmStatus == util.VM_STATUS_PARTIAL {
		// 	vmStatusPartial++
	} else {
		returnVmStatus = VM_STATUS_ETC
	}
	return returnVmStatus
}

func GetVmConnectionName(vmConnectionName string) string {
	return strings.ToLower(vmConnectionName)
}

// Json형태의 obj를 map으로 형 변환
func StructToMapByJson(obj interface{}) (newMap map[string]interface{}, err error) {
	data, err := json.Marshal(obj) // Convert to a json

	if err != nil {
		return
	}

	err = json.Unmarshal(data, &newMap) // Convert to a map
	return
}

// func StructToMap(i interface{}) (values url.Values) {
// 	values = map[string]
// 	iVal := reflect.ValueOf(i).Elem()
// 	typ := iVal.Type()
// 	for i := 0; i < iVal.NumField(); i++ {
// 		f := iVal.Field(i)
// 		// You ca use tags here...
// 		// tag := typ.Field(i).Tag.Get("tagname")
// 		// Convert each type into a string for the url.Values string map
// 		var v string
// 		switch f.Interface().(type) {
// 		case int, int8, int16, int32, int64:
// 			v = strconv.FormatInt(f.Int(), 10)
// 		case uint, uint8, uint16, uint32, uint64:
// 			v = strconv.FormatUint(f.Uint(), 10)
// 		case float32:
// 			v = strconv.FormatFloat(f.Float(), 'f', 4, 32)
// 		case float64:
// 			v = strconv.FormatFloat(f.Float(), 'f', 4, 64)
// 		case []byte:
// 			v = string(f.Bytes())
// 		case string:
// 			v = f.String()
// 		}
// 		values.Set(typ.Field(i).Name, v)
// 	}
// 	return
// }
