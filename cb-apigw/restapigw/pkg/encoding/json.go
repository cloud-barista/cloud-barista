package encoding

import (
	"io"

	"github.com/cloud-barista/cb-apigw/restapigw/pkg/core"
)

// ===== [ Constants and Variables ] =====

const (
	// JSON - JSON 인코딩 식별자
	JSON = "json"
)

// ===== [ Types ] =====

// ===== [ Implementations ] =====

// ===== [ Private Functions ] =====

// ===== [ Public Functions ] =====

// JSONDecoder - 지정한 Reader의 JSON 데이터에 대한 Decoder를 생성하고 Decode 처리
func JSONDecoder(r io.Reader, v *map[string]interface{}) error {
	return core.JSONDecode(r, v)
}

// JSONCollectionDecoder - 지정한 Reader의 JSON 데이터에 대한 Collection 으로 Decoder를 생성하고 Decode 처리 (최종 반환할 때 Array인 형태로 변횐해서 처리)
func JSONCollectionDecoder(r io.Reader, v *map[string]interface{}) error {
	var collection []interface{}
	if err := core.JSONDecode(r, &collection); nil != err {
		return err
	}
	// Backend 결과 Array를 처리하기 위한 식별자 설정
	*(v) = map[string]interface{}{core.CollectionTag: collection, core.WrappingTag: core.CollectionTag}
	return nil
}

// JSONWrapedCollectionDecoder - 지정한 Reader의 JSON 데이터에 대한 Collection 으로 Decoder를 생성하고 Decode 처리
func JSONWrapedCollectionDecoder(r io.Reader, v *map[string]interface{}) error {
	var collection []interface{}
	if err := core.JSONDecode(r, &collection); nil != err {
		return err
	}
	// Backend 결과 Array를 처리하기 위한 식별자 설정
	*(v) = map[string]interface{}{core.CollectionTag: collection}
	return nil
}

// NewJSONDecoder - Collection 여부에 따라서 JSON Decoder 생성
func NewJSONDecoder(isCollection bool, wrapCollectionToJSON bool) func(io.Reader, *map[string]interface{}) error {
	if isCollection {
		if wrapCollectionToJSON {
			return JSONWrapedCollectionDecoder
		}
		return JSONCollectionDecoder
	}
	return JSONDecoder
}
