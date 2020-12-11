// Package gin - GIN 기반의 Routing 기능 제공 패키지
package gin

import (
	"net/http"
	"strings"

	"github.com/cloud-barista/cb-apigw/restapigw/pkg/config"
	"github.com/cloud-barista/cb-apigw/restapigw/pkg/core"
	"github.com/cloud-barista/cb-apigw/restapigw/pkg/logging"
	cors "github.com/cloud-barista/cb-apigw/restapigw/pkg/middlewares/cors/gin"
	httpsecure "github.com/cloud-barista/cb-apigw/restapigw/pkg/middlewares/httpsecure/gin"
	"github.com/cloud-barista/cb-apigw/restapigw/pkg/proxy"
	"github.com/cloud-barista/cb-apigw/restapigw/pkg/router"
	"github.com/gin-gonic/gin"
)

// ===== [ Constants and Variables ] =====
// ===== [ Types ] =====

type (
	// Option - Server 인스턴스에 옵션을 설정하는 함수 형식
	Option func(*PipeConfig)

	// PipeConfig - 서비스  운영에 필요한 Pipeline 을 구성하기 위한 구조
	PipeConfig struct {
		engine         *router.DynamicRouter
		middlewares    []gin.HandlerFunc
		handlerFactory HandlerFactory
		proxyFactory   proxy.Factory
		logger         logging.Logger
	}
)

// ===== [ Implementations ] =====

// createRouter - Gin 기반의 Router 생성
func (pc *PipeConfig) createRouter(sConf *config.ServiceConfig) http.Handler {
	if !sConf.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.New()
	engine.Use(gin.Recovery())

	engine.RedirectTrailingSlash = true
	engine.RedirectFixedPath = true
	engine.HandleMethodNotAllowed = true

	// CORS Middleware 반영
	cors.New(sConf.Middleware, engine)

	// HTTPSecure Middleware 반영
	if err := httpsecure.Register(sConf.Middleware, engine); nil != err {
		pc.logger.Warning(err)
	}

	// TODO: 사전 처리되어야할 Middleware 추가는 여기서...

	// 기본 설정
	engine.Use(pc.middlewares...)
	if sConf.Debug {
		engine.Any("/__debug/*param", DebugHandler(pc.logger))
	}
	engine.NoRoute(func(c *gin.Context) {
		c.Header(router.CompleteResponseHeaderName, router.HeaderIncompleteResponseValue)
	})

	return engine
}

// createEngine - Gin Router Engine 인스턴스 생성
func (pc *PipeConfig) createEngine(sConf *config.ServiceConfig) {
	de := &router.DynamicRouter{}
	de.SetHandler(pc.createRouter(sConf))

	pc.engine = de
}

// registerAPIGroup - Bypass인 경우는 Group 단위로 Gin Engine에 Endpoint Handler 등록
func (pc PipeConfig) registerAPIGroup(path string, handler gin.HandlerFunc, totBackends int) {
	if 1 < totBackends {
		pc.logger.Errorf("[API G/W] Router > Bypass endpoint must have a single backend! Ignoring -> %s", path)
		return
	}

	// Bypass에 적합한 Group 정보 조정 및 Route 등록
	suffix := "/" + core.Bypass
	group := strings.TrimSuffix(path, suffix)

	engine := pc.engine.GetHandler().(*gin.Engine)

	groupRoute := engine.Group(group)
	groupRoute.Any(suffix, handler)
}

// registerAPI - 지정한 정보를 기준으로 Gin Engine에 Endpoint Handler 등록
func (pc PipeConfig) registerAPI(method, path string, handler gin.HandlerFunc, totBackends int) {
	method = strings.ToTitle(method)
	if method != http.MethodGet && 1 < totBackends {
		pc.logger.Errorf("[API G/W] Router > Method: %s, endpoints must have a single backend! Ignoring -> %s", method, path)
		return
	}

	engine := pc.engine.GetHandler().(*gin.Engine)

	switch method {
	case http.MethodGet:
		engine.GET(path, handler)
	case http.MethodPost:
		engine.POST(path, handler)
	case http.MethodPut:
		engine.PUT(path, handler)
	case http.MethodPatch:
		engine.PATCH(path, handler)
	case http.MethodDelete:
		engine.DELETE(path, handler)
	default:
		pc.logger.Errorf("[API G/W] Router > Unsupported method -> %s", method)
	}
}

// UpdateEngine - API 변경 적용을 위한 Gin Engine 재 생성
func (pc *PipeConfig) UpdateEngine(sConf *config.ServiceConfig) {
	gin := pc.createRouter(sConf)
	pc.engine.SetHandler(gin)
}

// Engine - Router 기능을 처리하는 Gin Engine 반환 (http.Handler)
func (pc *PipeConfig) Engine() http.Handler {
	return pc.engine
}

// RegisterAPIs - API Provider (Repository)에서 추출된 API 설정들을 Router로 등록
func (pc *PipeConfig) RegisterAPIs(sConf *config.ServiceConfig, defs []*config.EndpointConfig) error {
	pc.logger.Info("[API G/W] Loading API Endpoints")

	for _, def := range defs {
		// 정보 재구성
		err := def.AdjustValues(sConf)
		if nil != err {
			pc.logger.WithError(err).Error("[API G/W] Router > Can not adjust values for definition")
		}

		// 활성화된 경우만 적용
		if def.Active {
			// Endpoint에 연결되어 동작할 수 있도록 ProxyFactory의 Call chain에 대한 인스턴스 생성 (ProxyStack)
			proxyStack, err := pc.proxyFactory.New(def)
			if nil != err {
				pc.logger.WithError(err).Error("[API G/W] Router > Can not calling the ProxyFactory")
				continue
			}

			if def.IsBypass {
				// Bypass case
				pc.registerAPIGroup(def.Endpoint, pc.handlerFactory(def, proxyStack), len(def.Backend))
			} else {
				// Normal case
				pc.registerAPI(def.Method, def.Endpoint, pc.handlerFactory(def, proxyStack), len(def.Backend))
			}
		}
	}

	pc.logger.Info("[API G/W] API Endpoints loaded")
	return nil
}

// ===== [ Private Functions ] =====

// ===== [ Public Functions ] =====

// WithLogger - Logger 인스턴스 설정
func WithLogger(logger logging.Logger) Option {
	return func(pc *PipeConfig) {
		pc.logger = logger
	}
}

// WithHandlerFactory - Route Handler Factory 설정
func WithHandlerFactory(hf HandlerFactory) Option {
	return func(pc *PipeConfig) {
		pc.handlerFactory = hf
	}
}

// WithProxyFactory - Proxy Factory 설정
func WithProxyFactory(pf proxy.Factory) Option {
	return func(pc *PipeConfig) {
		pc.proxyFactory = pf
	}
}

// New - PipeConfig 인스턴스 생성
func New(sConf *config.ServiceConfig, opts ...Option) router.Router {
	pc := PipeConfig{}

	for _, opt := range opts {
		opt(&pc)
	}

	pc.createEngine(sConf)

	return &pc
}
