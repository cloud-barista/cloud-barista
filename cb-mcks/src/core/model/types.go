package model

const (
	KIND_STATUS       = "Status"
	KIND_CLUSTER      = "Cluster"
	KIND_CLUSTER_LIST = "ClusterList"
	KIND_NODE         = "Node"
	KIND_NODE_LIST    = "NodeList"
)

type Model struct {
	Name string `json:"name"`
	Kind string `json:"kind"`
}
type ListModel struct {
	Kind string `json:"kind"`
}
