package config

import (
	"fmt"
	"net/http"

	"github.com/cloud-barista/cb-dragonfly/pkg/config"
	"github.com/cloud-barista/cb-dragonfly/pkg/storage/cbstore"
	"github.com/cloud-barista/cb-dragonfly/pkg/types"

	"github.com/mitchellh/mapstructure"
)

// 모니터링 정책 설정
func SetMonConfig(newMonConfig config.Monitoring) (*config.Monitoring, int, error) {
	config.GetInstance().SetMonConfig(newMonConfig)

	var monConfigMap map[string]interface{}
	err := mapstructure.Decode(config.GetInstance().Monitoring, &monConfigMap)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	for key, val := range monConfigMap {
		if val != nil && val != 0 && val != "" {
			err := cbstore.GetInstance().StorePut(types.MonConfig+"/"+key, fmt.Sprintf("%v", val))
			if err != nil {
				return nil, http.StatusInternalServerError, err
			}
		}
	}

	monConfig := config.Monitoring{
		AgentInterval:           cbstore.GetInstance().StoreGetToInt(fmt.Sprintf("%s/%s", types.MonConfig, "agent_interval")),
		CollectorInterval:       cbstore.GetInstance().StoreGetToInt(fmt.Sprintf("%s/%s", types.MonConfig, "collector_interval")),
		MaxHostCount:            cbstore.GetInstance().StoreGetToInt(fmt.Sprintf("%s/%s", types.MonConfig, "max_host_count")),
		MonitoringPolicy:        cbstore.GetInstance().StoreGetToString(fmt.Sprintf("%s/%s", types.MonConfig, "monitoring_policy")),
		DefaultPolicy:           cbstore.GetInstance().StoreGetToString(fmt.Sprintf("%s/%s", types.MonConfig, "default_policy")),
		PullerInterval:          cbstore.GetInstance().StoreGetToInt(fmt.Sprintf("%s/%s", types.MonConfig, "puller_interval")),
		PullerAggregateInterval: cbstore.GetInstance().StoreGetToInt(fmt.Sprintf("%s/%s", types.MonConfig, "puller_aggregate_interval")),
		AggregateType:           cbstore.GetInstance().StoreGetToString(fmt.Sprintf("%s/%s", types.MonConfig, "aggregate_type")),
		DeployType:              cbstore.GetInstance().StoreGetToString(fmt.Sprintf("%s/%s", types.MonConfig, "deploy_type")),
	}

	return &monConfig, http.StatusOK, nil
}

// 모니터링 정책 조회
func GetMonConfig() (*config.Monitoring, int, error) {
	monConfig := config.Monitoring{
		AgentInterval:           cbstore.GetInstance().StoreGetToInt(fmt.Sprintf("%s/%s", types.MonConfig, "agent_interval")),
		CollectorInterval:       cbstore.GetInstance().StoreGetToInt(fmt.Sprintf("%s/%s", types.MonConfig, "collector_interval")),
		MaxHostCount:            cbstore.GetInstance().StoreGetToInt(fmt.Sprintf("%s/%s", types.MonConfig, "max_host_count")),
		MonitoringPolicy:        cbstore.GetInstance().StoreGetToString(fmt.Sprintf("%s/%s", types.MonConfig, "monitoring_policy")),
		DefaultPolicy:           cbstore.GetInstance().StoreGetToString(fmt.Sprintf("%s/%s", types.MonConfig, "default_policy")),
		PullerInterval:          cbstore.GetInstance().StoreGetToInt(fmt.Sprintf("%s/%s", types.MonConfig, "puller_interval")),
		PullerAggregateInterval: cbstore.GetInstance().StoreGetToInt(fmt.Sprintf("%s/%s", types.MonConfig, "puller_aggregate_interval")),
		AggregateType:           cbstore.GetInstance().StoreGetToString(fmt.Sprintf("%s/%s", types.MonConfig, "aggregate_type")),
		DeployType:              cbstore.GetInstance().StoreGetToString(fmt.Sprintf("%s/%s", types.MonConfig, "deploy_type")),
	}

	if monConfig.AgentInterval == -1 || monConfig.CollectorInterval == -1 || monConfig.MaxHostCount == -1 || monConfig.MonitoringPolicy == "" || monConfig.DefaultPolicy == "" || monConfig.PullerInterval == -1 || monConfig.PullerAggregateInterval == -1 || monConfig.AggregateType == "" || monConfig.DeployType == "" {
		return nil, http.StatusInternalServerError, nil
	}

	return &monConfig, http.StatusOK, nil
}

// 모니터링 정책 초기화 co
func ResetMonConfig() (*config.Monitoring, int, error) {
	defaultMonConfig := config.GetDefaultConfig().Monitoring

	var monConfigMap map[string]interface{}
	err := mapstructure.Decode(defaultMonConfig, &monConfigMap)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	for key, val := range monConfigMap {
		err := cbstore.GetInstance().StorePut(types.MonConfig+"/"+key, fmt.Sprintf("%v", val))
		if err != nil {
			return nil, http.StatusInternalServerError, err
		}
	}

	return &defaultMonConfig, http.StatusOK, nil
}
