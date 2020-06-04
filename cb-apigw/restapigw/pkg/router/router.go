package router

import (
	http "github.com/cloud-barista/cb-apigw/restapigw/pkg/transport/http/server"
)

// ===== [ Constants and Variables ] =====

const (
	// HeaderCompleteResponseValue - 정상적으로 종료된 Response들에 대한 Complete Header ("X-CB-RESTAPIGW-COMPLETED")로 설정될 값
	HeaderCompleteResponseValue = http.HeaderCompleteResponseValue
	// HeaderIncompleteResponseValue - 비 정상적으로 종료된 Response들에 대한 Complete Header로 설정될 값
	HeaderIncompleteResponseValue = http.HeaderIncompleteResponseValue
)

var (
	// MessageResponseHeaderName - Response 오류인 경우 클라이언트에 표시할 Header 명
	MessageResponseHeaderName = http.MessageResponseHeaderName
	// CompleteResponseHeaderName - 정상/비정상 종료에 대한 Header 정보를 클라이언트에 알리기 위한 Header 명
	CompleteResponseHeaderName = http.CompleteResponseHeaderName
	// HeadersToSend - Route로 전달된 Request에서 Proxy로 전달할 설정에 지정된 Header들 정보
	HeadersToSend = http.HeadersToSend
	// HeadersToNotSend - Router로 전달된 Request에서 Proxy로 전달죄지 않을 Header들 정보
	HeadersToNotSend = http.HeadersToNotSend
	// UserAgentHeaderValue - Proxy Request에 설정할 User-Agent Header 값
	UserAgentHeaderValue = http.UserAgentHeaderValue

	// DefaultToHTTPError - 항상 InternalServerError로 처리하는 기본 오류
	DefaultToHTTPError = http.DefaultToHTTPError
	// ErrInternalError - 문제가 발생했을 때 InternalServerError로 표현되는 오류
	ErrInternalError = http.ErrInternalError
)

// ===== [ Types ] =====

// ToHTTPError - HTTP StatusCode에 따라서 처리하는 오류 형식
type ToHTTPError http.ToHTTPError

// ===== [ Implementations ] =====

// ===== [ Private Functions ] =====

// ===== [ Public Functions ] =====
