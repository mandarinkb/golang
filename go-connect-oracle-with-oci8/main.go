package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-oci8"
)

type User struct {
	USER_ID   int    `db:"USER_ID"`
	USERNAME  string `db:"USERNAME"`
	PASSWORD  string `db:"PASSWORD"`
	USER_ROLE string `db:"USER_ROLE"`
}

func main() {
	// ("oci8", user/passowrd@host:port/service_name")

	//===== for sqlx conn =====//
	// db, err := sqlx.Open("oci8", "system/oracle@127.0.0.1:1521/xe")

	//===== for build-in conn =====//
	db, err := sql.Open("oci8", "system/oracle@127.0.0.1:1521/xe")

	if err != nil {
		fmt.Println(err)
		return
	}

	// addUser := User{
	// 	USER_ID:   6,
	// 	USERNAME:  "mod",
	// 	PASSWORD:  "dsafjlanrweor",
	// 	USER_ROLE: "admin",
	// }

	// updateUser := User{
	// 	USER_ID:   6,
	// 	USERNAME:  "666666-username",
	// 	PASSWORD:  "666666-password",
	// 	USER_ROLE: "666666-admin",
	// }

	// err = create(db, addUser)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// err = update(db, updateUser)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// err = delete(db, 6)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	users, err := read(db)
	// users, err := readById(db, 4)
	if err != nil {
		fmt.Println(err)
	}

	//=== use sqlx ===//

	// err = createX(db, addUser)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// err = updateX(db, updateUser)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// err = deleteX(db, 6)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// users, err := readX(db)
	// // users, err := readByIdX(db, 2)
	// if err != nil {
	// 	fmt.Println(err)
	// }

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
	query := "SELECT USER_ID, USERNAME,PASSWORD, USER_ROLE FROM USERS"
	err = db.Select(&users, query)
	if err != nil {
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

	query := "SELECT USER_ID, USERNAME, PASSWORD, USER_ROLE FROM USERS WHERE USER_ID=:id" //  ใช้ :id
	user := User{}
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

	query := "INSERT INTO USERS (USER_ID, USERNAME, PASSWORD, USER_ROLE) VALUES (:USER_ID,:USERNAME,:PASSWORD,:USER_ROLE)"
	result, err := tx.Exec(query, user.USER_ID, user.USERNAME, user.PASSWORD, user.USER_ROLE)
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
	query := "UPDATE USERS SET USERNAME=:USERNAME, PASSWORD=:PASSWORD, USER_ROLE=:USER_ROLE WHERE USER_ID=:USER_ID"
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
	// ตรวจสอบว่า server ได้เปิดอยู่หรือไม่
	err := db.Ping()
	if err != nil {
		return err
	}
	query := "DELETE FROM USERS WHERE USER_ID=:id"
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
	query := "SELECT USER_ID, USERNAME,PASSWORD, USER_ROLE FROM USERS"
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
	return users, nil
}

func readById(db *sql.DB, id int) (*User, error) {
	// ตรวจสอบว่า server ได้เปิดอยู่หรือไม่
	err := db.Ping()
	if err != nil {
		return nil, err
	}
	query := "SELECT USER_ID, USERNAME, PASSWORD, USER_ROLE FROM USERS WHERE USER_ID=:id" //  ใช้ :id
	row := db.QueryRow(query, id)
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
	query := "INSERT INTO USERS (USER_ID, USERNAME, PASSWORD, USER_ROLE) VALUES (:USER_ID,:USERNAME,:PASSWORD,:USER_ROLE)"
	result, err := db.Exec(query, user.USER_ID, user.USERNAME, user.PASSWORD, user.USER_ROLE)
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
	query := "UPDATE USERS SET USERNAME=:USERNAME, PASSWORD=:PASSWORD, USER_ROLE=:USER_ROLE WHERE USER_ID=:USER_ID"
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
	query := "DELETE FROM USERS WHERE USER_ID=:id"
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
