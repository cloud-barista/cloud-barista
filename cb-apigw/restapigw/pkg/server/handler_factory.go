// Package server -
package server

import (
	"github.com/cloud-barista/cb-apigw/restapigw/pkg/logging"
	"github.com/cloud-barista/cb-apigw/restapigw/pkg/middlewares/auth"
	ginMetrics "github.com/cloud-barista/cb-apigw/restapigw/pkg/middlewares/metrics/gin"
	ginOpencensus "github.com/cloud-barista/cb-apigw/restapigw/pkg/middlewares/opencensus/router/gin"
	ginRateLimit "github.com/cloud-barista/cb-apigw/restapigw/pkg/middlewares/ratelimit/router/gin"
	ginRouter "github.com/cloud-barista/cb-apigw/restapigw/pkg/router/gin"
)

// ===== [ Constants and Variables ] =====
// ===== [ Types ] =====
// ===== [ Implementations ] =====
// ===== [ Private Functions ] =====
// ===== [ Public Functions ] =====

// setupGinHandlerFactory - Gin Router 처리를 위한 핸들러 설정
func setupGinHandlerFactory(logger logging.Logger, mc *ginMetrics.Collector) ginRouter.HandlerFactory {
	// Rate Limit 처리용 Router Handler 구성
	handlerFactory := ginRateLimit.HandlerFactory(ginRouter.EndpointHandler, logger)

	// TODO: JWT Auth, JWT Rejector 처리용 Router Handler 구성

	// 임시로 HMAC을 활용한 Auth 인증 처리용 RouteHandlerFactory 구성
	handlerFactory = auth.HandlerFactory(handlerFactory, logger)

	// metrics Collector를 활용하는 RouteHandlerFactory 구성
	handlerFactory = mc.HandlerFactory(handlerFactory, logger)

	// Opencensus를 활용하는 RouteHandlerFactory 구성
	handlerFactory = ginOpencensus.HandlerFactory(handlerFactory, logger)

	return handlerFactory
}
