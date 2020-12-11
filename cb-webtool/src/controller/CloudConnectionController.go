package controller

import (
	"fmt"

	"net/http"

	"github.com/cloud-barista/cb-webtool/src/service"
	"github.com/labstack/echo"
	//"github.com/davecgh/go-spew/spew"
)

// Driver Contorller
func DriverRegController(c echo.Context) error {

	username := c.FormValue("username")
	description := c.FormValue("description")

	fmt.Println("NSRegController : ", username, description)
	return nil
}

func DriverRegForm(c echo.Context) error {
	comURL := GetCommonURL()
	apiInfo := AuthenticationHandler()
	if loginInfo := CallLoginInfo(c); loginInfo.Username != "" {
		return c.Render(http.StatusOK, "DriverRegister.html", map[string]interface{}{
			"LoginInfo": loginInfo,
			"comURL":    comURL,
			"apiInfo":   apiInfo,
		})
	}
	return c.Redirect(http.StatusTemporaryRedirect, "/login")
}

func DriverListForm(c echo.Context) error {
	fmt.Println("=============start NsListForm =============")
	comURL := GetCommonURL()
	apiInfo := AuthenticationHandler()
	loginInfo := CallLoginInfo(c)

	if loginInfo.Username != "" {
		//nsList := service.GetDriverList()
		return c.Render(http.StatusOK, "DriverList.html", map[string]interface{}{
			"LoginInfo": loginInfo,
			"comURL":    comURL,
			"apiInfo":   apiInfo,
			//"NSList": nsList,
		})
	}

	fmt.Println("LoginInfo : ", loginInfo)

	return c.Redirect(http.StatusTemporaryRedirect, "/login")
}

//Credential Controller
func CredertialRegForm(c echo.Context) error {
	comURL := GetCommonURL()
	apiInfo := AuthenticationHandler()
	if loginInfo := CallLoginInfo(c); loginInfo.Username != "" {
		return c.Render(http.StatusOK, "CredentialRegister.html", map[string]interface{}{
			"LoginInfo": loginInfo,
			"comURL":    comURL,
			"apiInfo":   apiInfo,
		})
	}
	return c.Redirect(http.StatusTemporaryRedirect, "/login")
}

func CredertialListForm(c echo.Context) error {

	fmt.Println("=============start CredertialRegForm =============")
	loginInfo := CallLoginInfo(c)
	comURL := GetCommonURL()
	apiInfo := AuthenticationHandler()
	if loginInfo.Username != "" {
		//nsList := service.GetCredentialList()
		return c.Render(http.StatusOK, "CredentialList.html", map[string]interface{}{
			"LoginInfo": loginInfo,
			"comURL":    comURL,
			"apiInfo":   apiInfo,
			// "NSList": nsList,
		})
	}

	fmt.Println("LoginInfo : ", loginInfo)
	return c.Redirect(http.StatusTemporaryRedirect, "/login")

}

//Region Controller
func RegionRegForm(c echo.Context) error {
	comURL := GetCommonURL()
	apiInfo := AuthenticationHandler()
	if loginInfo := CallLoginInfo(c); loginInfo.Username != "" {
		return c.Render(http.StatusOK, "RegionRegister.html", map[string]interface{}{
			"LoginInfo": loginInfo,
			"comURL":    comURL,
			"apiInfo":   apiInfo,
		})
	}
	// return c.Redirect(http.StatusPermanentRedirect, "/login")
	return c.Redirect(http.StatusTemporaryRedirect, "/login")
}

func RegionListForm(c echo.Context) error {
	comURL := GetCommonURL()
	loginInfo := CallLoginInfo(c)
	apiInfo := AuthenticationHandler()
	if loginInfo.Username != "" {
		nsList := service.GetRegionList()
		fmt.Println("REGION List : ", nsList)

		//spew.Dump(nsList)
		return c.Render(http.StatusOK, "RegionList.html", map[string]interface{}{
			"LoginInfo": loginInfo,
			"comURL":    comURL,
			"NSList":    nsList,
			"apiInfo":   apiInfo,
		})
	}

	fmt.Println("LoginInfo : ", loginInfo)
	return c.Redirect(http.StatusTemporaryRedirect, "/login")

}

//Connection Controller
func ConnectionRegForm(c echo.Context) error {
	comURL := GetCommonURL()
	apiInfo := AuthenticationHandler()
	if loginInfo := CallLoginInfo(c); loginInfo.Username != "" {

		return c.Render(http.StatusOK, "ConnectionRegister.html", map[string]interface{}{
			"LoginInfo": loginInfo,
			"comURL":    comURL,
			"apiInfo":   apiInfo,
		})
	}
	return c.Redirect(http.StatusTemporaryRedirect, "/login")
	//return c.Render(http.StatusOK, "RegionRegister.html", nil)
}

func ConnectionListForm(c echo.Context) error {
	comURL := GetCommonURL()
	apiInfo := AuthenticationHandler()
	loginInfo := CallLoginInfo(c)
	if loginInfo.Username != "" {
		cList := service.GetConnectionList()
		fmt.Println("=============info GetConnectionList =============", cList)
		return c.Render(http.StatusOK, "CloudConnection.html", map[string]interface{}{
			// return c.Render(http.StatusOK, "ConnectionList.html", map[string]interface{}{
			"LoginInfo": loginInfo,
			"cList":     cList,
			"comURL":    comURL,
			"apiInfo":   apiInfo,
		})
	}

	fmt.Println("LoginInfo : ", loginInfo)

	return c.Redirect(http.StatusTemporaryRedirect, "/login")
}

//Image Controller
func ImageRegForm(c echo.Context) error {
	comURL := GetCommonURL()
	apiInfo := AuthenticationHandler()
	if loginInfo := CallLoginInfo(c); loginInfo.Username != "" {
		return c.Render(http.StatusOK, "ImageRegister.html", map[string]interface{}{
			"LoginInfo": loginInfo,
			"comURL":    comURL,
			"apiInfo":   apiInfo,
		})
	}
	// return c.Redirect(http.StatusPermanentRedirect, "/login")
	return c.Redirect(http.StatusTemporaryRedirect, "/login")
}

func ImageListForm(c echo.Context) error {
	comURL := GetCommonURL()
	loginInfo := CallLoginInfo(c)
	apiInfo := AuthenticationHandler()
	if loginInfo.Username != "" {
		nsList := service.GetRegionList()
		fmt.Println("REGION List : ", nsList)

		//spew.Dump(nsList)
		return c.Render(http.StatusOK, "Resources_Image.html", map[string]interface{}{
			"LoginInfo": loginInfo,
			"comURL":    comURL,
			"NSList":    nsList,
			"apiInfo":   apiInfo,
		})
	}

	fmt.Println("LoginInfo : ", loginInfo)
	return c.Redirect(http.StatusTemporaryRedirect, "/login")

}

//VPC Controller
func VpcRegForm(c echo.Context) error {
	comURL := GetCommonURL()
	apiInfo := AuthenticationHandler()
	if loginInfo := CallLoginInfo(c); loginInfo.Username != "" {
		return c.Render(http.StatusOK, "VpcRegister.html", map[string]interface{}{
			"LoginInfo": loginInfo,
			"comURL":    comURL,
			"apiInfo":   apiInfo,
		})
	}
	// return c.Redirect(http.StatusPermanentRedirect, "/login")
	return c.Redirect(http.StatusTemporaryRedirect, "/login")
}

func VpcListForm(c echo.Context) error {
	comURL := GetCommonURL()
	loginInfo := CallLoginInfo(c)
	apiInfo := AuthenticationHandler()
	if loginInfo.Username != "" {
		nsList := service.GetRegionList()
		fmt.Println("REGION List : ", nsList)

		//spew.Dump(nsList)
		return c.Render(http.StatusOK, "Resources_Network.html", map[string]interface{}{
			"LoginInfo": loginInfo,
			"comURL":    comURL,
			"NSList":    nsList,
			"apiInfo":   apiInfo,
		})
	}

	fmt.Println("LoginInfo : ", loginInfo)
	return c.Redirect(http.StatusTemporaryRedirect, "/login")

}

// Controller
func SecurityGroupRegForm(c echo.Context) error {
	comURL := GetCommonURL()
	apiInfo := AuthenticationHandler()
	if loginInfo := CallLoginInfo(c); loginInfo.Username != "" {
		return c.Render(http.StatusOK, "SecurityGroupRegister.html", map[string]interface{}{
			"LoginInfo": loginInfo,
			"comURL":    comURL,
			"apiInfo":   apiInfo,
		})
	}
	// return c.Redirect(http.StatusPermanentRedirect, "/login")
	return c.Redirect(http.StatusTemporaryRedirect, "/login")
}

func SecurityGroupListForm(c echo.Context) error {
	comURL := GetCommonURL()
	loginInfo := CallLoginInfo(c)
	apiInfo := AuthenticationHandler()
	if loginInfo.Username != "" {
		nsList := service.GetRegionList()
		fmt.Println("REGION List : ", nsList)

		//spew.Dump(nsList)
		return c.Render(http.StatusOK, "Resources_Security.html", map[string]interface{}{
			"LoginInfo": loginInfo,
			"comURL":    comURL,
			"NSList":    nsList,
			"apiInfo":   apiInfo,
		})
	}

	fmt.Println("LoginInfo : ", loginInfo)
	return c.Redirect(http.StatusTemporaryRedirect, "/login")

}

// Controller
func SSHRegForm(c echo.Context) error {
	comURL := GetCommonURL()
	apiInfo := AuthenticationHandler()
	if loginInfo := CallLoginInfo(c); loginInfo.Username != "" {
		return c.Render(http.StatusOK, "SSHRegister.html", map[string]interface{}{
			"LoginInfo": loginInfo,
			"comURL":    comURL,
			"apiInfo":   apiInfo,
		})
	}
	// return c.Redirect(http.StatusPermanentRedirect, "/login")
	return c.Redirect(http.StatusTemporaryRedirect, "/login")
}

func SSHListForm(c echo.Context) error {
	comURL := GetCommonURL()
	loginInfo := CallLoginInfo(c)
	apiInfo := AuthenticationHandler()
	if loginInfo.Username != "" {
		nsList := service.GetRegionList()
		fmt.Println("REGION List : ", nsList)

		//spew.Dump(nsList)
		return c.Render(http.StatusOK, "Resources_Ssh.html", map[string]interface{}{
			"LoginInfo": loginInfo,
			"comURL":    comURL,
			"NSList":    nsList,
			"apiInfo":   apiInfo,
		})
	}

	fmt.Println("LoginInfo : ", loginInfo)
	return c.Redirect(http.StatusTemporaryRedirect, "/login")

}

// Controller
func SpecRegForm(c echo.Context) error {
	comURL := GetCommonURL()
	apiInfo := AuthenticationHandler()
	if loginInfo := CallLoginInfo(c); loginInfo.Username != "" {
		return c.Render(http.StatusOK, "SpecRegister.html", map[string]interface{}{
			"LoginInfo": loginInfo,
			"comURL":    comURL,
			"apiInfo":   apiInfo,
		})
	}
	// return c.Redirect(http.StatusPermanentRedirect, "/login")
	return c.Redirect(http.StatusTemporaryRedirect, "/login")
}

func SpecListForm(c echo.Context) error {
	comURL := GetCommonURL()
	loginInfo := CallLoginInfo(c)
	apiInfo := AuthenticationHandler()
	if loginInfo.Username != "" {
		nsList := service.GetRegionList()
		fmt.Println("REGION List : ", nsList)

		//spew.Dump(nsList)
		return c.Render(http.StatusOK, "Resources_Spec.html", map[string]interface{}{
			"LoginInfo": loginInfo,
			"comURL":    comURL,
			"NSList":    nsList,
			"apiInfo":   apiInfo,
		})
	}

	fmt.Println("LoginInfo : ", loginInfo)
	return c.Redirect(http.StatusTemporaryRedirect, "/login")

}
