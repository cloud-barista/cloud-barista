package pull

import (
	"fmt"
	"sync"
	"time"

	"github.com/cloud-barista/cb-dragonfly/pkg/api/core/agent"
	"github.com/cloud-barista/cb-dragonfly/pkg/config"
	"github.com/cloud-barista/cb-dragonfly/pkg/modules/procedure/pull/puller"
)

type PullManager struct {
	AgentList map[string]agent.AgentInfo
	WaitGroup *sync.WaitGroup
}

func NewPullManager() (*PullManager, error) {
	pullManager := PullManager{
		WaitGroup: &sync.WaitGroup{},
	}
	return &pullManager, nil
}

func (pm *PullManager) StartPullCaller() error {
	for {

		pullingInterval := time.Duration(config.GetInstance().Monitoring.PullerInterval)

		// PULL 콜러 모듈 실행
		err := pm.syncAgentList()
		if err != nil {
			fmt.Println(err)
			return err
		}

		// 에이전트가 없을 경우
		if len(pm.AgentList) == 0 {
			time.Sleep(pullingInterval * time.Second)
			continue
		}

		pullCaller, err := puller.NewPullCaller(pm.AgentList)
		if err != nil {
			fmt.Println(err)
			return err
		}
		go pullCaller.StartPull()

		time.Sleep(pullingInterval * time.Second)
	}
}

func (pm *PullManager) StopPullCaller() error {
	return nil
}

func (pm *PullManager) syncAgentList() error {
	syncedAgentList, err := agent.ListAgent()
	if err != nil {
		fmt.Println(err)
		return err
	}
	pm.AgentList = syncedAgentList
	return nil
}
