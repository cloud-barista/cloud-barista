package agent

import (
	"fmt"
	"net/http"

	"github.com/cloud-barista/cb-dragonfly/pkg/api/rest"
	"github.com/cloud-barista/cb-dragonfly/pkg/core/agent"
	"github.com/labstack/echo/v4"
)

var agentListManager agent.AgentListManager

func ListAgentMetadata(c echo.Context) error {
	// 에이전트 UUID 파라미터 값 추출
	nsId := c.QueryParam("ns")
	mcisId := c.QueryParam("mcisId")
	vmId := c.QueryParam("vmId")
	cspType := c.QueryParam("cspType")

	// 파라미터 값 체크
	if nsId == "" || mcisId == "" || vmId == "" || cspType == "" {
		agentMetadataList, err := agentListManager.GetAgentList()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, rest.SetMessage(fmt.Sprintf("failed to get metadata list, error=%s", err)))
		}
		return c.JSON(http.StatusOK, agentMetadataList)
	} else {
		agentUUID := agent.MakeAgentUUID(nsId, mcisId, vmId, cspType)
		agentMetadata, err := agentListManager.GetAgentInfo(agentUUID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, rest.SetMessage(fmt.Sprintf("failed to get metadata, error=%s", err)))
		}
		return c.JSON(http.StatusOK, agentMetadata)
	}
}
