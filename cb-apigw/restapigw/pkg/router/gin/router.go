package gin

import (
	"context"
	"net/http"
	"strings"

	"github.com/cloud-barista/cb-apigw/restapigw/pkg/config"
	"github.com/cloud-barista/cb-apigw/restapigw/pkg/core"
	"github.com/cloud-barista/cb-apigw/restapigw/pkg/logging"
	"github.com/cloud-barista/cb-apigw/restapigw/pkg/proxy"
	"github.com/cloud-barista/cb-apigw/restapigw/pkg/router"
	httpServer "github.com/cloud-barista/cb-apigw/restapigw/pkg/transport/http/server"
	"github.com/gin-gonic/gin"
)

// ===== [ Constants and Variables ] =====

// ===== [ Types ] =====

// PipeConfig - 서비스  운영에 필요한 Pipeline 을 구성하기 위한 구조
type PipeConfig struct {
	Context        context.Context
	Engine         *gin.Engine
	Middlewares    []gin.HandlerFunc
	HandlerFactory HandlerFactory
	ProxyFactory   proxy.Factory
	Logger         *logging.Logger
}

// ===== [ Implementations ] =====

// Run - 서비스 설정을 기준으로 Gin Engine 구성 / Route - Endpoint Handler 추가 / HTTP Server 구동
func (pc PipeConfig) Run(sConf config.ServiceConfig) {
	if !sConf.Debug {
		gin.SetMode(gin.ReleaseMode)
	} else {
		pc.Logger.Debug("Debug enabled")
	}

	pc.Engine.Use(pc.Middlewares...)

	if sConf.Debug {
		// Debug Endpoint Handler 구성
		pc.registerDebugEndpoints()
	}

	// Endpoint Handler 구성
	pc.registerEndpoints(sConf.Endpoints)

	pc.Engine.NoRoute(func(c *gin.Context) {
		c.Header(router.CompleteResponseHeaderName, router.HeaderIncompleteResponseValue)
	})

	// HTTP Server 환경 초기화
	httpServer.InitHTTPDefaultTransport(sConf)
	// Gin Engine을 Handler로 사용하는 HTTP Server 구동
	if err := httpServer.RunServer(pc.Context, sConf, pc.Engine); err != nil {
		pc.Logger.Error(err.Error())
	}
}

// registerDebugEndpoints - debug 옵션이 활성화된 경우에 `__debug` 로 호출되는 Enpoint Handler 구성
func (pc PipeConfig) registerDebugEndpoints() {
	handler := DebugHandler(*pc.Logger)
	pc.Engine.GET("/__debug/*param", handler)
	pc.Engine.POST("/__debug/*param", handler)
	pc.Engine.PUT("/__debug/*param", handler)
}

// registerEndpoints - 지정한 Endpoint 설정들을 기준으로 Endpoint Handler 구성
func (pc PipeConfig) registerEndpoints(eConfs []*config.EndpointConfig) {
	for _, eConf := range eConfs {
		// Endpoint에 연결되어 동작할 수 있도록 ProxyFactory의 Call chain에 대한 인스턴스 생성 (ProxyStack)
		proxyStack, err := pc.ProxyFactory.New(eConf)
		if err != nil {
			pc.Logger.Error("calling the ProxyFactory", err.Error())
			continue
		}

		if eConf.IsBypass {
			// Bypass case
			pc.registerGroup(eConf.Endpoint, pc.HandlerFactory(eConf, proxyStack), len(eConf.Backend))
		} else {
			// Normal case
			pc.registerEndpoint(eConf.Method, eConf.Endpoint, pc.HandlerFactory(eConf, proxyStack), len(eConf.Backend))
		}
	}
}

// registerGroup - Bypass인 경우는 Group 단위로 Gin Engine에 Endpoint Handler 등록
func (pc PipeConfig) registerGroup(path string, handler gin.HandlerFunc, totBackends int) {
	if totBackends > 1 {
		pc.Logger.Error("Bypass endpoint must have a single backend! Ignoring", path)
		return
	}

	// Bypass에 적합한 Group 정보 조정 및 Route 등록
	suffix := "/" + core.Bypass
	group := strings.TrimSuffix(path, suffix)

	groupRoute := pc.Engine.Group(group)
	groupRoute.Any(suffix, handler)
}

// registerEndpoint - 지정한 정보를 기준으로 Gin Engine에 Endpoint Handler 등록
func (pc PipeConfig) registerEndpoint(method, path string, handler gin.HandlerFunc, totBackends int) {
	method = strings.ToTitle(method)
	if method != http.MethodGet && totBackends > 1 {
		pc.Logger.Error(method, "endpoints must have a single backend! Ignoring", path)
		return
	}
	switch method {
	case http.MethodGet:
		pc.Engine.GET(path, handler)
	case http.MethodPost:
		pc.Engine.POST(path, handler)
	case http.MethodPut:
		pc.Engine.PUT(path, handler)
	case http.MethodPatch:
		pc.Engine.PATCH(path, handler)
	case http.MethodDelete:
		pc.Engine.DELETE(path, handler)
	default:
		pc.Logger.Error("Unsupported method", method)
	}
}

// ===== [ Private Functions ] =====

// ===== [ Public Functions ] =====
