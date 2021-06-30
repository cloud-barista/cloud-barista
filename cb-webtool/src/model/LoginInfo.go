package model

// client와 서버가 통신하면서 사용자정보로 이용할 객체
// 최초로그인 시 Username set.
// page redirect 될 때 default namespace 설정
// - namespace가 없으면 새로 생성
// - namespace가 1개면 default로 set
// - namespace가 2개 이상이면 화면에서 default set 하도록
type LoginInfo struct {
	UserID               string `json:"UserID"`
	Username             string `json:"Username"`
	DefaultNameSpaceID   string `json:"DefaultNameSpaceID"`
	DefaultNameSpaceName string `json:"DefaultNameSpaceName"`
	AccessToken          string `json:"AccessToken"`
}
