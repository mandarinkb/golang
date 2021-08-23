package repository

import (
	"database/sql"
	"errors"
)

type mysqlRepository struct {
	db *sql.DB
}

func NewMysqlRepo(db *sql.DB) database {
	return mysqlRepository{db: db}
}

func (r mysqlRepository) MariadbRead() ([]User, error) {
	// ตรวจสอบว่า server ได้เปิดอยู่หรือไม่
	err := r.db.Ping()
	if err != nil {
		return nil, err
	}

	query := "SELECT USER_ID, USERNAME,PASSWORD, ROLE FROM USERS"
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, nil
	}
	// defer rows.Close()

	users := []User{}
	for rows.Next() {
		user := User{}
		err = rows.Scan(&user.USER_ID, &user.USERNAME, &user.PASSWORD, &user.ROLE) //จะเรียงตามชื่อ field ที่ query
		if err != nil {
			return nil, nil
		}
		users = append(users, user)
	}

	return users, nil
}

func (r mysqlRepository) MariadbReadById(id int) (*User, error) {
	err := r.db.Ping()
	if err != nil {
		return nil, err
	}
	// mssql
	// query := "SELECT USER_ID, USERNAME, ROLE FROM USERS WHERE USER_ID=@id"
	// row := db.QueryRow(query, sql.Named("id",id))

	query := "SELECT USER_ID, USERNAME,PASSWORD, ROLE FROM USERS WHERE USER_ID=?" //mysql ใช้ ? mssql ใช้ @id
	row := r.db.QueryRow(query, id)                                               //ถ้ามี ? หลายตัว ก็ใส่ พารามิเตอร์ตาม index ไปเรื่อยๆ
	user := User{}
	err = row.Scan(&user.USER_ID, &user.USERNAME, &user.PASSWORD, &user.ROLE)
	if err != nil {
		return nil, nil
	}

	return &user, nil
}

func (r mysqlRepository) MariadbCreate(user User) error {
	query := "INSERT INTO USERS (USERNAME,PASSWORD,ROLE) VALUES (?,?,?)"
	result, err := r.db.Exec(query, user.USERNAME, user.PASSWORD, user.ROLE)
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

func (r mysqlRepository) MariadbUpdate(user User) error {
	query := "UPDATE USERS SET USERNAME=?,PASSWORD=?,ROLE=? WHERE USER_ID=?"
	result, err := r.db.Exec(query, user.USERNAME, user.PASSWORD, user.ROLE, user.USER_ID)
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

func (r mysqlRepository) MariadbDelete(id int) error {
	query := "DELETE FROM USERS WHERE USER_ID=?"
	result, err := r.db.Exec(query, id)
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
