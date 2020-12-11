package gin

import (
	"context"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/stats"
	"go.opencensus.io/tag"
)

// ===== [ Constants and Variables ] =====

// ===== [ Types ] =====

// trackingResponseWriter - Gin Framework response writer와 수집된 Trace 정보를 관리하는 구조
type trackingResponseWriter struct {
	gin.ResponseWriter
	ctx     context.Context
	reqSize int64
	start   time.Time
	endOnce sync.Once
}

// ===== [ Implementations ] =====

// end - Request 호출 결과를 Trace 정보로 처리
func (t *trackingResponseWriter) end() {
	t.endOnce.Do(func() {
		m := []stats.Measurement{
			ochttp.ServerLatency.M(float64(time.Since(t.start)) / float64(time.Millisecond)),
			ochttp.ServerResponseBytes.M(int64(t.Size())),
		}
		if 0 <= t.reqSize {
			m = append(m, ochttp.ServerRequestBytes.M(t.reqSize))
		}
		status := t.Status()
		if 0 == status {
			status = http.StatusOK
		}
		ctx, _ := tag.New(t.ctx, tag.Upsert(ochttp.StatusCode, strconv.Itoa(status)))
		stats.Record(ctx, m...)
	})
}

// ===== [ Private Functions ] =====

// ===== [ Public Functions ] =====
