package middleware

import (
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/mandarinkb/go-rest-api-jwt-mariadb/repository"
)

// var (
// 	ErrInvalidToken = errors.New("token is invalid")
// 	ErrExpiredToken = errors.New("token has expired")
// )

type TokenResponse struct {
	Token string `json:"token"`
}

type jwtMaker struct {
	secretKey string
}

func NewJWTMaker(secretKey string) jwtMaker {
	return jwtMaker{secretKey: secretKey}
}

func (maker jwtMaker) GenerateToken(user repository.User) (string, error) {
	// กำหนด claim name ขึ้นมาเอง
	atClaims := jwt.MapClaims{}
	atClaims["sub"] = user.Username
	atClaims["id"] = user.UserId
	atClaims["roles"] = user.UserRole
	atClaims["iat"] = time.Now().Unix()
	atClaims["exp"] = time.Now().Add(time.Minute * 2).Unix()

	// map claim and hash
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)

	// create token
	token, err := claims.SignedString([]byte(maker.secretKey))
	if err != nil {
		return "", err
	}
	return token, nil
}

func (maker jwtMaker) VerifyToken(token string) (bool, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		return []byte(maker.secretKey), nil
	}

	_, err := jwt.Parse(token, keyFunc)
	if err != nil {
		return false, err
	}

	// claims, ok := jwtToken.Claims.(jwt.MapClaims)
	// if !ok {
	// 	// fmt.Println("Can't convert token's claims to standard claims")
	// 	fmt.Println("error 3")
	// 	return false, ErrInvalidToken
	// }

	// var tmExp time.Time
	// switch exp := claims["exp"].(type) {
	// case float64:
	// 	sec, dec := math.Modf(exp)
	// 	tmExp = time.Unix(int64(sec), int64(dec*(1e9)))
	// 	tmExp = time.Unix(int64(exp), 0)
	// case json.Number:
	// 	v, _ := exp.Int64()
	// 	tmExp = time.Unix(v, 0)
	// }

	// // กำหนดเวลาปัจจุบันก่อน
	// tmNow := time.Unix(time.Now().Unix(), 0)
	// // กรณี token หมดอายุ
	// if tmExp.Before(tmNow) {
	// 	fmt.Println("token is expire")
	// 	fmt.Println("error 4")
	// 	return false, ErrExpiredToken
	// }
	return true, nil
}
