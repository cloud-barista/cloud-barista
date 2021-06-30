package tumblebug

type SpiderSpecInfo struct {
	Mem    string `json:"mem"`
	Name   string `json:"name"`
	Region string `json:"region"`

	Vcpu          SpiderVCpuInfo  `json:"vcpc"`
	Gpus          []SpiderGpuInfo `json:"gpu"`
	KeyValueInfos []KeyValueInfo  `json:"keyValueList"`
}
type SpiderSpecInfos []SpiderSpecInfo
