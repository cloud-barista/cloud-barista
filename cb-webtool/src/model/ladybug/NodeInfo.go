package ladybug

// ladybug
type NodeInfo struct {
	Credential string `json:"credential"`
	Csp        string `json:"csp"`
	Kind       string `json:"kind"` // Node 냐 cluster냐
	Name       string `json:"name"`
	PublicIp   string `json:"publicIp"`
	Role       string `json:"role"` // Control-plane냐, Worker냐
	Spec       string `json:"spec"`
	UID        string `json:"uid"`
}
type Nodes []NodeInfo
