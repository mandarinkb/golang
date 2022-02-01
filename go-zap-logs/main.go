package main

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	config := zap.NewProductionConfig()
	config.OutputPaths = []string{"./mylog/project-final.log"}
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	logger, err := config.Build()
	if err != nil {
		logger.Error(err.Error())
	}
	defer logger.Sync()

	const url = "http://example.com"
	logger.Info("login", zap.String("url", url))
	fmt.Println("end")
}
