// Package admin - Admin API 기능 지원 패키지
package admin

import (
	"fmt"
	"net/http"
	"net/http/pprof"

	"github.com/cloud-barista/cb-apigw/restapigw/pkg/api"
	"github.com/cloud-barista/cb-apigw/restapigw/pkg/config"

	ginAdapter "github.com/cloud-barista/cb-apigw/restapigw/pkg/core/adapters/gin"
	"github.com/cloud-barista/cb-apigw/restapigw/pkg/jwt"
	"github.com/cloud-barista/cb-apigw/restapigw/pkg/logging"
	"github.com/gin-gonic/gin"
	"github.com/rs/cors"
	ginCors "github.com/rs/cors/wrapper/gin"
)

// ===== [ Constants and Variables ] =====

const (
	// APIBasePath - Admin API 관리용 기본 Path
	APIBasePath = "/apis"
	// GroupBasePath - Group 관리용 기본 Path
	GroupBasePath = APIBasePath + "/group/"
)

// ===== [ Types ] =====

type (
	// Server - Admin API Server 관리 정보 형식
	Server struct {
		ConfigurationChan chan api.ConfigChangedMessage

		apiHandler *APIHandler
		logger     logging.Logger

		Port             int `mapstructure:"port"`
		Credentials      *config.CredentialsConfig
		TLS              *config.TLSConfig
		profilingEnabled bool
		profilingPublic  bool
	}
)

// ===== [ Implementations ] =====

// addInternalPublicRoutes - ADMIN Server의 외부에 노출되는 Routes 설정
func (s *Server) addInternalPublicRoutes(ge *gin.Engine) {
	// Admin Web Serve를 위한 Middleware 구성
	ge.Use(WebServe("/"))

	// Status Endpoints
	statusAPI := ge.Group("/status")
	statusAPI.GET("/", gin.WrapH(NewOverviewHandler(s.apiHandler.Configs, s.logger)))
	statusAPI.GET("/{name}", gin.WrapH(NewStatusHandler(s.apiHandler.Configs, s.logger)))
}

// addInternalAUthRoutes - ADMIN Server의 Auth 처리용 Routes 설정
func (s *Server) addInternalAuthRoutes(ge *gin.Engine, guard jwt.Guard) {
	handlers := jwt.Handler{Guard: guard}

	authGroup := ge.Group("/auth")
	{
		authGroup.POST("/login", gin.WrapH(handlers.Login(s.Credentials, s.logger)))
		authGroup.POST("/logout", gin.WrapH(handlers.Logout(s.Credentials, s.logger)))
		authGroup.GET("/refresh_token", gin.WrapH(handlers.Refresh()))
	}
}

// addInternalRoutes - API Definition 처리를 위한 Routes 설정
func (s *Server) addInternalRoutes(ge *gin.Engine, guard jwt.Guard) {
	s.logger.Debug("[ADMIN Server] Loading API Endpoints")

	// APIs endpoints
	groupAPI := ge.Group(APIBasePath)
	groupAPI.Use(ginAdapter.Wrap(jwt.NewMiddleware(guard).Handler))
	{
		// All group datas
		groupAPI.GET("/", gin.WrapH(s.apiHandler.GetDefinitions())) // Get All Groups (with definitions)
		groupAPI.PUT("/", gin.WrapH(s.apiHandler.ApplyGroups()))    // Apply Group changes to persistence
		// Groups
		groupAPI.POST("/group", gin.WrapH(s.apiHandler.AddGroup()))           // Add Group
		groupAPI.GET("/group/:gid", gin.WrapH(s.apiHandler.GetGroup()))       // Get Group
		groupAPI.DELETE("/group/:gid", gin.WrapH(s.apiHandler.RemoveGroup())) // Remove Group (all definitions in group)
		// Definitions for group
		groupAPI.POST("/group/:gid/definition", gin.WrapH(s.apiHandler.AddDefinition()))          // Add Definition
		groupAPI.PUT("/group/:gid/definition", gin.WrapH(s.apiHandler.UpdateDefinition()))        // Update Definition
		groupAPI.DELETE("/group/:gid/definition/:id", gin.WrapH(s.apiHandler.RemoveDefinition())) // Remove Definition
	}

	if s.profilingEnabled {
		groupProfiler := ge.Group("/debug/pprof")
		if !s.profilingPublic {
			groupProfiler.Use(ginAdapter.Wrap(jwt.NewMiddleware(guard).Handler))
		}
		{
			groupProfiler.GET("/", gin.WrapF(pprof.Index))
			groupProfiler.GET("/cmdline", gin.WrapF(pprof.Cmdline))
			groupProfiler.GET("/profile", gin.WrapF(pprof.Profile))
			groupProfiler.GET("/symbol", gin.WrapF(pprof.Symbol))
			groupProfiler.GET("/trace", gin.WrapF(pprof.Trace))
		}
	}
}

// isClosedChannel - 이미 채널이 종료되었는지 검증
func (s *Server) isClosedChannel(ch <-chan api.ConfigChangedMessage) bool {
	select {
	case <-ch:
		return true
	default:
	}
	return false
}

// listenAndServe - GIN Router를 기반으로 HTTP Server 구동
func (s *Server) listenAndServe(h http.Handler) error {
	address := fmt.Sprintf(":%v", s.Port)

	if s.TLS.IsHTTPS() {
		addressTLS := fmt.Sprintf(":%v", s.TLS.Port)
		if s.TLS.Redirect {
			go func() {
				s.logger.WithField("address", address).Info("[API SERVER] Listening HTTP redirects to HTTPS")
				s.logger.Fatal(http.ListenAndServe(address, RedirectHTTPS(s.TLS.Port)))
			}()
		}

		s.logger.WithField("address", addressTLS).Info("[API SERVER] LIstening HTTPS")
		return http.ListenAndServeTLS(addressTLS, s.TLS.PublicKey, s.TLS.PrivateKey, h)
	}

	s.logger.WithField("address", address).Info("[API SERVER] Certificate and certificate key where not found, defaulting to HTTP")
	return http.ListenAndServe(address, h)
}

// Start - Admin API Server 구동
func (s *Server) Start() error {
	s.logger.Info("[API SERVER] Admin API starting...")

	// Gin Router 생성
	engine := gin.Default()
	engine.RedirectTrailingSlash = true
	engine.RedirectFixedPath = true
	engine.HandleMethodNotAllowed = true

	// No Route 처리
	engine.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"code": "API_NOT_FOUND", "message": "API not found"})
	})

	// Router 설정 및 구동
	go s.listenAndServe(engine)

	s.AddRoutes(engine)

	s.logger.Info("[API SERVER] Admin API started.")
	return nil
}

// Stop - Admin API Server 종료
func (s *Server) Stop() {
	if nil == s {
		fmt.Println("server is null")
		return
	}

	if !s.isClosedChannel(s.ConfigurationChan) {
		close(s.ConfigurationChan)
	}
	s.logger.Info("[API G/W] Admin API stoped.")
}

// AddRoutes - ADMIN Routes 정보 구성
func (s *Server) AddRoutes(ge *gin.Engine) {
	guard := jwt.NewGuard(s.Credentials)

	// Cors 적용
	ge.Use(
		ginCors.New(cors.Options{
			AllowedOrigins:   []string{"*"},
			AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
			AllowedHeaders:   []string{"*"},
			AllowCredentials: true,
		}),
	)

	// Admin API Routes
	s.addInternalPublicRoutes(ge)
	s.addInternalAuthRoutes(ge, guard)
	s.addInternalRoutes(ge, guard)
}

// ===== [ Private Functions ] =====

// ===== [ Public Functions ] =====

// New - Admin Server 구동
func New(opts ...Option) *Server {
	configurationChan := make(chan api.ConfigChangedMessage)
	s := &Server{
		ConfigurationChan: configurationChan,
		apiHandler:        NewAPIHandler(configurationChan),
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}
