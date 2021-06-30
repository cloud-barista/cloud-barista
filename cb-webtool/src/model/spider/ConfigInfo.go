package spider

type ConfigInfo struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Value string `json:"value"`
}
type ConfigInfos []ConfigInfo
