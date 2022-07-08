package mcis

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cloud-barista/cb-dragonfly/pkg/types"
	"io/ioutil"
	"net/http"
)

const (
	Rtt  = "Rtt"
	Mrtt = "Mrtt"
)

type MCISMetric struct {
	client    http.Client
	data      types.CBMCISMetric
	mdata     types.MCBMCISMetric
	mrequest  types.Mrequest
	request   types.Request
	parameter types.Parameter
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
func GetMCISCommonMonInfo(nsId string, mcisId string, vmId string, agentIp string, metricName string) (*types.CBMCISMetric, int, error) {
	// MCIS Get 요청 API 생성
	resp, err := http.Get(fmt.Sprintf("http://%s:%d/cb-dragonfly/mcis/metric/%s", agentIp, AgentPort, metricName))
	if err != nil {
		return nil, http.StatusInternalServerError, errors.New("agent server is closed")
	}
	defer resp.Body.Close()

	var metricData types.CBMCISMetric
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
func GetMCISMonRTTInfo(nsId string, mcisId string, vmId string, agentIp string, rttParam types.Request) (*types.CBMCISMetric, int, error) {
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

	var metricData types.CBMCISMetric
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
func GetMCISMonMRTTInfo(nsId string, mcisId string, vmId string, agentIp string, mrttParam types.Mrequest) (*types.MCBMCISMetric, int, error) {
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

	var metricData types.MCBMCISMetric
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
