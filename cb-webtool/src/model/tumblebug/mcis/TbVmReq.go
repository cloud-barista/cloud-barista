package mcis

type TbVmReq struct {
	ConnectionName string `json:"connectionName"`
	Description    string `json:"description"`
	ImageID        string `json:"imageId"`
	Label          string `json:"label"`
	Name           string `json:"name"`

	SecurityGroupIDs []string `json:"securityGroupIds"`

	SpecID         string `json:"specId"`
	SshKeyID       string `json:"sshKeyId"`
	SubnetID       string `json:"subnetId"`
	VNetID         string `json:"vNetId"`
	VmGroupSize    string `json:"vmGroupSize"`
	VmUserAccount  string `json:"vmUserAccount"`
	VmUserPassword string `json:"vmUserPassword"`
}
