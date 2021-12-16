package mcir

type TbSshKeyReq struct {
	ConnectionName string `json:"connectionName"`
	Description    string `json:"description"`
	Name           string `json:"name"`
}
