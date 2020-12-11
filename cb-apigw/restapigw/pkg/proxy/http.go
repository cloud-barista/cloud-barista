package proxy

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/cloud-barista/cb-apigw/restapigw/pkg/config"
	"github.com/cloud-barista/cb-apigw/restapigw/pkg/core"
	"github.com/cloud-barista/cb-apigw/restapigw/pkg/encoding"
	"github.com/cloud-barista/cb-apigw/restapigw/pkg/logging"
	"github.com/cloud-barista/cb-apigw/restapigw/pkg/transport/http/client"
)

// ===== [ Constants and Variables ] =====

var (
	logger = logging.NewLogger()
)

// ===== [ Types ] =====

// responseError - Defines interface for response error
type responseError interface {
	Error() string
	Name() string
	StatusCode() int
}

// wrappedError - Defines interface for Wrapped response error
type wrappedError interface {
	Error() error
	StatusCode() int
	Message() string
}

// ===== [ Implementations ] =====

// ===== [ Private Functions ] =====

// ===== [ Public Functions ] =====

// NewHTTPProxyWithHTTPExecutor - 지정된 BackendConfig 와 HTTP Request Executor와 응답 처리에 사용할 Decoder를 설정한 Proxy 반환
func NewHTTPProxyWithHTTPExecutor(bconf *config.BackendConfig, hre client.HTTPRequestExecutor, dec encoding.Decoder) Proxy {
	if bconf.Encoding == encoding.NOOP {
		return NewHTTPProxyDetailed(bconf, hre, client.NoOpHTTPStatusHandler, NoOpHTTPResponseParser)
	}

	ef := NewEntityFormatter(bconf)
	rp := DefaultHTTPResponseParserFactory(HTTPResponseParserConfig{dec, ef})
	return NewHTTPProxyDetailed(bconf, hre, client.GetHTTPStatusHandler(bconf), rp)
}

// NewHTTPProxyDetailed - 지정된 BackendConfig와 HTTP Reqeust Executor와 응답 처리에 사용할 StatusHandler, Response Parser를 설정한 Proxy 반환
func NewHTTPProxyDetailed(bconf *config.BackendConfig, hre client.HTTPRequestExecutor, hsh client.HTTPStatusHandler, hrp HTTPResponseParser) Proxy {
	return func(ctx context.Context, req *Request) (*Response, error) {
		reqToBackend, err := http.NewRequest(strings.ToTitle(req.Method), req.URL.String(), req.Body)
		if nil != err {
			return nil, err
		}

		// Backend 호출에 필요한 Header 정보 설정
		reqToBackend.Header = make(map[string][]string, len(req.Headers))
		for k, v := range req.Headers {
			tmp := make([]string, len(v))
			copy(tmp, v)
			reqToBackend.Header[k] = tmp
		}

		// Body Size 정보 설정
		if nil != req.Body {
			if v, ok := req.Headers["Content-Length"]; ok && 1 == len(v) && "chunked" != v[0] {
				if size, err := strconv.Atoi(v[0]); nil == err {
					reqToBackend.ContentLength = int64(size)
				}
			}
		}

		// Backend 호출에 필요한 Query String 정보 설정
		if nil != req.Query {
			q := reqToBackend.URL.Query()
			for k, v := range req.Query {
				q.Add(k, v[0])
			}
			reqToBackend.URL.RawQuery = q.Encode()
		}

		logger.Debugf("[Backend Process Flow] Proxy(HTTP) > CallChain (%s)", reqToBackend.URL.String())

		// Backed 호출
		resp, err := hre(ctx, reqToBackend)
		if nil != reqToBackend.Body {
			reqToBackend.Body.Close()
		}

		select {
		case <-ctx.Done():
			// 호출에 문제가 생긴 경우
			return nil, ctx.Err()
		default:
		}

		// 응답이 없어서 상태 확인이 안되는 경우
		if nil == resp && nil != err {
			return nil, err
		}

		// Response Status 처리
		resp, err = hsh(ctx, resp)
		if nil != err {
			if t, ok := err.(responseError); ok {
				return &Response{
					Data: map[string]interface{}{
						fmt.Sprintf("error_%s", t.Name()): t,
					},
					Metadata: Metadata{StatusCode: t.StatusCode()},
				}, err
			} else if we, ok := err.(core.WrappedError); ok {
				return nil, we
			}
			return nil, err
		}

		// Response Parser 호출
		return hrp(ctx, resp)
	}
}

// NewRequestBuilderChain - Request 파라미터와 Backend Path를 설정한 Proxy 호출 체인을 생성한다.
func NewRequestBuilderChain(bConf *config.BackendConfig) CallChain {
	return func(next ...Proxy) Proxy {
		if 1 < len(next) {
			panic(ErrTooManyProxies)
		}
		return func(ctx context.Context, req *Request) (*Response, error) {
			r := req.Clone()
			// Bypass가 아닌 경우는 Path와 Method를 설정에 맞도록 재 구성
			if !req.IsBypass {
				r.GeneratePath(bConf.URLPattern)
				r.Method = bConf.Method
			}

			return next[0](ctx, &r)
		}
	}
}
