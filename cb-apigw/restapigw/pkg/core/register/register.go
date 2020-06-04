package register

import "github.com/cloud-barista/cb-apigw/restapigw/pkg/core/register/internal"

// ===== [ Constants and Variables ] =====

// ===== [ Types ] =====

// Untyped - 형식없는 맵을 관리하는 인터페이스
type Untyped interface {
	Register(name string, v interface{})
	Get(name string) (interface{}, bool)
	Clone() map[string]interface{}
}

// ===== [ Implementations ] =====

// ===== [ Private Functions ] =====

// ===== [ Public Functions ] =====

// NewUntyped - 형식없는 맵 관리 인스턴스 생성
func NewUntyped() Untyped {
	return internal.NewUntyped()
}
