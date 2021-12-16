package agent

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cloud-barista/cb-dragonfly/pkg/config"
	"github.com/cloud-barista/cb-dragonfly/pkg/storage/cbstore"
	"github.com/cloud-barista/cb-dragonfly/pkg/types"
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
	ServiceType           string `json:"service_type"`
}

func NewAgentInfo(nsId string, mcisId string, vmId string, cspType string, publicIp string, isHealth bool, serviceType string) AgentInfo {
	state := string(Disable)
	health := string(Unhealthy)
	if isHealth {
		state = string(Enable)
		health = string(Healthy)
	}
	return AgentInfo{
		NsId:                  nsId,
		McisId:                mcisId,
		VmId:                  vmId,
		CspType:               cspType,
		AgentType:             config.GetInstance().Monitoring.DefaultPolicy,
		AgentState:            state,
		AgentHealth:           health,
		AgentUnhealthyRespCnt: 0,
		PublicIp:              publicIp,
		ServiceType:           serviceType,
	}
}

func MakeAgentUUID(nsId string, mcisId string, vmId string, cspType string) string {
	//UUID := types.Agent + "/" + nsId + "/" + mcisId + "/" + vmId + "/" + cspType
	UUID := nsId + "_" + mcisId + "_" + vmId + "_" + cspType
	return UUID
}

// AgentListManager 에이전트 목록 관리

func DeleteAgent(nsId string, mcisId string, vmId string, cspType string) error {
	agentUUID := MakeAgentUUID(nsId, mcisId, vmId, cspType)
	if err := cbstore.GetInstance().StoreDelete(types.Agent + agentUUID); err != nil {
		return err
	}
	return nil
}

func ListAgent() (map[string]AgentInfo, error) {
	agentList := map[string]AgentInfo{}
	agentListByteMap := cbstore.GetInstance().StoreGetListMap(types.Agent, true)

	if len(agentListByteMap) != 0 {
		for uuid, bytes := range agentListByteMap {
			agent := AgentInfo{}
			if err := json.Unmarshal([]byte(bytes), &agent); err != nil {
				return nil, errors.New(fmt.Sprintf("failed to convert agent list, error=%s", err))
			}
			agentList[uuid] = agent
		}
	}
	return agentList, nil
}

func GetAgent(nsId string, mcisId string, vmId string, cspType string) (AgentInfo, error) {
	agentUUID := MakeAgentUUID(nsId, mcisId, vmId, cspType)
	agentInfo := AgentInfo{}
	agentInfoStr := cbstore.GetInstance().StoreGet(fmt.Sprintf(types.Agent + agentUUID))

	if agentInfoStr == "" {
		return agentInfo, errors.New(fmt.Sprintf("failed to get agent with UUID %s", agentUUID))
	}
	err := json.Unmarshal([]byte(agentInfoStr), &agentInfo)
	if err != nil {
		return agentInfo, errors.New(fmt.Sprintf("failed to convert agent info, error=%s", err))
	}
	return agentInfo, nil
}

//func PutAgent(agentUUID string, agentInfo AgentInfo) error {
//nsId string, mcisId string, vmId string, cspType string, publicIp string
func PutAgent(nsId string, mcisId string, vmId string, cspType string, publicIp string, isHealth bool, serviceType string) (string, AgentInfo, error) {
	agentUUID := MakeAgentUUID(nsId, mcisId, vmId, cspType)
	agentInfo := NewAgentInfo(nsId, mcisId, vmId, cspType, publicIp, isHealth, serviceType)
	agentInfoBytes, err := json.Marshal(agentInfo)
	if err != nil {
		return "", AgentInfo{}, errors.New(fmt.Sprintf("failed to convert metadata format to json, error=%s", err))
	}
	err = cbstore.GetInstance().StorePut(types.Agent+agentUUID, string(agentInfoBytes))
	if err != nil {
		return "", AgentInfo{}, errors.New(fmt.Sprintf("failed to put metadata, error=%s", err))
	}
	return agentUUID, agentInfo, nil
}
