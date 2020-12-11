package localstore

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cloud-barista/cb-dragonfly/pkg/api/rest"
	"github.com/labstack/echo/v4"
	"net/http"
)

const (
	ENABLE  = "Enable"
	DISABLE = "Disable"
)

type Metadata struct {
	Key   *string
	Value *AgentInfo
}

type AgentInfo struct {
	AgentState *string
	AgentType  *string
	Public_IP  *string
}

var metadata = &Metadata{}
var agentinfo = &AgentInfo{}

func setAgentInfo(agentstate *string, agenttype *string, public_ip *string) {
	getAgentInfo().AgentState = agentstate
	getAgentInfo().AgentType = agenttype
	getAgentInfo().Public_IP = public_ip
}
func getAgentInfo() *AgentInfo {
	return agentinfo
}
func setMetadata(uuid *string, agentinfo *AgentInfo) {
	getMetadata().Key = uuid
	getMetadata().Value = agentinfo
}

func getMetadata() *Metadata {
	return metadata
}

func AgentInstallationMetadata(nsId string, mcisId string, vmId string, cspType string, publicIp string) error {
	agent_State := ENABLE
	agent_Type := ENABLE
	setAgentInfo(&agent_State, &agent_Type, &publicIp)
	setMetadata(makeAgentUUID(nsId, mcisId, vmId, cspType), getAgentInfo())
	data, err := json.Marshal(getMetadata().Value)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to convert metadata format to json, error=%s", err))
	}

	err = GetInstance().StorePut(*getMetadata().Key, string(data))
	if err != nil {
		return errors.New(fmt.Sprintf("failed to put metadata, error=%s", err))
	}
	return nil
}

func AgentDeletionMetadata(nsId string, mcisId string, vmId string, cspType string, publicIp string) error {
	agent_State := DISABLE
	agent_Type := DISABLE
	setAgentInfo(&agent_State, &agent_Type, &publicIp)
	setMetadata(makeAgentUUID(nsId, mcisId, vmId, cspType), getAgentInfo())

	_, err := json.Marshal(getMetadata().Value)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to convert metadata format to json, error=%s", err))
	}

	err = GetInstance().StoreDelete(*getMetadata().Key)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to delete metadata, error=%s", err))
	}
	return nil
}

func makeAgentUUID(nsId string, mcisId string, vmId string, cspType string) *string {
	data := fmt.Sprintf(nsId + "/" + mcisId + "/" + vmId + "/" + cspType)
	return &data
}

func ShowMetadata(c echo.Context) error {
	//온디멘드 모니터링 Agent IP 파라미터 추출
	ns_id := c.Param("ns")
	mcis_id := c.Param("mcis_id")
	vm_id := c.Param("vm_id")
	csp_type := c.Param("csp_type")

	// Query 파라미터  값 체크
	if ns_id == "" || mcis_id == "" || vm_id == "" || csp_type == "" {
		return c.JSON(http.StatusInternalServerError, rest.SetMessage("failed to get package. query parameter is missing"))
	}
	value, err := GetInstance().Store.Get(fmt.Sprintf(ns_id + "/" + mcis_id + "/" + vm_id + "/" + csp_type))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errors.New(fmt.Sprintf("Get Data from CB-Store Error, err=%s", err)))
	}

	return c.JSON(http.StatusOK, value)
}
