package gin

import (
	"github.com/cloud-barista/cb-apigw/restapigw/pkg/config"
	corsMw "github.com/cloud-barista/cb-apigw/restapigw/pkg/middlewares/cors"
	"github.com/gin-gonic/gin"
	"github.com/rs/cors"
	ginCors "github.com/rs/cors/wrapper/gin"
)

// ===== [ Constants and Variables ] =====

// ===== [ Types ] =====

// ===== [ Implementations ] =====

// ===== [ Private Functions ] =====

// ===== [ Public Functions ] =====

// New - CORS 구성을 반영한 GIN HandlerFunc 반환
func New(mConf config.MWConfig, engine *gin.Engine) {
	conf := corsMw.ParseConfig(mConf)
	if nil == conf {
		return
	}

	engine.Use(ginCors.New(cors.Options{
		AllowedOrigins:   conf.AllowOrigins,
		AllowedMethods:   conf.AllowMethods,
		AllowedHeaders:   conf.AllowHeaders,
		ExposedHeaders:   conf.ExposeHeaders,
		AllowCredentials: conf.AllowCredentials,
		MaxAge:           int(conf.MaxAge.Seconds()),
	}))
}
