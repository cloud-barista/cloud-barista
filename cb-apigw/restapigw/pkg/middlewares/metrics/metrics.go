package metrics

import (
	"bytes"
	"context"
	"time"

	"github.com/cloud-barista/cb-apigw/restapigw/pkg/config"
	"github.com/cloud-barista/cb-apigw/restapigw/pkg/core"
	"github.com/cloud-barista/cb-apigw/restapigw/pkg/logging"
	"github.com/rcrowley/go-metrics"
	"gopkg.in/yaml.v3"
)

// ===== [ Constants and Variables ] =====

const (
	// MWNamespace - Middleware 식별자
	MWNamespace = "mw-metrics"
)

var (
	// defaultListenAddress - Metrics 수집을 위한 Pulling 용 Port
	//defaultListenAddress = ":9000"

	// Metrics 의 수행 구간 비율 측정용
	percentiles = []float64{0.1, 0.25, 0.5, 0.75, 0.9, 0.95, 0.99}
	// 기본 Metrics 샘플
	defaultSample = func() metrics.Sample { return metrics.NewUniformSample(1028) }
)

// ===== [ Types ] =====

type (
	// DummyRegistry - 설정이 없는 경우에 활용할 Registry 구조 정의
	DummyRegistry struct{}

	// Config - Metrics 운영에 필요한 설정 구조 정의
	Config struct {
		// Router Metrics 활성화 여부 (기본값: true)
		RouterEnabled bool `yaml:"router_enabled"`
		// Proxy Metrics 활성화 여부 (기본값: true)
		ProxyEnabled bool `yaml:"proxy_enabled"`
		// Backend Metrics 활성화 여부 (기본값: true)
		BackendEnabled bool `yaml:"backend_enabled"`
		// Metrics 수집 시간 (기본값: 0, Metrics 수집하지 않음)
		CollectionPeriod time.Duration `yaml:"collection_period"`
		// Metrics Expose를 위한 HTTP/Gin Server 구동 비활성화 여부 (기본값: false)
		ExposeMetrics bool `yaml:"expose_metrics"`
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

	// Producer - Metrics 정보 관리 구조
	Producer struct {
		Logger   *logging.Logger
		Config   *Config
		Proxy    *ProxyMetrics
		Router   *RouterMetrics
		Registry *metrics.Registry

		latestSnapshot Stats
	}
)

// ===== [ Implementations ] =====

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
func (mp *Producer) processMetrics(ctx context.Context, period time.Duration, log logging.Logger) {
	// 서비스 기준 Metric 설정 Child (GC, MEM)
	cr := metrics.NewPrefixedChildRegistry(*(mp.Registry), "service.")

	metrics.RegisterDebugGCStats(cr)
	metrics.RegisterRuntimeMemStats(cr)

	// 지정된 시간 간격으로 metric 수집 및 snapshot 처리 (Context 종료 시점까지 반복)
	go func() {
		ticker := time.NewTicker(period)
		for {
			select {
			case <-ticker.C:
				metrics.CaptureDebugGCStatsOnce(cr)
				metrics.CaptureRuntimeMemStatsOnce(cr)
				// Router Metrics 수집
				mp.Router.Aggregate()
				// Snapshot 처리
				mp.latestSnapshot = mp.TakeSnapshot()
			case <-ctx.Done():
				return
			}
		}
	}()
}

// Snapshot - 최종 처리된 Snapshot 반환 (InfluxDB 쪽에서 호출)
func (mp *Producer) Snapshot() Stats {
	return mp.latestSnapshot
}

// TakeSnapshot - 현재 상태에 대한 Snapshot 처리
func (mp *Producer) TakeSnapshot() Stats {
	stats := NewStats()

	(*mp.Registry).Each(func(key string, val interface{}) {
		switch metric := val.(type) {
		case metrics.Counter:
			stats.Counters[key] = metric.Count()
		case metrics.Gauge:
			stats.Gauges[key] = metric.Value()
		case metrics.Histogram:
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

// ===== [ Private Functions ] =====

// parseConfig - Metrics 운영을 위한 Configuration 설정
func parseConfig(mConf config.MWConfig) *Config {
	conf := new(Config)

	tmp, ok := mConf[MWNamespace]
	if !ok {
		return nil
	}

	buf := new(bytes.Buffer)
	yaml.NewEncoder(buf).Encode(tmp)
	if err := yaml.NewDecoder(buf).Decode(conf); nil != err {
		return nil
	}

	return conf
}

// ===== [ Public Functions ] =====

// NewDummyRegistry - 설정이 없는 경우의 빈 Registry 생성
func NewDummyRegistry() metrics.Registry {
	return DummyRegistry{}
}

// New - Metrics Collector에서 사용할 Metrics 정보 관리 Producer 인스턴스 생성
func New(ctx context.Context, mConf config.MWConfig, log logging.Logger) *Producer {
	// 설정 기반으로 Metric Producer 구성
	conf := parseConfig(mConf)
	if nil == conf {
		// 설정이 존재하지 않는 경우는 장애 방지를 위한 Dummy Registry 활용
		dummyRegistry := NewDummyRegistry()
		return &Producer{
			Registry: &dummyRegistry,
			Router:   &RouterMetrics{},
			Proxy:    &ProxyMetrics{},
		}
	}

	// Metrics 처리를 위한 Root Registry 생성
	registry := metrics.NewPrefixedRegistry(core.AppName + ".")

	metricProducer := Producer{
		Logger:         &log,
		Config:         conf,
		Router:         NewRouterMetrics(&registry),
		Proxy:          NewProxyMetrics(&registry),
		Registry:       &registry,
		latestSnapshot: NewStats(),
	}

	metricProducer.processMetrics(ctx, metricProducer.Config.CollectionPeriod, log)

	return &metricProducer
}
