package logging

import (
	"io"
	"io/ioutil"

	"github.com/cloud-barista/cb-apigw/restapigw/pkg/core"
	cblog "github.com/cloud-barista/cb-log"
	"github.com/sirupsen/logrus"
)

// ===== [ Constants and Variables ] =====

var (
	// logger - 시스템에서 사용할 기본 생성 Logger
	logger *Logger = nil
)

// ===== [ Types ] =====

type (
	// Fields - Logging 처리에 사용할 Field 정보 형식
	Fields map[string]interface{}

	// Logger - CB-LOG에서 사용하는 "logrus" Logger를 위한 Wrapper 구조
	Logger struct {
		*logrus.Logger
	}
)

// ===== [ Implementations ] =====

// SetOutput - 로그 출력기 설정
func (l *Logger) SetOutput(w io.Writer) {
	l.Logger.Out = w
}

// DisableOutput - 로그 출력 비활성화
func (l *Logger) DisableOutput() {
	l.SetOutput(ioutil.Discard)
}

// SetFormatter - 로그 포맷터 설정
func (l *Logger) SetFormatter(f logrus.Formatter) {
	l.Logger.Formatter = f
}

// SetLogLevel - 로그 레벨 설정
func (l *Logger) SetLogLevel(lv logrus.Level) {
	l.Logger.SetLevel(lv)
}

// SetFields - 로그에 사용할 Fields 정보 설정
func (l *Logger) SetFields(fields Fields) *Logger {
	l.WithFields(logrus.Fields(fields))
	return l
}

// ===== [ Private Functions ] =====

// init - 패키지 초기화
func init() {
}

// ===== [ Public Functions ] =====

// NewLogger - 초기화된 Logger의 인스턴스 생성
func NewLogger() *Logger {
	logger = &Logger{
		Logger: cblog.GetLogger(core.AppName),
	}

	return logger
}

// NewLoggerByName - 지정한 이름으로 구성된 Logger 인스턴스 생성
func NewLoggerByName(name string) *Logger {
	return &Logger{
		Logger: cblog.GetLogger(name),
	}
}

// GetLogger - 관리되고 있는 Logger 반환
func GetLogger() *Logger {
	return logger
}
