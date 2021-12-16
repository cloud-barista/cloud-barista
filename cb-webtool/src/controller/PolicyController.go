package controller

import (
	// "encoding/json"
	"fmt"
	"github.com/cloud-barista/cb-webtool/src/util"
	"log"
	"net/http"

	// model "github.com/cloud-barista/cb-webtool/src/model"
	// "github.com/cloud-barista/cb-webtool/src/model/dragonfly"
	// "github.com/cloud-barista/cb-webtool/src/model/tumblebug"

	"github.com/cloud-barista/cb-webtool/src/model/dragonfly"
	service "github.com/cloud-barista/cb-webtool/src/service"

	// util "github.com/cloud-barista/cb-webtool/src/util"

	echotemplate "github.com/foolin/echo-template"
	// echosession "github.com/go-session/echo-session"
	"github.com/labstack/echo"
	// echosession "github.com/go-session/echo-session"
)

// [MCIS] Auto control policy management (WIP) 참조

// PolicyMonitoring 등록화면
func MonitoringConfigPolicyRegForm(c echo.Context) error {
	fmt.Println("PolicyMonitoringRegForm ************ : ")

	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}
	defaultNameSpaceID := loginInfo.DefaultNameSpaceID

	// 최신 namespacelist 가져오기
	nsList, _ := service.GetNameSpaceList()
	log.Println(" nsList  ", nsList)

	//

	return echotemplate.Render(c, http.StatusOK,
		"operation/policies/monitoring/MonitoringPolicyCreate",
		map[string]interface{}{
			"LoginInfo":          loginInfo,
			"DefaultNameSpaceID": defaultNameSpaceID,
			"NameSpaceList":      nsList,
		})
}

// Policy Monitoring 관리 화면
func MonitoringConfigPolicyMngForm(c echo.Context) error {
	fmt.Println("PolicyMonitoringMngForm ************ : ")

	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	defaultNameSpaceID := loginInfo.DefaultNameSpaceID

	// 최신 namespacelist 가져오기
	nsList, _ := service.GetNameSpaceList()
	log.Println(" nsList  ", nsList)

	//MonitoringConfig
	monitoringConfig, _ := service.GetMonitoringConfig()
	// status, filepath, return params
	return echotemplate.Render(c, http.StatusOK,
		"/operation/policies/monitoring/MonitoringConfigPolicyMng", // 파일명
		map[string]interface{}{
			"LoginInfo":          loginInfo,
			"DefaultNameSpaceID": defaultNameSpaceID,
			"NameSpaceList":      nsList,
			"MonitoringConfig":   monitoringConfig,
		})
}

// MonitoringPolicy 목록 조회
func GetMonitoringConfigPolicyList(c echo.Context) error {
	log.Println("GetMonitoringConfigPolicyList : ")
	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	// store := echosession.FromContext(c)
	defaultNameSpaceID := loginInfo.DefaultNameSpaceID

	// monitoringPolicyList, respStatus := service.GetMonitoringPolicyList(defaultNameSpaceID)
	// if respStatus.StatusCode != 200 && respStatus.StatusCode != 201 {
	// 	return c.JSON(respStatus.StatusCode, map[string]interface{}{
	// 		"error":  respStatus.Message,
	// 		"status": respStatus.StatusCode,
	// 	})
	// }

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success",
		// "status":               respStatus.StatusCode,
		"DefaultNameSpaceID": defaultNameSpaceID,
		// "MonitoringPolicyList": monitoringPolicyList,
	})
}

// MonitoringPolicy 등록 처리
func MonitoringConfigPolicyPutProc(c echo.Context) error {
	log.Println("MonitoringConfigPolicyRegProc : ")
	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	monitoringConfigRegInfo := &dragonfly.MonitoringConfigReg{}
	if err := c.Bind(monitoringConfigRegInfo); err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "fail",
			"status":  "fail",
		})
	}
	log.Println(monitoringConfigRegInfo)
	taskKey := loginInfo.DefaultNameSpaceID + "||" + util.TASK_TYPE_MONITORING_POLICY
	resultMonitoringConfigInfo, respStatus := service.PutMonigoringConfig(monitoringConfigRegInfo)
	log.Println("MonitoringPolicyReg service returned")
	if respStatus.StatusCode != 200 && respStatus.StatusCode != 201 {
		service.StoreWebsocketMessage(util.TASK_TYPE_MONITORING_POLICY, taskKey, util.MONITORING_POLICY_STATUS_FAIL, util.TASK_STATUS_REQUEST, c)
		return c.JSON(respStatus.StatusCode, map[string]interface{}{
			"error":  respStatus.Message,
			"status": respStatus.StatusCode,
		})
	}
	service.StoreWebsocketMessage(util.TASK_TYPE_MONITORING_POLICY, taskKey, util.MONITORING_POLICY_STATUS_REG, util.TASK_STATUS_REQUEST, c)
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":          "success",
		"status":           respStatus.StatusCode,
		"MonitoringConfig": resultMonitoringConfigInfo,
	})
}

//PolicyThresholdMngForm
// PolicyThreshold 등록화면
func ThresholdPolicyRegForm(c echo.Context) error {
	fmt.Println("PolicyThresholdRegForm ************ : ")

	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}
	defaultNameSpaceID := loginInfo.DefaultNameSpaceID

	// 최신 namespacelist 가져오기
	nsList, _ := service.GetNameSpaceList()
	log.Println(" nsList  ", nsList)

	//

	return echotemplate.Render(c, http.StatusOK,
		"operation/policies/threshold/ThresholdPolicyCreate",
		map[string]interface{}{
			"LoginInfo":          loginInfo,
			"DefaultNameSpaceID": defaultNameSpaceID,
			"NameSpaceList":      nsList,
		})
}

// Policy Threshold 관리 화면
func MonitoringAlertPolicyMngForm(c echo.Context) error {
	fmt.Println("PolicyThresholdMngForm ************ : ")

	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	defaultNameSpaceID := loginInfo.DefaultNameSpaceID

	// 최신 namespacelist 가져오기
	nsList, _ := service.GetNameSpaceList()
	log.Println(" nsList  ", nsList)

	monitoringAlertPolicyList, _ := service.GetMonitoringAlertList()
	// if respStatus.StatusCode != 200 && respStatus.StatusCode != 201 {
	// 	return c.JSON(respStatus.StatusCode, map[string]interface{}{
	// 		"error":  respStatus.Message,
	// 		"status": respStatus.StatusCode,
	// 	})
	// }

	// Monitoring Alert Event Handler 호출
	monitoringAlertEventHandlerList, _ := service.GetMonitoringAlertEventHandlerList()
	// if respStatus.StatusCode != 200 && respStatus.StatusCode != 201 {
	// 	return c.JSON(respStatus.StatusCode, map[string]interface{}{
	// 		"error":  respStatus.Message,
	// 		"status": respStatus.StatusCode,
	// 	})
	// }

	// status, filepath, return params
	return echotemplate.Render(c, http.StatusOK,
		"/operation/policies/threshold/MonitoringAlertPolicyMng", // 파일명
		map[string]interface{}{
			"LoginInfo":                       loginInfo,
			"DefaultNameSpaceID":              defaultNameSpaceID,
			"NameSpaceList":                   nsList,
			"MonitoringAlertPolicyList":       monitoringAlertPolicyList,
			"MonitoringAlertEventHandlerList": monitoringAlertEventHandlerList,
		})
}

// Monitoring Threshold 목록 조회
func GetMonitoringAlertPolicyList(c echo.Context) error {
	fmt.Println("GetMonitoringAlertPolicyList ************ : ")

	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	// store := echosession.FromContext(c)
	defaultNameSpaceID := loginInfo.DefaultNameSpaceID

	monitoringAlertPolicyList, respStatus := service.GetMonitoringAlertList()
	if respStatus.StatusCode != 200 && respStatus.StatusCode != 201 {
		return c.JSON(respStatus.StatusCode, map[string]interface{}{
			"error":  respStatus.Message,
			"status": respStatus.StatusCode,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":                   "success",
		"status":                    respStatus.StatusCode,
		"DefaultNameSpaceID":        defaultNameSpaceID,
		"MonitoringAlertPolicyList": monitoringAlertPolicyList,
	})
}

// Monitoring Threshold 단건 조회
func GetMonitoringAlertPolicyData(c echo.Context) error {
	fmt.Println("GetMonitoringAlertPolicyData ************ : ")

	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	taskName := c.Param("alertName")
	// store := echosession.FromContext(c)
	defaultNameSpaceID := loginInfo.DefaultNameSpaceID

	monitoringAlertPolicyInfo, respStatus := service.GetMonitoringAlertData(taskName)
	if respStatus.StatusCode != 200 && respStatus.StatusCode != 201 {
		return c.JSON(respStatus.StatusCode, map[string]interface{}{
			"error":  respStatus.Message,
			"status": respStatus.StatusCode,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":                   "success",
		"status":                    respStatus.StatusCode,
		"DefaultNameSpaceID":        defaultNameSpaceID,
		"MonitoringAlertPolicyInfo": monitoringAlertPolicyInfo,
	})
}

// Threshold 등록 처리
func MonitoringAlertPolicyRegProc(c echo.Context) error {
	log.Println("MonitoringAlertPolicyRegProc : ")
	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	// defaultNameSpaceID := loginInfo.DefaultNameSpaceID
	// _, respStatus := service.RegMonitoringPolicy(defaultNameSpaceID, mCISInfo)
	// log.Println("MonitoringPolicyReg service returned")
	// if respStatus.StatusCode != 200 && respStatus.StatusCode != 201 {
	// 	return c.JSON(respStatus.StatusCode, map[string]interface{}{
	// 		"error":  respStatus.Message,
	// 		"status": respStatus.StatusCode,
	// 	})
	// }

	monitoringAlertRegInfo := &dragonfly.VmMonitoringAlertInfo{}
	if err := c.Bind(monitoringAlertRegInfo); err != nil {
		log.Println(err)

		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "fail",
			"status":  "fail",
		})
	}
	log.Println(monitoringAlertRegInfo)
	defaultNameSpaceID := loginInfo.DefaultNameSpaceID
	taskKey := defaultNameSpaceID + "||" + "threshold" + "||" + monitoringAlertRegInfo.AlertName

	resultMonitoringAlertInfo, respStatus := service.RegMonitoringAlert(monitoringAlertRegInfo)
	log.Println("MonitoringAlertPolicyReg service returned")
	if respStatus.StatusCode != 200 && respStatus.StatusCode != 201 {
		service.StoreWebsocketMessage(util.TASK_TYPE_MONITORING_POLICY, taskKey, util.MONITORING_POLICY_STATUS_REG, util.TASK_STATUS_FAIL, c) // session에 작업내용 저장
		return c.JSON(respStatus.StatusCode, map[string]interface{}{
			"error":  respStatus.Message,
			"status": respStatus.StatusCode,
		})
	}

	service.StoreWebsocketMessage(util.TASK_TYPE_MONITORING_POLICY, taskKey, util.MONITORING_POLICY_STATUS_REG, util.TASK_STATUS_REQUEST, c) // session에 작업내용 저장
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":          "success",
		"status":           respStatus.StatusCode,
		"MonitoringConfig": resultMonitoringAlertInfo,
	})

}

// Monitoring Threshold 삭제
func MonitoringAlertPolicyDelProc(c echo.Context) error {
	log.Println("MonitoringAlertPolicyDelProc : ")

	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	// store := echosession.FromContext(c)
	//defaultNameSpaceID := loginInfo.DefaultNameSpaceID

	paramMonitoringAlertID := c.Param("alertName")

	// 글로벌한 설정이라 namespace 없이 호출
	respBody, respStatus := service.DelMonitoringAlert(paramMonitoringAlertID)
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

// Monitoring Alert Event-Handler 목록 조회
func GetMonitoringAlertEventHandlerList(c echo.Context) error {
	fmt.Println("GetMonitoringAlertEventHandlerList ************ : ")

	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	// store := echosession.FromContext(c)
	defaultNameSpaceID := loginInfo.DefaultNameSpaceID

	// Monitoring Alert Event Handler 호출
	monitoringAlertEventHandlerList, respStatus := service.GetMonitoringAlertEventHandlerList()
	if respStatus.StatusCode != 200 && respStatus.StatusCode != 201 {
		return c.JSON(respStatus.StatusCode, map[string]interface{}{
			"error":  respStatus.Message,
			"status": respStatus.StatusCode,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":                         "success",
		"status":                          respStatus.StatusCode,
		"DefaultNameSpaceID":              defaultNameSpaceID,
		"MonitoringAlertEventHandlerList": monitoringAlertEventHandlerList,
	})
}

// Monitoring Alert Event-Handler 등록 처리
func MonitoringAlertEventHandlerRegProc(c echo.Context) error {
	log.Println("MonitoringAlertEventHandlerRegProc : ")
	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	defaultNameSpaceID := loginInfo.DefaultNameSpaceID
	// _, respStatus := service.RegMonitoringPolicy(defaultNameSpaceID, mCISInfo)
	// log.Println("MonitoringPolicyReg service returned")
	// if respStatus.StatusCode != 200 && respStatus.StatusCode != 201 {
	// 	return c.JSON(respStatus.StatusCode, map[string]interface{}{
	// 		"error":  respStatus.Message,
	// 		"status": respStatus.StatusCode,
	// 	})
	// }

	monitoringAlertEventHandlerRegInfo := &dragonfly.VmMonitoringAlertEventHandlerInfoReg{}
	if err := c.Bind(monitoringAlertEventHandlerRegInfo); err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "fail",
			"status":  "fail",
		})
	}
	log.Println(monitoringAlertEventHandlerRegInfo)

	taskKey := defaultNameSpaceID + "||" + "monitoring_threshold_eventHandler" + "||" + monitoringAlertEventHandlerRegInfo.Name
	resultMonitoringAlertEventHandlerInfo, respStatus := service.RegMonitoringAlertEventHandler(monitoringAlertEventHandlerRegInfo)
	log.Println("MonitoringAlertEventHandlerReg service returned")
	if respStatus.StatusCode != 200 && respStatus.StatusCode != 201 {
		service.StoreWebsocketMessage(util.TASK_TYPE_MONITORINGTHRESHOLD_EVENTHANDLER, taskKey, util.MONITORING_THRESHOLD_REG, util.TASK_STATUS_FAIL, c)
		return c.JSON(respStatus.StatusCode, map[string]interface{}{
			"error":  respStatus.Message,
			"status": respStatus.StatusCode,
		})
	}
	service.StoreWebsocketMessage(util.TASK_TYPE_MONITORINGTHRESHOLD_EVENTHANDLER, taskKey, util.MONITORING_THRESHOLD_REG, util.TASK_STATUS_REQUEST, c)
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":                "success",
		"status":                 respStatus.StatusCode,
		"MonitoringEventHandler": resultMonitoringAlertEventHandlerInfo,
	})

}

// Monitoring Alert Event-Handler 삭제
func MonitoringAlertEventHandlerDelProc(c echo.Context) error {
	log.Println("MonitoringAlertEventHandlerDelProc : ")

	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	// store := echosession.FromContext(c)
	//defaultNameSpaceID := loginInfo.DefaultNameSpaceID

	paramMonitoringAlertEvantHandlerType := c.Param("type")
	paramMonitoringAlertEvantHandlerName := c.Param("name")

	// 글로벌한 설정이라 namespace 없이 호출
	respBody, respStatus := service.DelMonitoringAlertEventHandler(paramMonitoringAlertEvantHandlerType, paramMonitoringAlertEvantHandlerName)
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

// PolicyPlacement 등록화면
func PlacementPolicyRegForm(c echo.Context) error {
	fmt.Println("PolicyPlacementRegForm ************ : ")

	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}
	defaultNameSpaceID := loginInfo.DefaultNameSpaceID

	// 최신 namespacelist 가져오기
	nsList, _ := service.GetNameSpaceList()
	log.Println(" nsList  ", nsList)

	//

	return echotemplate.Render(c, http.StatusOK,
		"operation/policies/placement/PlacementPolicyCreate",
		map[string]interface{}{
			"LoginInfo":          loginInfo,
			"DefaultNameSpaceID": defaultNameSpaceID,
			"NameSpaceList":      nsList,
		})
}

// Policy Monitoring 관리 화면
func PlacementPolicyMngForm(c echo.Context) error {
	fmt.Println("PolicyPlacementMngForm ************ : ")

	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	defaultNameSpaceID := loginInfo.DefaultNameSpaceID

	// 최신 namespacelist 가져오기
	nsList, _ := service.GetNameSpaceList()
	log.Println(" nsList  ", nsList)

	// status, filepath, return params
	return echotemplate.Render(c, http.StatusOK,
		"/operation/policies/placement/PlacementPolicyMng", // 파일명
		map[string]interface{}{
			"LoginInfo":          loginInfo,
			"DefaultNameSpaceID": defaultNameSpaceID,
			"NameSpaceList":      nsList,
		})
}

// Placement Policy 목록 조회
func GetPlacementPolicyList(c echo.Context) error {
	log.Println("GetMonitoringPolicyList : ")
	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	// store := echosession.FromContext(c)
	defaultNameSpaceID := loginInfo.DefaultNameSpaceID

	// monitoringPolicyList, respStatus := service.GetMonitoringPolicyList(defaultNameSpaceID)
	// if respStatus.StatusCode != 200 && respStatus.StatusCode != 201 {
	// 	return c.JSON(respStatus.StatusCode, map[string]interface{}{
	// 		"error":  respStatus.Message,
	// 		"status": respStatus.StatusCode,
	// 	})
	// }

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success",
		// "status":               respStatus.StatusCode,
		"DefaultNameSpaceID": defaultNameSpaceID,
		// "MonitoringPolicyList": monitoringPolicyList,
	})
}

// Placement 등록 처리
func PlacementPolicyRegProc(c echo.Context) error {
	log.Println("PlacementolicyRegProc : ")
	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	// defaultNameSpaceID := loginInfo.DefaultNameSpaceID
	// _, respStatus := service.RegMonitoringPolicy(defaultNameSpaceID, mCISInfo)
	// log.Println("MonitoringPolicyReg service returned")
	// if respStatus.StatusCode != 200 && respStatus.StatusCode != 201 {
	// 	return c.JSON(respStatus.StatusCode, map[string]interface{}{
	// 		"error":  respStatus.Message,
	// 		"status": respStatus.StatusCode,
	// 	})
	// }

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success",
		// "status":  respStatus.StatusCode,
	})
}

// MonitoringPolicy 목록 조회
func GetMonitoringAlertLogList(c echo.Context) error {
	log.Println("GetMonitoringAlertLogList : ")
	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	// store := echosession.FromContext(c)
	defaultNameSpaceID := loginInfo.DefaultNameSpaceID

	paramTaskName := c.Param("task_name")
	paramLevel := c.Param("level")

	monitoringAlertLogList, respStatus := service.GetMonitoringAlertLogList(paramTaskName, paramLevel)
	if respStatus.StatusCode != 200 && respStatus.StatusCode != 201 {
		return c.JSON(respStatus.StatusCode, map[string]interface{}{
			"error":  respStatus.Message,
			"status": respStatus.StatusCode,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success",
		// "status":               respStatus.StatusCode,
		"DefaultNameSpaceID":     defaultNameSpaceID,
		"MonitoringAlertLogList": monitoringAlertLogList,
	})
}
