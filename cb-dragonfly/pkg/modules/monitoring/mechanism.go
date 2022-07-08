package monitoring

import (
	"errors"
	"fmt"
	"sync"

	"github.com/cloud-barista/cb-dragonfly/pkg/config"
	"github.com/cloud-barista/cb-dragonfly/pkg/modules/monitoring/pull"
	"github.com/cloud-barista/cb-dragonfly/pkg/modules/monitoring/pull/puller"
	push_mcis "github.com/cloud-barista/cb-dragonfly/pkg/modules/monitoring/push/mcis"
	push_mck8s "github.com/cloud-barista/cb-dragonfly/pkg/modules/monitoring/push/mck8s"
	"github.com/cloud-barista/cb-dragonfly/pkg/storage/cbstore"
	"github.com/cloud-barista/cb-dragonfly/pkg/types"
	"github.com/cloud-barista/cb-dragonfly/pkg/util"
	"github.com/mitchellh/mapstructure"
)

func SetConfigurationToMemoryDB() {
	monConfigMap := map[string]interface{}{}
	err := mapstructure.Decode(config.GetInstance().Monitoring, &monConfigMap)
	if err != nil {
		util.GetLogger().Error(fmt.Sprintf("failed to decode monConfigMap, error=%s", err))
	}
	for key, val := range monConfigMap {
		err := cbstore.GetInstance().StorePut(types.MonConfig+"/"+key, fmt.Sprintf("%v", val))
		if err != nil {
			util.GetLogger().Error(fmt.Sprintf("failed to put monitoring configuration info, error=%s", err))
		}
	}
}

func startMCISPushModule(wg *sync.WaitGroup) error {

	// 콜렉터 매니저를 생성합니다.
	// 콜렉터 매니저는 collector 생성, 삭제 기능을 제공합니다.
	// 배포방식이 helm 일 경우, k8s와의 conn 및 configmap 을 생성합니다.
	cm, err := push_mcis.NewCollectorManager()
	if err != nil {
		util.GetLogger().Error("failed to initialize collector manager")
		return err
	}

	// 콜렉터 스케줄러를 생성합니다.
	// 콜렉터에게 분배할 topic 들을 관리하며 콜렉터의 배포 정책이 MaxAgentHost 일 경우,
	// 콜렉터 매니저의 콜렉터 생성 및 삭제 기능을 활용하여 콜렉터 스케일 인/아웃을 수행합니다.
	wg.Add(1)
	err = push_mcis.StartScheduler(wg, cm)
	if err != nil {
		return err
	}
	return nil
}

// startMCK8SPushModule MCK8S 수집 모듈 구동
func startMCK8SPushModule(wg *sync.WaitGroup) error {

	// 콜렉터 매니저 생성
	cm, err := push_mck8s.NewCollectorManager(wg)
	if err != nil {
		util.GetLogger().Error(err)
		return err
	}

	// 콜렉터 스케줄러 실행
	wg.Add(1)
	err = push_mck8s.StartScheduler(cm)
	if err != nil {
		util.GetLogger().Error(err)
		return err
	}
	return nil
}

func startMCISPullModule(wg *sync.WaitGroup) error {

	// PULL 매니저 생성
	pm, err := pull.NewPullManager()
	if err != nil {
		util.GetLogger().Error("Failed to initialize collector manager")
		return err
	}
	pa, err := puller.NewPullAggregator()
	if err != nil {
		util.GetLogger().Error("Failed to initialize Aggregator")
		return err
	}
	// PULL 콜러 실행
	wg.Add(1)
	go pm.StartPullCaller()
	// PULL Aggregator 실행
	wg.Add(1)
	go pa.StartAggregate()

	return nil
}

// TODO: MCK8S 환경 PULL 모듈 개발 시 활용
func startMCK8SPullModule(wg *sync.WaitGroup) error {
	return nil
}

func NewMechanism(wg *sync.WaitGroup) error {

	// Set Conf to InMemoryDB => Dragonfly의 config파일을 cb-store에 저장
	// cb-store의 기록 정보는 dragonfly의 모듈이 restart해도 지워지지 않습니다.
	SetConfigurationToMemoryDB()

	// Monitoring Policy => Push or Pull
	switch config.GetDefaultConfig().GetMonConfig().DefaultPolicy {
	case types.PushPolicy:
		// MCIS 수집 모듈 구동
		if err := startMCISPushModule(wg); err != nil {
			return err
		}
		// MCK8S 수집 모듈 구동
		if err := startMCK8SPushModule(wg); err != nil {
			return err
		}
		break
	case types.PullPolicy:
		// MCIS 수집 모듈 구동
		if err := startMCISPullModule(wg); err != nil {
			return err
		}
		break
		// TODO: MCK8S 수집 모듈 구동
	default:
		errMsg := "wrong monitoring mechanism config detected. change config to \"push\" or \"pull\""
		util.GetLogger().Error(errMsg)
		return errors.New(errMsg)
	}
	return nil
}
