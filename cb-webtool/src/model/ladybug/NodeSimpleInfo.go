package ladybug

// NodeInfo 의 간단버전  credential 등 제거
type NodeSimpleInfo struct {
	NodeIndex	   int `json:"nodeIndex"`
	TotalNodeCount int `json:"totalNodeCount"`
	NodeCsp        string `json:"csp"`
	NodeKind       string `json:"kind"`
	NodeName       string `json:"name"`
	NodePublicIp   string `json:"publicIp"`
	NodeRole       string `json:"role"`
	NodeSpec       string `json:"spec"`
	NodeUID        string `json:"uid"`
}
type NodeSimpleInfos []NodeSimpleInfo
