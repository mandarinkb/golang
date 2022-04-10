package utils

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	KafkaHost      string `mapstructure:"KAFKA_HOST"`
	DriverName     string `mapstructure:"DRIVER_NAME"`
	DatasourceName string `mapstructure:"DATASOURCE_NAME"`
	Elasticsearch  string `mapstructure:"ELASTICSEARCH"`
}

func LoadConfig(path string) (config *Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
		return
	}

	err = viper.Unmarshal(&config)
	return
}
