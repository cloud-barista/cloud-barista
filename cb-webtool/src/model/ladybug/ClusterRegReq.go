package ladybug

// ladybug
type ClusterRegReq struct {
	// ControlPlaneNodeCount int    `json:"controlPlaneNodeCount"`
	// ControlPlaneNodeSpec  string `json:"controlPlaneNodeSpec"`
	// Name                  string `json:"name"`
	// WorkerNodeCount       int    `json:"workerNodeCount"`
	// WorkerNodeSpec        string `json:"workerNodeSpec"`

	Name         string        `json:"name"`
	Config       ClusterConfig `json:"config"`
	ControlPlane []NodeConfig  `json:"controlPlane"`
	Worker       []NodeConfig  `json:"worker"`
}
