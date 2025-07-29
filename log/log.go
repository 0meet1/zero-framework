package log

import (
	"fmt"
	"io"
	"path"
	"strings"
	"time"

	cfg "github.com/0meet1/zero-framework/config"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"

	x0errors "github.com/pkg/errors"
)

const (
	DEBUG = "DEBUG"
	INFO  = "INFO"
	WARN  = "WARN"
	ERROR = "ERROR"
	FATAL = "FATAL"
	PANIC = "PANIC"

	CONSOLE_ENABLE  = "enable"
	CONSOLE_DISABLE = "disable"
)

type LogFormatter struct{}

func (s *LogFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	return []byte(fmt.Sprintf("[%s][%s] %s\n", time.Now().Format("2006-01-02 15:04:05.000"), strings.ToUpper(entry.Level.String()), entry.Message)), nil
}

type ConsoleWriter struct{}

func (writer *ConsoleWriter) Write(bytes []byte) (int, error) {
	fmt.Printf("%s", string(bytes))
	return len(bytes), nil
}

var supportLevels = map[string]logrus.Level{
	DEBUG: logrus.DebugLevel,
	INFO:  logrus.InfoLevel,
	WARN:  logrus.WarnLevel,
	ERROR: logrus.ErrorLevel,
	FATAL: logrus.FatalLevel,
	PANIC: logrus.PanicLevel,
}

type ZeroLogger struct {
	rootLogger *logrus.Logger

	asLogPath      string
	asLogName      string
	asRotationTime int
	asMaxAge       int
	asLevels       []string
	asConsole      string

	supports map[string]logrus.Level
}

func (aLogger *ZeroLogger) readLoggerConfig(prefix string) {
	aLogger.asLogName = cfg.StringValue(fmt.Sprintf("%s.name", prefix))
	if len(cfg.StringValue(fmt.Sprintf("%s.path", prefix))) > 0 {
		aLogger.asLogPath = cfg.StringValue(fmt.Sprintf("%s.path", prefix))
	} else {
		aLogger.asLogPath = path.Join(cfg.ServerAbsPath(), "logs")
	}
	aLogger.asConsole = cfg.StringValue(fmt.Sprintf("%s.console", prefix))
	aLogger.asMaxAge = cfg.IntValue(fmt.Sprintf("%s.maxAge", prefix))
	aLogger.asRotationTime = cfg.IntValue(fmt.Sprintf("%s.rotationTime", prefix))
	aLogger.asLevels = cfg.SliceStringValue(fmt.Sprintf("%s.level", prefix))
	aLogger.supports = make(map[string]logrus.Level)
	for _, enableLv := range aLogger.asLevels {
		if _, ok := supportLevels[enableLv]; !ok {
			panic(fmt.Errorf(" no support log level `%s` ", enableLv))
		}
		aLogger.supports[enableLv] = supportLevels[enableLv]
	}
}

func (aLogger *ZeroLogger) initWriter(moduleName string) (io.Writer, error) {
	return rotatelogs.New(
		path.Join(aLogger.asLogPath, aLogger.asLogName+"."+moduleName+".%Y%m%d%H.log"),
		rotatelogs.WithRotationTime(time.Hour*time.Duration(aLogger.asRotationTime)),
		rotatelogs.WithMaxAge(time.Hour*time.Duration(aLogger.asMaxAge)),
	)
}

func (aLogger *ZeroLogger) initLogger(prefix string) {
	aLogger.readLoggerConfig(prefix)
	aLogger.rootLogger = logrus.New()
	writers := lfshook.WriterMap{}
	for _, v := range aLogger.supports {
		_writer, err := aLogger.initWriter(v.String())
		if err != nil {
			panic(err)
		}
		writers[v] = _writer
	}
	aLogger.rootLogger.AddHook(lfshook.NewHook(writers, &LogFormatter{}))
	aLogger.rootLogger.SetFormatter(&LogFormatter{})
	if aLogger.asConsole == CONSOLE_ENABLE {
		aLogger.rootLogger.SetOutput(&ConsoleWriter{})
	} else {
		_writer, err := aLogger.initWriter(logrus.TraceLevel.String())
		if err != nil {
			panic(err)
		}
		aLogger.rootLogger.SetOutput(_writer)
	}
	aLogger.rootLogger.SetLevel(logrus.TraceLevel)
}

func NewLogger(prefix string) *ZeroLogger {
	nLogger := &ZeroLogger{}
	nLogger.initLogger(prefix)
	return nLogger
}

func (aLogger *ZeroLogger) Debug(message string) {
	if _, ok := aLogger.supports[DEBUG]; ok {
		aLogger.rootLogger.Debug(message)
	}
}

func (aLogger *ZeroLogger) Debugf(format string, p ...any) {
	if _, ok := aLogger.supports[DEBUG]; ok {
		aLogger.rootLogger.Debug(fmt.Sprintf(format, p...))
	}
}

func (aLogger *ZeroLogger) Info(message string) {
	if _, ok := aLogger.supports[INFO]; ok {
		aLogger.rootLogger.Info(message)
	}
}

func (aLogger *ZeroLogger) Infof(format string, p ...any) {
	if _, ok := aLogger.supports[INFO]; ok {
		aLogger.rootLogger.Info(fmt.Sprintf(format, p...))
	}
}

func (aLogger *ZeroLogger) Warn(message string) {
	if _, ok := aLogger.supports[WARN]; ok {
		aLogger.rootLogger.Warn(message)
	}
}

func (aLogger *ZeroLogger) Warnf(format string, p ...any) {
	if _, ok := aLogger.supports[WARN]; ok {
		aLogger.rootLogger.Warn(fmt.Sprintf(format, p...))
	}
}

func (aLogger *ZeroLogger) Error(message string) {
	if _, ok := aLogger.supports[ERROR]; ok {
		aLogger.rootLogger.Error(fmt.Sprintf("%+v", x0errors.WithStack(x0errors.New(message))))
	}
}

func (aLogger *ZeroLogger) ErrorS(err error) {
	if _, ok := aLogger.supports[ERROR]; ok {
		aLogger.rootLogger.Error(fmt.Sprintf("%+v", x0errors.WithStack(err)))
	}
}

func (aLogger *ZeroLogger) Errorf(format string, p ...any) {
	aLogger.Error(fmt.Sprintf(format, p...))
}

func (aLogger *ZeroLogger) Fatal(message string) {
	if _, ok := aLogger.supports[FATAL]; ok {
		aLogger.rootLogger.Fatal(message)
	}
}

func (aLogger *ZeroLogger) Panic(message string) {
	if _, ok := aLogger.supports[PANIC]; ok {
		aLogger.rootLogger.Panic(message)
	}
}
