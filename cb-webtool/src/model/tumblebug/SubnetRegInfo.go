package tumblebug

type SubnetRegInfo struct {
	Ipv4_CIDR     string         `json:"ipv4_CIDR"`
	Name          string         `json:"name"`
	KeyValueInfos []KeyValueInfo `json:"keyValueList"`
}

type SubnetRegInfos []SubnetRegInfo
