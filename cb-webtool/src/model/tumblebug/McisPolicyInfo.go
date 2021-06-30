package tumblebug

// Life Cycle command 전송용 : VM과 Lifecycle 이 다를 수 있으므로 각각 사용
type McisPolicyInfo struct {
	ID          string     `json:"Id"`
	Name        string     `json:"Name"`
	ActionLog   string     `json:"actionLog"`
	Description string     `json:"description"`
	McisPolicy  McisPolicy `json:"policy"`
}

type McisPolicy struct {
	AutoAction    AutoAction    `json:"autoAction"`
	AutoCondition AutoCondition `json:"autoCondition"`
	Status        string        `json:"status"`
}

type AutoAction struct {
	ActionType    string          `json:"actionType"`
	Placementalgo string          `json:"placement_algo"`
	PostCommand   McisCommandInfo `json:"postCommand"`
	Vm            VmInfo          `json:"vm"`
}

type AutoCondition struct {
	EvaluationPeriod string   `json:"evaluationPeriod"`
	EvaluationValues []string `json:"evaluationValue"`
	Metric           string   `json:"metric"`
	Operand          string   `json:"operand"`
	Operator         string   `json:"operator"`
}
