package proxy

import "github.com/cloud-barista/cb-apigw/restapigw/pkg/core/register"

// ===== [ Constants and Variables ] =====

// ===== [ Types ] =====

// combinerRegister - Response Merging 처리를 위한 Register 구조
type combinerRegister struct {
	data     register.Untyped
	fallback ResponseCombiner
}

// ===== [ Implementations ] =====

// GetRseponseCombiner - 지정한 이름에 해당하는 Response Combiner 반환
func (cr *combinerRegister) GetResponseCombiner(name string) (ResponseCombiner, bool) {
	v, ok := cr.data.Get(name)
	if !ok {
		return cr.fallback, ok
	}
	if rc, ok := v.(ResponseCombiner); ok {
		return rc, ok
	}
	return cr.fallback, ok
}

// ===== [ Private Functions ] =====

// newCombinerRegister - 지정된 데이터로 구성되는 Response Merging 처리용 Combiner 생성
func newCombinerRegister(data map[string]ResponseCombiner, fallback ResponseCombiner) *combinerRegister {
	r := register.NewUntyped()
	for k, v := range data {
		r.Register(k, v)
	}
	return &combinerRegister{r, fallback}
}

// ===== [ Public Functions ] =====
