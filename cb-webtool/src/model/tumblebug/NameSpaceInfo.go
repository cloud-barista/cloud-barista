package tumblebug

type NameSpaceInfo struct {
	ID          string `json:"id"`
	Name        string `json:"name" validate:"required`
	Description string `json:"description"`
}

type NameSpaceInfos []NameSpaceInfo
