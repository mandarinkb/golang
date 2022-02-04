package utils

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var TypeWeb string = "web"

func LogConf() (*zap.Logger, error) {
	config := zap.NewProductionConfig()
	config.OutputPaths = []string{"./mylog/project-final.log"}
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	logger, err := config.Build()
	if err != nil {
		return nil, err
	}
	return logger, nil
}

func Url(url string) zapcore.Field {
	return zap.String("url", url)
}
func User(user string) zapcore.Field {
	return zap.String("user", user)
}
func Type(ty string) zapcore.Field {
	return zap.String("typeLog", ty)
}
