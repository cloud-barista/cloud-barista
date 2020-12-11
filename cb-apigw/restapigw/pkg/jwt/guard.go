// Package jwt - JWT 관련 기능을 제공하는 패키지
package jwt

import (
	"time"

	"github.com/cloud-barista/cb-apigw/restapigw/pkg/config"
)

// ===== [ Constants and Variables ] =====
// ===== [ Types ] =====

type (
	// Guard - JWT 인증 기반 구조
	Guard struct {
		ParserConfig
		// JWT 인증 유효 시간 (기본값: 1h)
		Timeout time.Duration
		// 신규 토큰 서명을 정보 (알고리즘/키페어)
		SigningMethod SigningMethod
		// 최대 갱신 시간(기본값: 24h, 0 - 갱신하지 않음)
		MaxRefresh time.Duration
	}
)

// ===== [ Implementations ] =====
// ===== [ Private Functions ] =====
// ===== [ Public Functions ] =====

// NewGuard - 지정한 정보를 기준으로 인증 처리기 생성
func NewGuard(cc *config.CredentialsConfig) Guard {
	return Guard{
		ParserConfig: ParserConfig{
			SigningMethods: []SigningMethod{{Alg: cc.Algorithm, Key: cc.Secret}},
			TokenLookup:    "header:Authorization",
		},
		SigningMethod: SigningMethod{Alg: cc.Algorithm, Key: cc.Secret},
		Timeout:       cc.TokenTimeout,
		MaxRefresh:    time.Hour * 24,
	}
}
