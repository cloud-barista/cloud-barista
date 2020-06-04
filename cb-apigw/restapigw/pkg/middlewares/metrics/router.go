package metrics

import (
	"strings"

	gometrics "github.com/rcrowley/go-metrics"
)

// ===== [ Constants and Variables ] =====

// ===== [ Types ] =====

// RouterMetrics - Router에 대한 Metrics Collector 구조 정의
type RouterMetrics struct {
	register          gometrics.Registry
	connected         gometrics.Counter
	disconnected      gometrics.Counter
	connectedTotal    gometrics.Counter
	disconnectedTotal gometrics.Counter
	connectedGauge    gometrics.Gauge
	disconnectedGauge gometrics.Gauge
}

// ===== [ Implementations ] =====

// Connection - Router에 연결이 발생한 경우에 Connection counter 증가 처리
func (rm *RouterMetrics) Connection() {
	rm.connected.Inc(1)
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

// Counter - Metric Counter가 없는 경우는 등록하고 대상 Counter 반환
func (rm *RouterMetrics) Counter(labels ...string) gometrics.Counter {
	return gometrics.GetOrRegisterCounter(strings.Join(labels, "."), rm.register)
}

// Histogram - Metric Histogram이 없는 경우는 등록하고 대상 Histogram 반환
func (rm *RouterMetrics) Histogram(labels ...string) gometrics.Histogram {
	return gometrics.GetOrRegisterHistogram(strings.Join(labels, "."), rm.register, defaultSample())
}

// ===== [ Private Functions ] =====

// ===== [ Public Functions ] =====

// NewRouterMetrics - 지정된 Registry를 Parent로 사용하는 Router metrics 생성
func NewRouterMetrics(parent *gometrics.Registry) *RouterMetrics {
	r := gometrics.NewPrefixedChildRegistry(*parent, "router.")

	return &RouterMetrics{
		register:          r,
		connected:         gometrics.NewRegisteredCounter("connected", r),
		disconnected:      gometrics.NewRegisteredCounter("disconnected", r),
		connectedTotal:    gometrics.NewRegisteredCounter("connected-total", r),
		disconnectedTotal: gometrics.NewRegisteredCounter("disconnected-total", r),
		connectedGauge:    gometrics.NewRegisteredGauge("connected-gauge", r),
		disconnectedGauge: gometrics.NewRegisteredGauge("disconnected-gauge", r),
	}
}
