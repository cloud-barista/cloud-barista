package mcis

type TbMcisDynamicReq struct {
	Description     string           `json:"description"`
	InstallMonAgent string           `json:"installMonAgent"`
	Label           string           `json:"label"`
	Name            string           `json:"name"`
	SystemLabel     string           `json:"systemLabel"`
	Vm              []TbVmDynamicReq `json:"vm"`
}
