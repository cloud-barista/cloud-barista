// CB-Log: Logger for Cloud-Barista.
//
//      * Cloud-Barista: https://github.com/cloud-barista
//
// ref) https://github.com/sirupsen/logrus
// ref) https://github.com/natefinch/lumberjack
// ref) https://github.com/snowzach/rotatefilehook
// by CB-Log Team, 2019.08.

package cblog


import (
	"os"
	"fmt"
	"time"
	"strings"

        "github.com/sirupsen/logrus"
	"github.com/snowzach/rotatefilehook"
	"github.com/cloud-barista/cb-log/formatter"
)

type CBLogger struct {
	loggerName string
	logrus *logrus.Logger
}

// global var.
var (
	thisLogger *CBLogger
	thisFormatter *cblogformatter.Formatter
	cblogConfig CBLOGCONFIG
)

// You can set up with Framework Name, a Framework Name is one of loggerName.
func GetLogger(loggerName string) *logrus.Logger {
	if thisLogger != nil {
		return thisLogger.logrus
	}
	thisLogger = new(CBLogger)
	thisLogger.loggerName = loggerName
	thisLogger.logrus =  &logrus.Logger{
        Level: logrus.DebugLevel,
        Out:   os.Stderr,
        Hooks: make(logrus.LevelHooks),
        Formatter: getFormatter(loggerName),
	}

	// set config.
	setup(loggerName)
	return thisLogger.logrus
}

func setup(loggerName string) {
	cblogConfig = GetConfigInfos()
	thisLogger.logrus.SetReportCaller(true)

	if cblogConfig.CBLOG.LOOPCHECK {
		SetLevel(cblogConfig.CBLOG.LOGLEVEL)
		go levelSetupLoop(loggerName)
	} else {
		SetLevel(cblogConfig.CBLOG.LOGLEVEL)
	}

	if cblogConfig.CBLOG.LOGFILE {
		setRotateFileHook(loggerName, &cblogConfig)
	}
}

// Now, this method is busy wait. 
// @TODO must change this  with file watch&event.
// ref) https://github.com/fsnotify/fsnotify/blob/master/example_test.go
func levelSetupLoop(loggerName string) {
	for {
		cblogConfig = GetConfigInfos()
		SetLevel(cblogConfig.CBLOG.LOGLEVEL)
		time.Sleep(time.Second*2)
	}
}

func setRotateFileHook(loggerName string, logConfig *CBLOGCONFIG) {
	level, _ := logrus.ParseLevel(logConfig.CBLOG.LOGLEVEL)

        rotateFileHook, err := rotatefilehook.NewRotateFileHook(rotatefilehook.RotateFileConfig{
                Filename:   logConfig.LOGFILEINFO.FILENAME,
                MaxSize:    logConfig.LOGFILEINFO.MAXSIZE, // megabytes
                MaxBackups: logConfig.LOGFILEINFO.MAXBACKUPS,
                MaxAge:     logConfig.LOGFILEINFO.MAXAGE, //days
                Level:      level,
                Formatter: getFormatter(loggerName),
        })

        if err != nil {
                logrus.Fatalf("Failed to initialize file rotate hook: %v", err)
        }
        thisLogger.logrus.AddHook(rotateFileHook)
}

func SetLevel(strLevel string) {
	err := checkLevel(strLevel)
	if err != nil {
                logrus.Errorf("Failed to set log level: %v", err)
	}
	level, _ := logrus.ParseLevel(strLevel)
	thisLogger.logrus.SetLevel(level)
}

func checkLevel(lvl string) (error) {
	switch strings.ToLower(lvl) {
	case "error":
		return nil
	case "warn", "warning":
		return nil
	case "info":
		return nil
	case "debug":
		return nil
	}
	return fmt.Errorf("not a valid cblog Level: %q", lvl)
}

func GetLevel() string {
	return thisLogger.logrus.GetLevel().String()
}

func getFormatter(loggerName string) *cblogformatter.Formatter {

	if thisFormatter != nil {
		return thisFormatter
	}
	thisFormatter = &cblogformatter.Formatter{
            TimestampFormat: "2006-01-02 15:04:05",
            LogFormat:       "[" + loggerName + "]." + "[%lvl%]: %time% %func% - %msg%\n",
        }	
	return thisFormatter
}


