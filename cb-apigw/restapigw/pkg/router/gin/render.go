package gin

import (
	"io"
	"net/http"
	"sync"

	"github.com/cloud-barista/cb-apigw/restapigw/pkg/config"
	"github.com/cloud-barista/cb-apigw/restapigw/pkg/core"
	"github.com/cloud-barista/cb-apigw/restapigw/pkg/encoding"
	"github.com/cloud-barista/cb-apigw/restapigw/pkg/proxy"
	"github.com/gin-gonic/gin"
)

// ===== [ Constants and Variables ] =====

const (
	// NEGOTIATE - Backend의 Response 포맷과 Endpoint의 "OutputEncoding" 협상을 위한 식별자
	NEGOTIATE = "nogotiate"
)

var (
	mutex          = &sync.RWMutex{}
	emptyResponse  = gin.H{}
	renderRegister = map[string]Render{
		NEGOTIATE:       negotiatedRender,
		encoding.STRING: stringRender,
		encoding.JSON:   jsonRender,
		encoding.NOOP:   noopRender,
	}
)

// ===== [ Types ] =====

// Render - Proxy 수행의 결과인 Response에 대한 Encoding과 Rendering 함수 정의
type Render func(*gin.Context, *proxy.Response)

// ===== [ Implementations ] =====

// ===== [ Private Functions ] =====

// negotiatedRender - ginContext 기반에서 Response 정보를 기준으로 Render를 결정하고 처리
func negotiatedRender(c *gin.Context, response *proxy.Response) {
	switch c.NegotiateFormat(gin.MIMEJSON, gin.MIMEPlain, gin.MIMEXML) {
	case gin.MIMEXML:
		xmlRender(c, response)
	case gin.MIMEPlain:
		yamlRender(c, response)
	default:
		jsonRender(c, response)
	}
}

// getWithFallback - 지정한 이름에 해당하는 Render를 반환, 만일 존재하지 않는 경우는 지정한 Fallback Render 반환
func getWithFallback(key string, fallback Render) Render {
	mutex.RLock()
	r, ok := renderRegister[key]
	mutex.RUnlock()
	if !ok {
		return fallback
	}
	return r
}

// jsonRender - JSON 포맷에 대한 Render 처리
func jsonRender(c *gin.Context, res *proxy.Response) {
	status := c.Writer.Status()
	if nil == res {
		c.JSON(status, emptyResponse)
		return
	}

	// WrappingTag가 설정된 경우는 Array인 상태로 반환한다.
	var data interface{}
	if v, ok := res.Data[core.WrappingTag]; ok {
		delete(res.Data, core.WrappingTag)
		data = res.Data[v.(string)]
	} else {
		data = res.Data
	}

	c.JSON(status, data)
}

// stringRender - 단순 문자열에 대한 Render 처리
func stringRender(c *gin.Context, res *proxy.Response) {
	status := c.Writer.Status()

	if nil == res {
		c.String(status, "")
		return
	}
	d, ok := res.Data["content"]
	if !ok {
		c.String(status, "")
		return
	}
	msg, ok := d.(string)
	if !ok {
		c.String(status, "")
		return
	}
	c.String(status, msg)
}

// xmlRender - XML 포맷에 대한 Render 처리
func xmlRender(c *gin.Context, res *proxy.Response) {
	status := c.Writer.Status()
	if nil == res {
		c.XML(status, nil)
		return
	}
	d, ok := res.Data["content"]
	if !ok {
		c.XML(status, nil)
		return
	}
	c.XML(status, d)
}

// yamlRender - YAML 포맷에 대한 Render 처리
func yamlRender(c *gin.Context, res *proxy.Response) {
	status := c.Writer.Status()
	if nil == res {
		c.YAML(status, emptyResponse)
		return
	}
	c.YAML(status, res.Data)
}

// noopRender - 아무 변환도 없는 Render 처리
func noopRender(c *gin.Context, res *proxy.Response) {
	if nil == res {
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Status(res.Metadata.StatusCode)
	for k, vs := range res.Metadata.Headers {
		for _, v := range vs {
			c.Writer.Header().Add(k, v)
		}
	}
	io.Copy(c.Writer, res.Io)
}

// getRender - Endpoint 설정에 지정된 Backend의 "encoding" 을 기준 Encoding(fallback)을 설정하고 Endpoint 설정의 "output_encoding" 기준으로 운영되는 Render 반환
func getRender(eConf *config.EndpointConfig) Render {
	fallback := jsonRender
	if 1 == len(eConf.Backend) {
		fallback = getWithFallback(eConf.Backend[0].Encoding, fallback)
	}

	if "" == eConf.OutputEncoding {
		return fallback
	}

	return getWithFallback(eConf.OutputEncoding, fallback)
}

// ===== [ Public Functions ] =====
