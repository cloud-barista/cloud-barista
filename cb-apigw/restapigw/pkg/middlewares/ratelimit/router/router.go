// Package router - Endpoint 처리 구간에 Rate Limit 기능을 제공하는 Router 패키지
package router

import (
	"bytes"

	"github.com/cloud-barista/cb-apigw/restapigw/pkg/config"
	"gopkg.in/yaml.v3"
)

// ===== [ Constants and Variables ] =====

const (
	// MWNamespace - Middleware configuration 식별자
	MWNamespace = "mw-ratelimit"
)

// ===== [ Types ] =====

// Config - Rate Limit 적용을 위한 Router Middleware Configuration 구조
type Config struct {
	MaxRate       int64  `yaml:"max_rate"`        // 초당 허용할 요청 수
	ClientMaxRate int64  `yaml:"client_max_rate"` // 클라이언트별 초당 허용할 요청 수
	Strategy      string `yaml:"strategy"`        // 클라이언트 식별 방법 ('ip', 'header')
	Key           string `yaml:"key"`             // 클라이언트를 header로 식별할 경우의 header key 값 (ex. X-Private-Token, ...)
}

// ===== [ Implementations ] =====
// ===== [ Private Functions ] =====
// ===== [ Public Functions ] =====

// ParseConfig - RateLimit 운영을 위한 Configuration parser 처리
func ParseConfig(eConf *config.EndpointConfig) *Config {
	conf := new(Config)
	tmp, ok := eConf.Middleware[MWNamespace]
	if !ok {
		return nil
	}

	buf := new(bytes.Buffer)
	yaml.NewEncoder(buf).Encode(tmp)
	if err := yaml.NewDecoder(buf).Decode(&conf); nil != err {
		return nil
	}

	return conf
}
