package utils

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	RedisHost      string `mapstructure:"REDIS_HOST"`
	RedisPassword  string `mapstructure:"REDIS_PASSWORD"`
	DriverName     string `mapstructure:"DRIVER_NAME"`
	DatasourceName string `mapstructure:"DATASOURCE_NAME"`
	Secretkey      string `mapstructure:"SECRETKEY"`
}

func LoadConfig(path string) (config *Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")
	// viper.SetConfigName("app")
	// viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
		return
	}

	err = viper.Unmarshal(&config)
	return

	// err = godotenv.Load("app.env")
	// if err != nil {
	// 	fmt.Println(err)
	// 	return nil, err
	// }
	// config = &Config{
	// 	RedisHost:      os.Getenv("REDIS_HOST"),
	// 	RedisPassword:  os.Getenv("REDIS_PASSWORD"),
	// 	DriverName:     os.Getenv("DRIVER_NAME"),
	// 	DatasourceName: os.Getenv("DATASOURCE_NAME"),
	// 	Secretkey:      os.Getenv("SECRETKEY"),
	// }
	// return config, nil
}
