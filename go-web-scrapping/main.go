package main

import (
	"fmt"
	"time"
	"webscrapping/repository"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	///////////////////////
	// controller.Bigc()
	// controller.Makro()
	// controller.Lotus()
	// controller.All()
	// db, err := sql.Open("mysql", "root:mandarinkb@tcp(127.0.0.1)/TEST_DB?charset=utf8")
	// if err != nil {
	// 	fmt.Print(err)
	// }

	start := time.Now().Format(time.RFC3339)

	// makro, err := repository.NewMakro().Makro()
	// if err != nil {
	// 	fmt.Println(err)
	// }

	bigc, err := repository.NewBigC().Bigc()
	if err != nil {
		fmt.Println(err)
	}
	_ = bigc

	lotus, err := repository.NewLotus().Lotus()
	if err != nil {
		fmt.Println(err)
	}
	_ = lotus

	end := time.Now().Format(time.RFC3339)
	fmt.Println("start : ", start)
	fmt.Println("end   : ", end)
}
