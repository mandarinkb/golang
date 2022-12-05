package config

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var c *Root = &Root{}

func LoadConfig(path string, fileNames ...string) *Root {
	viper.AddConfigPath(path) // ระบุ path ของ config file
	for _, fileName := range fileNames {
		viper.SetConfigName(fileName) // ชื่อ config file
	}
	viper.AutomaticEnv() // อ่าน value จาก ENV variable

	// อ่าน config
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}

	viper.Unmarshal(c)

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Printf("config file change %s", e.Name)
		viper.Unmarshal(c)
	})
	return c
}

// get config
func C() *Root {
	return c
}
