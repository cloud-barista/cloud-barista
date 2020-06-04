package encoding

import (
	"io"
)

// ===== [ Constants and Variables ] =====

const (
	// NOOP - 응답을 변환없이 반환하기 위한 식별자
	NOOP = "no-op"
)

var (
	// 기본 Decoder들을 관리하는 Registry
	decoders = initDecoderRegister()
	// 기본 제공되는 Decoder들 정의
	defaultDecoders = map[string]func(bool, bool) func(io.Reader, *map[string]interface{}) error{
		JSON:   NewJSONDecoder,
		STRING: NewStringDecoder,
		NOOP:   noOpDecoderFactory,
	}
)

// ===== [ Types ] =====

// Decoder - 지정된 Reader의 데이터를 읽어서 맵으로 Decode 처리 함수 형식
type Decoder func(io.Reader, *map[string]interface{}) error

// DecoderFactory - CollectionDecoder나 EntityDecoder를 반환하는 함수 형식
type DecoderFactory func(bool, bool) func(io.Reader, *map[string]interface{}) error

// ===== [ Implementations ] =====

// ===== [ Private Functions ] =====

// noOpDecoderFactory - NO-OP Decoder를 반환하는 팩토리
func noOpDecoderFactory(_ bool, _ bool) func(io.Reader, *map[string]interface{}) error {
	return NoOpDecoder
}

// ===== [ Public Functions ] =====

// Register - 지정된 이름으로 Decoder 등록
func Register(name string, dec func(bool, bool) func(io.Reader, *map[string]interface{}) error) error {
	return decoders.Register(name, dec)
}

// NoOpDecoder - NO-OP Decoder를 처리를 위한 Dummy 반환
func NoOpDecoder(_ io.Reader, _ *map[string]interface{}) error { return nil }

// Get - 지정한 이름의 Decoder 반환
func Get(name string) DecoderFactory {
	return decoders.Get(name)
}
