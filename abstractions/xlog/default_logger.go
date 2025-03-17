package xlog

import (
	"fmt"
	"github.com/liangboceo/yuanboot/abstractions/platform/consolecolors"
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"time"
)

type XDefaultLogger struct {
	logger        *logrus.Logger
	dateFormat    string
	class         string
	logFormatter  func(interface{}) string
	fields        map[string]interface{}
	displayFields bool
	option        *LogOptions
}

func NewXLogger() *XDefaultLogger {
	logger := NewLoggerWith(log.New(os.Stdout, "", 0))
	return logger
}

func NewLoggerWith(log *log.Logger) *XDefaultLogger {
	logger := &XDefaultLogger{logger: logrus.New(), dateFormat: LoggerDefaultDateFormat}
	logger.SetCustomLogFormat(defaultLogFormatter)
	return logger
}

var LoggerDefaultDateFormat = "2006/01/02 15:04:05.00"

type LogLevel int

// MessageLevel
const (
	NOTSET  = iota
	DEBUG   = LogLevel(10 * iota) // DEBUG = 10
	INFO    = LogLevel(10 * iota) // INFO = 20
	WARNING = LogLevel(10 * iota) // WARNING = 30
	ERROR   = LogLevel(10 * iota) // ERROR = 40
)

var LevelString = map[LogLevel]string{
	DEBUG:   "DEBUG",
	INFO:    "INFO",
	WARNING: consolecolors.Yellow("WARNING"),
	ERROR:   consolecolors.Red("ERROR"),
}

func defaultLogFormatter(log interface{}) string {
	logInfo := log.(LogInfo)
	outLog := fmt.Sprintf("%s [%s] [%s] [%s] [%s] , %s",
		consolecolors.Yellow("[yuanboot]"), logInfo.StartTime, logInfo.Level, logInfo.Class, logInfo.Host, logInfo.Message)
	return outLog
}

func (log *XDefaultLogger) SetClass(className string) {
	log.class = className
}

func (log *XDefaultLogger) SetCustomLogFormat(logFormatterFunc func(logInfo interface{}) string) {
	log.logFormatter = logFormatterFunc
}

func (log *XDefaultLogger) SetDateFormat(format string) {
	log.dateFormat = format
}

func (log *XDefaultLogger) log(level LogLevel, format string, a ...interface{}) {
	hostName, _ := os.Hostname()
	message := format
	message = fmt.Sprintf(format, a...)

	start := time.Now()
	info := LogInfo{
		StartTime: start.Format(log.dateFormat),
		Level:     LevelString[level],
		Class:     log.class,
		Host:      hostName,
		Message:   message,
	}

	log.logger.Println(log.logFormatter(info))
}
func (log *XDefaultLogger) With(level LogLevel, fiedls map[string]interface{}) *logrus.Entry {

	//start := time.Now()

	fieldsMap := make(map[string]interface{})
	fieldsMap["prefix"] = "yuanboot"
	if fiedls != nil {
		fieldsMap = fiedls
	}

	if log.displayFields {
		fieldsMap["class"] = log.class
		hostName, _ := os.Hostname()
		fieldsMap["host"] = hostName
	}
	//fieldsMap["message"] = message

	return log.logger.WithFields(fieldsMap)
}

func (log *XDefaultLogger) Warning(format string, a ...interface{}) {
	log.With(WARNING, log.fields).Warnf(format, a...)
}

func (log *XDefaultLogger) Info(args ...interface{}) {
	log.With(INFO, log.fields).Info(args)
}

func (log *XDefaultLogger) Warn(args ...interface{}) {
	log.With(WARNING, log.fields).Warn(args)
}

func (log *XDefaultLogger) Error(args ...interface{}) {
	log.logger.Out = os.Stderr
	log.With(ERROR, log.fields).Error(args)
	log.logger.Out = os.Stdout
}

func (log *XDefaultLogger) Debug(args ...interface{}) {
	log.With(DEBUG, log.fields).Debug(args)
}

func (log *XDefaultLogger) Infof(fmt string, args ...interface{}) {
	log.With(INFO, log.fields).Infof(fmt, args)
}

func (log *XDefaultLogger) Warnf(fmt string, args ...interface{}) {
	log.With(WARNING, log.fields).Warnf(fmt, args)
}

func (log *XDefaultLogger) Errorf(fmt string, args ...interface{}) {
	log.logger.Out = os.Stderr
	log.With(ERROR, log.fields).Errorf(fmt, args)
	log.logger.Out = os.Stdout
}

func (log *XDefaultLogger) Debugf(fmt string, args ...interface{}) {
	log.With(DEBUG, log.fields).Debugf(fmt, args)
}

func (log *XDefaultLogger) GetOptions() *LogOptions {
	return log.option
}
