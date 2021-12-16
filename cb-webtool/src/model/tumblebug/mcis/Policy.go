package mcis

type Policy struct {
	AutoAction    AutoAction    `json:"autoAction"`
	AutoCondition AutoCondition `json:"autoCondition"`
	Status        string        `json:"status"`
}

type Policies []Policy
