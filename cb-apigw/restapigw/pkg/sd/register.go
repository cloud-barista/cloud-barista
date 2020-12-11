// Package sd -
package sd

import (
	"github.com/cloud-barista/cb-apigw/restapigw/pkg/config"
	"github.com/cloud-barista/cb-apigw/restapigw/pkg/core/register"
)

// ===== [ Constants and Variables ] =====
const ()

var (
	subscriberFactories = initRegister()
)

// ===== [ Types ] =====
type (
	// Register - Service Discovery 관리 구조
	Register struct {
		data register.Untyped
	}
)

// ===== [ Implementations ] =====

// Get - 관리 중인 Subscriber들중에 지정한 이름의 Subscriber Factory 추출
func (r *Register) Get(name string) SubscriberFactory {
	tmp, ok := r.data.Get(name)
	if !ok {
		return FixedSubscrberFactory
	}
	sf, ok := tmp.(SubscriberFactory)
	if !ok {
		return FixedSubscrberFactory
	}
	return sf
}

// ===== [ Private Functions ] =====

// initRegister - Subscriber 들을 관리하기 위한 Register 초기화
func initRegister() *Register {
	return &Register{register.NewUntyped()}
}

// ===== [ Public Functions ] =====

// GetSubscriber - 지정한 Backend 설정 정보에 해당하는 Subscriber 반환
func GetSubscriber(bConf *config.BackendConfig) Subscriber {
	// TODO: Service Discovery 기능 구현시 식별자 설정 필요
	return subscriberFactories.Get("")(bConf)
}
