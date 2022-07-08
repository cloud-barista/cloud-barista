package mcir

type TbSshKeyReq struct {
	ConnectionName string `json:"connectionName"`
	CspSshKeyId    string `json:"cspSshKeyId"`
	Description    string `json:"description"`
	Fingerprint    string `json:"fingerprint"`
	Name           string `json:"name"`
	PrivateKey     string `json:"privateKey"`
	PublicKey     string `json:"publicKey"`
	Username       string `json:"username"`
	VerifiedUsername string `json:"verifiedUsername"`
}
