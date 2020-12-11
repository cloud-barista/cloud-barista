package gin

import (
	"context"
	"fmt"
	"net/textproto"
	"strings"

	"github.com/cloud-barista/cb-apigw/restapigw/pkg/config"
	"github.com/cloud-barista/cb-apigw/restapigw/pkg/core"
	"github.com/cloud-barista/cb-apigw/restapigw/pkg/logging"
	"github.com/cloud-barista/cb-apigw/restapigw/pkg/proxy"
	"github.com/cloud-barista/cb-apigw/restapigw/pkg/router"
	"github.com/gin-gonic/gin"
)

// ===== [ Constants and Variables ] =====

const (
	requestParamsAsterisk = "*"
)

var (
	logger = logging.NewLogger()
)

// ===== [ Types ] =====

// responseError - Response에서 발생한 오류와 상태코드를 관리하는 오류
type responseError interface {
	error
	StatusCode() int
}

// HandlerFactory - 지정된 Endpoint 설정과 Proxy 를 기반으로 동작하는 Gin Framework handler 팩토리 정의
type HandlerFactory func(*config.EndpointConfig, proxy.Proxy) gin.HandlerFunc

// ===== [ Implementations ] =====

// ===== [ Private Functions ] =====

// ===== [ Public Functions ] =====

// EndpointHandler - 지정된 Endpoint 설정과 Proxy를 연계 호출하는 Gin Framework handler 생성
func EndpointHandler(eConf *config.EndpointConfig, proxy proxy.Proxy) gin.HandlerFunc {
	return CustomErrorEndpointHandler(eConf, proxy, router.DefaultToHTTPError)
}

// NewRequest - 제외할 Header 정보와 Query string 정보를 반영하는 Gin Context 기반의 Request 생성
func NewRequest(eConf *config.EndpointConfig) func(*gin.Context, []string) *proxy.Request {
	return func(c *gin.Context, exceptQueryStrings []string) *proxy.Request {
		params := make(map[string]string, len(c.Params))
		for _, param := range c.Params {
			params[strings.Title(param.Key)] = param.Value
		}

		// Header Check - All pass with blacklist
		headers := make(map[string][]string)
		for k, v := range c.Request.Header {
			except := false
			for i := range eConf.ExceptHeaders {
				key := textproto.CanonicalMIMEHeaderKey(eConf.ExceptHeaders[i])
				if k == key {
					except = true
					break
				}
			}
			if !except && !core.ContainsString(router.HeadersToNotSend, k) {
				tmp := make([]string, len(v))
				copy(tmp, v)
				headers[k] = tmp
			}
		}

		// 필수 Header 정보 설정
		for _, k := range router.HeadersToSend {
			key := textproto.CanonicalMIMEHeaderKey(k)
			// 현재 없는 경우
			if _, ok := headers[key]; !ok {
				// Requst Header에 있는 경우
				if h, ok := c.Request.Header[key]; ok {
					headers[k] = h
				}
			}
		}

		headers["X-Forwarded-For"] = []string{c.ClientIP()}
		// 전달되는 Header에 "User-Agent" 가 존재하지 않는 경우는 Router의 User-Agent 사용
		if _, ok := headers["User-Agent"]; !ok {
			headers["User-Agent"] = router.UserAgentHeaderValue
		} else {
			headers["X-Forwarded-Via"] = router.UserAgentHeaderValue
		}

		// QueryString Check - All pass with blacklist
		query := c.Request.URL.Query()

		// Black list 적용
		if nil != exceptQueryStrings {
			for i := range exceptQueryStrings {
				delete(query, exceptQueryStrings[i])
			}
		}

		// Bypass인 경우 실제 호출 경로 설정
		path := ""
		if eConf.IsBypass {
			path = c.Request.URL.Path
		}

		return &proxy.Request{
			IsBypass: eConf.IsBypass,
			Method:   c.Request.Method,
			Query:    query,
			Body:     c.Request.Body,
			Path:     path,
			Params:   params,
			Headers:  headers,
		}
	}
}

// CustomErrorEndpointHandler - 지정한 Endpoint 설정과 수행할 Proxy 정보 및 오류 처리를 위한 Gin Framework handler 생성
func CustomErrorEndpointHandler(eConf *config.EndpointConfig, proxy proxy.Proxy, errF router.ToHTTPError) gin.HandlerFunc {
	cacheControlHeaderValue := fmt.Sprintf("public, max-age=%d", int(eConf.CacheTTL.Seconds()))
	isCacheEnabled := 0 != eConf.CacheTTL.Seconds()
	requestGenerator := NewRequest(eConf)
	render := getRender(eConf)

	return func(c *gin.Context) {
		requestCtx, cancel := context.WithTimeout(c, eConf.Timeout)
		c.Header(core.AppHeaderName, fmt.Sprintf("Version %s", core.AppVersion))
		response, err := proxy(requestCtx, requestGenerator(c, eConf.ExceptQueryStrings))
		select {
		case <-requestCtx.Done():
			if nil == err {
				err = router.ErrInternalError
			}
		default:
		}

		complete := router.HeaderIncompleteResponseValue
		if nil != response && 0 < len(response.Data) {
			if response.IsComplete {
				complete = router.HeaderCompleteResponseValue
				if isCacheEnabled {
					c.Header("Cache-Control", cacheControlHeaderValue)
				}
			}

			for k, vs := range response.Metadata.Headers {
				for _, v := range vs {
					c.Writer.Header().Add(k, v)
				}
			}
		}
		c.Header(router.CompleteResponseHeaderName, complete)

		if nil != err {
			// Proxy 처리 중에 발생한 오류들을 Header로 설정
			c.Header(router.MessageResponseHeaderName, err.Error())
			logger.Errorf("[API G/W] Router > Endpoint Error Processing: %s", err.Error())
			c.Error(err)

			// Response가 없는 경우의 상태 코드 설정
			if nil == response {
				if t, ok := err.(responseError); ok {
					c.Status(t.StatusCode())
				} else if e, ok := err.(core.WrappedError); ok {
					c.Status(e.Code())
				} else {
					c.Status(errF(err))
				}
				cancel()
				return
			}
		} else {
			c.Header(router.MessageResponseHeaderName, "OK")
		}

		render(c, response)
		cancel()
	}
}
