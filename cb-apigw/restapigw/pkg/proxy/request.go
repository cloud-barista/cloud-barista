package proxy

import (
	"bytes"
	"io"
	"net/url"
)

// ===== [ Constants and Variables ] =====

// ===== [ Types ] =====

// Request - Proxy 구간에서 사용할 Request 구조
type Request struct {
	IsBypass bool
	Method   string
	URL      *url.URL
	Query    url.Values
	Path     string
	Body     io.ReadCloser
	Params   map[string]string
	Headers  map[string][]string
}

// ===== [ Implementations ] =====

// GeneratePath - Params의 정보를 이용해서 URLPattern에 존재하는 파라미터 설정을 실제 값으로 변경
func (r *Request) GeneratePath(urlPattern string) {
	// 전달된 Path Parameter가 존재하지 않는 경우
	if len(r.Params) == 0 {
		r.Path = urlPattern
		return
	}
	buff := []byte(urlPattern)
	for k, v := range r.Params {
		key := []byte{}
		key = append(key, "{{."...)
		key = append(key, k...)
		key = append(key, "}}"...)
		buff = bytes.Replace(buff, key, []byte(v), -1)
	}
	r.Path = string(buff)
}

// Clone - Request 복제 (단, Thread-safe가 아니므로 Thread-safe가 필요한 경우는 "CloneRequest" 사용)
func (r *Request) Clone() Request {
	return Request{
		IsBypass: r.IsBypass,
		Method:   r.Method,
		URL:      r.URL,
		Query:    r.Query,
		Path:     r.Path,
		Body:     r.Body,
		Params:   r.Params,
		Headers:  r.Headers,
	}
}

// ===== [ Private Functions ] =====

// ===== [ Public Functions ] =====

// CloneMapValues - map[string][]string 형식의 정보를 복제
func CloneMapValues(values map[string][]string) map[string][]string {
	m := make(map[string][]string, len(values))
	for k, v := range values {
		tmp := make([]string, len(v))
		copy(tmp, v)
		m[k] = tmp
	}
	return m
}

// CloneParams - 지정한 Parameter 정보 복제
func CloneParams(params map[string]string) map[string]string {
	m := make(map[string]string)
	for k, v := range params {
		m[k] = v
	}
	return m
}

// CloneRequest - 지정한 Request에 대한 Deep Copy를 처리
func CloneRequest(r *Request) *Request {
	clone := r.Clone()
	clone.Headers = CloneMapValues(r.Headers)
	clone.Params = CloneParams(r.Params)
	clone.Query = CloneMapValues(r.Query)
	return &clone
}
