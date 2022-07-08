package mcir

//
type FilterSpecsByRangeRequest struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	ProviderName   string `json:"providerName"`
	RegionName     string `json:"regionName"`
	Description    string `json:"description"`
	ConnectionName string `json:"connectionName"`
	CspSpecName    string `json:"cspSpecName"`
	OsType         string `json:"osType"`

	CostPerHour Range `json:"costPerHour"`
	EbsBwMbps   Range `json:"ebsBwMbps"`

	EvaluationScore01 Range  `json:"evaluationScore01"`
	EvaluationScore02 Range  `json:"evaluationScore02"`
	EvaluationScore03 Range  `json:"evaluationScore03"`
	EvaluationScore04 Range  `json:"evaluationScore04"`
	EvaluationScore05 Range  `json:"evaluationScore05"`
	EvaluationScore06 Range  `json:"evaluationScore06"`
	EvaluationScore07 Range  `json:"evaluationScore07"`
	EvaluationScore08 Range  `json:"evaluationScore08"`
	EvaluationScore09 Range  `json:"evaluationScore09"`
	EvaluationScore10 Range  `json:"evaluationScore10"`
	EvaluationStatus  string `json:"evaluationStatus"`

	GpuModel string `json:"gpuModel"`
	GpuP2p   string `json:"gpuP2p"`

	MaxNumStorage      Range `json:"maxNumStorage"`
	MaxTotalStorageTiB Range `json:"maxTotalStorageTiB"`
	MemGiB             Range `json:"memGiB"`

	NetBwGbps  Range `json:"netBwGbps"`
	NumCore    Range `json:"numCore"`
	NumGpu     Range `json:"numGpu"`
	NumStorage Range `json:"numStorage"`
	NumVCPU    Range `json:"numvCPU"`
	StorageGiB Range `json:"storageGiB"`
}
