// Package config - Configuration for Cloud-Barista's REST API Gateway and provides the required process
package config

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/cloud-barista/cb-apigw/restapigw/pkg/core"
	"github.com/cloud-barista/cb-apigw/restapigw/pkg/core/defaults"
	"github.com/cloud-barista/cb-apigw/restapigw/pkg/encoding"
	"github.com/cloud-barista/cb-apigw/restapigw/pkg/errors"
	"github.com/cloud-barista/cb-apigw/restapigw/pkg/logging"
)

// ===== [ Constants and Variables ] =====

const (
	// ConfigVersion - 설정 구조에 대한 버전
	ConfigVersion = 1
)

const (
	// Roundrobin - Roundrobin 방식의 Load Balancing
	Roundrobin LBModes = iota
	// Weight - Weight 방식의 Load Balancing
	Weight
)

var (
	jwtAlgorithms = []string{"HS256", "HS384", "HS512", "RS256", "RS384", "RS512"}
	encodings     = []string{"no-op", "json", "string"}

	errInvalidHost         = errors.New("invalid host")
	errInvalidNoOpEncoding = errors.New("can not use NoOp encoding with more than one backends connected to the same endpoint")

	// ErrNoHosts - Load Balancing 처리 대상 Host 가 지정되지 않은 경우 오류
	ErrNoHosts = errors.New("no available hosts")
)

// ===== [ Types ] =====

type (
	// IConfig - 설정관리용 struct의 기본 메서드 제공 인터페이스
	IConfig interface {
		InitializeDefaults()
		Validate()
	}

	// LBModes - Load Balancing Mode 유형 형식
	LBModes int

	// ServiceConfig - REST API Gateway 운영에 필요한 서비스 설정 형식
	ServiceConfig struct {
		// 서비스 식별 명 (기본값: '')
		Name string `mapstructure:"name"`
		// 기본 처리 시간 (기본값: 2s)
		Timeout time.Duration `mapstructure:"timeout" default:"2s"`
		// 종료시 잔여 요청을 처리하기 위한 대기 시간 (기본 값: 0s)
		GraceTimeout time.Duration `mapstructure:"grace_timeout"`
		// 디버그모드 여부 (기본값: false)
		Debug bool `mapstructure:"debug"`
		// GET 처리에 대한 캐시 TTL 기간 (기본값: 1h)
		CacheTTL time.Duration `mapstructure:"cache_ttl" default:"1h"`
		// 서비스에서 사용할 포트 (기본값: 8000)
		Port int `mapstructure:"port" default:"8000"`
		// 설정 파일 버전 (기본값: 1)
		Version int `mapstructure:"version" default:"1"`
		// 전체 요청을 읽기 위한 최대 허용 시간 (기본값: 0, 0이면 제한없음)
		ReadTimeout time.Duration `mapstructure:"read_timeout"`
		// 전체 응답을 출력하기 위한 최대 허용 시간 (기본값: 0, , 0이면 제한없음)
		WriteTimeout time.Duration `mapstructure:"write_timeout"`
		// Keep-alive 활성 상태에서 다음 요청까지의 최대 대기 시간 (기본값: 0, 0이면 제한없음)
		IdleTimeout time.Duration `mapstructure:"idle_timeout"`
		// 요청헤더를 읽기 위한 최대 허용 시간 (기본값: 0, 0이면 제한없음)
		ReadHeaderTimeout time.Duration `mapstructure:"read_header_timeout"`
		// 유휴연결(Keep-alive)들의 최대 유지 수 (기본값: 0, 0이면 제한없음)
		MaxIdleConnections int `mapstructure:"max_idle_connections"`
		// 호스트당 유휴연결(Keep-alive)들의 최대 유지 수 (기본값: 250, 0이면 250 사용)
		MaxIdleConnectionsPerHost int `mapstructure:"max_idle_connections_per_host" default:"250"`
		// 유휴연결(Keep-alive)의 최대 유효시간 (기본값: , 0이면 ReadTimeout 사용, 이것도 0이면 ReadHeaderTimeout 사용)
		IdleConnectionTimeout time.Duration `mapstructure:"idle_connection_timeout"`
		// 서비스 단위에서 적용할 Middleware 설정
		Middleware MWConfig `mapstructure:"middleware"`
		// 서비스에서 사용할 TLS 설정
		TLS *TLSConfig `mapstructure:"tls"`
		// TCP 연결에 사용할 대기시간 (기본값:0, 0이면 no timeout)
		DialerTimeout time.Duration `mapstructure:"dialer_timeout"`
		// 활성연결의 유지 시간 (기본값:없음, 0 지정시는 Keep-alive 비활성화)
		DialerKeepAlive time.Duration `mapstructure:"dialer_keep_alive"`
		// DualStack 활성화 시에 실패한 연결을 재 처리하는데 필요한 대기 시간 (기본값: 0, 0이면 no delay)
		DialerFallbackDelay time.Duration `mapstructure:"dialer_fallback_delay"`
		// 압축 비활성 여부 (기본값: false)
		DisableCompression bool `mapstructure:"disable_compression"`
		// 다른 요청에 TCP 연결을 재 사용하는 것의 비활성 여부 (기본값: false)
		DisableKeepAlives bool `mapstructure:"disable_keep_alives"`
		// Request 처리 후에 서버의 Response Header 정보를 기다리는 시간 (기본값: 0, 0이면 no timeout)
		ResponseHeaderTimeout time.Duration `mapstructure:"response_header_timeout"`
		// 서버의 첫번째 Response Header 정보를 기다리는 시간 (기본값: 0, 0이면 no timeout)
		ExpectContinueTimeout time.Duration `mapstructure:"expect_continue_timeout"`
		// DisableStrictREST - REST 강제 규칙 비활성화 여부 (기본값: false)
		DisableStrictREST bool `mapstructure:"disable_strict_rest"`
		// RouterEngine - Route 처리에 사용할 Engine 지정 (기본값: gin)
		RouterEngine string `mapstructure:"router_engine" default:"gin"`

		// Admin API 설정
		Admin *AdminConfig `mapstructure:"admin"`
		// Repository 설정
		Repository *RepositoryConfig `mapstructure:"repository"`
		// Cluster 설정
		Cluster *ClusterConfig `mapstructure:"cluster"`
	}

	// AdminConfig - Admin API 운영에 필요한 설정 형식
	AdminConfig struct {
		// Port - Admin Server 포트 (기본값: 8001)
		Port int `mapstructure:"port" default:"8001"`
		// Credentials - Admin Server (WEB)를 사용할 사용자 설정
		Credentials *CredentialsConfig `mapstructure:"credentials"`
		// TLS - Admin Server에서 사용할 TLS 설정
		TLS *TLSConfig `mapstructure:"tls"`
		// ProfilingEnabled - Admin Server Profile 처리 여부 (기본값: false)
		ProfilingEnabled bool `mapstructure:"profiling_enabled"`
		// ProfilingPublic - Admin Server Profile 정보 노출 여부 (기본값: false)
		ProfilingPublic bool `mapstructure:"profiling_public"`
	}

	// RepositoryConfig - Routing 정보 관리 형식
	RepositoryConfig struct {
		// DSN - Repository 연결 문자열 (기본값: "file://./conf", "cbstore://api/restapigw/conf" 설정 가능)
		DSN string `mapstructure:"dsn" default:"file://./conf"`
	}

	// ClusterConfig - Cluster 환경 정보 관리 형식
	ClusterConfig struct {
		// UpdateFrequency - Repository Polling 주기 (기본값: 10s, file repository가 아닌 경우)
		UpdateFrequency time.Duration `mapstructure:"update_frequency" default:"10s"`
	}

	// CredentialsConfig - Admin API 사용자 인증 정보 형식
	CredentialsConfig struct {
		// Algorithm - JWT 인증 알고리즘 (기본값: HS256)
		Algorithm string `mapstructure:"algorithm" default:"HS256"`
		// Secret - JWT 인증 비밀 키 (기본값: "testSecret")
		Secret string `mapstructure:"secret" default:"testSecret"`
		// TokenTimeout - JWT 인증 유효기간 (기본값: 3h)
		TokenTimeout time.Duration `mapstructure:"token_timeout" default:"3h"`
		// Basic - Admin 사용자 기본 인증 정보
		Basic *BasicAuthConfig `mapstructure:"basic"`
	}

	// BasicAuthConfig - Admin 사용자 기본 인증 정보 형식
	BasicAuthConfig struct {
		// Users - Admin 사용자 정보 (기본값: "admin:test@admin00")
		Users map[string]string `mapstructure:"users" default:"{\"admin\":\"test@admin00\"}"`
	}

	// EndpointConfig - 서비스 라우팅에 사용할 설정 구조
	EndpointConfig struct {
		// Name - 설정 식별 명
		Name string `yaml:"name" json:"name"`
		// Active - 설정 활성화 여부 (기본값: true)
		Active bool `yaml:"active" json:"active" default:"true"`
		// Endpoint - 클라이언트에 노출될 URL 패턴
		Endpoint string `yaml:"endpoint" json:"endpoint"`
		// Hosts - 전역으로 사용할 기본 Host 리스트 (기본값: "[]")
		Hosts []*HostConfig `yaml:"hosts" json:"hosts" default:"[]"`
		// Method - Endpoint에 대한 HTTP 메서드 (GET, POST, PUT, etc) (기본값: "GET"")
		Method string `yaml:"method" json:"method" default:"GET"`
		// Timeout - Endpoint 처리 시간 (기본값: "2s"")
		Timeout time.Duration `yaml:"timeout" json:"timeout" default:"2s"`
		// CacheTTL - GET 처리에 대한 캐시 TTL 기간 (기본값: "1h")
		CacheTTL time.Duration `yaml:"cache_ttl" json:"cache_ttl" default:"1h"`
		// OutputEncoding - 반환결과 처리에 사용할 인코딩 (기본값: "json")
		OutputEncoding string `yaml:"output_encoding" json:"output_encoding" default:"json"`
		// ExceptQueryStrings - Backend 에 전달되는 Query String에서 제외할 파라미터 Key 리스트
		ExceptQueryStrings []string `yaml:"except_querystrings" json:"except_querystrings" default:"[]"`
		// ExceptHeaders - Backend 에 전달되는 Header에서 제외할 파라미터 Key 리스트
		ExceptHeaders []string `yaml:"except_headers" json:"except_headers" default:"[]"`
		// Middleware - Endpoint 단위에서 적용할 Middleware 설정
		Middleware MWConfig `yaml:"middleware" json:"middleware"`
		// HealthCheck - Health Check 설정
		HealthCheck *HealthCheck `yaml:"health_check" json:"health_check"`
		// Backend - Endpoint에서 호출할 Backend API 서버 호출/응답 처리 설정 리스트
		Backend []*BackendConfig `yaml:"backend" json:"backend"`

		// Bypass 처리 여부 (내부 처리용)
		IsBypass bool `yaml:"-" json:"-"`
	}

	// BackendConfig - Backend API Server 연결과 응답 처리를 위한 설정 구조
	BackendConfig struct {
		// Hosts - Backend API Server의 Host URI (기본값: 없으면 EndPointConfig 정보 사용)
		Hosts []*HostConfig `yaml:"hosts" json:"hosts"`
		// Timeout - Backend 처리 시간 (기본값: 없으면 EndPointConfig 정보 사용)
		Timeout time.Duration `yaml:"timeout" json:"timeout"`
		// Method - Backend 호출에 사용할 HTTP Method (기본값: 없으면 EndPointConfig 정보 사용)
		Method string `yaml:"method" json:"method"`
		// URLPattern - Backend 호출에 사용할 URL Patthern (기본값: "")
		URLPattern string `yaml:"url_pattern" json:"url_pattern"`
		// Encoding - 인코딩 포맷 (기본값: "json")
		Encoding string `yaml:"encoding" json:"encoding" default:"json"`
		// Group - Backend 결과를 묶을 Group 명 (기본값: "")
		Group string `yaml:"group" json:"group"`
		// Blacklist - Backend 결과에서 생략할 필드명 리스트 (기본값: "[]", flatmap 적용 "." operation)
		Blacklist []string `yaml:"blacklist" json:"blacklist" default:"[]"`
		// Whitelist - Backend 결과에서 추출할 필드명 리스트 (기본값: "[]", flatmap 적용 "." operation)
		Whitelist []string `yaml:"whitelist" json:"whitelist" default:"[]"`
		// Mapping - Backend 결과에서 필드명을 변경할 리스트 맵 (기본값: "{}")
		Mapping map[string]string `yaml:"mapping" json:"mapping"`
		// IsCollection - Backend 결과가 컬랙션인지 여부 (기본값: false)
		IsCollection bool `yaml:"is_collection" json:"is_collection" default:"false"`
		// WrapCollectionToJSON - Backend 결과가 컬랙션인 경우에 core.CollectionTag ("collection") 으로 JSON 포맷을 할 것인지 여부
		// (True 면 core.CollectionTag ("collection") 으로 JSON 전환, false면 Array 상태로 반환) (기본값: false)
		WrapCollectionToJSON bool `yaml:"wrap_collection_to_json" json:"wrap_collection_to_json" default:"false"`
		// Target - Backend 결과 중에서 특정한 필드만 처리할 경우의 필드명 (기본값: "")
		Target string `yaml:"target" json:"target"`
		// Middleware - Backend 에서 동작할 Middleware 설정
		Middleware MWConfig `yaml:"middleware" json:"middleware"`
		// HostSanitizationDisabled - host 정보의 정제작업 비활성화 여부 (기본값: false)
		HostSanitizationDisabled bool `yaml:"disable_host_sanitize" json:"disable_host_sanitize" default:"false"`
		// BalanceMode - Backend Loadbalacing 모드 (기본값: "", "rr" - "roundrobin", "wrr" - "weighted roundrobin", "" - random)
		BalanceMode string `yaml:"lb_mode" json:"lb_mode" default:""`

		// API 호출의 응답을 파싱하기 위한 디코더 (내부 사용)
		Decoder encoding.Decoder `yaml:"-" json:"-"`
		// URLPattern에서 파라미터 변환에 사용할 키 관리 (내부 사용)
		URLKeys []string `yaml:"-" json:"-"`
	}

	// HostConfig - Backend Load balancing 처리를 위한 Host 구조
	HostConfig struct {
		// Host - Backend Service 호스트 정보 (기본값: "", 필수)
		Host string `mapstructure:"host"`
		// Weight - Weighted Roundrobin 선택 적용할 가중치 (기본값: 0)
		Weight int `mapstructure:"weight" default:"0"`
	}

	// TLSConfig - 서비스에서 사용할 TLS 설정 구조
	TLSConfig struct {
		// Port - 기본 포트 (기본값: 8443)
		Port int `mapstructure:"port"`
		// Redirect - TLS 리다이렉션 (기본값: true)
		Redirect bool `mapstructure:"redirect" default:"true"`
		// IsDiabled - TLS 비활성화 여부 (기본값: false)
		IsDisabled bool `mapstructure:"disabled"`
		// PublicKey - 공개 키 경로 (기본값: "")
		PublicKey string `mapstructure:"public_key"`
		// PrivateKey - 비밀 키 경로 (기본값: "")
		PrivateKey string `mapstructure:"private_key"`
		// MinVersion - TLS 최소 버전 (기본값: VersionTLS12)
		MinVersion string `mapstructure:"min_version"`
		// MaxVersion - TLS 최대 버전 (기본값: VersionTLS12)
		MaxVersion string `mapstructure:"max_version"`
		// Curve 설정들의 리스트 (use 23 for CurveP256, 24 for CurveP384 or 25 for CurveP521, 기본값: 모두 사용)
		CurvePreferences []uint16 `mapstructure:"curve_preferences"`
		// PreferServerCipherSuites - 서버에서 사용을 강제하는 Cipher Suite 리스트 (기본값: false)
		PreferServerCipherSuites bool `mapstructure:"prefer_server_cipher_suites"`
		// CipherSuites - Chiper Suite 리스트 (기본값: defaultCipherSuites 리스트 사용)
		CipherSuites []uint16 `mapstructure:"cipher_suites"`
	}

	// MWConfig - Middleware 설정을 저장하기 위한 맵 구조 (개별 Middlewares에서 설정 Parsing 적용)
	MWConfig map[string]interface{}

	// HealthCheck - Health Check 구조
	HealthCheck struct {
		// URL - Health Checking URL (기본값: "")
		URL string `mapstructure:"url" yaml:"url" json:"url" bson:"url"` // `mapstructure:"url" yaml:"url" json:"url" bson:"url" valid:"url"`
		// Timeout - 검증 제한 시간 (기본값: 0, 제한없음)
		Timeout time.Duration `mapstructure:"timeout" yaml:"timeout" json:"timeout" bson:"timeout"`
	}

	// UnsupportedVersionError - 설정 초기화 과정에서 버전 검증을 통해 반환할 오류 구조
	UnsupportedVersionError struct {
		Have int
		Want int
	}

	// WrongNumberOfParamsError - 파라미터의 IN/OUT 갯수가 다를 경우에 반환할 오류 구조
	WrongNumberOfParamsError struct {
		Endpoint     string
		Method       string
		Backend      int
		InputParams  []string
		OutputParams []string
	}

	// UndefinedOutputParamError - IN 파라미터에 대한 OUT 파라미터가 지정되지 않았을 경우에 반환할 오류 구조
	UndefinedOutputParamError struct {
		Endpoint     string
		Method       string
		Backend      int
		InputParams  []string
		OutputParams []string
		Param        string
	}

	// EndpointMatchError - Endpoint 패턴에 문제가 있을 경우에 반환할 오류 구조
	EndpointMatchError struct {
		Path   string
		Method string
		Err    error
	}

	// NoBackendsError - Backend가 지정되지 않았을 경우에 반환할 오류 구조
	NoBackendsError struct {
		Path   string
		Method string
	}

	// EndpointPathError - Endpoint 경로에 문제가 있을 경우에 반환할 오류 구조
	EndpointPathError struct {
		Path   string
		Method string
	}
)

// ===== [ Implementations ] =====

// IsHTTPS - HTTPS 활성화 여부 검증
func (t *TLSConfig) IsHTTPS() bool {
	return nil != t && "" != t.PublicKey && "" != t.PrivateKey
}

// InitializeDefaults - 기본 값 설정
func (t *TLSConfig) InitializeDefaults() error {
	// TLSConfig 초기화
	if err := defaults.Set(t); nil != err {
		return err
	}

	return nil
}

// Validate - 설정 검증
func (t *TLSConfig) Validate() error {
	return nil
}

// InitializeDefaults - 기본 값 설정
func (ba *BasicAuthConfig) InitializeDefaults() error {
	// BasicAuthConfig 초기화
	if err := defaults.Set(ba); nil != err {
		return err
	}

	return nil
}

// Validate - 설정 검증
func (ba *BasicAuthConfig) Validate() error {
	return nil
}

// InitializeDefaults - 기본값 설정
func (cc *CredentialsConfig) InitializeDefaults() error {
	// CredentialsConfig 초기화
	if err := defaults.Set(cc); nil != err {
		return err
	}

	if nil != cc.Basic {
		if err := cc.Basic.InitializeDefaults(); nil != err {
			return err
		}
	}

	return nil
}

// Validate - 설정 검증
func (cc *CredentialsConfig) Validate() error {
	if nil == cc.Basic {
		return errors.New("basic account for admin required")
	}
	if !core.ContainsString(jwtAlgorithms, cc.Algorithm) {
		return errors.New("invalid credential algorithm")
	}

	return nil
}

// InitializeDefaults - 기본값 설정
func (hc *HealthCheck) InitializeDefaults() error {
	// HealthCheck 초기화
	if err := defaults.Set(hc); nil != err {
		return err
	}

	return nil
}

// Validate - 설정 검증
func (hc *HealthCheck) Validate() error {
	return nil
}

// InitializeDefaults - 기본값 설정
func (h *HostConfig) InitializeDefaults() error {
	// HealthCheck 초기화
	if err := defaults.Set(h); nil != err {
		return err
	}

	return nil
}

// Validate - 설정 검증
func (h *HostConfig) Validate() error {
	return nil
}

// InitializeDefaults - 기본 값 설정
func (admin *AdminConfig) InitializeDefaults() error {
	// AdminConfig 초기화
	if err := defaults.Set(admin); nil != err {
		return err
	}

	if nil != admin.TLS {
		if err := admin.TLS.InitializeDefaults(); nil != err {
			return err
		}
	}

	if nil != admin.Credentials {
		if err := admin.Credentials.InitializeDefaults(); nil != err {
			return err
		}
	}

	return nil
}

// Validate - 설정 검증
func (admin *AdminConfig) Validate() error {
	if nil == admin.Credentials {
		return errors.New("crendential for admin required")
	}

	if err := admin.Credentials.Validate(); nil != err {
		return err
	}

	return nil
}

// InitializeDefaults - 기본 값 설정
func (r *RepositoryConfig) InitializeDefaults() error {
	// RepositoryConfig 초기화
	if err := defaults.Set(r); nil != err {
		return err
	}

	return nil
}

// Validate - 설정 검증
func (r *RepositoryConfig) Validate() error {
	if !strings.HasPrefix(r.DSN, "file://") && !strings.HasPrefix(r.DSN, "cbstore://") {
		return errors.New("invalid repository data soruce name format")
	}

	return nil
}

// InitializeDefaults - 기본 값 설정
func (c *ClusterConfig) InitializeDefaults() error {
	// ClusterConfig 초기화
	if err := defaults.Set(c); nil != err {
		return err
	}

	return nil
}

// Validate - 설정 검증
func (c *ClusterConfig) Validate() error {
	return nil
}

// Error - Endpoint 경로 오류 문자열 반환
func (e *EndpointPathError) Error() string {
	return "ERROR: the endpoint url path '" + e.Method + " " + e.Path + "' is not a valid one!!! Ignoring"
}

// Error - Backend 미지정 오류 문자열 반환
func (n *NoBackendsError) Error() string {
	return "WARNING: the '" + n.Method + " " + n.Path + "' endpoint has 0 backends defined! Ignoring"
}

// Error - Endpoint 패턴 오류 문자열 반환
func (e *EndpointMatchError) Error() string {
	return fmt.Sprintf("ERROR: parsing the endpoint url '%s %s': %s. Ignoring", e.Method, e.Path, e.Err.Error())
}

// Error - Output 파라미터 미지정 오류 문자열 반환
func (u *UndefinedOutputParamError) Error() string {
	return fmt.Sprintf(
		"Undefined output param '%s'! endpoint: %s %s, backend: %d. input: %v, output: %v",
		u.Param,
		u.Method,
		u.Endpoint,
		u.Backend,
		u.InputParams,
		u.OutputParams,
	)
}

// Error - IN/OUT 파라미터 갯수 문제 오류 문자열 반환
func (w *WrongNumberOfParamsError) Error() string {
	return fmt.Sprintf(
		"input and output params do not match. endpoint: %s %s, backend: %d. input: %v, output: %v",
		w.Method,
		w.Endpoint,
		w.Backend,
		w.InputParams,
		w.OutputParams,
	)
}

// sanitize - Middleware 구성을 정제해서 Namespace 기준의 맵으로 관리
func (mwConf *MWConfig) sanitize() {
	for module, mw := range *mwConf {
		switch mw := mw.(type) {
		case map[interface{}]interface{}:
			sanitized := map[string]interface{}{}
			for k, v := range mw {
				sanitized[fmt.Sprintf("%v", k)] = v
			}
			(*mwConf)[module] = sanitized
		}
	}
}

// Init - 설정에 대한 검사 및 초기화
func (sConf *ServiceConfig) Init() error {
	// 서비스 설정 초기화
	if err := sConf.InitializeDefaults(); nil != err {
		return err
	}

	return nil
}

// InitializeDefaults - 전역 설정 초기화
func (sConf *ServiceConfig) InitializeDefaults() error {
	// ServiceConfig 초기화
	if err := defaults.Set(sConf); nil != err {
		return err
	}

	if nil != sConf.TLS {
		// TLS 설정 초기화
		if err := sConf.TLS.InitializeDefaults(); nil != err {
			return err
		}
	}

	if nil != sConf.Admin {
		// Admin API 초기화
		if err := sConf.Admin.InitializeDefaults(); nil != err {
			return err
		}
	}

	if nil != sConf.Repository {
		// Repository 초기화
		if err := sConf.Repository.InitializeDefaults(); nil != err {
			return err
		}
	}

	if nil != sConf.Cluster {
		// Cluster 초기화
		if err := sConf.Cluster.InitializeDefaults(); nil != err {
			return err
		}
	}

	// 전역 설정에 대한 Middleware 구성 맵 관리
	sConf.Middleware.sanitize()

	return nil
}

// Validate - 설정 검증
func (sConf *ServiceConfig) Validate() error {
	// 설정 파일 버전 검증
	if sConf.Version != ConfigVersion {
		return &UnsupportedVersionError{
			Have: sConf.Version,
			Want: ConfigVersion,
		}
	}
	if "" == sConf.Name {
		return errors.New("service name reqired")
	}
	if nil == sConf.Admin {
		return errors.New("admin configuration required")
	}
	if nil == sConf.Repository {
		return errors.New("repository configuration required")
	}
	if nil != sConf.TLS {
		if err := sConf.TLS.Validate(); nil != err {
			return err
		}
	}
	if err := sConf.Admin.Validate(); nil != err {
		return err
	}
	if err := sConf.Repository.Validate(); nil != err {
		return err
	}
	if nil != sConf.Cluster {
		if err := sConf.Cluster.Validate(); nil != err {
			return err
		}
	}

	return nil
}

// InitializeDefaults - Endpoint에 미 설정된 항목들을 기본 값으로 초기화
func (eConf *EndpointConfig) InitializeDefaults() error {
	if err := defaults.Set(eConf); nil != err {
		return err
	}
	if nil != eConf.Hosts {
		for _, host := range eConf.Hosts {
			if err := host.InitializeDefaults(); nil != err {
				return err
			}
		}

	}
	if nil != eConf.HealthCheck {
		if err := eConf.HealthCheck.InitializeDefaults(); nil != err {
			return err
		}
	}
	if nil != eConf.Backend {
		for _, bConf := range eConf.Backend {
			if err := bConf.InitializeDefaults(); nil != err {
				return err
			}
		}
	}

	return nil
}

// AdjustValues - 설정 정보를 사용가능한 정보로 재 구성
func (eConf *EndpointConfig) AdjustValues(sConf *ServiceConfig) error {
	eConf.Endpoint = core.CleanPath(eConf.Endpoint)

	// Endpoint URL에 지정된 파라미터를 Input parameter로 설정
	inputParams := core.ExtractPlaceHoldersFromURLTemplate(eConf.Endpoint, core.ParameterExtractPattern(sConf.DisableStrictREST))
	inputSet := map[string]interface{}{}
	for ip := range inputParams {
		inputSet[inputParams[ip]] = nil
	}

	eConf.Endpoint = core.GetParameteredPath(eConf.Endpoint, inputParams)

	if eConf.OutputEncoding == encoding.NOOP && 1 < len(eConf.Backend) {
		return errInvalidNoOpEncoding
	}

	// Endpoint 설정의 Middleware 설정 맵 관리
	eConf.Middleware.sanitize()

	for bIdx, bConf := range eConf.Backend {
		if err := eConf.InitBackendDefaults(bIdx); nil != err {
			return err
		}

		if err := eConf.InitBackendURLMappings(bIdx, inputSet); nil != err {
			return err
		}

		// Backend 설정의 Middleware 설정 맵 관리
		bConf.Middleware.sanitize()
	}

	// 최종 Endpoint 정보 검사
	if err := eConf.Validate(); nil != err {
		return err
	}

	return nil
}

// InitBackendDefaults - Backend에 미 설정된 항목들을 기본 값으로 초기화
func (eConf *EndpointConfig) InitBackendDefaults(bIdx int) error {
	backend := eConf.Backend[bIdx]

	if err := backend.InitializeDefaults(); nil != err {
		return err
	}

	if 0 == len(backend.Hosts) {
		// URL 미 지정시 전역 URL 사용
		backend.Hosts = eConf.Hosts
	} else if !backend.HostSanitizationDisabled {
		cleanHosts(backend.Hosts)
	}
	// Method 미 지정시 Endpoint Method 사용
	if "" == backend.Method {
		backend.Method = eConf.Method
	}
	if 0 == backend.Timeout {
		backend.Timeout = eConf.Timeout
	}

	// Backend 처리 결과를 위한 Decoder 구성
	backend.Decoder = encoding.Get(strings.ToLower(backend.Encoding))(backend.IsCollection, backend.WrapCollectionToJSON)

	return nil
}

// InitBackendURLMappings - Backend에 지정된 파라미터 정보들을 이후에 사용할 수 있도록 초기화
func (eConf *EndpointConfig) InitBackendURLMappings(bIdx int, inputParams map[string]interface{}) error {
	backend := eConf.Backend[bIdx]

	backend.URLPattern = core.CleanPath(backend.URLPattern)
	// 파라미터 설정 추출 및 관리
	outputParams, outputSetSize := uniqueOutput(core.ExtractPlaceHoldersFromURLTemplate(backend.URLPattern, core.SimpleURLKeysPattern))

	ip := fromSetToSortedSlice(inputParams)

	if outputSetSize > len(ip) {
		return &WrongNumberOfParamsError{
			Endpoint:     eConf.Endpoint,
			Method:       eConf.Method,
			Backend:      bIdx,
			InputParams:  ip,
			OutputParams: outputParams,
		}
	}

	backend.URLKeys = []string{}
	for _, output := range outputParams {
		if !core.SequentialParamsPattern.MatchString(output) {
			if _, ok := inputParams[output]; !ok {
				return &UndefinedOutputParamError{
					Param:        output,
					Endpoint:     eConf.Endpoint,
					Method:       eConf.Method,
					Backend:      bIdx,
					InputParams:  ip,
					OutputParams: outputParams,
				}
			}
		}
		key := strings.Title(output)
		backend.URLPattern = strings.Replace(backend.URLPattern, "{"+output+"}", "{{."+key+"}}", -1)
		backend.URLKeys = append(backend.URLKeys, key)
	}
	return nil
}

// Validate - Endpoint 별 세부 필수 항목 검증
func (eConf *EndpointConfig) Validate() error {
	if "" == eConf.Name {
		return errors.New("endpoint name required")
	}

	matched, err := regexp.MatchString(core.DebugPattern, eConf.Endpoint)
	if nil != err {
		return &EndpointMatchError{
			Err:    err,
			Path:   eConf.Endpoint,
			Method: eConf.Method,
		}
	}
	if matched {
		return &EndpointPathError{Path: eConf.Endpoint, Method: eConf.Method}
	}

	if !core.ContainsString(encodings, eConf.OutputEncoding) {
		return errors.New("invalid output encoding")
	}

	if 0 == len(eConf.Backend) {
		return &NoBackendsError{Path: eConf.Endpoint, Method: eConf.Method}
	}
	for _, backend := range eConf.Backend {
		if err := backend.Validate(); nil != err {
			return err
		}
	}

	return nil
}

// InitializeDefaults - 설정 초기화
func (bConf *BackendConfig) InitializeDefaults() error {
	if err := defaults.Set(bConf); nil != err {
		return err
	}
	return nil
}

// Validate - 설정 검증
func (bConf *BackendConfig) Validate() error {
	if 0 == len(bConf.Hosts) {
		return ErrNoHosts
	}

	if "" == bConf.URLPattern {
		return errors.New("invalid backend url pattern")
	}

	if !core.ContainsString(encodings, bConf.Encoding) {
		return errors.New("invalid encoding for backend")
	}

	return nil
}

// Error - 비 호환 버전에 대한 오류 문자열 반환
func (u *UnsupportedVersionError) Error() string {
	return fmt.Sprintf("Unsupported version: %d (wanted: %d)", u.Have, u.Want)
}

// ===== [ Private Functions ] =====

// init - 패키지 초기화
func init() {}

// uniqueOutput - 지정된 파라미터 지정 정보를 기준으로 파라미터 구성 검증
func uniqueOutput(output []string) ([]string, int) {
	sort.Strings(output)
	j := 0
	outputSetSize := 0
	for i := 1; i < len(output); i++ {
		if output[j] == output[i] {
			continue
		}
		if !core.SequentialParamsPattern.MatchString(output[j]) {
			outputSetSize++
		}
		j++
		output[j] = output[i]
	}
	if j == len(output) {
		return output, outputSetSize
	}
	return output[:j+1], outputSetSize
}

// fromSetToSortedSlice - 지정된 맵 정보를 문자열 배열로 전환하고 정렬 처리
func fromSetToSortedSlice(set map[string]interface{}) []string {
	res := []string{}
	for element := range set {
		res = append(res, element)
	}
	sort.Strings(res)
	return res
}

// cleanHosts - Endpoint 및 Backend 설정에서 HostConfig 정보의 Host를 조정
func cleanHosts(hcs []*HostConfig) {
	for _, hc := range hcs {
		hc.Host = core.CleanHost(hc.Host)
	}
}

// ===== [ Public Functions ] =====

// NewDefinition - 시스템에서 사용할 Definition (Endpoint) 생성
func NewDefinition() *EndpointConfig {
	def := &EndpointConfig{}
	if err := def.InitializeDefaults(); nil != err {
		logging.GetLogger().WithError(err).Debug("[CONFIG] Failed to initialize definition (Endpoint)")
	}

	return def
}
