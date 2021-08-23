package main

import (
	"connect-database/repository"
	"database/sql"
	"encoding/json"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// mariadb ใช้ Library ของ mysql ในการเชื่อมต่อ
	db, _ := sql.Open("mysql", "root:mandarinkb@tcp(127.0.0.1)/TEST_DB?charset=utf8")

	// newUser := repository.User{
	// 	USERNAME: "test-inasert",
	// 	PASSWORD: "jflas;owreklfclsf",
	// 	ROLE:     "assistant",
	// }

	// updatUser := repository.User{47, "upadate", "updatet-password", "admin"}

	// err := repository.NewMysqlRepo(db).MariadbCreate(newUser)
	// err := repository.NewMysqlRepo(db).MariadbUpdate(updatUser)
	// err := repository.NewMysqlRepo(db).MariadbDelete(47)

	// if err != nil {
	// 	fmt.Println(err)
	// }

	users, err := repository.NewMysqlRepo(db).MariadbRead()
	// user, err := repository.NewMysqlRepo(db).MariadbReadById(1)
	if err != nil {
		fmt.Println(err)
	}
	j, _ := json.Marshal(users)

	fmt.Println(string(j))
	defer db.Close()
}
