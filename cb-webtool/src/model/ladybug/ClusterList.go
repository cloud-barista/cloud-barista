package ladybug

// ladybug
type ClusterList struct {
	Kind     string        `json:"kind"`
	Clusters []ClusterInfo `json:"items"`
}
