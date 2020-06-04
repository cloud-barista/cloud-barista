package opencensus

import (
	"context"
	"strings"

	"github.com/cloud-barista/cb-apigw/restapigw/pkg/config"
	"github.com/cloud-barista/cb-apigw/restapigw/pkg/core"
	"github.com/cloud-barista/cb-apigw/restapigw/pkg/logging"
	"github.com/cloud-barista/cb-apigw/restapigw/pkg/proxy"
	"go.opencensus.io/trace"
)

// ===== [ Constants and Variables ] =====

const (
	// errCtxCanceled - Context가 취소된 경우 오류
	errCtxCanceledMsg = "context canceled"
)

var (
	logger = logging.NewLogger()
)

// ===== [ Types ] =====

// ===== [ Implementations ] =====

// ===== [ Private Functions ] =====

// ===== [ Public Functions ] =====

// CallChain - 지정한 이름을 기준으로 순차적으로 실행되는 Call chain 함수 반환
func CallChain(tag string, isBypass bool, url string) proxy.CallChain {
	if !IsProxyEnabled() {
		return proxy.EmptyChain
	}
	return func(next ...proxy.Proxy) proxy.Proxy {
		if len(next) > 1 {
			panic(proxy.ErrTooManyProxies)
		}
		if len(next) < 1 {
			panic(proxy.ErrNotEnoughProxies)
		}

		return func(ctx context.Context, req *proxy.Request) (*proxy.Response, error) {
			var span *trace.Span
			var name = tag + " "
			if isBypass {
				name += req.Path
			} else {
				name += url
			}

			logger.Debugf("[Backend Process Flow] Opencensus > CallChain (%s) > %s", req.Path, tag)

			ctx, span = trace.StartSpan(trace.NewContext(ctx, fromContext(ctx)), name)
			resp, err := next[0](ctx, req)
			// 응답과 오류 여부에 따른 Trace 정보 설정
			if err != nil {
				if err.Error() != errCtxCanceledMsg {
					if resp != nil {
						span.SetStatus(trace.Status{Code: int32(resp.Metadata.StatusCode), Message: err.Error()})
					} else {
						if we, ok := err.(core.WrappedError); ok {
							if tag != "[backend]" {
								//span.SetStatus(trace.Status{Code: 500, Message: err.Error()})
								span.SetStatus(trace.Status{Code: int32(we.Code()), Message: we.GetError().Error()})
							} else {
								// Backend 호출의 경우는 원본 오류 메시지를 그대로 출력
								span.SetStatus(trace.Status{Code: int32(we.Code()), Message: we.Error()})
							}
						} else {
							span.SetStatus(trace.Status{Code: 500, Message: err.Error()})
						}
					}
				} else {
					span.AddAttributes(trace.BoolAttribute("error", true))
					span.AddAttributes(trace.BoolAttribute("canceled", true))
				}
			}
			span.AddAttributes(trace.BoolAttribute(("complete"), resp != nil && resp.IsComplete))
			span.End()

			return resp, err
		}
	}
}

// ProxyFactory - Opencensus Trace와 연동되는 Proxy Call chain 구성을 위한 팩토리
func ProxyFactory(pf proxy.Factory) proxy.FactoryFunc {
	if !IsProxyEnabled() {
		return pf.New
	}
	return func(eConf *config.EndpointConfig) (proxy.Proxy, error) {
		next, err := pf.New(eConf)
		if err != nil {
			return next, err
		}
		// Endpoint Bypass 검증 및 Flag 설정
		eConf.IsBypass = strings.HasSuffix(eConf.Endpoint, core.Bypass)
		return CallChain("[proxy]", eConf.IsBypass, eConf.Endpoint)(next), nil
	}
}

// BackendFactory - Opencensus Trace와 연동되는 Backend Call chain 구성을 위한 팩토리
func BackendFactory(bf proxy.BackendFactory) proxy.BackendFactory {
	if !IsBackendEnabled() {
		return bf
	}
	return func(bConf *config.BackendConfig) proxy.Proxy {
		// Backend Bypass 검증
		return CallChain("[backend]", strings.HasSuffix(bConf.URLPattern, core.Bypass), bConf.URLPattern)(bf(bConf))
	}
}
