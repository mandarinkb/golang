package main

import (
	"errors"
	"fmt"

	"github.com/mandarinkb/go-viper-config-and-response/assets"
	"github.com/mandarinkb/go-viper-config-and-response/config"
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
	config.LoadConfig("config", "config")
	assets.LoadAssets("assets", "error")
	// logger.InitialLogger()
	// mainLogger := logger.L().Named("main")

	// for {
	// 	mainLogger.Info("running success")
	// 	time.Sleep(1 * time.Second)
	// }

	alphabet := Alphabet{
		Thai:    "ก ขฃคฅฆงจฉชซฌญฎฏฐฑฒณดตถทธนบปผฝพฟภมยรลวศษสหฬอฮ อะแอเออะอาเอียะเอออ ิเอียอําอ ีเอือะใออึเอือไออือัวะเอาอัวฤอูอะฤๅเอะโอฦเออาะฦๅแอะออ",
		English: "A aBbCcDdEeFfGgHhIiJjKkLlMmNnOoPpQqRrSsTtUuVvWwXxYyZz",
	}
	ExampleValidate(alphabet)

}
