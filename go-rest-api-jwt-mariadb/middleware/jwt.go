package middleware

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/mandarinkb/go-rest-api-jwt-mariadb/repository"
)

type TokenResponse struct {
	Token string `json:"token"`
}

type jwtMaker struct {
	secretKey string
}

func NewJWTMaker(secretKey string) jwtMaker {
	return jwtMaker{secretKey: secretKey}
}

// สร้าง token
func (maker jwtMaker) GenerateToken(user repository.User) (string, error) {
	// claim ตามมาตรฐาน
	// – iss (issuer) : เว็บหรือบริษัทเจ้าของ token
	// – sub (subject) : subject ของ token (subject เอาไว้สำหรับ authenticate user.)
	// – aud (audience) : ผู้รับ token
	// – exp (expiration time) : เวลาหมดอายุของ token
	// – nbf (not before) : เป็นเวลาที่บอกว่า token จะเริ่มใช้งานได้เมื่อไหร่
	// – iat (issued at) : ใช้เก็บเวลาที่ token นี้เกิดปัญหา
	// – jti (JWT id) : เอาไว้เก็บไอดีของ JWT แต่ละตัว
	// – name (Full name) : เอาไว้เก็บชื่อ

	// กำหนด claim name ขึ้นมาเอง
	atClaims := jwt.MapClaims{}
	atClaims["sub"] = user.Username
	atClaims["id"] = user.UserId
	atClaims["roles"] = user.UserRole
	atClaims["iat"] = time.Now().Unix()
	atClaims["exp"] = time.Now().Add(time.Minute * 2).Unix()

	// map claim and hash algorithm
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)

	// create token
	token, err := claims.SignedString([]byte(maker.secretKey))
	if err != nil {
		return "", err
	}
	return token, nil
}

// ตรวจสอบ token
func (maker jwtMaker) VerifyToken(tokenStr string) (bool, error) {
	// นำ secret key ที่ตั้งไว้มาถอดรหัส
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		return []byte(maker.secretKey), nil
	}

	// นำ keyFunc ที่ได้มาทำการ ถอดรหัส token
	// กรณีเกิด error token ก็จะ return ค่าออกไป
	_, err := jwt.Parse(tokenStr, keyFunc)
	if err != nil {
		return false, err
	}
	return true, nil
}

// get subject ใน Payload จาก token
func (maker jwtMaker) GetSubjectToken(tokenStr string) (string, error) {
	// นำ secret key ที่ตั้งไว้มาถอดรหัส
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		return []byte(maker.secretKey), nil
	}

	// นำ keyFunc ที่ได้มาทำการ ถอดรหัส token
	jwtToken, err := jwt.Parse(tokenStr, keyFunc)
	if err != nil {
		return "", err
	}
	// ถอดรหัส sub ใน Payload จาก token
	claims := jwtToken.Claims.(jwt.MapClaims)
	var subStr string
	switch sub := claims["sub"].(type) {
	case string:
		subStr = sub
	}
	// กรณีเกิดข้อผิดพลาดถอด sub ไม่ได้
	if subStr == "" {
		return "", errors.New("cann't get subject token")
	}
	return subStr, nil
}
