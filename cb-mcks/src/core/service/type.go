package service

type SpecList struct {
	Kind    string    `json:"kind"`
	Config  string    `json:"connectionName"`
	Vmspecs []Vmspecs `json:"items"`
}

type Vmspecs struct {
	Name   string `json:"name"`   // output
	Memory string `json:"memory"` // output
	CPU    struct {
		Count string `json:"count"` // output
		Clock string `json:"clock"` // output - GHz
	} `json:"cpu"`
}
