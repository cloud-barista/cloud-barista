// Package backend - 문자열 Hash 처리 패키지
package backend

// ===== [ Constants and Variables ] =====

const (
	offset64 uint64 = 14695981039346656037
	prime64         = 1099511628211
)

// ===== [ Types ] =====
// ===== [ Implementations ] =====
// ===== [ Private Functions ] =====
// ===== [ Public Functions ] =====

// PseudoFNV64a - 지정한 문자열을 Hash 처리
func PseudoFNV64a(s string) uint64 {
	h := offset64
	for i := 0; i < len(s); i++ {
		h ^= uint64(s[i])
		h *= prime64
	}
	return h
}
