package tumblebug

type McisCommandInfo struct {
	Command  string `json:"command"`
	UserName string `json:"userName"`
	Ip       string `json:"ip"`
	SshKey   string `json:"sskKey"`
	McisID   string `json:"mcisId"`
	VmID     string `json:"vmId"`
}
