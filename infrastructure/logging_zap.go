package infrastructure

import (
	"fmt"
	"os"

	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"github.com/mkaiho/go-lambda-api-sample/util"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var _ util.Logger = (*zapLogger)(nil)

type zapLogger struct {
	core      logr.Logger
	debugCore logr.Logger
}

func (l *zapLogger) Debug(msg string, keysAndValues ...interface{}) {
	l.debugCore.Info(msg, keysAndValues...)
}

func (l *zapLogger) Info(msg string, keysAndValues ...interface{}) {
	l.core.Info(msg, keysAndValues...)
}

func (l *zapLogger) Error(err error, msg string, keysAndValues ...interface{}) {
	l.core.Error(err, msg, keysAndValues...)
}

func (l *zapLogger) Fatal(err error, msg string, keysAndValues ...interface{}) {
	l.Error(err, msg, keysAndValues...)
	os.Exit(1)
}

func (l *zapLogger) WithValues(keysAndValues ...interface{}) util.Logger {
	return &zapLogger{
		core:      l.core.WithValues(keysAndValues...),
		debugCore: l.debugCore.WithValues(keysAndValues...),
	}
}

func InjectZapLogger(conf util.LoggerConfig) {
	if l, err := NewZapLogger(conf); err != nil {
		panic(err)
	} else {
		util.InitLogger(l)
	}
}

func NewZapLogger(conf util.LoggerConfig) (util.Logger, error) {
	var (
		c zap.Config
		l logr.Logger
	)
	{
		level, err := parseZapcoreLevel(conf.Level())
		if err != nil {
			return nil, err
		}

		c = zap.NewProductionConfig()
		c.Level.SetLevel(*level)
		c.Encoding = "json"
		c.EncoderConfig.TimeKey = "time"
		c.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
		c.EncoderConfig.CallerKey = "src"
		c.OutputPaths = []string{"stdout"}
		zapLogger, err := c.Build()
		if err != nil {
			return nil, err
		}
		l = zapr.NewLogger(zapLogger)
	}

	return &zapLogger{
		core:      l,
		debugCore: l.V(1),
	}, nil
}

func parseZapcoreLevel(level util.LogLevel) (*zapcore.Level, error) {
	var _l zapcore.Level
	switch level {
	case util.DebugLevel:
		_l = zapcore.DebugLevel
	case util.InfoLevel:
		_l = zapcore.InfoLevel
	case util.ErrorLevel:
		_l = zapcore.ErrorLevel
	default:
		return nil, fmt.Errorf("not supported zapcore level was specified")
	}

	return &_l, nil
}
