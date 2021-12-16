// Cloud regioninfomanager Info. Manager of CB-Spider.
// The CB-Spider is a sub-Framework of the Cloud-Barista Multi-Cloud Project.
// The CB-Spider Mission is to connect all the clouds with a single interface.
//
//      * Cloud-Barista: https://github.com/cloud-barista
//
// by CB-Spider Team, 2019.09.

package regioninfomanager

import (
	"fmt"
	"strings"
	"github.com/cloud-barista/cb-store/config"
	icbs "github.com/cloud-barista/cb-store/interfaces"
	cim "github.com/cloud-barista/cb-spider/cloud-info-manager"

	"github.com/sirupsen/logrus"
)

var cblog *logrus.Logger

func init() {
	cblog = config.Cblogger
}

//====================================================================
type RegionInfo struct {
	RegionName       string          // ex) "region01"
	ProviderName     string          // ex) "GCP"
	KeyValueInfoList []icbs.KeyValue // ex) { {region, us-east1},
	//	 {zone, us-east1-c},
}

//====================================================================

func RegisterRegionInfo(rgnInfo RegionInfo) (*RegionInfo, error) {
	return RegisterRegion(rgnInfo.RegionName, rgnInfo.ProviderName, rgnInfo.KeyValueInfoList)
}

// 1. check params
// 2. insert them into cb-store
func RegisterRegion(regionName string, providerName string, keyValueInfoList []icbs.KeyValue) (*RegionInfo, error) {
	cblog.Info("call RegisterRegion()")

	cblog.Debug("check params")
	err := checkParams(regionName, providerName, keyValueInfoList)
	if err != nil {
		return nil, err

	}

        // trim user inputs
        regionName = strings.TrimSpace(regionName)
	providerName = strings.ToUpper(strings.TrimSpace(providerName))

	cblog.Debug("insert metainfo into store")

	err = insertInfo(regionName, providerName, keyValueInfoList)
	if err != nil {
		cblog.Error(err)
		return nil, err
	}

	rgnInfo := &RegionInfo{regionName, providerName, keyValueInfoList}
	return rgnInfo, nil
}

func ListRegion() ([]*RegionInfo, error) {
	cblog.Info("call ListRegion()")

	regionInfoList, err := listInfo()
	if err != nil {
		return nil, err
	}

	return regionInfoList, nil
}

// 1. check params
// 2. get CredentialInfo from cb-store
func GetRegion(regionName string) (*RegionInfo, error) {
	cblog.Info("call GetRegion()")

	if regionName == "" {
		return nil, fmt.Errorf("RegionName is empty!")
	}

	rgnInfo, err := getInfo(regionName)
	if err != nil {
		cblog.Error(err)
		return nil, err
	}

	return rgnInfo, err
}

func UnRegisterRegion(regionName string) (bool, error) {
	cblog.Info("call UnRegisterRegion()")

	if regionName == "" {
		return false, fmt.Errorf("RegionName is empty!")
	}

	result, err := deleteInfo(regionName)
	if err != nil {
		cblog.Error(err)
		return false, err
	}

	return result, nil
}

//----------------

func checkParams(regionName string, providerName string, keyValueInfoList []icbs.KeyValue) error {
	if regionName == "" {
		return fmt.Errorf("RegionName is empty!")
	}
	if providerName == "" {
		return fmt.Errorf("ProviderName is empty!")
	}
	if keyValueInfoList == nil {
		return fmt.Errorf("KeyValue List is nil!")
	}

	// get Provider's Meta Info
        cloudOSMetaInfo, err := cim.GetCloudOSMetaInfo(providerName)
        if err != nil {
                cblog.Error(err)
                return err
        }

	// validate the KeyValueList of Region Input
	err = cim.ValidateKeyValueList(keyValueInfoList, cloudOSMetaInfo.Region)
        if err != nil {
                cblog.Error(err)
                return err
        }

	return nil
}
