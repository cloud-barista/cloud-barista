package common

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/cloud-barista/cb-dragonfly/pkg/config"
	"github.com/cloud-barista/cb-dragonfly/pkg/storage/cbstore"
	"github.com/cloud-barista/cb-dragonfly/pkg/types"
	"github.com/cloud-barista/cb-dragonfly/pkg/util"
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
	ServiceType           string `json:"service_type"`
	NsId                  string `json:"ns_id"`
	McisId                string `json:"mcis_id"`
	VmId                  string `json:"vm_id"`
	CspType               string `json:"cspType"`
	AgentType             string `json:"agent_type"`
	AgentState            string `json:"agent_state"`
	AgentHealth           string `json:"agent_health"`
	AgentUnhealthyRespCnt int    `json:"agent_unhealthy_resp_cnt"`
	PublicIp              string `json:"public_ip"`
	Mck8sId               string `json:"mck8s_id"`
}

func MakeAgentUUID(info AgentInstallInfo) string {
	if util.CheckMCK8SType(info.ServiceType) {
		return fmt.Sprintf("%s_%s_%s", info.NsId, info.ServiceType, info.Mck8sId)
	}
	return fmt.Sprintf("%s_%s_%s_%s_%s", info.NsId, info.ServiceType, info.McisId, info.VmId, info.CspType)
}

// DeleteAgent 에이전트 메타데이터 삭제
func DeleteAgent(info AgentInstallInfo) (string, error) {
	agentUUID := MakeAgentUUID(info)
	if err := cbstore.GetInstance().StoreDelete(types.Agent + agentUUID); err != nil {
		return agentUUID, err
	}
	return agentUUID, nil
}

// DeleteAgentByUUID UUID 기준 에이전트 메타데이터 삭제
func DeleteAgentByUUID(agentUUID string) error {
	if err := cbstore.GetInstance().StoreDelete(types.Agent + agentUUID); err != nil {
		return err
	}
	return nil
}

// ListAgent 에이전트 메타데이터 목록 조회
func ListAgent() (map[string]AgentInfo, error) {
	agentList := map[string]AgentInfo{}
	agentListByteMap, err := cbstore.GetInstance().StoreGetListMap(types.Agent, true)
	if err != nil {
		return nil, err
	}

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

func GetAgent(info AgentInstallInfo) (*AgentInfo, error) {
	agentUUID := MakeAgentUUID(info)
	agentInfo := AgentInfo{}
	agentInfoStr, err := cbstore.GetInstance().StoreGet(fmt.Sprintf(types.Agent + agentUUID))
	if err != nil {
		return nil, err
	}

	if agentInfoStr == nil {
		return nil, errors.New(fmt.Sprintf("failed to get agent with UUID %s", agentUUID))
	}
	if err = json.Unmarshal([]byte(*agentInfoStr), &agentInfo); err != nil {
		return nil, errors.New(fmt.Sprintf("failed to convert agent info, error=%s", err))
	}
	return &agentInfo, nil
}

// GetAgentByUUID UUID 기준 에이전트 메타데이터 조회
func GetAgentByUUID(agentUUID string) (*AgentInfo, error) {
	agentInfo := AgentInfo{}
	agentInfoStr, err := cbstore.GetInstance().StoreGet(fmt.Sprintf(types.Agent + agentUUID))
	if err != nil {
		return nil, err
	}

	if agentInfoStr == nil {
		return nil, errors.New(fmt.Sprintf("failed to get agent with UUID %s", agentUUID))
	}
	if err = json.Unmarshal([]byte(*agentInfoStr), &agentInfo); err != nil {
		return nil, errors.New(fmt.Sprintf("failed to convert agent info, error=%s", err))
	}
	return &agentInfo, nil
}

// PutAgent 에이전트 메타데이터 수정
func PutAgent(info AgentInstallInfo, unHealthyRespCnt int, agentState AgentState, agentHealth AgentHealth) (string, AgentInfo, error) {
	agentUUID := MakeAgentUUID(info)
	agentInfo := AgentInfo{}
	if util.CheckMCK8SType(info.ServiceType) {
		// 에이전트 메타데이터 기본 정보 설정
		agentInfo = AgentInfo{
			ServiceType: info.ServiceType,
			NsId:        info.NsId,
			Mck8sId:     info.Mck8sId,
			AgentType:   config.GetInstance().Monitoring.DefaultPolicy,
			//AgentUnhealthyRespCnt: unHealthyRespCnt,
			//AgentState:            string(agentState),
			//AgentHealth:           string(agentHealth),
		}
	} else {
		agentInfo = AgentInfo{
			ServiceType: info.ServiceType,
			NsId:        info.NsId,
			McisId:      info.McisId,
			VmId:        info.VmId,
			CspType:     info.CspType,
			AgentType:   config.GetInstance().Monitoring.DefaultPolicy,
			PublicIp:    info.PublicIp,
			//AgentUnhealthyRespCnt: unHealthyRespCnt,
			//AgentState:            string(agentState),
			//AgentHealth:           string(agentHealth),
		}
	}

	// 에이전트 비정상 횟수, 설치 상태 정보, 헬스 정보 설정
	if unHealthyRespCnt != -1 {
		agentInfo.AgentUnhealthyRespCnt = unHealthyRespCnt
	}
	if agentState != "" {
		agentInfo.AgentState = string(agentState)
	}
	if agentHealth != "" {
		agentInfo.AgentHealth = string(agentHealth)
	}

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
