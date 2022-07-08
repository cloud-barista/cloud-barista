package mcis

type NetworkReq struct {
	EtcdEndpoints     []string   `json:"etcdEndpoints"`
	ServiceEndpoint   string     `json:"serviceEndpoint"`
}