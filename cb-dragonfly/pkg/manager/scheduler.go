package manager

import (
	"fmt"
	"strconv"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/sirupsen/logrus"

	"github.com/cloud-barista/cb-dragonfly/pkg/kafka"
	"github.com/cloud-barista/cb-dragonfly/pkg/localstore"
	"github.com/cloud-barista/cb-dragonfly/pkg/types"
	"github.com/cloud-barista/cb-dragonfly/pkg/util"
)

type CollectorScheduler struct {
	cm *CollectManager
	tm *TopicManager
}

func NewCollectorScheduler(manager *CollectManager) (*CollectorScheduler, error) {

	cScheduler := CollectorScheduler{
		tm: TopicMangerInstance(),
		cm: manager,
	}

	return &cScheduler, nil
}

func (cScheduler CollectorScheduler) Scheduler() error {

	kafkaAddr, err := kafka.GetInstance()
	if err != nil {
		return err
	}
	currentTopicsState := util.GetAllTopicBySort(kafkaAddr.GetAllTopics())
	beforeTopicsState := currentTopicsState
	beforeMaxHostCount, _ := strconv.Atoi(localstore.GetInstance().StoreGet(types.MONCONFIG + "/" + "max_host_count"))
	currentMaxHostCount := beforeMaxHostCount

	topicListChanged := !cmp.Equal(beforeTopicsState, currentTopicsState)
	maxHostCountChanged := !(beforeMaxHostCount == currentMaxHostCount)

	if cScheduler.cm.collectorPolicy == types.AGENTCOUNT {
		cScheduler.ScheduleBasedTheNumberOfCollector(currentTopicsState, beforeTopicsState, beforeMaxHostCount, currentMaxHostCount, topicListChanged, maxHostCountChanged)
	}
	if cScheduler.cm.collectorPolicy == types.CSP {
		cScheduler.ScheduleBasedCollectorCSPType(currentTopicsState, beforeTopicsState, topicListChanged)
	}
	return nil
}

//################# ScheduleBasedTheNumberOfCollector Start ################
func (cScheduler CollectorScheduler) ScheduleBasedTheNumberOfCollector(currentTopicsState []string, beforeTopicsState []string, beforeMaxHostCount int, currentMaxHostCount int, topicListChanged bool, maxHostCountChanged bool) {

	kafkaAddr, _ := kafka.GetInstance()
	// Init
	cScheduler.tm.SetTopicToCollectorBasedTheNumberOfAgent(currentTopicsState, currentMaxHostCount)
	cScheduler.NeedCollectorScaleInOut()
	cScheduler.SendTopicsToCollectors()

	for {
		aggreTime, _ := strconv.Atoi(localstore.GetInstance().StoreGet(types.MONCONFIG + "/" + "collector_interval"))
		time.Sleep(time.Duration(aggreTime) * time.Second)
		switch {
		case maxHostCountChanged:
			err := cScheduler.tm.DeleteAllTopicsInfo()
			if err != nil {
				logrus.Debug(err)
			}
			cScheduler.tm.SetTopicToCollectorBasedTheNumberOfAgent(currentTopicsState, currentMaxHostCount)
			cScheduler.NeedCollectorScaleInOut()
			break
		case topicListChanged:
			if !cScheduler.NeedRebalancedTopics(currentTopicsState, currentMaxHostCount) {
				deletedTopicList, newTopicList := cScheduler.ReturnDiffTopics(beforeTopicsState, currentTopicsState)
				err := cScheduler.tm.DeleteTopics(deletedTopicList)
				if err != nil {
					logrus.Debug(err)
				}
				err = cScheduler.tm.AddNewTopics(newTopicList, currentMaxHostCount)
				if err != nil {
					logrus.Debug(err)
				}
			}
			cScheduler.NeedCollectorScaleInOut()
			break
		}
		cScheduler.SendTopicsToCollectors()
		beforeTopicsState = currentTopicsState
		currentTopicsState = util.GetAllTopicBySort(kafkaAddr.GetAllTopics())
		fmt.Println(fmt.Sprintf("##### %s : %s #####", "All topics from kafka", currentTopicsState))
		beforeMaxHostCount = currentMaxHostCount
		currentMaxHostCount, _ = strconv.Atoi(localstore.GetInstance().StoreGet(types.MONCONFIG + "/" + "max_host_count"))

		topicListChanged = !cmp.Equal(beforeTopicsState, currentTopicsState)
		maxHostCountChanged = !(beforeMaxHostCount == currentMaxHostCount)
	}
}

func (cScheduler CollectorScheduler) NeedCollectorScaleInOut() {
	var err error
	var idealCollectorCnt int
	if len(cScheduler.tm.IdealCollectorPerAgentCntSlice) == 0 {
		idealCollectorCnt = 1
	} else {
		idealCollectorCnt = len(cScheduler.tm.IdealCollectorPerAgentCntSlice)
	}
	scaleCnt := idealCollectorCnt - len(cScheduler.cm.CollectorGroupManageMap)
	if scaleCnt != 0 {
		for needScalingCnt := scaleCnt; needScalingCnt != 0; {
			if needScalingCnt > 0 {
				err = cScheduler.cm.CreateCollectorGroup()
				needScalingCnt--
			} else {
				err = cScheduler.cm.StopCollectorGroup()
				needScalingCnt++
			}
			if err != nil {
				logrus.Debug(err)
			}
		}
	}
}

func (cScheduler CollectorScheduler) NeedRebalancedTopics(currentTopicsState []string, currentMaxHostCount int) bool {
	if len(cScheduler.tm.IdealCollectorPerAgentCntSlice) == ((len(currentTopicsState) / currentMaxHostCount) + 1) {
		return false
	} else {
		err := cScheduler.tm.DeleteAllTopicsInfo()
		if err != nil {
			logrus.Debug(err)
		}
		cScheduler.tm.SetTopicToCollectorBasedTheNumberOfAgent(currentTopicsState, currentMaxHostCount)
		return true
	}
}

//################# ScheduleBasedTheNumberOfCollector End ################

//################# ScheduleBasedCollectorCSPType Start ################
func (cScheduler CollectorScheduler) ScheduleBasedCollectorCSPType(currentTopicsState []string, beforeTopicsState []string, topicListChanged bool) {
	kafkaAddr, _ := kafka.GetInstance()
	// Init
	cScheduler.tm.SetTopicToCollectorBasedCSPTypeOfAgent(currentTopicsState)
	cScheduler.SendTopicsToCollectors()

	for {
		aggreTime, _ := strconv.Atoi(localstore.GetInstance().StoreGet(types.MONCONFIG + "/" + "collector_interval"))
		time.Sleep(time.Duration(aggreTime) * time.Second)
		switch {
		case topicListChanged:
			deletedTopicList, newTopicList := cScheduler.ReturnDiffTopics(beforeTopicsState, currentTopicsState)
			err := cScheduler.tm.DeleteTopics(deletedTopicList)
			if err != nil {
				logrus.Debug(err)
			}
			err = cScheduler.tm.AddNewTopicsOnCSPCollector(newTopicList)
			if err != nil {
				logrus.Debug(err)
			}
			break
		}
		cScheduler.SendTopicsToCollectors()
		beforeTopicsState = currentTopicsState
		currentTopicsState = util.GetAllTopicBySort(kafkaAddr.GetAllTopics())
		fmt.Println(fmt.Sprintf("##### %s : %s #####", "All topics from kafka", currentTopicsState))

		topicListChanged = !cmp.Equal(beforeTopicsState, currentTopicsState)
	}
}

//################# ScheduleBasedCollectorCSPType End ################

//################# Common methods ################
func (cScheduler CollectorScheduler) ReturnDiffTopics(beforeTopics []string, currentTopics []string) ([]string, []string) {
	return util.ReturnDiffTopicList(beforeTopics, currentTopics), util.ReturnDiffTopicList(currentTopics, beforeTopics)
}

func (cScheduler CollectorScheduler) SendTopicsToCollectors() {
	for idx, cAddrList := range cScheduler.cm.CollectorGroupManageMap {
		for _, cAddr := range cAddrList {
			(*cAddr).Ch <- localstore.GetInstance().StoreGet(fmt.Sprintf("%s/%d", types.COLLECTORGROUPTOPIC, idx))
		}
	}
}
