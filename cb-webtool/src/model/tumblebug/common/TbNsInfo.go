package common

type TbNsInfo struct {
	ID          string `json:"id"`
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
}

type TbNsInfos []TbNsInfo
