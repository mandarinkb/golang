package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type User struct {
	USER_ID   int    `json:"user_id"`
	USERNAME  string `json:"username"`
	PASSWORD  string `json:"password"`
	USER_ROLE string `json:"user_role"`
}

func main() {
	// ใช้ sqlx library
	// mariadb ใช้ Library ของ mysql ในการเชื่อมต่อ
	// db, _ := sqlx.Connect("mysql", "root:mandarinkb@tcp(127.0.0.1)/TEST_DB?charset=utf8")

	// newUser := User{
	// 	USERNAME:  "fc",
	// 	PASSWORD:  "dkeiotsfan",
	// 	USER_ROLE: "admin",
	// }
	// upUser := User{
	// 	USER_ID:   1,
	// 	USERNAME:  "aaa",
	// 	PASSWORD:  "bbb",
	// 	USER_ROLE: "ccc",
	// }

	// err := createX(db, newUser)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// err := updateX(db, upUser)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// err := deleteX(db, 6)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// user, err := readByIdX(db, 5)

	// users, err := readX(db)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// j, _ := json.Marshal(users)
	// fmt.Println(string(j))

	// mariadb ใช้ Library ของ mysql ในการเชื่อมต่อ
	db, _ := sql.Open("mysql", "root:mandarinkb@tcp(127.0.0.1)/TEST_DB?charset=utf8")

	// err := create(db, newUser)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// err := update(db, upUser)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// err := delete(db, 1)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// user, err := readById(db, 2)

	users, err := read(db)
	if err != nil {
		fmt.Println(err)
	}

	j, _ := json.Marshal(users)
	fmt.Println(string(j))

	defer db.Close()
}

//********** sqlx library **********//
func readX(db *sqlx.DB) ([]User, error) {
	// ตรวจสอบว่า server ได้เปิดอยู่หรือไม่
	err := db.Ping()
	if err != nil {
		return nil, err
	}

	users := []User{}
	query := "SELECT user_id, username, password, user_role FROM USERS" // ชื่อ filed ใน database ต้องใส่เป็นตัวเล็ก
	err = db.Select(&users, query)                                      // และ ชื่อ table ต้องใส่เป็นตัวใหญ่ตามที่ตั้งค่าไว้ใน database
	if err != nil {                                                     // มิฉะนั้นจะ error
		return nil, err
	}
	return users, nil
}

func readByIdX(db *sqlx.DB, id int) (*User, error) {
	// ตรวจสอบว่า server ได้เปิดอยู่หรือไม่
	err := db.Ping()
	if err != nil {
		return nil, err
	}

	user := User{}
	query := "SELECT user_id, username, password, user_role FROM USERS WHERE user_id=?" //mysql ใช้ ? mssql ใช้ @id
	err = db.Get(&user, query, id)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func createX(db *sqlx.DB, user User) error {
	// ตรวจสอบว่า server ได้เปิดอยู่หรือไม่
	err := db.Ping()
	if err != nil {
		return err
	}

	tx, err := db.Begin() // ใช้งานแบบ transection กรณีใช้ exec หลายๆคำสั่ง
	if err != nil {
		return err
	}

	query := "INSERT INTO USERS (username,password,user_role) VALUES (?,?,?)"
	result, err := tx.Exec(query, user.USERNAME, user.PASSWORD, user.USER_ROLE)
	if err != nil {
		return err
	}

	// รับค่ามาเพื่อตรวสสอบว่า insert สำเร็จหรือไม่
	affected, err := result.RowsAffected()
	if err != nil {
		tx.Rollback()
		return err
	}

	// กรณี insert ไม่สำเร็จ
	if affected <= 0 {
		return errors.New("cannot insert")
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func updateX(db *sqlx.DB, user User) error {
	// ตรวจสอบว่า server ได้เปิดอยู่หรือไม่
	err := db.Ping()
	if err != nil {
		return err
	}

	query := "UPDATE USERS SET USERNAME=?,PASSWORD=?,USER_ROLE=? WHERE USER_ID=?"
	result, err := db.Exec(query, user.USERNAME, user.PASSWORD, user.USER_ROLE, user.USER_ID)
	if err != nil {
		return err
	}

	// รับค่ามาเพื่อตรวสสอบว่า update สำเร็จหรือไม่
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	// กรณี update ไม่สำเร็จ
	if affected <= 0 {
		return errors.New("cannot update")
	}

	return nil
}

func deleteX(db *sqlx.DB, id int) error {
	query := "DELETE FROM USERS WHERE USER_ID=?"
	result, err := db.Exec(query, id)
	if err != nil {
		return err
	}

	// รับค่ามาเพื่อตรวสสอบว่า delete สำเร็จหรือไม่
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	// กรณี delete ไม่สำเร็จ
	if affected <= 0 {
		return errors.New("cannot delete")
	}

	return nil
}

//********** build-in library **********//
func read(db *sql.DB) ([]User, error) {
	// ตรวจสอบว่า server ได้เปิดอยู่หรือไม่
	err := db.Ping()
	if err != nil {
		return nil, err
	}

	query := "SELECT USER_ID, USERNAME, PASSWORD, USER_ROLE FROM USERS"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}

	users := []User{}
	for rows.Next() {
		user := User{}
		err = rows.Scan(&user.USER_ID, &user.USERNAME, &user.PASSWORD, &user.USER_ROLE) //จะเรียงตามชื่อ field ที่ query
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	defer rows.Close()

	return users, nil
}

func readById(db *sql.DB, id int) (*User, error) {
	// ตรวจสอบว่า server ได้เปิดอยู่หรือไม่
	err := db.Ping()
	if err != nil {
		return nil, err
	}

	query := "SELECT USER_ID, USERNAME,PASSWORD, USER_ROLE FROM USERS WHERE USER_ID=?" //mysql ใช้ ? mssql ใช้ @id
	row := db.QueryRow(query, id)                                                      //ถ้ามี ? หลายตัว ก็ใส่ พารามิเตอร์ตาม index ไปเรื่อยๆ
	user := User{}
	err = row.Scan(&user.USER_ID, &user.USERNAME, &user.PASSWORD, &user.USER_ROLE)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func create(db *sql.DB, user User) error {
	// ตรวจสอบว่า server ได้เปิดอยู่หรือไม่
	err := db.Ping()
	if err != nil {
		return err
	}
	query := "INSERT INTO USERS (USERNAME,PASSWORD,USER_ROLE) VALUES (?,?,?)"
	result, err := db.Exec(query, user.USERNAME, user.PASSWORD, user.USER_ROLE)
	if err != nil {
		return err
	}

	// รับค่ามาเพื่อตรวสสอบว่า insert สำเร็จหรือไม่
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	// กรณี insert ไม่สำเร็จ
	if affected <= 0 {
		return errors.New("cannot insert")
	}

	return nil
}

func update(db *sql.DB, user User) error {
	// ตรวจสอบว่า server ได้เปิดอยู่หรือไม่
	err := db.Ping()
	if err != nil {
		return err
	}
	query := "UPDATE USERS SET USERNAME=?,PASSWORD=?,USER_ROLE=? WHERE USER_ID=?"
	result, err := db.Exec(query, user.USERNAME, user.PASSWORD, user.USER_ROLE, user.USER_ID)
	if err != nil {
		return err
	}

	// รับค่ามาเพื่อตรวสสอบว่า update สำเร็จหรือไม่
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	// กรณี update ไม่สำเร็จ
	if affected <= 0 {
		return errors.New("cannot update")
	}

	return nil
}

func delete(db *sql.DB, id int) error {
	// ตรวจสอบว่า server ได้เปิดอยู่หรือไม่
	err := db.Ping()
	if err != nil {
		return err
	}

	query := "DELETE FROM USERS WHERE USER_ID=?"
	result, err := db.Exec(query, id)
	if err != nil {
		return err
	}

	// รับค่ามาเพื่อตรวสสอบว่า delete สำเร็จหรือไม่
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	// กรณี delete ไม่สำเร็จ
	if affected <= 0 {
		return errors.New("cannot delete")
	}

	return nil
}
