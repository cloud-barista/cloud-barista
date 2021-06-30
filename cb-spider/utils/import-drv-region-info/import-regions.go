// Cloud Control Manager's Rest Runtime of CB-Spider.
// The CB-Spider is a sub-Framework of the Cloud-Barista Multi-Cloud Project.
// The CB-Spider Mission is to connect all the clouds with a single interface.
//
//      * Cloud-Barista: https://github.com/cloud-barista
//
// by CB-Spider Team, 2020.05.

package main

import (
        "github.com/sirupsen/logrus"
        "github.com/cloud-barista/cb-store/config"

	icbs "github.com/cloud-barista/cb-store/interfaces"
	cim "github.com/cloud-barista/cb-spider/cloud-info-manager"
	rim "github.com/cloud-barista/cb-spider/cloud-info-manager/region-info-manager"

	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
	"fmt"
)


var cblog *logrus.Logger

func init() {
        cblog = config.Cblogger
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("================================================================")
		fmt.Println("    You can use this tool with Azure's Resource Group Name!!")
		fmt.Println("    Usage: import_regions.sh CB-GROUP-POWERKIM")
		fmt.Println("================================================================")
		return
	}
	args := os.Args[1:]

	resourceGroup := args[0]

	InsertRegionInfos(resourceGroup)

	//regionInfoList := InsertRegionInfos()
	//fmt.Println("%#v", regionInfoList)
}

// (1) get cloudos list
// (2) loop: 
// 		load RegionInfo List from all cloudos
// (3) insert
func InsertRegionInfos(resourceGroup string) ([]rim.RegionInfo) {

	// Set Environment Value of Project Root Path
	rootPath := os.Getenv("CBSPIDER_ROOT")
        if rootPath == "" {
                cblog.Error("$CBSPIDER_ROOT is not set!!")
                os.Exit(1)
        }
	regionInfoList := []rim.RegionInfo{}
	for _, cloudos := range cim.ListCloudOS() {
		regionFile, err := os.Open(rootPath + "/utils/import-drv-region-info/region-list/" + strings.ToLower(cloudos) + "-regions-list.json")
		if err != nil {
			if strings.Contains(err.Error(), "no such") {
				cblog.Info(err)
				continue
			} else {
				cblog.Error(err)
				continue
			}
		}
		defer regionFile.Close()

		var oneRegionInfoList []rim.RegionInfo 
		switch cloudos {
			case "AWS", "CLOUDIT", "OPENSTACK", "DOCKER":
				oneRegionInfoList = awsLoader(cloudos, regionFile)
			case "AZURE":
				oneRegionInfoList = azureLoader(cloudos, regionFile, resourceGroup)
			case "GCP":
				oneRegionInfoList = gcpLoader(cloudos, regionFile)
			case "ALIBABA":
				oneRegionInfoList = alibabaLoader(cloudos, regionFile)
			case "CLOUDTWIN":
			default:
				errmsg := cloudos + " is not a valid ProviderName!!"
				cblog.Error(errmsg)
				//return nil, fmt.Errorf(errmsg)
		}

		regionInfoList = append(regionInfoList, oneRegionInfoList...)
	}
	
	for _, regionInfo := range regionInfoList {
		_, err := rim.RegisterRegionInfo(regionInfo)
		if err != nil {
			cblog.Error(err)
		}
		//cblog.Info(fmt.Sprintf("###### Insert Region: %#v", regionInfo))
	}

	return regionInfoList
}


// for AWS, Cloudit, OpenStack, Docker
func awsLoader(cloudos string, regionFile *os.File) []rim.RegionInfo {

	type OrgRegions struct {
		Regions [] struct { 
			RegionName string 	`json:"RegionName"`
		}
	}

	byteValue, _ := ioutil.ReadAll(regionFile)

	var orgRegions OrgRegions
	json.Unmarshal(byteValue, &orgRegions)

	regionInfoList := []rim.RegionInfo{}
	for _, region := range orgRegions.Regions {
		keyValueList := []icbs.KeyValue{ {"Region", region.RegionName} }
		regionInfo := rim.RegionInfo{strings.ToLower(cloudos) + "-" + region.RegionName, 
						strings.ToUpper(cloudos), keyValueList}
		regionInfoList = append(regionInfoList, regionInfo)
	}

	return regionInfoList
}

func azureLoader(cloudos string, regionFile *os.File, resourceGroup string) []rim.RegionInfo {


	type Regions struct {
		Name string        `json:"name"`
	}

        byteValue, _ := ioutil.ReadAll(regionFile)

        var orgRegions []Regions
        json.Unmarshal(byteValue, &orgRegions)

        regionInfoList := []rim.RegionInfo{}
        for _, region := range orgRegions {

                keyValueList := []icbs.KeyValue{ {"location", region.Name}, {"ResourceGroup", resourceGroup} }
                regionInfo := rim.RegionInfo{strings.ToLower(cloudos) + "-" + region.Name,
                                                strings.ToUpper(cloudos), keyValueList}
                regionInfoList = append(regionInfoList, regionInfo)
        }

        return regionInfoList
}

func alibabaLoader(cloudos string, regionFile *os.File) []rim.RegionInfo {

        type OrgRegions struct {
                Regions struct {
			Region [] struct {
				RegionId string       `json:"RegionId"`
			}
                }
        }

        byteValue, _ := ioutil.ReadAll(regionFile)

        var orgRegions OrgRegions
        json.Unmarshal(byteValue, &orgRegions)

        regionInfoList := []rim.RegionInfo{}
        for _, region := range orgRegions.Regions.Region {
                keyValueList := []icbs.KeyValue{ {"Region", region.RegionId} , {"Zone", region.RegionId + "a"}}
                regionInfo := rim.RegionInfo{strings.ToLower(cloudos) + "-" + region.RegionId,
                                                strings.ToUpper(cloudos), keyValueList}
                regionInfoList = append(regionInfoList, regionInfo)
        }

        return regionInfoList
}

func gcpLoader(cloudos string, regionFile *os.File) []rim.RegionInfo {

        type Regions struct {
                Name string        `json:"name"`
        }

        byteValue, _ := ioutil.ReadAll(regionFile)

        var orgRegions []Regions
        json.Unmarshal(byteValue, &orgRegions)

        regionInfoList := []rim.RegionInfo{}
        for _, region := range orgRegions {

		runes := []rune(region.Name)
                keyValueList := []icbs.KeyValue{ {"Region", string(runes[0:len(region.Name)-2])}, {"Zone", region.Name} }
                regionInfo := rim.RegionInfo{strings.ToLower(cloudos) + "-" + region.Name,
                                                strings.ToUpper(cloudos), keyValueList}
                regionInfoList = append(regionInfoList, regionInfo)
        }

        return regionInfoList
}

