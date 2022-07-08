package mcis

type InspectResourceAllResult struct {
	AvailableConnection     int                      `json:"availableConnection"`
	CspOnlyOverview         InspectOverview          `json:"cspOnlyOverview"`
	ElapsedTime             int                      `json:"elapsedTime"`
	InspectResult           []InspectResourceResult  `json:"inspectResult"`
	RegisteredConnection    int                      `json:"registeredConnection"`
	TumblebugOverview       InspectOverview          `json:"tumblebugOverview"`
	
}