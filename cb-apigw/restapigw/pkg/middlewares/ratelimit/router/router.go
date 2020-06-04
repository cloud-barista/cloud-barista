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
	MaxRate       int64  `yaml:"maxRate"`
	ClientMaxRate int64  `yaml:"clientMaxRate"`
	Strategy      string `yaml:"strategy"`
	Key           string `yaml:"key"`
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
	if err := yaml.NewDecoder(buf).Decode(&conf); err != nil {
		return nil
	}

	return conf
}
