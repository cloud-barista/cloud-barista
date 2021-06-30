package controller

import (
	// "encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/cloud-barista/cb-webtool/src/model/ladybug"
	// "github.com/cloud-barista/cb-webtool/src/model/dragonfly"
	// "github.com/cloud-barista/cb-webtool/src/model/tumblebug"

	service "github.com/cloud-barista/cb-webtool/src/service"
	util "github.com/cloud-barista/cb-webtool/src/util"

	echotemplate "github.com/foolin/echo-template"
	// echosession "github.com/go-session/echo-session"
	"github.com/labstack/echo"
	// echosession "github.com/go-session/echo-session"
)

func McksRegForm(c echo.Context) error {
	fmt.Println("McksRegForm ************ : ")

	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}
	defaultNameSpaceID := loginInfo.DefaultNameSpaceID

	// 최신 namespacelist 가져오기
	nsList, _ := service.GetNameSpaceList()
	log.Println(" nsList  ", nsList)

	// connectionconfigList 가져오기
	cloudOsList, _ := service.GetCloudOSList()
	log.Println(" cloudOsList  ", cloudOsList)

	// regionList 가져오기
	regionList, _ := service.GetRegionList()
	log.Println(" regionList  ", regionList)

	cloudConnectionConfigInfoList, _ := service.GetCloudConnectionConfigList() // 등록된 모든 connection 정보
	log.Println("---------------------- GetCloudConnectionConfigList ", defaultNameSpaceID)

	clusterList, _ := service.GetClusterList(defaultNameSpaceID)

	return echotemplate.Render(c, http.StatusOK,
		"operation/manages/mcksmng/McksCreate", // 파일명
		map[string]interface{}{
			"LoginInfo":          loginInfo,
			"DefaultNameSpaceID": defaultNameSpaceID,
			"NameSpaceList":      nsList,
			"CloudOSList":        cloudOsList,
			"RegionList":         regionList,

			"CloudConnectionConfigInfoList": cloudConnectionConfigInfoList,
			"ClusterList":                   clusterList,
		})
}

func McksMngForm(c echo.Context) error {
	fmt.Println("McksMngForm ************ : ")

	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	defaultNameSpaceID := loginInfo.DefaultNameSpaceID
	
	mcksSimpleInfoList := []ladybug.ClusterSimpleInfo{} // 표에 뿌려줄 정보
	totalMcksStatusCountMap := make(map[string]int)   
	totalClusterCount := 0;

	// 최신 namespacelist 가져오기
	nsList, _ := service.GetNameSpaceList()
	log.Println(" nsList  ", nsList)

	// provider 별 연결정보 count
	cloudConnectionConfigInfoList, _ := service.GetCloudConnectionConfigList()
	connectionConfigCountMap, providerCount := service.GetCloudConnectionCountMap(cloudConnectionConfigInfoList)
	totalConnectionCount := len(cloudConnectionConfigInfoList)


	// 모든 MCKS 조회
	clusterList, clusterErr := service.GetClusterList(defaultNameSpaceID)
	if clusterErr.StatusCode != 200 && clusterErr.StatusCode != 201 {
		echotemplate.Render(c, http.StatusOK,
			"operation/manages/mcksmng/McksMng", // 파일명
			map[string]interface{}{
				"Message": clusterErr.Message,
				"Status":  clusterErr.StatusCode,
				"LoginInfo":          loginInfo,
				"DefaultNameSpaceID": defaultNameSpaceID,
				"NameSpaceList":      nsList,

				// cp count 영역
				"TotalProviderCount":         providerCount,
				"TotalConnectionConfigCount": totalConnectionCount,     // 총 connection 갯수
				"ConnectionConfigCountMap":   connectionConfigCountMap, // provider별 connection 수

				// "ClusterList": clusterList,
				"ClusterList":             mcksSimpleInfoList,
				"TotalMcksStatusCountMap": totalMcksStatusCountMap,
				"TotalClusterCount":       totalClusterCount,
			})
	}

	totalClusterCount = len(clusterList)
	if totalClusterCount == 0 {
		return c.Redirect(http.StatusTemporaryRedirect, "/operation/manages/mcksmng/regform")
	}

	totalMcksStatusCountMap = service.GetMcksStatusCountMap(clusterList)
	////////////// return value 에 set
	
	for _, mcksInfo := range clusterList {
		mcksSimpleInfo := ladybug.ClusterSimpleInfo{}
		mcksSimpleInfo.UID = mcksInfo.UID
		mcksSimpleInfo.Status = mcksInfo.Status
		mcksSimpleInfo.McksStatus = util.GetMcksStatus(mcksInfo.Status)
		mcksSimpleInfo.Name = mcksInfo.Name
		mcksSimpleInfo.ClusterConfig = mcksInfo.ClusterConfig
		mcksSimpleInfo.CpLeader = mcksInfo.CpLeader
		mcksSimpleInfo.Kind = mcksInfo.Kind
		mcksSimpleInfo.Mcis = mcksInfo.Mcis
		mcksSimpleInfo.NameSpace = mcksInfo.NameSpace
		mcksSimpleInfo.NetworkCni = mcksInfo.NetworkCni

		resultSimpleNodeList, resultSimpleNodeRoleCountMap := service.GetSimpleNodeCountMap(mcksInfo)

		mcksSimpleInfo.Nodes = resultSimpleNodeList
		mcksSimpleInfo.TotalNodeCount = len(resultSimpleNodeList) // 해당 mcks의 모든 node 갯수
		//nodeKindCountMapByMcks[mcksInfo.UID] = resultSimpleNodeKindCountMap // MCIS 내 vm 상태별 cnt
		mcksSimpleInfo.NodeCountMap = resultSimpleNodeRoleCountMap // MCKS UID 별 KindCountMap
		// mcksSimpleInfo.NodeSimpleList = resultSimpleNodeList

		// log.Println("**************")
		// mapValues, _ := util.StructToMapByJson(mcksSimpleInfo)
		// log.Println(mapValues)
		// log.Println("**************")

		mcksSimpleInfoList = append(mcksSimpleInfoList, mcksSimpleInfo)
	}

	// status, filepath, return params
	return echotemplate.Render(c, http.StatusOK,
		"operation/manages/mcksmng/McksMng", // 파일명
		map[string]interface{}{
			"Message": clusterErr.Message,
			"Status":  clusterErr.StatusCode,
			"LoginInfo":          loginInfo,
			"DefaultNameSpaceID": defaultNameSpaceID,
			"NameSpaceList":      nsList,

			// cp count 영역
			"TotalProviderCount":         providerCount,
			"TotalConnectionConfigCount": totalConnectionCount,     // 총 connection 갯수
			"ConnectionConfigCountMap":   connectionConfigCountMap, // provider별 connection 수

			// "ClusterList": clusterList,
			"ClusterList":             mcksSimpleInfoList,
			"TotalMcksStatusCountMap": totalMcksStatusCountMap,
			"TotalClusterCount":       totalClusterCount,
		})
}

// MCKS 목록 조회
func GetMcksList(c echo.Context) error {
	log.Println("GetMcksList : ")
	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	// store := echosession.FromContext(c)
	defaultNameSpaceID := loginInfo.DefaultNameSpaceID
	// TODO : defaultNameSpaceID 가 없으면 설정화면으로 보낼 것
	mcksList, respStatus := service.GetClusterList(defaultNameSpaceID)
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
		"McksList":           mcksList,
	})
}

// Cluster 등록 처리
func McksRegProc(c echo.Context) error {
	log.Println("McksRegProc : ")
	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	log.Println("get info")
	//&[]Person{}
	clusterReq := &ladybug.ClusterRegReq{}
	if err := c.Bind(clusterReq); err != nil {
		// if err := c.Bind(mCISInfoList); err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "fail",
			"status":  "fail",
		})
	}
	log.Println(clusterReq)

	defaultNameSpaceID := loginInfo.DefaultNameSpaceID
	// TODO : defaultNameSpaceID 가 없으면 설정화면으로 보낼 것
	clusterInfo, respStatus := service.RegCluster(defaultNameSpaceID, clusterReq)
	log.Println("RegMcis service returned")
	if respStatus.StatusCode != 200 && respStatus.StatusCode != 201 {
		return c.JSON(respStatus.StatusCode, map[string]interface{}{
			"error":  respStatus.Message,
			"status": respStatus.StatusCode,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":     "success",
		"status":      respStatus.StatusCode,
		"ClusterInfo": clusterInfo,
	})
}

// MCKS 삭제처리
func McksDelProc(c echo.Context) error {
	log.Println("McksDelProc : ")
	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	defaultNameSpaceID := loginInfo.DefaultNameSpaceID
	// TODO : defaultNameSpaceID 가 없으면 설정화면으로 보낼 것

	//clusteruID := c.Param("clusteruID")
	clusterName := c.Param("clusterName")
	log.Println("clusterName= " + clusterName)

	resultStatusInfo, respStatus := service.DelCluster(defaultNameSpaceID, clusterName)
	log.Println("DelMCKS service returned")
	if respStatus.StatusCode != 200 && respStatus.StatusCode != 201 {
		return c.JSON(respStatus.StatusCode, map[string]interface{}{
			"error":  respStatus.Message,
			"status": respStatus.StatusCode,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":    "success",
		"status":     respStatus.StatusCode,
		"StatusInfo": resultStatusInfo,
	})
}

// Node 등록 form
func McksNodeRegForm(c echo.Context) error {
	fmt.Println("McksNodeRegForm ************ : ")

	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}
	defaultNameSpaceID := loginInfo.DefaultNameSpaceID

	clusterUID := c.Param("clusterUID")
	clusterName := c.Param("clusterName")

	// 최신 namespacelist 가져오기
	nsList, _ := service.GetNameSpaceList()
	log.Println(" nsList  ", nsList)

	// connectionconfigList 가져오기
	cloudOsList, _ := service.GetCloudOSList()
	log.Println(" cloudOsList  ", cloudOsList)

	// regionList 가져오기
	regionList, _ := service.GetRegionList()
	log.Println(" regionList  ", regionList)

	cloudConnectionConfigInfoList, _ := service.GetCloudConnectionConfigList() // 등록된 모든 connection 정보
	log.Println("---------------------- GetCloudConnectionConfigList ", defaultNameSpaceID)

	nodeList, _ := service.GetNodeList(defaultNameSpaceID, clusterName)
	nodeListLength := len(nodeList.Items)
	log.Println("---------------------- nodeListLength ", nodeListLength)

	return echotemplate.Render(c, http.StatusOK,
		"operation/manages/mcksmng/NodeCreate", // 파일명
		map[string]interface{}{
			"LoginInfo":          loginInfo,
			"DefaultNameSpaceID": defaultNameSpaceID,
			"NameSpaceList":      nsList,
			"CloudOSList":        cloudOsList,
			"RegionList":         regionList,

			"CloudConnectionConfigInfoList": cloudConnectionConfigInfoList,
			"NodeList":                      nodeList,
			"McksID":                        clusterUID,
			"McksName":                      clusterName,
		})
}

// Node 등록 처리
func NodeRegProc(c echo.Context) error {
	log.Println("NodeRegProc : ")
	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	clusteruID := c.Param("clusteruID")
	clusterName := c.Param("clusterName")

	nodeRegReq := &ladybug.NodeRegReq{}
	// nodeRegReq := &ladybug.NodeReq{}
	if err := c.Bind(nodeRegReq); err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "fail",
			"status":  "fail",
		})
	}
	log.Println(nodeRegReq)

	defaultNameSpaceID := loginInfo.DefaultNameSpaceID

	nodeInfo, respStatus := service.RegNode(defaultNameSpaceID, clusterName, nodeRegReq)
	log.Println("RegMcis service returned")
	if respStatus.StatusCode != 200 && respStatus.StatusCode != 201 {
		return c.JSON(respStatus.StatusCode, map[string]interface{}{
			"error":  respStatus.Message,
			"status": respStatus.StatusCode,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":    "success",
		"status":     respStatus.StatusCode,
		"ClusteruID": clusteruID,
		"NodeInfo":   nodeInfo,
	})
}

// Node 삭제 처리
func NodeDelProc(c echo.Context) error {
	log.Println("NodeRegProc : ")
	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	//clusteruID := c.Param("clusteruID")
	clusterName := c.Param("clusterName")
	//nodeID := c.Param("nodeID")
	nodeName := c.Param("nodeName")

	defaultNameSpaceID := loginInfo.DefaultNameSpaceID

	resultStatusInfo, respStatus := service.DelNode(defaultNameSpaceID, clusterName, nodeName)
	log.Println("DelMCKS service returned")
	if respStatus.StatusCode != 200 && respStatus.StatusCode != 201 {
		return c.JSON(respStatus.StatusCode, map[string]interface{}{
			"error":  respStatus.Message,
			"status": respStatus.StatusCode,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":    "success",
		"status":     respStatus.StatusCode,
		"StatusInfo": resultStatusInfo,
	})
}
