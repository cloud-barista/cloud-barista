package mcis

type McisPolicyInfo struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	ActionLog   string `json:"actionLog"`
	Description string `json:"description"`
	policy      Policy `json:"policy"`
}
