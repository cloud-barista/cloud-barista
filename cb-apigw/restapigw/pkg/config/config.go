// Package config - Configuration for Cloud-Barista's REST API Gateway and provides the required process
package config

import (
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/cloud-barista/cb-apigw/restapigw/pkg/encoding"
)

// ===== [ Constants and Variables ] =====

const (
	// BracketsRouterPatternBuilder - 라우터 파라미터 구분으로 "{" 을 사용하는 경우 식별자
	BracketsRouterPatternBuilder = iota
	// ColonRouterPatternBuilder - 라우터 파라미터 구분으로 ":" 을 사용하는 경우 식별자
	ColonRouterPatternBuilder
	// DefaultMaxIdleConnectionsPerHost - 기본 값으로 사용할 호스트당 Max Idle connection 허용 수
	DefaultMaxIdleConnectionsPerHost = 250
	// DefaultTimeout - 기본 값으로 사용할 처리 허용 시간
	DefaultTimeout = 2 * time.Second
	// ConfigVersion - 설정 구조에 대한 버전
	ConfigVersion = 1
)

var (
	// 일반적 URL 패턴
	simpleURLKeysPattern = regexp.MustCompile(`\{([a-zA-Z\-_0-9]+)\}`)
	// Sequential 처리용 URL 패턴 (전 수행의 결과를 파라미터로 지정하는 경우)
	sequentialParamsPattern = regexp.MustCompile(`^resp[\d]+_.*$`)
	// Debug URL 패턴
	debugPattern = "^[^/]|/__debug(/.*)?$"

	errInvalidHost         = errors.New("invalid host")
	errInvalidNoOpEncoding = errors.New("can not use NoOp encoding with more than one backends connected to the same endpoint")
	defaultPort            = 8000
)

// ===== [ Types ] =====

// ServiceConfig - REST API Gateway 운영에 필요한 서비스 설정 구조
type ServiceConfig struct {
	// 서비스 식별 명 (기본값: '')
	Name string `mapstructure:"name"`
	// 기본 처리 시간 (기본값: 2s)
	Timeout time.Duration `mapstructure:"timeout"`
	// 디버그모드 여부 (기본값: false)
	Debug bool `mapstructure:"debug"`
	// GET 처리에 대한 캐시 TTL 기간 (기본값: 1h)
	CacheTTL time.Duration `mapstructure:"cache_ttl"`
	// 전역으로 사용할 기본 Host 리스트
	Host []string `mapstructure:"host"`
	// 서비스에서 사용할 포트 (기본값: 8000)
	Port int `mapstructure:"port"`
	// 설정 파일 버전
	Version int `mapstructure:"version"`
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
	MaxIdleConnectionsPerHost int `mapstructure:"max_idle_connections_per_host"`
	// 유휴연결(Keep-alive)의 최대 유효시간 (기본값: , 0이면 ReadTimeout 사용, 이것도 0이면 ReadHeaderTimeout 사용)
	IdleConnectionTimeout time.Duration `mapstructure:"idle_connection_timeout"`
	// 라우팅에 사용할 Endpoint 설정 리스트
	Endpoints []*EndpointConfig `mapstructure:"endpoints"`
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
	// OutputEncoding - Endpoint의 Response 처리에 사용할 기본 Encoding (기본값: "JSON")
	OutputEncoding string `mapstructure:"output_encoding"`
	// DisableStrictREST - REST 강제 규칙 비활성화 여부 (기본값: false)
	DisableStrictREST bool `mapstructure:"disable_strict_rest"`
	// URI Parser (internal use)
	uriParser URIParser
}

// EndpointConfig - 서비스 라우팅에 사용할 설정 구조
type EndpointConfig struct {
	// Bypass 처리 여부
	IsBypass bool

	// 클라이언트에 노출될 URL 패턴
	Endpoint string `mapstructure:"endpoint"`
	// Endpoint에 대한 HTTP 메서드 (GET, POST, PUT, etc) (기본값: GET)
	Method string `mapstructure:"method"`
	// Endpoint 처리 시간 (기본값: 서비스 값 사용)
	Timeout time.Duration `mapstructure:"timeout"`
	// GET 처리에 대한 캐시 TTL 기간 (기본값: 서비스 값 사용)
	CacheTTL time.Duration `mapstructure:"cache_ttl"`
	// 반환결과 처리에 사용할 인코딩
	OutputEncoding string `mapstructure:"output_encoding"`
	// Backend 에 전달되는 Query String에서 제외할 파라미터 Key 리스트
	ExceptQueryStrings []string `mapstructure:"except_querystrings"`
	// Backend 에 전달되는 Header에서 제외할 파라미터 Key 리스트
	ExceptHeaders []string `mapstructure:"except_headers"`
	// Endpoint 단위에서 적용할 Middleware 설정
	Middleware MWConfig `mapstructure:"middleware"`
	// Endpoint에서 호출할 Backend API 서버 호출/응답 처리 설정 리스트
	Backend []*BackendConfig `mapstructure:"backend"`
}

// BackendConfig - Backend API Server 연결과 응답 처리를 위한 설정 구조
type BackendConfig struct {
	// Backend API Server의 Host URI (서비스에서 전역으로 설정한 경우는 생략 가능)
	Host []string `mapstructure:"host"`
	// Backend 처리 시간 (기본값: )
	Timeout time.Duration
	// Backend 호출에 사용할 HTTP Method
	Method string `mapstructure:"method"`
	// Backend 호출에 사용할 URL Patthern
	URLPattern string `mapstructure:"url_pattern"`
	// 인코딩 포맷
	Encoding string `mapstructure:"encoding"`
	// Backend 결과를 묶을 Group 명 (기본값: )
	Group string `mapstructure:"group"`
	// Backend 결과에서 생략할 필드명 리스트 (flatmap 적용 "." operation)
	Blacklist []string `mapstructure:"blacklist"`
	// Backend 결과에서 추출할 필드명 리스트 (flatmap 적용 "." operation)
	Whitelist []string `mapstructure:"whitelist"`
	// Backend 결과에서 필드명을 변경할 리스트 맵
	Mapping map[string]string `mapstructure:"mapping"`
	// Backend 결과가 컬랙션인지 여부
	IsCollection bool `mapstructure:"is_collection"`
	// Backend 결과가 컬랙션인 경우에 core.CollectionTag ("collection") 으로 JSON 포맷을 할 것인지 여부 (True 면 core.CollectionTag ("collection") 으로 JSON 전환, false면 Array 상태로 반환)
	WrapCollectionToJSON bool `mapstructure:"wrap_collection_to_json"`
	// Backend 결과 중에서 특정한 필드만 처리할 경우의 필드명
	Target string `mapstructure:"target"`
	// Backend 에서 동작할 Middleware 설정
	Middleware MWConfig `mapstructure:"middleware"`
	// HostSanitizationDisabled - host 정보의 정제작업 비활성화 여부
	HostSanitizationDisabled bool `mapstructure:"disable_host_sanitize"`
	// API 호출의 응답을 파싱하기 위한 디코더 (내부 사용)
	Decoder encoding.Decoder `json:"-"`
	// URLPattern에서 파라미터 변환에 사용할 키 관리 (내부 사용)
	URLKeys []string
}

// TLSConfig - 서비스에서 사용할 TLS 설정 구조
type TLSConfig struct {
	// TLS 비활성화 여부
	IsDisabled bool `mapstructure:"disabled"`
	// 공개 키 경로
	PublicKey string `mapstructure:"public_key"`
	// 비밀 키 경로
	PrivateKey string `mapstructure:"private_key"`
	// TLS 최소 버전
	MinVersion string `mapstructure:"min_version"`
	// TLS 최대 버전
	MaxVersion string `mapstructure:"max_version"`
	// Curve 설정들의 리스트 (use 23 for CurveP256, 24 for CurveP384 or 25 for CurveP521)
	CurvePreferences []uint16 `mapstructure:"curve_preferences"`
	// 서버에서 사용을 강제하는 Cipher Suite 리스트
	PreferServerCipherSuites bool `mapstructure:"prefer_server_cipher_suites"`
	// Chiper Suite 리스트
	CipherSuites []uint16 `mapstructure:"cipher_suites"`
}

// MWConfig - Middleware 설정을 저장하기 위한 맵 구조 (개별 Middlewares에서 설정 Parsing 적용)
type MWConfig map[string]interface{}

// UnsupportedVersionError - 설정 초기화 과정에서 버전 검증을 통해 반환할 오류 구조
type UnsupportedVersionError struct {
	Have int
	Want int
}

// WrongNumberOfParamsError - 파라미터의 IN/OUT 갯수가 다를 경우에 반환할 오류 구조
type WrongNumberOfParamsError struct {
	Endpoint     string
	Method       string
	Backend      int
	InputParams  []string
	OutputParams []string
}

// UndefinedOutputParamError - IN 파라미터에 대한 OUT 파라미터가 지정되지 않았을 경우에 반환할 오류 구조
type UndefinedOutputParamError struct {
	Endpoint     string
	Method       string
	Backend      int
	InputParams  []string
	OutputParams []string
	Param        string
}

// EndpointMatchError - Endpoint 패턴에 문제가 있을 경우에 반환할 오류 구조
type EndpointMatchError struct {
	Path   string
	Method string
	Err    error
}

// NoBackendsError - Backend가 지정되지 않았을 경우에 반환할 오류 구조
type NoBackendsError struct {
	Path   string
	Method string
}

// EndpointPathError - Endpoint 경로에 문제가 있을 경우에 반환할 오류 구조
type EndpointPathError struct {
	Path   string
	Method string
}

// ===== [ Implementations ] =====

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
	sConf.uriParser = NewURIParser()

	// 설정 파일 버전 검증
	if sConf.Version != ConfigVersion {
		return &UnsupportedVersionError{
			Have: sConf.Version,
			Want: ConfigVersion,
		}
	}

	// 전역변수 초기화
	sConf.initGlobalParams()

	// Endpoint 초기화
	return sConf.initEndPoints()
}

// initGlobalParams - 전역 설정 초기화
func (sConf *ServiceConfig) initGlobalParams() {
	if sConf.Port == 0 {
		sConf.Port = defaultPort
	}
	if sConf.MaxIdleConnectionsPerHost == 0 {
		sConf.MaxIdleConnectionsPerHost = DefaultMaxIdleConnectionsPerHost
	}
	if sConf.Timeout == 0 {
		sConf.Timeout = DefaultTimeout
	}

	// Host 조정
	sConf.Host = sConf.uriParser.CleanHosts(sConf.Host)

	// 전역 설정에 대한 Middleware 구성 맵 관리
	sConf.Middleware.sanitize()
}

// extractPlaceHoldersFromURLTemplate - URL에 설정되어 있는 Placeholder를 추출한 문자열들 반환
func (sConf *ServiceConfig) extractPlaceHoldersFromURLTemplate(subject string, pattern *regexp.Regexp) []string {
	matches := pattern.FindAllStringSubmatch(subject, -1)
	keys := make([]string, len(matches))
	for k, v := range matches {
		keys[k] = v[1]
	}
	return keys
}

// paramExtractionPattern - 파라미터 추출에 사용할 패턴 반환
func (sConf *ServiceConfig) paramExtractionPattern() *regexp.Regexp {
	if sConf.DisableStrictREST {
		return simpleURLKeysPattern
	}
	return endpointURLKeysPattern
}

// initEndpointDefaults - Endpoint에 미 설정된 항목들을 기본 값으로 초기화
func (sConf *ServiceConfig) initEndpointDefaults(epIdx int) {
	endpoint := sConf.Endpoints[epIdx]
	if endpoint.Method == "" {
		endpoint.Method = "GET"
	}
	if sConf.CacheTTL != 0 && endpoint.CacheTTL == 0 {
		endpoint.CacheTTL = sConf.CacheTTL
	}
	if sConf.Timeout != 0 && endpoint.Timeout == 0 {
		endpoint.Timeout = sConf.Timeout
	}
	// 기본 출력 포맷 설정 (JSON)
	if endpoint.OutputEncoding == "" {
		if sConf.OutputEncoding != "" {
			endpoint.OutputEncoding = sConf.OutputEncoding
		} else {
			endpoint.OutputEncoding = encoding.JSON
		}
	}
}

// initBackendDefaults - Backend에 미 설정된 항목들을 기본 값으로 초기화
func (sConf *ServiceConfig) initBackendDefaults(epIdx, bIdx int) {
	endpoint := sConf.Endpoints[epIdx]
	backend := endpoint.Backend[bIdx]
	if len(backend.Host) == 0 {
		// URL 미 지정시 전역 URL 사용
		backend.Host = sConf.Host
	} else if !backend.HostSanitizationDisabled {
		backend.Host = sConf.uriParser.CleanHosts(backend.Host)
	}
	// Method 미 지정시 Endpoint Method 사용
	if backend.Method == "" {
		backend.Method = endpoint.Method
	}
	backend.Timeout = endpoint.Timeout
	backend.Decoder = encoding.Get(strings.ToLower(backend.Encoding))(backend.IsCollection, backend.WrapCollectionToJSON)
}

// initBackendURLMappings - Backend에 지정된 파라미터 정보들을 이후에 사용할 수 있도록 초기화
func (sConf *ServiceConfig) initBackendURLMappings(epIdx, bIdx int, inputParams map[string]interface{}) error {
	backend := sConf.Endpoints[epIdx].Backend[bIdx]

	backend.URLPattern = sConf.uriParser.CleanPath(backend.URLPattern)
	// 파라미터 설정 추출 및 관리
	outputParams, outputSetSize := uniqueOutput(sConf.extractPlaceHoldersFromURLTemplate(backend.URLPattern, simpleURLKeysPattern))

	ip := fromSetToSortedSlice(inputParams)

	if outputSetSize > len(ip) {
		return &WrongNumberOfParamsError{
			Endpoint:     sConf.Endpoints[epIdx].Endpoint,
			Method:       sConf.Endpoints[epIdx].Method,
			Backend:      bIdx,
			InputParams:  ip,
			OutputParams: outputParams,
		}
	}

	backend.URLKeys = []string{}
	for _, output := range outputParams {
		if !sequentialParamsPattern.MatchString(output) {
			if _, ok := inputParams[output]; !ok {
				return &UndefinedOutputParamError{
					Param:        output,
					Endpoint:     sConf.Endpoints[epIdx].Endpoint,
					Method:       sConf.Endpoints[epIdx].Method,
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

// initEndpoints - Endpoint 초기화 처리
func (sConf *ServiceConfig) initEndPoints() error {
	for epIdx, eConf := range sConf.Endpoints {
		eConf.Endpoint = sConf.uriParser.CleanPath(eConf.Endpoint)

		if err := eConf.validate(); err != nil {
			return err
		}

		// Endpoint URL에 지정된 파라미터를 Input parameter로 설정
		inputParams := sConf.extractPlaceHoldersFromURLTemplate(eConf.Endpoint, sConf.paramExtractionPattern())
		inputSet := map[string]interface{}{}
		for ip := range inputParams {
			inputSet[inputParams[ip]] = nil
		}

		eConf.Endpoint = sConf.uriParser.GetEndpointPath(eConf.Endpoint, inputParams)

		sConf.initEndpointDefaults(epIdx)

		if eConf.OutputEncoding == encoding.NOOP && len(eConf.Backend) > 1 {
			return errInvalidNoOpEncoding
		}

		// Endpoint 설정의 Middleware 설정 맵 관리
		eConf.Middleware.sanitize()

		for bIdx, bConf := range eConf.Backend {
			sConf.initBackendDefaults(epIdx, bIdx)

			if err := sConf.initBackendURLMappings(epIdx, bIdx, inputSet); err != nil {
				return err
			}

			// Backend 설정의 Middleware 설정 맵 관리
			bConf.Middleware.sanitize()
		}
	}
	return nil
}

// validate - Endpoint 별 세부 필수 항목 검증
func (e *EndpointConfig) validate() error {
	matched, err := regexp.MatchString(debugPattern, e.Endpoint)
	if err != nil {
		return &EndpointMatchError{
			Err:    err,
			Path:   e.Endpoint,
			Method: e.Method,
		}
	}
	if matched {
		return &EndpointPathError{Path: e.Endpoint, Method: e.Method}
	}

	if len(e.Backend) == 0 {
		return &NoBackendsError{Path: e.Endpoint, Method: e.Method}
	}
	return nil
}

// Error - 비 호환 버전에 대한 오류 문자열 반환
func (u *UnsupportedVersionError) Error() string {
	return fmt.Sprintf("Unsupported version: %d (wanted: %d)", u.Have, u.Want)
}

// ===== [ Private Functions ] =====

// uniqueOutput - 지정된 파라미터 지정 정보를 기준으로 파라미터 구성 검증
func uniqueOutput(output []string) ([]string, int) {
	sort.Strings(output)
	j := 0
	outputSetSize := 0
	for i := 1; i < len(output); i++ {
		if output[j] == output[i] {
			continue
		}
		if !sequentialParamsPattern.MatchString(output[j]) {
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

// ===== [ Public Functions ] =====
