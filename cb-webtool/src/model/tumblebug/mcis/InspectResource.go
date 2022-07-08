package mcis

type InspectResource struct {
	ConnectionName     string                    `json:"connectionName"`
	ResourceOverview   ResourceCountOverview     `json:"resourceOverview"`
	ResourceType       string                    `json:"resourceType"`
	Resources          ResourcesByManageType     `json:"resources"`
	SystemMessage      string                    `json:"systemMessage"`
}