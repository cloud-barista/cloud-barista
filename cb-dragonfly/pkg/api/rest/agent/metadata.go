package agent

import (
	"fmt"
	"net/http"

	"github.com/cloud-barista/cb-dragonfly/pkg/api/core/agent/common"
	"github.com/cloud-barista/cb-dragonfly/pkg/api/rest"
	"github.com/cloud-barista/cb-dragonfly/pkg/types"
	"github.com/cloud-barista/cb-dragonfly/pkg/util"
	"github.com/labstack/echo/v4"
)

type MetaDataListType struct {
	Id common.AgentInfo `json:"id(ns_id/mcis_id/vm_id/csp_type)"`
}

// ListAgentMetadata 에이전트 메타데이터 조회
// @Summary List Agent Metadata
// @Description 에이전트 메타데이터 목록 조회
// @Tags [Agent] Monitoring Agent
// @Accept  json
// @Produce  json
// @Success 200 {object}  rest.JSONResult{[DEFAULT]=[]MetaDataListType,[ID]=common.AgentInfo} "Different return structures by the given param"
// @Failure 404 {object} rest.SimpleMsg
// @Failure 500 {object} rest.SimpleMsg
// @Router /agents/metadata [get]
func ListAgentMetadata(c echo.Context) error {
	// 에이전트 UUID 파라미터 값 추출

	// 파라미터 값 체크
	agentMetadataList, err := common.ListAgent()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, rest.SetMessage(fmt.Sprintf("failed to get metadata list, error=%s", err)))
	}
	return c.JSON(http.StatusOK, agentMetadataList)
}

// GetAgentMetadata 에이전트 메타데이터 조회
// @Summary Get Agent Metadata
// @Description 에이전트 메타데이터 단일 조회
// @Tags [Agent] Monitoring Agent
// @Accept  json
// @Produce  json
// @Param ns_id query string false "네임스페이스 아이디" Enums(test_ns)
// @Param service_type query string false "서비스 타입" Enums(mcis)
// @Param service_id query string false "서비스 아이디" Enums(mcis_id)
// @Success 200 {object}  rest.JSONResult{[DEFAULT]=[]MetaDataListType,[ID]=common.AgentInfo} "Different return structures by the given param"
// @Failure 404 {object} rest.SimpleMsg
// @Failure 500 {object} rest.SimpleMsg
// @Router /agent/metadata [get]
func GetAgentMetadata(c echo.Context) error {
	// 에이전트 UUID 파라미터 값 추출
	nsId := c.QueryParam("ns")
	serviceType := c.QueryParam("service_type")
	serviceId := c.QueryParam("service_id")
	var vmId, cspType string

	if !checkEmptyFormParam(serviceType) {
		return c.JSON(http.StatusBadRequest, rest.SetMessage("empty agent type parameter"))
	}

	requestInfo := common.AgentInstallInfo{
		ServiceType: serviceType,
		NsId:        nsId,
	}

	if util.CheckMCK8SType(serviceType) {
		if !checkEmptyFormParam(nsId, serviceId) {
			return c.JSON(http.StatusBadRequest, rest.SetMessage("bad request parameter to get mck8s agent metadata"))
		}
		requestInfo.Mck8sId = serviceId
	}

	if util.CheckMCISType(serviceType) {
		vmId = c.QueryParam("vm_id")
		cspType = c.QueryParam("csp_type")
		if !checkEmptyFormParam(nsId, serviceId, vmId, cspType) {
			return c.JSON(http.StatusBadRequest, rest.SetMessage("bad request parameter to get mcis agent metadata"))
		}
		requestInfo.Mck8sId = serviceId
		requestInfo.VmId = vmId
		requestInfo.CspType = cspType
	}

	agentMetadata, err := common.GetAgent(requestInfo)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, rest.SetMessage(fmt.Sprintf("failed to get metadata, error=%s", err)))
	}
	return c.JSON(http.StatusOK, agentMetadata)
}

// PutAgentMetadata 에이전트 메타데이터 수정
// @Summary Put Agent Metadata
// @Description 에이전트 메타데이터 수정
// @Tags [Agent] Monitoring Agent
// @Accept  json
// @Produce  json
// @Param ns query string false "네임스페이스 아이디" Enums(test_ns)
// @Param mcisId query string false "MCIS 아이디" Enums(test_mcis)
// @Param vmId query string false "VM 아이디" Enums(test_vm)
// @Param cspType query string false "VM의 CSP 정보" Enums(aws)
// @Param mck8sId query string false "MCK8S 아이디" Enums(test_mck8s)
// @Success 200 {object}  rest.JSONResult{[DEFAULT]=[]MetaDataListType,[ID]=common.AgentInfo} "Different return structures by the given param"
// @Failure 404 {object} rest.SimpleMsg
// @Failure 500 {object} rest.SimpleMsg
// @Router /agent/metadata [put]
func PutAgentMetadata(c echo.Context) error {
	// 에이전트 UUID 파라미터 값 추출
	params := &rest.AgentType{}
	if err := c.Bind(params); err != nil {
		return err
	}

	if !checkEmptyFormParam(params.ServiceType) {
		return c.JSON(http.StatusBadRequest, rest.SetMessage("empty agent type parameter"))
	}

	if util.CheckMCK8SType(params.ServiceType) {
		// 토큰 값이 비어있을 경우
		if !checkEmptyFormParam(params.NsId, params.Mck8sId) {
			return c.JSON(http.StatusBadRequest, rest.SetMessage("bad request parameter to update mck8s agent metadata"))
		}
	}
	// MCIS 에이전트 form 파라미터 값 체크
	if util.CheckMCISType(params.ServiceType) {
		// MCIS 에이전트 form 파라미터 값 체크
		if !checkEmptyFormParam(params.NsId, params.McisId, params.VmId, params.CspType, params.PublicIp) {
			return c.JSON(http.StatusBadRequest, rest.SetMessage("bad request parameter to update mcis agent metadata"))
		}
	}

	requestInfo := common.AgentInstallInfo{
		ServiceType: params.ServiceType,
		NsId:        params.NsId,
		McisId:      params.McisId,
		VmId:        params.VmId,
		PublicIp:    params.PublicIp,
		CspType:     params.CspType,
		Mck8sId:     params.Mck8sId,
	}

	// 메타데이터 조회
	existAgentMetadata, err := common.GetAgent(requestInfo)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, rest.SetMessage(fmt.Sprintf("failed to get metadata before update metadata, error=%s", err)))
	}

	// 메타데이터 수정
	agentUUID, agentMetadata, err := common.PutAgent(requestInfo,
		existAgentMetadata.AgentUnhealthyRespCnt,
		common.AgentState(existAgentMetadata.AgentState),
		common.AgentHealth(existAgentMetadata.AgentHealth))

	errQue := util.RingQueuePut(types.TopicAdd, agentUUID)
	if err != nil || errQue != nil {
		return c.JSON(http.StatusInternalServerError, rest.SetMessage(fmt.Sprintf("failed to update metadata, error=%s", err)))
	}
	return c.JSON(http.StatusOK, agentMetadata)
}
