package mcis

type ParameterKeyVal struct {
	Key string   `json:"key"`
	Val []string `json:"val"` //  >=, =, <=, ==
}
