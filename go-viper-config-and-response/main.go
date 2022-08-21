package main

import (
	"fmt"

	"github.com/mandarinkb/go-viper-config-and-response/assets"
	"github.com/mandarinkb/go-viper-config-and-response/config"
)

func main() {
	c := config.LoadConfig("config", "config")
	a := assets.LoadAssets("assets", "error")

	fmt.Println(c.MariaDb.DriverName)
	fmt.Println(c.MariaDb.DataSourceName)
	fmt.Println(a.Success)

}
