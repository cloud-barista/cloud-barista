// Package admin -
package admin

import (
	"github.com/cloud-barista/cb-apigw/restapigw/pkg/api"
	"github.com/cloud-barista/cb-apigw/restapigw/pkg/config"
	"github.com/cloud-barista/cb-apigw/restapigw/pkg/logging"
)

// ===== [ Constants and Variables ] =====
// ===== [ Types ] =====

type (
	// Option - Admin Server 운영을 위한 옵션 설정 함수 형식
	Option func(*Server)
)

// ===== [ Implementations ] =====
// ===== [ Private Functions ] =====
// ===== [ Public Functions ] =====

// WithLog - Logging에 사용할 Logger 적용
func WithLog(log logging.Logger) Option {
	return func(s *Server) {
		s.logger = log
	}
}

// WithConfigurations - Memory상에 동작하고 있는 Configuration 옵션 설정
func WithConfigurations(configs *api.Configuration) Option {
	return func(s *Server) {
		s.apiHandler.Configs = configs
	}
}

// WithPort - 서비스를 위한 포트 옵션 설정
func WithPort(port int) Option {
	return func(s *Server) {
		s.Port = port
	}
}

// WithCredentials - 서비스 사용을 위한 사용자 옵션 설정
func WithCredentials(credential *config.CredentialsConfig) Option {
	return func(s *Server) {
		s.Credentials = credential
	}
}

// WithTLS - TLS 옵션 설정
func WithTLS(tls *config.TLSConfig) Option {
	return func(s *Server) {
		s.TLS = tls
	}
}

// WithProfiler - Profiler 적용 옵션 설정
func WithProfiler(enabled, public bool) Option {
	return func(s *Server) {
		s.profilingEnabled = enabled
		s.profilingPublic = public
	}
}
