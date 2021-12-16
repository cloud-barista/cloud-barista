package controller

import (
	// "encoding/json"
	"fmt"

	// "github.com/cloud-barista/cb-webtool/src/model/tumblebug"
	tbcommon "github.com/cloud-barista/cb-webtool/src/model/tumblebug/common"
	// tbmcir "github.com/cloud-barista/cb-webtool/src/model/tumblebug/mcir"
	// tbmcis "github.com/cloud-barista/cb-webtool/src/model/tumblebug/mcis"

	service "github.com/cloud-barista/cb-webtool/src/service"

	// util "github.com/cloud-barista/cb-webtool/src/util"

	"log"
	"net/http"

	echotemplate "github.com/foolin/echo-template"
	echosession "github.com/go-session/echo-session"
	"github.com/labstack/echo"
)

// deprecated
// func NsRegController(c echo.Context) error {
// 	username := c.FormValue("username")
// 	description := c.FormValue("description")

// 	fmt.Println("NSRegController : ", username, description)
// 	return nil
// }

// func NsRegForm(c echo.Context) error {
func NameSpaceRegForm(c echo.Context) error {
	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		// Login 정보가 없으므로 login화면으로
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	// 	comURL := service.GetCommonURL()
	// 	apiInfo := service.AuthenticationHandler()
	// 	if loginInfo := service.CallLoginInfo(c); loginInfo.Username != "" {
	// return c.Render(http.StatusOK, "NSRegister.html", map[string]interface{}{
	return echotemplate.Render(c, http.StatusOK, "setting/namespaces/NSRegister.html", map[string]interface{}{
		"LoginInfo": loginInfo,
		// 			"comURL":    comURL,
		// 			"apiInfo":   apiInfo,
	})
	// 	}
	// 	return c.Redirect(http.StatusTemporaryRedirect, "/login")
}

// namespace 등록 처리
func NameSpaceRegProc(c echo.Context) error {

	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		// Login 정보가 없으므로 login화면으로
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	// namespace := c.FormValue("name")
	// description := c.FormValue("description")
	// fmt.Println("namespace : " + namespace + " , description :" + description)
	// nameSpaceInfo := new(model.NameSpaceInfo)
	// nameSpaceInfo.Name = namespace
	// nameSpaceInfo.Description = description

	nameSpaceInfo := new(tbcommon.TbNsInfo)
	if err := c.Bind(nameSpaceInfo); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	// if err = c.Validate(nameSpaceInfo); err != nil {
	// 	return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	// }
	fmt.Println("nameSpaceInfo : ", nameSpaceInfo)

	// Tubblebug 호출하여 namespace 생성

	// person := Person{"Alex", 10}
	// pbytes, _ := json.Marshal(person)
	respBody, respStatus := service.RegNameSpace(nameSpaceInfo)
	fmt.Println("=============respBody =============", respBody)
	if respStatus.StatusCode != 200 && respStatus.StatusCode != 201 {
		return c.JSON(respStatus.StatusCode, map[string]interface{}{
			"error":  respStatus.Message,
			"status": respStatus.StatusCode,
		})
	}

	// 저장 성공하면 namespace 목록 조회
	nameSpaceList, nsStatus := service.GetNameSpaceList()
	if nsStatus.StatusCode == 500 {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message":       respStatus.Message,
			"status":        respStatus.StatusCode,
			"NameSpaceList": nil,
		})
	}
	storeNameSpaceErr := service.SetStoreNameSpaceList(c, nameSpaceList)
	if storeNameSpaceErr != nil {
		log.Println("Store NameSpace Err")
	}
	// return namespace 목록
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":       "success",
		"status":        "200",
		"NameSpaceList": nameSpaceList,
	})
}

// Namespace 수정
func NameSpaceUpdateProc(c echo.Context) error {
	log.Println("NameSpaceUpdateProc : ")

	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		// Login 정보가 없으므로 login화면으로
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}
	defaultNameSpaceID := loginInfo.DefaultNameSpaceID

	nameSpaceInfo := new(tbcommon.TbNsInfo)
	if err := c.Bind(nameSpaceInfo); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	respBody, respStatus := service.UpdateNameSpace(defaultNameSpaceID, nameSpaceInfo)
	fmt.Println("=============respBody =============", respBody)
	if respStatus.StatusCode != 200 && respStatus.StatusCode != 201 {
		return c.JSON(respStatus.StatusCode, map[string]interface{}{
			"error":  respStatus.Message,
			"status": respStatus.StatusCode,
		})
	}

	// 저장 성공하면 namespace 목록 조회
	nameSpaceList, nsStatus := service.GetNameSpaceList()
	if nsStatus.StatusCode == 500 {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message":       nsStatus.Message,
			"status":        nsStatus.StatusCode,
			"NameSpaceList": nil,
		})
	}

	storeNameSpaceErr := service.SetStoreNameSpaceList(c, nameSpaceList)
	if storeNameSpaceErr != nil {
		log.Println("Store NameSpace Err")
	}

	// return namespace 목록
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":       "success",
		"status":        respStatus.StatusCode,
		"NameSpaceList": nameSpaceList,
	})

	// return c.JSON(http.StatusOK, map[string]interface{}{
	// 	"message": "success",
	// 	"status":  "200",
	// })
}

// NameSpace 삭제
func NameSpaceDelProc(c echo.Context) error {
	log.Println("NameSpaceDelProc : ")

	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		// Login 정보가 없으므로 login화면으로
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	paramNameSpaceID := c.Param("nameSpaceID")
	defaultNameSpaceID := loginInfo.DefaultNameSpaceID
	log.Println(paramNameSpaceID)
	if paramNameSpaceID == defaultNameSpaceID {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":  "Cannot delete current namespace",
			"status": 4001,
		})
	}

	respBody, respStatus := service.DelNameSpace(paramNameSpaceID)
	fmt.Println("=============respBody =============", respBody)
	if respStatus.StatusCode != 200 && respStatus.StatusCode != 201 {
		return c.JSON(respStatus.StatusCode, map[string]interface{}{
			"error":  respStatus.Message,
			"status": respStatus.StatusCode,
		})
	}

	// 저장 성공하면 namespace 목록 조회
	nameSpaceList, nsStatus := service.GetNameSpaceList()
	if nsStatus.StatusCode == 500 {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message":       nsStatus.Message,
			"status":        nsStatus.StatusCode,
			"NameSpaceList": nil,
		})
	}

	storeNameSpaceErr := service.SetStoreNameSpaceList(c, nameSpaceList)
	if storeNameSpaceErr != nil {
		log.Println("Store NameSpace Err")
	}

	// return namespace 목록
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":       "success",
		"status":        respStatus,
		"NameSpaceList": nameSpaceList,
	})

	// return c.JSON(http.StatusOK, map[string]interface{}{
	// 	"message": "success",
	// 	"status":  "200",
	// })
}

// NsListForm -> NameSpaceMngForm으로 변경
//func NsListForm(c echo.Context) error {
func NameSpaceMngForm(c echo.Context) error {
	fmt.Println("=============start NameSpaceMngForm =============")
	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		// Login 정보가 없으므로 login화면으로
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	fmt.Println("=============start GetNSList =============")

	nameSpaceList, _ := service.GetNameSpaceList()
	fmt.Println("=============start GetNSList =============", nameSpaceList)

	storeNameSpaceErr := service.SetStoreNameSpaceList(c, nameSpaceList)
	if storeNameSpaceErr != nil {
		log.Println("Store NameSpace Err")
	}

	// status, filepath(vies기준), return params
	return echotemplate.Render(c, http.StatusOK,
		"setting/namespaces/NameSpaceMng", // 파일명
		map[string]interface{}{
			"LoginInfo":     loginInfo,
			"NameSpaceList": nameSpaceList,
		})

}

// 사용자의 namespace 목록 조회
func GetNameSpaceList(c echo.Context) error {
	fmt.Println("====== GET NAMESPACE LIST ========")
	// store := echosession.FromContext(c)
	// nameSpaceInfoList, nsStatus := service.GetNameSpaceList()

	optionParam := c.QueryParam("option")

	if optionParam == "id" {
		nameSpaceInfoList, nsStatus := service.GetNameSpaceListByOptionID(optionParam)
		if nsStatus.StatusCode == 500 {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"message": nsStatus.Message,
				"status":  nsStatus.StatusCode,
			})
		}
		return c.JSON(http.StatusOK, nameSpaceInfoList)
	} else {
		nameSpaceInfoList, nsStatus := service.GetNameSpaceListByOption(optionParam)
		if nsStatus.StatusCode == 500 {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"message": nsStatus.Message,
				"status":  nsStatus.StatusCode,
			})
		}
		return c.JSON(http.StatusOK, nameSpaceInfoList)
	}
	//nameSpaceInfoList, nsStatus := service.GetStoredNameSpaceList(c)

}

// 기본 namespace set. set default Namespace
func SetNameSpace(c echo.Context) error {
	fmt.Println("====== SET SELECTED NAME SPACE ========")
	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		return c.JSON(http.StatusNotAcceptable, map[string]string{
			"message": "Not Login",
			"status":  "403",
		})
	}
	store := echosession.FromContext(c)

	result, ok := store.Get(loginInfo.UserID)
	if !ok {
		return c.JSON(http.StatusNotAcceptable, map[string]string{
			"message": "Not Login",
			"status":  "403",
		})
	}
	storedUser := result.(map[string]string)

	nameSpaceID := c.Param("nameSpaceID")
	loginInfo.DefaultNameSpaceID = nameSpaceID

	// nsResult, nsOk := store.Get("namespacelist")
	// fmt.Println("nsResult : ", nsResult)
	// if !nsOk {
	// 	fmt.Println("nsOk : ", nsOk)
	// 	return c.JSON(http.StatusOK, map[string]interface{}{
	// 		"message": "저장된 namespace가 없습니다.",
	// 		"status":  "403",
	// 	})
	nsList, nsStatus := service.GetStoredNameSpaceList(c)
	if nsStatus.StatusCode != 200 && nsStatus.StatusCode != 201 {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "NameSpace 조회 오류.",
			"status":  "403",
		})
	} else {
		fmt.Println("______________")

		// nsList := nsResult.([]tbcommon.TbNsInfo)
		// fmt.Println("nsList ", nsList)
		for _, nsInfo := range nsList {
			fmt.Println(nsInfo.ID + " :  " + nameSpaceID)
			if nsInfo.ID == nameSpaceID {
				loginInfo.DefaultNameSpaceID = nsInfo.ID
				loginInfo.DefaultNameSpaceName = nsInfo.Name
				fmt.Println(nsInfo.ID + " :  " + nameSpaceID + " found " + nsInfo.Name)
				storedUser["defaultnameSpacename"] = nsInfo.ID
				storedUser["defaultnamespaceid"] = nsInfo.Name
				break
			}
		}
	}

	// storedUser["defaultnamespaceid"] = nameSpaceID
	fmt.Println("storedUser : ", storedUser)
	store.Set(loginInfo.UserID, storedUser)

	storeErr := store.Save()
	if storeErr != nil {
		return c.JSON(http.StatusNotAcceptable, map[string]string{
			"message": storeErr.Error(),
			"status":  "403",
		})
	}

	mcisList, _ := service.GetMcisListByID(loginInfo.DefaultNameSpaceID)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":   "success",
		"status":    "200",
		"LoginInfo": loginInfo,
		"McisList":  mcisList,
	})

}

// 기본 namespace get. get default Namespace
func GetNameSpace(c echo.Context) error {
	fmt.Println("====== GET SELECTED NAME SPACE ========")
	// store := echosession.FromContext(c)

	// getInfo, ok := store.Get("namespace")
	// if !ok {
	// 	return c.JSON(http.StatusNotAcceptable, map[string]string{
	// 		"message": "Not Exist",
	// 	})
	// }
	// nsId := getInfo.(string)

	loginInfo := service.CallLoginInfo(c)

	res := map[string]string{
		"message": "success",
		"nsID":    loginInfo.DefaultNameSpaceID,
	}

	return c.JSON(http.StatusOK, res)
}
