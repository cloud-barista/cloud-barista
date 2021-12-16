package controller

import (
	// "bytes"
	// "encoding/json"
	"fmt"
	// "github.com/davecgh/go-spew/spew"
	// "io/ioutil"
	"log"
	"net/http"
	"os"
	// "strings"
	"time"

	//"github.com/foolin/echo-template"
	// echosession "github.com/go-session/echo-session"

	"github.com/dgrijalva/jwt-go"

	//"github.com/twinj/uuid"
	// "github.com/google/uuid"

	"github.com/labstack/echo"

	//db "mzc/src/databases/store"
	"github.com/cloud-barista/cb-webtool/src/model"
	"github.com/cloud-barista/cb-webtool/src/service"
	"github.com/cloud-barista/cb-webtool/src/util"
)

// Api를 통한 로그인
// Token을 생성하여 해당 token 반환
func ApiLogin(c echo.Context) error {
	// username := c.FormValue("username")
	// password := c.FormValue("password")

	params := make(map[string]string)
	if err := c.Bind(&params); err != nil {
		fmt.Println("err = ", err) // bind Error는 나지만 크게 상관없는 듯.
	}
	fmt.Println(params)

	storedUser, ok := util.GetUserInfo(c, params["userID"])
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{ //401
			"message": "Try again",
			"status":  "fail",
		})
	}
	if params["userID"] != storedUser["userid"] || params["password"] != storedUser["password"] {
		return echo.ErrUnauthorized
	}

	// Token 생성
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = params["userID"]
	claims["name"] = params["userID"]
	claims["admin"] = true
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	t, err := token.SignedString([]byte(os.Getenv("LoginAccessSecret")))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]string{
		"token": t,
	})
}

func Accessible(c echo.Context) error {
	return c.String(http.StatusOK, "Accessible")
}

// 접속후 확인용 : token이 생성된 경우 /api/auth/restricted/user", method: "get"  호출시 Welcome + name 을 return 함.
func Restricted(c echo.Context) error {
	log.Println("Restricted1")
	user := c.Get("user").(*jwt.Token)
	log.Println("Restricted2")
	claims := user.Claims.(jwt.MapClaims)
	log.Println("Restricted3")
	name := claims["name"].(string)
	log.Println("Restricted4")
	return c.String(http.StatusOK, "Welcome "+name+"!")
}

// name이 admin이고 claims["admin"] == true 면 namespace 목록을 가져온다.
func ApiUserInfo(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	id := claims["id"].(string)
	name := claims["name"].(string)
	isAdmin := claims["admin"].(bool)

	userInfo := model.LoginInfo{
		UserID:   id,
		Username: name,
	}

	fmt.Println("name="+name+", isAdmin %v ", isAdmin)
	if name == "admin" && isAdmin {
		fmt.Println("admin userinfo set")

		// 	nsList, nsStatus := service.GetNameSpaceList()
		// 	return c.JSON(http.StatusOK, map[string]interface{}{
		// 		"message":       nsStatus.Message,
		// 		"status":        nsStatus.StatusCode,
		// 		"user":          userInfo,
		// 		"NameSpaceList": nsList,
		// 	})

	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success",
		"status":  "200",
		"user":    userInfo,
	})

	// return Restricted(c)
}

func ApiNamespaceList(c echo.Context) error {
	nsList, nsStatus := service.GetNameSpaceList()
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":       nsStatus.Message,
		"status":        nsStatus.StatusCode,
		"NamespaceList": nsList,
	})

}
