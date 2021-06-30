package ladybug

// ladybug
// Node의 array이나 parameter 추가된 것(kind)이 있음.   Nodes 와 NodeList는 다른것임.
type NodeList struct {
	Items []NodeInfo `json:"items"`
	Kind  string     `json:"kind"`
}
