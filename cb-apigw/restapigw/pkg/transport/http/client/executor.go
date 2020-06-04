package client

import (
	"context"
	"net/http"
)

// ===== [ Constants and Variables ] =====

// ===== [ Types ] =====

// HTTPRequestExecutor - 지정한 Context 기반에서 HTTP를 사용해서 Backend와 내부 API 들과 통신을 위한 함수 정의
type HTTPRequestExecutor func(ctx context.Context, req *http.Request) (*http.Response, error)

// HTTPClientFactory - 지정한 Context 기반에서 동작하는 http client 생성을 위한 함수 정의
type HTTPClientFactory func(ctx context.Context) *http.Client

// ===== [ Implementations ] =====

// ===== [ Private Functions ] =====

// ===== [ Public Functions ] =====

// DefaultHTTPRequestExecutor - 지정한 HTTPClientFactory를 통해서 HTTPRequestExecutor 생성
func DefaultHTTPRequestExecutor(hcf HTTPClientFactory) HTTPRequestExecutor {
	return func(ctx context.Context, req *http.Request) (*http.Response, error) {
		return hcf(ctx).Do(req.WithContext(ctx))
	}
}

// NewHTTPClient - 기본 HTTP Client 생성
func NewHTTPClient(ctx context.Context) *http.Client {
	return &http.Client{}
}
