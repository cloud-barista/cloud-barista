package model

// 호출 후 상태값을 return하기 위한 객체
type WebStatus struct {
	StatusCode int    `json:"code"`
	Message    string `json:"message"`
}
