package httpcache

import (
	"bytes"
	"context"
	"net/http"

	"github.com/cloud-barista/cb-apigw/restapigw/pkg/config"
	"github.com/cloud-barista/cb-apigw/restapigw/pkg/transport/http/client"
	"github.com/gregjones/httpcache"
	"gopkg.in/yaml.v2"
)

// ===== [ Constants and Variables ] =====

var (
	// MWNamespace - Middleware 설정 식별자
	MWNamespace = "mw-httpcache"

	memTransport = httpcache.NewMemoryCacheTransport()
	memClient    = http.Client{Transport: memTransport}
)

// ===== [ Types ] =====

// Config - HTTPCache 운영을 위한 설정 구조
type Config struct {
	// httpcache 사용 여부
	Enabled bool `yaml:"enabled"`
}

// ===== [ Implementations ] =====

// ===== [ Private Functions ] =====

// ParseConfig - HTTPCache 운영을 위한 Configuration parsing 처리
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

// ===== [ Public Functions ] =====

// NewHTTPClient - 메모리 관리용 HTTP client 를 사용하는 HTTPClientFactory 생성, 설정이 없거나 enabled=false인 경우는 일반 http.Client 사용
func NewHTTPClient(bConf *config.BackendConfig) client.HTTPClientFactory {
	conf := ParseConfig(bConf.Middleware)
	if conf == nil || !conf.Enabled {
		return client.NewHTTPClient
	}
	return func(_ context.Context) *http.Client {
		return &memClient
	}
}
