// Package observability - Request ID 추적을 제공하는 패키지
package observability

import (
	"context"

	"github.com/cloud-barista/cb-apigw/restapigw/pkg/core"
)

// ===== [ Constants and Variables ] =====
// ===== [ Types ] =====
// ===== [ Implementations ] =====
// ===== [ Private Functions ] =====
// ===== [ Public Functions ] =====

// RequestIDToContext - 지정한 컨텍스트에 Request ID 정보 설정
func RequestIDToContext(ctx context.Context, requestID string) context.Context {
	if nil == ctx {
		panic("Can not put request ID to empty context")
	}

	return context.WithValue(ctx, core.RequestIDKey, requestID)
}

// RequestIDFromContext - 지정한 컨텍스트에서 Request ID를 추출 (없다면 "")
func RequestIDFromContext(ctx context.Context) string {
	if nil == ctx {
		panic("Can not get request ID from empty context")
	}

	if requestID, ok := ctx.Value(core.RequestIDKey).(string); ok {
		return requestID
	}

	return ""
}
