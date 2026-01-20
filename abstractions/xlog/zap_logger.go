package xlog

import (
	"github.com/liangboceo/yuanboot/utils"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"strconv"
	"syscall"
)

type ZapLogger struct {
	logger        *zap.Logger
	sugar         *zap.SugaredLogger
	dateFormat    string
	fields        map[string]interface{}
	displayFields bool
	class         string
	LogPath       string
}

func NewZapLogger(options *LogOptions) ILogger {
	lw := &HourlySplit{
		Dir:           options.LogPath,
		FileFormat:    options.AppName + "_2006-01-02T15",
		MaxFileNumber: int64(options.LogMaxFileNum),
		MaxDiskUsage:  options.LogMaxDiskUsage,
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.TimeEncoderOfLayout(LoggerDefaultDateFormat),
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	core := zapcore.NewTee(
		zapcore.NewCore(zapcore.NewConsoleEncoder(encoderConfig), zapcore.AddSync(lw), getZapLogLevel(options.LogLevel)),
		zapcore.NewCore(zapcore.NewConsoleEncoder(encoderConfig), zapcore.AddSync(os.Stdout), getZapLogLevel(options.LogLevel)),
	)
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	defer func(logger *zap.Logger) {
		_ = logger.Sync()
	}(logger)

	return &ZapLogger{
		logger:        logger,
		sugar:         logger.Sugar(),
		LogPath:       options.LogPath,
		dateFormat:    LoggerDefaultDateFormat,
		displayFields: true,
		fields:        make(map[string]interface{}),
	}
}

func GetZapClassLogger(class string, options *LogOptions) ILogger {
	lw := &HourlySplit{
		Dir:           options.LogPath,
		FileFormat:    options.AppName + "_2006-01-02T15",
		MaxFileNumber: int64(options.LogMaxFileNum),
		MaxDiskUsage:  options.LogMaxDiskUsage,
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.TimeEncoderOfLayout(LoggerDefaultDateFormat),
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	core := zapcore.NewTee(
		zapcore.NewCore(zapcore.NewConsoleEncoder(encoderConfig), zapcore.AddSync(lw), getZapLogLevel(options.LogLevel)),
		zapcore.NewCore(zapcore.NewConsoleEncoder(encoderConfig), zapcore.AddSync(os.Stdout), getZapLogLevel(options.LogLevel)),
	)

	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	defer func(logger *zap.Logger) {
		_ = logger.Sync()
	}(logger)

	return &ZapLogger{
		logger:        logger,
		sugar:         logger.Sugar(),
		LogPath:       options.LogPath,
		dateFormat:    LoggerDefaultDateFormat,
		class:         class,
		displayFields: true,
		fields:        make(map[string]interface{}),
	}
}

func getZapLogLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn", "warning":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "fatal":
		return zapcore.FatalLevel
	case "panic":
		return zapcore.PanicLevel
	default:
		return zapcore.InfoLevel
	}
}

func (log *ZapLogger) With(level LogLevel, fields map[string]interface{}) *zap.SugaredLogger {
	var fieldsMap []interface{}
	fieldsMap = append(fieldsMap, zap.Any("level", LevelString[level]))
	fieldsMap = append(fieldsMap, zap.Any("prefix", "yuanboot-nio-"+strconv.Itoa(syscall.Getpid())+"-"+utils.GoId()))
	if fields != nil {
		for k, v := range fields {
			fieldsMap = append(fieldsMap, zap.Any(k, v))
		}
	}
	if log.displayFields {
		fieldsMap = append(fieldsMap, zap.Any("class", log.class))
		hostName, _ := os.Hostname()
		fieldsMap = append(fieldsMap, zap.Any("host", hostName))
	}
	return log.sugar.With(fieldsMap...)
}

func (log *ZapLogger) Warning(format string, a ...interface{}) {
	log.With(WARNING, log.fields).Warnf(format, a...)
}

func (log *ZapLogger) Info(args ...interface{}) {
	log.With(INFO, log.fields).Info(args...)
}

func (log *ZapLogger) Warn(args ...interface{}) {
	log.With(WARNING, log.fields).Warn(args...)
}

func (log *ZapLogger) Error(args ...interface{}) {
	log.With(ERROR, log.fields).Error(args...)
}

func (log *ZapLogger) Debug(args ...interface{}) {
	log.With(DEBUG, log.fields).Debug(args...)
}

func (log *ZapLogger) Infof(fmt string, args ...interface{}) {
	if len(args) <= 0 {
		log.With(INFO, log.fields).Info(fmt)
	} else {
		log.With(INFO, log.fields).Infof(fmt, args...)
	}
}

func (log *ZapLogger) Warnf(fmt string, args ...interface{}) {
	if len(args) <= 0 {
		log.With(WARNING, log.fields).Warn(fmt)
	} else {
		log.With(WARNING, log.fields).Warnf(fmt, args...)
	}
}

func (log *ZapLogger) Errorf(fmt string, args ...interface{}) {
	if len(args) <= 0 {
		log.With(ERROR, log.fields).Error(fmt)
	} else {
		log.With(ERROR, log.fields).Errorf(fmt, args...)
	}
}

func (log *ZapLogger) Debugf(fmt string, args ...interface{}) {
	if len(args) <= 0 {
		log.With(DEBUG, log.fields).Debug(fmt)
	} else {
		log.With(DEBUG, log.fields).Debugf(fmt, args...)
	}
}

func (log *ZapLogger) SetClass(className string) {
	log.class = className
}

func (log *ZapLogger) SetCustomLogFormat(logFormatterFunc func(logInfo interface{}) string) {
	log.displayFields = false
}

func (log *ZapLogger) SetDateFormat(format string) {
	log.dateFormat = format
}

func (log *ZapLogger) GetLogPath() string {
	return log.LogPath
}
