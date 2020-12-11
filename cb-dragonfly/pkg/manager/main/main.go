package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"runtime"
	"sync"
	"time"

	grpc "github.com/cloud-barista/cb-dragonfly/pkg/api/grpc/server"
	"github.com/cloud-barista/cb-dragonfly/pkg/core/alert/template"
	"github.com/cloud-barista/cb-dragonfly/pkg/manager"
)

func main() {

	time.Sleep(5 * time.Second)

	// 로그 파일 설정
	logrus.SetLevel(logrus.DebugLevel)
	logFileName := "cb-dragonfly.log"
	f, err := os.OpenFile(logFileName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	formatter := new(logrus.TextFormatter)
	logrus.SetFormatter(formatter)
	if err != nil {
		fmt.Println(err)
	} else {
		logrus.SetOutput(f)
	}

	// 멀티 CPU 기반 고루틴 병렬 처리 활성화
	runtime.GOMAXPROCS(runtime.NumCPU())

	template.RegisterTemplate()

	var wg sync.WaitGroup

	// 콜렉터 매니저 생성
	cm, err := manager.NewCollectorManager()
	if err != nil {
		logrus.Error("Failed to initialize collector manager")
		panic(err)
	}

	// 실시간 모니터링 데이터 초기화
	cm.FlushMonitoringData()
	err = cm.StartCollectorGroup(&wg)
	if err != nil {
		panic(err)
	}

	// 모니터링 콜렉터 스케일 인/아웃 스케줄러 실행
	wg.Add(1)
	err = cm.StartScheduler(&wg)
	if err != nil {
		panic(err)
	}

	// 모니터링 API 서버 실행
	wg.Add(1)
	apiServer, err := manager.NewAPIServer()
	if err != nil {
		logrus.Error("Failed to initialize api server")
		panic(err)
	}
	go apiServer.StartAPIServer(&wg)

	// 모니터링 gRPC 서버 실행
	wg.Add(1)
	go grpc.StartGRPCServer()

	// 모든 고루틴이 종료될 때까지 대기
	wg.Wait()
}
