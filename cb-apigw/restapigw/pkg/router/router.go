package router

import (
	"net/http"
	"sync"

	"github.com/cloud-barista/cb-apigw/restapigw/pkg/config"
	"github.com/cloud-barista/cb-apigw/restapigw/pkg/core"
	"github.com/cloud-barista/cb-apigw/restapigw/pkg/errors"
	"github.com/cloud-barista/cb-apigw/restapigw/pkg/proxy"
)

// ===== [ Constants and Variables ] =====

const (
	// HeaderCompleteResponseValue - 정상적으로 종료된 Response들에 대한 Complete Header ("X-CB-RESTAPIGW-COMPLETED")로 설정될 값
	HeaderCompleteResponseValue = "true"
	// HeaderIncompleteResponseValue - 비 정상적으로 종료된 Response들에 대한 Complete Header로 설정될 값
	HeaderIncompleteResponseValue = "false"
)

var (
	// MessageResponseHeaderName - Response 오류인 경우 클라이언트에 표시할 Header 명
	MessageResponseHeaderName = "X-" + core.AppName + "-Messages"
	// CompleteResponseHeaderName - 정상/비정상 종료에 대한 Header 정보를 클라이언트에 알리기 위한 Header 명
	CompleteResponseHeaderName = "X-" + core.AppName + "-Completed"
	// HeadersToSend - Route로 전달된 Request에서 Proxy로 전달할 설정에 지정된 Header들 정보
	HeadersToSend = []string{"Content-Type"}
	// HeadersToNotSend - Router로 전달된 Request에서 Proxy로 전달죄지 않을 Header들 정보
	HeadersToNotSend = []string{"Accept-Encoding"}
	// UserAgentHeaderValue - Proxy Request에 설정할 User-Agent Header 값
	UserAgentHeaderValue = []string{core.AppUserAgent}

	// DefaultToHTTPError - 항상 InternalServerError로 처리하는 기본 오류
	DefaultToHTTPError = func(_ error) int { return http.StatusInternalServerError }
	// ErrInternalError - 문제가 발생했을 때 InternalServerError로 표현되는 오류
	ErrInternalError = errors.New("internal server error")
)

// ===== [ Types ] =====

type (
	// HandlerFactory - 사용할 Router에 Proxy 연계를 위한 함수
	HandlerFactory func(*config.EndpointConfig, proxy.Proxy) http.HandlerFunc

	// ToHTTPError - HTTP StatusCode에 따라서 처리하는 오류 형식
	ToHTTPError func(error) int

	// Constructor - Middleware 형식의 운영을 위한 함수 형식
	Constructor func(http.Handler) http.Handler

	// Router - Router 운영에 필요한 인페이스
	Router interface {
		Engine() http.Handler
		UpdateEngine(sConf *config.ServiceConfig)
		RegisterAPIs(sConf *config.ServiceConfig, defs []*config.EndpointConfig) error
	}

	// DynamicRouter - 동적 라우팅 구성을 위한 Routing Engine 구조
	DynamicRouter struct {
		handler http.Handler
		mu      sync.Mutex
	}
)

// ===== [ Implementations ] =====

// ServeHTTP - HTTP 요청 처리
func (de *DynamicRouter) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	de.handler.ServeHTTP(rw, req)
}

// SetHandler - 지정한 Handler로 핸들러 교체
func (de *DynamicRouter) SetHandler(h http.Handler) {
	de.mu.Lock()
	defer de.mu.Unlock()

	de.handler = h
}

// GetHandler - 관리 중인 Handler 반환
func (de *DynamicRouter) GetHandler() http.Handler {
	return de.handler
}

// ===== [ Private Functions ] =====

// ===== [ Public Functions ] =====
