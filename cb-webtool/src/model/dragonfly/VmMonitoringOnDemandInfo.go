package dragonfly

// 멀티 클라우드 인프라 VM 온디맨드 모니터링 정보 결과용
type VmMonitoringOnDemandInfo struct {
	NameSpaceID        string                    `json:"name"`
	VmMonitoringTags   VmMonitoringTagOnDemand   `json:"tags"`
	Time               string                    `json:"time"`
	VmMonitoringValues VmMonitoringValueOnDemand `json:"values"`
}

type VmMonitoringTagOnDemand struct {
	McisId string `json:"mcisId"`
	NsId   string `json:"nsId"`
	VmId   string `json:"vmId"`
}

type VmMonitoringValueOnDemand struct {
	BytesIn  int `json:"bytes_in"`
	BytesOut int `json:"bytes_out"`
	PktsIn   int `json:"pkts_in"`
	PktsOut  int `json:"pkts_out"`
}
