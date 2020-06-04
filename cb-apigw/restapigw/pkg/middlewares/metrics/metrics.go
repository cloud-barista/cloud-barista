package metrics

import (
	"bytes"
	"context"
	"time"

	"github.com/cloud-barista/cb-apigw/restapigw/pkg/config"
	"github.com/cloud-barista/cb-apigw/restapigw/pkg/core"
	"github.com/cloud-barista/cb-apigw/restapigw/pkg/logging"
	gometrics "github.com/rcrowley/go-metrics"
	"gopkg.in/yaml.v3"
)

// ===== [ Constants and Variables ] =====

const (
	// MWNamespace - Middleware 식별자
	MWNamespace = "mw-metrics"
)

var (
	// defaultListenAddress - Metrics 수집을 위한 Pulling 용 Port
	defaultListenAddress = ":9000"

	// Metrics 의 수행 구간 비율 측정용
	percentiles = []float64{0.1, 0.25, 0.5, 0.75, 0.9, 0.95, 0.99}
	// 기본 Metrics 샘플
	defaultSample = func() gometrics.Sample { return gometrics.NewUniformSample(1028) }
)

// ===== [ Types ] =====

// DummyRegistry - 설정이 없는 경우에 활용할 Registry 구조 정의
type DummyRegistry struct{}

// Config - Metrics 운영에 필요한 설정 구조 정의
type Config struct {
	// Router Metrics 활성화 여부 (기본값: true)
	RouterEnabled bool `yaml:"router_enabled"`
	// Metrics Expose를 위한 HTTP/Gin Server 구동 비활성화 여부 (기본값: false)
	ExposeMetrics bool `yaml:"expose_metrics"`
	// Proxy Metrics 활성화 여부 (기본값: true)
	ProxyEnabled bool `yaml:"proxy_enabled"`
	// Backend Metrics 활성화 여부 (기본값: true)
	BackendEnabled bool `yaml:"backend_enabled"`
	// Metrics 수집 시간 (기본값: 0, Metrics 수집하지 않음)
	CollectionPeriod time.Duration `yaml:"collection_period"`
	// Metrics Expose 처리를 위한 URL (기본값: )
	ListenAddress string `yaml:"listen_address"`
	InfluxDB      struct {
		// InfluxDB Address
		Address string `yaml:"address"`
		// 수집된 Metrics를 InfluxDB로 전송하는 주기 (기본값: 0, 전송하지 않음)
		ReportingPeriod time.Duration `yaml:"reporting_period"`
		// 데이터 전송에 사용할 Buffer 크기 (기본값: 0, 크기 제한 없음)
		BufferSize int `yaml:"buffer_size"`
		// 접속에 사용할 사용자 명
		UserName string `yaml:"user_name"`
		// 접속에 사용할 비밀번호
		Password string `yaml:"password"`
		// 접속할 데이터베이스 명
		Database string `yaml:"database"`
	} `yaml:"influxdb"`
}

// Metrics - 전체 Metrics 관리를 위한 구조 정의
type Metrics struct {
	logger         *logging.Logger
	Config         *Config
	Proxy          *ProxyMetrics
	Router         *RouterMetrics
	Registry       *gometrics.Registry
	latestSnapshot Stats
}

// ===== [ Implementations ] =====

// Snapshot - 최종 처리된 Snapshot 반환 (InfluxDB 쪽에서 호출)
func (m *Metrics) Snapshot() Stats {
	return m.latestSnapshot
}

// TakeSnapshot - 현재 상태에 대한 Snapshot 처리
func (m *Metrics) TakeSnapshot() Stats {
	stats := NewStats()

	(*m.Registry).Each(func(key string, val interface{}) {
		switch metric := val.(type) {
		case gometrics.Counter:
			stats.Counters[key] = metric.Count()
		case gometrics.Gauge:
			stats.Gauges[key] = metric.Value()
		case gometrics.Histogram:
			stats.Histograms[key] = HistogramData{
				Max:         metric.Max(),
				Min:         metric.Min(),
				Mean:        metric.Mean(),
				Stddev:      metric.StdDev(),
				Variance:    metric.Variance(),
				Percentiles: metric.Percentiles(percentiles),
			}
			metric.Clear()
		}
	})
	return stats
}

// Each - Dummy Each
func (dr DummyRegistry) Each(_ func(string, interface{})) {}

// Get - Dummy Get
func (dr DummyRegistry) Get(_ string) interface{} { return nil }

// GetAll - Dummy GetAll
func (dr DummyRegistry) GetAll() map[string]map[string]interface{} {
	return map[string]map[string]interface{}{}
}

// GetOrRegister - Dummy GetOrRegister
func (dr DummyRegistry) GetOrRegister(_ string, i interface{}) interface{} { return i }

// Register - Dummy Register
func (dr DummyRegistry) Register(_ string, _ interface{}) error { return nil }

// RunHealthchecks - Dummy RunHealthchecks
func (dr DummyRegistry) RunHealthchecks() {}

// Unregister - Dummy Unregister
func (dr DummyRegistry) Unregister(_ string) {}

// UnregisterAll - Dummy UnregisterAll
func (dr DummyRegistry) UnregisterAll() {}

// processMetrics - 지정된 시간을 기준으로 Metrics snapshot 처리
func (m *Metrics) processMetrics(ctx context.Context) {
	// 서비스 기준 Metric 설정 Child (GC, MEM)
	r := gometrics.NewPrefixedChildRegistry(*(m.Registry), "service.")
	gometrics.RegisterDebugGCStats(r)
	gometrics.RegisterRuntimeMemStats(r)

	// 지정된 시간 간격으로 metric 수집 및 snapshot 처리 (Context 종료 시점까지 반복)
	go func() {
		ticker := time.NewTicker(m.Config.CollectionPeriod)
		for {
			select {
			case <-ticker.C:
				gometrics.CaptureDebugGCStatsOnce(r)
				gometrics.CaptureRuntimeMemStatsOnce(r)
				// Router Metrics 수집
				m.Router.Aggregate()
				// Snapshot 처리
				m.latestSnapshot = m.TakeSnapshot()
			case <-ctx.Done():
				return
			}
		}
	}()
}

// ===== [ Private Functions ] =====

// parseConfig - Metrics 운영을 위한 Configuration parsing 처리
func parseConfig(sConf config.ServiceConfig) *Config {
	conf := new(Config)
	tmp, ok := sConf.Middleware[MWNamespace]
	if !ok {
		return nil
	}

	buf := new(bytes.Buffer)
	yaml.NewEncoder(buf).Encode(tmp)
	if err := yaml.NewDecoder(buf).Decode(conf); err != nil {
		return nil
	}

	return conf
}

// ===== [ Public Functions ] =====

// NewDummyRegistry - 설정이 없는 경우의 빈 Registry 생성
func NewDummyRegistry() gometrics.Registry {
	return DummyRegistry{}
}

// SetupAndCreate - Metrics Producer 생성
func SetupAndCreate(ctx context.Context, sConfig config.ServiceConfig, logger logging.Logger) *Metrics {
	// Root Registry 생성 (cb-restapigw)
	registry := gometrics.NewPrefixedRegistry(core.AppName + ".")

	conf := parseConfig(sConfig)
	if conf == nil {
		// no middleware configuration defined or error for the metrics module, use DummyRegistry
		registry = NewDummyRegistry()
		return &Metrics{
			logger:   &logger,
			Registry: &registry,
			Router:   &RouterMetrics{},
			Proxy:    &ProxyMetrics{},
		}
	}

	// Metric Producer 생성
	m := Metrics{
		logger:         &logger,
		Config:         conf,
		Router:         NewRouterMetrics(&registry), // Router Metrics
		Proxy:          NewProxyMetrics(&registry),  // Pipeline Metrics
		Registry:       &registry,
		latestSnapshot: NewStats(), // 최종 Snapshot 관리용
	}

	// Run metrics
	m.processMetrics(ctx)

	return &m
}
