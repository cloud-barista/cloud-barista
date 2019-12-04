package controller

import (
	"net/http"

	"github.com/labstack/echo"
)

func MornitoringListForm(c echo.Context) error {
	if loginInfo := CallLoginInfo(c); loginInfo.Username != "" {
		namespace := GetNameSpaceToString(c)
		return c.Render(http.StatusOK, "Monitoring.html", map[string]interface{}{
			"LoginInfo": loginInfo,
			"NameSpace": namespace,
		})

	}

	//return c.Render(http.StatusOK, "MCISlist.html", nil)
	return c.Redirect(http.StatusTemporaryRedirect, "/login")
}
