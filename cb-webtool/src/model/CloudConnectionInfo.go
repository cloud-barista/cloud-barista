package model

// ConnectionConfigData와 뭐가 다른건지... TODO : ConnectionConfigData와 합칠 수 있는지 검토할 것 --> 사용하지 않는것으로 보임
type CloudConnectionInfo struct {
	ID             string `json:"id"`
	ConfigName     string `json:"ConfigName"`
	ProviderName   string `json:"ProviderName"`
	DriverName     string `json:"DriverName"`
	CredentialName string `json:"CredentialName"`
	RegionName     string `json:"RegionName"`
	Description    string `json:"description"`
}
type CloudConnectionInfos []CloudConnectionInfo
