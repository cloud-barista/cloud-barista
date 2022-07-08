package mcis

import (
	tbcommon "github.com/cloud-barista/cb-webtool/src/model/tumblebug/common"
)

type SpiderVMInfo struct {
	CspId            string                `json:"cspid"`
	IID              tbcommon.TbIID        `json:"iid"`
	ImageIID         tbcommon.TbIID        `json:"imageIId"`
	KeyPairIID       tbcommon.TbIID        `json:"keyPairIId"`
	KeyPairName      string                `json:"keyPairName"`
	KeyValueList     []tbcommon.TbKeyValue `json:"keyValueList"`
	Name             string                `json:"name"`
	NetworkInterface string                `json:"networkInterface"`

	PrivateDns         string           `json:"networkInterface"`
	PrivateIP          string           `json:"privateIP"`
	PublicDns          string           `json:"publicDns"`
	PublicIP           string           `json:"publicIP"`
	Region             RegionInfo       `json:"region"`
	SecurityGroupIIds  []tbcommon.TbIID `json:"securityGroupIIds"`
	SecurityGroupNames []string         `json:"securityGroupNames"`

	SshaccessPoint string `json:"sshaccessPoint"`
	StartTime      string `json:"startTime"`

	SubnetIID  tbcommon.TbIID `json:"subnetIID"`
	SubnetName string         `json:"subnetName"`

	VmblockDisk  string `json:"vmblockDisk"`
	VmbootDisk   string `json:"vmbootDisk"`
	VmspecName   string `json:"vmspecName"`
	VmuserId     string `json:"vmuserId"`
	VmuserPasswd string `json:"vmuserPasswd"`

	RootDeviceName string `json:"rootDeviceName"`
	RootDiskSize   string `json:"rootDiskSize"`
	RootDiskType   string `json:"rootDiskType"`

	VpcIID  tbcommon.TbIID `json:"vpcIID"`
	VpcName string         `json:"vpcName"`
}
