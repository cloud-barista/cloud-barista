package service

import (
	// "encoding/base64"
	"fmt"
	"log"

	// "log"
	// "io"
	// "net/http"

	echosession "github.com/go-session/echo-session"
	"github.com/labstack/echo"

	"github.com/cloud-barista/cb-webtool/src/model"
	"github.com/cloud-barista/cb-webtool/src/model/spider"
	// "github.com/cloud-barista/cb-webtool/src/model/tumblebug"
	tbcommon "github.com/cloud-barista/cb-webtool/src/model/tumblebug/common"
	// tbmcir "github.com/cloud-barista/cb-webtool/src/model/tumblebug/mcir"
	// tbmcis "github.com/cloud-barista/cb-webtool/src/model/tumblebug/mcis"
	util "github.com/cloud-barista/cb-webtool/src/util"
)

// 로그인할 때, NameSpace 저장(Create, Delete, Update) 외에는 이 funtion 사용
// 없으면 tb 조회
func GetStoredNameSpaceList(c echo.Context) ([]tbcommon.TbNsInfo, model.WebStatus) {
	fmt.Println("====== GET STORED NAME SPACE ========")
	nameSpaceList := []tbcommon.TbNsInfo{}
	nameSpaceErr := model.WebStatus{}
	store := echosession.FromContext(c)

	storedNameSpaceList, isExist := store.Get(util.STORE_NAMESPACELIST)
	if !isExist { // 존재하지 않으면 TB 조회
		nameSpaceList, nameSpaceErr = GetNameSpaceList()
		setError := SetStoreNameSpaceList(c, nameSpaceList)
		if setError != nil {
			log.Println("Set Namespace failed")
			nameSpaceErr.StatusCode = 4000
		}
	} else {
		log.Println(storedNameSpaceList)
		nameSpaceList = storedNameSpaceList.([]tbcommon.TbNsInfo)
		nameSpaceErr.StatusCode = 200
	}
	return nameSpaceList, nameSpaceErr
}

func SetStoreNameSpaceList(c echo.Context, nameSpaceList []tbcommon.TbNsInfo) error {
	fmt.Println("====== SET NAME SPACE ========")
	store := echosession.FromContext(c)
	store.Set(util.STORE_NAMESPACELIST, nameSpaceList)
	err := store.Save()
	return err
}

//GetCloudOSList
func GetStoredCloudOSList(c echo.Context) ([]string, model.WebStatus) {
	fmt.Println("====== GET STORED CloudOS ========")
	cloudOSList := []string{}
	cloudOsErr := model.WebStatus{}
	store := echosession.FromContext(c)

	storedCloudOSList, isExist := store.Get(util.STORE_CLOUDOSLIST)
	if !isExist { // 존재하지 않으면 TB 조회
		cloudOSList, cloudOsErr = GetCloudOSList()
		setError := SetStoreCloudOSList(c, cloudOSList)
		if setError != nil {
			log.Println("Set cloudOS failed")
		}
	} else {
		log.Println(storedCloudOSList)
		cloudOSList = storedCloudOSList.([]string)
		cloudOsErr.StatusCode = 200
	}
	return cloudOSList, cloudOsErr
}

func SetStoreCloudOSList(c echo.Context, cloudOSList []string) error {
	fmt.Println("====== SET cloudOS ========")
	store := echosession.FromContext(c)
	store.Set(util.STORE_CLOUDOSLIST, cloudOSList)
	err := store.Save()
	return err
}

//GetRegionList
func GetStoredRegionList(c echo.Context) ([]spider.RegionInfo, model.WebStatus) {
	fmt.Println("====== GET STORED Region ========")
	regionList := []spider.RegionInfo{}
	regionErr := model.WebStatus{}
	store := echosession.FromContext(c)

	storedRegionList, isExist := store.Get(util.STORE_REGIONLIST)
	if !isExist { // 존재하지 않으면 TB 조회
		regionList, regionErr = GetRegionList()
		setError := SetStoreRegionList(c, regionList)
		if setError != nil {
			log.Println("Set Region failed")
		}
	} else {
		log.Println(storedRegionList)
		regionList = storedRegionList.([]spider.RegionInfo)
		regionErr.StatusCode = 200
	}
	return regionList, regionErr
}

func SetStoreRegionList(c echo.Context, regionList []spider.RegionInfo) error {
	fmt.Println("====== SET Region ========")
	store := echosession.FromContext(c)
	store.Set(util.STORE_REGIONLIST, regionList)
	err := store.Save()
	return err
}

// GetCredentialList
func GetStoredCredentialList(c echo.Context) ([]spider.CredentialInfo, model.WebStatus) {
	fmt.Println("====== GET STORED Region ========")
	credentialList := []spider.CredentialInfo{}
	credentialErr := model.WebStatus{}
	store := echosession.FromContext(c)

	storedCredentialList, isExist := store.Get(util.STORE_CREDENTIALLIST)
	if !isExist { // 존재하지 않으면 TB 조회
		credentialList, credentialErr = GetCredentialList()
		setError := SetStoreCredentialList(c, credentialList)
		if setError != nil {
			log.Println("Set Credential failed")
		}
	} else {
		log.Println(storedCredentialList)
		credentialList = storedCredentialList.([]spider.CredentialInfo)
		credentialErr.StatusCode = 200
	}
	return credentialList, credentialErr
}

func SetStoreCredentialList(c echo.Context, credentialList []spider.CredentialInfo) error {
	fmt.Println("====== SET Credential ========")
	store := echosession.FromContext(c)
	store.Set(util.STORE_CREDENTIALLIST, credentialList)
	err := store.Save()
	return err
}

// GetDriverList
func GetStoredDriverList(c echo.Context) ([]spider.DriverInfo, model.WebStatus) {
	fmt.Println("====== GET STORED Driver ========")
	driverList := []spider.DriverInfo{}
	driverErr := model.WebStatus{}
	store := echosession.FromContext(c)

	storedDriverList, isExist := store.Get(util.STORE_DRIVERLIST)
	if !isExist { // 존재하지 않으면 TB 조회
		driverList, driverErr = GetDriverList()
		setError := SetStoreDriverList(c, driverList)
		if setError != nil {
			log.Println("Set Driver failed")
		}
	} else {
		log.Println(storedDriverList)
		driverList = storedDriverList.([]spider.DriverInfo)
		driverErr.StatusCode = 200
	}
	return driverList, driverErr
}

func SetStoreDriverList(c echo.Context, driverList []spider.DriverInfo) error {
	fmt.Println("====== SET Driver ========")
	store := echosession.FromContext(c)
	store.Set(util.STORE_DRIVERLIST, driverList)
	err := store.Save()
	return err
}

//GetCloudConnectionConfigList
func GetStoredCloudConnectionConfigList(c echo.Context) ([]spider.CloudConnectionConfigInfo, model.WebStatus) {
	fmt.Println("====== GET STORED CloudConnectionConfigList ========")
	connectionConfigList := []spider.CloudConnectionConfigInfo{}
	connectionConfigErr := model.WebStatus{}
	store := echosession.FromContext(c)

	storedConnectionConfigList, isExist := store.Get(util.STORE_CLOUDCONNECTIONCONFIGLIST)
	if !isExist { // 존재하지 않으면 TB 조회
		connectionConfigList, connectionConfigErr = GetCloudConnectionConfigList()
		setError := SetStoreCloudConnectionConfigList(c, connectionConfigList)
		if setError != nil {
			log.Println("Set ConnectionConfigList failed")
		}
	} else {
		log.Println(storedConnectionConfigList)
		connectionConfigList = storedConnectionConfigList.([]spider.CloudConnectionConfigInfo)
		connectionConfigErr.StatusCode = 200
	}
	return connectionConfigList, connectionConfigErr
}

func SetStoreCloudConnectionConfigList(c echo.Context, connectionConfigList []spider.CloudConnectionConfigInfo) error {
	fmt.Println("====== SET CloudConnectionConfigList ========")
	store := echosession.FromContext(c)
	store.Set(util.STORE_CLOUDCONNECTIONCONFIGLIST, connectionConfigList)
	err := store.Save()
	return err
}

// move to NameSpaceController.go
// func GetNameSpace(c echo.Context) error {
// 	fmt.Println("====== GET NAME SPACE ========")
// 	store := echosession.FromContext(c)

// 	getInfo, ok := store.Get("namespace")
// 	if !ok {
// 		return c.JSON(http.StatusNotAcceptable, map[string]string{
// 			"message": "Not Exist",
// 		})
// 	}
// 	nsId := getInfo.(string)

// 	res := map[string]string{
// 		"message": "success",
// 		"nsID":    nsId,
// 	}

// 	return c.JSON(http.StatusOK, res)
// }
