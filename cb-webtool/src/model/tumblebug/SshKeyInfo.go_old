package tumblebug

type SshKeyInfo struct {
	ConnectionName string `json:"connectionName"`
	CspSshKeyName  string `json:"cspSshKeyName"`

	Description string `json:"description"`
	Fingerprint string `json:"fingerprint"`
	ID          string `json:"id"`
	Name        string `json:"name"`
	PrivateKey  string `json:"privateKey"`
	PublicKey   string `json:"publicKey"`
	Username    string `json:"username"` //TODO : UserID 바꿔야 할텐데...

	KeyValueInfos []KeyValueInfo `json:"keyValueList"`
}

type SshKeyInfos []SshKeyInfo
