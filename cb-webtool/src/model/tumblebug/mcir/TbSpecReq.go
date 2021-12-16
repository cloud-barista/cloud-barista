package mcir

type TbSpecReq struct {
	ConnectionName string `json:"connectionName"`
	CspSpecName    string `json:"cspSpecName"`
	Description    string `json:"description"`
	Name           string `json:"name"`
}
