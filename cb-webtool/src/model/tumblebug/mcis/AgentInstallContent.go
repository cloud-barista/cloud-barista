package mcis

type AgentInstallContent struct {
	McisID string `json:"mcisId"`
	Result string `json:"result"`
	VmID   string `json:"vmId"`
	VmIp   string `json:"vmIp"`
}
