// Package jwt -
package jwt

import (
	"github.com/cloud-barista/cb-apigw/restapigw/pkg/core"
	"github.com/cloud-barista/cb-apigw/restapigw/pkg/errors"

	"github.com/dgrijalva/jwt-go"
)

// ===== [ Constants and Variables ] =====
const ()

var ()

// ===== [ Types ] =====
type (
	// AppClaims - JWT 검증을 위한 Claim 처리 구조
	AppClaims struct {
		jwt.MapClaims
		leeway int64
	}
)

// ===== [ Implementations ] =====

// UnmarshalJSON - AppClaims에 대한 Claims를 unmarshal 처리
func (ac *AppClaims) UnmarshalJSON(text []byte) error {
	return core.JSONUnmarshal(text, &ac.MapClaims)
}

// Valid - 현재는 시간정보(exp, iat, nbf, ...Valid)를 기준으로 하며, 이런 정보가 없다고 해도 Valid 한 것으로 처리
func (ac *AppClaims) Valid() error {
	vErr := new(jwt.ValidationError)
	now := jwt.TimeFunc().Unix()

	if !ac.VerifyExpiresAt(now, false) {
		vErr.Inner = errors.New("token is expired")
		vErr.Errors |= jwt.ValidationErrorExpired
	}

	if !ac.VerifyIssuedAt(now, false) {
		vErr.Inner = errors.New("token used before issued")
		vErr.Errors |= jwt.ValidationErrorIssuedAt
	}

	if !ac.VerifyNotBefore(now, false) {
		vErr.Inner = errors.New("token is not valid yet")
		vErr.Errors |= jwt.ValidationErrorNotValidYet
	}

	if 0 == vErr.Errors {
		return nil
	}

	return vErr
}

// VerifyExpiresAt - 만기여부와 Leeway 검증을 위해 jwt.StandardClaims.VerifyExpiresAt() 재 정의
func (ac *AppClaims) VerifyExpiresAt(cmp int64, req bool) bool {
	return ac.MapClaims.VerifyExpiresAt(cmp-ac.leeway, req)
}

// VerifyIssuedAt - 발금과 Leeway 검증을 위해 jwt.StandardClaims.VerifyIssuedAt() 재 정의
func (ac *AppClaims) VerifyIssuedAt(cmp int64, req bool) bool {
	return ac.MapClaims.VerifyIssuedAt(cmp+ac.leeway, req)
}

// VerifyNotBefore 검증여부와 Leeway 검증을 위해 jwt.StandardClaims.VerifyNotBefore() 재 정의
func (ac *AppClaims) VerifyNotBefore(cmp int64, req bool) bool {
	return ac.MapClaims.VerifyNotBefore(cmp+ac.leeway, req)
}

// ===== [ Private Functions ] =====
// ===== [ Public Functions ] =====

// NewAppClaims - ADMIN 에서 사용할 Claims 생성
func NewAppClaims(leeway int64) *AppClaims {
	return &AppClaims{MapClaims: jwt.MapClaims{}, leeway: leeway}
}
