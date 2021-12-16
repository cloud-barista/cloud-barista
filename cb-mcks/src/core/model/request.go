package model

type ClusterReq struct {
	Name            string       `json:"name" example:"cluster-01"`
	ControlPlane    []NodeConfig `json:"controlPlane"`
	Worker          []NodeConfig `json:"worker"`
	Config          Config       `json:"config"`
	Label           string       `json:"label"`
	InstallMonAgent string       `json:"installMonAgent" example:"no" default:"yes"`
	Description     string       `json:"description"`
}

type NodeReq struct {
	ControlPlane []NodeConfig `json:"controlPlane"`
	Worker       []NodeConfig `json:"worker"`
}

type NodeConfig struct {
	Connection string `json:"connection" example:"config-aws-ap-northeast-2"`
	Count      int    `json:"count" example:"3"`
	Spec       string `json:"spec" example:"t2.medium"`
}

type Config struct {
	Kubernetes Kubernetes `json:"kubernetes"`
}

type Kubernetes struct {
	NetworkCni       string `json:"networkCni" example:"canal" enums:"canal,kilo"`
	PodCidr          string `json:"podCidr" example:"10.244.0.0/16"`
	ServiceCidr      string `json:"serviceCidr" example:"10.96.0.0/12"`
	ServiceDnsDomain string `json:"serviceDnsDomain" example:"cluster.local"`
}
