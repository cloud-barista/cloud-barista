package config

import (
	"os"

	"github.com/cloud-barista/cb-ladybug/src/utils/lang"
)

type conf struct {
	RunMode      string
	SpiderUrl    string
	TumblebugUrl string
	BasePath     string
	Username     string
	Password     string
	TargetPath   string
	AppRootPath  string
}

var Config = &conf{}

func Setup() {

	Config.SpiderUrl = lang.NVL(os.Getenv("SPIDER_URL"), "http://localhost:1024/spider")
	Config.TumblebugUrl = lang.NVL(os.Getenv("TUMBLEBUG_URL"), "http://localhost:1323/tumblebug")
	Config.BasePath = lang.NVL(os.Getenv("BASE_PATH"), "/ladybug")
	Config.Username = lang.NVL(os.Getenv("API_USERNAME"), "default")
	Config.Password = lang.NVL(os.Getenv("API_PASSWORD"), "default")
	Config.TargetPath = lang.NVL(os.Getenv("TARGET_PATH"), "/tmp")

	if len(os.Getenv("APP_ROOT")) == 0 {
		if pwd, err := os.Getwd(); err == nil {
			Config.AppRootPath = pwd
		}
	} else {
		Config.AppRootPath = os.Getenv("APP_ROOT")
	}

}
