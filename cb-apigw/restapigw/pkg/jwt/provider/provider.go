// Package provider - JWT Provider 기능 제공 패키지
package provider

import (
	"net/http"
	"sync"

	"github.com/cloud-barista/cb-apigw/restapigw/pkg/config"
	"github.com/dgrijalva/jwt-go"
)

// ===== [ Constants and Variables ] =====
const ()

var (
	providers *sync.Map
)

// ===== [ Types ] =====
type (
	// Provider - Auth Provider 인터페이스
	Provider interface {
		Verifier
		Build(conf *config.CredentialsConfig) Provider
		GetClaims(hc *http.Client) (jwt.MapClaims, error)
	}

	// Factory - Provider들에 대한 Factory 구조
	Factory struct{}
)

// ===== [ Implementations ] =====

// Build - 지정한 정보를 기준으로 Auth 처리용 Provider 구성
func (f *Factory) Build(providerName string, cc *config.CredentialsConfig) Provider {
	provider, ok := providers.Load(providerName)
	if !ok {
		provider, _ = providers.Load("basic")
	}

	p := provider.(Provider)
	return p.Build(cc)
}

// ===== [ Private Functions ] =====

// init - 패키지 로드 시점 초기화
func init() {
	providers = new(sync.Map)
}

// ===== [ Public Functions ] =====

// Register - 지정한 정보로 Auth Provider 등록
func Register(name string, provider Provider) {
	providers.Store(name, provider)
}

// GetProviders - 관리 중인 Auth Provider들 반환
func GetProviders() *sync.Map {
	return providers
}
