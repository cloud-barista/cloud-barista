package proxy

import (
	"github.com/cloud-barista/cb-apigw/restapigw/pkg/config"
	"github.com/cloud-barista/cb-apigw/restapigw/pkg/logging"
)

// ===== [ Constants and Variables ] =====

// ===== [ Types ] =====

// defaultFactory - logging과 BackendFactory로 구성된 기본 팩토리 구조
type defaultFactory struct {
	backendFactory BackendFactory
	logger         logging.Logger
}

// Factory - 지정된 Endpoint 기준으로 Proxy 호출을 위한 함수를 생성하는 팩토리 인터페이스
type Factory interface {
	New(eConf *config.EndpointConfig) (Proxy, error)
}

// FactoryFunc - 지정된 Endpoint 기준으로 Proxy 호출을 위한 함수 Adapter
type FactoryFunc func(*config.EndpointConfig) (Proxy, error)

// ===== [ Implementations ] =====

// New - Proxy 를 생성하는 인터페이스 구현
func (f FactoryFunc) New(eConf *config.EndpointConfig) (Proxy, error) { return f(eConf) }

// newMulti - 여러 Backend로 구성되고 반환 결과를 Merge 처리하는 Proxy 구성
func (df defaultFactory) newMulti(eConf *config.EndpointConfig) (p Proxy, err error) {
	backendProxy := make([]Proxy, len(eConf.Backend))
	for i, backend := range eConf.Backend {
		backendProxy[i] = df.newStack(backend)
	}
	p = NewMergeDataChain(eConf)(backendProxy...)
	return
}

// newSingle - 단일 Backend로 구성되는 Proxy 구성
func (df defaultFactory) newSingle(eConf *config.EndpointConfig) (Proxy, error) {
	return df.newStack(eConf.Backend[0]), nil
}

// newStack - Backend 호출을 위한 Proxy 구성
func (df defaultFactory) newStack(bConf *config.BackendConfig) (p Proxy) {
	p = df.backendFactory(bConf)

	// TODO: Loadbalancer


	// Backend 호출을 위한 Request Call chain 구성
	p = NewRequestBuilderChain(bConf)(p)
	return
}

// New - 지정된 Endpoint 설정을 기준으로 동작하는 Proxy 생성 (Backend 설정 갯수에 따라서 single/multi 구분)
func (df defaultFactory) New(cfg *config.EndpointConfig) (p Proxy, err error) {
	switch len(cfg.Backend) {
	case 0:
		err = ErrNoBackends
	case 1:
		p, err = df.newSingle(cfg)
	default:
		p, err = df.newMulti(cfg)
	}
	if err != nil {
		return
	}

	// TODO: Static Content Middleware
	return
}

// ===== [ Private Functions ] =====

// ===== [ Public Functions ] =====

// NewDefaultFactory - 전달된 BackendFactory를 사용하는 기본 ProxyFactory 반환
func NewDefaultFactory(bf BackendFactory, logger logging.Logger) Factory {
	return defaultFactory{
		backendFactory: bf,
		logger:         logger,
	}
}
