package mcis

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"time"

	que "github.com/Workiva/go-datastructures/queue"
	"github.com/cloud-barista/cb-dragonfly/pkg/api/core/agent/common"
	"github.com/cloud-barista/cb-dragonfly/pkg/config"
	"github.com/cloud-barista/cb-dragonfly/pkg/storage/cbstore"
	"github.com/cloud-barista/cb-dragonfly/pkg/types"
	"github.com/cloud-barista/cb-dragonfly/pkg/util"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// InMemoryTopic Struct
// TopicMap : 콜렉터 마다 가지고 있는 Topic 현황
// CollectorPerAgentCnt : 콜렉터 마다 가지고 있는 Topic 개수 현황
type InMemoryTopic struct {
	TopicMap             map[int][]string
	CollectorPerAgentCnt []int
}

// CollectorScheduler Struct
// cm : 콜렉터 생성 및 시작, 중지 및 삭제 를 위한 매니저 객체
// inMemoryTopicMap : 토픽 관련 정보를 다루기 위한 inMemoryTopic 구조체
// topicQue : Topic 추가 및 삭제 처리 요청을 담아두는 Queue
type CollectorScheduler struct {
	cm               *CollectManager
	inMemoryTopicMap *InMemoryTopic
	topicQue         *que.Queue
}

func StartScheduler(wg *sync.WaitGroup, manager *CollectManager) error {
	// WaitGroup Start, Initialize Collector
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
	// scheduler 는 토픽 관련 정보들을 in-memory 및 cb-store 로 관리합니다.
	manager.WaitGroup = wg
	c := cbstore.GetInstance()
	// InMemoryTopic Struct
	// TopicMap : 콜렉터 마다 가지고 있는 Topic 현황
	// CollectorPerAgentCnt : 콜렉터 마다 가지고 있는 Topic 개수 현황
	inMemoryTopic := InMemoryTopic{
		TopicMap:             map[int][]string{},
		CollectorPerAgentCnt: []int{},
	}
	// 초기화 작업때, 기존에 생성되어 있는 Topic, TopicMap 을 로드합니다.
	if config.GetInstance().Monitoring.DeployType == types.Helm {
		// Helm 일 경우, configmap 을 통하여 데이터를 로드합니다. (To InMemoryTopic)
		configMap, err := manager.K8sClientSet.CoreV1().ConfigMaps(config.GetInstance().Dragonfly.HelmNamespace).Get(context.TODO(), "cb-dragonfly-collector-configmap", metav1.GetOptions{})
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
		// docker-compose 일 경우, cb-store 를 통하여 데이터를 로드합니다. (To InMemoryTopic)
		_ = c.StoreDelList(types.Topic)
		cPolicy, _ := c.StoreGet(types.CollectorPolicy)
		if cPolicy != nil {
			if *cPolicy == manager.CollectorPolicy {
				getCMapFromStore, _ := c.StoreGet(fmt.Sprintf("%s", types.CollectorTopicMap))
				if getCMapFromStore != nil {
					_ = json.Unmarshal([]byte(*getCMapFromStore), &inMemoryTopic)
					for collectorIdx, topicSlice := range inMemoryTopic.TopicMap {
						for i := 0; i < len(topicSlice); i++ {
							_ = c.StorePut(fmt.Sprintf("%s/%s", types.Topic, topicSlice[i]), strconv.Itoa(collectorIdx))
						}
					}
				}
			}
		}
		_ = c.StorePut(types.CollectorPolicy, manager.CollectorPolicy)
	}
	// CollectorScheduler Struct
	// cm : 콜렉터 생성 및 시작, 중지 및 삭제 를 위한 매니저 객체
	// inMemoryTopicMap : 토픽 관련 정보를 다루기 위한 inMemoryTopic 구조체
	// topicQue : Topic 추가 및 삭제 처리 요청을 담아두는 Queue
	cScheduler := CollectorScheduler{
		cm:               manager,
		inMemoryTopicMap: &inMemoryTopic,
		topicQue:         util.GetRingQueue(), // Global RingBuffer
	}
	return &cScheduler, nil
}

func (cScheduler CollectorScheduler) Scheduler() error {

	interval, _ := cbstore.GetInstance().StoreGet(types.MonConfig + "/" + "mcis_collector_interval")
	aggreTime, _ := strconv.Atoi(*interval)
	topicQue := cScheduler.topicQue
	cPolicy := cScheduler.cm.CollectorPolicy

	for {
		// Aggregate 인터벌을 주기로 계속 수행됩니다.
		time.Sleep(time.Duration(aggreTime) * time.Second)
		// cScheduler.topicQue 에 담겨 있는 Topic 추가, 삭제 처리 요청들을 각각 addTopicList 와 delTopicList 에 담습니다.
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

			addTopicList = util.GetAllTopicBySort(util.Unique(addTopicList, true))
			delTopicList = util.GetAllTopicBySort(util.Unique(util.ReturnDiffTopicList(delTopicList, addTopicList), true))
		}

		curTime := time.Now().Format(time.RFC3339)
		fmt.Printf("[%s] <MCIS> collector scheduler - Now Scheduling ###\n", curTime)
		fmt.Printf("[%s] <MCIS> Add Topics Queue ## : %s\n", curTime, addTopicList)
		fmt.Printf("[%s] <MCIS> Del Topics Queue ## : %s\n", curTime, delTopicList)

		// collector 운용 정책에 따라 addTopicList 와 delTopicList 를 아래 메소드를 통해 수행합니다.
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
// SchedulePolicyBasedCollector
//  === InMemory 연산 Start ===
//  - AddTopicsToCollector: cScheduler.inMemoryTopicMap (콜렉터의 토픽 현황 in-memory 객체) 에 추가 토픽들을 더합니다.
//  - DeleteTopicsToCollector: cScheduler.inMemoryTopicMap 에 삭제 토픽들을 제외합니다.
//  - BalanceTopicsToCollector: cScheduler.inMemoryTopicMap 에 Topic 최적화 배치 작업을 수행합니다.
//  === InMemory 연산 End ===
//  === InMemory 연산 결과값 기반 기능 동작 수행 Start ===
//  - ScaleInOutCollector: cScheduler.inMemoryTopicMap 와 현재 실제로 생성되어 있는 collector 개수를 비교하여 현상 유지(개수 유지), scale in 또는 out 을 수행합니다.
//  - DistributeTopicsToCollector: cScheduler.inMemoryTopicMap 의 topic 설정 값대로 현재 생성되어 있는 콜렉터들에게 Topic 리스트를 전달합니다.
//  - WriteCollectorMapToInMemoryDB: cScheduler.inMemoryTopicMap.TopicMap 을 configmap 또는 cb-store 에 저장합니다.
//  === InMemory 연산 결과값 기반 기능 동작 수행 End ===
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

// AddTopicsToCollector
//  스케줄러의 Topic 관리는 cScheduler.inMemoryTopicMap 와 cb-store 로 이루어 집니다.
// - cScheduler.inMemoryTopicMap
//   > TopicMap => 실제 컬렉터 당 가지고 있는 토픽 맵 리스트
//   > CollectorPerAgentCnt => 컬렉터당 가지고 있는 토픽 개수 array
// - cb-store
//   > types.CollectorTopicMap < /push/collectorTopicMap/{토픽 맵 리스트} > => cScheduler.inMemoryTopicMap.TopicMap 을 고스란히 저장
//   > types.Topic < /push/topic/{토픽} > => 해당 토픽이 배치된 collector idx 값 저장
// 설명 ppt 링크 :
func (cScheduler CollectorScheduler) AddTopicsToCollector(addTopicList []string, maxHostCount int) {
	cMap := *(cScheduler.inMemoryTopicMap)
	c := cbstore.GetInstance()

	for i := 0; i < len(addTopicList); i++ {
		topic := addTopicList[i]
		// cb-store 의 types.Topic < /push/topic/{토픽} > 경로 조회했을 때 이미 등록되어 있을 경우 continue
		topicMessage, _ := c.StoreGet(fmt.Sprintf("%s/%s", types.Topic, topic))

		if topicMessage != nil {
			continue
		}
		// 콜렉터가 아예 생성이 안되어있을때(초기 상태)
		//  - topicMap 에 idx 값이 0 인 콜렉터 생성 (in-memory 상 생성) 및 topic append
		//  - cb-store 경로 /push/topic/{토픽} 에 idx 값 0 쓰기
		//  - CollectorPerAgentCnt 에 토픽 개수 현황 초기화 (array: [1] - 콜렉터 1개, topic 1개)
		if len(cMap.CollectorPerAgentCnt) == 0 {
			cMap.TopicMap[0] = []string{topic}
			_ = c.StorePut(fmt.Sprintf("%s/%s", types.Topic, topic), "0")
			cMap.CollectorPerAgentCnt = []int{1}
			continue
		} else {
			// 콜렉터가 생성되어 있을 때
			//  - CollectorPerAgentCnt 의 콜렉터 당 토픽 개수 현황과 maxHostCount 컨피그 값을 토대로 topic 배치 가능 여부 판단 (in-memory 상 배치)
			needToScaleOut := true
			for collectorIdx, cnt := range cMap.CollectorPerAgentCnt {
				if cnt < maxHostCount {
					// 배치가 가능할 경우, topicMap 의 해당 idx 에 topic 배치
					// cb-store 경로 /push/topic/{토픽} 에 해당 idx 쓰기
					// CollectorPerAgentCnt 에 토픽 개수 현황 업데이트
					cMap.TopicMap[collectorIdx] = append(cMap.TopicMap[collectorIdx], topic)
					_ = c.StorePut(fmt.Sprintf("%s/%s", types.Topic, topic), strconv.Itoa(collectorIdx))
					cMap.CollectorPerAgentCnt[collectorIdx] += 1
					needToScaleOut = false
					break
				}
			}
			// 현재 topicMap 의 모든 콜렉터에 topic 배치가 불가능할 경우(topic 이 꽉찼을 경우) collector 추가 (in-memory 상)
			if needToScaleOut {
				// topicMap 의 해당 idx 에 topic 배치
				// cb-store 경로 /push/topic/{토픽} 에 해당 idx 쓰기
				// CollectorPerAgentCnt 에 토픽 개수 현황 업데이트
				cMap.TopicMap[len(cMap.CollectorPerAgentCnt)] = []string{topic}
				_ = c.StorePut(fmt.Sprintf("%s/%s", types.Topic, topic), strconv.Itoa(len(cMap.CollectorPerAgentCnt)))
				cMap.CollectorPerAgentCnt = append(cMap.CollectorPerAgentCnt, 1)
			}
		}
	}
	*(cScheduler.inMemoryTopicMap) = cMap
}

// DeleteTopicsToCollector
//  - cb-store 의 { key: /push/topic/{토픽}, value: 콜렉터 idx 값 } 토대로 topicMap 에 등록된 topic 을 삭제합니다.
//  - 해당 topic 을 가지고 있는 agent 의 상태를 unhealthy 로 업데이트 합니다.
// 설명 ppt 링크 :
func (cScheduler CollectorScheduler) DeleteTopicsToCollector(delTopicList []string) {
	cMap := *(cScheduler.inMemoryTopicMap)
	if len(cMap.TopicMap) == 0 {
		return
	}
	c := cbstore.GetInstance()

	deleteTopicsMap := map[int][]string{}
	for i := 0; i < len(delTopicList); i++ {
		delTopic := delTopicList[i]
		collectorIdxStr, _ := c.StoreGet(fmt.Sprintf("%s/%s", types.Topic, delTopic))
		collectorIdx, _ := strconv.Atoi(*collectorIdxStr)
		// cb-store 경로 /push/topic/{토픽} 삭제
		_ = c.StoreDelete(fmt.Sprintf("%s/%s", types.Topic, delTopic))

		/* 에이전트 상태 unhealthy 업데이트 처리 요청 Start */
		agentInfo := common.AgentInfo{}
		agentInfoBytes, _ := c.StoreGet(types.Agent + delTopic)
		if agentInfoBytes != nil {
			_ = json.Unmarshal([]byte(*agentInfoBytes), &agentInfo)
			agentInfo.AgentHealth = string(common.Unhealthy)
			agentInfo.AgentState = string(common.Disable)
			recentAgentInfoBytes, _ := json.Marshal(agentInfo)
			_ = c.StorePut(types.Agent+delTopic, string(recentAgentInfoBytes))
		}
		/* 에이전트 상태 unhealthy 업데이트 처리 요청 End */
		// 토픽 삭제가 필요한 콜렉터 map 현황 생성
		deleteTopicsMap[collectorIdx] = append(deleteTopicsMap[collectorIdx], delTopic)
	}
	for cIdx, deleteTopicsArray := range deleteTopicsMap {
		// {현재 topicMap} 과 {토픽 삭제가 필요한 콜렉터 map} 을 비교하여 현재 topicMap 최신화
		processedTopicsArray := util.ReturnDiffTopicList(cMap.TopicMap[cIdx], util.GetAllTopicBySort(deleteTopicsArray))
		cMap.TopicMap[cIdx] = processedTopicsArray
		cMap.CollectorPerAgentCnt[cIdx] = len(cMap.TopicMap[cIdx])
	}
	*(cScheduler.inMemoryTopicMap) = cMap
}

// BalanceTopicsToCollector
// 토픽 추가 & 삭제 처리가 완료된 topicMap 을 기준으로 최적화 배치 수행
func (cScheduler CollectorScheduler) BalanceTopicsToCollector(maxHostCount int) {
	cMap := *(cScheduler.inMemoryTopicMap)

	totalTopicsCnt := 0
	for _, topicCnt := range cMap.CollectorPerAgentCnt {
		totalTopicsCnt += topicCnt
	}
	// example)
	// topicMap 현황: [ 2, 3, 1 ]
	// 전체 토픽 개수: 6
	// 콜렉터당 토픽 개수(환경 설정 값): 5
	idealCollectorCnt := util.CalculateNumberOfCollector(totalTopicsCnt, maxHostCount) // totalTopicsCnt: 6, maxHostCount: 5
	// < idealCollectorCnt: [5, 1] > versus < topicMap 현황: [ 2, 3, 1 ] > Length Comparison
	if len(cMap.TopicMap) == idealCollectorCnt {
		// 두 배열의 비교 크기 값이 같다면 최적화 배치 수행 X
		// 아래 if 문은 init 케이스 때만 수행
		if len(cMap.CollectorPerAgentCnt) == 1 && cMap.CollectorPerAgentCnt[0] == 0 {
			cScheduler.inMemoryTopicMap.CollectorPerAgentCnt = []int{}
			cScheduler.inMemoryTopicMap.TopicMap = map[int][]string{}
		}
		return
	} else {
		// 두 배열의 비교 크기 값이 다르다면 최적화 배치 수행
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

// ScaleInOutCollector
//   - 토픽 추가, 삭제, 최적화 배치가 끝난 topicMap 을 기준으로 collector 개수 동기화
//   - 현재 생성되어 있는 collector 개수 < topicMap collector 개수 => scale out
//   - 현재 생성되어 있는 collector 개수 > topicMap collector 개수 => scale in
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

// DistributeTopicsToCollector
//   - 토픽 추가, 삭제, 최적화 배치가 끝난 topicMap 을 기준으로 각 collector 에게 topic 분배
//   - helm 의 경우 topicMap 현황을 configmap 에 작성
//   - docker-compose 의 경우 topicMap 현황을 go-channel 을 통해 collector(go-routine) 에게 분배
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
			Data:       collectorUUIDMapData,
			BinaryData: topicMapData,
			ObjectMeta: metav1.ObjectMeta{
				Name: types.ConfigMapName,
			}}
		configMapsClient := cScheduler.cm.K8sClientSet.CoreV1().ConfigMaps(config.GetInstance().Dragonfly.HelmNamespace)
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

// WriteCollectorMapToInMemoryDB
//   - 최종 연산이 끝난 topicMap 을 cb-store 에 저장
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
// ScheduleCSPBasedCollector
//   - types.TotalCspCnt 만큼 콜렉터 생성
//   - types 에 정의된 csp 별로 collector 미리 생성
//   - agent 는 topic 및 cspType 정보를 같이 넘겨주며, DF는 agentCspType 별로 collector 에게 분배
//   - 해당 csp 정책의 동작 테스트는 해봐야합니다. 로직 구현은 해놓았으나 정상 동작 유무는 확인이 필요합니다.
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

// ProvisioningCollectorByCSP
// types.TotalCspCnt 만큼 콜렉터 생성
func (cScheduler CollectorScheduler) ProvisioningCollectorByCSP() {
	provisioningCollectorCnt := types.TotalCspCnt - len(cScheduler.cm.CollectorAddrSlice)
	for i := 0; i < provisioningCollectorCnt; i++ {
		if err := cScheduler.cm.CreateCollector(); err != nil {
			fmt.Println("Fail to Create Collector")
		}
	}
}

// AddTopicsToCSPCollector
// csp type 별로 collector 에게 topic 을 추가 (topicMap)
func (cScheduler CollectorScheduler) AddTopicsToCSPCollector(addTopicList []string) {
	cMap := *(cScheduler.inMemoryTopicMap)
	c := cbstore.GetInstance()

	for i := 0; i < len(addTopicList); i++ {
		topic := addTopicList[i]
		topicMsg, _ := c.StoreGet(fmt.Sprintf("%s/%s", types.Topic, topic))
		if topicMsg != nil {
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

// DeleteTopicsToCSPCollector
// csp type 별로 collector 에게 topic 을 삭제 (topicMap)
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

		agentInfo := common.AgentInfo{}
		agentInfoBytes, _ := c.StoreGet(types.Agent + delTopic)
		if agentInfoBytes == nil {
			return
		}
		_ = json.Unmarshal([]byte(*agentInfoBytes), &agentInfo)
		agentInfo.AgentHealth = string(common.Unhealthy)
		agentInfo.AgentState = string(common.Disable)
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
