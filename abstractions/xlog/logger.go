package xlog

type ILogger interface {
	SetClass(className string)
	Debug(a ...interface{})
	Info(a ...interface{})
	Warning(format string, a ...interface{})
	Warn(args ...interface{})
	Error(a ...interface{})
	Debugf(fmt string, args ...interface{})
	Infof(fmt string, args ...interface{})
	Warnf(fmt string, args ...interface{})
	Errorf(fmt string, args ...interface{})
	SetCustomLogFormat(logFormatterFunc func(logInfo interface{}) string)
	SetDateFormat(format string)
	GetOptions() *LogOptions
}
