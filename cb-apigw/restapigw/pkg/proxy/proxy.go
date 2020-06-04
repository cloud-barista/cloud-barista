// Package proxy - Proxy 기능과 사용될 Middleware 인터페이스 정의와 구현체들을 제공
package proxy

import (
	"context"
	"errors"
	"io"

	"github.com/cloud-barista/cb-apigw/restapigw/pkg/config"
)

// ===== [ Constants and Variables ] =====
const (
	// 미들웨어 식별 자
	MWNamespace = "mw-proxy"
)

var (
	// ErrNoBackends - Backend 미 지정 오류
	ErrNoBackends = errors.New("all endpoints must have at least one backend")
	// ErrTooManyBackends - 너무 많은 Backend 지정 오류
	ErrTooManyBackends = errors.New("too many backends for this proxy")
	// ErrTooManyProxies - ProxyChain 에 너무 많은 ProxyChain 설정 오류
	ErrTooManyProxies = errors.New("too many proxes for this proxy chain")
	// ErrNotEnoughProxies - Proxy 가 부족한 경우 오류
	ErrNotEnoughProxies = errors.New("not enough proxies for this endpoint")
)

// ===== [ Types ] =====

// readCloserWrapper - Context 기반의 종료나 취소 시점에 닫을 수 있는 io.Reader 관리를 위한 구조 정의
type readCloserWrapper struct {
	ctx context.Context
	rc  io.ReadCloser
}

// Proxy - Context 기반에서 Request를 처리하고 Response 또는 error 를 반환하는 함수 형식
type Proxy func(ctx context.Context, req *Request) (*Response, error)

// Metadata - http.Response의 Headers와 StatusCode 관리를 위한 Metadata 구조 정의
type Metadata struct {
	Headers    map[string][]string
	StatusCode int
}

// Response - Backend 처리를 통해 반환된 http.Response 처리를 위한 Response 구조 정의
type Response struct {
	Data       map[string]interface{}
	IsComplete bool
	Metadata   Metadata
	Io         io.Reader
}

// BackendFactory - 지정된 BackendConfig 정보를 기반으로 수행할 Proxy를 반환하는 함수 형식
type BackendFactory func(bConf *config.BackendConfig) Proxy

// CallChain - 설정에 따른 중첩 Proxy 연계 구조를 위한 함수 형식
type CallChain func(next ...Proxy) Proxy

// ===== [ Implementations ] =====

// Read - 지정된 데이터 읽기 구현
func (rcw readCloserWrapper) Read(bytes []byte) (int, error) {
	return rcw.rc.Read(bytes)
}

// closeOnCancel - Context가 취소나 종료된 경우에 io.Reader를 닫기 위한 Wrrapper 처리
func (rcw readCloserWrapper) closeOnCancel() {
	<-rcw.ctx.Done()
	rcw.rc.Close()
}

// ===== [ Private Functions ] =====

// ===== [ Public Functions ] =====

// EmptyChain - 테스트나 오류 처리를 위한 빈 Proxy Chain 생성
func EmptyChain(next ...Proxy) Proxy {
	if len(next) > 1 {
		panic(ErrTooManyProxies)
	}
	return next[0]
}

// NewReadCloserWrapper - 닫기가 가능한 io.Reader 생성
func NewReadCloserWrapper(ctx context.Context, rc io.ReadCloser) io.Reader {
	wrapper := readCloserWrapper{ctx, rc}
	go wrapper.closeOnCancel()
	return wrapper
}

// DummyProxy - 테스트나 오류 방지를 위한 Dummcy Proxy 생성
func DummyProxy(_ context.Context, _ *Request) (*Response, error) {
	return nil, nil
}
