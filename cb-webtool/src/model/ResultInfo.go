package model

// 호출 후 결과값을 return하기 위한 객체
// Create, Delete 등의 결과로 message만 넘기는 경우가 있어 필요함.
type ResultInfo struct {
	Message string `json:"message"`
	Result string `json:"result"`
}
