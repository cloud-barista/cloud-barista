package ladybug

// ClusterInfo의 간단버전.
type ClusterSimpleInfo struct {
	ClusterConfig  string           `json:"clusterConfig"`
	CpLeader       string           `json:"cpLeader"`
	Kind           int              `json:"kind"`
	Mcis           string           `json:"mcis"`
	Name           string           `json:"name"`
	NameSpace      string           `json:"namespace"`
	NetworkCni     string           `json:"networkCni"`
	Status         string           `json:"status"`
	McksStatus     string           `json:"mcksStatus"` // icon(RUNNING, STOPPED, TERMINATED)를 나타내기 위한 변수
	UID            string           `json:"uid"`
	Nodes          []NodeSimpleInfo `json:"nodes"`
	TotalNodeCount int              `json:"totalNodeCount"`
	NodeCountMap   map[string]int   `json:"nodeCountMap"`
}
