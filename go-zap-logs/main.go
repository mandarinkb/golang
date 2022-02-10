package main

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func LogConf() (*zap.Logger, error) {
	config := zap.NewProductionConfig()
	// [for production] กรณี build ขึ้น docker จะเก็บไว้ที่ /home ใน container
	// config.OutputPaths = []string{"../home/go-zap-logs.log"}
	// [for vscode run]
	config.OutputPaths = []string{"./mylog/go-zap-logs.log"}
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

func main() {
	logger, err := LogConf()
	if err != nil {
		logger.Error(err.Error())
	}
	defer logger.Sync()

	const url = "http://example.com"
	logger.Info("login", Url(url), User("mandarinkb"), Type("[demo]"))
	fmt.Println("end")
}
