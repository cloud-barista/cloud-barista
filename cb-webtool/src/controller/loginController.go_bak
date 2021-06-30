package controller

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	echosession "github.com/go-session/echo-session"
	"github.com/labstack/echo"
)

type ReqInfo struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

func LoginForm(c echo.Context) error {
	fmt.Println("============== Login Form ===============")
	comURL := GetCommonURL()
	apiInfo := AuthenticationHandler()
	return c.Render(http.StatusOK, "login.html", map[string]interface{}{
		"comURL":  comURL,
		"apiInfo": apiInfo,
	})
}

func LogoutForm(c echo.Context) error {
	fmt.Println("============== Logout Form ===============")
	//comURL := GetCommonURL()
	return c.Render(http.StatusOK, "logout.html", nil)
}

func RegUserConrtoller(c echo.Context) error {
	//comURL := GetCommonURL()

	user := os.Getenv("LoginEmail")
	pass := os.Getenv("LoginPassword")

	store := echosession.FromContext(c)
	obj := map[string]string{
		"username": user,
		"password": pass,
	}
	store.Set(user, obj)
	err := store.Save()
	if err != nil {
		return c.JSON(http.StatusServiceUnavailable, map[string]string{
			"message": "Fail",
		})
	}
	return c.JSON(http.StatusOK, map[string]string{
		"message": "SUCCESS",
		"user":    user,
	})

}

func LoginController(c echo.Context) error {
	store := echosession.FromContext(c)
	reqInfo := new(ReqInfo)
	comURL := GetCommonURL()
	apiInfo := AuthenticationHandler()
	if err := c.Bind(reqInfo); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "fail",
			"comURL":  comURL,
			"apiInfo": apiInfo,
		})
	}
	getUser := strings.TrimSpace(reqInfo.UserName)
	getPass := strings.TrimSpace(reqInfo.Password)
	fmt.Println("getUser & getPass : ", getUser, getPass)

	get, ok := store.Get(getUser)
	fmt.Println("GEt USER:", get)
	if !ok {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"message": " 정보가 없으니 다시 등록 해라",
			"status":  "fail",
			"comURL":  comURL,
			"apiInfo": apiInfo,
		})
	}
	//result := map[string]string{}
	result := get.(map[string]string)
	fmt.Println("result mapping : ", result)
	// for k, v := range get.(map[string]string) {
	// 	fmt.Println(k, v)
	// 	result[k] = v

	// }

	fmt.Println("result : ", result["password"])
	if result["password"] == getPass && result["username"] == getUser {
		store.Set("username", result["username"])
		store.Save()
		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "Login Success",
			//	"nameSpace": result["namespace"],
			"status":  "success",
			"comURL":  comURL,
			"apiInfo": apiInfo,
		})
	} else {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "wrong password of ID",
			"status":  "fail",
			"comURL":  comURL,
			"apiInfo": apiInfo,
		})
	}

	// var result map[string]string
	// for k, item := range getObj {
	// 	fmt.Println("GetItem : ", item)
	// 	result[k] = item
	// }
	fmt.Println("getObj :", get)
	// if sesEmail := session.Get(getUser); sesEmail != nil {
	// 	if sesEmail == getUser {

	// 	}
	// }

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success",
		"status":  "200",
		"comURL":  comURL,
		"apiInfo": apiInfo,
	})
}

// 여기서 둘다 다 되게 처리 해야 한다.
// 로그인체크와, ns check
// func LoginProc(c echo.Context) error {

// 	inputName := c.FormValue("username")
// 	inputPass := c.FormValue("password")
// 	//username에저장되어 있는 크리덴셜 정보를 가져 오자.
// 	credentialInfo := GetCredentialInfo(c, inputName)
// 	if credentialInfo.Username == inputName && credentialInfo.Password == inputPass {

// 	}
// }
