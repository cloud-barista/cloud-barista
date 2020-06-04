package gin

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/unrolled/secure"
	"github.com/cloud-barista/cb-apigw/restapigw/pkg/config"
	"github.com/cloud-barista/cb-apigw/restapigw/pkg/middlewares/httpsecure"
)

// ===== [ Constants and Variables ] =====

var (
	errNoConfig = errors.New("no config present for the httpsecure middleware")
)

// ===== [ Types ] =====

// ===== [ Implementations ] =====

// ===== [ Private Functions ] =====

// secureHandler - Gin Engine에 적용할 httpsecure middleware handler 설정
func secureHandler(opts secure.Options) gin.HandlerFunc {
	secMw := secure.New(opts)
	return func(c *gin.Context) {
		err := secMw.Process(c.Writer, c.Request)
		if err != nil {
			c.Abort()
			return
		}
		// Redirect 관련 상태 코드인 경우
		if status := c.Writer.Status(); status > 300 && status < 399 {
			c.Abort()
		}
	}
}

// ===== [ Public Functions ] =====

// Register - Gin Engine 에 httpsecure middleware 설정을 등록
func Register(mwConf config.MWConfig, engine *gin.Engine) error {
	conf := httpsecure.ParseConfig(mwConf)
	if conf == nil {
		return errNoConfig
	}
	engine.Use(secureHandler(secure.Options{
		AllowedHosts:            conf.AllowedHosts,
		HostsProxyHeaders:       conf.HostsProxyHeaders,
		STSSeconds:              conf.STSSeconds,
		CustomFrameOptionsValue: conf.CustomFrameOptionsValue,
		ContentSecurityPolicy:   conf.ContentSecurityPolicy,
		PublicKey:               conf.PublicKey,
		SSLHost:                 conf.SSLHost,
		ReferrerPolicy:          conf.ReferrerPolicy,
		ContentTypeNosniff:      conf.ContentTypeNosniff,
		BrowserXssFilter:        conf.BrowserXSSFilter,
		IsDevelopment:           conf.IsDevelopment,
		STSIncludeSubdomains:    conf.STSIncludeSubdomains,
		FrameDeny:               conf.FrameDeny,
		SSLRedirect:             conf.SSLRedirect,
	}))
	return nil
}
