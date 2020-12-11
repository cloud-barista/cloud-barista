// Package server -
package server

import (
	"context"

	"github.com/cloud-barista/cb-apigw/restapigw/pkg/config"
	"github.com/cloud-barista/cb-apigw/restapigw/pkg/logging"
	"github.com/cloud-barista/cb-apigw/restapigw/pkg/middlewares/httpcache"
	ginMetrics "github.com/cloud-barista/cb-apigw/restapigw/pkg/middlewares/metrics/gin"
	"github.com/cloud-barista/cb-apigw/restapigw/pkg/middlewares/opencensus"
	ratelimitProxy "github.com/cloud-barista/cb-apigw/restapigw/pkg/middlewares/ratelimit/proxy"
	"github.com/cloud-barista/cb-apigw/restapigw/pkg/proxy"
	"github.com/cloud-barista/cb-apigw/restapigw/pkg/transport/http/client"
)

// ===== [ Constants and Variables ] =====
// ===== [ Types ] =====
// ===== [ Implementations ] =====
// ===== [ Private Functions ] =====
// ===== [ Public Functions ] =====

// setupGinBackendFactoryWithContext - 지정한 Context 기반으로 Middleware 들이 적용된 BackendFactory 설정
func setupGinBackendFactoryWithContext(ctx context.Context, logger logging.Logger, mc *ginMetrics.Collector) proxy.BackendFactory {
	requestExecutorFactory := func(bConf *config.BackendConfig) client.HTTPRequestExecutor {
		var clientFactory client.HTTPClientFactory
		// TODO: Backend Auth

		// HTTPCache 가 적용된 HTTP Client
		clientFactory = httpcache.NewHTTPClient(bConf)
		// Opencensus 와 연계된 HTTP Request Executor
		return opencensus.HTTPRequestExecutor(clientFactory)
	}

	// Opencensus HTTPRequestExecutor를 사용하는 Base BackendFactory 설정
	backendFactory := func(bConf *config.BackendConfig) proxy.Proxy {
		return proxy.NewHTTPProxyWithHTTPExecutor(bConf, requestExecutorFactory(bConf), bConf.Decoder)
	}

	// TODO: Martian for Backend

	// Backend 호출에 대한 Rate Limit Middleware 설정
	backendFactory = ratelimitProxy.BackendFactory(backendFactory)

	// TODO: Circuit-Breaker for Backend

	// Metrics 연동 기반의 BackendFactory 설정
	backendFactory = mc.BackendFactory("backend", backendFactory)

	// Opencensus 연동 기반의 BackendFactory 설정
	backendFactory = opencensus.BackendFactory(backendFactory)
	return backendFactory
}
