package controller

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	echosession "github.com/go-session/echo-session"
	"github.com/labstack/echo"
)

type LoginInfo struct {
	Username  string
	NameSpace string
}

type CredentialInfo struct {
	Username string
	Password string
}

func GetCredentialInfo(c echo.Context, username string) CredentialInfo {
	store := echosession.FromContext(c)
	getObj, ok := store.Get(username)
	if !ok {
		return CredentialInfo{}
	}
	result := getObj.(map[string]string)
	credentialInfo := CredentialInfo{
		Username: result["username"],
		Password: result["password"],
	}
	return credentialInfo
}

// func SetLoginInfo(c echo.Context) LoginInfo {
// 	store := echosession.FromContext(c)
// 	nsList := service.GetNSList()
// 	store.Set("username")
// }

func SetNameSpace(c echo.Context) error {
	fmt.Println("====== SET NAME SPACE ========")
	store := echosession.FromContext(c)
	ns := c.Param("nsid")
	fmt.Println("SetNameSpaceID : ", ns)
	store.Set("namespace", ns)
	err := store.Save()
	res := map[string]string{
		"message": "success",
	}
	if err != nil {
		res["message"] = "fail"
		return c.JSON(http.StatusNotAcceptable, res)
	}
	return c.JSON(http.StatusOK, res)
}

func GetNameSpace(c echo.Context) error {
	fmt.Println("====== GET NAME SPACE ========")
	store := echosession.FromContext(c)

	getInfo, ok := store.Get("namespace")
	if !ok {
		return c.JSON(http.StatusNotAcceptable, map[string]string{
			"message": "Not Exist",
		})
	}
	nsId := getInfo.(string)

	res := map[string]string{
		"message": "success",
		"nsID":    nsId,
	}

	return c.JSON(http.StatusOK, res)
}

func GetNameSpaceToString(c echo.Context) string {
	fmt.Println("====== GET NAME SPACE ========")
	store := echosession.FromContext(c)

	getInfo, ok := store.Get("namespace")
	if !ok {
		return ""
	}
	nsId := getInfo.(string)

	return nsId
}

func CallLoginInfo(c echo.Context) LoginInfo {
	store := echosession.FromContext(c)
	getUser, ok := store.Get("username")
	if !ok {
		fmt.Println("========= CallLoginInfo Nothing =========")
		return LoginInfo{}
	}
	fmt.Println("GETUSER : ", getUser.(string))
	getObj, ok := store.Get(getUser.(string))

	if !ok {
		return LoginInfo{}
	}

	result := getObj.(map[string]string)
	loginInfo := LoginInfo{
		Username:  result["username"],
		NameSpace: result["namespace"],
	}
	getNs, ok := store.Get("namespace")
	if !ok {
		return loginInfo
	}
	loginInfo.NameSpace = getNs.(string)

	return loginInfo

}

func LoginCheck(c echo.Context) bool {
	store := echosession.FromContext(c)

	inputName := c.FormValue("username")
	inputPass := c.FormValue("password")

	getInfo, ok := store.Get(inputName)
	if !ok {
		return false
	}
	result := getInfo.(map[string]string)
	if result["password"] == inputPass && result["username"] == inputName {
		return true
	}

	return false
}

func MakeNameSpace(name string) string {
	now := time.Now()
	nanos := strconv.FormatInt(now.UnixNano(), 10)

	result := name + "-" + nanos
	fmt.Println("makeNameSpace : ", result)
	return result
}

// func RequestTumBleBug(method string, url string, s ) {
// 	proxyReq, err := http.NewRequest(method, url, nil)
// 	if err != nil {
// 		//log.Fatal(err)
// 	}
// 	client := &http.Client{}
// 	proxyRes, err := client.Do(proxyReq)
// 	if err != nil {
// 		//log.Fatal(err)
// 	}

// 	defer proxyRes.Body.Close()
// 	var cInfo []connectionInfo
// 	e := json.NewDecoder(proxyRes.Body).Decode(&cInfo)
// 	if e != nil {
// 		//http.Error(w, e.Error(), http.StatusBadRequest)
// 		//log.Fatal(e)
// 	}
// 	fmt.Println("bind :", cInfo[0])
// 	spew.Dump(cInfo)
// }

// func requestPost
