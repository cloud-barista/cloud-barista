package spider

type DriverInfo struct {
	//user(username, password, email)
	DriverName        string `json:"DriverName"`
	ProviderName      string `json:"ProviderName"`
	DriverLibFileName string `json:"DriverLibFileName"`
}
type DriverInfos []DriverInfo
