// Package proxy - Backend의 반환 결과를 Router에서 처리하기 위한 Response 객체로 변환처리하는 패키지
package proxy

import (
	"context"
	"net/http"

	"github.com/cloud-barista/cb-apigw/restapigw/pkg/encoding"
)

// ===== [ Constants and Variables ] =====

// ===== [ Types ] =====

// HTTPResponseParser - Backend에서 반환된 http.Response를 Proxy에서 관리하는 Response 로 변환하기 위한 함수 형식
type HTTPResponseParser func(context.Context, *http.Response) (*Response, error)

// HTTPResponseParserConfig - HTTPResponseParser 운영을 위한 설정 형식
type HTTPResponseParserConfig struct {
	Decoder         encoding.Decoder
	EntityFormatter EntityFormatter
}

// ===== [ Implementations ] =====

// ===== [ Private Functions ] =====

// ===== [ Public Functions ] =====

// NoOpHTTPResponseParser - http.Response 를 변경없이 전달하는 용도의 ResponseParser
func NoOpHTTPResponseParser(ctx context.Context, resp *http.Response) (*Response, error) {
	return &Response{
		Data:       map[string]interface{}{},
		IsComplete: true,
		Io:         NewReadCloserWrapper(ctx, resp.Body),
		Metadata: Metadata{
			StatusCode: resp.StatusCode,
			Headers:    resp.Header,
		},
	}, nil
}

// DefaultHTTPResponseParserFactory - NoOpResponseParser를 사용하지 않는 모든 경우에 사용할 ResponseParser
func DefaultHTTPResponseParserFactory(conf HTTPResponseParserConfig) HTTPResponseParser {
	return func(ctx context.Context, resp *http.Response) (*Response, error) {
		var data map[string]interface{}
		err := conf.Decoder(resp.Body, &data)
		resp.Body.Close()
		if nil != err {
			return nil, err
		}

		newResponse := Response{Data: data, IsComplete: true}
		newResponse = conf.EntityFormatter.Format(newResponse)
		return &newResponse, nil
	}
}
