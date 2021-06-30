package main

import (
	"errors"
	"fmt"
	"path/filepath"

	cblog "github.com/cloud-barista/cb-log"
	"github.com/sirupsen/logrus"
)

var cblogger *logrus.Logger

func init() {
	// cblog is a global variable.
	filePath := filepath.Join("..", "conf", "log_conf.yaml")
	cblogger = cblog.GetLoggerWithConfigPath("TEST", filePath)
}

func main() {

	fmt.Printf("\n####LogLevel: %s\n", cblog.GetLevel())
	cblogger.Trace("Trace() test: Hello CBLogger")
	cblogger.Tracef("Tracef() test: Hello CBLogger from %s", "Cloud-Barista")

	cblogger.Debug("Debug() test: Hello CBLogger")
	cblogger.Debugf("Debugf() test: Hello CBLogger from %s", "Cloud-Barista")

	cblogger.Info("Info() test: Hello CBLogger")
	cblogger.Infof("Infof() test: Hello CBLogger from %s", "Cloud-Barista")

	cblogger.Warn("Warn() test: Hello CBLogger")
	cblogger.Warnf("Warnf() test: Hello CBLogger from %s", "Cloud-Barista")

	cblogger.Error("Error() test: Hello CBLogger")
	cblogger.Errorf("Errorf() test: Hello CBLogger - %s", getDummyErrorMsg())

	// Fatal and panic should be tested one by one because the two include the exit process.
	cblogger.Fatal("Fatal() test: Hello CBLogger")
	//cblogger.Fatalf("Fatalf() test: Hello CBLogger from %s", "Cloud-Barista")

	//cblogger.Panic("Panic() test: Hello CBLogger")
	//cblogger.Panicf("Panicf() test: Hello CBLogger from %s", "Cloud-Barista")

	fmt.Printf("\n####LogLevel: %s\n", cblog.GetLevel())
	// WithField test
	cblogger.WithField("TestField", "test").Debug("WithField test")
	// WithFields test
	cblogger.WithFields(logrus.Fields{"Field1": "value1", "Field2": "value2", "Field3": "value3"}).Debug("WithFields test")
	// WithError test
	cblogger.WithError(errors.New("Test-error")).Debug("WithError test")

	//	fmt.Printf("===> %#v\n", cblogger.Hooks[0][0])
}

func getDummyErrorMsg() error {
	return fmt.Errorf("internal error message")
}
