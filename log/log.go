package log

import (
	"fmt"
	"io"
	"path"
	"runtime"
	"strings"
	"time"

	cfg "github.com/0meet1/zero-framework/config"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"

	x0errors "github.com/pkg/errors"
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

const (
	CONSOLE_ENABLE  = "enable"
	CONSOLE_DISABLE = "disable"
)

var (
	sysLogger *logrus.Logger

	sysLogPath      string
	sysLogName      string
	sysRotationTime int
	sysMaxAge       int
	sysLvs          []string
	sysConsole      string
)

func initLoggerConfig() {
	sysLogName = cfg.StringValue("zero.log.name")
	if len(cfg.StringValue("zero.log.path")) > 0 {
		sysLogPath = cfg.StringValue("zero.log.path")
	} else {
		sysLogPath = path.Join(cfg.ServerAbsPath(), "logs")
	}
	sysConsole = cfg.StringValue("zero.log.console")
	sysMaxAge = cfg.IntValue("zero.log.maxAge")
	sysRotationTime = cfg.IntValue("zero.log.rotationTime")
	sysLvs = cfg.SliceStringValue("zero.log.level")

}

func makeLogWriter(level string) io.Writer {
	writer, err := rotatelogs.New(
		path.Join(sysLogPath, sysLogName+"."+level+".%Y%m%d%H.log"),
		rotatelogs.WithRotationTime(time.Hour*time.Duration(sysRotationTime)),
		rotatelogs.WithMaxAge(time.Hour*time.Duration(sysMaxAge)),
	)
	if err != nil {
		panic(err)
	}
	return writer
}

func InitLogger() *ZeroLogger {
	initLoggerConfig()
	sysLogger = logrus.New()
	wMap := lfshook.WriterMap{}
	for i := 0; i < len(sysLvs); i++ {
		level, err := logrus.ParseLevel(sysLvs[i])
		if err != nil {
			panic(err)
		}
		wMap[level] = makeLogWriter(level.String())
	}
	(*sysLogger).AddHook(lfshook.NewHook(wMap, new(LogFormatter)))
	(*sysLogger).SetFormatter(new(LogFormatter))
	if sysConsole == CONSOLE_ENABLE {
		(*sysLogger).SetOutput(new(ConsoleWriter))
	} else {
		(*sysLogger).SetOutput(makeLogWriter(logrus.TraceLevel.String()))
	}
	(*sysLogger).SetLevel(logrus.TraceLevel)

	return &ZeroLogger{}
}

type ZeroLogger struct {
}

func (logger *ZeroLogger) Debug(message string) {
	sysLogger.Debug(message)
}

func (logger *ZeroLogger) Debugf(format string, p ...string) {
	sysLogger.Debug(fmt.Sprintf(format, p))
}

func (logger *ZeroLogger) Info(message string) {
	sysLogger.Info(message)
}

func (logger *ZeroLogger) Infof(format string, p ...string) {
	sysLogger.Info(fmt.Sprintf(format, p))
}

func (logger *ZeroLogger) Warn(message string) {
	sysLogger.Warn(message)
}

func (logger *ZeroLogger) Warnf(format string, p ...string) {
	sysLogger.Warn(fmt.Sprintf(format, p))
}

func (logger *ZeroLogger) Error(message string) {
	sysLogger.Error(fmt.Sprintf("%+v", x0errors.WithStack(x0errors.New(message))))
}

func (logger *ZeroLogger) ErrorS(err error) {
	sysLogger.Error(fmt.Sprintf("%+v", x0errors.WithStack(err)))
}

func (logger *ZeroLogger) Errorf(format string, p ...string) {
	logger.Error(fmt.Sprintf(format, p))
}

func (logger *ZeroLogger) Fatal(message string) {
	sysLogger.Fatal(message)
}

func (logger *ZeroLogger) Panic(message string) {
	sysLogger.Panic(message)
}

func (logger *ZeroLogger) CallerInfosMaxLine(maxLine int, skip int) string {
	pc := make([]uintptr, maxLine)
	n := runtime.Callers(skip, pc)
	cf := runtime.CallersFrames(pc[:n])
	infos := ""
	for {
		frame, more := cf.Next()
		if !more {
			break
		} else {
			if len(frame.Func.Name()) > 0 {
				infos = fmt.Sprintf("%s\t at %s (%s:%d)\n", infos, frame.Func.Name(), frame.File, frame.Line)
			} else {
				infos = fmt.Sprintf("%s\t at (%s:%d)\n", infos, frame.File, frame.Line)
			}
		}
	}
	return infos
}

func (logger *ZeroLogger) CallerInfos() string {
	return logger.CallerInfosMaxLine(64, 3)
}
