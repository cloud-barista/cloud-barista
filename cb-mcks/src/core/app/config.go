package app

import (
	"flag"
	"os"

	"github.com/cloud-barista/cb-mcks/src/utils/lang"
	"github.com/sirupsen/logrus"

	logger "github.com/sirupsen/logrus"
)

type conf struct {
	SpiderCallMethod    *string
	TumblebugCallMethod *string
	SpiderUrl           *string
	TumblebugUrl        *string
	RootURL             *string
	Username            *string
	Password            *string
	AppRootPath         *string
	LoglevelHTTP        *bool
}

var Config *conf

func Setup() {

	var logLevel *string

	Config = &conf{
		AppRootPath:         flag.String("app-root", lang.NVL(os.Getenv("APP_ROOT"), ""), "application root path"),
		RootURL:             flag.String("root-url", lang.NVL(os.Getenv("BASE_URL"), "/mcks"), "root url"),
		SpiderCallMethod:    flag.String("spider-call-method", lang.NVL(os.Getenv("SPIDER_CALL_METHOD"), "REST"), "Method of calling CB-Spider (REST/gRPC)"),
		TumblebugCallMethod: flag.String("tumblebug-call-method", lang.NVL(os.Getenv("TUMBLEBUG_CALL_METHOD"), "REST"), "Method of calling CB-Tumblebug (REST/gRPC)"),
		SpiderUrl:           flag.String("spider-url", lang.NVL(os.Getenv("SPIDER_URL"), "http://localhost:1024/spider"), "cb-spider service end-point url"),
		TumblebugUrl:        flag.String("tumblebug-url", lang.NVL(os.Getenv("TUMBLEBUG_URL"), "http://localhost:1323/tumblebug"), "cb-tumblebug service end-point url"),
		Username:            flag.String("basic-auth-username", lang.NVL(os.Getenv("BASIC_AUTH_USERNAME"), "default"), "rest-api basic auth usernmae"),
		Password:            flag.String("basic-auth-password", lang.NVL(os.Getenv("BASIC_AUTH_PASSWORD"), "default"), "rest-api basic auth password"),
		LoglevelHTTP:        flag.Bool("log-http", os.Getenv("LOG_HTTP") == "true", "The logging http data"),
	}
	logLevel = flag.String("log-level", lang.NVL(os.Getenv("LOG_LEVEL"), "debug"), "The log level")

	flag.Parse()

	//logger
	logger.SetFormatter(&logrus.TextFormatter{
		DisableColors: false,
		FullTimestamp: false,
	})
	//logger.SetReportCaller(true)	//Logging Method Name
	logger.SetOutput(os.Stderr)

	level, err := logrus.ParseLevel(*logLevel)
	if err != nil {
		logger.Fatal(err)
	} else if level != logrus.GetLevel() {
		logger.SetLevel(level)
	} else {
		logger.SetLevel(logrus.InfoLevel)
	}

	// app root path
	if len(*Config.AppRootPath) == 0 {
		if pwd, err := os.Getwd(); err == nil {
			Config.AppRootPath = &pwd
		}
	}

}
