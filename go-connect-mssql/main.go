package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/jmoiron/sqlx"
)

type User struct {
	USER_ID   int    `db:"USER_ID"`
	USERNAME  string `db:"USERNAME"`
	PASSWORD  string `db:"PASSWORD"`
	USER_ROLE string `db:"USER_ROLE"`
}

func main() {
	//===== ใช้ sqlx library =====//
	// db, _ := sqlx.Open("sqlserver", "server=localhost;user id=sa;password=P@ssw0rd;port=1433;database=master;")

	// newUser := User{
	// 	USER_ID:   6,
	// 	USERNAME:  "fc",
	// 	PASSWORD:  "dkeiotsfan",
	// 	USER_ROLE: "admin",
	// }

	// upUser := User{
	// 	USER_ID:   6,
	// 	USERNAME:  "ssssssss",
	// 	PASSWORD:  "dddddddd",
	// 	USER_ROLE: "ffffffff",
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

	// users, err := readByIdX(db, 5)

	// users, err := readX(db)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	//===== ใช้ build-in library =====//
	db, _ := sql.Open("sqlserver", "server=localhost;user id=sa;password=P@ssw0rd;port=1433;database=master;")

	// err := create(db, newUser)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// err := update(db, upUser)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// err := delete(db, 6)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// users, err := readById(db, 2)

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
	query := "SELECT USER_ID, USERNAME, PASSWORD, USER_ROLE FROM USERS"
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

	user := User{}
	query := "SELECT USER_ID, USERNAME,PASSWORD, USER_ROLE FROM USERS WHERE USER_ID=@id" //mysql ใช้ ? mssql ใช้ @id
	err = db.Get(&user, query, sql.Named("id", id))
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

	query := "INSERT INTO USERS (USER_ID,USERNAME,PASSWORD,USER_ROLE) VALUES (@user_id,@username,@password,@user_role)"
	result, err := tx.Exec(query,
		sql.Named("user_id", user.USER_ID),
		sql.Named("username", user.USERNAME),
		sql.Named("password", user.PASSWORD),
		sql.Named("user_role", user.USER_ROLE))
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

	query := "UPDATE USERS SET USERNAME=@username, PASSWORD=@password, USER_ROLE=@user_role WHERE USER_ID=@user_id"
	result, err := db.Exec(query,
		sql.Named("username", user.USERNAME),
		sql.Named("password", user.PASSWORD),
		sql.Named("user_role", user.USER_ROLE),
		sql.Named("user_id", user.USER_ID))
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
	query := "DELETE FROM USERS WHERE USER_ID=@id"
	result, err := db.Exec(query, sql.Named("id", id))
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

	query := "SELECT USER_ID, USERNAME,PASSWORD, USER_ROLE FROM USERS WHERE USER_ID=@id" //mysql ใช้ ? mssql ใช้ @id
	row := db.QueryRow(query, sql.Named("id", id))                                       //ถ้ามี ? หลายตัว ก็ใส่ พารามิเตอร์ตาม index ไปเรื่อยๆ
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
	query := "INSERT INTO USERS (USER_ID,USERNAME,PASSWORD,USER_ROLE) VALUES (@user_id,@username,@password,@user_role)"
	result, err := db.Exec(query,
		sql.Named("user_id", user.USER_ID),
		sql.Named("username", user.USERNAME),
		sql.Named("password", user.PASSWORD),
		sql.Named("user_role", user.USER_ROLE))
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
	query := "UPDATE USERS SET USERNAME=@username, PASSWORD=@password, USER_ROLE=@user_role WHERE USER_ID=@user_id"
	result, err := db.Exec(query,
		sql.Named("username", user.USERNAME),
		sql.Named("password", user.PASSWORD),
		sql.Named("user_role", user.USER_ROLE),
		sql.Named("user_id", user.USER_ID))
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

	query := "DELETE FROM USERS WHERE USER_ID=@id"
	result, err := db.Exec(query, sql.Named("id", id))
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
