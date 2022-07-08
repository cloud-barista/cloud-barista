package mcis

import (
	tbcommon "github.com/cloud-barista/cb-webtool/src/model/tumblebug/common"
)

type TbVmInfo struct {
	ConnectionName string `json:"connectionName"`
	CreatedTime    string `json:"createdTime"`

	CspViewVmDetail SpiderVMInfo `json:"cspViewVmDetail"`

	Description string      `json:"description"`
	ID          string      `json:"id"`
	IdByCSP     string      `json:"idByCSP"`
	ImageID     string      `json:"imageId"`
	Label       string      `json:"label"`
	Location    tbcommon.TbGeoLocation `json:"location"`

	MonAgentStatus string `json:"monAgentStatus"`

	Name                     string     `json:"name"`
	NetworkAgentStatus       string     `json:"networkAgentStatus"`
	PrivateDns               string     `json:"privateDns"`
	PrivateIP                string     `json:"privateIP"`
	PublicDNS                string     `json:"publicDNS"`
	PublicIP                 string     `json:"publicIP"`
	Region                   RegionInfo `json:"region"`

	RootDeviceName     string `json:"rootDeviceName"`
	RootDiskSize       string `json:"rootDiskSize"`
	RootDiskType       string `json:"rootDiskType"`

	SecurityGroupIDs []string `json:"securityGroupIds"`

	SpecID         string `json:"specId"`
	SshKeyID       string `json:"sshKeyId"`
	SshPort        string `json:"sshPort"`
	Status         string `json:"status"`
	SubnetID       string `json:"subnetId"`
	SystemMessage  string `json:"systemMessage"`
	TargetAction   string `json:"targetAction"`
	TargetStatus   string `json:"targetStatus"`
	VNetID         string `json:"vNetId"`
	VmBlockDisk    string `json:"vmBlockDisk"`
	VmBootDisk     string `json:"vmBootDisk"`
	VmGroupID      string `json:"vmGroupId"`
	VmUserAccount  string `json:"vmUserAccount"`
	VmUserPassword string `json:"vmUserPassword"`
}
