package mcis

type TbVmStatusInfo struct {
	CreatedTime string `json:"createdTime"`
	CspVmID     string `json:"cspVmId"`

	ID       string      `json:"id"`
	Location GeoLocation `json:"location"`

	MonAgentStatus string `json:"monAgentStatus"`

	Name          string `json:"name"`
	NativeStatus  string `json:"nativeStatus"`
	PrivateIP     string `json:"privateIP"`
	PublicIP      string `json:"publicIP"`
	SshPort       string `json:"sshPort"`
	Status        string `json:"status"`
	SystemMessage string `json:"systemMessage"`
	TargetAction  string `json:"targetAction"`
	TargetStatus  string `json:"targetStatus"`
}
