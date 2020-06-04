package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/cloud-barista/cb-apigw/restapigw/pkg/config"
	"github.com/cloud-barista/cb-apigw/restapigw/pkg/logging"
	"github.com/cloud-barista/cb-apigw/restapigw/pkg/middlewares/auth"
	ginCors "github.com/cloud-barista/cb-apigw/restapigw/pkg/middlewares/cors/gin"
	httpcache "github.com/cloud-barista/cb-apigw/restapigw/pkg/middlewares/httpcache"
	httpsecure "github.com/cloud-barista/cb-apigw/restapigw/pkg/middlewares/httpsecure/gin"
	ginMetrics "github.com/cloud-barista/cb-apigw/restapigw/pkg/middlewares/metrics/gin"
	"github.com/cloud-barista/cb-apigw/restapigw/pkg/middlewares/metrics/influxdb"
	opencensus "github.com/cloud-barista/cb-apigw/restapigw/pkg/middlewares/opencensus"
	ginOpencensus "github.com/cloud-barista/cb-apigw/restapigw/pkg/middlewares/opencensus/router/gin"
	ratelimitProxy "github.com/cloud-barista/cb-apigw/restapigw/pkg/middlewares/ratelimit/proxy"
	ratelimit "github.com/cloud-barista/cb-apigw/restapigw/pkg/middlewares/ratelimit/router/gin"
	"github.com/cloud-barista/cb-apigw/restapigw/pkg/proxy"
	ginRouter "github.com/cloud-barista/cb-apigw/restapigw/pkg/router/gin"
	"github.com/cloud-barista/cb-apigw/restapigw/pkg/transport/http/client"
	"github.com/gin-gonic/gin"

	// Opencensus 연동을 위한 Exporter 로드 및 초기화
	_ "github.com/cloud-barista/cb-apigw/restapigw/pkg/middlewares/opencensus/exporters/jaeger"
)

// ===== [ Constants and Variables ] =====

// ===== [ Types ] =====

// ===== [ Implementations ] =====

// ===== [ Private Functions ] =====

// contextWithSignal - System Interrupt Signal을 반영한 Context 생성
func contextWithSignal(ctx context.Context) context.Context {
	newCtx, cancel := context.WithCancel(ctx)
	signals := make(chan os.Signal)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		select {
		case <-signals:
			cancel()
			close(signals)
		}
	}()
	return newCtx
}

// newHandlerFactory - Middleware들과 Opencensus 처리를 반영한 Gin Endpoint Handler 생성
func newHandlerFactory(logger logging.Logger, metricsProducer *ginMetrics.Metrics) ginRouter.HandlerFactory {
	// Rate Limit 처리용 RouterHandlerFactory 구성
	handlerFactory := ratelimit.HandlerFactory(ginRouter.EndpointHandler, logger)
	// TODO: JWT Auth, JWT Rejector

	// 임시로 HMAC을 활용한 Auth 인증 처리용 RouteHandlerFactory 구성
	handlerFactory = auth.HandlerFactory(handlerFactory, logger)

	// metricsProducer 활용하는 RouteHandlerFactory 구성
	handlerFactory = metricsProducer.HandlerFactory(handlerFactory, logger)
	// Opencensus를 활용하는 RouteHandlerFactory 구성
	handlerFactory = ginOpencensus.HandlerFactory(handlerFactory, logger)
	return handlerFactory
}

// newEngine - HTTP Server 운영과 연관되는 Middleware 구성을 처리한 Gin Engine 생성
func newEngine(sConf config.ServiceConfig, logger logging.Logger) *gin.Engine {
	engine := gin.Default()
	engine.RedirectTrailingSlash = true
	engine.RedirectFixedPath = true
	engine.HandleMethodNotAllowed = true

	// CORS Middleware 반영
	ginCors.New(sConf.Middleware, engine)

	// HTTPSecure Middleware 반영
	if err := httpsecure.Register(sConf.Middleware, engine); err != nil {
		logger.Warning(err)
	}

	return engine
}

// newProxyFactory - 지정된 BackendFactory를 기반으로 동작하는 ProxyFactory 생성
func newProxyFactory(logger logging.Logger, backendFactory proxy.BackendFactory, metricsProducer *ginMetrics.Metrics) proxy.Factory {
	proxyFactory := proxy.NewDefaultFactory(backendFactory, logger)

	// Metrics 연동 기반의 ProxyFactory 설정
	proxyFactory = metricsProducer.ProxyFactory("pipe", proxyFactory)

	// Opencensus 연동 기반의 ProxyFactory 설정
	proxyFactory = opencensus.ProxyFactory(proxyFactory)
	return proxyFactory
}

// newBackendFactoryWithContext - 지정된 Context 기반에서 활용가능한 middleware들이 적용된 BackendFactory 생성
func newBackendFactoryWithContext(ctx context.Context, logger logging.Logger, metricsProducer *ginMetrics.Metrics) proxy.BackendFactory {
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
	backendFactory = metricsProducer.BackendFactory("backend", backendFactory)

	// Opencensus 연동 기반의 BackendFactory 설정
	backendFactory = opencensus.BackendFactory(backendFactory)
	return backendFactory
}

// ===== [ Public Functions ] =====

// SetupAndRun - API Gateway 서비스를 위한 Router 및 Pipeline 구성과 HTTP Server 구동
func SetupAndRun(ctx context.Context, sConf config.ServiceConfig) error {
	// Sets up the Logger (CB-LOG)
	logger := logging.NewLogger()

	// Sets up the Metrics
	metricsProducer := ginMetrics.SetupAndRun(ctx, sConf, *logger)

	if metricsProducer.Config != nil {
		// Sets up the InfluxDB Client for Metrics
		influxdb.SetupAndRun(ctx, metricsProducer.Config.InfluxDB, func() interface{} { return metricsProducer.Snapshot() }, logger)
	} else {
		logger.Warn("Skip the influxdb setup and running because the no metrics configuration or incorrect")
	}

	// Sets up the Opencensus
	if err := opencensus.Setup(ctx, sConf); err != nil {
		logger.Fatal(err)
	}

	// Sets up the Pipeline (Router (Endpoint Handler) - Proxy - Backend)
	pipeConfig := ginRouter.PipeConfig{
		Engine:         newEngine(sConf, *logger),
		ProxyFactory:   newProxyFactory(*logger, newBackendFactoryWithContext(ctx, *logger, metricsProducer), metricsProducer),
		Middlewares:    []gin.HandlerFunc{},
		Logger:         logger,
		HandlerFactory: newHandlerFactory(*logger, metricsProducer),
		Context:        contextWithSignal(ctx),
	}

	// PipeConfig 정보를 기준으로 HTTP Server 실행 (Gin Router + Endpoint Handler, Pipeline)
	pipeConfig.Run(sConf)

	return nil
}
