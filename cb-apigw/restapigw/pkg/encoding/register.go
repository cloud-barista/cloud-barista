package encoding

import (
	"io"

	"github.com/cloud-barista/cb-apigw/restapigw/pkg/core/register"
)

// ===== [ Constants and Variables ] =====

// ===== [ Types ] =====

// DecoderRegister - 형식없는 맵 방식으로 Decoder를 관리하기 위한 등록 구조
type DecoderRegister struct {
	data register.Untyped
}

// ===== [ Implementations ] =====

// Register - 지정한 이름으로 Decorder 등록 구현
func (r *DecoderRegister) Register(name string, dec func(bool, bool) func(io.Reader, *map[string]interface{}) error) error {
	r.data.Register(name, dec)
	return nil
}

// Get - 지정한 이름에 해당하는 Decoder 반환
func (r *DecoderRegister) Get(name string) func(bool, bool) func(io.Reader, *map[string]interface{}) error {
	for _, n := range []string{name, JSON} {
		if v, ok := r.data.Get(n); ok {
			if dec, ok := v.(func(bool, bool) func(io.Reader, *map[string]interface{}) error); ok {
				return dec
			}
		}
	}
	return NewJSONDecoder
}

// ===== [ Private Functions ] =====

// initDecoderRegister - 기본 Decoder들을 관리하는 Decoder 관리기 생성 및 초기화
func initDecoderRegister() *DecoderRegister {
	r := &DecoderRegister{register.NewUntyped()}
	for k, v := range defaultDecoders {
		r.Register(k, v)
	}
	return r
}

// ===== [ Public Functions ] =====
