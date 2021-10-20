package middleware

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/mandarinkb/go-rest-api-jwt-mariadb/database"
	"github.com/mandarinkb/go-rest-api-jwt-mariadb/repository"
)

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	AccessUuid   string
	RefreshUuid  string
	AtExpires    int64
	RtExpires    int64
}

type ClaimsDetails struct {
	Subject     string  `json:"sub"`
	Id          float64 `json:"id"`
	Roles       string  `json:"roles"`
	AccessUuid  string  `json:"access_uuid"`
	RefreshUuid string  `json:"refresh_uuid"`
}

type jwtMaker struct {
	secretKey string
}

func NewJWTMaker(secretKey string) jwtMaker {
	return jwtMaker{secretKey: secretKey}
}

var ctx = context.Background()

// สร้าง token
func (maker jwtMaker) GenerateToken(user repository.User) (*TokenDetails, error) {
	td := &TokenDetails{}
	td.AtExpires = time.Now().Add(time.Minute * 15).Unix()
	td.AccessUuid = uuid.New().String()

	td.RtExpires = time.Now().Add(time.Hour * 24).Unix()
	td.RefreshUuid = uuid.New().String()

	// claim ตามมาตรฐาน
	// – iss (issuer) : เว็บหรือบริษัทเจ้าของ token
	// – sub (subject) : subject ของ token (subject เอาไว้สำหรับ authenticate user.)
	// – aud (audience) : ผู้รับ token
	// – exp (expiration time) : เวลาหมดอายุของ token
	// – nbf (not before) : เป็นเวลาที่บอกว่า token จะเริ่มใช้งานได้เมื่อไหร่
	// – iat (issued at) : ใช้เก็บเวลาที่ token นี้เกิดปัญหา
	// – jti (JWT id) : เอาไว้เก็บไอดีของ JWT แต่ละตัว
	// – name (Full name) : เอาไว้เก็บชื่อ

	//Creating Access Token
	// กำหนด claim name ขึ้นมาเอง
	atClaims := jwt.MapClaims{}
	atClaims["sub"] = user.Username
	atClaims["id"] = user.UserId
	atClaims["access_uuid"] = td.AccessUuid
	atClaims["roles"] = user.UserRole
	atClaims["exp"] = td.AtExpires

	// map claim and hash algorithm
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)

	// create token
	accessToken, err := at.SignedString([]byte(maker.secretKey))
	td.AccessToken = accessToken
	if err != nil {
		return nil, err
	}

	//Creating Refresh Token
	rtClaims := jwt.MapClaims{}
	rtClaims["sub"] = user.Username
	rtClaims["id"] = user.UserId
	rtClaims["refresh_uuid"] = td.RefreshUuid
	rtClaims["roles"] = user.UserRole
	rtClaims["exp"] = td.RtExpires

	// map claim and hash algorithm
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)

	// create token
	refreshToken, err := rt.SignedString([]byte(maker.secretKey))
	td.RefreshToken = refreshToken
	if err != nil {
		return nil, err
	}

	// connect redis
	rdb := database.NewDatabase().RedisConn()
	defer rdb.Close()

	// จัดเก็บ access token ลง redis
	atExp := time.Unix(td.AtExpires, 0) //converting Unix to UTC(to Time object)
	now := time.Now()
	err = rdb.Set(ctx, td.AccessUuid, user.UserId, atExp.Sub(now)).Err()
	if err != nil {
		return nil, err
	}

	// จัดเก็บ refresh token ลง redis
	rtExp := time.Unix(td.RtExpires, 0)
	err = rdb.Set(ctx, td.RefreshUuid, user.UserId, rtExp.Sub(now)).Err()
	if err != nil {
		return nil, err
	}

	return td, nil
}

// ตรวจสอบ token
func (maker jwtMaker) VerifyAccessToken(tokenStr string) (bool, error) {
	// นำ secret key ที่ตั้งไว้มาถอดรหัส
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(maker.secretKey), nil
	}

	// นำ keyFunc ที่ได้มาทำการ ถอดรหัส token
	// กรณีเกิด error เช่น token หมดอายุ ก็จะ return ค่าออกไป
	jwtToken, err := jwt.Parse(tokenStr, keyFunc)
	if err != nil {
		return false, err
	}

	// ตรวจสอบว่าใช่ access token หรือไม่
	isAt, err := maker.isAccessToken(jwtToken)
	if err != nil {
		return false, err
	}
	if !isAt {
		return false, errors.New("access token not found")
	}

	// connect redis
	rdb := database.NewDatabase().RedisConn()
	defer rdb.Close()
	// ดึงค่า claims จาก token
	cd, err := maker.GetClaimsToken(tokenStr)
	if err != nil {
		return false, err
	}

	// ตรวจสอบ access uuid ใน redis
	val, err := rdb.Get(ctx, cd.AccessUuid).Result()
	_ = val
	// กรณี เกิดerror แสดงว่าไม่พบ access uuid ซึ่งแปลว่า token ถูกยกเลิกการใช้งานแล้ว
	if err != nil {
		return false, errors.New("access token is revoked")
	}
	return true, nil
}

// get claims ใน Payload จาก token
func (maker jwtMaker) GetClaimsToken(tokenStr string) (*ClaimsDetails, error) {
	// นำ secret key ที่ตั้งไว้มาถอดรหัส
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(maker.secretKey), nil
	}

	// นำ keyFunc ที่ได้มาทำการ ถอดรหัส token
	jwtToken, err := jwt.Parse(tokenStr, keyFunc)
	if err != nil {
		return nil, err
	}

	cl := &ClaimsDetails{}
	claims := jwtToken.Claims.(jwt.MapClaims)
	cl.Subject = claims["sub"].(string)
	cl.Id = claims["id"].(float64)
	cl.Roles = claims["roles"].(string)

	aUuid := claims["access_uuid"]
	if aUuid != nil {
		cl.AccessUuid = claims["access_uuid"].(string)
	}

	rUuid := claims["refresh_uuid"]
	if rUuid != nil {
		cl.RefreshUuid = claims["refresh_uuid"].(string)
	}

	return cl, nil
}

// get access_uuid ใน Payload จาก token
// ใช้ภายในไฟล์นี้ i เลยขึ้นต้นด้วยตัวเล็ก
func (maker jwtMaker) isAccessToken(jwtToken *jwt.Token) (bool, error) {
	// ถอดรหัส access_uuid ใน Payload จาก token
	claims := jwtToken.Claims.(jwt.MapClaims)
	aUuid := claims["access_uuid"]
	// กรณี aUuid == nil แสดงว่า ไม่พบ claims access_uuid
	return aUuid != nil, nil
}

// get refresh_uuid ใน Payload จาก token
func (maker jwtMaker) VerifyRefreshToken(tokenStr string) (bool, error) {
	// นำ secret key ที่ตั้งไว้มาถอดรหัส
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(maker.secretKey), nil
	}

	// นำ keyFunc ที่ได้มาทำการ ถอดรหัส token
	jwtToken, err := jwt.Parse(tokenStr, keyFunc)
	if err != nil {
		return false, err
	}

	// ถอดรหัส refresh_uuid ใน Payload จาก token
	claims := jwtToken.Claims.(jwt.MapClaims)
	rUuid := claims["refresh_uuid"]
	// กรณี rUuid == nil แสดงว่า ไม่พบ claims access_uuid
	return rUuid != nil, nil
}
