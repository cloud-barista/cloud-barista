package mcis

type BenchmarkInfo struct {
	Desc        string             `json:"desc"`
	Elapsed     string             `json:"elapsed"`
	Result      string             `json:"result"`
	ResultArray StructElementCycle `json:"struct-element-cycle"` //?
	SpecID      string             `json:"specid"`
	Unit        string             `json:"unit"`
}

type BenchmarkInfos []BenchmarkInfo

type StructElementCycle struct {
	Cycle string `json:"cycle"`
}
