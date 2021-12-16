package mcir

type TbVNetReq struct {
	CidrBlock      string        `json:"cidrBlock"`
	ConnectionName string        `json:"connectionName"`
	Description    string        `json:"description"`
	Name           string        `json:"name"`
	SubnetInfoList []TbSubnetReq `json:"subnetInfoList"`
}
