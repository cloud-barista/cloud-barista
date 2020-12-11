// Package sd - Service Discovery 기능 제공 패키지
package sd

import "github.com/cloud-barista/cb-apigw/restapigw/pkg/config"

// ===== [ Constants and Variables ] =====
const ()

var ()

// ===== [ Types ] =====
type (
	// FixedSubscriber - Backend Hosts 정보 관리 형식 (갱신되지 않는 형식)
	FixedSubscriber struct {
		mode  string
		hosts []*config.HostConfig
	}

	// Subscriber - Backend Hosts 정보 관리 인터페이스
	Subscriber interface {
		Mode() string
		Hosts() ([]*config.HostConfig, error)
	}

	// SubscriberFunc - Subcriber 사용을 위한 함수 어뎁터
	SubscriberFunc func() ([]*config.HostConfig, error)

	// SubscriberFactory - 지정된 Backend 설정에 따른 Subscriber 구성 팩토리 함수 형식
	SubscriberFactory func(*config.BackendConfig) Subscriber
)

// ===== [ Implementations ] =====

// Hosts - Subscriber에서 관리되는 Hosts 반환
func (sf SubscriberFunc) Hosts() ([]*config.HostConfig, error) { return sf() }

// Mode - Load Balancing Mode 반환
func (fs FixedSubscriber) Mode() string {
	return fs.mode
}

// Hosts - FixedSubcriber에서 관리되는 Hosts 반환
func (fs FixedSubscriber) Hosts() ([]*config.HostConfig, error) {
	return fs.hosts, nil
}

// ===== [ Private Functions ] =====
// ===== [ Public Functions ] =====

// FixedSubscrberFactory - 지정된 Backend 설정에 따른 Fixed Subscriber를 구성
func FixedSubscrberFactory(bConf *config.BackendConfig) Subscriber {
	return FixedSubscriber{mode: bConf.BalanceMode, hosts: bConf.Hosts}
}
