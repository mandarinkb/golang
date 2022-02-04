package middleware

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mandarinkb/go-api-project-final/database"
	"github.com/mandarinkb/go-api-project-final/repository"
	"github.com/mandarinkb/go-api-project-final/utils"
)

// เอาไว้ใช้ใน log
var Username string

func JWTAuth(c *gin.Context) {
	// กำหนด path ที่ไม่ต้องทำการ authenticate
	permitPath := NewPermitPathConfig(c).Path("/v1/authenticate", "/v1/token/refresh")
	// กรณีที่เซ็ต path ที่ไม่ต้อง authenticate ไว้
	// ไปทำคำสั่ง handler func อื่นต่อได้เลย
	if permitPath {
		c.Next()
	} else {
		// config zap log
		logger, err := utils.LogConf()
		if err != nil {
			logger.Error(err.Error(), utils.Url(c.Request.URL.Path),
				utils.User("-"), utils.Type(utils.TypeWeb))
		}
		// close log
		defer logger.Sync()

		// ต้องทำการตรวจสอบ token ก่อน ถึงจะทำคำสั่ง handler func อื่นต่อ
		// ดึง token จาก header
		token, err := GetToken(c.Request)
		if err != nil {
			logger.Error(err.Error(), utils.Url(c.Request.URL.Path),
				utils.User("-"), utils.Type(utils.TypeWeb))
			c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
			// ใช้คำสั่ง c.Abort() จะหยุดการเรียก HandlerFunc อื่นต่อจากนี้
			c.Abort()
			return
		}

		// ดึง username ที่อยู่ใน token เพื่อไว้เก็บ log
		claimsDetail, err := GetClaimsToken(token)
		if err != nil {
			logger.Error(err.Error(), utils.Url(c.Request.URL.Path),
				utils.User("-"), utils.Type(utils.TypeWeb))
			c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
			c.Abort()
			return
		}
		Username = claimsDetail.Subject

		// ตรวจสอบ token
		isToken, err := verifyAccessToken(token)
		if err != nil {
			logger.Error(err.Error(), utils.Url(c.Request.URL.Path),
				utils.User(Username), utils.Type(utils.TypeWeb))
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
	// config zap log
	logger, err := utils.LogConf()
	if err != nil {
		logger.Error(err.Error(), utils.Url(c.Request.URL.Path),
			utils.User("-"), utils.Type(utils.TypeWeb))
	}
	// close log
	defer logger.Sync()

	// กำหนด type ของ key value
	mapToken := map[string]string{}
	if err := c.ShouldBindJSON(&mapToken); err != nil {
		logger.Error(err.Error(), utils.Url(c.Request.URL.Path),
			utils.User("-"), utils.Type(utils.TypeWeb))
		c.IndentedJSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}
	// กำหนดค่า key เป็น refresh_token
	refreshToken := mapToken["refresh_token"]

	// ตรวจสอบ refresh token
	_, err2 := verifyRefreshToken(refreshToken)
	if err2 != nil {
		logger.Error(err2.Error(), utils.Url(c.Request.URL.Path),
			utils.User("-"), utils.Type(utils.TypeWeb))
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	// ดึงค่า claims จาก token
	claimsDetail, err := GetClaimsToken(refreshToken)
	if err != nil {
		logger.Error(err.Error(), utils.Url(c.Request.URL.Path),
			utils.User("-"), utils.Type(utils.TypeWeb))
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
		logger.Error(err.Error(), utils.Url(c.Request.URL.Path),
			utils.User(claimsDetail.Subject), utils.Type(utils.TypeWeb))
		c.IndentedJSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}
	// เซ็ตค่า เพื่อส่ง token แสดงผ่าน api
	resToken := TokenResponse{
		AccessToken:  td.AccessToken,
		RefreshToken: td.RefreshToken,
	}

	// connect redis
	rdb := database.RedisConn()
	defer rdb.Close()
	// ลบ(revoke) access token และ refresh token เก่า ออกจาก redis
	ctx := context.Background()
	rdb.Del(ctx, claimsDetail.RefAccessUuid, claimsDetail.RefreshUuid)

	logger.Info("create new token", utils.Url(c.Request.URL.Path),
		utils.User(claimsDetail.Subject), utils.Type(utils.TypeWeb))
	c.IndentedJSON(http.StatusCreated, resToken)
}
