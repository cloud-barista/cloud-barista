package mcis

type AutoCondition struct {
	EvaluationPeriod string   `json:"evaluationPeriod"`
	EvaluationValue  []string `json:"evaluationValue"`

	Metric   string `json:"metric"`
	Operand  string `json:"operand"`
	Operator string `json:"operator"`
}
