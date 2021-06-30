package tumblebug

// cloud에 등록된 resource 조회 응답용
type InspectResourcesResponse struct {
	ResourcesOnCsp       ResourcesOnCsp         `json:"resourcesOnCsp"`
	ResourcesOnSpider    ResourcesOnSpider      `json:"resourcesOnSpider"`
	ResourcesOnTumblebug []ResourcesOnTumblebug `json:"resourcesOnTumblebug"`
}

type ResourcesOnCsp struct {
	cspNativeID string `json:"cspNativeId"`
	ID          string `json:"id"`
}
type ResourcesOnSpider struct {
	cspNativeID string `json:"cspNativeId"`
	ID          string `json:"id"`
}
type ResourcesOnTumblebug struct {
	CspNativeID string `json:"cspNativeId"`
	ID          string `json:"id"`
	McisID      string `json:"mcisId"`
	NsID        string `json:"nsId"`
	ObjectKey   string `json:"objectKey"`
	Type        string `json:"type"`
}
