// Package jwt -
package jwt

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

// ===== [ Constants and Variables ] =====
const ()

var ()

// ===== [ Types ] =====
type (
	// AccessToken - Token 정보 구조
	AccessToken struct {
		Type    string `yaml:"token_type" json:"token_type"`
		Token   string `yaml:"access_token" json:"access_token"`
		Expires int64  `yaml:"expires_in" json:"expires_in"`
	}
)

// ===== [ Implementations ] =====
// ===== [ Private Functions ] =====
// ===== [ Public Functions ] =====

// IssueAdminToken - API 접근을 위한 ADMIN JWT 발급
func IssueAdminToken(signingMethod SigningMethod, claims jwt.MapClaims, expireIn time.Duration) (*AccessToken, error) {
	// 지정한 알고리즘으로 토큰 생성
	token := jwt.New(jwt.GetSigningMethod(signingMethod.Alg))
	exp := time.Now().Add(expireIn).Unix()

	token.Claims = claims
	claims["exp"] = exp
	claims["iat"] = time.Now().Unix()

	// 토큰 서명
	accessToken, err := token.SignedString([]byte(signingMethod.Key))
	if nil != err {
		return nil, err
	}

	// 토큰 정보 반환
	return &AccessToken{
		Type:    "Bearer",
		Token:   accessToken,
		Expires: exp,
	}, nil
}
