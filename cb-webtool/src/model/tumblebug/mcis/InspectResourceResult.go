package mcis

type InspectResourceResult struct {
	ConnectionName          string                   `json:"connectionName"`
	CspOnlyOverview         InspectOverview          `json:"cspOnlyOverview"`
	ElapsedTime             int                      `json:"elapsedTime"`
	SystemMessage           string                   `json:"systemMessage"`
	TumblebugOverview       InspectOverview          `json:"tumblebugOverview"`
}