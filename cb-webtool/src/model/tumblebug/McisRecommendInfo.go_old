package tumblebug

type McisRecommendInfo struct {
	PlacementAlgo  string       `json:"placement_algo"`
	PlacementParam KeyValueInfo `json:"placement_param"`
	VmRecommend    VmRecommend  `json:"vm_recommend"`
}

type VmRecommend struct {
	PlacementAlgo  string       `json:"placement_algo"`
	PlacementParam KeyValueInfo `json:"placement_param"`

	VmPriority VmPriority `json:"vm_priority"`
}

type VmPriority struct {
	Priority string     `json:"priority"`
	VmSpec   VmSpecInfo `json:"vm_spec"`
}
