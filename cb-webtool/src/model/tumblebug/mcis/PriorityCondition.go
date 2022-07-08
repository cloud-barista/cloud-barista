package mcis

type PriorityCondition struct {
	Metric    string          `json:"metric"`
	Parameter []ParameterKeyVal `json:"parameter"`
	Weight    string         `json:"weight"`
}
