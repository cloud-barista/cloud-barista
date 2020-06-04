package main

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/cloud-barista/cb-log"
)

var cblogger *logrus.Logger

func init() {
	// cblog is a global variable.
	cblogger = cblog.GetLogger("CB-SPIDER")
}

func main() {

	for {
		fmt.Printf("\n####LogLevel: %s\n", cblog.GetLevel())
		cblogger.Info("Log Info message")
		cblogger.Warning("Log Waring message")
		cblogger.Error("Log Error message")
		cblogger.Errorf("Log Error message:%s", errorMsg())
		time.Sleep(time.Second*2)
	}
}

func errorMsg() error {
	return fmt.Errorf("internal error message")
}
