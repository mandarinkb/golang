package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/mandarinkb/go-viper-config-and-response/assets"
	"github.com/mandarinkb/go-viper-config-and-response/config"
	"github.com/mandarinkb/go-viper-config-and-response/database"
	"github.com/mandarinkb/go-viper-config-and-response/logger"
	"github.com/mandarinkb/go-viper-config-and-response/repository"
	"github.com/mandarinkb/go-viper-config-and-response/utils"
)

type Alphabet struct {
	Thai    string `json:"thai" validate:"tha"`
	English string `json:"english" validate:"eng"`
}

func ExampleValidate(req Alphabet) {
	if err := utils.V().Struct(req); err != nil {
		fmt.Println(err)
		MapStrucValidationError(err)
	}
}
func MapStrucValidationError(err error) error {
	var newErr error
	for _, v := range utils.ToValidationError(err) {
		switch v.Field() {
		case "Thai":
			newErr = errors.New("validat Thai")
			fmt.Println(newErr)
		case "English":
			newErr = errors.New("validat English")
			fmt.Println(newErr)
		}

	}
	return newErr
}

func main() {
	cfg := config.LoadConfig("config", "config")
	assets.LoadAssets("assets", "error")
	mainLog := logger.InitialLogger()
	database.NewClient(cfg.Redis)
	// mainLogger := logger.L().Named("main")

	// for {
	// 	mainLogger.Info("running success")
	// 	time.Sleep(1 * time.Second)
	// }

	// alphabet := Alphabet{
	// 	Thai:    "ก ขฃคฅฆงจฉชซฌญฎฏฐฑฒณดตถทธนบปผฝพฟภมยรลวศษสหฬอฮ อะแอเออะอาเอียะเอออ ิเอียอําอ ีเอือะใออึเอือไออือัวะเอาอัวฤอูอะฤๅเอะโอฦเออาะฦๅแอะออ",
	// 	English: "A aBbCcDdEeFfGgHhIiJjKkLlMmNnOoPpQqRrSsTtUuVvWwXxYyZz",
	// }
	// ExampleValidate(alphabet)

	rd := repository.NewRedisRepository()

	alphabet := Alphabet{
		Thai:    "กขคง",
		English: "abcd",
	}
	err := rd.SaveRedis(context.Background(), cfg.RedisOption.MyOption.KeyFormat, alphabet, cfg.RedisOption.MyOption.TTL)
	if err != nil {
		mainLog.Errorf("save redis error key %v error: %v", cfg.RedisOption.MyOption.KeyFormat, err)
	}
	mainLog.Info("save redis success")

	// rdData := Alphabet{}
	// err := rd.GetRedis(context.Background(), cfg.RedisOption.MyOption.KeyFormat, &rdData)
	// if err != nil {
	// 	mainLog.Errorf("get redis error key %v error: %v", cfg.RedisOption.MyOption.KeyFormat, err)
	// }
	// js, _ := json.Marshal(rdData)
	// fmt.Println(string(js))

	// err := rd.RemoveRedis(context.Background(), cfg.RedisOption.MyOption.KeyFormat)
	// if err != nil {
	// 	mainLog.Errorf("remove redis error key %v error: %v", cfg.RedisOption.MyOption.KeyFormat, err)
	// }
}
