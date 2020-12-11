// Package proxy - Backend 호출 구간에 대한 Rate Limit 적용 패키지
package proxy

import (
	"bytes"
	"context"

	"github.com/cloud-barista/cb-apigw/restapigw/pkg/config"
	"github.com/cloud-barista/cb-apigw/restapigw/pkg/logging"
	"github.com/cloud-barista/cb-apigw/restapigw/pkg/middlewares/ratelimit"
	"github.com/cloud-barista/cb-apigw/restapigw/pkg/proxy"
	"gopkg.in/yaml.v3"
)

// ===== [ Constants and Variables ] =====

const (
	// MWNamespace - Middleware configuration 식별자
	MWNamespace = "mw-ratelimit"
)

var (
	logger = logging.NewLogger()
)

// ===== [ Types ] =====

// Config - Rate Limit 구성을 위한 Configuration 구조
type Config struct {
	MaxRate  float64 `yaml:"max_rate"` // 초당 하용할 요청 수
	Capacity int64   `yaml:"capacity"` // 초당 허용할 최대 요청 수
}

// ===== [ Implementations ] =====
// ===== [ Private Functions ] =====
// ===== [ Public Functions ] =====

// ParseConfig - HTTPCache 운영을 위한 Configuration parsing 처리
func ParseConfig(mwConf config.MWConfig) *Config {
	conf := new(Config)
	tmp, ok := mwConf[MWNamespace]
	if !ok {
		return nil
	}

	buf := new(bytes.Buffer)
	yaml.NewEncoder(buf).Encode(tmp)
	if err := yaml.NewDecoder(buf).Decode(conf); nil != err {
		return nil
	}

	return conf
}

// NewBackendLimiter - Backend 호출에 대한 Rate Limit 기능을 제공하는 Middleware 생성
func NewBackendLimiter(bConf *config.BackendConfig) proxy.CallChain {
	conf := ParseConfig(bConf.Middleware)
	if nil == conf || 0 >= conf.MaxRate {
		return proxy.EmptyChain
	}
	backendLimiter := ratelimit.NewLimiterWithRate(conf.MaxRate, conf.Capacity)
	return func(next ...proxy.Proxy) proxy.Proxy {
		if 1 < len(next) {
			panic(proxy.ErrTooManyProxies)
		}
		return func(ctx context.Context, req *proxy.Request) (*proxy.Response, error) {
			logger.Debugf("[Backend Process Flow] RateLimit > CallChain (%s)", req.Path)
			// TokenBucket 검증
			if !backendLimiter.Allow() {
				logger.Debugf("[Backend Process Flow] RateLimit > CallChain (%s) ::: STOPPED!!", req.Path)
				return nil, ratelimit.ErrProxyLimited
			}
			logger.Debugf("[Backend Process Flow] RateLimit > CallChain (%s) ::: CONTINUE TO NEXT STEP!!", req.Path)
			return next[0](ctx, req)
		}
	}
}

// BackendFactory - Proxy에서 운영될 Rate Limit 기능이 적용된 Backend Factory 생성
func BackendFactory(next proxy.BackendFactory) proxy.BackendFactory {
	return func(bConf *config.BackendConfig) proxy.Proxy {
		return NewBackendLimiter(bConf)(next(bConf))
	}
}
