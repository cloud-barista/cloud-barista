// Package render - 응답 결과를 출력하는 기능을 제공하는 패키지
package render

import (
	"bytes"
	"net/http"

	jsoniter "github.com/json-iterator/go"
)

// ===== [ Constants and Variables ] =====
// ===== [ Types ] =====

type (
	// M - 단순 Map 처리 인터페이스
	M map[string]interface{}
)

// ===== [ Implementations ] =====
// ===== [ Private Functions ] =====
// ===== [ Public Functions ] =====

// JSON - Content-Type 을 `application/json`으로 설정하고 HTML escape 처리를 제공
func JSON(rw http.ResponseWriter, code int, v interface{}) {
	buf := &bytes.Buffer{}
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(true)
	if err := enc.Encode(v); nil != err {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(code)
	rw.Write(buf.Bytes())
}
