package ladybug

// ladybug
type NodeReq struct {
	Config          string `json:"config"`
	WorkerNodeCount int    `json:"workerNodeCount"`
	WorkerNodeSpec  string `json:"workerNodeSpec"`
}
