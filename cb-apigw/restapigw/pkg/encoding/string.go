package encoding

import (
	"io"
	"io/ioutil"
)

// ===== [ Constants and Variables ] =====

const (
	// STRING - 단순 문자열 인코딩 식별자
	STRING = "string"
)

// ===== [ Types ] =====

// ===== [ Implementations ] =====

// ===== [ Private Functions ] =====

// ===== [ Public Functions ] =====

// StringDecoder - 지정한 Reader의 문자열 데이터를 기준으로 Decode 처리
func StringDecoder(r io.Reader, v *map[string]interface{}) error {
	data, err := ioutil.ReadAll(r)
	if nil != err {
		return err
	}
	*(v) = map[string]interface{}{"content": string(data)}
	return nil
}

// NewStringDecoder - 단순 문자열을 위한 Decoder 생성
func NewStringDecoder(_ bool, _ bool) func(io.Reader, *map[string]interface{}) error {
	return StringDecoder
}
