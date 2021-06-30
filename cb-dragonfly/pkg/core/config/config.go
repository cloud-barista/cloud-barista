package config

import (
	"fmt"
	"net/http"

	"github.com/cloud-barista/cb-dragonfly/pkg/cbstore"
	"github.com/cloud-barista/cb-dragonfly/pkg/config"
	"github.com/cloud-barista/cb-dragonfly/pkg/types"

	"github.com/mitchellh/mapstructure"
)

const (
	MonConfigKey = "/mon/config"
)

// 모니터링 정책 설정
func SetMonConfig(newMonConfig config.Monitoring) (*config.Monitoring, int, error) {
	config.GetInstance().SetMonConfig(newMonConfig)

	var monConfigMap map[string]interface{}
	err := mapstructure.Decode(config.GetInstance().Monitoring, &monConfigMap)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	var defaultMonConfigMap map[string]interface{}
	err = mapstructure.Decode(config.GetDefaultConfig().Monitoring, &defaultMonConfigMap)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	for key, val := range monConfigMap {
		if val == nil || val == 0 || val == "" {
			val = defaultMonConfigMap[key]
		}
		cbstore.GetInstance().StorePut(types.MoNConfig+"/"+key, fmt.Sprintf("%v", val))
	}

	monConfig := config.Monitoring{
		AgentInterval:     cbstore.GetInstance().StoreGetToInt(fmt.Sprintf("%s/%s", types.MoNConfig, "agent_interval")),
		CollectorInterval: cbstore.GetInstance().StoreGetToInt(fmt.Sprintf("%s/%s", types.MoNConfig, "collector_interval")),
		MaxHostCount:      cbstore.GetInstance().StoreGetToInt(fmt.Sprintf("%s/%s", types.MoNConfig, "max_host_count")),
		MonitoringPolicy:  cbstore.GetInstance().StoreGetToString(fmt.Sprintf("%s/%s", types.MoNConfig, "monitoring_policy")),
	}

	return &monConfig, http.StatusOK, nil
}

// 모니터링 정책 조회
func GetMonConfig() (*config.Monitoring, int, error) {

	monConfig := config.Monitoring{
		AgentInterval:     cbstore.GetInstance().StoreGetToInt(fmt.Sprintf("%s/%s", types.MoNConfig, "agent_interval")),
		CollectorInterval: cbstore.GetInstance().StoreGetToInt(fmt.Sprintf("%s/%s", types.MoNConfig, "collector_interval")),
		MaxHostCount:      cbstore.GetInstance().StoreGetToInt(fmt.Sprintf("%s/%s", types.MoNConfig, "max_host_count")),
		MonitoringPolicy:  cbstore.GetInstance().StoreGetToString(fmt.Sprintf("%s/%s", types.MoNConfig, "monitoring_policy")),
	}

	if monConfig.AgentInterval == -1 || monConfig.CollectorInterval == -1 || monConfig.MaxHostCount == -1 || monConfig.MonitoringPolicy == "" {
		return nil, http.StatusInternalServerError, nil
	}

	return &monConfig, http.StatusOK, nil
}

// 모니터링 정책 초기화 co
func ResetMonConfig() (*config.Monitoring, int, error) {
	defaultMonConfig := config.GetDefaultConfig().Monitoring

	var monConfigMap map[string]interface{}
	err := mapstructure.Decode(config.GetDefaultConfig().Monitoring, &monConfigMap)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	mapstructure.Decode(config.GetInstance().Monitoring, &monConfigMap)
	for key, val := range monConfigMap {
		cbstore.GetInstance().StorePut(types.MoNConfig+"/"+key, fmt.Sprintf("%v", val))
	}
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return &defaultMonConfig, http.StatusOK, nil
}
