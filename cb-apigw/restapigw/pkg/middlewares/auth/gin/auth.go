// Package gin - API Gateway 접근인증을 위한 GIN 과 연동되는 기능 제공
package gin

import (
	"bytes"
	"net/http"
	"strings"

	"github.com/cloud-barista/cb-apigw/restapigw/pkg/config"
	"github.com/cloud-barista/cb-apigw/restapigw/pkg/errors"
	"github.com/cloud-barista/cb-apigw/restapigw/pkg/logging"
	"github.com/cloud-barista/cb-apigw/restapigw/pkg/proxy"
	"github.com/cloud-barista/cb-apigw/restapigw/pkg/router"
	ginRouter "github.com/cloud-barista/cb-apigw/restapigw/pkg/router/gin"
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
)

// ===== [ Constants and Variables ] =====

var (
	// MWNamespace - Middleware 설정 식별자
	MWNamespace = "mw-auth"
)

// ===== [ Types ] =====

// Config - AUTH 운영을 위한 설정 구조
type Config struct {
	// HMAC 구성을 위한 보안 키
	SecureKey string `yaml:"secure_key"`
	// 접속을 허용하는 ID 리스트
	AcessIds []string `yaml:"access_ids"`
}

// ===== [ Implementations ] =====

// ===== [ Private Functions ] =====

// ParseConfig - HTTPCache 운영을 위한 Configuration parsing 처리
func ParseConfig(mwConf config.MWConfig) *Config {
	conf := new(Config)
	tmp, ok := mwConf[MWNamespace]
	if !ok {
		return nil
	}

	buf := new(bytes.Buffer)
	yaml.NewEncoder(buf).Encode(tmp)
	if err := yaml.NewDecoder(buf).Decode(conf); nil != err {
		return nil
	}

	return conf
}

// validateToken - Header로 전달된 Auth Token 검증
func validateToken(conf *Config, req *http.Request) error {
	token := req.Header.Get("Authorization")
	if "" == token {
		return errors.New("Authorization token not found")
	}

	token = strings.Replace(token, "Bearer ", "", 1)

	return ValidateToken(conf.SecureKey, token, conf.AcessIds)
}

// ===== [ Public Functions ] =====

// TokenValidator - Auth Token Validation 처리를 수행하는 Route Handler 구성
func TokenValidator(hf ginRouter.HandlerFactory, logger logging.Logger) ginRouter.HandlerFactory {
	return func(eConf *config.EndpointConfig, proxy proxy.Proxy) gin.HandlerFunc {
		handler := hf(eConf, proxy)

		conf := ParseConfig(eConf.Middleware)
		if nil == conf {
			return handler
		}

		return func(c *gin.Context) {
			if err := validateToken(conf, c.Request); nil != err {
				c.Header(router.CompleteResponseHeaderName, "false")
				c.Header(router.MessageResponseHeaderName, err.Error())
				c.AbortWithError(http.StatusUnauthorized, err)
				return
			}

			handler(c)
		}
	}
}
