package mcir

type TbFirewallRuleInfo struct {
	Cidr       string `json:"cidr"`
	Direction  string `json:"direction" validate:"required"`
	FromPort   string `json:"fromPort" validate:"required"`
	ToPort     string `json:"toPort" validate:"required"`
	IpProtocol string `json:"ipprotocol" validate:"required"`
}
