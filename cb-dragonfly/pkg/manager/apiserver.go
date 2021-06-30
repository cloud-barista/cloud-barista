package manager

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/cloud-barista/cb-dragonfly/pkg/util"

	"github.com/cloud-barista/cb-dragonfly/pkg/api/rest/agent"
	"github.com/cloud-barista/cb-dragonfly/pkg/api/rest/alert"
	restconfig "github.com/cloud-barista/cb-dragonfly/pkg/api/rest/config"
	"github.com/cloud-barista/cb-dragonfly/pkg/api/rest/healthcheck"
	"github.com/cloud-barista/cb-dragonfly/pkg/api/rest/metric"
	"github.com/cloud-barista/cb-dragonfly/pkg/config"
	"github.com/cloud-barista/cb-dragonfly/pkg/core/alert/eventhandler"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type APIServer struct {
	echo *echo.Echo
}

// NewAPIServer API 서버 초기화
func NewAPIServer() (*APIServer, error) {
	e := echo.New()
	apiServer := APIServer{
		echo: e,
	}
	return &apiServer, nil
}

// StartAPIServer 모니터링 API 서버 실행
func (apiServer *APIServer) StartAPIServer(wg *sync.WaitGroup) error {
	defer wg.Done()
	util.GetLogger().Info("start CB-Dragonfly Framework API Server")

	// 모니터링 API 라우팅 룰 설정
	apiServer.SetRoutingRule(apiServer.echo)

	eventhandler.InitializeEventTypes()

	// 모니터링 API 서버 실행
	return apiServer.echo.Start(fmt.Sprintf(":%d", config.GetInstance().APIServer.Port))
}

func (apiServer *APIServer) SetRoutingRule(e *echo.Echo) {
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	}))

	dragonfly := e.Group("/dragonfly")

	// 멀티 클라우드 인프라 서비스 모니터링/실시간 모니터링 정보 조회
	//dragonfly.GET("/ns/:ns_id/mcis/:mcis_id/info", metric.GetMCISMonInfo)
	//dragonfly.GET("/ns/:ns_id/mcis/:mcis_id/rt-info", metric.GetMCISRealtimeMonInfo)

	// 멀티 클라우드 인프라 VM 모니터링/실시간 모니터링 정보 조회
	dragonfly.GET("/ns/:ns_id/mcis/:mcis_id/vm/:vm_id/metric/:metric_name/info", metric.GetVMMonInfo)

	// 멀티 클라우드 모니터링 정책 설정
	dragonfly.PUT("/config", restconfig.SetMonConfig)
	dragonfly.GET("/config", restconfig.GetMonConfig)
	dragonfly.PUT("/config/reset", restconfig.ResetMonConfig)

	// 에이전트 설치
	dragonfly.POST("/agent", agent.InstallTelegraf)
	// 에이전트 삭제
	dragonfly.DELETE("/agent", agent.UninstallAgent)
	// MCIS 모니터링 (Milkyway)
	dragonfly.GET("/ns/:ns_id/mcis/:mcis_id/vm/:vm_id/agent_ip/:agent_ip/mcis_metric/:metric_name/mcis-monitoring-info", metric.GetMCISMetric)
	// 멀티클라우드 인프라 VM 온디멘드 모니터링
	dragonfly.GET("/ns/:ns/mcis/:mcis_id/vm/:vm_id/agent_ip/:agent_ip/metric/:metric_name/ondemand-monitoring-info", metric.GetVMOnDemandMetric)

	// windows 에이전트 config, package 파일 다운로드
	dragonfly.GET("/installer/cbinstaller.zip", agent.GetWindowInstaller)
	//dragonfly.GET("/file/agent/conf", agent.GetTelegrafConfFile)
	//dragonfly.GET("/file/agent/pkg", agent.GetTelegrafPkgFile)

	// 에이전트 메타데이터 조회
	dragonfly.GET("/agent/metadata", agent.ListAgentMetadata)

	// 알람 조회, 생성, 삭제
	dragonfly.GET("/alert/tasks", alert.ListAlertTask)
	dragonfly.POST("/alert/task", alert.CreateAlertTask)
	dragonfly.GET("/alert/task/:task_id", alert.GetAlertTask)
	dragonfly.PUT("/alert/task/:task_id", alert.UpdateAlertTask)
	dragonfly.DELETE("/alert/task/:task_id", alert.DeleteAlertTask)

	// 알람 이벤트 핸들러 조회, 생성, 삭제
	dragonfly.GET("/alert/eventhandlers", alert.ListEventHandler)
	dragonfly.POST("/alert/eventhandler", alert.CreateEventHandler)
	dragonfly.GET("/alert/eventhandler/type/:type/event/:name", alert.GetEventHandler)
	dragonfly.PUT("/alert/eventhandler/type/:type/event/:name", alert.UpdateEventHandler)
	dragonfly.DELETE("/alert/eventhandler/type/:type/event/:name", alert.DeleteEventHandler)

	// 알람 이벤트 로그 조회, 생성
	dragonfly.POST("/alert/event", alert.CreateEventLog)
	dragonfly.GET("/alert/task/:task_id/events", alert.ListEventLog)

	// 헬스체크
	dragonfly.GET("/healthcheck", healthcheck.Ping)
}
