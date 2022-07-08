package mcis

// ResourcesOnCsp 와 ResourcesOnSpider 가 동일하여  ResourecesOnCspOrSpider 로 사용하는 것 같음.
type ResourceOnCsp struct {
	Count  int `json:"count"`
	Info   []ResourceOnCspInfo `json:"info"`
}
