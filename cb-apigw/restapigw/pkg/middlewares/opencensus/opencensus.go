// Package opencensus - Provide functions for Opencensus integration
package opencensus

import (
	"bytes"
	"context"
	"errors"
	"sync"
	"time"

	"github.com/cloud-barista/cb-apigw/restapigw/pkg/config"
	"github.com/cloud-barista/cb-apigw/restapigw/pkg/core"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/trace"
	"gopkg.in/yaml.v3"
)

// ===== [ Constants and Variables ] =====

const (
	// MWNamespace - Middleware 식별자
	MWNamespace = "mw-opencensus"
	// ContextKey - Request trace 식별자
	ContextKey = core.AppName + "-request-span"
)

var (
	// DefaultViews - Opencensus Trace와 통합되기 위한 View들 정의
	DefaultViews = []*view.View{
		ochttp.ClientSentBytesDistribution,
		ochttp.ClientReceivedBytesDistribution,
		ochttp.ClientRoundtripLatencyDistribution,
		ochttp.ClientCompletedCount,

		ochttp.ServerRequestCountView,
		ochttp.ServerRequestBytesView,
		ochttp.ServerResponseBytesView,
		ochttp.ServerLatencyView,
		ochttp.ServerRequestCountByMethod,
		ochttp.ServerResponseCountByStatusCode,
	}

	// register - Opencensus 연동을 위한 정보를 관리
	register = composableRegister{
		viewExporter:       registerViewExporter,
		traceExporter:      registerTraceExporter,
		setDefaultSampler:  setDefaultSampler,
		setReportingPeriod: setReportingPeriod,
		registerViews:      registerViews,
	}

	exporterFactories = []ExporterFactory{}

	errNoMWConfig                         = errors.New("no middleware configuration defined for the opencensus module")
	errSingletonExporterFactoriesRegister = errors.New("expecting only one exporter factory registration per instance")

	registerOnce  = new(sync.Once)
	enabledLayers EnabledLayers
	mu            = new(sync.RWMutex)
)

// ===== [ Types ] =====

// composableRegister - Opencensus 연동을 위한 정보 구조
type composableRegister struct {
	viewExporter       func(exporters ...view.Exporter)
	traceExporter      func(exporters ...trace.Exporter)
	registerViews      func(views ...*view.View) error
	setDefaultSampler  func(rate int)
	setReportingPeriod func(d time.Duration)
}

// ExporterFactory - Opencenus Exporter관리를 위한 팩토리 함수 정의
type ExporterFactory func(context.Context, Config) (interface{}, error)

// EnabledLayers - Opencensus와 연동하기 위한 레이어 활성화 구조
type EnabledLayers struct {
	// Router 활성화 여부
	Router bool `yaml:"router"`
	// Proxy 활성화 여부
	Proxy bool `yaml:"proxy"`
	// Backend 활성화 여부
	Backend bool `yaml:"backend"`
}

// Config - Opencencus 연동을 위한 설정 구조 정의
type Config struct {
	// 샘플링 비율 (기본값: 0, no sampling, 100 이면 전수 검사)
	SampleRate int `yaml:"sample_rate"`
	// Exporter로 데이터 처리를 위한 주기
	ReportingPeriod time.Duration `yaml:"reporting_period"`
	// 레이어별 활성화 여부
	EnabledLayers *EnabledLayers `yaml:"enabled_layers"`
	// 데이터 처리를 위한 Exporter 설정
	Exporters struct {
		Jaeger *struct {
			// Jaeger 서버의 Endpoint
			Endpoint string `yaml:"endpoint"`
			// Jaeger 서버의 식별을 위한 서비스 명
			ServiceName string `yaml:"service_name"`
		} `yaml:"jaeger"`
	} `yaml:"exporters"`
}

// ===== [ Implementations ] =====

// ExporterFactories - 지정된 정보를 기준으로 Opencensus 연동 Registry에 Exporter 등록 처리
func (c *composableRegister) ExporterFactories(ctx context.Context, conf Config, fs []ExporterFactory) {
	viewExporters := []view.Exporter{}
	traceExporters := []trace.Exporter{}

	for _, f := range fs {
		e, err := f(ctx, conf)
		if err != nil {
			continue
		}
		if ve, ok := e.(view.Exporter); ok {
			viewExporters = append(viewExporters, ve)
		}
		if te, ok := e.(trace.Exporter); ok {
			traceExporters = append(traceExporters, te)
		}
	}

	c.viewExporter(viewExporters...)
	c.traceExporter(traceExporters...)
}

// Reegister - 지정된 정보를 기준으로 Opencensus 연동 Register에 Views 와 옵션 정보 (Sampling Rate, Reporting Period) 등록 처리
func (c composableRegister) Register(ctx context.Context, conf Config, vs []*view.View) error {
	if len(vs) == 0 {
		vs = DefaultViews
	}

	c.setDefaultSampler(conf.SampleRate)
	c.setReportingPeriod(conf.ReportingPeriod)

	return c.registerViews(vs...)
}

// ===== [ Private Functions ] =====

// registerViewExporter - Opencensus에 View Exporter들 등록
func registerViewExporter(exporters ...view.Exporter) {
	for _, e := range exporters {
		view.RegisterExporter(e)
	}
}

// registerTraceExporter - Opencensus에 Trace Exporter들 등록
func registerTraceExporter(exporters ...trace.Exporter) {
	for _, e := range exporters {
		trace.RegisterExporter(e)
	}
}

// 지정한 비율로 Opencensus에 Sampling 비율 설정
func setDefaultSampler(rate int) {
	var sampler trace.Sampler
	switch {
	case rate <= 0:
		sampler = trace.NeverSample()
	case rate >= 100:
		sampler = trace.AlwaysSample()
	default:
		sampler = trace.ProbabilitySampler(float64(rate) / 100.0)
	}
	trace.ApplyConfig(trace.Config{DefaultSampler: sampler})
}

// 지정한 시간을 Opencensus의 Reporting 주기 등록
func setReportingPeriod(d time.Duration) {
	view.SetReportingPeriod(d)
}

// registerViews - 지정된 View들을 Opencensus에 등록
func registerViews(views ...*view.View) error {
	return view.Register(views...)
}

// fromContext - 지정한 Context 기반으로 동작하는 Trace Span 구성 및 반환
func fromContext(ctx context.Context) *trace.Span {
	span := trace.FromContext(ctx)
	if span == nil {
		span, _ = ctx.Value(ContextKey).(*trace.Span)
	}
	return span
}

// parseConfig - 서비스 설정에서 Opencensus Trace에 관련된 설정 Parsing
func parseConfig(sConf config.ServiceConfig) (*Config, error) {
	conf := new(Config)
	tmp, ok := sConf.Middleware[MWNamespace]
	if !ok {
		return nil, errNoMWConfig
	}

	buf := new(bytes.Buffer)
	yaml.NewEncoder(buf).Encode(tmp)
	if err := yaml.NewDecoder(buf).Decode(conf); err != nil {
		return nil, err
	}

	return conf, nil
}

// ===== [ Public Functions ] =====

// Setup - Opencensus Trace 관련 설정을 검증하고 Exporter들을 Registry에 등록하고 Trace 대상 Layer들 구성
func Setup(ctx context.Context, sConf config.ServiceConfig, vs ...*view.View) error {
	cfg, err := parseConfig(sConf)
	if err != nil {
		return err
	}

	err = errSingletonExporterFactoriesRegister
	registerOnce.Do(func() {
		register.ExporterFactories(ctx, *cfg, exporterFactories)

		err = register.Register(ctx, *cfg, vs)
		if err != nil {
			return
		}

		if cfg.EnabledLayers != nil {
			enabledLayers = *cfg.EnabledLayers
			return
		}

		enabledLayers = EnabledLayers{true, true, true}
	})

	return err
}

// RegisterExporterFactories - Opencensus Exporter 팩토리들을 등록
func RegisterExporterFactories(ef ExporterFactory) {
	mu.Lock()
	exporterFactories = append(exporterFactories, ef)
	mu.Unlock()
}

// IsRouterEnabled - Router Layer 활성화 여부
func IsRouterEnabled() bool {
	return enabledLayers.Router
}

// IsProxyEnabled - Proxy Layer 활성화 여부
func IsProxyEnabled() bool {
	return enabledLayers.Proxy
}

// IsBackendEnabled - Backend Layer 활성화 여부
func IsBackendEnabled() bool {
	return enabledLayers.Backend
}
