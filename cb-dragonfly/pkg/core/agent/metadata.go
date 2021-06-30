package agent

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/cloud-barista/cb-dragonfly/pkg/cbstore"
	"github.com/cloud-barista/cb-dragonfly/pkg/config"
)

const (
	AgentListKey = "agentlist"
)

// AgentType 에이전트 동작 메커니즘 유형 (Push, Pull)
type AgentType string

const (
	Push AgentType = "push"
	Pull AgentType = "pull"
)

// AgentState 에이전트 설치 상태 (설치, 제거)
type AgentState string

const (
	Enable  AgentState = "enable"
	Disable AgentState = "disable"
)

// AgentHealth 에이전트 구동 상태 (정상, 비정상)
type AgentHealth string

const (
	Healthy   AgentHealth = "healthy"
	Unhealthy AgentHealth = "unhealthy"
)

// AgentInfo 에이전트 상세 정보
type AgentInfo struct {
	NsId                  string `json:"ns_id"`
	McisId                string `json:"mcis_id"`
	VmId                  string `json:"vm_id"`
	CspType               string `json:"csp_type"`
	AgentType             string `json:"agent_type"`
	AgentState            string `json:"agent_state"`
	AgentHealth           string `json:"agent_health"`
	AgentUnhealthyRespCnt int    `json:"agent_unhealthy_resp_cnt"`
	PublicIp              string `json:"public_ip"`
}

func newAgentInfo(nsId string, mcisId string, vmId string, cspType string, publicIp string) AgentInfo {
	return AgentInfo{
		NsId:                  nsId,
		McisId:                mcisId,
		VmId:                  vmId,
		CspType:               cspType,
		AgentType:             config.GetInstance().Monitoring.DefaultPolicy,
		AgentState:            string(Enable),
		AgentHealth:           string(Healthy),
		AgentUnhealthyRespCnt: 0,
		PublicIp:              publicIp,
	}
}

// AgentListManager 에이전트 목록 관리
type AgentListManager struct{}

func (a AgentListManager) getAgentListFromStore() (map[string]AgentInfo, error) {
	agentList := map[string]AgentInfo{}
	agentListStr := cbstore.GetInstance().StoreGet(AgentListKey)
	if agentListStr != "" {
		err := json.Unmarshal([]byte(agentListStr), &agentList)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("failed to convert agent list, error=%s", err))
		}
	}
	return agentList, nil
}

func (a AgentListManager) putAgentListToStore(agentList map[string]AgentInfo) error {
	agentListBytes, err := json.Marshal(agentList)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to convert agentList format to json, error=%s", err))
	}
	err = cbstore.GetInstance().Store.Put(AgentListKey, string(agentListBytes))
	if err != nil {
		return errors.New(fmt.Sprintf("failed to put agentList, error=%s", err))
	}
	return nil
}

func (a AgentListManager) PutAgent(uuid string, agentInfo AgentInfo) error {
	agentList, err := a.getAgentListFromStore()
	if err != nil {
		return err
	}
	agentList[uuid] = agentInfo
	return a.putAgentListToStore(agentList)
}

func (a AgentListManager) DeleteAgent(uuid string) error {
	agentList, err := a.getAgentListFromStore()
	if err != nil {
		return err
	}

	if _, ok := agentList[uuid]; !ok {
		return errors.New(fmt.Sprintf("failed to update agent, agent with UUID %s not exist", uuid))
	}
	delete(agentList, uuid)

	return a.putAgentListToStore(agentList)
}

func (a AgentListManager) GetAgentList() (map[string]AgentInfo, error) {
	return a.getAgentListFromStore()
}

func (a AgentListManager) GetAgentInfo(uuid string) (AgentInfo, error) {
	agentInfo := AgentInfo{}
	agentInfoStr := cbstore.GetInstance().StoreGet(uuid)

	if agentInfoStr == "" {
		return AgentInfo{}, errors.New(fmt.Sprintf("failed to get agent with UUID %s", uuid))
	}
	err := json.Unmarshal([]byte(agentInfoStr), &agentInfo)
	if err != nil {
		return AgentInfo{}, errors.New(fmt.Sprintf("failed to convert agent info, error=%s", err))
	}
	return agentInfo, nil
}

func PutAgentMetadataToStore(agentUUID string, agentInfo AgentInfo) error {
	// 에이전트 메타데이터 업데이트
	agentInfoBytes, err := json.Marshal(agentInfo)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to convert metadata format to json, error=%s", err))
	}
	err = cbstore.GetInstance().StorePut(agentUUID, string(agentInfoBytes))
	if err != nil {
		return errors.New(fmt.Sprintf("failed to put metadata, error=%s", err))
	}
	// 에이전트 목록 수정
	var agentListManager AgentListManager
	err = agentListManager.PutAgent(agentUUID, agentInfo)
	if err != nil {
		return err
	}
	return nil
}

func SetMetadataByAgentInstall(nsId string, mcisId string, vmId string, cspType string, publicIp string) error {
	agentUUID := MakeAgentUUID(nsId, mcisId, vmId, cspType)
	agentInfo := newAgentInfo(nsId, mcisId, vmId, cspType, publicIp)

	// 에이전트 메타데이터 업데이트
	err := PutAgentMetadataToStore(agentUUID, agentInfo)
	if err != nil {
		return err
	}
	// 에이전트 목록 추가
	var agentListManager AgentListManager
	err = agentListManager.PutAgent(agentUUID, agentInfo)
	if err != nil {
		return err
	}
	return nil
}

func SetMetadataByAgentUninstall(nsId string, mcisId string, vmId string, cspType string) error {
	agentUUID := MakeAgentUUID(nsId, mcisId, vmId, cspType)

	// 에이전트 정보 조회
	var agentListManager AgentListManager
	deletedAgentInfo, err := agentListManager.GetAgentInfo(agentUUID)
	if err != nil {
		return err
	}

	// 에이전트 메타데이터 업데이트 (에이전트 비활성화 처리)
	deletedAgentInfo.AgentState = string(Disable)
	deletedAgentInfo.AgentHealth = string(Unhealthy)
	err = PutAgentMetadataToStore(agentUUID, deletedAgentInfo)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to delete metadata, error=%s", err))
	}
	// 에이전트 목록 수정
	err = agentListManager.PutAgent(agentUUID, deletedAgentInfo)
	if err != nil {
		return err
	}
	return nil
}

func MakeAgentUUID(nsId string, mcisId string, vmId string, cspType string) string {
	UUID := nsId + "/" + mcisId + "/" + vmId + "/" + cspType
	return UUID
}
