package tumblebug

type SecurityGroupInfo struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	ConnectionName string `json:"connectionName"`
	Description    string `json:"description"`
	VNetID         string `json:"vNetID"`

	CspSecurityGroupID   string `json:"cspSecurityGroupId"`
	CspSecurityGroupName string `json:"cspSecurityGroupName"`

	FirewallRules []FirewallRule `json:"firewallRules"`

	KeyValueInfos []KeyValueInfo `json:"keyValueList"`
}

type SecurityGroupInfos []SecurityGroupInfo
