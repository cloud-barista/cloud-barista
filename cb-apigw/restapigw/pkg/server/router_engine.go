// Package server -
package server

import (
	"context"

	"github.com/cloud-barista/cb-apigw/restapigw/pkg/config"
	"github.com/cloud-barista/cb-apigw/restapigw/pkg/logging"
	ginMetrics "github.com/cloud-barista/cb-apigw/restapigw/pkg/middlewares/metrics/gin"
	"github.com/cloud-barista/cb-apigw/restapigw/pkg/middlewares/metrics/influxdb"
	"github.com/cloud-barista/cb-apigw/restapigw/pkg/middlewares/opencensus"
	"github.com/cloud-barista/cb-apigw/restapigw/pkg/router"
	ginRouter "github.com/cloud-barista/cb-apigw/restapigw/pkg/router/gin"
)

// ===== [ Constants and Variables ] =====
// ===== [ Types ] =====
// ===== [ Implementations ] =====
// ===== [ Private Functions ] =====

// setupGinRouter - Gin 기반으로 동작하는 Router 설정
func setupGinRouter(ctx context.Context, sConf *config.ServiceConfig, logger logging.Logger) router.Router {
	// Setup the metrics
	mc := ginMetrics.New(ctx, sConf.Middleware, logger, sConf.Debug)

	// Setup the InfluxDB client for metrics
	if nil != mc.Config {
		influxdb.SetupAndRun(ctx, mc.Config.InfluxDB, func() interface{} { return mc.Snapshot() }, &logger)
	} else {
		logger.Warn("Skip the influxdb setup and running because the no metrics configuration or incorrect.")
	}

	// Setup the Opencensus
	if err := opencensus.Setup(ctx, *sConf); nil != err {
		logger.Fatal(err)
	}

	// API G/W 구동을 위한 Factory 구성 반환
	return ginRouter.New(sConf,
		ginRouter.WithLogger(logger),
		ginRouter.WithHandlerFactory(setupGinHandlerFactory(logger, mc)),
		ginRouter.WithProxyFactory(setupGinProxyFactory(logger, setupGinBackendFactoryWithContext(ctx, logger, mc), mc)),
	)
}

// ===== [ Public Functions ] =====

// SetupRouter - API G/W 운영을 위한 Router 설정
func SetupRouter(ctx context.Context, sConf *config.ServiceConfig, logger logging.Logger) router.Router {
	switch sConf.RouterEngine {
	default:
		return setupGinRouter(ctx, sConf, logger)
	}
}
