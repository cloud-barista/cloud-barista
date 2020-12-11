package cmd

import (
	"context"

	"github.com/cloud-barista/cb-apigw/restapigw/pkg/api"
	"github.com/cloud-barista/cb-apigw/restapigw/pkg/config"
	"github.com/cloud-barista/cb-apigw/restapigw/pkg/core"
	"github.com/cloud-barista/cb-apigw/restapigw/pkg/logging"
	"github.com/cloud-barista/cb-apigw/restapigw/pkg/server"

	// Opencensus 연동을 위한 Exporter 로드 및 초기화
	_ "github.com/cloud-barista/cb-apigw/restapigw/pkg/middlewares/opencensus/exporters/jaeger"

	// 각종 필요 패키지 로드 및 초기화
	_ "github.com/cloud-barista/cb-apigw/restapigw/pkg/jwt/basic"
)

// ===== [ Constants and Variables ] =====

// ===== [ Types ] =====

// ===== [ Implementations ] =====

// ===== [ Private Functions ] =====

// setupRepository - HTTP Server와 Admin Server 운영에 사용할 API Route 정보 리파지토리 구성
func setupRepository(sConf *config.ServiceConfig, log logging.Logger) (api.Repository, error) {
	// API Routing 정보를 관리하는 Repository 구성
	repo, err := api.BuildRepository(sConf, sConf.Cluster.UpdateFrequency)
	if nil != err {
		return nil, err
	}

	return repo, nil
}

// ===== [ Public Functions ] =====

// SetupAndRun - API Gateway 서비스를 위한 Router 및 Pipeline 구성과 구동
func SetupAndRun(ctx context.Context, sConf *config.ServiceConfig) error {
	// Sets up the Logger (CB-LOG)
	logger := logging.NewLogger()

	// API 운영을 위한 라우팅 리파지토리 구성
	repo, err := setupRepository(sConf, *logger)
	if nil != err {
		logger.WithError(err).Error("[SERVER] Terminate with Errors")
		return nil
	}
	defer repo.Close()

	// API G/W Server 구동
	svr := server.New(
		server.WithServiceConfig(sConf),
		server.WithLogger(*logger),
		server.WithRepository(repo),
	)

	ctx = core.ContextWithSignal(ctx)
	svr.StartWithContext(ctx)

	svr.Wait()
	logger.Info("[Server] Shutting down")

	return nil
}
