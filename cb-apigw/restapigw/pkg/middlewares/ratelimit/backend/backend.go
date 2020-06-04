// Package backend - Rate Limit 기능을 제공하는 저장소 패키지
package backend

import "time"

// ===== [ Constants and Variables ] =====

var (
	// DataTTL - 데이터 유지 시간
	DataTTL = 10 * time.Minute

	now           = time.Now
	shards uint64 = 2048
)

// ===== [ Types ] =====

// Hasher - 지정된 문자열을 Hash 처리하는 함수 형식
type Hasher func(string) uint64

// IBackend - 데이터 저장소 인터페이스
type IBackend interface {
	Load(string, func() interface{}) interface{}
	Store(string, interface{}) error
}

// ===== [ Implementations ] =====
// ===== [ Private Functions ] =====
// ===== [ Public Functions ] =====
