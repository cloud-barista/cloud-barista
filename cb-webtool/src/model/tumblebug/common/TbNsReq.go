package common

type TbNsReq struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
}
