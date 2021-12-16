package common

type TbRegion struct {
	ProviderName     string       `json:"providerName"`
	RegionName       string       `regionName:"regionName"`
	KeyValueInfoList []TbKeyValue `json:"keyValueInfoList"`
}

type TbRegions []TbRegion
