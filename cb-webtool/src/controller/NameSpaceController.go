package controller

import (
	"fmt"
	_ "fmt"
	"net/http"
	_ "net/http"

	"github.com/cloud-barista/cb-webtool/src/service"
	"github.com/labstack/echo"
)

func NsRegController(c echo.Context) error {
	username := c.FormValue("username")
	description := c.FormValue("description")

	fmt.Println("NSRegController : ", username, description)
	return nil
}

func NsRegForm(c echo.Context) error {
	comURL := GetCommonURL()
	apiInfo := AuthenticationHandler()
	if loginInfo := CallLoginInfo(c); loginInfo.Username != "" {
		return c.Render(http.StatusOK, "NSRegister.html", map[string]interface{}{
			"LoginInfo": loginInfo,
			"comURL":    comURL,
			"apiInfo":   apiInfo,
		})
	}
	return c.Redirect(http.StatusTemporaryRedirect, "/login")
	//return c.Render(http.StatusOK, "NSRegister.html", nil)
}

func NsListForm(c echo.Context) error {
	fmt.Println("=============start NsListForm =============")
	comURL := GetCommonURL()
	apiInfo := AuthenticationHandler()
	loginInfo := CallLoginInfo(c)
	if loginInfo.Username != "" {
		fmt.Println("=============start GetNSList =============")
		nsList := service.GetNSList()
		fmt.Println("=============start GetNSList =============", nsList)
		if nsList != nil {
			return c.Render(http.StatusOK, "NameSpace.html", map[string]interface{}{
				"LoginInfo": loginInfo,
				"NSList":    nsList,
				"comURL":    comURL,
				"apiInfo":   apiInfo,
			})
		} else {
			return c.Redirect(http.StatusTemporaryRedirect, "/NS/reg")
		}

	}

	fmt.Println("LoginInfo : ", loginInfo)
	//return c.Redirect(http.StatusPermanentRedirect, "/login")
	return c.Redirect(http.StatusTemporaryRedirect, "/login")
}
