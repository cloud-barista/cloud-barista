package manager

import (
	"github.com/sirupsen/logrus"
)

type CollectorScheduler struct {
	cm *CollectManager
}

func NewCollectorScheduler(cm *CollectManager) CollectorScheduler {
	return CollectorScheduler{cm: cm}
}

// 콜렉터 스케질 인/아웃 조건 체크
func (c CollectorScheduler) CheckScaleCondition() error {

	totalHostCnt := len(c.cm.AgentQueueTTL)
	isScaling := false
	scaleCnt := 0
	scalingEvent := ""

	maxHostCount := c.cm.Config.Monitoring.MaxHostCount
	currentCollectorN := len(c.cm.CollectorIdx)
	collectorAddr := c.cm.CollectorUUIDAddr

	if maxHostCount*currentCollectorN < totalHostCnt {

		scaleCnt = totalHostCnt/maxHostCount - currentCollectorN

		if totalHostCnt%maxHostCount != 0 {
			scaleCnt += 1
		}
		isScaling = true
		scalingEvent = "out"
	}

	for _, colAddr := range collectorAddr {

		if currentCollectorN != 1 && len((*colAddr).MarkingAgent) == 0 {
			isScaling = true
			scalingEvent = "in"
			break
		}
	}

	if isScaling {
		var err error
		if scalingEvent == "in" {
			err = c.ScaleIn()
		} else if scalingEvent == "out" {
			err = c.ScaleOut(scaleCnt)
		}
		if err != nil {
			return err
		}
	}

	return nil
}

func (c CollectorScheduler) ScaleIn() error {

	collectorIdx := c.cm.CollectorIdx
	collectorUUIDAddr := c.cm.CollectorUUIDAddr

	for idx, uuid := range collectorIdx {

		if len((*collectorUUIDAddr[uuid]).MarkingAgent) == 0 {
			c.cm.CollectorIdx = c.cm.CollectorIdx[:idx]
			c.cm.StopCollector(uuid)
		}
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
