package types

const (
	PushPolicy = "push"
	PullPolicy = "pull"
)

// CB-Store key
const (
	Agent             = "/monitoring/agents/"
	MonConfig         = "/monitoring/configs"
	EventLog          = "/monitoring/eventLogs"
	CollectorPolicy   = "/monitoring/collectorPolicy"
	Topic             = "/push/topic"
	CollectorTopicMap = "/push/collectorTopicMap"
)

const (
	NsId    = "nsId"
	McisId  = "mcisId"
	VmId    = "vmId"
	OsType  = "osType"
	CspType = "cspType"
)

const (
	AgentCntCollectorPolicy = "AGENTCOUNT"
	CSPCollectorPolicy      = "CSP"
)

const (
	Alibaba     = "ALIBABA"
	Aws         = "AWS"
	Azure       = "AZURE"
	Cloudit     = "CLOUDIT"
	Cloudtwin   = "CLOUDTWIN"
	Docker      = "DOCKER"
	Gcp         = "GCP"
	Openstack   = "OPENSTACK"
	TotalCspCnt = 8
)

const (
	KafkaDefaultPort     = 9092
	InfluxDefaultPort    = 8086
	KapacitorDefaultPort = 9092
)

const (
	Dev     = "dev"
	Helm    = "helm"
	Compose = "compose"
)

const (
	TopicAdd = "TopicAdd"
	TopicDel = "TopicDel"
)

const (
	ConfigMapName  = "cb-dragonfly-collector-configmap"
	DeploymentName = "cb-dragonfly-collector-"
)

const (
	LabelKey  = "name"
	Namespace = "dragonfly"
	//CollectorImage = "docker.io/hojun121/collector:latest"
	CollectorImage = "cloudbaristaorg/cb-dragonfly:0.5.0-collector"
)

const (
	TBRestAPIURL = "http://localhost:1323/tumblebug"
)
