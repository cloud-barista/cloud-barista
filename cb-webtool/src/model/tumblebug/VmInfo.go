package tumblebug

import (
	"time"
)

type VmInfo struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	ConnectionName string `json:"connectionName"`

	SpecID  string `json:"specId"`
	ImageID string `json:"imageId"`

	VNetID           string   `json:"vNetId"`
	SubnetID         string   `json:"subnetId"`
	SecurityGroupIDs []string `json:"securityGroupIds"`

	SshKeyID string `json:"sshKeyId"`

	VmUserAccount  string `json:"vmUserAccount"`
	VmUserPassword string `json:"vmUserPassword"`

	Description string `json:"description"`
	Label       string `json:"label"`

	Location   LocationInfo   `json:"location"`
	RegionZone RegionZoneInfo `json:"region"` // tumblebug에서 결과에 region에 region/zone을 같이보내는가?

	PublicIP   string `json:"publicIP"`
	PublicDNS  string `json:"publicDNS"`
	PrivateIP  string `json:"privateIP"`
	PrivateDNS string `json:"privateDNS"`

	VmBootDisk  string `json:"vmBootDisk"`
	VmBlockDisk string `json:"vmBlockDisk"`

	Status       string `json:"status"`
	TargetStatus string `json:"targetStatus"`
	TargetAction string `json:"targetAction"`

	MonAgentStatus string `json:"monAgentStatus"` // "monAgentStatus": "[installed, notInstalled, failed]",

	CspViewVmDetail CspViewVmDetailInfo `json:"cspViewVmDetail"`
}

type CspViewVmDetailInfo struct {
	Name               string   `json:"name"`
	ImageName          string   `json:"imageName"`
	Vpcname            string   `json:"vpcname"`
	SubnetName         string   `json:"subnetName"`
	SecurityGroupNames []string `json:"securityGroupNames"`

	KeyPairName string `json:"keyPairName"`

	VmspecName string `json:"vmspecName"`

	VmuserID     string `json:"vmuserId"`
	VmuserPasswd string `json:"vmuserPasswd"`

	ConnectionName    string             `json:"connectionName"`
	IID               NameSystemIdInfo   `json:"iid"`
	ImageIID          NameSystemIdInfo   `json:"imageIId"`
	VpcIID            NameSystemIdInfo   `json:"vpcIID"`
	SubnetIID         NameSystemIdInfo   `json:"subnetIID"`
	SecurityGroupIIDs []NameSystemIdInfo `json:"securityGroupIIds"`
	KeyPairIID        NameSystemIdInfo   `json:"keyPairIId"`

	StartTime  time.Time      `json:"startTime"`
	RegionZone RegionZoneInfo `json:"region"`

	NetworkInterface string `json:"networkInterface"`

	PublicIP   string `json:"publicIP"`
	PublicDNS  string `json:"publicDNS"`
	PrivateIP  string `json:"privateIP"`
	PrivateDNS string `json:"privateDNS"`

	VmbootDisk  string `json:"vmbootDisk"`
	VmblockDisk string `json:"vmblockDisk"`

	KeyValueInfos []KeyValueInfo `json:"keyValueList"`
}

type LocationInfo struct {
	BriefAddr string `json:"briefAddr"`
	CloudType string `json:"cloudType"`
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
	// Latitude     float64 `json:"latitude"`
	// Longitude    float64 `json:"longitude"`
	NativeRegion string `json:"nativeRegion"`
}

//
type NameSystemIdInfo struct {
	NameID   string `json:"nameId"`
	SystemID string `json:"systemId"`
}

type RegionZoneInfo struct {
	Region string `json:"Region"`
	Zone   string `json:"Zone"`
}
