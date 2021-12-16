package dragonfly

// Monitoring 수신 Data  return용

type VmMonitoringInfoByCpu struct {
	Name       string                   `json:"name"`
	Tags       VmMonitoringTag          `json:"tags"`
	ValuesList []VmMonitoringValueByCpu `json:"values"`
}

type VmMonitoringValueByCpu struct {
	CpuGuest       float64 `json:"cpu_guest"`
	CpuGuestNice   float64 `json:"cpu_guest_nice"`
	CpuHintr       float64 `json:"cpu_hintr"`
	CpuIdle        float64 `json:"cpu_idle"`
	CpuIowait      float64 `json:"cpu_iowait"`
	CpuNice        float64 `json:"cpu_nice"`
	CpuSintr       float64 `json:"cpu_sintr"`
	CpuSteal       float64 `json:"cpu_steal"`
	CpuSystem      float64 `json:"cpu_system"`
	CpuUser        float64 `json:"cpu_user"`
	CpuUtilization float64 `json:"cpu_utilization"`
	Time           string  `json:"time"`
}
