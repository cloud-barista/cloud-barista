package mcis

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/Workiva/go-datastructures/queue"
	"github.com/cloud-barista/cb-dragonfly/pkg/config"
	"github.com/cloud-barista/cb-dragonfly/pkg/storage/cbstore"
	"github.com/cloud-barista/cb-dragonfly/pkg/types"
	"github.com/cloud-barista/cb-dragonfly/pkg/util"
)

type InMemoryTopic struct {
	TopicMap map[string][]string // 콜렉터 별 토픽 현황
}

// CollectorScheduler 콜렉터에게 토픽을 분배하는 역할을 담당하는 콜렉터 스케줄러
type CollectorScheduler struct {
	cm               CollectManager
	inMemoryTopicMap InMemoryTopic
	topicQueue       *queue.Queue
}

// StartScheduler MCK8S 콜렉터 스케줄러 구동
func StartScheduler(collectManager CollectManager) error {

	// 콜렉터 스케줄러 생성
	scheduler, err := NewCollectorScheduler(collectManager)
	if err != nil {
		errMsg := fmt.Sprintf("failed to initialize MCK8S collector scheduler, error=%s", err.Error())
		util.GetLogger().Error(errMsg)
		return errors.New(errMsg)
	}

	// 콜렉터 스케줄러 구동
	go func() {
		err = scheduler.DoSchedule()
		if err != nil {
			errMsg := fmt.Sprintf("failed to run goroutine, error=%s", err.Error())
			util.GetLogger().Error(errMsg)
		}
	}()

	return nil
}

// NewCollectorScheduler 콜렉터 스케줄러 생성
func NewCollectorScheduler(cm CollectManager) (*CollectorScheduler, error) {

	inMemoryTopic := InMemoryTopic{TopicMap: map[string][]string{}}

	// 배포 방식에 따라 콜렉터 스케줄러 구동
	deployType := config.GetInstance().Monitoring.DeployType
	// TODO: Helm 모드 개발 및 분리
	if deployType == types.Dev || deployType == types.Compose || deployType == types.Helm {
		// 배포 방식이 개발 모드이거나 도커 모드일 경우

		cbStore := cbstore.GetInstance()

		// 기존에 저장된 개별 토픽 정보 초기화
		_ = cbstore.GetInstance().StoreDelList(types.MCK8STopic)

		// 기존에 저장된 토픽 목록 정보가 있을 경우 동기화
		if topicListData, _ := cbStore.StoreGet(types.MCK8SCollectorTopicMap); topicListData != nil {
			// 콜렉터 정책 검사 (기존 구동 정책과 현재 구동 정책이 동일한 지 확인)
			collectorPolicy, _ := cbstore.GetInstance().StoreGet(types.CollectorPolicy)
			if collectorPolicy != nil {
				if *collectorPolicy == cm.CollectorPolicy {
					// 콜렉터 목록 로드
					err := json.Unmarshal([]byte(*topicListData), &inMemoryTopic)
					if err != nil {
						util.GetLogger().Error("failed to load collector topic map, error=", err.Error())
						return nil, err
					}
					// 개별 토픽 정보 로드
					for key, topicSlice := range inMemoryTopic.TopicMap {
						for i := 0; i < len(topicSlice); i++ {
							_ = cbStore.StorePut(fmt.Sprintf("%s/%s", types.MCK8STopic, topicSlice[i]), key)
						}
					}
				}
			}
			_ = cbstore.GetInstance().StorePut(types.CollectorPolicy, cm.CollectorPolicy)
		}
	}
	//else if deployType == types.Helm {
	//	// 배포 방식이 헬름 모드일 경우
	//	// TODO: 헬름 환경 기반 구동
	//}

	collectorScheduler := &CollectorScheduler{
		cm:               cm,
		inMemoryTopicMap: inMemoryTopic,
		topicQueue:       util.GetMCK8SRingQueue(),
	}
	return collectorScheduler, nil
}

// DoSchedule 콜렉터 스케줄러 구동
func (cScheduler CollectorScheduler) DoSchedule() error {
	interval, _ := cbstore.GetInstance().StoreGet(types.MonConfig + "/" + "mck8s_collector_interval")
	if interval == nil {
		errMsg := "failed to schedule collectors, err: no collector interval configuration data"
		util.GetLogger().Error(errMsg)
		return errors.New(errMsg)
	}
	aggregateInterval, err := strconv.Atoi(*interval)
	if err != nil {
		errMsg := fmt.Sprintf("failed to collector_interval configuration data, error=%s", err.Error())
		util.GetLogger().Error(errMsg)
		return errors.New(errMsg)
	}

	cPolicy := cScheduler.cm.CollectorPolicy

	for {

		// 설정된 스케줄러 주기 기준 동작
		time.Sleep(time.Duration(aggregateInterval) * time.Second)

		// cScheduler.topicQueue 에 담겨 있는 Topic 추가, 삭제 처리 요청들을 각각 addTopicList 와 delTopicList 에 담습니다.
		var addTopicList []string
		var delTopicList []string
		if cScheduler.topicQueue.Len() != 0 {
			topicBytesList, err := cScheduler.topicQueue.Get(cScheduler.topicQueue.Len())
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
		fmt.Printf("[%s] <MCK8S> collector scheduler - Now Scheduling ###\n", curTime)
		fmt.Printf("[%s] <MCK8S> Add Topics Queue ## : %s\n", curTime, addTopicList)
		fmt.Printf("[%s] <MCK8S> Del Topics Queue ## : %s\n", curTime, delTopicList)

		// 콜렉터 구동
		switch cPolicy {
		case types.AgentCntCollectorPolicy:
			cScheduler.SchedulePolicyBasedCollector(addTopicList, delTopicList)
			break
		case types.CSPCollectorPolicy:
			break
		}
	}
}

// SchedulePolicyBasedCollector 쿠버네티스 서비스(MCK8S) 에이전트와 콜렉터를 1:1로 스케줄링
func (cScheduler CollectorScheduler) SchedulePolicyBasedCollector(addTopicList []string, delTopicList []string) {
	provisioningOnce.Do(cScheduler.ProvisioningCollector)

	if len(addTopicList) != 0 {
		cScheduler.AddTopicsToCollector(addTopicList)
	}
	if len(delTopicList) != 0 {
		cScheduler.DeleteTopicsToCollector(delTopicList)
	}
	cScheduler.WriteCollectorMapToInMemoryDB()

	cScheduler.TriggerCollector()
}

// ProvisioningCollector 기존 토픽 맵에 등록된 콜렉터 로드
var provisioningOnce sync.Once

func (cScheduler CollectorScheduler) ProvisioningCollector() {
	for topic, _ := range cScheduler.inMemoryTopicMap.TopicMap {
		// 콜렉터 생성
		err := cScheduler.cm.CreateCollector(topic)
		if err != nil {
			errMsg := fmt.Sprintf("failed to create mck8s collector with topic %s, error=%s", topic, err.Error())
			fmt.Println(errMsg)
			util.GetLogger().Error(errMsg)
			continue
		}
	}
}

func (cScheduler CollectorScheduler) TriggerCollector() {
	for key, _ := range cScheduler.cm.CollectorAddrMap {
		cScheduler.cm.CollectorAddrMap[key].Ch <- key
	}
}

// AddTopicsToCollector 신규 토픽에 대한 콜렉터 생성
func (cScheduler CollectorScheduler) AddTopicsToCollector(addTopicList []string) {
	cbStore := cbstore.GetInstance()
	updatedTopicMap := cScheduler.inMemoryTopicMap

	for i := 0; i < len(addTopicList); i++ {
		topic := addTopicList[i]

		topicMsg, _ := cbStore.StoreGet(fmt.Sprintf("%s/%s", types.MCK8STopic, topic))
		if topicMsg != nil {
			continue
		}

		// 콜렉터 생성
		err := cScheduler.cm.CreateCollector(topic)
		if err != nil {
			errMsg := fmt.Sprintf("failed to create mck8s collector with topic %s, error=%s", topic, err.Error())
			fmt.Println(errMsg)
			util.GetLogger().Error(errMsg)
			continue
		}

		updatedTopicMap.TopicMap[topic] = []string{topic}
	}

	// 토픽 맵 최신화
	cScheduler.inMemoryTopicMap = updatedTopicMap
}

// DeleteTopicsToCollector 삭제 토픽에 대한 콜렉터 삭제
func (cScheduler CollectorScheduler) DeleteTopicsToCollector(delTopicList []string) {
	if len(cScheduler.inMemoryTopicMap.TopicMap) == 0 {
		return
	}

	cbStore := cbstore.GetInstance()
	updatedTopicMap := cScheduler.inMemoryTopicMap
	for i := 0; i < len(delTopicList); i++ {
		topic := delTopicList[i]

		topicMsg, _ := cbStore.StoreGet(fmt.Sprintf("%s/%s", types.MCK8STopic, topic))
		if topicMsg == nil {
			continue
		}

		// 콜렉터 삭제
		err := cScheduler.cm.DeleteCollector(topic)
		if err != nil {
			errMsg := fmt.Sprintf("failed to delete mck8s collector with topic %s, error=%s", topic, err.Error())
			fmt.Println(errMsg)
			util.GetLogger().Error(errMsg)
			continue
		}
		delete(updatedTopicMap.TopicMap, topic)
	}

	// 토픽 맵 최신화
	cScheduler.inMemoryTopicMap = updatedTopicMap
}

func (cScheduler CollectorScheduler) WriteCollectorMapToInMemoryDB() {
	inMemoryTopic := InMemoryTopic{
		TopicMap: cScheduler.inMemoryTopicMap.TopicMap,
	}
	cMapBytes, _ := json.Marshal(inMemoryTopic)
	_ = cbstore.GetInstance().StorePut(fmt.Sprintf("%s", types.MCK8SCollectorTopicMap), string(cMapBytes))

	for _, topic := range inMemoryTopic.TopicMap {
		if err := cbstore.GetInstance().StorePut(fmt.Sprintf("%s/%s", types.MCK8STopic, topic), "0"); err != nil {
			errMsg := fmt.Sprintf("[%s] MCK8S: Failed to save topic data in cbstore, topic: %s, error=%s", time.Now().Format(time.RFC3339), topic, err.Error())
			fmt.Println(errMsg)
			util.GetLogger().Error(errMsg)
			return
		}
	}
}
