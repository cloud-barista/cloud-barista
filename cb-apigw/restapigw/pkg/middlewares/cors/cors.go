// Package cors - CORS (Cross Origin Resource Sharing) 을 지원하는 미들웨어 패키지
package cors

import (
	"bytes"
	"time"

	"github.com/cloud-barista/cb-apigw/restapigw/pkg/config"
	"gopkg.in/yaml.v3"
)

// ===== [ Constants and Variables ] =====

const (
	// MWNamespace - Middleware 설정 식별자
	MWNamespace = "mw-cors"
)

// ===== [ Types ] =====

// Config - Middleware 설정 정보 형식
type Config struct {
	AllowOrigins     []string      `yaml:"allow_origins"`
	AllowMethods     []string      `yaml:"allow_methods"`
	AllowHeaders     []string      `yaml:"allow_headers"`
	ExposeHeaders    []string      `yaml:"expose_headers"`
	AllowCredentials bool          `yaml:"allow_credentials"`
	MaxAge           time.Duration `yaml:"max_age"`
}

// ===== [ Implementations ] =====

// ===== [ Private Functions ] =====

// ===== [ Public Functions ] =====

// ParseConfig - CORS 운영을 위한 Configuration parsing 처리
func ParseConfig(mwConf config.MWConfig) *Config {
	conf := new(Config)
	tmp, ok := mwConf[MWNamespace]
	if !ok {
		return nil
	}

	buf := new(bytes.Buffer)
	yaml.NewEncoder(buf).Encode(tmp)
	if err := yaml.NewDecoder(buf).Decode(conf); err != nil {
		return nil
	}

	return conf
}
