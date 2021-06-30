package spider

// ConnectionConfigData -> CloudConnectionConfigInfo 로 변경
type CloudConnectionConfigInfo struct {
	//user(username, password, email)
	ConfigName     string `json:"ConfigName"`
	ProviderName   string `json:"ProviderName"`
	DriverName     string `json:"DriverName"`
	CredentialName string `json:"CredentialName"`
	RegionName     string `json:"RegionName"`
}
type CloudConnectionConfigInfos []CloudConnectionConfigInfo
