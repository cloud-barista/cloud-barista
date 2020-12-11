package metrics

import (
	"crypto/tls"

	gometrics "github.com/rcrowley/go-metrics"
)

// ===== [ Constants and Variables ] =====

// ===== [ Types ] =====

// RouterMetrics - Router에 대한 Metrics Collector 구조 정의
type RouterMetrics struct {
	ProxyMetrics
	connected         gometrics.Counter
	disconnected      gometrics.Counter
	connectedTotal    gometrics.Counter
	disconnectedTotal gometrics.Counter
	connectedGauge    gometrics.Gauge
	disconnectedGauge gometrics.Gauge
}

// ===== [ Implementations ] =====

// Connection - Router에 연결이 발생한 경우에 Connection counter 증가 처리
func (rm *RouterMetrics) Connection(TLS *tls.ConnectionState) {
	rm.connected.Inc(1)
	if nil == TLS {
		return
	}

	rm.Counter("tls_version", tlsVersion[TLS.Version], "count").Inc(1)
	rm.Counter("tls_cipher", tlsCipherSuite[TLS.CipherSuite], "count").Inc(1)
}

// Disconnection - Router 연결이 종료된 경우에 Disconnection counter 증가 처리
func (rm *RouterMetrics) Disconnection() {
	rm.disconnected.Inc(1)
}

// Aggregate - Router Metrics에 정보 취합 처리
func (rm *RouterMetrics) Aggregate() {
	conCount := rm.connected.Count()
	rm.connectedGauge.Update(conCount)
	rm.connectedTotal.Inc(conCount)
	rm.connected.Clear()

	disconCount := rm.disconnected.Count()
	rm.disconnectedGauge.Update(disconCount)
	rm.disconnectedTotal.Inc(disconCount)
	rm.disconnected.Clear()
}

// RegisterResponseWriterMetrics - gin.ResponseWriter에 연동되는 Metric 설정
func (rm *RouterMetrics) RegisterResponseWriterMetrics(name string) {
	rm.Counter("response", name, "status")
	rm.Histogram("response", name, "size")
	rm.Histogram("response", name, "time")
}

// ===== [ Private Functions ] =====

// ===== [ Public Functions ] =====

// NewRouterMetrics - 지정된 Registry를 Parent로 사용하는 Router metrics 생성
func NewRouterMetrics(parentRegistry *gometrics.Registry) *RouterMetrics {
	r := gometrics.NewPrefixedChildRegistry(*parentRegistry, "router.")

	// Registry는 ProxyMetrics 와 공유
	return &RouterMetrics{
		ProxyMetrics:      ProxyMetrics{r},
		connected:         gometrics.NewRegisteredCounter("connected", r),
		disconnected:      gometrics.NewRegisteredCounter("disconnected", r),
		connectedTotal:    gometrics.NewRegisteredCounter("connected-total", r),
		disconnectedTotal: gometrics.NewRegisteredCounter("disconnected-total", r),
		connectedGauge:    gometrics.NewRegisteredGauge("connected-gauge", r),
		disconnectedGauge: gometrics.NewRegisteredGauge("disconnected-gauge", r),
	}
}
