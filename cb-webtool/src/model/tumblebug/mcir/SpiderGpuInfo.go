package mcir

type SpiderGpuInfo struct {
	Count string `json:"count"`
	Mem   string `json:"mem"`
	Mfr   string `json:"mfr"`
	Model string `json:"model"`
}
type SpiderGpuInfos []SpiderGpuInfo
