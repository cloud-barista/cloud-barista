package webtool

// Life Cycle command 전송용 : VM과 Lifecycle 이 다를 수 있으므로 각각 사용
type McisLifeCycle struct {
	NameSpaceID string   `json:"nameSpaceID"`
	McisID      string   `json:"mcisID"`
	QueryParams []string `json:"queryParams"` // queryParams에 들어올 수 있는 값: action, force / action : create, suspend, resume, terminate, delete  : Const.MCIS_LIFECYCLE_xxx / force : false, true
}
