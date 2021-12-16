package mcis

type DeploymentPlan struct {
	Filter   FilterInfo   `json:"filter"`
	Limit    string       `json:"limit"`
	Priority PriorityInfo `json:"priority"`
}
