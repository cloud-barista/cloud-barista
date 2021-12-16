package dragonfly

// Monitoring 수신 Data  return용
type VmMonitoringInfoByDisk struct {
	Name       string                    `json:"name"`
	Tags       VmMonitoringTag           `json:"tags"`
	ValuesList []VmMonitoringValueByDisk `json:"values"`
}
type VmMonitoringValueByDisk struct {
	Free        float64 `json:"free"`
	ReadBytes   float64 `json:"read_bytes"`
	ReadTime    float64 `json:"read_time"`
	Reads       float64 `json:"reads"`
	Total       float64 `json:"total"`
	Used        float64 `json:"used"`
	UsedPercent float64 `json:"used_percent"`
	WriteBytes  float64 `json:"write_bytes"`
	WriteTime   float64 `json:"write_time"`
	Writes      float64 `json:"writes"`
	Time        string  `json:"time"`
}
