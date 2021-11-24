package utils

import "github.com/spf13/viper"

type Config struct {
	RedisHost      string `mapstructure:"REDIS_HOST"`
	RedisPassword  string `mapstructure:"REDIS_PASSWORD"`
	DriverName     string `mapstructure:"DRIVER_NAME"`
	DatasourceName string `mapstructure:"DATASOURCE_NAME"`
	Secretkey      string `mapstructure:"SECRETKEY"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
