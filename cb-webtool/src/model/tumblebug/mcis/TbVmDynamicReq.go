package mcis

type TbVmDynamicReq struct {
	CommonImage string `json:"commonImage"`
	CommonSpec  string `json:"commonSpec"`
	Description string `json:"description"`
	Label       string `json:"label"`
	Name        string `json:"name"`
	VmGroupSize string `json:"VmGroupSize"`
}
