package model

type ClusterReq struct {
	Name         string       `json:"name"`
	ControlPlane []NodeConfig `json:"controlPlane"`
	Worker       []NodeConfig `json:"worker"`
	Config       Config       `json:"config"`
}

type NodeReq struct {
	ControlPlane []NodeConfig `json:"controlPlane"`
	Worker       []NodeConfig `json:"worker"`
}

type NodeConfig struct {
	Connection string `json:"connection"`
	Count      int    `json:"count"`
	Spec       string `json:"spec"`
}

type Config struct {
	Kubernetes Kubernetes `json:"kubernetes"`
}

type Kubernetes struct {
	NetworkCni       string `json:"networkCni"`
	PodCidr          string `json:"podCidr"`
	ServiceCidr      string `json:"serviceCidr"`
	ServiceDnsDomain string `json:"serviceDnsDomain"`
}
