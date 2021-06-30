// CB-Log: Logger for Cloud-Barista.
//
//      * Cloud-Barista: https://github.com/cloud-barista
//
// load and set config file
//
// ref) https://github.com/go-yaml/yaml/tree/v3
//	https://godoc.org/gopkg.in/yaml.v3
//
// by CB-Log Team, 2019.08.

package cblog

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type CBLOGCONFIG struct {
	CBLOG struct {
		LOOPCHECK bool
		LOGLEVEL  string
		LOGFILE   bool
	}

	LOGFILEINFO struct {
		FILENAME   string
		MAXSIZE    int
		MAXBACKUPS int
		MAXAGE     int
	}
}

func load(filePath string) ([]byte, error) {
	data, err := ioutil.ReadFile(filePath)
	return data, err
}

func GetConfigInfos(configFilePath string) CBLOGCONFIG {
	var filePath string
	cblogRootPath := os.Getenv("CBLOG_ROOT")

	if cblogRootPath == "" && configFilePath == "" {
		log.Fatalf("Both $CBLOG_ROOT and configPath are not set!!")
	}

	if cblogRootPath != "" {
		filePath = filepath.Join(cblogRootPath, "conf", "log_conf.yaml")
	} else {
		filePath = configFilePath
	}

	data, err := load(filePath)

	if err != nil {
		log.Fatalf("error: %v", err)
	}

	configInfos := CBLOGCONFIG{}
	err = yaml.Unmarshal([]byte(data), &configInfos)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	configInfos.LOGFILEINFO.FILENAME = ReplaceEnvPath(configInfos.LOGFILEINFO.FILENAME)
	return configInfos
}

// $ABC/def ==> /abc/def
func ReplaceEnvPath(str string) string {
	if strings.Index(str, "$") == -1 {
		return str
	}

	// ex) input "$CBSTORE_ROOT/meta_db/dat"
	strList := strings.Split(str, "/")
	for n, one := range strList {
		if strings.Index(one, "$") != -1 {
			cbstoreRootPath := os.Getenv(strings.Trim(one, "$"))
			if cbstoreRootPath == "" {
				log.Fatal(one + " is not set!")
			}
			strList[n] = cbstoreRootPath
		}
	}

	var resultStr string
	for _, one := range strList {
		resultStr = resultStr + one + "/"
	}
	// ex) "/root/go/src/github.com/cloud-barista/cb-spider/meta_db/dat/"
	resultStr = strings.TrimRight(resultStr, "/")
	resultStr = strings.ReplaceAll(resultStr, "//", "/")
	return resultStr
}

func GetConfigString(configInfos *CBLOGCONFIG) string {
	d, err := yaml.Marshal(configInfos)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	return string(d)
}
