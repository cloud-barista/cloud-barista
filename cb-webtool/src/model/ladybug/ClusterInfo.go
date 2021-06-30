package ladybug

// ladybug
type ClusterInfo struct {
	ClusterConfig string     `json:"clusterConfig"`
	CpLeader      string     `json:"cpLeader"`
	Kind          int        `json:"kind"`
	Mcis          string     `json:"mcis"`
	Name          string     `json:"name"`
	NameSpace     string     `json:"namespace"`
	NetworkCni    string     `json:"networkCni"`
	Status        string     `json:"status"`
	UID           string     `json:"uid"`
	Nodes         []NodeInfo `json:"nodes"`
}
