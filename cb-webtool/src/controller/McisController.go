package controller

import (
	// "encoding/json"
	"fmt"
	"log"
	"net/http"

	// model "github.com/cloud-barista/cb-webtool/src/model"
	"github.com/cloud-barista/cb-webtool/src/model"
	"github.com/cloud-barista/cb-webtool/src/model/dragonfly"
	spider "github.com/cloud-barista/cb-webtool/src/model/spider"
	"github.com/cloud-barista/cb-webtool/src/model/tumblebug"

	service "github.com/cloud-barista/cb-webtool/src/service"
	util "github.com/cloud-barista/cb-webtool/src/util"

	echotemplate "github.com/foolin/echo-template"
	echosession "github.com/go-session/echo-session"
	"github.com/labstack/echo"
	// echosession "github.com/go-session/echo-session"
)

// type SecurityGroup struct {
// 	Id []string `form:"sg"`
// }

func McisRegForm(c echo.Context) error {
	fmt.Println("McisRegForm ************ : ")

	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}
	defaultNameSpaceID := loginInfo.DefaultNameSpaceID

	// namespacelist 가져오기
	// nsList, _ := service.GetNameSpaceList()
	nsList, _ := service.GetStoredNameSpaceList(c)
	log.Println(" nsList  ", nsList)

	// connectionconfigList 가져오기
	cloudOsList, _ := service.GetCloudOSList()
	log.Println(" cloudOsList  ", cloudOsList)

	// regionList 가져오기
	regionList, _ := service.GetRegionList()
	log.Println(" regionList  ", regionList)

	cloudConnectionConfigInfoList, _ := service.GetCloudConnectionConfigList() // 등록된 모든 connection 정보
	log.Println("---------------------- GetCloudConnectionConfigList ", defaultNameSpaceID)

	//// namespace에 등록 된 resource 정보들 //////
	virtualMachineImageInfoList, _ := service.GetVirtualMachineImageInfoList(defaultNameSpaceID)
	vmSpecInfoList, _ := service.GetVmSpecInfoList(defaultNameSpaceID)
	vNetInfoList, _ := service.GetVnetList(defaultNameSpaceID)
	securityGroupInfoList, _ := service.GetSecurityGroupList(defaultNameSpaceID)
	sshKeyInfoList, _ := service.GetSshKeyInfoList(defaultNameSpaceID)

	return echotemplate.Render(c, http.StatusOK,
		"operation/manages/mcismng/McisCreate", // 파일명
		map[string]interface{}{
			"LoginInfo":          loginInfo,
			"DefaultNameSpaceID": defaultNameSpaceID,
			"NameSpaceList":      nsList,
			"CloudOSList":        cloudOsList,
			"RegionList":         regionList,

			"CloudConnectionConfigInfoList": cloudConnectionConfigInfoList,
			"VMImageList":                   virtualMachineImageInfoList,
			"VMSpecList":                    vmSpecInfoList,
			"VNetList":                      vNetInfoList,
			"SecurityGroupList":             securityGroupInfoList,
			"SshKeyList":                    sshKeyInfoList,
		})
}

// MCIS 관리 화면 McisListForm 에서 이름 변경 McisMngForm으로
// func McisListForm(c echo.Context) error {
func McisMngForm(c echo.Context) error {
	fmt.Println("McisMngForm ************ : ")

	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	defaultNameSpaceID := loginInfo.DefaultNameSpaceID
	selectedMcisID := c.QueryParam("mcisid") // Dashboard 등에서 선택한 mcis가 있는경우 mng 화면에 해당 mcis만 보이기 위해(실제로는 filterling을 위해서만 사용)
	store := echosession.FromContext(c)

	// 최신 namespacelist 가져오기
	//nsList, _ := service.GetNameSpaceList()
	nsList, _ := service.GetStoredNameSpaceList(c)
	log.Println(" nsList  ", nsList)

	// provider 별 연결정보 count(MCIS 무관)
	cloudConnectionConfigInfoList, _ := service.GetCloudConnectionConfigList()
	connectionConfigCountMap, providerCount := service.GetCloudConnectionCountMap(cloudConnectionConfigInfoList)
	totalConnectionCount := len(cloudConnectionConfigInfoList)

	mcisList := []tumblebug.McisInfo{}
	mcisErr := model.WebStatus{}

	totalMcisCount := 0 // mcis 갯수
	totalVmCount := 0

	totalMcisStatusCountMap := make(map[string]int)             // 모든 MCIS의 상태 Map
	mcisStatusCountMapByMcis := make(map[string]map[string]int) // MCIS ID별 mcis status
	totalVmStatusCountMap := make(map[string]int)               // 모든 VM의 상태 Map
	vmStatusCountMapByMcis := make(map[string]map[string]int)   // MCIS ID 별 vmStatusMap
	mcisSimpleInfoList := []tumblebug.McisSimpleInfo{}          // 표에 뿌려줄 mics summary 정보

	cloudOsList := []string{}
	regionInfoList := []spider.RegionInfo{}
	regionErr := model.WebStatus{}
	virtualMachineImageInfoList := []tumblebug.VirtualMachineImageInfo{}
	vmSpecInfoList := []tumblebug.VmSpecInfo{}
	vNetInfoList := []tumblebug.VNetInfo{}
	securityGroupInfoList := []tumblebug.SecurityGroupInfo{}

	mcisList, mcisErr = service.GetMcisList(defaultNameSpaceID)
	log.Println(" mcisList  ", mcisList, mcisErr.StatusCode)
	if mcisErr.StatusCode != 200 && mcisErr.StatusCode != 201 {
		//return render(c, "ErrorPage.html", map[string]interface{}{"Message": mcisErr.Message, "StatusCode": mcisErr.StatusCode, "PageUrl": "/main"})
		//return render(c, "ErrorPage.html", map[string]interface{}{"Message": messages, "StatusCode": msgtype, "PageUrl": template.HTML(pageHtml)})
		//return c.Render(http.StatusNotFound, tmpl.ErrNotFoundTpl, tmpl.NotFoundMessage)

		// errPage := fmt.Sprintf("%d.html", mcisErr.StatusCode)
		// if err := c.File(errPage); err != nil {
		// 	log.Println("Return error page")
		// 	return c.HTML(mcisErr.StatusCode, "Return Error page Message : "+mcisErr.Message)
		// }
		return echotemplate.Render(c, http.StatusOK,
			"operation/manages/mcismng/McisMng", // 파일명
			map[string]interface{}{
				"Message":            mcisErr.Message,
				"Status":             mcisErr.StatusCode,
				"LoginInfo":          loginInfo,
				"DefaultNameSpaceID": defaultNameSpaceID,
				"SelectedMcisID":     selectedMcisID, // 선택한 MCIS ID
				"NameSpaceList":      nsList,

				// mcis count 영역
				"TotalMcisCount":          totalMcisCount,
				"TotalMcisStatusCountMap": totalMcisStatusCountMap, // 모든 MCIS의 상태 Map

				// server count 영역
				"TotalVmCount":          totalVmCount,
				"TotalVmStatusCountMap": totalVmStatusCountMap, // 모든 VmStatus 별 count Map(MCIS 무관)

				// cp count 영역
				"TotalProviderCount":         providerCount,            // VM이 등록 된 provider 목록
				"TotalConnectionConfigCount": totalConnectionCount,     // 총 connection 갯수
				"ConnectionConfigCountMap":   connectionConfigCountMap, // provider별 connection 수
				// mcis list
				"McisList":               mcisSimpleInfoList,     // 표에 뿌려줄 mics summary 정보
				"VmStatusCountMapByMcis": vmStatusCountMapByMcis, // MCIS ID 별 vmStatusMap

				"CloudOSList":                   cloudOsList,
				"RegionList":                    regionInfoList,
				"CloudConnectionConfigInfoList": cloudConnectionConfigInfoList,
				"VMImageList":                   virtualMachineImageInfoList,
				"VMSpecList":                    vmSpecInfoList,
				"VNetList":                      vNetInfoList,
				"SecurityGroupList":             securityGroupInfoList,
			})
	}

	// 모든 MCIS 조회
	// mcisList, mcisErr := service.GetMcisList(defaultNameSpaceID)
	// log.Println(" mcisList  ", mcisList, mcisErr.StatusCode)
	// if mcisErr.StatusCode != 200 && mcisErr.StatusCode != 201 {
	// 	//return render(c, "ErrorPage.html", map[string]interface{}{"Message": mcisErr.Message, "StatusCode": mcisErr.StatusCode, "PageUrl": "/main"})
	// 	//return render(c, "ErrorPage.html", map[string]interface{}{"Message": messages, "StatusCode": msgtype, "PageUrl": template.HTML(pageHtml)})
	// 	//return c.Render(http.StatusNotFound, tmpl.ErrNotFoundTpl, tmpl.NotFoundMessage)

	// 	errPage := fmt.Sprintf("%d.html", mcisErr.StatusCode)
	// 	if err := c.File(errPage); err != nil {
	// 		log.Println("Return error page")
	// 		return c.HTML(mcisErr.StatusCode, "Return Error page Message : "+mcisErr.Message)
	// 	}
	// }

	totalMcisCount = len(mcisList) // mcis 갯수
	totalVmCount = 0               // 모든 vm 갯수

	if totalMcisCount == 0 {
		return c.Redirect(http.StatusTemporaryRedirect, "/operation/manages/mcismng/regform")
	}

	// TODO : store에 MCIS내 VM정보를 저장했다가 상세정보 조회시 사용
	// loginInfo.vMList
	// store := echosession.FromContext(c)
	// result, ok := store.Get(loginInfo.UserID)
	// if !ok {
	// 	// user의 mcis내 vm
	// 	store.Set("", nsList) // 이게 유효한가?? 쓸모없을 듯
	// 	store.Save()
	// }

	// result := getObj.(map[string]string)
	// loginInfo := LoginInfo{
	// 	Username: "admin",
	// 	//Username:  result["username"],
	// 	NameSpace: result["namespace"],
	// }
	// getNs, ok := store.Get("namespace")
	// if !ok {
	// 	return loginInfo
	// }
	// loginInfo.NameSpace = getNs.(string)

	// vmList := result

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
		resultSimpleVmList, resultVmStatusCountMap := service.GetSimpleVmWithStatusCountMap(mcisInfo)

		resultVmStatusNames := ""
		for _, vmSimpleObj := range resultSimpleVmList {
			resultVmStatusNames += vmSimpleObj.VmID + "|" + vmSimpleObj.VmName + "@"
		}

		log.Println("before " + resultVmStatusNames)
		if len(resultVmStatusNames) > 0 {
			resultVmStatusNames = resultVmStatusNames[:len(resultVmStatusNames)-1]
		}
		log.Println("after " + resultVmStatusNames)

		// UI에서 보여 줄 VM_STATUS로 Count. (가져온 Key중에 UI에서 보여줄 Key가 없을 수 있으므로)
		for i, _ := range util.VM_STATUS_ARRAY {
			// status_array는 고정값이므로 없는 경우 default로 '0'으로 set
			_, exists := resultVmStatusCountMap[util.VM_STATUS_ARRAY[i]]
			if !exists {
				resultVmStatusCountMap[util.VM_STATUS_ARRAY[i]] = 0
			}
			totalVmStatusCountMap[util.VM_STATUS_ARRAY[i]] += resultVmStatusCountMap[util.VM_STATUS_ARRAY[i]]
		}
		// UI manage mcis > server 영역에서는 run/stopped/terminated 만 있음. etc를 stopped에 추가한다.
		totalVmStatusCountMap[util.VM_STATUS_STOPPED] = totalVmStatusCountMap[util.VM_STATUS_STOPPED] + resultVmStatusCountMap[util.VM_STATUS_ETC]

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
		mcisSimpleInfo := tumblebug.McisSimpleInfo{}
		mcisSimpleInfo.ID = mcisInfo.ID
		mcisSimpleInfo.Status = mcisInfo.Status
		mcisSimpleInfo.McisStatus = util.GetMcisStatus(mcisInfo.Status)
		mcisSimpleInfo.Name = mcisInfo.Name
		mcisSimpleInfo.Description = mcisInfo.Description

		mcisSimpleInfo.InstallMonAgent = mcisInfo.InstallMonAgent
		mcisSimpleInfo.Label = mcisInfo.Label

		mcisSimpleInfo.VmCount = totalVmCountByMcis // 해당 mcis의 모든 vm 갯수
		mcisSimpleInfo.VmSimpleList = resultSimpleVmList
		mcisSimpleInfo.VmStatusNames = resultVmStatusNames
		mcisSimpleInfo.VmStatusCountMap = resultVmStatusCountMap
		// mcisSimpleInfo.VmRunningCount = vmStatusCountMap[util.STATUS_ARRAY[0]]    //running
		// mcisSimpleInfo.VmStoppedCount = vmStatusCountMap[util.STATUS_ARRAY[1]]    //stopped
		// mcisSimpleInfo.VmTerminatedCount = vmStatusCountMap[util.STATUS_ARRAY[2]] //terminated

		mcisSimpleInfo.ConnectionConfigProviderMap = mcisConnectionMap     // 해당 MCIS 등록된 connection의 provider 목록
		mcisSimpleInfo.ConnectionConfigProviderNames = mcisConnectionNames // 해당 MCIS 등록된 connection의 provider 목록을 String
		mcisSimpleInfo.ConnectionConfigProviderCount = len(mcisConnectionMap)
		// mcisConnectionMap.ConnectionCount = mcisConnectionMap

		mcisSimpleInfoList = append(mcisSimpleInfoList, mcisSimpleInfo)

	}

	// log.Println(" totoalMcisCount  ", totoalMcisCount)
	// log.Println(" totoalVmCount  ", totoalVmCount)

	// // mcis 별 vmCnt
	// // mcisSimpleInfos = model.McisSimpleInfos{}
	// connectionCountTotal := 0
	// connectionCountByMcis := 0
	// vmCountTotal := 0
	// vmRunningCountByMcis := 0
	// vmStoppedCountByMcis := 0
	// vmTerminatedCountByMcis := 0
	// for mcisIndex, mcisInfo := range mcisList {
	// 	// mcis.ID, mcis.status, mcis.name, mcis.description
	// 	// csp : 해당 mcis의 connection cnt
	// 	// vm_cnt : 해당 mcis의 vm cnt
	// 	// vm_run_cnt, vm_stop_cnt
	// 	vmList := mcisInfo.VMs
	// 	mcisConnectionCountMap := make(map[string]int)
	// 	mcisVmStatusCountMap := make(map[string]int)
	// 	for vmIndex, vmInfo := range vmList {
	// 		locationInfo := vmInfo.LocationInfo
	// 		cloudType := locationInfo.CloudType // CloudConnection
	// 		providerCount := 0
	// 		val, exists := mcisConnectionCountMap[util.GetProviderName(locationInfo.CloudType)]
	// 		if !exists {
	// 			providerCount = 1
	// 		} else {
	// 			providerCount = val + 1
	// 		}
	// 		mcisConnectionCountMap[util.GetProviderName(locationInfo.CloudType)] = providerCount

	// 		vmStatus := vmInfo.Status
	// 		vnStatusCount := 0
	// 		val2, exists2 := mcisVmStatusCountMap[util.GetVmStatus(vmInfo.Status)]
	// 		if !exists2 {
	// 			vnStatusCount = 1
	// 		} else {
	// 			vnStatusCount = val2 + 1
	// 		}
	// 		mcisVmStatusCountMap[util.GetVmStatus(vmInfo.Status)] = vnStatusCount
	// 	}

	// 	mcisSimpleInfo := model.McisSimpleInfo{}
	// 	mcisSimpleInfo.ID = mcisInfo.ID
	// 	mcisSimpleInfo.Status = mcisInfo.Status
	// 	mcisSimpleInfo.Name = mcisInfo.Name
	// 	mcisSimpleInfo.Description = mcisInfo.Description

	// 	mcisSimpleInfo.VmCount = len(vmList)
	// 	mcisSimpleInfo.VmRunningCount = mcisVmStatusCountMap[util.VM_STATUS_RUNNING]
	// 	mcisSimpleInfo.VmStoppedCount = mcisVmStatusCountMap[util.VM_STATUS_RUNNING]
	// 	mcisSimpleInfo.VmTerminatedCount = mcisVmStatusCountMap[util.VM_STATUS_RUNNING]
	// }

	storedCloudOsList, ok := store.Get("cloudoslist")
	if !ok {
		cloudOsList, _ = service.GetCloudOSList()
		log.Println(" cloudOsList  ", cloudOsList)
		//mcisList = storedMcisList.([]tumblebug.McisInfo)
		store.Set("cloudoslist", cloudOsList)
	} else {
		cloudOsList = storedCloudOsList.([]string)
	}
	// cloudOsList, _ := service.GetCloudOSList() // provider
	log.Println("---------------------- GetCloudOSList ", defaultNameSpaceID)

	// Region 목록

	storedRegionList, ok := store.Get("regionlist")
	if !ok {
		regionInfoList, regionErr = service.GetRegionList()
		log.Println(" storedRegionList  ", regionInfoList, regionErr.StatusCode)

		store.Set("regionlist", regionInfoList)
	} else {
		regionInfoList = storedRegionList.([]spider.RegionInfo)
	}
	// regionInfoList, _ := service.GetRegionList() // region
	log.Println("---------------------- GetRegionList ", defaultNameSpaceID)
	// cloudConnectionConfigInfoList, _ := service.GetCloudConnectionConfigList() // 등록된 모든 connection 정보
	// log.Println("---------------------- GetCloudConnectionConfigList ", defaultNameSpaceID);
	//// namespace에 등록 된 resource 정보들 //////

	// VM Image 목록
	// virtualMachineImageErr := model.WebStatus{}

	// storedvirtualMachineImageList, ok := store.Get("MCIS_VMIMAGE_" + defaultNameSpaceID)
	// if !ok {
	// 	virtualMachineImageInfoList, virtualMachineImageErr = service.GetVirtualMachineImageInfoList(defaultNameSpaceID)
	// 	log.Println(" virtualMachineImageInfoList  ", virtualMachineImageInfoList, virtualMachineImageErr.StatusCode)

	// 	store.Set("MCIS_VMIMAGE_"+defaultNameSpaceID, virtualMachineImageInfoList)
	// } else {
	// 	virtualMachineImageInfoList = storedvirtualMachineImageList.([]tumblebug.VirtualMachineImageInfo)
	// }
	// virtualMachineImageInfoList, _ := service.GetVirtualMachineImageInfoList(defaultNameSpaceID)

	// VMSpec 목록

	// vmSpecErr := model.WebStatus{}

	// storedVmSpecList, ok := store.Get("MCIS_VMSPEC_" + defaultNameSpaceID)
	// if !ok {
	// 	vmSpecInfoList, vmSpecErr = service.GetVmSpecInfoList(defaultNameSpaceID)
	// 	log.Println(" virtualMachineImageInfoList  ", vmSpecInfoList, vmSpecErr.StatusCode)

	// 	store.Set("MCIS_VMSPEC_"+defaultNameSpaceID, vmSpecInfoList)
	// } else {
	// 	vmSpecInfoList = storedVmSpecList.([]tumblebug.VmSpecInfo)
	// }
	// vmSpecInfoList, _ := service.GetVmSpecInfoList(defaultNameSpaceID)

	// vNet 목록

	// vNetErr := model.WebStatus{}

	// storedVnetList, ok := store.Get("MCIS_VNET_" + defaultNameSpaceID)
	// if !ok {
	// 	vNetInfoList, vNetErr = service.GetVnetList(defaultNameSpaceID)
	// 	log.Println(" vNetInfoList  ", vNetInfoList, vNetErr.StatusCode)

	// 	store.Set("MCIS_VNET_"+defaultNameSpaceID, vNetInfoList)
	// } else {
	// 	vNetInfoList = storedVnetList.([]tumblebug.VNetInfo)
	// }
	// vNetInfoList, _ := service.GetVnetList(defaultNameSpaceID)

	// SecurityGroup

	// securityGroupErr := model.WebStatus{}

	// storedSecurityGroupList, ok := store.Get("MCIS_SECURITYGROUP_" + defaultNameSpaceID)
	// if !ok {
	// 	securityGroupInfoList, securityGroupErr = service.GetSecurityGroupList(defaultNameSpaceID)
	// 	log.Println(" securityGroupInfoList  ", storedSecurityGroupList, securityGroupErr.StatusCode)

	// 	store.Set("MCIS_SECURITYGROUP_"+defaultNameSpaceID, securityGroupInfoList)
	// } else {
	// 	securityGroupInfoList = storedSecurityGroupList.([]tumblebug.SecurityGroupInfo)
	// }
	// securityGroupInfoList, _ := service.GetSecurityGroupList(defaultNameSpaceID)

	// status, filepath, return params
	return echotemplate.Render(c, http.StatusOK,
		"operation/manages/mcismng/McisMng", // 파일명
		map[string]interface{}{
			"Message":            mcisErr.Message,
			"Status":             mcisErr.StatusCode, // 주요한 객체 return message 를 사용
			"LoginInfo":          loginInfo,
			"DefaultNameSpaceID": defaultNameSpaceID,
			"SelectedMcisID":     selectedMcisID, // 선택한 MCIS ID
			"NameSpaceList":      nsList,
			// "McisList":           mcisList,	// mcisSimpleInfoList 로 대체
			// "McisIDList":         mcisIdArr,
			// "VmIDList":           vmIdArr,
			// "VMStatusList":  vmStatusArr,
			// "MCISStatusMap":            mcisStatusMap,
			// "MCISCount":                totoalMcisCount,
			// "VMStatusMap":              vmStatusMap,
			// "VMCount":                  totoalVmCount,
			// "ConnectionConfigCountMap": connectionConfigCountMap,
			// "ConnectionCount":          totalConnectionCount,
			// "ProviderCount":            providerCount,

			// mcis count 영역
			"TotalMcisCount":          totalMcisCount,
			"TotalMcisStatusCountMap": totalMcisStatusCountMap, // 모든 MCIS의 상태 Map

			// server count 영역
			"TotalVmCount":          totalVmCount,
			"TotalVmStatusCountMap": totalVmStatusCountMap, // 모든 VmStatus 별 count Map(MCIS 무관)

			// cp count 영역
			"TotalProviderCount":         providerCount,            // VM이 등록 된 provider 목록
			"TotalConnectionConfigCount": totalConnectionCount,     // 총 connection 갯수
			"ConnectionConfigCountMap":   connectionConfigCountMap, // provider별 connection 수
			// mcis list
			"McisList":               mcisSimpleInfoList,     // 표에 뿌려줄 mics summary 정보
			"VmStatusCountMapByMcis": vmStatusCountMapByMcis, // MCIS ID 별 vmStatusMap

			"CloudOSList":                   cloudOsList,
			"RegionList":                    regionInfoList,
			"CloudConnectionConfigInfoList": cloudConnectionConfigInfoList,
			"VMImageList":                   virtualMachineImageInfoList,
			"VMSpecList":                    vmSpecInfoList,
			"VNetList":                      vNetInfoList,
			"SecurityGroupList":             securityGroupInfoList,
		})
}

// MCIS 목록 조회
func GetMcisList(c echo.Context) error {
	log.Println("GetMcisList : ")
	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	// store := echosession.FromContext(c)
	defaultNameSpaceID := loginInfo.DefaultNameSpaceID
	// TODO : defaultNameSpaceID 가 없으면 설정화면으로 보낼 것
	mcisList, respStatus := service.GetMcisList(defaultNameSpaceID)
	if respStatus.StatusCode != 200 && respStatus.StatusCode != 201 {
		return c.JSON(respStatus.StatusCode, map[string]interface{}{
			"error":  respStatus.Message,
			"status": respStatus.StatusCode,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":            "success",
		"status":             respStatus.StatusCode,
		"DefaultNameSpaceID": defaultNameSpaceID,
		"McisList":           mcisList,
	})
}

// MCIS 등록
func McisRegProc(c echo.Context) error {
	log.Println("McisRegProc : ")
	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	// json_map := make(map[string]interface{})
	// err := json.NewDecoder(c.Request().Body).Decode(&json_map)
	// if err != nil {
	// 	return err
	// } else {
	// 	log.Println(json_map)
	// }

	// map[description:bb installMonAgent:yes name:aa vm:[map[connectionName:gcp-asia-east1 description:dd imageId:gcp-jsyoo-ubuntu name:cc provider:GCP securityGroupIds:[gcp-jsyoo-sg-01] specId:gcp-jsyoo-01 sshKeyId:gcp-jsyoo-sshkey subnetId:jsyoo-gcp-sub-01 vNetId:jsyoo-gcp-01 vm_add_cnt:0 vm_cnt:]]]
	log.Println("get info")
	//&[]Person{}
	mCISInfo := &tumblebug.McisInfo{}
	if err := c.Bind(mCISInfo); err != nil {
		// if err := c.Bind(mCISInfoList); err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "fail",
			"status":  "fail",
		})
	}
	log.Println(mCISInfo) // 여러개일 수 있음.

	// store := echosession.FromContext(c)
	defaultNameSpaceID := loginInfo.DefaultNameSpaceID
	// TODO : defaultNameSpaceID 가 없으면 설정화면으로 보낼 것
	_, respStatus := service.RegMcis(defaultNameSpaceID, mCISInfo)
	log.Println("RegMcis service returned")
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

// MCIS 삭제
func McisDelProc(c echo.Context) error {
	log.Println("McisDelProc : ")
	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	// store := echosession.FromContext(c)
	defaultNameSpaceID := loginInfo.DefaultNameSpaceID
	// TODO : defaultNameSpaceID 가 없으면 설정화면으로 보낼 것

	mcisID := c.Param("mcisID")
	log.Println("mcisID= " + mcisID)
	_, respStatus := service.DelMcis(defaultNameSpaceID, mcisID)
	log.Println("RegMcis service returned")
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

// server instance 등록

// func McisListFormWithParam(c echo.Context) error {
// 	mcis_id := c.Param("mcis_id")
// 	mcis_name := c.Param("mcis_name")
// 	comURL := service.GetCommonURL()
// 	apiInfo := service.AuthenticationHandler()
// 	if mcis_id == "" && mcis_name == "" {
// 		mcis_id = ""
// 		mcis_name = ""
// 	}
// 	if loginInfo := service.CallLoginInfo(c); loginInfo.Username != "" {
// 		namespace := service.GetNameSpaceToString(c)
// 		return c.Render(http.StatusOK, "Manage_Mcis.html", map[string]interface{}{
// 			"LoginInfo": loginInfo,
// 			"NameSpace": namespace,
// 			"McisID":    mcis_id,
// 			"McisName":  mcis_name,
// 			"comURL":    comURL,
// 			"apiInfo":   apiInfo,
// 		})

// 	}

// 	//return c.Render(http.StatusOK, "MCISlist.html", nil)
// 	return c.Redirect(http.StatusTemporaryRedirect, "/login")
// }

// MCIS에 VM 추가
func McisVmRegForm(c echo.Context) error {
	log.Println("McisVmRegForm : ")
	mcisId := c.Param("mcisID")
	mcisName := c.Param("mcisName")

	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	defaultNameSpaceID := loginInfo.DefaultNameSpaceID

	// 최신 namespacelist 가져오기
	nsList, _ := service.GetNameSpaceList()
	log.Println(" nsList  ", nsList)

	// MCIS 정보는 받은것으로

	// MCIS 조회
	resultMcisInfo, _ := service.GetMcisData(defaultNameSpaceID, mcisId) // TODO : store에 있는 것 꺼내쓰도록.  주기적으로 store 갱신.
	log.Println(" resultMcisInfo  ", resultMcisInfo)

	// vm List
	vmList := resultMcisInfo.Vms

	///////// 등록을 위한 정보 ////////////
	cloudOsList, _ := service.GetCloudOSList() // provider
	log.Println("---------------------- GetCloudOSList ", defaultNameSpaceID)
	regionInfoList, _ := service.GetRegionList() // region
	log.Println("---------------------- GetRegionList ", defaultNameSpaceID)
	cloudConnectionConfigInfoList, _ := service.GetCloudConnectionConfigList() // 등록된 모든 connection 정보
	log.Println("---------------------- GetCloudConnectionConfigList ", defaultNameSpaceID)

	//// namespace에 등록 된 resource 정보들 //////
	virtualMachineImageInfoList, _ := service.GetVirtualMachineImageInfoList(defaultNameSpaceID)
	vmSpecInfoList, _ := service.GetVmSpecInfoList(defaultNameSpaceID)
	vNetInfoList, _ := service.GetVnetList(defaultNameSpaceID)
	securityGroupInfoList, _ := service.GetSecurityGroupList(defaultNameSpaceID)
	sshKeyInfoList, _ := service.GetSshKeyInfoList(defaultNameSpaceID)

	// status, filepath, return params
	return echotemplate.Render(c, http.StatusOK,
		"operation/manages/mcismng/McisVmCreate", // 파일명
		map[string]interface{}{
			"LoginInfo":          loginInfo,
			"DefaultNameSpaceID": defaultNameSpaceID,
			"NameSpaceList":      nsList,
			"McisID":             mcisId,
			"McisName":           mcisName,
			"VMList":             vmList,

			"CloudOSList":                   cloudOsList,
			"RegionList":                    regionInfoList,
			"CloudConnectionConfigInfoList": cloudConnectionConfigInfoList,
			"VMImageList":                   virtualMachineImageInfoList,
			"VMSpecList":                    vmSpecInfoList,
			"VNetList":                      vNetInfoList,
			"SecurityGroupList":             securityGroupInfoList,
			"SshKeyList":                    sshKeyInfoList,
		})
}

// GetMcisInfoData
// 특정 MCIS의 상세정보를 가져온다.
func GetMcisInfoData(c echo.Context) error {
	log.Println("GetMcisInfoData")
	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/login") // 조회기능에서 바로 login화면으로 돌리지말고 return message로 하는게 낫지 않을까?
	}
	defaultNameSpaceID := loginInfo.DefaultNameSpaceID

	mcisID := c.Param("mcisID")
	log.Println("mcisID= " + mcisID)

	resultMcisInfo, _ := service.GetMcisData(defaultNameSpaceID, mcisID)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":  "success",
		"status":   200,
		"McisInfo": resultMcisInfo,
	})
}

// MCIS에 VM 추가 등록
func VmRegProc(c echo.Context) error {
	log.Println("VmRegProc : ")
	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	// mCISInfo := &tumblebug.McisInfo{}
	vmInfo := &tumblebug.VmInfo{}
	if err := c.Bind(vmInfo); err != nil {
		// if err := c.Bind(mCISInfoList); err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "fail",
			"status":  "fail",
		})
	}
	log.Println(vmInfo)

	defaultNameSpaceID := loginInfo.DefaultNameSpaceID
	mcisID := c.Param("mcisID")

	// 일반 호출 : return 값 수신방식
	_, respStatus := service.RegVm(defaultNameSpaceID, mcisID, vmInfo)
	log.Println("RegVM service returned")
	if respStatus.StatusCode != 200 && respStatus.StatusCode != 201 {
		return c.JSON(respStatus.StatusCode, map[string]interface{}{
			"error":  respStatus.Message,
			"status": respStatus.StatusCode,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": respStatus.Message,
		"status":  respStatus.StatusCode,
	})

	// go 루틴
	// go service.AsyncRegVm(defaultNameSpaceID, mcisID, vmInfo)

	// return c.JSON(http.StatusOK, map[string]interface{}{
	// 	"message": "Call success",
	// 	"status":  200,
	// })

}

// MCIS 의 특정 VM의 정보를 가져온다. 단. 텀블벅 조회가 아니라 이미 저장되어 있는 store에서 꺼낸다.
func GetVmInfoData(c echo.Context) error {
	log.Println("GetVmInfoData")
	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/login") // 조회기능에서 바로 login화면으로 돌리지말고 return message로 하는게 낫지 않을까?
	}
	defaultNameSpaceID := loginInfo.DefaultNameSpaceID

	mcisID := c.Param("mcisID")
	vmID := c.Param("vmID")
	log.Println("mcisID= " + mcisID + " , vmID= " + vmID)

	// store := echosession.FromContext(c)
	// mcisObj, ok := store.Get("MCIS_" +loginInfo.UserID)
	// if !ok {
	// 	return c.JSON(http.StatusOK, map[string]interface{}{
	// 		"message": "fail",
	// 		"status":  500,
	// 	})
	// }

	// log.Println("stored key = " + "MCIS_" + loginInfo.UserID)
	// mcisList := mcisObj.([]model.McisInfo)
	// mcisInfo := model.McisInfo{}
	// for _, keyMcisInfo := range mcisList {
	// 	if keyMcisInfo.ID == mcisID {
	// 		mcisInfo = keyMcisInfo
	// 		break;
	// 	}
	// }

	// vmList := mcisInfo.Vms
	// returnVmInfo := model.VMInfo{}
	// if len(vmList) > 0 {
	// 	for _, keyVmInfo := range vmList {
	// 		if keyVmInfo.ID == vmID {
	// 			log.Println("found vm " , keyVmInfo)
	// 			returnVmInfo = keyVmInfo
	// 			break
	// 		}
	// 	}
	// }
	returnVmInfo, respStatus := service.GetVMofMcisData(defaultNameSpaceID, mcisID, vmID)

	if respStatus.StatusCode != 200 && respStatus.StatusCode != 201 {
		return c.JSON(respStatus.StatusCode, map[string]interface{}{
			"error":  respStatus.Message,
			"status": respStatus.StatusCode,
		})
	}

	connectionName := returnVmInfo.ConnectionName
	cloudConnectionConfigInfo, _ := service.GetCloudConnectionConfigData(connectionName)
	// credential Info by connection

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":              respStatus.Message,
		"status":               respStatus.StatusCode,
		"VmInfo":               returnVmInfo,
		"ConnectionConfigInfo": cloudConnectionConfigInfo,
	})
}

// MCIS의 status변경
func McisLifeCycle(c echo.Context) error {
	log.Println("McisLifeCycle : ")
	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	mcisLifeCycle := &tumblebug.McisLifeCycle{}
	if err := c.Bind(mcisLifeCycle); err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "fail",
			"status":  "fail",
		})
	}
	log.Println(mcisLifeCycle)

	// store := echosession.FromContext(c)
	defaultNameSpaceID := loginInfo.DefaultNameSpaceID

	if defaultNameSpaceID != mcisLifeCycle.NameSpaceID {
		// 변경할 Namespace 정보가 다르므로 변경 불가
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "선택된 Namespace가 아닙니다. Namespace를 임의로 변경하여 호출하면 안됨.",
			"status":  "400", // TODO : custom으로 만드는 resultCode 정리 필요
		})
	}

	// TODO : defaultNameSpaceID 가 없으면 설정화면으로 보낼 것
	_, respStatus := service.McisLifeCycle(mcisLifeCycle)
	log.Println("McisLifeCycle service returned")
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

// VM의 LifeCycle status변경
func McisVmLifeCycle(c echo.Context) error {
	log.Println("McisVmLifeCycle : ")
	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}
	log.Println("bind")
	vmLifeCycle := &tumblebug.VmLifeCycle{}
	if err := c.Bind(vmLifeCycle); err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "fail",
			"status":  "fail",
		})
	}
	log.Println(vmLifeCycle)

	// store := echosession.FromContext(c)
	defaultNameSpaceID := loginInfo.DefaultNameSpaceID
	vmLifeCycle.NameSpaceID = defaultNameSpaceID
	// TODO : defaultNameSpaceID 가 없으면 설정화면으로 보낼 것
	_, respStatus := service.McisVmLifeCycle(vmLifeCycle)
	log.Println("McisVmLifeCycle service returned")
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

// VM 통게보기
func GetVmMonitoring(c echo.Context) error {
	log.Println("GetVmInfoData")
	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/login") // 조회기능에서 바로 login화면으로 돌리지말고 return message로 하는게 낫지 않을까?
	}
	defaultNameSpaceID := loginInfo.DefaultNameSpaceID

	// mcisID := c.Param("mcisID")
	// vmID := c.Param("vmID")
	// metric := c.Param("metric")
	// log.Println("mcisID= " + mcisID + " , vmID= " + vmID)

	vmMonitoring := &dragonfly.VmMonitoring{}
	if err := c.Bind(vmMonitoring); err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "fail",
			"status":  "fail",
		})
	}

	vmMonitoring.NameSpaceID = defaultNameSpaceID
	// vmMonitoring.McisID = mcisID
	// vmMonitoring.VmID = vmID
	// vmMonitoring.Metric = metric

	returnVMMonitoringInfo, respStatus := service.GetVmMonitoring(vmMonitoring)
	if respStatus.StatusCode != 200 && respStatus.StatusCode != 201 {
		return c.JSON(respStatus.StatusCode, map[string]interface{}{
			"error":  respStatus.Message,
			"status": respStatus.StatusCode,
		})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":          "success",
		"status":           respStatus.StatusCode,
		"VMMonitoringInfo": returnVMMonitoringInfo,
	})
}

// MCIS에 Command 전송
func CommandMcis(c echo.Context) error {
	log.Println("CommandMcis : ")
	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	mcisCommand := new(tumblebug.McisCommandInfo)
	if err := c.Bind(mcisCommand); err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "fail",
			"status":  "fail",
		})
	}
	log.Println(mcisCommand)

	mcisID := c.Param("mcisID")
	log.Println("mcisID= " + mcisID)

	// store := echosession.FromContext(c)
	defaultNameSpaceID := loginInfo.DefaultNameSpaceID

	// command는 bind 되어있을 것이고.
	mcisCommand.McisID = mcisID
	respMessage, respStatus := service.CommandMcis(defaultNameSpaceID, mcisCommand)
	log.Println("CommandMcis result")
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

// Vm에 Command 전송
func CommandVmOfMcis(c echo.Context) error {
	log.Println("CommandVmOfMcis : ")
	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	// vmCommand := &tumblebug.McisCommandInfo{}
	vmCommand := new(tumblebug.McisCommandInfo)
	if err := c.Bind(vmCommand); err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "fail",
			"status":  "fail",
		})
	}
	log.Println(vmCommand)

	mcisID := c.Param("mcisID")
	vmID := c.Param("vmID")
	log.Println("mcisID= " + mcisID + " , vmID= " + vmID)

	// store := echosession.FromContext(c)
	defaultNameSpaceID := loginInfo.DefaultNameSpaceID

	// 여기에서... sshKey 는 id로 값을 찾아 보내려 했으나 TB에서 알아서 처리 함.
	// sshKeyInfo, _ := service.GetSshKeyData(defaultNameSpaceID, remoteCommandInfo.SshKeyID)
	// remoteCommandInfo.SshKey = sshKeyInfo.PrivateKey
	// remoteCommandInfo.UserName = sshKeyInfo.Username
	// PrivateKey  string `json:"privateKey"`
	// PublicKey   string `json:"publicKey"`
	// Username    string `json:"username"`

	// command는 bind 되어있을 것이고.
	vmCommand.McisID = mcisID
	vmCommand.VmID = vmID
	respMessage, respStatus := service.CommandVmOfMcis(defaultNameSpaceID, vmCommand)
	log.Println("CommandVmOfMcis result")
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
