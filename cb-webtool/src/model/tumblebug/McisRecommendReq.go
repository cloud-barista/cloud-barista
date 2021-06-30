package tumblebug

type McisRecommendReq struct {
	MaxResultNum   string       `json:"max_result_num"`
	PlacementAlgo  string       `json:"placement_algo"`
	PlacementParam KeyValueInfo `json:"placement_param"`
	VmReq          VmReq        `json:"vm_req"`
}
