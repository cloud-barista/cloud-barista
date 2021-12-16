package dragonfly

// Monitoring 수신 Data  return용
type VmMonitoringInfoByMemory struct {
	Name       string                      `json:"name"`
	Tags       VmMonitoringTag             `json:"tags"`
	ValuesList []VmMonitoringValueByMemory `json:"values"`
}
type VmMonitoringValueByMemory struct {
	MemCached      float64 `json:"mem_cached"`
	MemFree        float64 `json:"mem_free"`
	MemShared      float64 `json:"mem_shared"`
	MemTotal       float64 `json:"mem_total"`
	MemUsed        float64 `json:"mem_used"`
	MemUtilization float64 `json:"mem_utilization"`
	Time           string  `json:"time"`
}
