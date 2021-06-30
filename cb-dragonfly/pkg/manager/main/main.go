package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/cloud-barista/cb-dragonfly/pkg/metricstore/influxdb"
	v1 "github.com/cloud-barista/cb-dragonfly/pkg/metricstore/influxdb/v1"
	"github.com/cloud-barista/cb-dragonfly/pkg/puller"
	"github.com/cloud-barista/cb-dragonfly/pkg/util"

	grpc "github.com/cloud-barista/cb-dragonfly/pkg/api/grpc/server"
	"github.com/cloud-barista/cb-dragonfly/pkg/config"
	"github.com/cloud-barista/cb-dragonfly/pkg/core/alert/template"
	"github.com/cloud-barista/cb-dragonfly/pkg/manager"
	"github.com/cloud-barista/cb-dragonfly/pkg/types"
)

func startPushModule(wg *sync.WaitGroup) {
	// 콜렉터 매니저 생성
	cm, err := manager.NewCollectorManager()
	if err != nil {
		util.GetLogger().Error("failed to initialize collector manager")
		panic(err)
	}

	// CB-Store 모니터링 토픽 정보 초기화
	cm.SetConfigurationToMemoryDB()
	err = cm.StartCollectorGroup(wg)
	if err != nil {
		panic(err)
	}

	// 모니터링 콜렉터 스케일 인/아웃 스케줄러 실행
	wg.Add(1)
	err = cm.StartScheduler(wg)
	if err != nil {
		panic(err)
	}
}

func startPullModule(wg *sync.WaitGroup) {
	// PULL 매니저 생성
	pm, err := manager.NewPullManager()
	if err != nil {
		util.GetLogger().Error("Failed to initialize collector manager")
		panic(err)
	}
	pa, err := puller.NewPullAggregator()
	if err != nil {
		util.GetLogger().Error("Failed to initialize Aggregator")
		panic(err)
	}
	// PULL 콜러 실행
	wg.Add(1)
	go pm.StartPullCaller()
	// PULL Aggregator 실행
	wg.Add(1)
	go pa.StartAggregate()
}

func main() {

	time.Sleep(5 * time.Second)

	// 멀티 CPU 기반 고루틴 병렬 처리 활성화
	runtime.GOMAXPROCS(runtime.NumCPU())

	// 알람 모듈 템플릿 등록
	template.RegisterTemplate()

	// InfluxDB 클라이언트 설정
	influxDBConfig := config.GetDefaultConfig().GetInfluxDBConfig()
	influxDBAddr := fmt.Sprintf("%s:%d", influxDBConfig.EndpointUrl, influxDBConfig.ExternalPort)
	influxDBClientConfig := v1.Config{
		Addr:     influxDBAddr,
		Username: influxDBConfig.UserName,
		Password: influxDBConfig.Password,
	}
	err := influxdb.NewStorage(influxdb.V1, influxDBClientConfig)
	if err != nil {
		util.GetLogger().Error(fmt.Sprintf("failed to initialize influxDB, error=%s", err))
		panic(err)
	}

	// CB-Dragonfly config 정보 설정
	monConfig := config.GetDefaultConfig().GetMonConfig()

	// Push, Pull 메커니즘 기반 모니터링 모듈 동작
	var wg sync.WaitGroup
	if monConfig.DefaultPolicy == types.PushPolicy {
		startPushModule(&wg)
	} else if monConfig.DefaultPolicy == types.PullPolicy {
		startPullModule(&wg)
	}

	// 모니터링 API 서버 실행
	wg.Add(1)
	apiServer, err := manager.NewAPIServer()
	if err != nil {
		util.GetLogger().Error(fmt.Sprintf("failed to initialize api server, error=%s", err))
		panic(err)
	}
	go apiServer.StartAPIServer(&wg)

	// 모니터링 gRPC 서버 실행
	wg.Add(1)
	go grpc.StartGRPCServer()

	// 모든 고루틴이 종료될 때까지 대기
	wg.Wait()
}
