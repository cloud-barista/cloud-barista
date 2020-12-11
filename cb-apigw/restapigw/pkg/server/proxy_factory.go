// Package server -
package server

import (
	"github.com/cloud-barista/cb-apigw/restapigw/pkg/logging"
	ginMetrics "github.com/cloud-barista/cb-apigw/restapigw/pkg/middlewares/metrics/gin"
	"github.com/cloud-barista/cb-apigw/restapigw/pkg/middlewares/opencensus"
	"github.com/cloud-barista/cb-apigw/restapigw/pkg/proxy"
)

// ===== [ Constants and Variables ] =====
// ===== [ Types ] =====
// ===== [ Implementations ] =====
// ===== [ Private Functions ] =====
// ===== [ Public Functions ] =====

// setupGinProxyFactory - 지정한 Backend Factory를 연계하는 Proxy Factory 설정
func setupGinProxyFactory(logger logging.Logger, bf proxy.BackendFactory, mc *ginMetrics.Collector) proxy.Factory {
	proxyFactory := proxy.NewDefaultFactory(bf, logger)

	// Metrics 연동 기반의 ProxyFactory 설정
	proxyFactory = mc.ProxyFactory("pipe", proxyFactory)

	// Opencensus 연동 기반의 ProxyFactory 설정
	proxyFactory = opencensus.ProxyFactory(proxyFactory)
	return proxyFactory
}
