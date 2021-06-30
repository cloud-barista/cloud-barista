package ladybug

type NodeConfig struct {
	Connection string `json:"connection"`
	Count      int    `json:"count"`
	Spec       string `json:"spec"`
}
