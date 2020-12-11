package core

import (
	"fmt"
	"regexp"
	"strings"
)

// ===== [ Constants and Variables ] =====

var (
	// SimpleURLKeysPattern - 일반적 URL 패턴
	SimpleURLKeysPattern = regexp.MustCompile(`\{([a-zA-Z\-_0-9]+)\}`)
	// SequentialParamsPattern - Sequential 처리용 URL 패턴 (전 수행의 결과를 파라미터로 지정하는 경우)
	SequentialParamsPattern = regexp.MustCompile(`^resp[\d]+_.*$`)
	// DebugPattern - Debug URL 패턴
	DebugPattern = "^[^/]|/__debug(/.*)?$"
	// URLKeyPattern - URL 에 대한 파라미터 처리를 위한 패턴 "{...}"
	URLKeyPattern = regexp.MustCompile(`/\{([a-zA-Z\-_0-9]+)\}`)

	// HOSTPattern - Host URL 패턴
	HOSTPattern = regexp.MustCompile(`(https?://)?([a-zA-Z0-9\._\-]+)(:[0-9]{2,6})?/?`)
)

// ===== [ Types ] =====

// ===== [ Implementations ] =====

// ===== [ Private Functions ] =====

// ===== [ Public Functions ] =====

// CleanHosts - 지정된 호스트들에 대해 호스트 패턴 처리
func CleanHosts(hosts []string) []string {
	cleaned := []string{}
	for i := range hosts {
		cleaned = append(cleaned, CleanHost(hosts[i]))
	}
	return cleaned
}

// CleanHost - 지정된 호스트에 대해 패턴 검증 및 정리
func CleanHost(host string) string {
	matches := HOSTPattern.FindAllStringSubmatch(host, -1)
	if 1 != len(matches) {
		panic(fmt.Errorf("invalid host: %s", host))
	}
	keys := matches[0][1:]
	if "" == keys[0] {
		keys[0] = "http://"
	}
	return strings.Join(keys, "")
}

// CleanPath - 지정된 URI Path에 대해 정리
func CleanPath(path string) string {
	return "/" + strings.TrimPrefix(path, "/")
}

// GetParameteredPath - 지정된 URI Path에 대해 파라미터를 설정
func GetParameteredPath(path string, params []string) string {
	result := path

	for p := range params {
		parts := strings.Split(result, "?")
		// URI 에 지정된 `{xxx}` 파라미터를 URI에 맞는 문법으로 변경
		parts[0] = strings.Replace(parts[0], "{"+params[p]+"}", ":"+params[p], -1)
		result = strings.Join(parts, "?")
	}

	return result
}

// ParameterExtractPattern - 파라미터 추출에 사용할 패턴 반환
func ParameterExtractPattern(disableStrictREST bool) *regexp.Regexp {
	if disableStrictREST {
		return SimpleURLKeysPattern
	}
	return URLKeyPattern
}

// ExtractPlaceHoldersFromURLTemplate - URL에 설정되어 있는 Placeholder를 추출한 문자열들 반환
func ExtractPlaceHoldersFromURLTemplate(subject string, pattern *regexp.Regexp) []string {
	matches := pattern.FindAllStringSubmatch(subject, -1)
	keys := make([]string, len(matches))
	for k, v := range matches {
		keys[k] = v[1]
	}
	return keys
}
