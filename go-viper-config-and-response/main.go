package main

import (
	"fmt"

	"github.com/mandarinkb/go-viper-config-and-response/assets"
	"github.com/mandarinkb/go-viper-config-and-response/config"
	"github.com/mandarinkb/go-viper-config-and-response/repository"
)

type user struct {
	Username string `json:"username"`
	UserRole string `json:"userRole"`
}

func main() {
	config.LoadConfig("config", "config")
	assets.LoadAssets("assets", "error")
	user := user{
		Username: "test user",
		UserRole: "test user role",
	}
	es := repository.NewElasticSearchRepository()

	msg, err := es.ESBuildJsonToByte(user)
	if err != nil {
		fmt.Println(err)
		return
	}
	data, err := es.ESPost(config.C().Elastic.Index, config.C().Elastic.Type, msg)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(data))

}
