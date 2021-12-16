package mcis

type FilterCondition struct {
	Condition []Operation `json:"condition"`
	Metric    string      `json:"metric"`
}
