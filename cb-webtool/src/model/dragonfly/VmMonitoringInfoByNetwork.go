package dragonfly

// Monitoring 수신 Data  return용
type VmMonitoringInfoByNetwork struct {
	Name       string                       `json:"name"`
	Tags       VmMonitoringTag              `json:"tags"`
	ValuesList []VmMonitoringValueByNetwork `json:"values"`
}
type VmMonitoringValueByNetwork struct {
	BytesIn  float64 `json:"bytes_in"`
	BytesOut float64 `json:"bytes_out"`
	PktsIn   float64 `json:"pkts_in"`
	PktsOut  float64 `json:"pkts_out"`
	Time     string  `json:"time"`
}
