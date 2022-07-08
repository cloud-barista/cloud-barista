package mcis

import (
	tbcommon "github.com/cloud-barista/cb-webtool/src/model/tumblebug/common"
)

type TbVmStatusInfo struct {
	CreatedTime string `json:"createdTime"`
	CspVmID     string `json:"cspVmId"`

	ID       string      `json:"id"`
	Location tbcommon.TbGeoLocation `json:"location"`

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
