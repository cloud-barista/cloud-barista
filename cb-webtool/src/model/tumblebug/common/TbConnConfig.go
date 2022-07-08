package common

// TB의 common.ConfigInfo -> TBConfigInfo로   ConfigInfo가 spider에도 있으므로  common package는 모두 TB를 붙이도록 한다
type TbConConfig struct {
	ConfigName     string       `json:"configName"`
	ProviderName   string       `json:"providerName"`
	CredentialName string       `json:"credentialName"`
	DriverName     string       `json:"driverName"`
	RegionName     string       `json:"regionName"`
	Location       TbGeoLocation  `json:"location"`
}

type TbConConfigs []TbConConfig
