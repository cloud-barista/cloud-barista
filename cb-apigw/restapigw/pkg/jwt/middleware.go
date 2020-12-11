// Package jwt -
package jwt

import (
	"net/http"

	"github.com/cloud-barista/cb-apigw/restapigw/pkg/logging"
	"github.com/cloud-barista/cb-apigw/restapigw/pkg/render"
)

// ===== [ Constants and Variables ] =====
const ()

var ()

// ===== [ Types ] =====
type (
	// Payload - JWT 연관 Context Key 구조
	Payload struct{}

	// User - 사용자 정보 구조
	User struct {
		Username string
		Email    string
	}

	// Middleware - JWT Middleware 설정 구조
	Middleware struct {
		Guard Guard
	}
)

// ===== [ Implementations ] =====

// Handler - Request 핸들러 처리
func (m *Middleware) Handler(h http.Handler) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		parser := Parser{m.Guard.ParserConfig}
		_, err := parser.ParseFromRequest(req)
		if nil != err {
			logging.GetLogger().WithError(err).Debug("[MIDDLEWARE] JWT > failed to parse the token")
			render.JSON(rw, http.StatusUnauthorized, err.Error())
			return
		}

		h.ServeHTTP(rw, req)
	}
}

// ===== [ Private Functions ] =====
// ===== [ Public Functions ] =====

// NewMiddleware - JWT 관련 Middleware 구성
func NewMiddleware(g Guard) *Middleware {
	return &Middleware{g}
}
