package main

import (
	"time"

	"github.com/mandarinkb/go-viper-config-and-response/assets"
	"github.com/mandarinkb/go-viper-config-and-response/config"
	"github.com/mandarinkb/go-viper-config-and-response/logger"
)

func main() {
	config.LoadConfig("config", "config")
	assets.LoadAssets("assets", "error")
	logger.InitialLogger()
	mainLogger := logger.L().Named("main")

	for {
		mainLogger.Info("running success")
		time.Sleep(1 * time.Second)
	}

}
