package config

import (
	"fmt"
	"os"
	"sync"

	"github.com/spf13/viper"
)

type Config struct {
	InfluxDB
	Kapacitor
	Kafka
	Agent
	Dragonfly
	Monitoring
}

type InfluxDB struct {
	EndpointUrl             string `json:"endpoint_url" mapstructure:"endpoint_url"`
	HelmPort                int    `json:"helm_port" mapstructure:"helm_port"`
	Database                string
	UserName                string `json:"user_name" mapstructure:"user_name"`
	Password                string
	RetentionPolicyDuration string `json:"rpDuration" mapstructure:"rpDuration"`
}

type Kapacitor struct {
	EndpointUrl string `json:"endpoint_url" mapstructure:"endpoint_url"`
	HelmPort    int    `json:"helm_port" mapstructure:"helm_port"`
}

type Kafka struct {
	EndpointUrl string `json:"endpoint_url" mapstructure:"endpoint_url"`
	HelmPort    int    `json:"helm_port" mapstructure:"helm_port"`
}

type Dragonfly struct {
	DragonflyIP   string `json:"dragonfly_ip" mapstructure:"dragonfly_ip"`
	Port          int    `json:"port" mapstructure:"port"`
	HelmPort      int    `json:"helm_port" mapstructure:"helm_port"`
	HelmNamespace string `json:"helm_namespace" mapstructure:"helm_namespace"`
}

type Agent struct {
	ServiceAccount string `json:"mck8s_serviceaccount" mapstructure:"mck8s_serviceaccount"` // MCK8S 에이전트 클러스터 시스템 계정
	Namespace      string `json:"mck8s_namespace" mapstructure:"mck8s_namespace"`           // MCK8S 에이전트 클러스터 네임스페이스
	Image          string `json:"image" mapstructure:"image"`
}

type Monitoring struct {
	AgentInterval           int    `json:"agent_interval" mapstructure:"agent_interval"`                     // 모니터링 에이전트 수집주기
	MCISCollectorInterval   int    `json:"mcis_collector_interval" mapstructure:"mcis_collector_interval"`   // MCIS 모니터링 콜렉터 Aggregate 주기
	MCK8SCollectorInterval  int    `json:"mck8s_collector_interval" mapstructure:"mck8s_collector_interval"` // MCK8S 모니터링 콜렉터 Aggregate 주기
	MonitoringPolicy        string `json:"monitoring_policy" mapstructure:"monitoring_policy"`               // 모니터링 콜렉터 정책
	MaxHostCount            int    `json:"max_host_count" mapstructure:"max_host_count"`                     // 모니터링 콜렉터 수
	DefaultPolicy           string `json:"default_policy" mapstructure:"default_policy"`                     // 모니터링 기본 정책
	PullerInterval          int    `json:"puller_interval" mapstructure:"puller_interval"`                   // 모니터링 puller 실행 주기
	PullerAggregateInterval int    `json:"puller_aggregate_interval" mapstructure:"puller_aggregate_interval"`
	AggregateType           string `json:"aggregate_type" mapstructure:"aggregate_type"`
	DeployType              string `json:"deploy_type" mapstructure:"deploy_type"`
}

var once sync.Once
var config Config

func GetInstance() *Config {
	once.Do(func() {
		loadConfigFromYAML(&config)
	})
	return &config
}

func GetDefaultConfig() *Config {
	var defaultMonConfig Config
	loadConfigFromYAML(&defaultMonConfig)
	return &defaultMonConfig
}

func (config *Config) SetMonConfig(newMonConfig Monitoring) {
	config.Monitoring = newMonConfig
}

func (config *Config) GetMonConfig() Monitoring {
	return config.Monitoring
}

func loadConfigFromYAML(config *Config) {
	configPath := os.Getenv("CBMON_ROOT") + "/conf"

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configPath)

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	err = viper.Unmarshal(&config)
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
}
