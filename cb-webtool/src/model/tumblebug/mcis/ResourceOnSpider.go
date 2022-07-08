package mcis

// ResourcesOnCsp 와 ResourcesOnSpider 가 동일하여  ResourecesOnCspOrSpider 로 사용하는 것 같음.
type ResourceOnSpider struct {
	Count  int `json:"count"`
	Info   []ResourceOnSpiderInfo `json:"info"`
}
