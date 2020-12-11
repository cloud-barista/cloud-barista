// Package basic - BASIC Auth 기능 제공 패키지
package basic

import (
	"net/http"

	"github.com/cloud-barista/cb-apigw/restapigw/pkg/config"
	"github.com/cloud-barista/cb-apigw/restapigw/pkg/jwt/provider"
	"github.com/dgrijalva/jwt-go"
)

// ===== [ Constants and Variables ] =====
const ()

var ()

// ===== [ Types ] =====
type (
	// Provider - Basic Auth 처리를 위한 구조
	Provider struct {
		provider.Verifier
	}
)

// ===== [ Implementations ] =====

// Build - Auth Provider 구성
func (p *Provider) Build(conf *config.CredentialsConfig) provider.Provider {
	return &Provider{
		Verifier: provider.NewVerifierBasket(NewPasswordVerifier(userConfigToTeam(conf.Basic.Users))),
	}
}

// GetClaims - JWT 토큰의 Claim 정보 반환
func (p *Provider) GetClaims(hc *http.Client) (jwt.MapClaims, error) {
	return jwt.MapClaims{}, nil
}

// ===== [ Private Functions ] =====

// userConfigToTeam - 지정한 정보를 기준으로 사용자정보 반환
func userConfigToTeam(userConf map[string]string) []*user {
	users := []*user{}
	for u, p := range userConf {
		users = append(users, &user{
			Username: u,
			Password: p,
		})
	}
	return users
}

// init - 패키지 초기화
func init() {
	provider.Register("basic", &Provider{})
}

// ===== [ Public Functions ] =====
