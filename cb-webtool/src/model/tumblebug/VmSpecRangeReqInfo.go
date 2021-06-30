package tumblebug

// Spec filterling 할 때 range로 조회하기 위함( 요청용으로 int로 하면 default값이 0으로 설정되어 날아가므로 string으로 전송)
type VmSpecRangeReqInfo struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	OsType         string `json:"os_type"`
	GpuModel       string `json:"gpu_model"`
	GpuP2p         string `json:"gpu_p2p"`
	ConnectionName string `json:"connectionName"`
	CspSpecName    string `json:"cspSpecName"`
	Description    string `json:"description"`

	CostPerHour        RangeMinMaxReq `json:"cost_per_hour"`
	EbsBwMbps          RangeMinMaxReq `json:"ebs_bw_Mbps"`
	GpumemGiB          RangeMinMaxReq `json:"gpumem_GiB"`
	MaxNumStorage      RangeMinMaxReq `json:"max_num_storage"`
	MaxTotalStorageTiB RangeMinMaxReq `json:"max_total_storage_TiB"`
	MemGiB             RangeMinMaxReq `json:"mem_GiB"`
	NetBwGbps          RangeMinMaxReq `json:"net_bw_Gbps"`
	NumCore            RangeMinMaxReq `json:"num_core"`
	NumGpu             RangeMinMaxReq `json:"num_gpu"`
	NumStorage         RangeMinMaxReq `json:"num_storage"`
	NumVCPU            RangeMinMaxReq `json:"num_vCPU"`
	StorageGiB         RangeMinMaxReq `json:"storage_GiB"`
	EvaluationStatus   string         `json:"evaluationStatus"`
	EvaluationScore01  RangeMinMaxReq `json:"evaluationScore_01"`
	EvaluationScore02  RangeMinMaxReq `json:"evaluationScore_02"`
	EvaluationScore03  RangeMinMaxReq `json:"evaluationScore_03"`
	EvaluationScore04  RangeMinMaxReq `json:"evaluationScore_04"`
	EvaluationScore05  RangeMinMaxReq `json:"evaluationScore_05"`
	EvaluationScore06  RangeMinMaxReq `json:"evaluationScore_06"`
	EvaluationScore07  RangeMinMaxReq `json:"evaluationScore_07"`
	EvaluationScore08  RangeMinMaxReq `json:"evaluationScore_08"`
	EvaluationScore09  RangeMinMaxReq `json:"evaluationScore_09"`
	EvaluationScore10  RangeMinMaxReq `json:"evaluationScore_10"`
}

// Range 조회를 위해
type RangeMinMaxReq struct {
	MAX float32 `json:"max"`
	MIN float32 `json:"min"`
}
