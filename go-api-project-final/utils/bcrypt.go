package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// เข้ารหัส
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

// เช็คตรวจสอบรหัสว่าตรงกับที่รับมาไหม
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	// กรณีไม่มี error แสดงว่ารหัสตรงกัน
	// เป็นการเขียนแบบ shot แบบเต็มคือ if err == nil{return true}else{return false}
	return err == nil
}
