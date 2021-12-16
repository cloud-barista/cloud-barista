package push

import (
	"context"
	"encoding/json"
	"fmt"
	que "github.com/Workiva/go-datastructures/queue"
	"github.com/cloud-barista/cb-dragonfly/pkg/api/core/agent"
	"github.com/cloud-barista/cb-dragonfly/pkg/config"
	"github.com/cloud-barista/cb-dragonfly/pkg/storage/cbstore"
	"github.com/cloud-barista/cb-dragonfly/pkg/types"
	"github.com/cloud-barista/cb-dragonfly/pkg/util"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strconv"
	"sync"
	"time"
)

type InMemoryTopic struct {
	TopicMap             map[int][]string
	CollectorPerAgentCnt []int
}

type CollectorScheduler struct {
	cm               *CollectManager
	inMemoryTopicMap *InMemoryTopic
	topicQue         *que.Queue
}

func StartScheduler(wg *sync.WaitGroup, manager *CollectManager) error {
	//// WaitGroup Start, Initialize Collector
	scheduler, err := NewCollectorScheduler(wg, manager)
	if err != nil {
		util.GetLogger().Error("Failed to initialize influxDB")
		return err
	}

	// Start Scheduler(Go-routine)
	go func() {
		err = scheduler.Scheduler()
	}()
	if err != nil {
		util.GetLogger().Error("Failed to make scheduler")
		return err
	}
	// WaitGroup End
	defer wg.Done()
	return nil
}

func NewCollectorScheduler(wg *sync.WaitGroup, manager *CollectManager) (*CollectorScheduler, error) {
	manager.WaitGroup = wg
	c := cbstore.GetInstance()
	inMemoryTopic := InMemoryTopic{
		TopicMap:             map[int][]string{},
		CollectorPerAgentCnt: []int{},
	}
	if config.GetInstance().Monitoring.DeployType == types.Helm {
		configMap, err := manager.K8sClientSet.CoreV1().ConfigMaps(types.Namespace).Get(context.TODO(), "cb-dragonfly-collector-configmap", metav1.GetOptions{})
		if err != nil {
			fmt.Println("Fail to Get ConfigMap")
			fmt.Println(err)
			return &CollectorScheduler{}, err
		}
		if err := json.Unmarshal(configMap.BinaryData["topicMap"], &inMemoryTopic.TopicMap); err != nil {
			fmt.Println("Fail to unMarshal ConfigMap Object Data")
		}
		for collectorIdx, topicSlice := range inMemoryTopic.TopicMap {
			for i := 0; i < len(topicSlice); i++ {
				_ = c.StorePut(fmt.Sprintf("%s/%s", types.Topic, topicSlice[i]), strconv.Itoa(collectorIdx))
			}
		}
	} else {
		_ = c.StoreDelList(types.Topic)
		cPolicy := c.StoreGet(types.CollectorPolicy)
		if cPolicy == manager.CollectorPolicy {
			if getCMapFromStore := c.StoreGet(fmt.Sprintf("%s", types.CollectorTopicMap)); getCMapFromStore != "" {
				_ = json.Unmarshal([]byte(getCMapFromStore), &inMemoryTopic)
				for collectorIdx, topicSlice := range inMemoryTopic.TopicMap {
					for i := 0; i < len(topicSlice); i++ {
						_ = c.StorePut(fmt.Sprintf("%s/%s", types.Topic, topicSlice[i]), strconv.Itoa(collectorIdx))
					}
				}
			}
		}
		_ = c.StorePut(types.CollectorPolicy, manager.CollectorPolicy)
	}
	cScheduler := CollectorScheduler{
		cm:               manager,
		inMemoryTopicMap: &inMemoryTopic,
		topicQue:         util.GetRingQueue(), // Set Global RingBuffer
	}
	return &cScheduler, nil
}

func (cScheduler CollectorScheduler) Scheduler() error {

	aggreTime, _ := strconv.Atoi(cbstore.GetInstance().StoreGet(types.MonConfig + "/" + "collector_interval"))
	topicQue := cScheduler.topicQue
	cPolicy := cScheduler.cm.CollectorPolicy

	for {
		time.Sleep(time.Duration(aggreTime) * time.Second)
		var addTopicList []string
		var delTopicList []string
		if topicQue.Len() != 0 {
			topicBytesList, err := topicQue.Get(topicQue.Len())
			if err != nil {
				return err
			}
			for _, topicBytes := range topicBytesList {
				topicStructure := util.TopicStructure{}
				if err := json.Unmarshal(topicBytes.([]byte), &topicStructure); err != nil {
					return err
				}
				if topicStructure.Policy == types.TopicAdd {
					addTopicList = append(addTopicList, topicStructure.Topic)
				} else if topicStructure.Policy == types.TopicDel {
					delTopicList = append(delTopicList, topicStructure.Topic)
				}
			}

			addTopicList = util.GetAllTopicBySort(util.Unique(addTopicList))
			delTopicList = util.GetAllTopicBySort(util.Unique(util.ReturnDiffTopicList(delTopicList, addTopicList)))
		}
		fmt.Println("### Now Scheduling ###")
		fmt.Println("## Add Topics Queue ##", addTopicList)
		fmt.Println("## Del Topics Queue ##", delTopicList)
		switch cPolicy {
		case types.AgentCntCollectorPolicy:
			cScheduler.SchedulePolicyBasedCollector(addTopicList, delTopicList)
			break
		case types.CSPCollectorPolicy:
			cScheduler.ScheduleCSPBasedCollector(addTopicList, delTopicList)
			break
		}
	}
}

/** ### AgentCnt Policy Start ### */
func (cScheduler CollectorScheduler) SchedulePolicyBasedCollector(addTopicList []string, delTopicList []string) {
	maxHostCount := config.GetInstance().GetMonConfig().MaxHostCount
	if len(addTopicList) != 0 {
		cScheduler.AddTopicsToCollector(addTopicList, maxHostCount)
	}
	if len(delTopicList) != 0 {
		cScheduler.DeleteTopicsToCollector(delTopicList)
	}
	cScheduler.BalanceTopicsToCollector(maxHostCount)
	cScheduler.ScaleInOutCollector()
	cScheduler.DistributeTopicsToCollector()
	cScheduler.WriteCollectorMapToInMemoryDB()
	return
}

func (cScheduler CollectorScheduler) AddTopicsToCollector(addTopicList []string, maxHostCount int) {
	cMap := *(cScheduler.inMemoryTopicMap)
	c := cbstore.GetInstance()

	for i := 0; i < len(addTopicList); i++ {
		topic := addTopicList[i]
		if c.StoreGet(fmt.Sprintf("%s/%s", types.Topic, topic)) != "" {
			continue
		}

		if len(cMap.CollectorPerAgentCnt) == 0 {
			cMap.TopicMap[0] = []string{topic}
			_ = c.StorePut(fmt.Sprintf("%s/%s", types.Topic, topic), "0")
			cMap.CollectorPerAgentCnt = []int{1}
			continue
		} else {
			needToScaleOut := true
			for collectorIdx, cnt := range cMap.CollectorPerAgentCnt {
				if cnt < maxHostCount {
					cMap.TopicMap[collectorIdx] = append(cMap.TopicMap[collectorIdx], topic)
					_ = c.StorePut(fmt.Sprintf("%s/%s", types.Topic, topic), strconv.Itoa(collectorIdx))
					cMap.CollectorPerAgentCnt[collectorIdx] += 1
					needToScaleOut = false
					break
				}
			}
			if needToScaleOut {
				cMap.TopicMap[len(cMap.CollectorPerAgentCnt)] = []string{topic}
				_ = c.StorePut(fmt.Sprintf("%s/%s", types.Topic, topic), strconv.Itoa(len(cMap.CollectorPerAgentCnt)))
				cMap.CollectorPerAgentCnt = append(cMap.CollectorPerAgentCnt, 1)
			}
		}
	}
	*(cScheduler.inMemoryTopicMap) = cMap
}

func (cScheduler CollectorScheduler) DeleteTopicsToCollector(delTopicList []string) {
	cMap := *(cScheduler.inMemoryTopicMap)
	if len(cMap.TopicMap) == 0 {
		return
	}
	c := cbstore.GetInstance()

	deleteTopicsMap := map[int][]string{}
	for i := 0; i < len(delTopicList); i++ {
		delTopic := delTopicList[i]
		collectorIdxStr := c.StoreGet(fmt.Sprintf("%s/%s", types.Topic, delTopic))
		collectorIdx, _ := strconv.Atoi(collectorIdxStr)
		_ = c.StoreDelete(fmt.Sprintf("%s/%s", types.Topic, delTopic))

		agentInfo := agent.AgentInfo{}
		agentInfoBytes := c.StoreGet(types.Agent + delTopic)
		if agentInfoBytes != "" {
			_ = json.Unmarshal([]byte(agentInfoBytes), &agentInfo)
			agentInfo.AgentHealth = string(agent.Unhealthy)
			agentInfo.AgentState = string(agent.Disable)
			recentAgentInfoBytes, _ := json.Marshal(agentInfo)
			_ = c.StorePut(types.Agent+delTopic, string(recentAgentInfoBytes))
		}
		deleteTopicsMap[collectorIdx] = append(deleteTopicsMap[collectorIdx], delTopic)
	}
	for cIdx, deleteTopicsArray := range deleteTopicsMap {
		processedTopicsArray := util.ReturnDiffTopicList(cMap.TopicMap[cIdx], util.GetAllTopicBySort(deleteTopicsArray))
		cMap.TopicMap[cIdx] = processedTopicsArray
		cMap.CollectorPerAgentCnt[cIdx] = len(cMap.TopicMap[cIdx])
	}
	*(cScheduler.inMemoryTopicMap) = cMap
}

func (cScheduler CollectorScheduler) BalanceTopicsToCollector(maxHostCount int) {
	cMap := *(cScheduler.inMemoryTopicMap)

	totalTopicsCnt := 0
	for _, topicCnt := range cMap.CollectorPerAgentCnt {
		totalTopicsCnt += topicCnt
	}
	idealCollectorCnt := util.CalculateNumberOfCollector(totalTopicsCnt, maxHostCount)
	if len(cMap.TopicMap) == idealCollectorCnt {
		if len(cMap.CollectorPerAgentCnt) == 1 && cMap.CollectorPerAgentCnt[0] == 0 {
			cScheduler.inMemoryTopicMap.CollectorPerAgentCnt = []int{}
			cScheduler.inMemoryTopicMap.TopicMap = map[int][]string{}
		}
		return
	} else {
		var totalTopics []string

		for _, topicsSlice := range cMap.TopicMap {
			totalTopics = append(totalTopics, topicsSlice...)
		}
		idealTm, idealCpc := util.MakeCollectorTopicMap(totalTopics, maxHostCount)

		cMap.TopicMap = idealTm
		cMap.CollectorPerAgentCnt = idealCpc

		for collectorIdx, collectorTopics := range idealTm {
			for i := 0; i < len(collectorTopics); i++ {
				_ = cbstore.GetInstance().StorePut(fmt.Sprintf("%s/%s", types.Topic, collectorTopics[i]), strconv.Itoa(collectorIdx))
			}
		}
		*(cScheduler.inMemoryTopicMap) = cMap
	}
}

func (cScheduler CollectorScheduler) ScaleInOutCollector() {
	collectorScale := len(cScheduler.inMemoryTopicMap.CollectorPerAgentCnt) - len(cScheduler.cm.CollectorAddrSlice)
	if collectorScale > 0 {
		for i := 0; i < collectorScale; i++ {
			err := cScheduler.cm.CreateCollector()
			if err != nil {
				fmt.Println("Fail to Create(Deploy) Collector")
				fmt.Println("errMsg: ", err)
			}
		}
	}
	if collectorScale < 0 {
		for i := 0; i < -collectorScale; i++ {
			err := cScheduler.cm.DeleteCollector()
			if err != nil {
				fmt.Println("Fail to Delete Collector")
				fmt.Println("errMsg: ", err)
			}
		}
	}
}

func (cScheduler CollectorScheduler) DistributeTopicsToCollector() {
	if config.GetInstance().Monitoring.DeployType == types.Helm {
		topicMapData := map[string][]byte{}
		topicMapBytes, _ := json.Marshal(cScheduler.inMemoryTopicMap.TopicMap)
		topicMapData["topicMap"] = topicMapBytes
		collectorUUIDMapData := map[string]string{}
		if len(cScheduler.cm.CollectorAddrSlice) != 0 {
			for _, collectorUUID := range cScheduler.cm.CollectorAddrSlice {
				collectorUUIDString := fmt.Sprintf("%p", collectorUUID)
				collectorUUIDMapData[collectorUUIDString] = "alive"
			}
		}
		configMap := &apiv1.ConfigMap{
			Data: collectorUUIDMapData,
			BinaryData: topicMapData,
			ObjectMeta: metav1.ObjectMeta{
				Name: types.ConfigMapName,
			}}
		configMapsClient := cScheduler.cm.K8sClientSet.CoreV1().ConfigMaps(types.Namespace)
		result, err := configMapsClient.Update(context.TODO(), configMap, metav1.UpdateOptions{})
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("updated ConfigMap: ", result.GetObjectMeta().GetName())
		}
	} else {
		for idx, topics := range cScheduler.inMemoryTopicMap.TopicMap {
			(*cScheduler.cm.CollectorAddrSlice[idx]).Ch <- topics
		}
	}
}

func (cScheduler CollectorScheduler) WriteCollectorMapToInMemoryDB() {
	inMemoryTopic := InMemoryTopic{
		TopicMap:             cScheduler.inMemoryTopicMap.TopicMap,
		CollectorPerAgentCnt: cScheduler.inMemoryTopicMap.CollectorPerAgentCnt,
	}
	cMapBytes, _ := json.Marshal(inMemoryTopic)
	_ = cbstore.GetInstance().StorePut(fmt.Sprintf("%s", types.CollectorTopicMap), string(cMapBytes))
}
/** ### AgentCnt Policy End ### */

/** ### CSP Policy Start ### */
func (cScheduler CollectorScheduler) ScheduleCSPBasedCollector(addTopicList []string, delTopicList []string) {
	cScheduler.ProvisioningCollectorByCSP()
	if len(addTopicList) != 0 {
		cScheduler.AddTopicsToCSPCollector(addTopicList)
	}
	if len(delTopicList) != 0 {
		cScheduler.DeleteTopicsToCSPCollector(delTopicList)
	}
	cScheduler.DistributeTopicsToCollector()
	cScheduler.WriteCollectorMapToInMemoryDB()
	return

}

func (cScheduler CollectorScheduler) ProvisioningCollectorByCSP() {
	provisioningCollectorCnt := types.TotalCspCnt - len(cScheduler.cm.CollectorAddrSlice)
	for i := 0; i < provisioningCollectorCnt; i++ {
		if err := cScheduler.cm.CreateCollector(); err != nil {
			fmt.Println("Fail to Create Collector")
		}
	}
}

func (cScheduler CollectorScheduler) AddTopicsToCSPCollector(addTopicList []string) {
	cMap := *(cScheduler.inMemoryTopicMap)
	c := cbstore.GetInstance()

	for i := 0; i < len(addTopicList); i++ {
		topic := addTopicList[i]
		if c.StoreGet(fmt.Sprintf("%s/%s", types.Topic, topic)) != "" {
			continue
		}
		collectorIdx := util.GetCspCollectorIdx(topic)
		if len(cMap.CollectorPerAgentCnt) == 0 {
			cMap.TopicMap[collectorIdx] = []string{topic}
			initCollectorPerAgentCnt := []int{}
			for i = 0; i < types.TotalCspCnt; i++ {
				topicCnt := 0
				if i == collectorIdx {
					topicCnt = 1
				}
				initCollectorPerAgentCnt = append(initCollectorPerAgentCnt, topicCnt)
			}
			cMap.CollectorPerAgentCnt = initCollectorPerAgentCnt
			continue
		}
		cMap.TopicMap[collectorIdx] = append(cMap.TopicMap[collectorIdx], topic)
		cMap.CollectorPerAgentCnt[collectorIdx] += 1
	}
	*(cScheduler.inMemoryTopicMap) = cMap
}

func (cScheduler CollectorScheduler) DeleteTopicsToCSPCollector(delTopicList []string) {
	cMap := *(cScheduler.inMemoryTopicMap)
	if len(cMap.TopicMap) == 0 {
		return
	}
	c := cbstore.GetInstance()

	deleteTopicsMap := map[int][]string{}
	for i := 0; i < len(delTopicList); i++ {
		delTopic := delTopicList[i]
		collectorIdx := util.GetCspCollectorIdx(delTopic)

		agentInfo := agent.AgentInfo{}
		agentInfoBytes := c.StoreGet(types.Agent + delTopic)
		_ = json.Unmarshal([]byte(agentInfoBytes), &agentInfo)
		agentInfo.AgentHealth = string(agent.Unhealthy)
		agentInfo.AgentState = string(agent.Disable)
		recentAgentInfoBytes, _ := json.Marshal(agentInfo)
		_ = c.StorePut(types.Agent+delTopic, string(recentAgentInfoBytes))

		deleteTopicsMap[collectorIdx] = append(deleteTopicsMap[collectorIdx], delTopic)
	}
	for cIdx, deleteTopicsArray := range deleteTopicsMap {
		processedTopicsArray := util.ReturnDiffTopicList(cMap.TopicMap[cIdx], util.GetAllTopicBySort(deleteTopicsArray))
		if len(processedTopicsArray) == 0 {
			delete(cMap.TopicMap, cIdx)
		} else {
			cMap.TopicMap[cIdx] = processedTopicsArray
		}
		cMap.CollectorPerAgentCnt[cIdx] = len(processedTopicsArray)
	}
	*(cScheduler.inMemoryTopicMap) = cMap
}
/** ### CSP Policy End ### */
