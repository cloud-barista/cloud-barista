package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/cloud-barista/cb-webtool/src/service"
	"github.com/labstack/echo"
)

type RespPublicIPInfo struct {
	PublicIPInfo []struct {
		PublicIp string `json:"publicIP"`
		Status   string `json:"status"`
		VMID     string `json:"id"`
		VMName   string `json:"name"`
	} `json:"vm"`
}

func GlobalDashBoard(c echo.Context) error {
	comURL := GetCommonURL()
	apiInfo := AuthenticationHandler()
	nsCnt := service.GetNSCnt()
	fmt.Println("=========== DashBoard start ==============")
	if loginInfo := CallLoginInfo(c); loginInfo.Username != "" {
		nameSpace := GetNameSpaceToString(c)
		if nameSpace != "" {
			fmt.Println("Namespace : ", nameSpace)
			return c.Render(http.StatusOK, "Dashboard_Global.html", map[string]interface{}{
				"LoginInfo": loginInfo,
				"NameSpace": nameSpace,
				"comURL":    comURL,
				"apiInfo":   apiInfo,
				"nsCnt":     nsCnt,
			})
		} else {
			return c.Redirect(http.StatusTemporaryRedirect, "/NS/reg")
		}

	}

	return c.Redirect(http.StatusTemporaryRedirect, "/login")
}
func DashBoard(c echo.Context) error {
	comURL := GetCommonURL()
	apiInfo := AuthenticationHandler()
	fmt.Println("=========== DashBoard start ==============")
	if loginInfo := CallLoginInfo(c); loginInfo.Username != "" {
		nameSpace := GetNameSpaceToString(c)
		if nameSpace != "" {
			fmt.Println("Namespace : ", nameSpace)
			return c.Render(http.StatusOK, "dashboard.html", map[string]interface{}{
				"LoginInfo": loginInfo,
				"NameSpace": nameSpace,
				"comURL":    comURL,
				"apiInfo":   apiInfo,
			})
		} else {
			return c.Redirect(http.StatusTemporaryRedirect, "/NS/reg")
		}

	}

	return c.Redirect(http.StatusTemporaryRedirect, "/login")
}

func NSDashBoard(c echo.Context) error {
	comURL := GetCommonURL()
	apiInfo := AuthenticationHandler()
	nsCnt := service.GetNSCnt()
	fmt.Println("=========== DashBoard start ==============")
	if loginInfo := CallLoginInfo(c); loginInfo.Username != "" {
		nameSpace := GetNameSpaceToString(c)
		if nameSpace != "" {
			fmt.Println("Namespace : ", nameSpace)
			return c.Render(http.StatusOK, "Dashboard_Ns.html", map[string]interface{}{
				"LoginInfo": loginInfo,
				"NameSpace": nameSpace,
				"comURL":    comURL,
				"apiInfo":   apiInfo,
				"nsCnt":     nsCnt,
			})
		} else {
			return c.Redirect(http.StatusTemporaryRedirect, "/NS/reg")
		}

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

func Map(c echo.Context) error {
	fmt.Println("=========== MAP ==============")

	return c.Render(http.StatusOK, "Map.html", nil)
}

func GeoInfo(c echo.Context) error {
	//goroutine sync wg
	var wg sync.WaitGroup
	nameSpace := GetNameSpaceToString(c)
	comURL := GetCommonURL()
	//apiInfo := AuthenticationHandler()
	tumble_url := comURL.TumbleBugURL

	mcis_id := c.Param("mcis_id")
	url := tumble_url + "/ns/" + nameSpace + "/mcis/" + mcis_id
	fmt.Println("===========")
	fmt.Println("=========== GetGeoINFO ==============")
	fmt.Println("=========== GetGeoINFO request URL : ", url)
	body := service.HttpGetHandler(url)

	defer body.Close()

	publicIpInfo := RespPublicIPInfo{}
	json.NewDecoder(body).Decode(&publicIpInfo)
	fmt.Println("================mcis ID info ===============")
	fmt.Println("Public Info : ", publicIpInfo)
	var ipStackInfo []service.IPStackInfo

	for _, item := range publicIpInfo.PublicIPInfo {
		wg.Add(1)

		go service.GetGeoMetryInfo(&wg, item.PublicIp, item.Status, item.VMID, item.VMName, &ipStackInfo)

	}
	wg.Wait()
	fmt.Println("DashBoard ipStackInfo  : ", ipStackInfo)
	return c.JSON(http.StatusOK, ipStackInfo)
}
