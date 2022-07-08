package mcis

type TbMcisInfo struct {
	ConfigureCloudAdaptiveNetwork string `json:"configureCloudAdaptiveNetwork"`
	ID              string          `json:"id"`
	Description     string          `json:"description"`
	InstallMonAgent string          `json:"installMonAgent"`
	Label           string          `json:"label"`
	SystemLabel     string          `json:"systemLabel"`
	Name            string          `json:"name"`
	PlacementAlgo   string          `json:"placementAlgo"`
	Status          string          `json:"status"`
	StatusCount     StatusCountInfo `json:"statusCount"`
	TargetAction    string          `json:"targetAction"`
	TargetStatus    string          `json:"targetStatus"`

	Vm []TbVmInfo `json:"vm"`
}
