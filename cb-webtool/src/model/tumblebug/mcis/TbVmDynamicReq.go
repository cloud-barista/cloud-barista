package mcis

type TbVmDynamicReq struct {
	CommonImage    string `json:"commonImage"`
	CommonSpec     string `json:"commonSpec"`
	ConnectionName string `json:"connectionName"`
	Description    string `json:"description"`
	Label          string `json:"label"`
	Name           string `json:"name"`
	RootDiskSize   string `json:"rootDiskSize"`
	RootDiskType   string `json:"rootDiskType"`
	VmGroupSize    string `json:"VmGroupSize"`
}
