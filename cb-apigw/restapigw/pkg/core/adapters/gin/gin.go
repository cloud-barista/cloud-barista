// Package gin - GIN 기반으로 표준 HTTP 연계 기능을 제공하는 패키지
package gin

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ===== [ Constants and Variables ] =====
const ()

var ()

// ===== [ Types ] =====
type (
	// connectHandler - 다음 처리 handler 연계용 구조
	connectHandler struct{}

	// middlewareCtx - Gin Context와 다음 호출 여부 관리 구조
	middlewareCtx struct {
		ctx         *gin.Context
		childCalled bool
	}
)

// ===== [ Implementations ] =====

// ServeHTTP - Request로 부터 Gin Context를 추출해서 다음 호출을 처리
func (ch *connectHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	state := req.Context().Value(ch).(*middlewareCtx)
	defer func(req *http.Request) { state.ctx.Request = req }(state.ctx.Request)
	state.ctx.Request = req
	state.childCalled = true
	state.ctx.Next()
}

// ===== [ Private Functions ] =====
// ===== [ Public Functions ] =====

// New - Gin Engine에서 Middleware를 사용하기 위한 형식으로 HTTP Handler와 gin Handler 인스턴스 생성
func New() (http.Handler, func(h http.Handler) gin.HandlerFunc) {
	nextHandler := new(connectHandler)
	makeGinHandler := func(h http.Handler) gin.HandlerFunc {
		return func(c *gin.Context) {
			state := &middlewareCtx{ctx: c}
			ctx := context.WithValue(c.Request.Context(), nextHandler, state)
			h.ServeHTTP(c.Writer, c.Request.WithContext(ctx))
			if !state.childCalled {
				c.Abort()
			}
		}
	}
	return nextHandler, makeGinHandler
}

// Wrap - Gin Middleware Handler를 표준 HTTP Middleware 형식으로 구성
func Wrap(f func(h http.Handler) http.HandlerFunc) gin.HandlerFunc {
	next, adapter := New()
	return adapter(f(next))
}

// URLParam - Gin Context 기반의 URL Parameter 정보 추출
func URLParam(req *http.Request, key string) string {
	nextHandler := new(connectHandler)
	state := req.Context().Value(nextHandler).(*middlewareCtx)
	if nil != state {
		return state.ctx.Param(key)
	}
	return ""
}
