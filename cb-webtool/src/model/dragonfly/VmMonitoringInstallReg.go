package dragonfly

// Dragonfly에서 사용하는 AgentInstall
type VmMonitoringInstallReg struct {
	NameSpaceID string `json:"ns_id"`     // 네임스페이스 아이디
	McisID      string `json:"mcis_id"`   // MCIS 아이디
	VmID        string `json:"vm_id"`     // VM 아이디
	PublicIp    string `json:"public_ip"` // VM의 퍼블릭 아이피
	UserName    string `json:"user_name"` // VM의 SSH 접속 계정(cb-user)
	SshKey      string `json:"ssh_key"`   // VM의 SSH 접근 키
	SshKeyName  string `json:"ssh_key_name"`
	CspType     string `json:"cspType"` // VM의 CSP 정보 (aws ... 등)
	Port        string `json:"port"`
	//Port        int    `json:"port"`
}
