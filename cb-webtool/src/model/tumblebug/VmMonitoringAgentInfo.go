package tumblebug

type VmMonitoringAgentInfo struct {
	ConnectionName string         `json:"connectionName"`
	CspSshKeyName  string         `json:"cspSshKeyName"`
	Description    string         `json:"description"`
	Fingerprint    string         `json:"fingerprint"`
	ID             string         `json:"id"` // ? 무엇의 ID인가?
	Name           string         `json:"name"`
	PrivateKey     string         `json:"privateKey"`
	PublicKey      string         `json:"publicKey"`
	Nsername       string         `json:"username"`
	KeyValueInfos  []KeyValueInfo `json:"keyValueList"`
}
