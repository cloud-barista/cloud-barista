package controller

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
)

func DashBoard(c echo.Context) error {
	fmt.Println("=========== DashBoard start ==============")
	if loginInfo := CallLoginInfo(c); loginInfo.Username != "" {
		nameSpace := GetNameSpaceToString(c)
		fmt.Println("Namespace : ", nameSpace)
		return c.Render(http.StatusOK, "dashboard.html", map[string]interface{}{
			"LoginInfo": loginInfo,
			"NameSpace": nameSpace,
		})

	}

	return c.Redirect(http.StatusTemporaryRedirect, "/login")
}

func IndexController(c echo.Context) error {
	fmt.Println("=========== DashBoard start ==============")
	if loginInfo := CallLoginInfo(c); loginInfo.Username != "" {

		return c.Redirect(http.StatusTemporaryRedirect, "/dashboard")

	}
	fmt.Println("=========== Index Controller nothing ==============")
	return c.Redirect(http.StatusTemporaryRedirect, "/login")
}
func PopSpec(c echo.Context) error {
	fmt.Println("=========== popup ==============")

	return c.Render(http.StatusOK, "PopSpec.html", nil)
}
