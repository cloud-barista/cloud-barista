package controller

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
)

func ResourceBoard(c echo.Context) error {
	fmt.Println("=========== ResourceBoard start ==============")
	if loginInfo := CallLoginInfo(c); loginInfo.Username != "" {
		nameSpace := GetNameSpaceToString(c)
		fmt.Println("Namespace : ", nameSpace)
		return c.Render(http.StatusOK, "ResourceBoard.html", map[string]interface{}{
			"LoginInfo": loginInfo,
			"NameSpace": nameSpace,
		})

	}

	return c.Redirect(http.StatusTemporaryRedirect, "/login")
}
