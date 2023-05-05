package logger

import (
	"context"
	"fmt"
	"log"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	LOG_FIELD_REQ_ID = "request_id"
	LOG_FIELD_PATH   = "path"
	LOG_FIELD_KEY    = "logFields"
	DebugLevel       = zap.DebugLevel
	InfoLevel        = zap.InfoLevel
	WarnLevel        = zap.WarnLevel
	ErrorLevel       = zap.ErrorLevel
	DPanicLevel      = zap.DPanicLevel
	PanicLevel       = zap.PanicLevel
	FatalLevel       = zap.FatalLevel
)

type (
	Level zapcore.Level
)

type Logger struct {
	*zap.Logger
}

var logger *Logger = &Logger{zap.L()}

var withContextHandler func(*Logger, context.Context) *Logger = DefaultWithContextHandler

func IsLevelEnabled(lv1 zapcore.Level) bool {
	return logger.Core().Enabled(lv1)
}

func InitialLogger() *Logger {
	var err error
	zapLogger, err := NewProduction(zap.AddCaller(), zap.AddCallerSkip(1))
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	logger = &Logger{zapLogger}
	zap.ReplaceGlobals(zapLogger)
	defer logger.Sync()
	return logger
}

func NewProduction(options ...zap.Option) (*zap.Logger, error) {
	return NewProductionConfig().Build(options...)
}

func NewProductionConfig() zap.Config {
	var lv = zapcore.DebugLevel
	if envValue, isSet := os.LookupEnv("LOG_LEVEL"); isSet {
		if err := (&lv).Set(envValue); err != nil {
			log.Fatalf("can't parse LOG_LEVEL: %s.", envValue)
		}
	}
	return zap.Config{
		Level:       zap.NewAtomicLevelAt(lv),
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:          "json",
		EncoderConfig:     NewProductionEncoderConfig(),
		OutputPaths:       []string{"stderr", "myproject.log"},
		ErrorOutputPaths:  []string{"stderr"},
		DisableCaller:     true,
		DisableStacktrace: true,
	}
}

func NewProductionEncoderConfig() zapcore.EncoderConfig {
	config := zapcore.EncoderConfig{
		LevelKey:       "severity",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    "func",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.MillisDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	if _, isKibana := os.LookupEnv("LOG_FORMAT_ONPREM"); isKibana {
		config.TimeKey = "@timestamp"
		config.MessageKey = "@message"
		config.LevelKey = "level"
	} else {
		config.TimeKey = "time"
	}
	return config
}

func L() *Logger {
	return logger
}

func SetWithContextHandler(fn func(*Logger, context.Context) *Logger) {
	withContextHandler = fn
}
func (log *Logger) WithContext(ctx context.Context) *Logger {
	return withContextHandler(log, ctx)
}

func (log *Logger) With(fields ...zap.Field) *Logger {
	return &Logger{log.Logger.With(fields...)}
}
func (log *Logger) Named(s string) *Logger {
	return &Logger{log.Logger.Named(s)}
}

func (log *Logger) Infof(format string, a ...interface{}) {
	log.Info(fmt.Sprintf(format, a...))
}

func (log *Logger) Debugf(format string, a ...interface{}) {
	log.Debug(fmt.Sprintf(format, a...))
}

func (log *Logger) Warnf(format string, a ...interface{}) {
	log.Warn(fmt.Sprintf(format, a...))
}

func (log *Logger) Errorf(format string, a ...interface{}) {
	log.Error(fmt.Sprintf(format, a...))
}

func (log *Logger) Fatalf(format string, a ...interface{}) {
	log.Fatal(fmt.Sprintf(format, a...))
}

func DefaultWithContextHandler(log *Logger, ctx context.Context) *Logger {
	if ctx == nil {
		return log
	}
	if fs, ok := ctx.Value(LOG_FIELD_KEY).([]zap.Field); ok {
		return log.With(fs...)

	} else {
		return log
	}
}

// func GinMiddleware(c *gin.Context){
// 	id := uuid.New().string
// 	c.Set(LOG_FIELD_KEY,[]zap.Field{
// 		zap.String(LOG_FIELD_PATH,c.Request.URL.Path),
// 		zap.String(LOG_FIELD_REQ_ID,id),
// 	})
// 	c.Set(LOG_FIELD_PATH,c.Request.URL.Path)
// 	c.Set(LOG_FIELD_REQ_ID,id)
// 	c.Next()
// }
