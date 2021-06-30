package dragonfly

// 멀티 클라우드 인프라 MCIS 온디맨드 모니터링 정보 결과용
type McisMonitoringOnDemandInfo struct {
	Result  string `json:"result"`
	Unit    string `json:"unit"`
	Desc    string `json:"desc"`
	Elapsed string `json:"elapsed"`
	Specid  string `json:"specid"`
}
