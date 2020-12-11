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

type (
	// Collector - Metrics 기반의 Collector 구조
	Collector struct {
		*metrics.Producer
	}

	// responseWriter - Gin 의 ResponseWriter에 Metrics 처리를 위한 구조
	responseWriter struct {
		gin.ResponseWriter
		name  string
		begin time.Time
		rm    *metrics.RouterMetrics
	}
)

// ===== [ Implementations ] =====

// end - gin.ResponseWriter 와 Metrics를 end 시점 처리
func (rw *responseWriter) end() {
	duration := time.Since(rw.begin)

	rw.rm.Counter("response", rw.name, "status", strconv.Itoa(rw.Status()), "count").Inc(1)
	rw.rm.Histogram("response", rw.name, "size").Update(int64(rw.Size()))
	rw.rm.Histogram("response", rw.name, "time").Update(int64(duration))
}

// HandlerFactory - 전달된 HandlerFactory 수행 전에 필요한 Metric 관련 처리를 수행하는 HandlerFactory 구성
func (c *Collector) HandlerFactory(hf ginRouter.HandlerFactory, log logging.Logger) ginRouter.HandlerFactory {
	if nil == c.Config || !c.Config.RouterEnabled {
		return hf
	}

	return NewHTTPHandlerFactory(c.Router, hf, log)
}

// RunEndpoint - Gin Engine 기반으로 http.Server 구동
func (c *Collector) RunEndpoint(ctx context.Context, engine *gin.Engine, log logging.Logger) {
	server := &http.Server{
		Addr:    c.Config.ListenAddress,
		Handler: engine,
	}

	go func() {
		// http.Server 실행 및 오류 logging
		log.Error(server.ListenAndServe())
	}()

	go func() {
		// http.Server 종료 처리
		<-ctx.Done()
		// shutting down the stats handler
		log.Info("[METRICS] shutting down the metrics stats handler.")
		ctx, cancel := context.WithTimeout(ctx, time.Second)
		server.Shutdown(ctx)
		cancel()
	}()
}

// NewEngine - Stats 처리를 위한 Endpoint 역할을 담당하는 Gin Engine 생성
func (c *Collector) NewEngine(debugMode bool) *gin.Engine {
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

	engine.GET("/metrics", c.NewExportHandler())
	return engine
}

// NewExportHandler - 수집된 Metrics를 JSON 포맷으로 노출하는 go-metrics에서 제공되는 http.Handler 생성
func (c *Collector) NewExportHandler() gin.HandlerFunc {
	return gin.WrapH(gometricsExp.ExpHandler(*c.Registry))
}

// ===== [ Private Functions ] =====

// ===== [ Public Functions ] =====

// NewHTTPHandlerFactory - Router 단위의 Metrics 처리를 수행하는 HandlerFactory 생성
func NewHTTPHandlerFactory(rm *metrics.RouterMetrics, hf ginRouter.HandlerFactory, log logging.Logger) ginRouter.HandlerFactory {
	return func(eConf *config.EndpointConfig, p proxy.Proxy) gin.HandlerFunc {
		next := hf(eConf, p)

		return func(c *gin.Context) {
			name := eConf.Endpoint
			// Bypass인 경우 실제 호출 URL로 처리
			if eConf.IsBypass {
				name = c.Request.URL.Path
			}

			// Endpoint 응답관련 Metric 처리기 등록
			rm.RegisterResponseWriterMetrics(name)

			// Metric 처리기를 반영한 Gin Response Writer 생성
			rw := &responseWriter{c.Writer, name, time.Now(), rm}
			c.Writer = rw

			// Router 연결 Metric 처리
			rm.Connection(c.Request.TLS)

			next(c)

			rw.end()

			// Router 종료 Metric 처리
			rm.Disconnection()
		}
	}
}

// New -  Metric Collector를 생성하고 Collection 처리를 위한 Endpoint Server 구동
func New(ctx context.Context, mConf config.MWConfig, log logging.Logger, debugMode bool) *Collector {
	mc := Collector{metrics.New(ctx, mConf, log)}

	// 설정에 따라서 Metrics에 대한 Endpoint 설정
	if nil != mc.Config && mc.Config.ExposeMetrics {
		mc.RunEndpoint(ctx, mc.NewEngine(debugMode), log)
	}

	return &mc
}
