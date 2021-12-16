package common

type RestInspectResourcesRequest struct {
	ConnectionName string `json:"connectionName"`
	Type           string `json:"type"`
}
