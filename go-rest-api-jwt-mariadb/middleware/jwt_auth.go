package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mandarinkb/go-rest-api-jwt-mariadb/database"
	"github.com/mandarinkb/go-rest-api-jwt-mariadb/repository"
)

func JWTAuth(c *gin.Context) {
	// กำหนด path ที่ไม่ต้องทำการ authenticate
	permitPath := NewPermitPathConfig(c).Path("/v1/authenticate", "/v1/token/refresh")
	// กรณีที่เซ็ต path ที่ไม่ต้อง authenticate ไว้
	// ไปทำคำสั่ง handler func อื่นต่อได้เลย
	if permitPath {
		c.Next()
	} else {
		// ต้องทำการตรวจสอบ token ก่อน ถึงจะทำคำสั่ง handler func อื่นต่อ
		// ดึง token จาก header
		token, err := GetToken(c.Request)
		if err != nil {
			c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
			// ใช้คำสั่ง c.Abort() จะหยุดการเรียก HandlerFunc อื่นต่อจากนี้
			c.Abort()
			return
		}
		// ตรวจสอบ token
		isToken, err := verifyAccessToken(token)
		if err != nil {
			c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
			c.Abort()
			return
		}
		// เมื่อ token ถูกต้องจะเรียก Handler อื่นต่อจากนี้
		if isToken {
			c.Next()
		}

	}
}

func JWTRefresh(c *gin.Context) {
	// กำหนด type ของ key value
	mapToken := map[string]string{}
	if err := c.ShouldBindJSON(&mapToken); err != nil {
		c.IndentedJSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}
	// กำหนดค่า key เป็น refresh_token
	refreshToken := mapToken["refresh_token"]

	// ตรวจสอบ refresh token
	_, err := verifyRefreshToken(refreshToken)
	if err != nil {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	// ดึงค่า claims จาก token
	claimsDetail, err := GetClaimsToken(refreshToken)
	if err != nil {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}
	// นำค่า claims เดิมมาเพื่อสร้าง token ใหม่อีกรอบ
	user := repository.User{
		UserId:   int(claimsDetail.Id),
		Username: claimsDetail.Subject,
		UserRole: claimsDetail.Roles,
	}
	// สร้าง token ขึ้นมาใหม่
	td, err := GenerateToken(user)
	if err != nil {
		c.IndentedJSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}
	// เซ็ตค่า เพื่อส่ง token แสดงผ่าน api
	resToken := TokenResponse{
		AccessToken:  td.AccessToken,
		RefreshToken: td.RefreshToken,
	}

	fmt.Println("delete old access  token: ", claimsDetail.RefAccessUuid)
	fmt.Println("delete old refresh token: ", claimsDetail.RefreshUuid)

	// connect redis
	rdb := database.RedisConn()
	defer rdb.Close()
	// ลบ(revoke) access token และ refresh token เก่า ออกจาก redis
	ctx := context.Background()
	rdb.Del(ctx, claimsDetail.RefAccessUuid, claimsDetail.RefreshUuid)

	c.IndentedJSON(http.StatusCreated, resToken)

}
