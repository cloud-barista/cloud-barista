package spider

type RegionInfo struct {
	RegionName       string             `json:"RegionName"`
	ProviderName     string             `json:"ProviderName"`
	KeyValueInfoList []KeyValueInfoList `json:"KeyValueInfoList"`
}

type RegionInfos []RegionInfo
