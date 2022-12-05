package main

import (
	"context"

	"github.com/go-cronjob-with-crontab-ui/assets"
	"github.com/go-cronjob-with-crontab-ui/config"
	"github.com/go-cronjob-with-crontab-ui/logger"
	"github.com/go-cronjob-with-crontab-ui/repository"
	"github.com/go-cronjob-with-crontab-ui/service"
)

func main() {
	config.LoadConfig("config", "config")
	assets.LoadAssets("assets", "error")
	mainlog := logger.InitialLogger()
	mainlog.Info("main log")

	cronJobRepo := repository.NewCronJobRepository()
	service.NewCronJobService(cronJobRepo).RunJobService(context.TODO())
}
