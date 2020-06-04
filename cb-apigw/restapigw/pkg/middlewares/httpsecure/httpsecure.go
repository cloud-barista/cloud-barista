package httpsecure

import (
	"bytes"

	"gopkg.in/yaml.v3"
	"github.com/cloud-barista/cb-apigw/restapigw/pkg/config"
)

// ===== [ Constants and Variables ] =====

const (
	// MWNamespace - 미들웨어 설정 식별자
	MWNamespace = "mw-httpsecure"
)

// ===== [ Types ] =====

// Config - HTTPSecure 운영을 위한 설정 구조
type Config struct {
	// 요청이 허용되는 도메인들
	AllowedHosts []string `yaml:"allowed_hosts"`
	// 요청에 대한 Proxy된 호스트 이름 값을 보유할 수 있는 Header의 키들
	HostsProxyHeaders []string `yaml:"host_proxy_headers"`
	// STS (Strict-Transport-Security) Header의 유지 시간 (기본값: 0, 0이면 Header 미 포함)
	STSSeconds int64 `yaml:"sts_seconds"`
	// X-Frame-Options Header 값을 재 정의할 값으로 FrameDeny 옵션보다 우선권이 있다.
	CustomFrameOptionsValue string `yaml:"custom_frame_options_value"`
	// Content-Security-Policy Header 값을 재 정의할 값
	ContentSecurityPolicy string `yaml:"content_security_policy"`
	// 위조된 인증서로 MITM 공격을 방지하기 위한 HPKP 구현
	PublicKey string `yaml:"public_key"`
	// HTTP Request를 HTTPS 로 Redirection하는데 사용할 Host 명
	SSLHost string `yaml:"ssl_host"`
	// Referrer-Policy 헤더를 재 정의할 값
	ReferrerPolicy string `yaml:"referrer_policy"`
	// true로 설정되면 'nosniff' 값을 X-Content-Type-Options Header에 추가 (기본값: false)
	ContentTypeNosniff bool `yaml:"content_type_nosniff"`
	// true로 설정되면 1 값을 X-XSS-Protection Header에 추가 (기본값: false)
	BrowserXSSFilter bool `yaml:"browser_xss_filter"`
	// AllowedHosts, SSLRedirect 및 STSSeconds / STSIncludeSubdomains 옵션을 무시하고 개발 모드로 처리 (Production 에서는 반드시 false로 변경 필요) (기본값: false)
	IsDevelopment bool `yaml:"is_development"`
	// true로 설정되면 `includeSubdomains` 값을 Strict-Transport-Security Header에 추가 (기본값: false)
	STSIncludeSubdomains bool `yaml:"sts_include_subdomains"`
	// true로 설정되면 'DENY' 값을 X-Frame-Options Header에 추가 (기본값: false)
	FrameDeny bool `yaml:"frame_deny"`
	// true로 설정되면 HTTPS 요청만 허용 (기본값: false)
	SSLRedirect bool `yaml:"ssl_redirect"`
}

// ===== [ Implementations ] =====

// ===== [ Private Functions ] =====

// ===== [ Public Functions ] =====

// ParseConfig - HTTPSecure 운영을 위한 Configuration parsing 처리
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
