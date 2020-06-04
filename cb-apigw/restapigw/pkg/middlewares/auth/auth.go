// Package auth - API Gateway 접근인증을 위한 기능 제공
package auth

import (
	"github.com/cloud-barista/cb-apigw/restapigw/pkg/logging"
	ginAuth "github.com/cloud-barista/cb-apigw/restapigw/pkg/middlewares/auth/gin"
	ginRouter "github.com/cloud-barista/cb-apigw/restapigw/pkg/router/gin"
)

// ===== [ Constants and Variables ] =====

// ===== [ Types ] =====

// ===== [ Implementations ] =====

// ===== [ Private Functions ] =====

// ===== [ Public Functions ] =====

// HandlerFactory - Auth 기능을 수행하는 Route Handler Factory 구성
func HandlerFactory(hf ginRouter.HandlerFactory, logger logging.Logger) ginRouter.HandlerFactory {
	return ginAuth.TokenValidator(hf, logger)
}
