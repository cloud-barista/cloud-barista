package mcir

type SpiderSecurityRuleInfo struct {
	Cidr       string `json:"cidr"`
	Direction  string `json:"direction"`
	FromPort   string `json:"fromPort"`
	ToPort     string `json:"toPort"`
	IpProtocol string `json:"ipProtocol"`
}
