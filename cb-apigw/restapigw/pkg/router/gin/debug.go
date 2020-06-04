package gin

import (
	"io/ioutil"

	"github.com/cloud-barista/cb-apigw/restapigw/pkg/logging"
	"github.com/gin-gonic/gin"
)

// ===== [ Constants and Variables ] =====

// ===== [ Types ] =====

// ===== [ Implementations ] =====

// ===== [ Private Functions ] =====

// ===== [ Public Functions ] =====

// DebugHandler - debug 테스트를 위한 log 출력용 Handler 구성
func DebugHandler(logger logging.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		body, _ := ioutil.ReadAll(c.Request.Body)
		c.Request.Body.Close()
		logger.Debugf("Method: %v, URL: %v, Query: %v, Params: %v, Headers: %v, Body: %v", c.Request.Method, c.Request.RequestURI, c.Request.URL.Query(), c.Params, c.Request.Header, string(body))
		c.JSON(200, gin.H{"message": "pong"})
	}
}
