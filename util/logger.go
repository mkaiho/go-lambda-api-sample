package util

import (
	"fmt"
)

var defaultLogger *Logger

func InitLogger(logger Logger) {
	defaultLogger = &logger
}

func GetLogger() Logger {
	if defaultLogger == nil {
		panic("logger is not Initialize")
	}
	return *defaultLogger
}

/** LogLevel **/
type LogLevel int

func (l LogLevel) validate() error {
	switch l {
	case DebugLevel, InfoLevel, WarnLevel, ErrorLevel:
		return nil
	default:
		return fmt.Errorf("not supported log level was specified: %d", l)
	}
}

const (
	DebugLevel LogLevel = iota - 1
	InfoLevel
	WarnLevel
	ErrorLevel
)

/** LoggerConfig **/
type LoggerConfig interface {
	Level() LogLevel
}

type loggerConfig struct {
	level LogLevel
}

func (c *loggerConfig) Level() LogLevel {
	return c.level
}

func (l *loggerConfig) validate() error {
	if err := l.level.validate(); err != nil {
		return err
	}
	return nil
}

func NewLoggerConfig(level LogLevel) (LoggerConfig, error) {
	loggerConfig := loggerConfig{
		level: level,
	}
	if err := loggerConfig.validate(); err != nil {
		return nil, err
	}
	return &loggerConfig, nil
}

/** Logger **/
type Logger interface {
	Debug(msg string, keysAndValues ...interface{})
	Info(msg string, keysAndValues ...interface{})
	Error(err error, msg string, keysAndValues ...interface{})
	Fatal(err error, msg string, keysAndValues ...interface{})
	WithValues(keysAndValues ...interface{}) Logger
}
