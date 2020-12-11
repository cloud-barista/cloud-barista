// Package admin -
package admin

import (
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"

	"github.com/cloud-barista/cb-apigw/restapigw/pkg/logging"
	"github.com/cloud-barista/cb-apigw/restapigw/pkg/static"
	"github.com/gin-gonic/gin"
)

// ===== [ Constants and Variables ] =====
// ===== [ Types ] =====
// ===== [ Implementations ] =====
// ===== [ Private Functions ] =====

// isWebPath - 전달된 경로가 Admin Web Path에 해당하는지 검증
func isWebPath(path string) bool {
	paths := []string{"/image", "/_nuxt", "/scripts"}
	for _, s := range paths {
		if strings.HasPrefix(path, s) {
			return true
		}
	}
	return false
}

// ===== [ Public Functions ] =====

// WebServe - Admin Web Serve
func WebServe(urlPrefix string) gin.HandlerFunc {
	staticHandler := http.StripPrefix("/", http.FileServer(static.Dir(false, "/web/dist")))

	return func(c *gin.Context) {
		_path := c.Request.URL.Path
		if _path == "/" {
			// 초기 페이지 반환
			vueApp, err := static.FSString(false, "/web/dist/index.html")
			if nil != err {
				c.AbortWithError(-1, err)
			}
			//c.Writer.WriteString(vueApp)
			c.Data(http.StatusOK, "text/html", []byte(vueApp))
			c.Abort()
		} else if isWebPath(_path) {
			staticHandler.ServeHTTP(c.Writer, c.Request)
			c.Abort()
		}
	}
}

// RedirectHTTPS - HTTP로 전달된 Request를 HTTPS로 전달
func RedirectHTTPS(port int) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		host, _, _ := net.SplitHostPort(req.Host)

		target := url.URL{
			Scheme: "https",
			Host:   fmt.Sprintf("%s:%v", host, port),
			Path:   req.URL.Path,
		}
		if 0 < len(req.URL.RawQuery) {
			target.RawQuery += "?" + req.URL.RawQuery
		}
		logging.GetLogger().Infof("[API SERVER] Redirect to: %s", target.String())
		http.Redirect(w, req, target.String(), http.StatusTemporaryRedirect)
	})
}
