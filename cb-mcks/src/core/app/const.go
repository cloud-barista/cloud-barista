package app

type CSP string
type ROLE string
type Kind string
type NetworkCni string
type StatusCode int

const (
	CSP_AWS       CSP = "aws"
	CSP_GCP       CSP = "gcp"
	CSP_AZURE     CSP = "azure"
	CSP_ALIBABA   CSP = "alibaba"
	CSP_TENCENT   CSP = "tencent"
	CSP_OPENSTACK CSP = "openstack"
	CSP_IBM       CSP = "ibm"
	CSP_CLOUDIT   CSP = "cloudit"

	CONTROL_PLANE ROLE = "control-plane"
	WORKER        ROLE = "worker"

	KIND_STATUS       Kind = "Status"
	KIND_CLUSTER      Kind = "Cluster"
	KIND_CLUSTER_LIST Kind = "ClusterList"
	KIND_NODE         Kind = "Node"
	KIND_NODE_LIST    Kind = "NodeList"

	STATUS_UNKNOWN  = 0
	STATUS_SUCCESS  = 200
	STATUS_NOTFOUND = 404

	NETWORKCNI_KILO  NetworkCni = "kilo"
	NETWORKCNI_CANAL NetworkCni = "canal"

	POD_CIDR       = "10.244.0.0/16"
	SERVICE_CIDR   = "10.96.0.0/12"
	SERVICE_DOMAIN = "cluster.local"

	LABEL_KEY_CSP    = "topology.cloud-barista.github.io/csp"
	LABEL_KEY_REGION = "topology.kubernetes.io/region"
	LABEL_KEY_ZONE   = "topology.kubernetes.io/zone"

	MCIS_LABEL       = "mcks"
	MCIS_SYSTEMLABEL = "Managed by MCKS"
)

type Status struct {
	Kind    Kind   `json:"kind"`
	Code    int    `json:"code"`
	Message string `json:"message" example:"Any message"`
}

type ClusterReq struct {
	Name            string           `json:"name" example:"cluster-01"`
	ControlPlane    []NodeSetReq     `json:"controlPlane"`
	Worker          []NodeSetReq     `json:"worker"`
	Config          ClusterConfigReq `json:"config"`
	Label           string           `json:"label"`
	InstallMonAgent string           `json:"installMonAgent" example:"no" default:"yes"`
	Description     string           `json:"description"`
}

type NodeReq struct {
	ControlPlane []NodeSetReq `json:"controlPlane"`
	Worker       []NodeSetReq `json:"worker"`
}

type NodeSetReq struct {
	Connection string `json:"connection" example:"config-aws-ap-northeast-2"`
	Count      int    `json:"count" example:"3"`
	Spec       string `json:"spec" example:"t2.medium"`
}

type ClusterConfigReq struct {
	Kubernetes ClusterConfigKubernetesReq `json:"kubernetes"`
}
type ClusterConfigKubernetesReq struct {
	NetworkCni       NetworkCni `json:"networkCni" example:"kilo" enums:"canal,kilo" default1:"kilo"`
	PodCidr          string     `json:"podCidr" example:"10.244.0.0/16"`
	ServiceCidr      string     `json:"serviceCidr" example:"10.96.0.0/12"`
	ServiceDnsDomain string     `json:"serviceDnsDomain" example:"cluster.local"`
}
