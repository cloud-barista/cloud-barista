package config

import (
	"fmt"
	"os"
	"sync"

	"github.com/spf13/viper"
)

type Config struct {
	InfluxDB       InfluxDB
	CollectManager CollectManager
	APIServer      APIServer
	Monitoring     Monitoring
	Kapacitor      Kapacitor
	Kafka          Kafka
	GrpcServer     GrpcServer
}

type Kapacitor struct {
	EndpointUrl string `json:"endpoint_url" mapstructure:"endpoint_url"`
}

type Kafka struct {
	EndpointUrl         string `json:"endpoint_url" mapstructure:"endpoint_url"`
	ExternalIP          string `json:"external_ip" mapstructure:"external_ip"`
	DeployType          string `json:"deploy_type" mapstructure:"deploy_type"`
	HelmExternalPort    int    `json:"helm_external_port" mapstructure:"helm_external_port"`
	ComposeExternalPort int    `json:"compose_external_port" mapstructure:"compose_external_port"`
	InternalPort        int    `json:"internal_port" mapstructure:"internal_port"`
}

type InfluxDB struct {
	EndpointUrl             string `json:"endpoint_url" mapstructure:"endpoint_url"`
	InternalPort            int    `json:"internal_port" mapstructure:"internal_port"`
	ExternalPort            int    `json:"external_port" mapstructure:"external_port"`
	Database                string
	UserName                string `json:"user_name" mapstructure:"user_name"`
	Password                string
	RetentionPolicyDuration string `json:"rpDuration" mapstructure:"rpDuration"`
}

type CollectManager struct {
	CollectorIP       string `json:"collector_ip" mapstructure:"collector_ip"`
	CollectorGroupCnt int    `json:"collectorGroup_count" mapstructure:"collector_group_count"`
}

type APIServer struct {
	Port int
}

type Monitoring struct {
	AgentInterval           int    `json:"agent_interval" mapstructure:"agent_interval"`         // 모니터링 에이전트 수집주기
	CollectorInterval       int    `json:"collector_interval" mapstructure:"collector_interval"` // 모니터링 콜렉터 Aggregate 주기
	MonitoringPolicy        string `json:"monitoring_policy" mapstructure:"monitoring_policy"`   // 모니터링 콜렉터 정책
	MaxHostCount            int    `json:"max_host_count" mapstructure:"max_host_count"`         // 모니터링 콜렉터 수
	DefaultPolicy           string `json:"default_policy" mapstructure:"default_policy"`         // 모니터링 기본 정책
	PullerInterval          int    `json:"puller_interval" mapstructure:"puller_interval"`       // 모니터링 puller 실행 주기
	PullerAggregateInterval int    `json:"puller_aggregate_interval" mapstructure:"puller_aggregate_interval"`
	AggregateType           string `json:"aggregate_type" mapstructure:"aggregate_type"`
}

type GrpcServer struct {
	Port int
}

func (kapacitor Kapacitor) GetKapacitorEndpointUrl() string {
	return kapacitor.EndpointUrl
}

func (kafka Kafka) GetKafkaEndpointUrl() string {
	return kafka.EndpointUrl
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

func (config *Config) GetInfluxDBConfig() InfluxDB {
	return config.InfluxDB
}

func (config *Config) GetKapacitorConfig() Kapacitor {
	return config.Kapacitor
}

func (config *Config) GetKafkaConfig() Kafka {
	return config.Kafka
}

func (config *Config) GetGrpcConfig() GrpcServer {
	return config.GrpcServer
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
