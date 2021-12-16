package agent

import (
	"fmt"
	"github.com/cloud-barista/cb-dragonfly/pkg/types"
	"github.com/cloud-barista/cb-dragonfly/pkg/util"
	"net/http"

	"github.com/cloud-barista/cb-dragonfly/pkg/api/core/agent"
	"github.com/cloud-barista/cb-dragonfly/pkg/api/rest"
	"github.com/labstack/echo/v4"
)

type MetaDataListType struct {
	Id agent.AgentInfo `json:"id(ns_id/mcis_id/vm_id/csp_type)"`
}

// ListAgentMetadata 에이전트 메타데이터 조회
// @Summary List agent metadata
// @Description 에이전트 메타데이터 조회
// @Tags [Agent] Monitoring Agent
// @Accept  json
// @Produce  json
// @Param ns query string false "네임스페이스 아이디" Enums(test_ns)
// @Param mcisId query string false "MCIS 아이디" Enums(test_mcis)
// @Param vmId query string false "VM 아이디" Enums(test_vm)
// @Param cspType query string false "VM의 CSP 정보" Enums(aws)
// @Success 200 {object}  rest.JSONResult{[DEFAULT]=[]MetaDataListType,[ID]=AgentInfo} "Different return structures by the given param"
// @Failure 404 {object} rest.SimpleMsg
// @Failure 500 {object} rest.SimpleMsg
// @Router /agent/metadata [get]
func ListAgentMetadata(c echo.Context) error {
	// 에이전트 UUID 파라미터 값 추출

	// 파라미터 값 체크
	agentMetadataList, err := agent.ListAgent()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, rest.SetMessage(fmt.Sprintf("failed to get metadata list, error=%s", err)))
	}
	return c.JSON(http.StatusOK, agentMetadataList)
}

func GetAgentMetadata(c echo.Context) error {
	// 에이전트 UUID 파라미터 값 추출
	nsId := c.Param("ns")
	mcisId := c.Param("mcis_id")
	vmId := c.Param("vm_id")
	cspType := c.Param("csp_type")

	// 파라미터 값 체크
	if nsId == "" || mcisId == "" || vmId == "" || cspType == "" {
		return c.JSON(http.StatusInternalServerError, rest.SetMessage("failed to get metadata"))
	} else {
		agentMetadata, err := agent.GetAgent(nsId, mcisId, vmId, cspType)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, rest.SetMessage(fmt.Sprintf("failed to get metadata, error=%s", err)))
		}
		return c.JSON(http.StatusOK, agentMetadata)
	}
}

func PutAgentMetadata(c echo.Context) error {
	nsId := c.Param("ns")
	mcisId := c.Param("mcis_id")
	vmId := c.Param("vm_id")
	cspType := c.Param("csp_type")
	agentIp := c.Param("agent_ip")

	// 파라미터 값 체크
	if nsId == "" || mcisId == "" || vmId == "" || cspType == "" {
		return c.JSON(http.StatusInternalServerError, rest.SetMessage("failed to update metadata. Check the Params"))
	} else {
		serviceType := "mcis"
		existAgentMetadata, err := agent.GetAgent(nsId, mcisId, vmId, cspType)
		if err != nil {
			serviceType = existAgentMetadata.ServiceType
		}
		agentUUID, agentMetadata, err := agent.PutAgent(nsId, mcisId, vmId, cspType, agentIp, true, serviceType)
		errQue := util.RingQueuePut(types.TopicAdd, agentUUID)
		if err != nil || errQue != nil {
			return c.JSON(http.StatusInternalServerError, rest.SetMessage(fmt.Sprintf("failed to update metadata, error=%s", err)))
		}
		return c.JSON(http.StatusOK, agentMetadata)
	}
}
