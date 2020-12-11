package gin

import (
	"net/http"
	"time"

	"github.com/cloud-barista/cb-apigw/restapigw/pkg/config"
	"github.com/cloud-barista/cb-apigw/restapigw/pkg/logging"
	"github.com/cloud-barista/cb-apigw/restapigw/pkg/middlewares/opencensus"
	"github.com/cloud-barista/cb-apigw/restapigw/pkg/proxy"
	ginRouter "github.com/cloud-barista/cb-apigw/restapigw/pkg/router/gin"
	"github.com/gin-gonic/gin"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/plugin/ochttp/propagation/b3"
	"go.opencensus.io/stats"
	"go.opencensus.io/tag"
	"go.opencensus.io/trace"
	"go.opencensus.io/trace/propagation"
)

// ===== [ Constants and Variables ] =====

// ===== [ Types ] =====

// handler - Gin Framework 에서 사용할 Trace 정보 처리 구조
type handler struct {
	name             string
	propagation      propagation.HTTPFormat
	Handler          gin.HandlerFunc
	StartOptions     trace.StartOptions
	IsPublicEndpoint bool
	IsByPassEndpoint bool
}

// ===== [ Implementations ] =====

// HandlerFunc - 지정된 Gin Context 기반으로 Request 호출 전후로 Trace 와 Stats 정보 처리 (Start/Stop)
func (h *handler) HandlerFunc(c *gin.Context) {
	var traceEnd, statsEnd func()
	c.Request, traceEnd = h.startTrace(c.Writer, c.Request)
	c.Writer, statsEnd = h.startStats(c.Writer, c.Request)

	c.Set(opencensus.ContextKey, trace.FromContext(c.Request.Context()))
	h.Handler(c)

	statsEnd()
	traceEnd()
}

// startTrace - 지정된 Gin ResponseWriter와 Http Request 정보를 Trace 정보와 연동될 수 있도록 처리
func (h *handler) startTrace(_ gin.ResponseWriter, r *http.Request) (*http.Request, func()) {
	// Bypass인 경우는 Handler Name을 실제 호출 경로로 변경한다.
	if h.IsByPassEndpoint {
		h.name = "[endpoint] " + r.URL.Path
	}

	ctx := r.Context()
	var span *trace.Span
	sc, ok := h.extractSpanContext(r)
	if ok && !h.IsPublicEndpoint {
		ctx, span = trace.StartSpanWithRemoteParent(ctx, h.name, sc, trace.WithSampler(h.StartOptions.Sampler), trace.WithSpanKind(h.StartOptions.SpanKind))
	} else {
		ctx, span = trace.StartSpan(ctx, h.name, trace.WithSampler(h.StartOptions.Sampler), trace.WithSpanKind(h.StartOptions.SpanKind))
		if ok {
			span.AddLink(trace.Link{
				TraceID:    sc.TraceID,
				SpanID:     sc.SpanID,
				Type:       trace.LinkTypeChild,
				Attributes: nil,
			})
		}
	}

	span.AddAttributes(requestAttrs(r)...)
	return r.WithContext(ctx), span.End
}

// extractSpanContext - 지정된 http Request 정보에서 Trace를 위한 Span Context 반환
func (h *handler) extractSpanContext(r *http.Request) (trace.SpanContext, bool) {
	return h.propagation.SpanContextFromRequest(r)
}

// startStats - 지정된 Gin ResponseWriter와 Http Request 정보를 Trace Stats 정보와 연동될 수 있도록 처리
func (h *handler) startStats(w gin.ResponseWriter, r *http.Request) (gin.ResponseWriter, func()) {
	ctx, _ := tag.New(r.Context(),
		tag.Upsert(ochttp.Host, r.URL.Host),
		tag.Upsert(ochttp.Path, r.URL.Path),
		tag.Upsert(ochttp.Method, r.Method))
	track := &trackingResponseWriter{
		start:          time.Now(),
		ctx:            ctx,
		ResponseWriter: w,
	}
	if nil == r.Body {
		track.reqSize = -1
	} else if 0 < r.ContentLength {
		track.reqSize = r.ContentLength
	}
	stats.Record(ctx, ochttp.ServerRequestCount.M(1))
	return track, track.end
}

// ===== [ Private Functions ] =====

// requestAttrs - Trace에 Request 관련 정보를 속성으로 추가
func requestAttrs(r *http.Request) []trace.Attribute {
	return []trace.Attribute{
		trace.StringAttribute(ochttp.PathAttribute, r.URL.Path),
		trace.StringAttribute(ochttp.HostAttribute, r.URL.Host),
		trace.StringAttribute(ochttp.MethodAttribute, r.Method),
		trace.StringAttribute(ochttp.UserAgentAttribute, r.UserAgent()),
	}
}

// responseAttrs - Trace에 Responst 관련 정보를 속성으로 추가
func responseAttrs(resp *http.Response) []trace.Attribute {
	return []trace.Attribute{
		trace.Int64Attribute(ochttp.StatusCodeAttribute, int64(resp.StatusCode)),
	}
}

// ===== [ Public Functions ] =====

// HandlerFunc - 지정된 Endpoint 설정과 Gin Framework handler를 기준으로 Trace 구성을 처리하는 Handler 구성 반환
func HandlerFunc(eConf *config.EndpointConfig, next gin.HandlerFunc, hf propagation.HTTPFormat) gin.HandlerFunc {
	if !opencensus.IsRouterEnabled() {
		return next
	}
	if nil == hf {
		hf = &b3.HTTPFormat{}
	}
	h := &handler{
		name:        "[endpoint] " + eConf.Endpoint,
		propagation: hf,
		Handler:     next,
		StartOptions: trace.StartOptions{
			SpanKind: trace.SpanKindServer,
		},
		IsByPassEndpoint: eConf.IsBypass,
	}
	return h.HandlerFunc
}

// HandlerFactory - 전달된 HandlerFactory를 실행하기 전에 Opencesus Trace 작업을 처리하는 HandlerFactory 구성
func HandlerFactory(hf ginRouter.HandlerFactory, logger logging.Logger) ginRouter.HandlerFactory {
	return func(eConf *config.EndpointConfig, p proxy.Proxy) gin.HandlerFunc {
		return HandlerFunc(eConf, hf(eConf, p), nil)
	}
}
