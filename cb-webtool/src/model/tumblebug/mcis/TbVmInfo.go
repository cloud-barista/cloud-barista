package mcis

type TbVmInfo struct {
	ConnectionName string `json:"connectionName"`
	CreatedTime    string `json:"createdTime"`

	CspViewVmDetail SpiderVMInfo `json:"cspViewVmDetail"`

	Description string      `json:"description"`
	ID          string      `json:"id"`
	ImageID     string      `json:"imageId"`
	Label       string      `json:"label"`
	Location    GeoLocation `json:"location"`

	MonAgentStatus string `json:"monAgentStatus"`

	Name       string     `json:"name"`
	PrivateDns string     `json:"privateDns"`
	PrivateIP  string     `json:"privateIP"`
	PublicDNS  string     `json:"publicDNS"`
	PublicIP   string     `json:"publicIP"`
	Region     RegionInfo `json:"region"`

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
