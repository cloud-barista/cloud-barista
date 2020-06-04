package config

import (
	"fmt"
	"regexp"
	"strings"
)

// ===== [ Constants and Variables ] =====

var (
	// Endpoint 파라미터 처리를 위한 패턴 "{...}"
	endpointURLKeysPattern = regexp.MustCompile(`/\{([a-zA-Z\-_0-9]+)\}`)
	// Host URL 패턴
	hostPattern = regexp.MustCompile(`(https?://)?([a-zA-Z0-9\._\-]+)(:[0-9]{2,6})?/?`)
	// RoutingPattern - 경로에 파라미터가 지정된 경우에 ":" 로 패턴 처리
	RoutingPattern = ColonRouterPatternBuilder
)

// ===== [ Types ] =====

// URIParser - URI 처리를 위한 인터페이스
type URIParser interface {
	CleanHosts([]string) []string
	CleanHost(string) string
	CleanPath(string) string
	GetEndpointPath(string, []string) string
}

// URI - URIParser 구현을 위한 데이터 관리 형식
type URI int

// ===== [ Implementations ] =====

// CleanHosts - 지정된 호스트들에 대해 호스트 패턴 처리
func (u URI) CleanHosts(hosts []string) []string {
	cleaned := []string{}
	for i := range hosts {
		cleaned = append(cleaned, u.CleanHost(hosts[i]))
	}
	return cleaned
}

// CleanHost - 지정된 호스트에 대해 패턴 검증 및 정리
func (URI) CleanHost(host string) string {
	matches := hostPattern.FindAllStringSubmatch(host, -1)
	if len(matches) != 1 {
		panic(fmt.Errorf("invalid host: %s", host))
	}
	keys := matches[0][1:]
	if keys[0] == "" {
		keys[0] = "http://"
	}
	return strings.Join(keys, "")
}

// CleanPath - 지정된 URI Path에 대해 정리
func (URI) CleanPath(path string) string {
	return "/" + strings.TrimPrefix(path, "/")
}

// GetEndpointPath - 지정된 URI Path에 대해 파라미터를 설정
func (u URI) GetEndpointPath(path string, params []string) string {
	result := path
	if u == ColonRouterPatternBuilder {
		for p := range params {
			parts := strings.Split(result, "?")
			// URI 에 지정된 `{xxx}` 파라미터를 URI에 맞는 문법으로 변경
			parts[0] = strings.Replace(parts[0], "{"+params[p]+"}", ":"+params[p], -1)
			result = strings.Join(parts, "?")
		}
	}
	return result
}

// ===== [ Private Functions ] =====

// ===== [ Public Functions ] =====

// NewURIParser - URIParser 인스턴스 생성
func NewURIParser() URIParser {
	return URI(RoutingPattern)
}
