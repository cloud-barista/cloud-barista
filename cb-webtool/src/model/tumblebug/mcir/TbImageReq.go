package mcir

type TbImageReq struct {
	ConnectionName string `json:"connectionName"`
	CspImageId     string `json:"cspImageId"`
	Name           string `json:"name"`
}
