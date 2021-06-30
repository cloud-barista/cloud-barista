package tumblebug

// 등록용과 수신용의 포맷이 달라 Reg 용으로 사용
type VNetRegInfo struct {
	Name           string          `json:"name"`
	ConnectionName string          `json:"connectionName"`
	CidrBlock      string          `json:"cidrBlock"`
	Description    string          `json:"description"`
	SubnetRegInfo  []SubnetRegInfo `json:"subnetInfoList"`
}

type VNetRegInfos []VNetRegInfo
