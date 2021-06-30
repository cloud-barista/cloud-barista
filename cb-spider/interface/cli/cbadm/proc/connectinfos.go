// gRPC Runtime of CB-Spider.
// The CB-Spider is a sub-Framework of the Cloud-Barista Multi-Cloud Project.
// The CB-Spider Mission is to connect all the clouds with a single interface.
//
//      * Cloud-Barista: https://github.com/cloud-barista
//
// by CB-Spider Team, 2020.09.

package proc

import (
	"encoding/json"
	"fmt"

	gc "github.com/cloud-barista/cb-spider/api-runtime/grpc-runtime/common"
	"github.com/cloud-barista/cb-spider/interface/api"
)

// ===== [ Constants and Variables ] =====

const (
	// ConfigVersion - 설정 구조에 대한 버전
	ConfigVersion = 1
)

// ===== [ Types ] =====

// ConnectInfosConfig -
type ConnectInfosConfig struct {
	Version         int           `yaml:"Version" json:"Version"`
	ConnectInfoList []ConnectInfo `yaml:"ConnectInfos" json:"ConnectInfos"`
}

// ConnectInfo -
type ConnectInfo struct {
	ConfigName   string         `yaml:"ConfigName" json:"ConfigName"`
	ProviderName string         `yaml:"ProviderName" json:"ProviderName"`
	Driver       DriverInfo     `yaml:"Driver" json:"Driver"`
	Credential   CredentialInfo `yaml:"Credential" json:"Credential"`
	Region       RegionInfo     `yaml:"Region" json:"Region"`
}

// DriverInfo -
type DriverInfo struct {
	DriverName        string `yaml:"DriverName" json:"DriverName"`
	DriverLibFileName string `yaml:"DriverLibFileName" json:"DriverLibFileName"`
}

// CredentialInfo -
type CredentialInfo struct {
	CredentialName   string         `yaml:"CredentialName" json:"CredentialName"`
	KeyValueInfoList []KeyValueInfo `yaml:"KeyValueInfoList" json:"KeyValueInfoList"`
}

// RegionInfo -
type RegionInfo struct {
	RegionName       string         `yaml:"RegionName" json:"RegionName"`
	KeyValueInfoList []KeyValueInfo `yaml:"KeyValueInfoList" json:"KeyValueInfoList"`
}

// KeyValueInfo -
type KeyValueInfo struct {
	Key   string `yaml:"Key" json:"Key"`
	Value string `yaml:"Value" json:"Value"`
}

// ===== [ Implementations ] =====

// ===== [ Private Functions ] =====

// ===== [ Public Functions ] =====

// ListConnectInfos - 연결정보 목록 통합 제공
func ListConnectInfos(cim *api.CIMApi) (string, error) {

	holdType, _ := cim.GetOutType()
	cim.SetOutType("json")
	defer cim.SetOutType(holdType)

	result, err := cim.ListConnectionConfig()
	if err != nil {
		return "", err
	}

	connectConfigList := make(map[string]interface{})
	err = json.Unmarshal([]byte(result), &connectConfigList)
	if err != nil {
		return "", err
	}

	connectInfoList := []ConnectInfo{}
	if connectConfigList["connectionconfig"] != nil {
		for _, m := range connectConfigList["connectionconfig"].([]interface{}) {

			connectConfig := m.(map[string]interface{})

			connectInfo := ConnectInfo{}
			connectInfo.ConfigName = fmt.Sprintf("%v", connectConfig["ConfigName"])
			connectInfo.ProviderName = fmt.Sprintf("%v", connectConfig["ProviderName"])

			result, err := cim.GetCloudDriverByParam(fmt.Sprintf("%v", connectConfig["DriverName"]))
			if err != nil {
				return "", err
			}

			driverItem := make(map[string]interface{})
			err = json.Unmarshal([]byte(result), &driverItem)
			if err != nil {
				return "", err
			}

			connectInfo.Driver.DriverName = fmt.Sprintf("%v", driverItem["DriverName"])
			connectInfo.Driver.DriverLibFileName = fmt.Sprintf("%v", driverItem["DriverLibFileName"])

			result, err = cim.GetCredentialByParam(fmt.Sprintf("%v", connectConfig["CredentialName"]))
			if err != nil {
				return "", err
			}

			credentialItem := make(map[string]interface{})
			err = json.Unmarshal([]byte(result), &credentialItem)
			if err != nil {
				return "", err
			}

			connectInfo.Credential.CredentialName = fmt.Sprintf("%v", credentialItem["CredentialName"])
			err = gc.CopySrcToDest(credentialItem["KeyValueInfoList"], &connectInfo.Credential.KeyValueInfoList)
			if err != nil {
				return "", err
			}

			result, err = cim.GetRegionByParam(fmt.Sprintf("%v", connectConfig["RegionName"]))
			if err != nil {
				return "", err
			}

			regionItem := make(map[string]interface{})
			err = json.Unmarshal([]byte(result), &regionItem)
			if err != nil {
				return "", err
			}

			connectInfo.Region.RegionName = fmt.Sprintf("%v", regionItem["RegionName"])
			err = gc.CopySrcToDest(regionItem["KeyValueInfoList"], &connectInfo.Region.KeyValueInfoList)
			if err != nil {
				return "", err
			}

			connectInfoList = append(connectInfoList, connectInfo)
		}
	}

	var cfg ConnectInfosConfig
	cfg.Version = ConfigVersion
	cfg.ConnectInfoList = connectInfoList

	return gc.ConvertToOutput(holdType, &cfg)
}

// GetConnectInfos - 연결정보 통합 제공
func GetConnectInfos(cim *api.CIMApi, configName string) (string, error) {

	holdType, _ := cim.GetOutType()
	cim.SetOutType("json")
	defer cim.SetOutType(holdType)

	result, err := cim.GetConnectionConfigByParam(configName)
	if err != nil {
		return "", err
	}

	connectConfig := make(map[string]interface{})
	err = json.Unmarshal([]byte(result), &connectConfig)
	if err != nil {
		return "", err
	}

	connectInfoList := []ConnectInfo{}

	connectInfo := ConnectInfo{}
	connectInfo.ConfigName = fmt.Sprintf("%v", connectConfig["ConfigName"])
	connectInfo.ProviderName = fmt.Sprintf("%v", connectConfig["ProviderName"])

	result, err = cim.GetCloudDriverByParam(fmt.Sprintf("%v", connectConfig["DriverName"]))
	if err != nil {
		return "", err
	}

	driverItem := make(map[string]interface{})
	err = json.Unmarshal([]byte(result), &driverItem)
	if err != nil {
		return "", err
	}

	connectInfo.Driver.DriverName = fmt.Sprintf("%v", driverItem["DriverName"])
	connectInfo.Driver.DriverLibFileName = fmt.Sprintf("%v", driverItem["DriverLibFileName"])

	result, err = cim.GetCredentialByParam(fmt.Sprintf("%v", connectConfig["CredentialName"]))
	if err != nil {
		return "", err
	}

	credentialItem := make(map[string]interface{})
	err = json.Unmarshal([]byte(result), &credentialItem)
	if err != nil {
		return "", err
	}

	connectInfo.Credential.CredentialName = fmt.Sprintf("%v", credentialItem["CredentialName"])
	err = gc.CopySrcToDest(credentialItem["KeyValueInfoList"], &connectInfo.Credential.KeyValueInfoList)
	if err != nil {
		return "", err
	}

	result, err = cim.GetRegionByParam(fmt.Sprintf("%v", connectConfig["RegionName"]))
	if err != nil {
		return "", err
	}

	regionItem := make(map[string]interface{})
	err = json.Unmarshal([]byte(result), &regionItem)
	if err != nil {
		return "", err
	}

	connectInfo.Region.RegionName = fmt.Sprintf("%v", regionItem["RegionName"])
	err = gc.CopySrcToDest(regionItem["KeyValueInfoList"], &connectInfo.Region.KeyValueInfoList)
	if err != nil {
		return "", err
	}

	connectInfoList = append(connectInfoList, connectInfo)

	var cfg ConnectInfosConfig
	cfg.Version = ConfigVersion
	cfg.ConnectInfoList = connectInfoList

	return gc.ConvertToOutput(holdType, &cfg)
}
