package model

// 사용하지 않는것으로 보임
type McisRequest struct {
	VmSpec           []string `form:"vmspec"`
	NameSpace        string   `form:"namespace"`
	McisName         string   `form:"mcis_name"`
	VmName           []string `form:"vmName"`
	Provider         []string `form:"provider"`
	SecurityGroupIds []string `form:"sg"`
}
