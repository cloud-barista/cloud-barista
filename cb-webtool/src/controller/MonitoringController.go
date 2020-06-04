package controller

import (
	"net/http"

	"github.com/labstack/echo"
)

func MornitoringListForm(c echo.Context) error {
	comURL := GetCommonURL()
	if loginInfo := CallLoginInfo(c); loginInfo.Username != "" {
		namespace := GetNameSpaceToString(c)
		return c.Render(http.StatusOK, "Monitoring.html", map[string]interface{}{
			"LoginInfo": loginInfo,
			"NameSpace": namespace,
			"comURL":    comURL,
		})

	}

	//return c.Render(http.StatusOK, "MCISlist.html", nil)
	return c.Redirect(http.StatusTemporaryRedirect, "/login")
}

func AgentRegForm(c echo.Context) error {
	comURL := GetCommonURL()
	mcis_id := c.Param("mcis_id")
	vm_id := c.Param("vm_id")
	public_ip := c.Param("public_ip")

	if loginInfo := CallLoginInfo(c); loginInfo.Username != "" {
		namespace := GetNameSpaceToString(c)
		return c.Render(http.StatusOK, "InstallAgent.html", map[string]interface{}{
			"LoginInfo": loginInfo,
			"NameSpace": namespace,
			"comURL":    comURL,
			"mcis_id":   mcis_id,
			"vm_id":     vm_id,
			"public_ip": public_ip,
		})

	}

	//return c.Render(http.StatusOK, "MCISlist.html", nil)
	return c.Redirect(http.StatusTemporaryRedirect, "/login")
}
