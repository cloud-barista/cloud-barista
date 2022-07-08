package common

type RestRegisterCspNativeResourcesRequest struct {
	ConnectionName string `json:"connectionName"`
	McisName   string `json:"mcisName"`
	NsId   string `json:"nsId"`
}