package mcis

type InspectOverview struct {
	SecurityGroup   int  `json:"securityGroup"`
	SshKey          int  `json:"sshKey"`
	VNet            int  `json:"vNet"`
	Vm              int  `json:"vm"`
}