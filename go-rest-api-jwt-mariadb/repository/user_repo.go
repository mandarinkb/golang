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

func (r userRepository) Authenticate(username string) (*User, error) {
	err := r.db.Ping()
	if err != nil {
		return nil, err
	}
	// mysql ใช้ ?
	// ถ้ามี ? หลายตัว ก็ใส่ พารามิเตอร์ตาม index ไปเรื่อยๆ
	query := "SELECT USER_ID, USERNAME, PASSWORD, USER_ROLE FROM USERS WHERE USERNAME=?"
	row := r.db.QueryRow(query, username)
	user := User{}
	err = row.Scan(&user.UserId, &user.Username, &user.Password, &user.UserRole)
	if err != nil {
		return nil, err
	}

	return &user, nil
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
		//จะเรียงตามชื่อ field ที่ query
		err = rows.Scan(&user.UserId, &user.Username, &user.Password, &user.UserRole)
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
	// mysql ใช้ ?
	// ถ้ามี ? หลายตัว ก็ใส่ พารามิเตอร์ตาม index ไปเรื่อยๆ
	query := "SELECT USER_ID, USERNAME, PASSWORD, USER_ROLE FROM USERS WHERE USER_ID=?"
	row := r.db.QueryRow(query, id)
	user := User{}
	err = row.Scan(&user.UserId, &user.Username, &user.Password, &user.UserRole)
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
	result, err := r.db.Exec(query, user.UserId, user.Username, user.Password, user.UserRole)
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
	result, err := r.db.Exec(query, user.Username, user.Password, user.UserRole, user.UserId)
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
