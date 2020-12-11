// CB-Log: Logger for Cloud-Barista.
//
//      * Cloud-Barista: https://github.com/cloud-barista
//
// by CB-Log Team, 2019.08.
// ref) https://github.com/t-tomalak/logrus-easy-formatter

package cblogformatter

import (
	"fmt"

	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	// Default log format will output [INFO]: 2006-01-02T15:04:05Z07:00 - Log message
	defaultLogFormat       = "[%lvl%]: %time% %func% - %msg%\n"
	defaultTimestampFormat = time.RFC3339
)

// Formatter implements logrus.Formatter interface.
type Formatter struct {
	TimestampFormat string
	LogFormat       string
}

// Format building log message.
func (f *Formatter) Format(entry *logrus.Entry) ([]byte, error) {
	output := f.LogFormat
	if output == "" {
		output = defaultLogFormat
	}

	timestampFormat := f.TimestampFormat
	if timestampFormat == "" {
		timestampFormat = defaultTimestampFormat
	}

	level := strings.ToUpper(entry.Level.String())
	output = strings.Replace(output, "%lvl%", level, 1)

	output = strings.Replace(output, "%time%", entry.Time.Format(timestampFormat), 1)

	if entry.HasCaller() {
		fileVal := fmt.Sprintf("%s:%d", shortFilePathName(entry.Caller.File), entry.Caller.Line)
		funcVal := fmt.Sprintf("%s()", entry.Caller.Function)

		funcInfo := fileVal + ", " + funcVal

		output = strings.Replace(output, "%func%", funcInfo, 1)
	} else {
		output = strings.Replace(output, "%func%", "", 1)
	}

	output = strings.Replace(output, "%msg%", entry.Message, 1)

	// Key-Value pair 정보 처리 (WithField, WithFields, WithError, ...)
	if len(entry.Data) > 0 {
		keyValues := make([]string, 0)
		for k := range entry.Data {
			keyValues = append(keyValues, fmt.Sprintf("%s=%v", k, entry.Data[k]))
		}
		if len(keyValues) > 0 {
			output = strings.Replace(output, "%keyvalues%", strings.Join(keyValues, ","), 1)
		} else {
			output = strings.Replace(output, "\t[%keyvalues%]", "", 1)
		}
	} else {
		output = strings.Replace(output, "\t[%keyvalues%]", "", 1)
	}

	return []byte(output), nil
}

func shortFilePathName(filePath string) string {
	strArray := strings.Split(filePath, "/")

	return strArray[len(strArray)-1]
}
