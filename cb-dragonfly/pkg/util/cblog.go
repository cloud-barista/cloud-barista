package util

import (
	"sync"

	cblog "github.com/cloud-barista/cb-log"
	"github.com/sirupsen/logrus"
)

var once sync.Once
var cbLogger *logrus.Logger

func GetLogger() *logrus.Logger {
	once.Do(func() {
		cbLogger = cblog.GetLogger("CB-DRAGONFLY")
	})
	return cbLogger
}
