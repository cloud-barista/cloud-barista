package manager

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"strconv"
)

type CollectorScheduler struct {
	cm *CollectManager
}

func NewCollectorScheduler(cm *CollectManager) CollectorScheduler {
	return CollectorScheduler{cm: cm}
}

// 콜렉터 스케질 인/아웃 조건 체크
func (c CollectorScheduler) CheckScaleCondition() error {

	// 전체 호스트 수 가져오기

	totalHostCnt := c.cm.HostCnt
	if totalHostCnt == 0 {
		hostNode, err := c.cm.Etcd.ReadMetric("/host")
		if err != nil {
			logrus.Error("failed to get total host count")
			return err
		}
		totalHostCnt = hostNode.Nodes.Len()
	}

	// 호스트 수 기준 스케일 인/아웃 기준 체크
	isScaling := false
	scalingEvent := ""
	scaleCnt := 0

	fmt.Println("collector len = " + strconv.Itoa(len(c.cm.CollectorList)))
	fmt.Println("aggregate len = " + strconv.Itoa(len(c.cm.CollectorChan)))

	if c.cm.Config.Monitoring.MaxHostCount*len(c.cm.CollectorList) < totalHostCnt {
		isScaling = true
		scalingEvent = "out"
		// 스케일 아웃 콜렉터 수 계산
		var collectCnt int
		if totalHostCnt%c.cm.Config.Monitoring.MaxHostCount == 0 {
			collectCnt = totalHostCnt / c.cm.Config.Monitoring.MaxHostCount
		} else {
			collectCnt = totalHostCnt/c.cm.Config.Monitoring.MaxHostCount + 1
		}
		scaleCnt = collectCnt - len(c.cm.CollectorList)
	} else if c.cm.Config.Monitoring.MaxHostCount*(len(c.cm.CollectorList)-1) >= totalHostCnt {
		isScaling = true
		scalingEvent = "in"

		// 1개 미만으로 떨어질 경우 default 1개
		if len(c.cm.CollectorList) == 1 {
			isScaling = false
			fmt.Println("default collector Count = 1")
		}

		// 스케일 인 콜렉터 수 계산
		var collectCnt int
		if totalHostCnt%c.cm.Config.Monitoring.MaxHostCount == 0 {
			collectCnt = totalHostCnt / c.cm.Config.Monitoring.MaxHostCount
		} else {
			collectCnt = totalHostCnt/c.cm.Config.Monitoring.MaxHostCount + 1
		}
		scaleCnt = len(c.cm.CollectorList) - collectCnt
	}

	// 콜렉터 스케일 인/아웃 이벤트 처리
	if isScaling {
		var err error
		if scalingEvent == "in" {
			err = c.ScaleIn(scaleCnt)
		} else if scalingEvent == "out" {
			err = c.ScaleOut(scaleCnt)
		}
		if err != nil {
			return err
		}
	}

	return nil
}

func (c CollectorScheduler) ScaleIn(scaleCnt int) error {
	for collectorId := range c.cm.CollectorList {
		if scaleCnt == 0 {
			break
		}
		if err := c.cm.StopCollector(collectorId); err != nil {
			logrus.Error("failed to stop collector")
			continue
		}
		scaleCnt--
	}
	return nil
}

func (c CollectorScheduler) ScaleOut(scaleCnt int) error {
	for i := 0; i < scaleCnt; i++ {
		if err := c.cm.CreateCollector(); err != nil {
			logrus.Error("failed to create collector")
			continue
		}
	}
	return nil
}
