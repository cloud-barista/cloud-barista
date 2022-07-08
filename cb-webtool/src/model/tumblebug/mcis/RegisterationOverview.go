package mcis

type RegisterationOverview struct {
	Failed          int  `json:"failed"`
	SecurityGroup   int  `json:"securityGroup"`
	SshKey          int  `json:"sshKey"`
	VNet            int  `json:"vNet"`
	Vm              int  `json:"vm"`
}