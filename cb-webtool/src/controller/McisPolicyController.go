package controller

import (
	"fmt"
	"log"
	"net/http"

	// "github.com/cloud-barista/cb-webtool/src/model/tumblebug"
	// tbcommon "github.com/cloud-barista/cb-webtool/src/model/tumblebug/common"
	// tbmcir "github.com/cloud-barista/cb-webtool/src/model/tumblebug/mcir"
	tbmcis "github.com/cloud-barista/cb-webtool/src/model/tumblebug/mcis"

	service "github.com/cloud-barista/cb-webtool/src/service"
	"github.com/labstack/echo"

	echotemplate "github.com/foolin/echo-template"
	echosession "github.com/go-session/echo-session"
)

///// McisPolicyController ////////////

// McisPolishMngForm 화면
func McisPolicyMngForm(c echo.Context) error {
	fmt.Println("McisPolishMngForm ************ : ")

	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	defaultNameSpaceID := loginInfo.DefaultNameSpaceID

	store := echosession.FromContext(c)

	// 최신 namespacelist 가져오기
	nsList, _ := service.GetNameSpaceList()
	store.Set("namespace", nsList)
	store.Save()
	log.Println(" nsList  ", nsList)

	// 해당 Namespace의 모든 MCIS 조회
	mcisPolicyList, _ := service.GetMcisPolicyList(defaultNameSpaceID)
	log.Println(" mcisList  ", mcisPolicyList)

	return echotemplate.Render(c, http.StatusOK,
		//"operation/monitorings/McisPolishMng", // 파일명
		"operation/monitorings/mcismonitoring/McisMonitoringMng", // 파일명
		map[string]interface{}{
			"LoginInfo":          loginInfo,
			"DefaultNameSpaceID": defaultNameSpaceID,
			"NameSpaceList":      nsList,
			"McisPolicyList":     mcisPolicyList,
		})

}

func GetMcisPolicyInfoList(c echo.Context) error {
	log.Println("GetMcisPolicyInfoList")
	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/login") // 조회기능에서 바로 login화면으로 돌리지말고 return message로 하는게 낫지 않을까?
	}

	defaultNameSpaceID := loginInfo.DefaultNameSpaceID

	mcisID := c.Param("mcisID")
	log.Println("mcisID= " + mcisID)

	resultMcisPolicyInfo, _ := service.GetMcisPolicyList(defaultNameSpaceID)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":        "success",
		"status":         200,
		"McisPolicyInfo": resultMcisPolicyInfo,
	})
}

// GetMcisPolishInfoData
// 특정 MCIS의 Polish 상세정보를 가져온다.
func GetMcisPolicyInfoData(c echo.Context) error {
	log.Println("GetMcisPolishInfoData")
	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/login") // 조회기능에서 바로 login화면으로 돌리지말고 return message로 하는게 낫지 않을까?
	}

	defaultNameSpaceID := loginInfo.DefaultNameSpaceID

	mcisID := c.Param("mcisID")
	log.Println("mcisID= " + mcisID)

	resultMcisPolicyInfo, _ := service.GetMcisPolicyInfoData(defaultNameSpaceID, mcisID)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":        "success",
		"status":         200,
		"McisPolicyInfo": resultMcisPolicyInfo,
	})
}

func McisPolicyRegProc(c echo.Context) error {
	log.Println("McisPolicyRegProc : ")
	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	mcisPolicyInfo := &tbmcis.McisPolicyInfo{}
	if err := c.Bind(mcisPolicyInfo); err != nil {
		// if err := c.Bind(mCISInfoList); err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "fail",
			"status":  "fail",
		})
	}
	log.Println(mcisPolicyInfo) // 여러개일 수 있음.

	defaultNameSpaceID := loginInfo.DefaultNameSpaceID
	mcisID := c.Param("mcisID")

	// TODO : defaultNameSpaceID 가 없으면 설정화면으로 보낼 것
	resultMcisPolicyInfo, respStatus := service.RegMcisPolicy(defaultNameSpaceID, mcisID, mcisPolicyInfo)
	log.Println("RegMcisPolicy service returned")
	if respStatus.StatusCode != 200 && respStatus.StatusCode != 201 {
		return c.JSON(respStatus.StatusCode, map[string]interface{}{
			"error":          respStatus.Message,
			"status":         respStatus.StatusCode,
			"McisPolicyInfo": resultMcisPolicyInfo,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success",
		"status":  respStatus.StatusCode,
	})
}

// 삭제
func McisPolicyDelProc(c echo.Context) error {
	log.Println("McisPolishDelProc : ")

	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	// store := echosession.FromContext(c)
	defaultNameSpaceID := loginInfo.DefaultNameSpaceID

	paramMcisID := c.Param("mcisID")

	respBody, respStatus := service.DelMcisPolicy(defaultNameSpaceID, paramMcisID)
	fmt.Println("=============respBody =============", respBody)
	// error 났을 때만 Message가 set 되고 정상인 경우에는 success로 return.

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
