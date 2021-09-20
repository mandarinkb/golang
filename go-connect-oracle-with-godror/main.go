package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	_ "github.com/godror/godror"
	"github.com/jmoiron/sqlx"
)

type User struct {
	USER_ID   int    `db:"USER_ID"`
	USERNAME  string `db:"USERNAME"`
	PASSWORD  string `db:"PASSWORD"`
	USER_ROLE string `db:"USER_ROLE"`
}

func main() {
	// ("godror", `user="your user" password="your passowrd" connectString="host:port/service_name"`)

	//===== for sqlx conn =====//
	// db, err := sqlx.Open("godror", `user="system" password="oracle" connectString="127.0.0.1:1521/xe"`)

	//===== for build-in conn =====//
	db, err := sql.Open("godror", `user="system" password="oracle" connectString="127.0.0.1:1521/xe"`)

	if err != nil {
		fmt.Println(err)
		return
	}
	//========== for test ru ==========//
	// db, err := sqlx.Open("godror", `user="scenter01" password="scenter01new" connectString="10.2.1.98:1571/RUBRAM" timezone=local`)

	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// err = db.Ping()
	// if err != nil {
	// 	fmt.Println(err)
	// } else {
	// 	fmt.Println("connect success")
	// }
	// user, err := GetAll(db)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// j, _ := json.Marshal(user)
	// fmt.Println(string(j))
	//========== end test ru ==========//

	// addUser := User{
	// 	USER_ID:   6,
	// 	USERNAME:  "mod",
	// 	PASSWORD:  "dsafjlanrweor",
	// 	USER_ROLE: "admin",
	// }
	// updateUser := User{
	// 	USER_ID:   6,
	// 	USERNAME:  "333update-username",
	// 	PASSWORD:  "333update-password",
	// 	USER_ROLE: "333update-admin",
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
	// users, err := readById(db, 6)
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
	// // users, err := readByIdX(db, 3)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	j, _ := json.Marshal(users)
	fmt.Println(string(j))
	defer db.Close()
}

//========== for test ru ==========//
// type User struct {
// 	UserId            string `db:"USER_ID"`
// 	UserName          string `db:"USERNAME"`
// 	UserRoleType      string `db:"ROLE_TYPE"`
// 	UserUpdateDate    string `db:"UPDATE_DATE"`
// 	UserActivity      string `db:"ACTIVATED"`
// 	UserFacultyNumber string `db:"FACULTY_NO"`
// }

// func GetAll(db *sqlx.DB) ([]User, error) {

// 	// ทำการสร้าง instant ค่าเป็น slide User{} ที่ได้สร้างไว้ที่  (user.go file) type User struct
// 	users := []User{}
// 	query := "SELECT USER_ID,USERNAME,ROLE_TYPE,TO_CHAR(UPDATE_DATE,'ddMONyyyy', 'NLS_CALENDAR=''THAI BUDDHA'' NLS_DATE_LANGUAGE=THAI')UPDATE_DATE,ACTIVATED,FACULTY_NO FROM SCENTER01.SV_USERS ORDER BY USER_ID"

// 	// u.db.Select เป็นการใช้ attribute ของ type userRepositoryDB struct { db *sqlx.DB }
// 	// ซึ่งเป็นการใช้คุณสมบัติ คำสั่ง query ต่างๆของตัวแปร db ที่มี Data type เป็น *sqlx.DB
// 	err := db.Select(&users, query)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return users, nil

// }
//========== end test ru ==========//

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
