package controller

import (
	"fmt"
	"log"
	"net/http"

	// "github.com/cloud-barista/cb-webtool/src/model"

	// dragonfly "github.com/cloud-barista/cb-webtool/src/model/dragonfly"
	// ladybug "github.com/cloud-barista/cb-webtool/src/model/ladybug"
	spider "github.com/cloud-barista/cb-webtool/src/model/spider"
	// tumblebug "github.com/cloud-barista/cb-webtool/src/model/tumblebug"
	tbcommon "github.com/cloud-barista/cb-webtool/src/model/tumblebug/common"
	// tbmcir "github.com/cloud-barista/cb-webtool/src/model/tumblebug/mcir"
	// tbmcis "github.com/cloud-barista/cb-webtool/src/model/tumblebug/mcis"

	service "github.com/cloud-barista/cb-webtool/src/service"

	"github.com/cloud-barista/cb-webtool/src/util"
	"github.com/labstack/echo"

	//"github.com/davecgh/go-spew/spew"
	echotemplate "github.com/foolin/echo-template"
	echosession "github.com/go-session/echo-session"
)

// 등록 form
// GET LIST
// GET Data
// POST 등록
// DELETE

// func CloudOSListForm
// CloudOS(Provider) 목록
func GetCloudOSList(c echo.Context) error {
	cloudOsList, respStatus := service.GetCloudOSList()
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success",
		"status":  respStatus,
		"cloudos": cloudOsList,
	})
}

// func GetCloudOS
// func CloudOSRegProc
// func CloudOSDelProc

// func ConnectionConfigList(c echo.Context) error {
func CloudConnectionConfigMngForm(c echo.Context) error {
	fmt.Println("ConnectionConfigList ************ : ")

	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		// Login 정보가 없으므로 login화면으로
		// return c.JSON(http.StatusBadRequest, map[string]interface{}{
		// 	"message": "invalid tumblebug connection",
		// 	"status":  "403",
		// })
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}
	fmt.Println("loginInfo ", loginInfo)

	store := echosession.FromContext(c)
	// result, ok := store.Get(paramUser)
	// storedUser := result.(map[string]string)

	cloudOsList, _ := service.GetCloudOSList()
	store.Set("cloudos", cloudOsList)
	log.Println(" cloudOsList  ", cloudOsList)

	// connectionconfigList 가져오기
	cloudConnectionConfigInfoList, cloudConnectionConfigErr := service.GetCloudConnectionConfigList()
	setCloudConnectionErr := service.SetStoreCloudConnectionConfigList(c, cloudConnectionConfigInfoList)
	if setCloudConnectionErr != nil {
		log.Println(" setCloudConnectionErr  ", setCloudConnectionErr)
	}
	// if ( cloudConnectionConfigStatus.StatusCode != 200 & cloudConnectionConfigStatus.StatusCode != 201 ) {
	// store.Set("cloudconnectionconfig", cloudConnectionConfigInfoList)
	// log.Println(" cloudconnectionconfig  ", cloudConnectionConfigInfoList)

	// regionList 가져오기
	regionList, _ := service.GetRegionList()
	setRegionErr := service.SetStoreRegionList(c, regionList)
	if setRegionErr != nil {
		log.Println(" setRegionErr  ", setRegionErr)
	}
	// store.Set("region", regionList)
	log.Println(" regionList  ", regionList)

	// credentialList 가져오기
	credentialList, _ := service.GetCredentialList()
	setCredentialErr := service.SetStoreCredentialList(c, credentialList)
	if setCredentialErr != nil {
		log.Println(" setCredentialErr  ", setCredentialErr)
	}
	// store.Set("credential", credentialList)
	log.Println(" credentialList  ", credentialList)

	// driverList 가져오기
	driverList, _ := service.GetDriverList()
	setDriverErr := service.SetStoreDriverList(c, driverList)
	if setDriverErr != nil {
		log.Println(" setDriverErr  ", setDriverErr)
	}
	// store.Set("driver", driverList)
	log.Println(" driverList  ", driverList)

	// 최신 namespacelist 가져오기
	// nsList, _ := service.GetNameSpaceList()
	nsList, _ := service.GetStoredNameSpaceList(c)
	store.Set("namespace", nsList)
	_ = store.Save()
	log.Println(" nsList  ", nsList)

	// status, filepath, return params
	return echotemplate.Render(c, http.StatusOK,
		"setting/connections/CloudConnectionConfigMng", // 파일명
		map[string]interface{}{
			"Message":                   cloudConnectionConfigErr.Message,
			"Status":                    cloudConnectionConfigErr.StatusCode,
			"LoginInfo":                 loginInfo,
			"CloudOSList":               cloudOsList,
			"NameSpaceList":             nsList,
			"CloudConnectionConfigList": cloudConnectionConfigInfoList,
			"RegionList":                regionList,
			"CredentialList":            credentialList,
			"DriverList":                driverList,
		})
}

// 현재 설정된 connection 목록
func GetCloudConnectionConfigList(c echo.Context) error {
	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	cloudConnectionConfigList, respStatus := service.GetCloudConnectionConfigList()
	setCloudConnectionConfigErr := service.SetStoreCloudConnectionConfigList(c, cloudConnectionConfigList)
	if setCloudConnectionConfigErr != nil {
		log.Println("setCloudConnectionConfigErr ", setCloudConnectionConfigErr)
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":          "success",
		"status":           respStatus,
		"ConnectionConfig": cloudConnectionConfigList,
	})
}

// cloud connection 상세정보
func GetCloudConnectionConfigData(c echo.Context) error {
	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	paramConfigName := c.Param("configName")
	cloudConnectionConfigInfo, respStatus := service.GetCloudConnectionConfigData(paramConfigName)
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":          "success",
		"status":           respStatus,
		"ConnectionConfig": cloudConnectionConfigInfo,
	})
}

// cloud connection 등록 :
func CloudConnectionConfigRegProc(c echo.Context) error {
	log.Println("ConnectionConfigRegProc : ")
	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	cloudConnectionConfigInfo := new(spider.CloudConnectionConfigInfo)
	if err := c.Bind(cloudConnectionConfigInfo); err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "fail",
			"status":  "fail",
		})
	}
	log.Println(cloudConnectionConfigInfo)
	respBody, respStatus := service.RegCloudConnectionConfig(cloudConnectionConfigInfo)
	fmt.Println("=============respBody =============", respBody)
	if respStatus.StatusCode != 200 && respStatus.StatusCode != 201 {
		return c.JSON(respStatus.StatusCode, map[string]interface{}{
			"error":  respStatus.Message,
			"status": respStatus.StatusCode,
		})
	}

	cloudConnectionConfigErr := GetCloudConnectionConfigList(c)
	if cloudConnectionConfigErr != nil {
		log.Println(" cloudConnectionConfigErr ", cloudConnectionConfigErr)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success",
		"status":  respStatus,
	})
}

// 삭제
func CloudConnectionConfigDelProc(c echo.Context) error {
	log.Println("ConnectionConfigDelProc : ")

	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		// Login 정보가 없으므로 login화면으로
		// return c.JSON(http.StatusBadRequest, map[string]interface{}{
		// 	"message": "invalid tumblebug connection",
		// 	"status":  "403",
		// })
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	paramConfigName := c.Param("configName")
	log.Println(paramConfigName)

	respBody, respStatus := service.DelCloudConnectionConfig(paramConfigName)
	fmt.Println("=============respBody =============", respBody)
	if respStatus.StatusCode != 200 && respStatus.StatusCode != 201 {
		return c.JSON(respStatus.StatusCode, map[string]interface{}{
			"error":  respStatus.Message,
			"status": respStatus.StatusCode,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success",
		"status":  respStatus,
	})
}

// func RegionListForm // Region 등록 form : maing 화면에서 popup형태로 뜸

// 현재 설정된 region 목록
func GetRegionList(c echo.Context) error {
	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		// Login 정보가 없으므로 login화면으로
		// return c.JSON(http.StatusBadRequest, map[string]interface{}{
		// 	"message": "invalid tumblebug connection",
		// 	"status":  "403",
		// })
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	regionList, respStatus := service.GetRegionList()
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success",
		"status":  respStatus,
		"Region":  regionList,
	})
}

// region 상세정보
func GetRegion(c echo.Context) error {
	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		// Login 정보가 없으므로 login화면으로
		// return c.JSON(http.StatusBadRequest, map[string]interface{}{
		// 	"message": "invalid tumblebug connection",
		// 	"status":  "403",
		// })
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	paramRegion := c.Param("region")
	resionInfo, respStatus := service.GetRegionData(paramRegion)
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success",
		"status":  respStatus,
		"Region":  resionInfo,
	})
}

// region 등록
func RegionRegProc(c echo.Context) error {
	log.Println("RegionRegProc : ")
	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		// Login 정보가 없으므로 login화면으로
		// return c.JSON(http.StatusBadRequest, map[string]interface{}{
		// 	"message": "invalid tumblebug connection",
		// 	"status":  "403",
		// })
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	regionInfo := new(spider.RegionInfo)
	if err := c.Bind(regionInfo); err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "fail",
			"status":  "fail",
		})
	}
	log.Println(regionInfo)
	respBody, respStatus := service.RegRegion(regionInfo)
	fmt.Println("=============respBody =============", respBody)
	if respStatus.StatusCode != 200 && respStatus.StatusCode != 201 {
		return c.JSON(respStatus.StatusCode, map[string]interface{}{
			"error":  respStatus.Message,
			"status": respStatus.StatusCode,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success",
		"status":  "200",
	})
}

// Region 삭제
func RegionDelProc(c echo.Context) error {
	log.Println("RegionRegProc : ")
	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		// Login 정보가 없으므로 login화면으로
		// return c.JSON(http.StatusBadRequest, map[string]interface{}{
		// 	"message": "invalid tumblebug connection",
		// 	"status":  "403",
		// })
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	paramRegion := c.Param("region")
	log.Println(paramRegion)

	respBody, respStatus := service.DelRegion(paramRegion)
	fmt.Println("=============respBody =============", respBody)
	if respStatus.StatusCode != 200 && respStatus.StatusCode != 201 {
		return c.JSON(respStatus.StatusCode, map[string]interface{}{
			"error":  respStatus.Message,
			"status": respStatus.StatusCode,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success",
		"status":  respStatus,
	})
}

// func CredentialListForm // Credential 등록 form : maing 화면에서 popup형태로 뜸

// 현재 설정된 Credential 목록 : 목록에서는 key의 value는 보여주지 않는다. ... 표시 (상세정보에서는 표시)
func GetCredentialList(c echo.Context) error {
	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		// Login 정보가 없으므로 login화면으로
		// return c.JSON(http.StatusBadRequest, map[string]interface{}{
		// 	"message": "invalid tumblebug connection",
		// 	"status":  "403",
		// })
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	// Spider호출이므로 optionParam 사용X
	credentialList, respStatus := service.GetCredentialList()
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":    "success",
		"status":     respStatus,
		"Credential": credentialList,
	})
}

// Credential 상세정보
func GetCredential(c echo.Context) error {
	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		// Login 정보가 없으므로 login화면으로
		// return c.JSON(http.StatusBadRequest, map[string]interface{}{
		// 	"message": "invalid tumblebug connection",
		// 	"status":  "403",
		// })
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	paramCredential := c.Param("credential")
	credentialInfo, respStatus := service.GetCredentialData(paramCredential)
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":    "success",
		"status":     respStatus,
		"Credential": credentialInfo,
	})
}

// Credential 등록
func CredentialRegProc(c echo.Context) error {
	log.Println("CredentialRegProc : ")
	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		// Login 정보가 없으므로 login화면으로
		// return c.JSON(http.StatusBadRequest, map[string]interface{}{
		// 	"message": "invalid tumblebug connection",
		// 	"status":  "403",
		// })
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	credentialInfo := new(spider.CredentialInfo)
	if err := c.Bind(credentialInfo); err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "fail",
			"status":  "fail",
		})
	}
	log.Println(credentialInfo)
	respBody, respStatus := service.RegCredential(credentialInfo)
	fmt.Println("=============respBody =============", respBody)
	if respStatus.StatusCode != 200 && respStatus.StatusCode != 201 {
		return c.JSON(respStatus.StatusCode, map[string]interface{}{
			"error":  respStatus.Message,
			"status": respStatus.StatusCode,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success",
		"status":  respStatus,
	})
}

// Credential 삭제
func CredentialDelProc(c echo.Context) error {
	log.Println("CredentialDelProc : ")
	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		// Login 정보가 없으므로 login화면으로
		// return c.JSON(http.StatusBadRequest, map[string]interface{}{
		// 	"message": "invalid tumblebug connection",
		// 	"status":  "403",
		// })
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	paramCredential := c.Param("credential")
	log.Println(paramCredential)

	respMessage, respStatus := service.DelCredential(paramCredential)
	fmt.Println("=============respMessage =============", respMessage)

	if respStatus.StatusCode != 200 && respStatus.StatusCode != 201 {
		return c.JSON(respStatus.StatusCode, map[string]interface{}{
			"error":  respStatus.Message,
			"status": respStatus.StatusCode,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": respMessage.Message,
		"status":  respMessage.StatusCode,
	})
}

// func DriverListForm // Driver 등록 form : maing 화면에서 popup형태로 뜸

// 현재 설정된 driver 목록
func GetDriverList(c echo.Context) error {
	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		// Login 정보가 없으므로 login화면으로
		// return c.JSON(http.StatusBadRequest, map[string]interface{}{
		// 	"message": "invalid tumblebug connection",
		// 	"status":  "403",
		// })
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	driverList, respStatus := service.GetDriverList()
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success",
		"status":  respStatus,
		"Driver":  driverList,
	})
}

// Driver 조회
func GetDriver(c echo.Context) error {
	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		// Login 정보가 없으므로 login화면으로
		// return c.JSON(http.StatusBadRequest, map[string]interface{}{
		// 	"message": "invalid tumblebug connection",
		// 	"status":  "403",
		// })
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	paramDriver := c.Param("driver")
	driverInfo, respStatus := service.GetDriverData(paramDriver)
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success",
		"status":  respStatus,
		"Driver":  driverInfo,
	})
}

// Driver 등록
// func DriverRegProc
func DriverRegProc(c echo.Context) error {
	log.Println("DriverRegProc : ")
	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		// Login 정보가 없으므로 login화면으로
		// return c.JSON(http.StatusBadRequest, map[string]interface{}{
		// 	"message": "invalid tumblebug connection",
		// 	"status":  "403",
		// })
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	driverInfo := new(spider.DriverInfo)
	if err := c.Bind(driverInfo); err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "fail",
			"status":  "fail",
		})
	}
	log.Println(driverInfo)
	respBody, respStatus := service.RegDriver(driverInfo)
	fmt.Println("=============respBody =============", respBody)
	if respStatus.StatusCode != 200 && respStatus.StatusCode != 201 {
		return c.JSON(respStatus.StatusCode, map[string]interface{}{
			"error":  respStatus.Message,
			"status": respStatus.StatusCode,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success",
		"status":  respStatus.StatusCode,
	})
}

// Driver 삭제
func DriverDelProc(c echo.Context) error {
	log.Println("DriverDelProc : ")
	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		// Login 정보가 없으므로 login화면으로
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	paramDriver := c.Param("driver")
	log.Println(paramDriver)

	respBody, respStatus := service.DelDriver(paramDriver)
	fmt.Println("=============respBody =============", respBody)

	if respStatus.StatusCode != 200 && respStatus.StatusCode != 201 {
		return c.JSON(respStatus.StatusCode, map[string]interface{}{
			"error":  respStatus.Message,
			"status": respStatus.StatusCode,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success",
		"status":  respStatus.StatusCode,
	})
}

///////// Config ////////////

// 현재 설정된 Config 목록
func GetConfigList(c echo.Context) error {
	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		// Login 정보가 없으므로 login화면으로
		// return c.JSON(http.StatusBadRequest, map[string]interface{}{
		// 	"message": "invalid tumblebug connection",
		// 	"status":  "403",
		// })
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	configList, respStatus := service.GetConfigList()
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success",
		"status":  respStatus,
		"Config":  configList,
	})
}

// Config 조회
func GetConfig(c echo.Context) error {
	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		// Login 정보가 없으므로 login화면으로
		// return c.JSON(http.StatusBadRequest, map[string]interface{}{
		// 	"message": "invalid tumblebug connection",
		// 	"status":  "403",
		// })
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	paramConfigID := c.Param("configID")
	configInfo, respStatus := service.GetDriverData(paramConfigID)
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success",
		"status":  respStatus,
		"Config":  configInfo,
	})
}

// Config 등록
// func ConfigRegProc
func ConfigRegProc(c echo.Context) error {
	log.Println("ConfigRegProc : ")
	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		// Login 정보가 없으므로 login화면으로
		// return c.JSON(http.StatusBadRequest, map[string]interface{}{
		// 	"message": "invalid tumblebug connection",
		// 	"status":  "403",
		// })
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	configInfo := new(spider.ConfigInfo)
	if err := c.Bind(configInfo); err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "fail",
			"status":  "fail",
		})
	}
	log.Println(configInfo)
	respBody, respStatus := service.RegConfig(configInfo)
	fmt.Println("=============respBody =============", respBody)
	if respStatus.StatusCode != 200 && respStatus.StatusCode != 201 {
		return c.JSON(respStatus.StatusCode, map[string]interface{}{
			"error":  respStatus.Message,
			"status": respStatus.StatusCode,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success",
		"status":  respStatus.StatusCode,
	})
}

// Config 삭제
func ConfigDelProc(c echo.Context) error {
	log.Println("ConfigDelProc : ")
	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		// Login 정보가 없으므로 login화면으로
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	paramConfigID := c.Param("configID")
	log.Println(paramConfigID)

	respBody, respStatus := service.DelConfig(paramConfigID)
	fmt.Println("=============respBody =============", respBody)
	if respStatus.StatusCode != 200 && respStatus.StatusCode != 201 {
		return c.JSON(respStatus.StatusCode, map[string]interface{}{
			"error":  respStatus.Message,
			"status": respStatus.StatusCode,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success",
		"status":  respStatus.StatusCode,
	})
}

///////// condif //////////
// Cloud 연결정보 표시(driver)
func ConnectionList(c echo.Context) error {
	fmt.Println("ConnectionList ************ : ")

	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		// Login 정보가 없으므로 login화면으로
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	nsList, nsStatus := service.GetNameSpaceList()
	if nsStatus.StatusCode == 500 {
		log.Println(" nsErr  ", nsStatus)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": nsStatus.Message,
			"status":  nsStatus.StatusCode,
		})
	}

	// namespace 가 없으면 1개를 기본으로 생성한다.
	if len(nsList) == 0 {
		// create default namespace
		nameSpaceInfo := new(tbcommon.TbNsInfo)
		nameSpaceInfo.Name = "NS-01" // default namespace name
		nameSpaceInfo.Description = "default name space name"
		respBody, respStatus := service.RegNameSpace(nameSpaceInfo)
		log.Println(" respBody  ", respBody)
		if respStatus.StatusCode != 200 && respStatus.StatusCode != 201 {
			log.Println(" respStatus  ", respStatus)
			return c.JSON(respStatus.StatusCode, map[string]interface{}{
				"error":  respStatus.Message,
				"status": respStatus.StatusCode,
			})
		}

		// 처음생성했으므로 connection부터
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	} else if len(nsList) == 1 {
		defaultNameSpace := nsList[0]
		// for _, item := range nsList {
		// 	fmt.Println("ID : ", item.ID)
		// }
		loginInfo.DefaultNameSpaceID = defaultNameSpace.ID
		loginInfo.DefaultNameSpaceName = defaultNameSpace.Name
	}

	return echotemplate.Render(c, http.StatusOK,
		"setting/connections/cloudconnectionconfig", //파일명
		map[string]interface{}{
			"LoginInfo":     loginInfo,
			"NameSpaceList": nsList,
		})
	// return echotemplate.Render(c, http.StatusOK, "CloudConnection", nil)// -> file not found 남. 경로 제대로 적을 것.
}

// Driver Contorller
func DriverRegController(c echo.Context) error {
	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		// Login 정보가 없으므로 login화면으로
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	userID := c.FormValue("userID")
	description := c.FormValue("description")

	fmt.Println("DriverRegController : ", userID, description)
	return nil
}

// deprecated
// func DriverRegForm(c echo.Context) error {
// 	return c.Redirect(http.StatusTemporaryRedirect, "/login")
// }

// deprecated
// func DriverListForm(c echo.Context) error {
// 	return c.Redirect(http.StatusTemporaryRedirect, "/login")
// }

// deprecated
//Credential Controller
// func CredentialRegForm(c echo.Context) error {
// 	return c.Redirect(http.StatusTemporaryRedirect, "/login")
// }

// deprecated
// func CredertialListForm(c echo.Context) error {
// 	return c.Redirect(http.StatusTemporaryRedirect, "/login")
// }

// deprecated
//Region Controller
// func RegionRegForm(c echo.Context) error {
// 	return c.Redirect(http.StatusTemporaryRedirect, "/login")
// }

// deprecated
// func RegionListForm(c echo.Context) error {
// 	return c.Redirect(http.StatusTemporaryRedirect, "/login")
// }

// deprecated
//Connection Controller
// func ConnectionRegForm(c echo.Context) error {
// 	return c.Redirect(http.StatusTemporaryRedirect, "/login")
// }

// deprecated
// func ConnectionListForm(c echo.Context) error {
// 	return c.Redirect(http.StatusTemporaryRedirect, "/login")
// }

//Image Controller
func ImageRegForm(c echo.Context) error {
	comURL := service.GetCommonURL()
	apiInfo := util.AuthenticationHandler()
	if loginInfo := service.CallLoginInfo(c); loginInfo.UserID != "" {
		return c.Render(http.StatusOK, "ImageRegister.html", map[string]interface{}{
			"LoginInfo": loginInfo,
			"comURL":    comURL,
			"apiInfo":   apiInfo,
		})
	}
	// return c.Redirect(http.StatusPermanentRedirect, "/login")
	return c.Redirect(http.StatusTemporaryRedirect, "/login")
}

// VPC Controller
// func VpcRegForm(c echo.Context) error {
// 	comURL := service.GetCommonURL()
// 	apiInfo := service.AuthenticationHandler()
// 	if loginInfo := service.CallLoginInfo(c); loginInfo.UserID != "" {
// 		return c.Render(http.StatusOK, "VpcRegister.html", map[string]interface{}{
// 			"LoginInfo": loginInfo,
// 			"comURL":    comURL,
// 			"apiInfo":   apiInfo,
// 		})
// 	}
// 	// return c.Redirect(http.StatusPermanentRedirect, "/login")
// 	return c.Redirect(http.StatusTemporaryRedirect, "/login")
// }

// func ImageListForm(c echo.Context) error {
// 	comURL := service.GetCommonURL()
// 	loginInfo := service.CallLoginInfo(c)
// 	apiInfo := service.AuthenticationHandler()
// 	if loginInfo.UserID != "" {
// 		nsList := service.GetRegionList()
// 		fmt.Println("REGION List : ", nsList)

// 		//spew.Dump(nsList)
// 		return c.Render(http.StatusOK, "Resources_Image.html", map[string]interface{}{
// 			"LoginInfo": loginInfo,
// 			"comURL":    comURL,
// 			"NSList":    nsList,
// 			"apiInfo":   apiInfo,
// 		})
// 	}

// 	fmt.Println("LoginInfo : ", loginInfo)
// 	return c.Redirect(http.StatusTemporaryRedirect, "/login")

// }

// Controller
// func SecurityGroupRegForm(c echo.Context) error {
// 	comURL := service.GetCommonURL()
// 	apiInfo := service.AuthenticationHandler()
// 	if loginInfo := service.CallLoginInfo(c); loginInfo.UserID != "" {
// 		return c.Render(http.StatusOK, "SecurityGroupRegister.html", map[string]interface{}{
// 			"LoginInfo": loginInfo,
// 			"comURL":    comURL,
// 			"apiInfo":   apiInfo,
// 		})
// 	}
// 	// return c.Redirect(http.StatusPermanentRedirect, "/login")
// 	return c.Redirect(http.StatusTemporaryRedirect, "/login")
// }

// func SecurityGroupListForm(c echo.Context) error {
// 	comURL := service.GetCommonURL()
// 	loginInfo := service.CallLoginInfo(c)
// 	apiInfo := service.AuthenticationHandler()
// 	if loginInfo.UserID != "" {
// 		nsList := service.GetRegionList()
// 		fmt.Println("REGION List : ", nsList)

// 		//spew.Dump(nsList)
// 		return c.Render(http.StatusOK, "Resources_Security.html", map[string]interface{}{
// 			"LoginInfo": loginInfo,
// 			"comURL":    comURL,
// 			"NSList":    nsList,
// 			"apiInfo":   apiInfo,
// 		})
// 	}

// 	fmt.Println("LoginInfo : ", loginInfo)
// 	return c.Redirect(http.StatusTemporaryRedirect, "/login")

// }

// Controller
func SSHRegForm(c echo.Context) error {
	comURL := service.GetCommonURL()
	apiInfo := util.AuthenticationHandler()
	if loginInfo := service.CallLoginInfo(c); loginInfo.UserID != "" {
		return c.Render(http.StatusOK, "SSHRegister.html", map[string]interface{}{
			"LoginInfo": loginInfo,
			"comURL":    comURL,
			"apiInfo":   apiInfo,
		})
	}
	// return c.Redirect(http.StatusPermanentRedirect, "/login")
	return c.Redirect(http.StatusTemporaryRedirect, "/login")
}

// func SSHListForm(c echo.Context) error {
// 	comURL := service.GetCommonURL()
// 	loginInfo := service.CallLoginInfo(c)
// 	apiInfo := service.AuthenticationHandler()
// 	if loginInfo.UserID != "" {
// 		nsList := service.GetRegionList()
// 		fmt.Println("REGION List : ", nsList)

// 		//spew.Dump(nsList)
// 		return c.Render(http.StatusOK, "Resources_Ssh.html", map[string]interface{}{
// 			"LoginInfo": loginInfo,
// 			"comURL":    comURL,
// 			"NSList":    nsList,
// 			"apiInfo":   apiInfo,
// 		})
// 	}

// 	fmt.Println("LoginInfo : ", loginInfo)
// 	return c.Redirect(http.StatusTemporaryRedirect, "/login")

// }

// Controller
// func SpecRegForm(c echo.Context) error {
// 	comURL := service.GetCommonURL()
// 	apiInfo := service.AuthenticationHandler()
// 	if loginInfo := service.CallLoginInfo(c); loginInfo.UserID != "" {
// 		return c.Render(http.StatusOK, "SpecRegister.html", map[string]interface{}{
// 			"LoginInfo": loginInfo,
// 			"comURL":    comURL,
// 			"apiInfo":   apiInfo,
// 		})
// 	}
// 	// return c.Redirect(http.StatusPermanentRedirect, "/login")
// 	return c.Redirect(http.StatusTemporaryRedirect, "/login")
// }

// func SpecListForm(c echo.Context) error {
// 	comURL := service.GetCommonURL()
// 	loginInfo := service.CallLoginInfo(c)
// 	apiInfo := service.AuthenticationHandler()
// 	if loginInfo.UserID != "" {
// 		nsList := service.GetRegionList()
// 		fmt.Println("REGION List : ", nsList)

// 		//spew.Dump(nsList)
// 		return c.Render(http.StatusOK, "Resources_Spec.html", map[string]interface{}{
// 			"LoginInfo": loginInfo,
// 			"comURL":    comURL,
// 			"NSList":    nsList,
// 			"apiInfo":   apiInfo,
// 		})
// 	}

// 	fmt.Println("LoginInfo : ", loginInfo)
// 	return c.Redirect(http.StatusTemporaryRedirect, "/login")

// }

func ChangeCloudBaristaConfig(c echo.Context) error {
	configReq := &tbcommon.TbConfigReq{}
	if err := c.Bind(configReq); err != nil {
		// if err := c.Bind(mCISInfoList); err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "fail",
			"status":  "5001",
		})
	}
	return nil
}
