package ladybug

// Node만 추가시에는 ControlPlane 없이 Worker만 사용한다.
type NodeOnlyRegReq struct {
	// ControlPlane []NodeConfig `json:"controlPlane"`
	//Worker []NodeReq `json:"worker"`
	Worker []NodeReq `json:"worker"`
}
