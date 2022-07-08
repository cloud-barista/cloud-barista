package mcis

type ResourceOnTumblebug struct {
	Count  int `json:"count"`
	Info   []ResourceOnTumblebugInfo `json:"info"`
}
