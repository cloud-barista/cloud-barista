package mcis

type MonResultSimpleResponse struct {
	McisID         string          `json:"mcisId"`
	McisMonitoring MonResultSimple `json:"mcisMonitoring"` // yes, no
	NamespaceID    string          `json:"nsId"`
}
