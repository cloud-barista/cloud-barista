package metrics

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/cloud-barista/cb-apigw/restapigw/pkg/config"
	"github.com/cloud-barista/cb-apigw/restapigw/pkg/logging"
	"github.com/cloud-barista/cb-apigw/restapigw/pkg/proxy"
	"github.com/rcrowley/go-metrics"
	gometrics "github.com/rcrowley/go-metrics"
)

// ===== [ Constants and Variables ] =====

var (
	logger = logging.NewLogger()
)

// ===== [ Types ] =====

// ProxyMetrics - Proxy 에 대한 Metrics Collector 구조 정의
type ProxyMetrics struct {
	register gometrics.Registry
}

// ===== [ Implementations ] =====

// NewProxyCallChain - 지정한 레이어에 대한 Proxy 호출 체인 구성
func (m *Metrics) NewProxyCallChain(layer, name string) proxy.CallChain {
	return NewProxyCallChain(layer, name, m.Proxy)
}

// ProxyFactory - Metrics 처리를 수행하는 ProxyFactory 생성
func (m *Metrics) ProxyFactory(segmentName string, next proxy.Factory) proxy.FactoryFunc {
	if m.Config == nil || !m.Config.ProxyEnabled {
		return next.New
	}

	return proxy.FactoryFunc(func(eConf *config.EndpointConfig) (proxy.Proxy, error) {
		next, err := next.New(eConf)
		if err != nil {
			return proxy.DummyProxy, err
		}
		return m.NewProxyCallChain(segmentName, eConf.Endpoint)(next), nil
	})
}

// BackendFactory - Metrics 처리를 수행하는 BackendFactory 생성
func (m *Metrics) BackendFactory(segmentName string, next proxy.BackendFactory) proxy.BackendFactory {
	if m.Config == nil || !m.Config.BackendEnabled {
		return next
	}

	return func(bConf *config.BackendConfig) proxy.Proxy {
		return m.NewProxyCallChain(segmentName, bConf.URLPattern)(next(bConf))
	}
}

// Counter - Metric Counter가 없는 경우는 등록하고 대상 Counter 반환
func (pm *ProxyMetrics) Counter(labels ...string) gometrics.Counter {
	return gometrics.GetOrRegisterCounter(strings.Join(labels, "."), pm.register)
}

// Histogram - Metric Histogram이 없는 경우는 등록하고 대상 Histogram 반환
func (pm *ProxyMetrics) Histogram(labels ...string) gometrics.Histogram {
	return gometrics.GetOrRegisterHistogram(strings.Join(labels, "."), pm.register, defaultSample())
}

// ===== [ Private Functions ] =====

// registerProxyCallChainMetrics - Metrics 처리를 위한 Proxy CallChain 등록
func registerProxyCallChainMetrics(layer, name string, pm *ProxyMetrics) {
	labels := "layer." + layer + ".name." + name
	for _, complete := range []string{"true", "false"} {
		for _, errored := range []string{"true", "false"} {
			pm.Counter("requests." + labels + ".complete." + complete + ".error." + errored)
			pm.Histogram("latency." + labels + ".complete." + complete + ".error." + errored)
		}
	}
}

// ===== [ Public Functions ] =====

// NewProxyMetrics - 지정된 Register를 Parent로 사용하는 Proxy metrics 생성
func NewProxyMetrics(parent *metrics.Registry) *ProxyMetrics {
	m := metrics.NewPrefixedChildRegistry(*parent, "proxy.")
	return &ProxyMetrics{m}
}

// NewProxyCallChain - Metrics 처리를 수행하는 Proxy 호출 체인 생성
func NewProxyCallChain(layer, name string, pm *ProxyMetrics) proxy.CallChain {
	return func(next ...proxy.Proxy) proxy.Proxy {
		if len(next) > 1 {
			panic(proxy.ErrTooManyProxies)
		}

		return func(ctx context.Context, request *proxy.Request) (*proxy.Response, error) {
			// Bypass Backend URLPattern을 실제 URL Path로 변경
			urlPath := name
			if request.IsBypass || layer == "pipe" {
				urlPath = request.Path
			}

			// Metric 처리를 위한 Proxy 호출 정보 등록
			registerProxyCallChainMetrics(layer, urlPath, pm)

			logger.Debugf("[Backend Process Flow] Metrics > Proxy CallChain > %s layer > %s name", layer, name)

			begin := time.Now()
			resp, err := next[0](ctx, request)

			// Metric 처리를 위한 호출 결과 정보 등록
			go func(duration int64, resp *proxy.Response, err error) {
				errored := strconv.FormatBool(err != nil)
				completed := strconv.FormatBool(resp != nil && resp.IsComplete)
				labels := "layer." + layer + ".name." + urlPath + ".complete." + completed + ".error." + errored
				pm.Counter("requests." + labels).Inc(1)
				pm.Histogram("latency." + labels).Update(duration)
			}(time.Since(begin).Nanoseconds(), resp, err)

			return resp, err
		}
	}
}
