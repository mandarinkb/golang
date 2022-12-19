package config

type Root struct {
	Crontab        CrontabConfig   `mapstructure:"crontab"`
	CronJobService CronJobServices `mapstructure:"cronjob_service"`
}

type CrontabConfig struct {
	Path string `mapstructure:"path"`
}

type CronJobServices struct {
	ServiceOne string `mapstructure:"service_one"`
	ServiceTwo string `mapstructure:"service_two"`
}
