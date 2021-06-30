package types

type Metric string

const (
	Cpu          Metric = "cpu"
	CpuFrequency Metric = "cpufreq"
	Memory       Metric = "memory"
	Disk         Metric = "disk"
	DiskIO       Metric = "diskio"
	Network      Metric = "network"
	None         Metric = "none"
)

func (m Metric) ToString() string {
	if m == "" {
		return "none"
	}
	return string(m)
}

func (m Metric) ToAgentMetricKey() string {
	//var metricKey string
	if m == Cpu || m == CpuFrequency || m == Disk || m == DiskIO {
		return m.ToString()
	} else if m == Memory {
		return "mem"
	} else if m == Network {
		return "net"
	} else {
		return "none"
	}
}
