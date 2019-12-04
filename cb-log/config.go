// CB-Log: Logger for Cloud-Barista.
//
//      * Cloud-Barista: https://github.com/cloud-barista
//
// load and set config file
//
// ref) https://github.com/go-yaml/yaml/tree/v3
//	https://godoc.org/gopkg.in/yaml.v3
//
// by powerkim@powerkim.co.kr, 2019.08.


package cblog

import (
    "os"
    "io/ioutil"
    "log"

    "gopkg.in/yaml.v3"
)

type CBLOGCONFIG struct {
        CBLOG struct {
                LOOPCHECK bool
                LOGLEVEL string
                LOGFILE bool
        }

        LOGFILEINFO struct {
                FILENAME string
                MAXSIZE int
                MAXBACKUPS int
                MAXAGE int
        }
}

func load(filePath string) ([]byte, error) {
        data, err := ioutil.ReadFile(filePath)
        return data, err
}

func GetConfigInfos() CBLOGCONFIG {
        cblogRootPath := os.Getenv("CBLOG_ROOT")
        data, err := load(cblogRootPath + "/conf/log_conf.yaml")

        if err != nil {
                log.Fatalf("error: %v", err)
        }

        configInfos := CBLOGCONFIG{}
        err = yaml.Unmarshal([]byte(data), &configInfos)
        if err != nil {
                log.Fatalf("error: %v", err)
        }

	return configInfos
}

func GetConfigString(configInfos *CBLOGCONFIG) string {
        d, err := yaml.Marshal(configInfos)
        if err != nil {
                log.Fatalf("error: %v", err)
        }
	return string(d)
}
