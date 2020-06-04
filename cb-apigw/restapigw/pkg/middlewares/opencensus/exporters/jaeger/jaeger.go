package jaeger

import (
	"context"
	"errors"

	"contrib.go.opencensus.io/exporter/jaeger"
	"github.com/cloud-barista/cb-apigw/restapigw/pkg/middlewares/opencensus"
)

// ===== [ Constants and Variables ] =====

var (
	errDisabled = errors.New("opencensus-jaeger exporter disabled")
)

// ===== [ Types ] =====

// ===== [ Implementations ] =====

// ===== [ Private Functions ] =====

// init - 패키지 로드 시점의 초기화 작업 진행
func init() {
	opencensus.RegisterExporterFactories(func(ctx context.Context, conf opencensus.Config) (interface{}, error) {
		return Exporter(ctx, conf)
	})
}

// ===== [ Public Functions ] =====

// Exporter - Opencensus와 연동하기 위한 Jaeger 의 Exporter 반환
func Exporter(ctx context.Context, conf opencensus.Config) (*jaeger.Exporter, error) {
	if conf.Exporters.Jaeger == nil {
		return nil, errDisabled
	}
	exporter, err := jaeger.NewExporter(jaeger.Options{
		CollectorEndpoint: conf.Exporters.Jaeger.Endpoint,
		ServiceName:       conf.Exporters.Jaeger.ServiceName,
	})
	if err != nil {
		return exporter, err
	}

	go func() {
		<-ctx.Done()
		exporter.Flush()
	}()

	return exporter, nil
}
