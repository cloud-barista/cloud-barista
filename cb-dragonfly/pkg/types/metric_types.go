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
	MCK8S_NODE   Metric = "node"
	MCK8S_POD    Metric = "pod"
)

const (
	MCIS       string = "mcis"
	VM         string = "vm"
	MCK8S      string = "mck8s"
	KUBERNETES string = "kubernetes"
	K8S        string = "k8s"
	Cluster    string = "cluster"
	Node       string = "node"
	Namespace  string = "namespace"
	ALL        string = "all"
)

type MCK8SReqInfo struct {
	GroupBy   string // Node: Cluster or Node  // Pod: Node or Namespace or Pod
	Node      string
	Namespace string
	Pod       string
}

type DBMetricRequestInfo struct {
	NsID                string
	ServiceID           string
	ServiceType         string
	MonitoringMechanism bool // push: true
	VMID                string
	MetricName          string
	MCK8SReqInfo
	Period       string
	AggegateType string
	Duration     string
}

func (m Metric) ToString() string {
	if m == "" {
		return "none"
	}
	return string(m)
}

func (m Metric) ToAgentMetricKey() string {
	//var metricKey string
	if m == Cpu || m == CpuFrequency || m == Disk || m == DiskIO || m == MCK8S_NODE || m == MCK8S_POD {
		return m.ToString()
	} else if m == Memory {
		return "mem"
	} else if m == Network {
		return "net"
	} else {
		return "none"
	}
}

func GetMetricType(input string) Metric {
	switch input {
	case Cpu.ToString():
		return Cpu
	case CpuFrequency.ToString():
		return CpuFrequency
	case Memory.ToString(), Memory.ToAgentMetricKey():
		return Memory
	case Network.ToString(), Network.ToAgentMetricKey():
		return Network
	case Disk.ToString():
		return Disk
	case DiskIO.ToString():
		return DiskIO
	default:
		return None
	}
}

// CBMCISMetric 단일 MCIS Milkyway 메트릭
type CBMCISMetric struct {
	Result  string `json:"result"`
	Unit    string `json:"unit"`
	Desc    string `json:"desc"`
	Elapsed string `json:"elapsed"`
	SpecId  string `json:"specid"`
}

// MCBMCISMetric 멀티 MCIS Milkyway 메트릭
type MCBMCISMetric struct {
	ResultArray []CBMCISMetric `json:"resultarray"`
}

// Request GET Request 단일 Body 정보
type Request struct {
	Host string `json:"host"`
	Spec string `json:"spec"`
}

// Mrequest GET Request 멀티 Body 정보
type Mrequest struct {
	MultiHost []Request `json:"multihost"`
}

type Parameter struct {
	agent_ip    string
	mcis_metric string
}

type DBData struct {
	Name    string            `json:"name"`
	Tags    map[string]string `json:"tags"`
	Columns []string          `json:"columns"`
	Values  [][]interface{}   `json:"values"`
}
