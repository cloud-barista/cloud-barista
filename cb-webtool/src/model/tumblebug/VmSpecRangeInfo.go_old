package tumblebug

// Spec filterling 할 때 range로 조회하기 위함
type VmSpecRangeInfo struct {
	CostPerHour        RangeMinMax `json:"cost_per_hour"`
	EbsBwMbps          RangeMinMax `json:"ebs_bw_Mbps"`
	GpumemGiB          RangeMinMax `json:"gpumem_GiB"`
	MaxNumStorage      RangeMinMax `json:"max_num_storage"`
	MaxTotalStorageTiB RangeMinMax `json:"max_total_storage_TiB"`
	MemGiB             RangeMinMax `json:"mem_GiB"`
	NetBwGbps          RangeMinMax `json:"net_bw_Gbps"`
	NumCore            RangeMinMax `json:"num_core"`
	NumGpu             RangeMinMax `json:"num_gpu"`
	NumStorage         RangeMinMax `json:"num_storage"`
	NumVCPU            RangeMinMax `json:"num_vCPU"`
	StorageGiB         RangeMinMax `json:"storage_GiB"`
	EvaluationStatus   string      `json:"evaluationStatus"`
	EvaluationScore01  RangeMinMax `json:"evaluationScore_01"`
	EvaluationScore02  RangeMinMax `json:"evaluationScore_02"`
	EvaluationScore03  RangeMinMax `json:"evaluationScore_03"`
	EvaluationScore04  RangeMinMax `json:"evaluationScore_04"`
	EvaluationScore05  RangeMinMax `json:"evaluationScore_05"`
	EvaluationScore06  RangeMinMax `json:"evaluationScore_06"`
	EvaluationScore07  RangeMinMax `json:"evaluationScore_07"`
	EvaluationScore08  RangeMinMax `json:"evaluationScore_08"`
	EvaluationScore09  RangeMinMax `json:"evaluationScore_09"`
	EvaluationScore10  RangeMinMax `json:"evaluationScore_10"`
}

// Range 조회를 위해
type RangeMinMax struct {
	MAX int `json:"max"`
	MIN int `json:"min"`
}
