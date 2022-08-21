package assets

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var e *RootError = &RootError{}

func LoadAssets(path string, fileNames ...string) *RootError {
	viper.AddConfigPath(path) // ระบุ path ของ config file
	for _, fileName := range fileNames {
		viper.SetConfigName(fileName) // ชื่อ config file
	}
	viper.AutomaticEnv() // อ่าน value จาก ENV variable

	// อ่าน config
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}

	viper.Unmarshal(e)

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Printf("config file change %s", e.Name)
		viper.Unmarshal(e)
	})
	return e
}

// get error
func E() *RootError {
	return e
}
