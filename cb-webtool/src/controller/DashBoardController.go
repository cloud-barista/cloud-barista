package controller

import (
	// "encoding/json"
	"fmt"
	"log"
	"net/http"

	// "sync"

	// model "github.com/cloud-barista/cb-webtool/src/model"
	// "github.com/cloud-barista/cb-webtool/src/model/tumblebug"
	// tbcommon "github.com/cloud-barista/cb-webtool/src/model/tumblebug/common"
	// tbmcir "github.com/cloud-barista/cb-webtool/src/model/tumblebug/mcir"
	// tbmcis "github.com/cloud-barista/cb-webtool/src/model/tumblebug/mcis"

	webtool "github.com/cloud-barista/cb-webtool/src/model/webtool"

	"github.com/cloud-barista/cb-webtool/src/service"
	util "github.com/cloud-barista/cb-webtool/src/util"
	"github.com/labstack/echo"

	echotemplate "github.com/foolin/echo-template"
	echosession "github.com/go-session/echo-session"
)

type RespPublicIPInfo struct {
	PublicIPInfo []struct {
		PublicIp string `json:"publicIP"`
		Status   string `json:"status"`
		VmID     string `json:"id"`
		VmName   string `json:"name"`
	} `json:"vm"`
}

// 특정 Namespace의 Dashboard -- > 모든 Namespace의 Dashboard도 있음.
func DashBoardByNameSpaceMngForm(c echo.Context) error {
	fmt.Println("DashBoardByNameSpaceMngForm ************ : ")

	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	defaultNameSpaceID := loginInfo.DefaultNameSpaceID

	store := echosession.FromContext(c)

	//totalMcisStatusCountMap := make(map[string]int)             // 모든 MCIS의 상태 Map
	//mcisStatusCountMapByMcis := make(map[string]map[string]int) // MCIS ID별 mcis status
	//totalVmStatusCountMap := make(map[string]int)               // 모든 VM의 상태 Map
	//vmStatusCountMapByMcis := make(map[string]map[string]int)   // MCIS ID 별 vmStatusMap [{mcis+status, count},{mcis+status, count}...]
	//mcisSimpleInfoList := []webtool.McisSimpleInfo{}            // mics summary 정보
	//
	//totalMcisCount := 0 // mcis 갯수
	//totalVmCount := 0   // 모든 vm 갯수
	//providerCount := 0
	//totalConnectionCount := 0
	//connectionConfigCountMap := make(map[string]int)

	// 최신 namespacelist 가져오기
	nsList, _ := service.GetNameSpaceList()
	store.Set("namespace", nsList)
	_ = store.Save()
	log.Println(" nsList  ", nsList)

	//// 해당 Namespace의 모든 MCIS 조회 -> UI Onload 시 가져오는 것으로 변경
	//optionParam := c.QueryParam("option")
	//mcisList, mcisErr := service.GetMcisList(defaultNameSpaceID, optionParam)
	//log.Println(" mcisList  ", mcisList)
	//
	//if mcisErr.StatusCode != 200 && mcisErr.StatusCode != 201 {
	//	return echotemplate.Render(c, http.StatusOK,
	//		"operation/dashboards/DashboardByNameSpaceMng", // 파일명
	//		map[string]interface{}{
	//			"LoginInfo":             loginInfo,
	//			"DefaultNameSpaceID":    defaultNameSpaceID,
	//			"NameSpaceList":         nsList,
	//			"TotalVmCount":          totalVmCount,
	//			"TotalVmStatusCountMap": totalVmStatusCountMap, // 모든 VmStatus 별 count Map(MCIS 무관)
	//
	//			// cp count 영역
	//			"TotalProviderCount":         providerCount,            // VM이 등록 된 provider 목록
	//			"TotalConnectionConfigCount": totalConnectionCount,     // 총 connection 갯수
	//			"ConnectionConfigCountMap":   connectionConfigCountMap, // provider별 connection 수
	//
	//			// mcis count 영역
	//			"TotalMcisCount":          totalMcisCount,
	//			"TotalMcisStatusCountMap": totalMcisStatusCountMap, // 모든 MCIS의 상태 Map
	//
	//			// mcis list
	//			"McisList":               mcisSimpleInfoList,     // 표에 뿌려줄 mics summary 정보
	//			"VmStatusCountMapByMcis": vmStatusCountMapByMcis, // MCIS ID 별 vmStatusMap
	//		})
	//}
	//
	//totalMcisCount = len(mcisList) // mcis 갯수
	//
	//// 등록된 mcis가 없으면 mcis생성화면으로 이동한다.
	//if len(mcisList) == 0 {
	//	return c.Redirect(http.StatusTemporaryRedirect, "/operation/manages/mcismng/regform")
	//}
	//
	//for _, mcisInfo := range mcisList {
	//	resultMcisStatusCountMap := service.GetMcisStatusCountMap(mcisInfo)
	//
	//	for mcisStatusKey, mcisStatusCountVal := range resultMcisStatusCountMap {
	//		if mcisStatusKey == "TOTAL" { // Total까지 오므로 Total은 제외
	//			continue
	//		}
	//
	//		val, exists := totalMcisStatusCountMap[mcisStatusKey]
	//		if exists {
	//			totalMcisStatusCountMap[mcisStatusKey] = val + mcisStatusCountVal
	//		} else {
	//			totalMcisStatusCountMap[mcisStatusKey] = mcisStatusCountVal
	//		}
	//	}
	//
	//	mcisStatusCountMapByMcis[mcisInfo.ID] = resultMcisStatusCountMap // 각 MCIS의 status별 cnt
	//
	//	//////////// vm status area
	//	resultSimpleVmList, resultVmStatusCountMap := service.GetSimpleVmWithStatusCountMap(mcisInfo)
	//
	//	resultVmStatusNames := ""
	//	for _, vmSimpleObj := range resultSimpleVmList {
	//		resultVmStatusNames += vmSimpleObj.VmID + "|" + vmSimpleObj.VmStatus + "@"
	//	}
	//
	//	log.Println("before " + resultVmStatusNames)
	//	if len(resultVmStatusNames) > 0 {
	//		resultVmStatusNames = resultVmStatusNames[:len(resultVmStatusNames)-1]
	//	}
	//	log.Println("after " + resultVmStatusNames)
	//
	//	// UI에서 보여 줄 STATUS로 Count. (가져온 Key중에 UI에서 보여줄 Key가 없을 수 있으므로)
	//	for i, _ := range util.STATUS_ARRAY {
	//		// status_array는 고정값이므로 없는 경우 default로 '0'으로 set
	//		_, exists := resultVmStatusCountMap[util.STATUS_ARRAY[i]]
	//		if !exists {
	//			resultVmStatusCountMap[util.STATUS_ARRAY[i]] = 0
	//		}
	//		totalVmStatusCountMap[util.STATUS_ARRAY[i]] += resultVmStatusCountMap[util.STATUS_ARRAY[i]]
	//	}
	//	// UI manage mcis > server 영역에서는 run/stopped/terminated 만 있음. etc를 stopped에 추가한다.
	//	totalVmStatusCountMap[util.VM_STATUS_STOPPED] = totalVmStatusCountMap[util.VM_STATUS_STOPPED] + resultVmStatusCountMap[util.VM_STATUS_ETC]
	//
	//	totalVmCount += resultVmStatusCountMap["TOTAL"] // 모든 vm의 갯수
	//
	//	totalVmCountByMcis := resultVmStatusCountMap["TOTAL"]        // 모든 vm의 갯수
	//	vmStatusCountMapByMcis[mcisInfo.ID] = resultVmStatusCountMap // MCIS 내 vm 상태별 cnt
	//
	//	// Provider 별 connection count (Location 내에 있는 provider로 갯수 셀 것.)
	//	mcisConnectionMap := service.GetVmConnectionCountByMcis(mcisInfo) // 해당 MCIS의 각 provider별 connection count
	//	log.Println(mcisConnectionMap)
	//
	//	mcisConnectionNames := ""
	//	for connectKey, _ := range mcisConnectionMap {
	//		mcisConnectionNames += connectKey + " "
	//	}
	//	////////////// return value 에 set
	//	mcisSimpleInfo := webtool.McisSimpleInfo{}
	//	mcisSimpleInfo.ID = mcisInfo.ID
	//	mcisSimpleInfo.Status = mcisInfo.Status
	//	mcisSimpleInfo.McisStatus = util.GetMcisStatus(mcisInfo.Status)
	//	mcisSimpleInfo.Name = mcisInfo.Name
	//	mcisSimpleInfo.Description = mcisInfo.Description
	//
	//	mcisSimpleInfo.InstallMonAgent = mcisInfo.InstallMonAgent
	//	mcisSimpleInfo.Label = mcisInfo.Label
	//
	//	mcisSimpleInfo.VmCount = totalVmCountByMcis // 해당 mcis의 모든 vm 갯수
	//
	//	mcisSimpleInfo.VmSimpleList = resultSimpleVmList
	//	mcisSimpleInfo.VmStatusNames = resultVmStatusNames
	//	mcisSimpleInfo.VmStatusCountMap = resultVmStatusCountMap
	//
	//	mcisSimpleInfo.ConnectionConfigProviderMap = mcisConnectionMap     // 해당 MCIS 등록된 connection의 provider 목록
	//	mcisSimpleInfo.ConnectionConfigProviderNames = mcisConnectionNames // 해당 MCIS 등록된 connection의 provider 목록을 String
	//	mcisSimpleInfo.ConnectionConfigProviderCount = len(mcisConnectionMap)
	//
	//	mcisSimpleInfoList = append(mcisSimpleInfoList, mcisSimpleInfo)
	//}
	//
	//// provider 별 연결정보 count(MCIS 무관)
	//cloudConnectionConfigInfoList, _ := service.GetCloudConnectionConfigList()
	//connectionConfigCountMap, providerCount = service.GetCloudConnectionCountMap(cloudConnectionConfigInfoList)
	//totalConnectionCount = len(cloudConnectionConfigInfoList)

	return echotemplate.Render(c, http.StatusOK,
		"operation/dashboards/DashboardByNameSpaceMng", // 파일명
		map[string]interface{}{
			"LoginInfo":          loginInfo,
			"DefaultNameSpaceID": defaultNameSpaceID,
			"NameSpaceList":      nsList,
			// "McisList":           mcisList,

			//// server count 영역
			//"TotalVmCount":          totalVmCount,
			//"TotalVmStatusCountMap": totalVmStatusCountMap, // 모든 VmStatus 별 count Map(MCIS 무관)
			//
			//// cp count 영역
			//"TotalProviderCount":         providerCount,            // VM이 등록 된 provider 목록
			//"TotalConnectionConfigCount": totalConnectionCount,     // 총 connection 갯수
			//"ConnectionConfigCountMap":   connectionConfigCountMap, // provider별 connection 수
			//
			//// mcis count 영역
			//"TotalMcisCount":          totalMcisCount,
			//"TotalMcisStatusCountMap": totalMcisStatusCountMap, // 모든 MCIS의 상태 Map
			//
			//// mcis list
			//"McisList":               mcisSimpleInfoList,     // 표에 뿌려줄 mics summary 정보
			//"VmStatusCountMapByMcis": vmStatusCountMapByMcis, // MCIS ID 별 vmStatusMap
		})

}

func GlobalDashBoardMngForm(c echo.Context) error {
	fmt.Println("GlobalDashBoardMngForm ************ : ")

	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	defaultNameSpaceID := loginInfo.DefaultNameSpaceID

	store := echosession.FromContext(c)

	// 최신 namespacelist 가져오기
	nsList, _ := service.GetNameSpaceList()
	store.Set("namespace", nsList)
	_ = store.Save()
	log.Println(" nsList  ", nsList)

	// provider 별 연결정보 count(MCIS 무관)
	cloudConnectionConfigInfoList, _ := service.GetCloudConnectionConfigList()
	connectionConfigCountMap, providerCount := service.GetCloudConnectionCountMap(cloudConnectionConfigInfoList)
	totalConnectionCount := len(cloudConnectionConfigInfoList)

	// 모든 MCIS 조회
	// mcisList, _ := service.GetMcisList(defaultNameSpaceID)
	optionParam := c.QueryParam("option")
	mcisList, _ := service.GetMcisList(defaultNameSpaceID, optionParam)
	log.Println(" mcisList  ", mcisList)

	// totalMcisCount := len(mcisList) // mcis 갯수
	totalVmCount := 0 // 모든 vm 갯수

	totalMcisStatusCountMap := make(map[string]int)             // 모든 MCIS의 상태 Map
	mcisStatusCountMapByMcis := make(map[string]map[string]int) // MCIS ID별 mcis status
	totalVmStatusCountMap := make(map[string]int)               // 모든 VM의 상태 Map
	vmStatusCountMapByMcis := make(map[string]map[string]int)   // MCIS ID 별 vmStatusMap [{mcis+status, count},{mcis+status, count}...]
	mcisSimpleInfoList := []webtool.McisSimpleInfo{}            // mics summary 정보

	for _, mcisInfo := range mcisList {
		resultMcisStatusCountMap := service.GetMcisStatusCountMap(mcisInfo)
		// mcisStatusMap["RUNNING"] = mcisStatusRunning
		// mcisStatusMap["STOPPED"] = mcisStatusStopped
		// mcisStatusMap["TERMINATED"] = mcisStatusTerminated
		// mcisStatusMap["TOTAL"] = mcisStatusRunning + mcisStatusStop

		for mcisStatusKey, mcisStatusCountVal := range resultMcisStatusCountMap {
			if mcisStatusKey == "TOTAL" { // Total까지 오므로 Total은 제외
				continue
			}

			val, exists := totalMcisStatusCountMap[mcisStatusKey]
			if exists {
				totalMcisStatusCountMap[mcisStatusKey] = val + mcisStatusCountVal
			} else {
				totalMcisStatusCountMap[mcisStatusKey] = mcisStatusCountVal
			}
		}

		mcisStatusCountMapByMcis[mcisInfo.ID] = resultMcisStatusCountMap // 각 MCIS의 status별 cnt
		// connectionConfigCountMap[util.GetProviderName(connectionInfo.ProviderName)] = count

		//////////// vm status area
		resultVmSimpleList, resultVmStatusCountMap := service.GetSimpleVmWithStatusCountMap(mcisInfo)

		resultVmStatusNames := ""
		for _, vmSimpleObj := range resultVmSimpleList {
			resultVmStatusNames += vmSimpleObj.VmID + "|" + vmSimpleObj.VmName + "@"
		}

		log.Println("before " + resultVmStatusNames)
		if len(resultVmStatusNames) > 0 {
			resultVmStatusNames = resultVmStatusNames[:len(resultVmStatusNames)-1]
		}
		log.Println("after " + resultVmStatusNames)

		// UI에서 보여 줄 STATUS로 Count. (가져온 Key중에 UI에서 보여줄 Key가 없을 수 있으므로)
		for i, _ := range util.STATUS_ARRAY {
			// status_array는 고정값이므로 없는 경우 default로 '0'으로 set
			_, exists := resultVmStatusCountMap[util.STATUS_ARRAY[i]]
			if !exists {
				resultVmStatusCountMap[util.STATUS_ARRAY[i]] = 0
			}
			totalVmStatusCountMap[util.STATUS_ARRAY[i]] += resultVmStatusCountMap[util.STATUS_ARRAY[i]]
		}
		// UI manage mcis > server 영역에서는 run/stopped/terminated 만 있음. etc를 stopped에 추가한다.
		totalVmStatusCountMap["stopped"] = totalVmStatusCountMap["stopped"] + resultVmStatusCountMap[util.VM_STATUS_ETC]

		totalVmCount += resultVmStatusCountMap["TOTAL"] // 모든 vm의 갯수

		totalVmCountByMcis := resultVmStatusCountMap["TOTAL"]        // 모든 vm의 갯수
		vmStatusCountMapByMcis[mcisInfo.ID] = resultVmStatusCountMap // MCIS 내 vm 상태별 cnt

		// Provider 별 connection count (Location 내에 있는 provider로 갯수 셀 것.)
		mcisConnectionMap := service.GetVmConnectionCountByMcis(mcisInfo) // 해당 MCIS의 각 provider별 connection count
		log.Println(mcisConnectionMap)

		mcisConnectionNames := ""
		for connectKey, _ := range mcisConnectionMap {
			mcisConnectionNames += connectKey + " "
		}
		////////////// return value 에 set
		mcisSimpleInfo := webtool.McisSimpleInfo{}
		mcisSimpleInfo.ID = mcisInfo.ID
		mcisSimpleInfo.Status = mcisInfo.Status
		mcisSimpleInfo.McisStatus = util.GetMcisStatus(mcisInfo.Status)
		mcisSimpleInfo.Name = mcisInfo.Name
		mcisSimpleInfo.Description = mcisInfo.Description

		mcisSimpleInfo.InstallMonAgent = mcisInfo.InstallMonAgent
		mcisSimpleInfo.Label = mcisInfo.Label

		mcisSimpleInfo.VmCount = totalVmCountByMcis // 해당 mcis의 모든 vm 갯수

		mcisSimpleInfo.VmSimpleList = resultVmSimpleList
		mcisSimpleInfo.VmStatusNames = resultVmStatusNames
		mcisSimpleInfo.VmStatusCountMap = resultVmStatusCountMap

		mcisSimpleInfo.ConnectionConfigProviderMap = mcisConnectionMap     // 해당 MCIS 등록된 connection의 provider 목록
		mcisSimpleInfo.ConnectionConfigProviderNames = mcisConnectionNames // 해당 MCIS 등록된 connection의 provider 목록을 String
		mcisSimpleInfo.ConnectionConfigProviderCount = len(mcisConnectionMap)

		mcisSimpleInfoList = append(mcisSimpleInfoList, mcisSimpleInfo)

	}

	return echotemplate.Render(c, http.StatusOK,
		"operation/dashboards/DashboardGlobalMng", // 파일명
		map[string]interface{}{
			"LoginInfo":          loginInfo,
			"DefaultNameSpaceID": defaultNameSpaceID,
			"NameSpaceList":      nsList,
			// "McisList":           mcisList,	// mcisSimpleInfoList 로 대체

			// server count 영역
			"TotalVmCount":          totalVmCount,
			"TotalVMStatusCountMap": totalVmStatusCountMap, // 모든 VmStatus 별 count Map(MCIS 무관)

			// cp count 영역
			"TotalProviderCount":         providerCount,            // VM이 등록 된 provider 목록
			"TotalConnectionConfigCount": totalConnectionCount,     // 총 connection 갯수
			"ConnectionConfigCountMap":   connectionConfigCountMap, // provider별 connection 수

			// mcis count 영역
			// "TotalMCISCount":          totalMcisCount,
			// "TotalMCISStatusCountMap": totalMcisStatusCountMap, // 모든 MCIS의 상태 Map

			// mcis list
			"McisList":               mcisSimpleInfoList,     // 표에 뿌려줄 mics summary 정보
			"VmStatusCountMapByMcis": vmStatusCountMapByMcis, // MCIS ID 별 vmStatusMap
		})

	// comURL := service.GetCommonURL()
	// apiInfo := service.AuthenticationHandler()
	// nsCnt := service.GetNSCnt()
	// fmt.Println("=========== DashBoard start ==============")
	// if loginInfo := service.CallLoginInfo(c); loginInfo.UserID != "" {
	// 	nameSpace := service.GetNameSpaceToString(c)
	// 	if nameSpace != "" {
	// 		fmt.Println("Namespace : ", nameSpace)
	// 		return c.Render(http.StatusOK, "Dashboard_Global.html", map[string]interface{}{
	// 			"LoginInfo": loginInfo,
	// 			"NameSpace": nameSpace,
	// 			"comURL":    comURL,
	// 			"apiInfo":   apiInfo,
	// 			"nsCnt":     nsCnt,
	// 		})
	// 	} else {
	// 		return c.Redirect(http.StatusTemporaryRedirect, "/NS/reg")
	// 	}

	// }

	// return c.Redirect(http.StatusTemporaryRedirect, "/login")
}

// func DashBoard(c echo.Context) error {
// 	comURL := service.GetCommonURL()
// 	apiInfo := service.AuthenticationHandler()
// 	fmt.Println("=========== DashBoard start ==============")
// 	if loginInfo := service.CallLoginInfo(c); loginInfo.UserID != "" {
// 		nameSpace := service.GetNameSpaceToString(c)
// 		if nameSpace != "" {
// 			fmt.Println("Namespace : ", nameSpace)
// 			return c.Render(http.StatusOK, "dashboard.html", map[string]interface{}{
// 				"LoginInfo": loginInfo,
// 				"NameSpace": nameSpace,
// 				"comURL":    comURL,
// 				"apiInfo":   apiInfo,
// 			})
// 		} else {
// 			return c.Redirect(http.StatusTemporaryRedirect, "/NS/reg")
// 		}

// 	}

// 	return c.Redirect(http.StatusTemporaryRedirect, "/login")
// }

func NSDashBoard(c echo.Context) error {
	// comURL := service.GetCommonURL()
	// apiInfo := service.AuthenticationHandler()
	// nsCnt := service.GetNSCnt()
	// fmt.Println("=========== DashBoard start ==============")
	// if loginInfo := service.CallLoginInfo(c); loginInfo.UserID != "" {
	// 	nameSpace := service.GetNameSpaceToString(c)
	// 	if nameSpace != "" {
	// 		fmt.Println("Namespace : ", nameSpace)
	// 		return c.Render(http.StatusOK, "Dashboard_Ns.html", map[string]interface{}{
	// 			"LoginInfo": loginInfo,
	// 			"NameSpace": nameSpace,
	// 			"comURL":    comURL,
	// 			"apiInfo":   apiInfo,
	// 			"nsCnt":     nsCnt,
	// 		})
	// 	} else {
	// 		return c.Redirect(http.StatusTemporaryRedirect, "/NS/reg")
	// 	}

	// }

	return c.Redirect(http.StatusTemporaryRedirect, "/login")
}

func IndexController(c echo.Context) error {

	fmt.Println("=========== DashBoard start ==============")
	if loginInfo := service.CallLoginInfo(c); loginInfo.UserID != "" {

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

// func GeoInfo(c echo.Context) error {
// 	//goroutine sync wg
// 	var wg sync.WaitGroup
// 	nameSpace := service.GetNameSpaceToString(c)
// 	comURL := service.GetCommonURL()
// 	//apiInfo := service.AuthenticationHandler()
// 	tumble_url := comURL.TumbleBugURL

// 	mcis_id := c.Param("mcis_id")
// 	url := tumble_url + "/ns/" + nameSpace + "/mcis/" + mcis_id
// 	fmt.Println("===========")
// 	fmt.Println("=========== GetGeoINFO ==============")
// 	fmt.Println("=========== GetGeoINFO request URL : ", url)
// 	body := service.HttpGetHandler(url)

// 	defer body.Close()

// 	publicIpInfo := RespPublicIPInfo{}
// 	json.NewDecoder(body).Decode(&publicIpInfo)
// 	fmt.Println("================mcis ID info ===============")
// 	fmt.Println("Public Info : ", publicIpInfo)
// 	var ipStackInfo []service.IPStackInfo

// 	for _, item := range publicIpInfo.PublicIPInfo {
// 		wg.Add(1)

// 		go service.GetGeoMetryInfo(&wg, item.PublicIp, item.Status, item.VmID, item.VmName, &ipStackInfo)

// 	}
// 	wg.Wait()
// 	fmt.Println("DashBoard ipStackInfo  : ", ipStackInfo)
// 	return c.JSON(http.StatusOK, ipStackInfo)
// }
