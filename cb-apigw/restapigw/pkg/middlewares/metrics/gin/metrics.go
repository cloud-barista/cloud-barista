package gin

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/cloud-barista/cb-apigw/restapigw/pkg/config"
	"github.com/cloud-barista/cb-apigw/restapigw/pkg/logging"
	"github.com/cloud-barista/cb-apigw/restapigw/pkg/middlewares/metrics"
	"github.com/cloud-barista/cb-apigw/restapigw/pkg/proxy"
	ginRouter "github.com/cloud-barista/cb-apigw/restapigw/pkg/router/gin"
	"github.com/gin-gonic/gin"
	gometricsExp "github.com/rcrowley/go-metrics/exp"
)

// ===== [ Constants and Variables ] =====

// ===== [ Types ] =====

// ginResponseWriter - Gin의 ResponseWriter에 Metrics 처리를 위한 구조 정의
type ginResponseWriter struct {
	gin.ResponseWriter
	name  string
	begin time.Time
	rm    *metrics.RouterMetrics
}

// Metrics - gin 기반의 모든 Metrics 를 관리하는 구조 정의
type Metrics struct {
	*metrics.Metrics
}

// ===== [ Implementations ] =====

// end - gin.ResponseWriter 와 Metrics를 end 시점 처리
func (grw *ginResponseWriter) end() {
	duration := time.Since(grw.begin)
	grw.rm.Counter("response", grw.name, "status", strconv.Itoa(grw.Status()), "count").Inc(1)
	grw.rm.Histogram("response", grw.name, "size").Update(int64(grw.Size()))
	grw.rm.Histogram("response", grw.name, "time").Update(int64(duration))
}

// HandlerFactory - 전달된 HandlerFactory 수행 전에 필요한 Metric 관련 처리를 수행하는 HandlerFactory 구성
func (m *Metrics) HandlerFactory(hf ginRouter.HandlerFactory, logger logging.Logger) ginRouter.HandlerFactory {
	if m.Config == nil || !m.Config.RouterEnabled {
		return hf
	}

	return NewHTTPHandlerFactory(m.Router, hf, logger)
}

// RunEndpoint - Gin Engine 기반으로 http.Server 구동
func (m *Metrics) RunEndpoint(ctx context.Context, engine *gin.Engine, logger logging.Logger) {
	server := &http.Server{
		Addr:    m.Config.ListenAddress,
		Handler: engine,
	}

	go func() {
		// http.Server 실행 및 오류 logging
		logger.Error(server.ListenAndServe())
	}()

	go func() {
		// http.Server 종료 처리
		<-ctx.Done()
		// shutting down the stats handler
		ctx, cancel := context.WithTimeout(ctx, time.Second)
		server.Shutdown(ctx)
		cancel()
	}()
}

// NewEngine - Stats 처리를 위한 Endpoint 역할을 담당하는 Gin Engine 생성
func (m *Metrics) NewEngine(debugMode bool, logger logging.Logger) *gin.Engine {
	// Sets up the Gin engine for exports the metrics
	if debugMode {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.RedirectTrailingSlash = true
	engine.RedirectFixedPath = true
	engine.HandleMethodNotAllowed = true

	engine.GET("/metrics", m.NewExportHandler())
	return engine
}

// NewExportHandler - 수집된 Metrics를 JSON 포맷으로 노출하는 go-metrics에서 제공되는 http.Handler 생성
func (m *Metrics) NewExportHandler() gin.HandlerFunc {
	return gin.WrapH(gometricsExp.ExpHandler(*m.Registry))
}

// ===== [ Private Functions ] =====

// ===== [ Public Functions ] =====

// NewHTTPHandlerFactory - RouterMetrics 처리를 수행하는 HandlerFactory 생성
func NewHTTPHandlerFactory(rm *metrics.RouterMetrics, hf ginRouter.HandlerFactory, logger logging.Logger) ginRouter.HandlerFactory {
	return func(eConf *config.EndpointConfig, p proxy.Proxy) gin.HandlerFunc {
		next := hf(eConf, p)
		// Endpoint 응답관련 Metric 처리기 등록
		//rm.RegisterResponseWriterMetrics(eConf.Endpoint)
		return func(c *gin.Context) {
			name := eConf.Endpoint
			// Bypass인 경우 실제 호출 URL로 처리
			if eConf.IsBypass {
				name = c.Request.URL.Path
			}
			rm.RegisterResponseWriterMetrics(name)

			// Metric 처리기를 반영한 Gin Response Writer 생성
			rw := &ginResponseWriter{c.Writer, name, time.Now(), rm}
			c.Writer = rw
			// Router 연결 Metric 처리
			rm.Connection()

			next(c)

			rw.end()
			// Router 종료 Metric 처리
			rm.Disconnection()
		}
	}
}

// SetupAndRun - Gin 기반으로 동작하는 Metric Producer를 생성하고 Collector 처리를 위한 Gin 기반의 Endpoint Server 구동
func SetupAndRun(ctx context.Context, sConf config.ServiceConfig, logger logging.Logger) *Metrics {
	// Metrics Producer 설정 및 생성
	metricProducer := Metrics{metrics.SetupAndCreate(ctx, sConf, logger)}

	if metricProducer.Config != nil && metricProducer.Config.ExposeMetrics {
		// Gin 기반 Server 구동 및 Endpoint 처리 설정
		metricProducer.RunEndpoint(ctx, metricProducer.NewEngine(sConf.Debug, logger), logger)
	}

	return &metricProducer
}
