package common

type RestInspectResourcesRequest struct {
	ConnectionName string `json:"connectionName"`
	ResourceType   string `json:"resourceType"`
}
