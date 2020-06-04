package opencensus

import (
	"context"
	"net/http"

	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/trace"

	transport "github.com/cloud-barista/cb-apigw/restapigw/pkg/transport/http/client"
)

// ===== [ Constants and Variables ] =====

// ===== [ Types ] =====

// ===== [ Implementations ] =====

// ===== [ Private Functions ] =====

// ===== [ Public Functions ] =====

// HTTPRequestExecutor - Backend에 대한 Trace 가 활성화된 경우에 사용할 Opencensus의 Trace 정보를 사용하는 HTTP Request Executor 반환
func HTTPRequestExecutor(cf transport.HTTPClientFactory) transport.HTTPRequestExecutor {
	if !IsBackendEnabled() {
		return transport.DefaultHTTPRequestExecutor(cf)
	}

	return func(ctx context.Context, req *http.Request) (*http.Response, error) {
		// Factory를 통해서 http client 생성
		client := cf(ctx)
		if _, ok := client.Transport.(*ochttp.Transport); !ok {
			// Opencensus를 사용하는 Transport 사용
			client.Transport = &ochttp.Transport{Base: client.Transport}
		}

		logger.Debugf("[Backend Process Flow] Opencensus(HTTP Client) > %s", req.URL.String())

		return client.Do(req.WithContext(trace.NewContext(ctx, fromContext(ctx))))
	}
}
