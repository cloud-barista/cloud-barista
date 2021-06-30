package metric

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"io/ioutil"
	"net/http"
)

const (
	Rtt  = "Rtt"
	Mrtt = "Mrtt"
)

type MCISMetric struct {
	client    http.Client
	data      CBMCISMetric
	mdata     MCBMCISMetric
	mrequest  Mrequest
	request   Request
	parameter Parameter
}

func GetMCISMonInfo() (interface{}, error) {
	// TODO: MCIS 서비스 모니터링 정보 조회 기능 개발
	return nil, nil
}

func GetMCISRealtimeMonInfo(nsId string, mcisId string) (interface{}, error) {
	// TODO: MCIS 서비스 실시간 모니터링 정보 조회 기능 개발
	return nil, nil
}

// GetMCISCommonMonInfos ...
func GetMCISCommonMonInfo(nsId string, mcisId string, vmId string, agentIp string, metricName string) (*CBMCISMetric, int, error) {
	// MCIS Get 요청 API 생성
	resp, err := http.Get(fmt.Sprintf("http://%s:%d/cb-dragonfly/mcis/metric/%s", agentIp, AgentPort, metricName))
	if err != nil {
		return nil, http.StatusInternalServerError, errors.New("agent server is closed")
	}
	defer resp.Body.Close()

	var metricData CBMCISMetric
	body, err := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		var errmsg map[string]interface{}
		json.Unmarshal(body, &errmsg)
		return nil, resp.StatusCode, errors.New(errmsg["message"].(string))
	}

	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	err = json.Unmarshal(body, &metricData)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	return &metricData, http.StatusOK, nil
}

// GetMCISMonRTTInfo ...
func GetMCISMonRTTInfo(nsId string, mcisId string, vmId string, agentIp string, rttParam Request) (*CBMCISMetric, int, error) {
	// MCIS Get 요청 API 생성
	payload, err := json.Marshal(rttParam)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	req, err := http.NewRequest("GET", fmt.Sprintf("http://%s:%d/cb-dragonfly/mcis/metric/%s", agentIp, AgentPort, Rtt), bytes.NewBuffer(payload))
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	req.Header.Add("Content-Type", "application/json")

	httpClient := http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	defer resp.Body.Close()

	var metricData CBMCISMetric
	body, err := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		var errmsg map[string]interface{}
		json.Unmarshal(body, &errmsg)
		return nil, resp.StatusCode, errors.New(errmsg["message"].(string))
	}

	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	err = json.Unmarshal(body, &metricData)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	return &metricData, http.StatusOK, nil
}

// GetMCISMonMRTTInfo ...
func GetMCISMonMRTTInfo(nsId string, mcisId string, vmId string, agentIp string, mrttParam Mrequest) (*MCBMCISMetric, int, error) {
	// MCIS Get 요청 API 생성
	payload, err := json.Marshal(mrttParam)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	req, err := http.NewRequest("GET", fmt.Sprintf("http://%s:%d/cb-dragonfly/mcis/metric/%s", agentIp, AgentPort, Mrtt), bytes.NewBuffer(payload))
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	req.Header.Add("Content-Type", "application/json")

	httpClient := http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	defer resp.Body.Close()

	var metricData MCBMCISMetric
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	if resp.StatusCode != http.StatusOK {
		var errmsg map[string]interface{}
		json.Unmarshal(body, &errmsg)
		return nil, resp.StatusCode, errors.New(errmsg["message"].(string))
	}

	err = json.Unmarshal(body, &metricData)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	return &metricData, http.StatusOK, nil
}

// GetMCISMonMRTTInfo ...
func (mc *MCISMetric) GetMCISMonMRTTInfo(c echo.Context) error {
	// API Body 데이터 추출
	if err := c.Bind(&mc.mrequest); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	// API 기반 필요 파라미터 추출
	_ = mc.CheckParameter(c)

	// MCIS Get 요청 API 생성
	payload, err := json.Marshal(mc.mrequest)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("http://%s:%d/cb-dragonfly/mcis/metric/%s", mc.parameter.agent_ip, AgentPort, mc.parameter.mcis_metric), bytes.NewBuffer(payload))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	req.Header.Add("Content-Type", "application/json")

	resp, err := mc.client.Do(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	err = json.Unmarshal(body, &mc.mdata)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, &mc.mdata)
}

func (mc *MCISMetric) CheckParameter(c echo.Context) error {
	mc.parameter.agent_ip = c.Param("agent_ip")

	// Query Agent IP 값 체크
	if mc.parameter.agent_ip == "" {
		err := errors.New("No Agent IP in API")
		return c.JSON(http.StatusInternalServerError, err)
	}
	// MCIS 모니터링 메트릭 파라미터 추출
	mc.parameter.mcis_metric = c.Param("mcis_metric_name")
	if mc.parameter.mcis_metric == "" {
		err := errors.New("No Metric Type in API")
		return c.JSON(http.StatusInternalServerError, err)
	}
	return nil
}
