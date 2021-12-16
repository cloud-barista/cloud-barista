package dragonfly

// Monitoring 수신 Data  return용. parameter에 따라 return받는 객체가 다름..
//type VmMonitoringInfo struct {
//	Name       string              `json:"name"`
//	Tags       VmMonitoringTag     `json:"tags"`
//	ValuesList []interface{} `json:"values"`
//	//ValuesList []VmMonitoringValue `json:"values"`
//}

type VmMonitoringInfo struct {
	ValuesByCpu     VmMonitoringInfoByCpu
	ValuesByDisk    VmMonitoringInfoByDisk
	ValuesByMemory  VmMonitoringInfoByMemory
	ValuesByNetwork VmMonitoringInfoByNetwork
}
