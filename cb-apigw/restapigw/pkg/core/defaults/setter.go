// Package defaults -
package defaults

// ===== [ Constants and Variables ] =====
// ===== [ Types ] =====

type (
	// Setter - Struct에 기본 값 설정이 가능하도록 메서드를 제공하는 인터페이스
	Setter interface {
		SetDefaults()
	}
)

// ===== [ Implementations ] =====
// ===== [ Private Functions ] =====

// callSetter - 지정한 값의 Setter 인터페이스 구현인지를 검증하고 메서드 호출
func callSetter(v interface{}) {
	if ds, ok := v.(Setter); ok {
		ds.SetDefaults()
	}
}

// ===== [ Public Functions ] ===== 
