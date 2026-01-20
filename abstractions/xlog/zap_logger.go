package xlog

import (
	"github.com/liangboceo/yuanboot/utils"
	"go.uber.org/zap"
	"go.uber.org/zap/buffer"
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
	w := zapcore.AddSync(lw)
	zapcore.Lock(w)
	core := zapcore.NewTee(
		zapcore.NewCore(getProdEncoder(), w, getZapLogLevel(options.LogLevel)),
		zapcore.NewCore(getProdEncoder(), zapcore.Lock(os.Stdout), getZapLogLevel(options.LogLevel)),
	)
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	if options.PrintStack {
		logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1), zap.AddStacktrace(zapcore.ErrorLevel))
	}
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
	w := zapcore.AddSync(lw)
	zapcore.Lock(w)
	core := zapcore.NewTee(
		zapcore.NewCore(getProdEncoder(), w, getZapLogLevel(options.LogLevel)),
		zapcore.NewCore(getProdEncoder(), zapcore.Lock(os.Stdout), getZapLogLevel(options.LogLevel)),
	)
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	if options.PrintStack {
		logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1), zap.AddStacktrace(zapcore.ErrorLevel))
	}
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

func (log *ZapLogger) Warning(format string, a ...interface{}) {
	log.sugar.Warnf(format, a...)
}

func (log *ZapLogger) Info(args ...interface{}) {
	log.sugar.Info(args...)
}

func (log *ZapLogger) Warn(args ...interface{}) {
	log.sugar.Warn(args...)
}

func (log *ZapLogger) Error(args ...interface{}) {
	log.sugar.Error(args...)
}

func (log *ZapLogger) Debug(args ...interface{}) {
	log.sugar.Debug(args...)
}

func (log *ZapLogger) Infof(fmt string, args ...interface{}) {
	if len(args) <= 0 {
		log.sugar.Info(fmt)
	} else {
		log.sugar.Infof(fmt, args...)
	}
}

func (log *ZapLogger) Warnf(fmt string, args ...interface{}) {
	if len(args) <= 0 {
		log.sugar.Warn(fmt)
	} else {
		log.sugar.Warnf(fmt, args...)
	}
}

func (log *ZapLogger) Errorf(fmt string, args ...interface{}) {
	if len(args) <= 0 {
		log.sugar.Error(fmt)
	} else {
		log.sugar.Errorf(fmt, args...)
	}
}

func (log *ZapLogger) Debugf(fmt string, args ...interface{}) {
	if len(args) <= 0 {
		log.sugar.Debug(fmt)
	} else {
		log.sugar.Debugf(fmt, args...)
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

type prefixEncoder struct {
	zapcore.Encoder
	prefix  string
	bufPool buffer.Pool
}

func (e *prefixEncoder) EncodeEntry(entry zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	buf := e.bufPool.Get()

	buf.AppendString("[yuanboot-nio-" + strconv.Itoa(syscall.Getpid()) + "-" + utils.GoId() + "]")
	buf.AppendString(" ")

	logEntry, err := e.Encoder.EncodeEntry(entry, fields)
	if err != nil {
		return nil, err
	}

	_, err = buf.Write(logEntry.Bytes())
	if err != nil {
		return nil, err
	}

	return buf, nil
}
func getCustomerConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
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
}

func getDevEncoder() zapcore.Encoder {
	encoderConfig := getCustomerConfig()
	return &prefixEncoder{
		Encoder: zapcore.NewConsoleEncoder(encoderConfig),
		prefix:  "[yuanboot]",
		bufPool: buffer.NewPool(),
	}
}

func getProdEncoder() zapcore.Encoder {
	encoderConfig := getCustomerConfig()
	return &prefixEncoder{
		Encoder: zapcore.NewConsoleEncoder(encoderConfig),
		prefix:  "[yuanboot]",
		bufPool: buffer.NewPool(),
	}
}
