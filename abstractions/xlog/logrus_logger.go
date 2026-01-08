package xlog

import (
	"bytes"
	"fmt"
	logrus "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io"
	"os"
	"syscall"
)

type LogrusLogger struct {
	logger        *logrus.Logger
	dateFormat    string
	fields        map[string]interface{}
	displayFields bool
	class         string
	LogPath       string
}

type LogOptions struct {
	LogPath         string `mapstructure:"log_path"`
	LogLevel        string `mapstructure:"log_level"` // trace, debug, info, warn[ing], error, fatal, panic
	LogMaxDiskUsage int64  `mapstructure:"log_max_disk_usage"`
	LogMaxFileNum   int    `mapstructure:"log_max_file_num"`
	AppName         string `mapstructure:"app_name"`
}

func GetXLogger(class string) ILogger {
	configViper := viper.New()
	configViper.SetConfigFile("conf/log.yml")
	fileData, err := Fs.ReadFile(configViper.ConfigFileUsed())
	if err == nil {
		err = configViper.ReadConfig(bytes.NewReader(fileData))
	} else {
		configViper.SetConfigFile("./log.yml")
		err = configViper.ReadInConfig()
	}
	var option *LogOptions
	if err == nil {
		err = configViper.Sub("yuanboot.log").Unmarshal(&option)
	}
	if err != nil {
		logPath := "/mnt/data/log/"
		appName := fmt.Sprintf("app_%d", syscall.Getpid())
		logLevel := "debug"
		option = &LogOptions{LogLevel: logLevel, LogPath: logPath, LogMaxDiskUsage: 102400000, LogMaxFileNum: 50, AppName: appName}
	}
	logger := GetClassLogger(class, option) // NewXLogger()
	return logger
}
func GetXLoggerByLogLevel(class string, logLevel string) ILogger {
	configViper := viper.New()
	configViper.SetConfigFile("conf/log.yml")
	fileData, err := Fs.ReadFile(configViper.ConfigFileUsed())
	if err == nil {
		err = configViper.ReadConfig(bytes.NewReader(fileData))
	} else {
		configViper.SetConfigFile("./log.yml")
		err = configViper.ReadInConfig()
	}
	var option *LogOptions
	if err == nil {
		err = configViper.Sub("yuanboot.log").Unmarshal(&option)
	}
	if err != nil {
		logPath := "/mnt/data/log/"
		appName := fmt.Sprintf("app_%d", syscall.Getpid())
		option = &LogOptions{LogLevel: logLevel, LogPath: logPath, LogMaxDiskUsage: 102400000, LogMaxFileNum: 50, AppName: appName}
	} else {
		option.LogLevel = logLevel
	}
	logger := GetClassLogger(class, option) // NewXLogger()
	return logger
}

func GetXLoggerWithFields(class string, fields map[string]interface{}) ILogger {
	logger := NewXLogger()
	logger.class = class
	return logger
}

func GetXLoggerWith(logger ILogger) ILogger {
	return logger
}

func NewLogger(options *LogOptions) ILogger {
	logger := logrus.New()
	lw := &HourlySplit{
		Dir:           options.LogPath,
		FileFormat:    options.AppName + "_2006-01-02T15",
		MaxFileNumber: int64(options.LogMaxFileNum),
		MaxDiskUsage:  options.LogMaxDiskUsage,
	}
	multiWriter := io.MultiWriter(os.Stdout, lw)
	defer func(lw *HourlySplit) {
		_ = lw.Close()
	}(lw)
	logger.SetReportCaller(true)
	logger.SetOutput(multiWriter)
	lv, err := logrus.ParseLevel(options.LogLevel)
	if err != nil {
		lv = logrus.WarnLevel
	}
	logger.SetLevel(lv)
	return &LogrusLogger{logger: logger, LogPath: options.LogPath, dateFormat: LoggerDefaultDateFormat}
}

func GetClassLogger(class string, options *LogOptions) ILogger {
	logger := logrus.New()
	lw := &HourlySplit{
		Dir:           options.LogPath,
		FileFormat:    options.AppName + "_2006-01-02T15",
		MaxFileNumber: int64(options.LogMaxFileNum),
		MaxDiskUsage:  options.LogMaxDiskUsage,
	}
	defer func(lw *HourlySplit) {
		_ = lw.Close()
	}(lw)
	multiWriter := io.MultiWriter(os.Stdout, lw)
	logger.SetReportCaller(true)
	logger.SetOutput(multiWriter)
	lv, err := logrus.ParseLevel(options.LogLevel)
	if err != nil {
		lv = logrus.WarnLevel
	}
	logger.SetLevel(lv)
	logger.Formatter = &TextFormatter{
		DisableColors:   false,
		ForceColors:     false,
		TimestampFormat: LoggerDefaultDateFormat,
		FullTimestamp:   true,
		ForceFormatting: true,
	}
	return &LogrusLogger{logger: logger, LogPath: options.LogPath, class: class, dateFormat: LoggerDefaultDateFormat, displayFields: true}
}

func (log *LogrusLogger) With(level LogLevel, fiedls map[string]interface{}) *logrus.Entry {

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

func (log *LogrusLogger) Warning(format string, a ...interface{}) {
	log.With(WARNING, log.fields).Warnf(format, a...)
}

func (log *LogrusLogger) Info(args ...interface{}) {
	log.With(INFO, log.fields).Info(args...)
}

func (log *LogrusLogger) Warn(args ...interface{}) {
	log.With(WARNING, log.fields).Warn(args...)
}

func (log *LogrusLogger) Error(args ...interface{}) {
	log.logger.Out = os.Stderr
	log.With(ERROR, log.fields).Error(args...)
	log.logger.Out = os.Stdout
}

func (log *LogrusLogger) Debug(args ...interface{}) {
	log.With(DEBUG, log.fields).Debug(args)
}

func (log *LogrusLogger) Infof(fmt string, args ...interface{}) {
	if len(args) <= 0 {
		log.With(INFO, log.fields).Info(fmt)
	} else {
		log.With(INFO, log.fields).Infof(fmt, args...)
	}

}

func (log *LogrusLogger) Warnf(fmt string, args ...interface{}) {
	if len(args) <= 0 {
		log.With(WARNING, log.fields).Warn(fmt)
	} else {
		log.With(WARNING, log.fields).Warnf(fmt, args...)
	}
}

func (log *LogrusLogger) Errorf(fmt string, args ...interface{}) {
	log.logger.Out = os.Stderr
	if len(args) <= 0 {
		log.With(ERROR, log.fields).Error(fmt)
	} else {
		log.With(ERROR, log.fields).Errorf(fmt, args...)
	}
	log.logger.Out = os.Stdout
}

func (log *LogrusLogger) Debugf(fmt string, args ...interface{}) {
	if len(args) <= 0 {
		log.With(DEBUG, log.fields).Debug(fmt)
	} else {
		log.With(DEBUG, log.fields).Debugf(fmt, args...)
	}
}

func (log *LogrusLogger) SetClass(className string) {
	log.class = className
}

func (log *LogrusLogger) SetCustomLogFormat(logFormatterFunc func(logInfo interface{}) string) {
	log.displayFields = false
}

func (log *LogrusLogger) SetDateFormat(format string) {
	//log.dateFormat = format
}

func (log *LogrusLogger) GetLogPath() string {
	return log.LogPath
}
