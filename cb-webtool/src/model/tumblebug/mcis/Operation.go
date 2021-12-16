package mcis

type Operation struct {
	Operand  string `json:"operand"`
	Operator string `json:"operator"` //  >=, =, <=, ==
}
