package tumblebug

type VNetInfo struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	ConnectionName string `json:"connectionName"`
	CidrBlock      string `json:"cidrBlock"`
	CspVNetID      string `json:"cspVNetID"`
	CspVNetName    string `json:"cspVNetName"`
	Description    string `json:"description"`

	Status        string         `json:"status"`
	KeyValueInfos []KeyValueInfo `json:"keyValueList"`
	SubnetInfos   []SubnetInfo   `json:"subnetInfoList"`
}

type VNetInfos []VNetInfo
