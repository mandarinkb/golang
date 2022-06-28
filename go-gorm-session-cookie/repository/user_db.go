package repository

import (
	"gorm.io/gorm"
)

type userORM struct {
	db *gorm.DB
}

func NewUserORM(db *gorm.DB) UserRepository {
	return userORM{db}
}

// รับแค่ username ส่วน password เผื่อเข้ารหัสไว้ จะตรวจสอบ password ใน service แทน
func (r userORM) Authenticate(username string) (*User, error) {
	user := User{}
	tx := r.db.Where("USERNAME=?", username).First(&user)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &user, nil
}
func (r userORM) Read() ([]User, error) {
	users := []User{}
	tx := r.db.Find(&users)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return users, nil
}
func (r userORM) ReadById(id int) (*User, error) {
	user := User{}
	tx := r.db.First(&user, id)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &user, nil
}
func (r userORM) Create(user User) error {
	tx := r.db.Create(&user)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
func (r userORM) Update(user User) error {
	// ส่ง model struct ที่ได้สร้างไว้
	tx := r.db.Model(&User{}).Where("USER_ID=?", user.ID).Updates(user)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
func (r userORM) Delete(id int) error {
	// จะลบโดยอ้างตาม model struct ที่ได้สร้างไว้
	tx := r.db.Delete(&User{}, id)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
