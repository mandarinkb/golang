package repository

import (
	"database/sql"
	"errors"
)

type userRepository struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) UserRepository {
	return userRepository{db: db}
}

func (r userRepository) Read() ([]User, error) {
	// ตรวจสอบว่า server ได้เปิดอยู่หรือไม่
	err := r.db.Ping()
	if err != nil {
		return nil, err
	}

	query := "SELECT USER_ID, USERNAME, PASSWORD, USER_ROLE FROM USERS"
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

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

func (r userRepository) ReadById(id int) (*User, error) {
	err := r.db.Ping()
	if err != nil {
		return nil, err
	}

	query := "SELECT USER_ID, USERNAME, PASSWORD, USER_ROLE FROM USERS WHERE USER_ID=?" //mysql ใช้ ?
	row := r.db.QueryRow(query, id)                                                     //ถ้ามี ? หลายตัว ก็ใส่ พารามิเตอร์ตาม index ไปเรื่อยๆ
	user := User{}
	err = row.Scan(&user.USER_ID, &user.USERNAME, &user.PASSWORD, &user.USER_ROLE)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r userRepository) Create(user User) error {
	err := r.db.Ping()
	if err != nil {
		return err
	}
	query := "INSERT INTO USERS (USER_ID,USERNAME,PASSWORD,USER_ROLE) VALUES (?,?,?,?)"
	result, err := r.db.Exec(query, user.USER_ID, user.USERNAME, user.PASSWORD, user.USER_ROLE)
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

func (r userRepository) Update(user User) error {
	err := r.db.Ping()
	if err != nil {
		return err
	}
	query := "UPDATE USERS SET USERNAME=?, PASSWORD=?, USER_ROLE=? WHERE USER_ID=?"
	result, err := r.db.Exec(query, user.USERNAME, user.PASSWORD, user.USER_ROLE, user.USER_ID)
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

func (r userRepository) Delete(id int) error {
	err := r.db.Ping()
	if err != nil {
		return err
	}
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
