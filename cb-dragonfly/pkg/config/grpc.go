package config

import (
	"fmt"
	"os"
	"sync"

	"github.com/spf13/viper"
)

type GrpcConfig struct {
	GrpcServer GrpcServer
}

type GrpcServer struct {
	Ip   string
	Port int
}

var gOnce sync.Once
var gConfig GrpcConfig

func GetGrpcInstance() *GrpcConfig {
	gOnce.Do(func() {
		loadGrpcConfigFromYAML(&gConfig)
		if gConfig.GrpcServer.Ip == "" {
			gConfig.GrpcServer.Ip = "0.0.0.0"
		}
	})
	return &gConfig
}

func loadGrpcConfigFromYAML(config *GrpcConfig) {
	configPath := os.Getenv("CBMON_ROOT") + "/conf"

	viper.SetConfigName("grpc_conf")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configPath)

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error g file: %s \n", err))
	}
	err = viper.Unmarshal(&config)
	if err != nil {
		panic(fmt.Errorf("Fatal error g file: %s \n", err))
	}
}
