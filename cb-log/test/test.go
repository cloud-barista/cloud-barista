package main

import (
	"errors"
	"fmt"

	cblog "github.com/cloud-barista/cb-log"
	"github.com/sirupsen/logrus"
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

	cblog.SetLevel("debug")
	fmt.Printf("\n####LogLevel: %s\n", cblog.GetLevel())
	// WithField 테스트
	cblogger.WithField("TestField", "test").Debug("WithField 테스트")
	// WithFields 테스트
	cblogger.WithFields(logrus.Fields{"Field1": "value1", "Field2": "value2", "Field3": "value3"}).Debug("WithFields 테스트")
	// WithError 테스트
	cblogger.WithError(errors.New("테스트 오류")).Debug("WithError 테스트")

	//	fmt.Printf("===> %#v\n", cblogger.Hooks[0][0])
}

func errorMsg() error {
	return fmt.Errorf("internal error message")
}
