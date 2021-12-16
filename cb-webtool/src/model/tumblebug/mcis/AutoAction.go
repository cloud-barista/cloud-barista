package mcis

type AutoAction struct {
	ActionType    string     `json:"actionType"`
	PlacementAlgo string     `json:"placementAlgo"`
	PostCommand   McisCmdReq `json:"postCommand"`
	Vm            TbVmInfo   `json:"vm"`
}
