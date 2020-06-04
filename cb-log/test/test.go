package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/cloud-barista/cb-log"
)

var cblogger *logrus.Logger

func init() {
	// cblog is a global variable.
	cblogger = cblog.GetLogger("CB-SPIDER")
	cblog.SetLevel("error")
	cblog.SetLevel("warn")
	cblog.SetLevel("info")
}

func main() {

	fmt.Printf("\n####LogLevel: %s\n", cblog.GetLevel())
	cblogger.Info("Log Info message")
	cblogger.Warning("Log Waring message")
	cblogger.Error("Log Error message")
	cblogger.Errorf("Log Error message:%s", errorMsg())

	cblog.SetLevel("warn")
	fmt.Printf("\n####LogLevel: %s\n", cblog.GetLevel())
	cblogger.Info("Log Info message")
	cblogger.Warning("Log Waring message")
	cblogger.Error("Log Error message")
	cblogger.Errorf("Log Error message:%s", errorMsg())

	cblog.SetLevel("error")
	fmt.Printf("\n####LogLevel: %s\n", cblog.GetLevel())
	cblogger.Info("Log Info message")
	cblogger.Warning("Log Waring message")
	cblogger.Error("Log Error message")
	cblogger.Errorf("Log Error message:%s", errorMsg())
	
//	fmt.Printf("===> %#v\n", cblogger.Hooks[0][0])
}

func errorMsg() error {
	return fmt.Errorf("internal error message")
}
