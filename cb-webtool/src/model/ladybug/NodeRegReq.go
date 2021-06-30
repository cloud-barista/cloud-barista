package ladybug

// ladybug
type NodeRegReq struct {
	// ControlPlane []NodeConfig `json:"controlPlane"`
	Worker []NodeConfig `json:"worker"`
}
