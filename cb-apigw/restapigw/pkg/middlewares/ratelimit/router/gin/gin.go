// Package gin - Gin Route (Endpoint) 처리 구간에 Rate Limit 기능을 제공하는 Gin Router 패키지
package gin

import (
	"net/http"
	"strings"

	"github.com/cloud-barista/cb-apigw/restapigw/pkg/config"
	"github.com/cloud-barista/cb-apigw/restapigw/pkg/core"
	"github.com/cloud-barista/cb-apigw/restapigw/pkg/logging"
	"github.com/cloud-barista/cb-apigw/restapigw/pkg/middlewares/ratelimit"
	"github.com/cloud-barista/cb-apigw/restapigw/pkg/middlewares/ratelimit/router"
	"github.com/cloud-barista/cb-apigw/restapigw/pkg/proxy"
	ginRouter "github.com/cloud-barista/cb-apigw/restapigw/pkg/router/gin"
	"github.com/gin-gonic/gin"
)

// ===== [ Constants and Variables ] =====

var (
	// logger - Logging
	logger = logging.NewLogger()
)

// ===== [ Types ] =====

// TokenExtractor - 각 요청에서 Token 정보를 추출하는 함수 형식
type TokenExtractor func(*gin.Context) string

// RateLimitMiddleware - Rate limit 처리가 적용된 Handler Func 반환 형식
type RateLimitMiddleware func(gin.HandlerFunc) gin.HandlerFunc

// ===== [ Implementations ] =====
// ===== [ Private Functions ] =====
// ===== [ Public Functions ] =====

// IPTokenExtractor - Request에서 IP 정보 추출
func IPTokenExtractor(c *gin.Context) string {
	//return strings.Split(c.ClientIP(), ":")[0]
	ip := strings.Split(c.ClientIP(), ":")[0]
	if ip == "" {
		cIP, err := core.GetClientIPHelper(c.Request)
		if err != nil {
			logger.Error(err)
		}
		return cIP
	}
	return ip
}

// HeaderTokenExtractor - Request에서 Header 정보 추출
func HeaderTokenExtractor(header string) TokenExtractor {
	return func(c *gin.Context) string { return c.Request.Header.Get(header) }
}

// NewTokenLimiter - 지정된 Token Extractor와 Limiter Store를 기반으로 Token Bucket 기준 Rate limit 처리를 수행하는 Middleware 생성
func NewTokenLimiter(te TokenExtractor, ls ratelimit.LimiterStore) RateLimitMiddleware {
	return func(next gin.HandlerFunc) gin.HandlerFunc {
		return func(c *gin.Context) {
			tokenKey := te(c)
			// TokenBucket 검증
			if tokenKey == "" || !ls(tokenKey).Allow() {
				c.AbortWithError(http.StatusTooManyRequests, ratelimit.ErrClientLimited)
				return
			}
			next(c)
		}
	}
}

// NewIPLimiter - Rate Limit 기능을 IP 기준으로 처리하는 Handler Func 생성 (By Client IP 단위 검증)
func NewIPLimiter(maxRate float64, capacity int64) RateLimitMiddleware {
	return NewTokenLimiter(IPTokenExtractor, ratelimit.NewMemoryStore(maxRate, capacity))
}

// NewHeaderLimiter - Rate Limit 기능을 Header 기준으로 처리하는 Handler Func 생성 (By Client Header 단위 검증)
func NewHeaderLimiter(header string, maxRate float64, capacity int64) RateLimitMiddleware {
	return NewTokenLimiter(HeaderTokenExtractor(header), ratelimit.NewMemoryStore(maxRate, capacity))
}

// NewEndpointRateLimiter - Rate Limit 처리를 포함하는 Handler Func 생성 (By Call 단위 검증)
func NewEndpointRateLimiter(tb ratelimit.RateLimiter) RateLimitMiddleware {
	return func(next gin.HandlerFunc) gin.HandlerFunc {
		return func(c *gin.Context) {
			// TokenBucket 검증
			if !tb.Allow() {
				c.AbortWithError(503, ratelimit.ErrRouteLimited)
				return
			}
			next(c)
		}
	}
}

// HandlerFactory - RateLimit 기능을 수행하는 Handler Factory 구성
func HandlerFactory(next ginRouter.HandlerFactory, logger logging.Logger) ginRouter.HandlerFactory {
	return func(eConf *config.EndpointConfig, p proxy.Proxy) gin.HandlerFunc {
		handlerFunc := next(eConf, p)

		conf := router.ParseConfig(eConf)
		if conf != nil {
			if conf.MaxRate <= 0 && conf.ClientMaxRate <= 0 {
				//TODO: Waring log
				return handlerFunc
			}

			if conf.MaxRate > 0 {
				handlerFunc = NewEndpointRateLimiter(ratelimit.NewLimiterWithRate(float64(conf.MaxRate), conf.MaxRate))(handlerFunc)
			}
			if conf.ClientMaxRate > 0 {
				switch strings.ToLower(conf.Strategy) {
				case "ip":
					handlerFunc = NewIPLimiter(float64(conf.ClientMaxRate), conf.ClientMaxRate)(handlerFunc)
				case "header":
					handlerFunc = NewHeaderLimiter(conf.Key, float64(conf.ClientMaxRate), conf.ClientMaxRate)(handlerFunc)
				}
			}
		}

		return handlerFunc
	}
}
