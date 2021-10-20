package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mandarinkb/go-rest-api-jwt-mariadb/middleware"
	"github.com/mandarinkb/go-rest-api-jwt-mariadb/repository"
)

type jwtHandler struct{}

func NewJwtHandler() jwtHandler {
	return jwtHandler{}
}

var secretKey = `cTq46<pSE8o;jD>~,H*an1_>uKj!nc1#S:+K&./_2uAiPr?N&.2c.m|^$HUZj0_`

func (jwtHandler) JWTAuth(c *gin.Context) {
	// กำหนด path ที่ไม่ต้องทำการ authenticate
	permitPath := middleware.NewPermitPathConfig(c).Path("/v1/authenticate", "/v1/token/refresh")
	// กรณีที่เซ็ต path ที่ไม่ต้อง authenticate ไว้
	// ไปทำคำสั่ง handler func อื่นต่อได้เลย
	if permitPath {
		c.Next()
	} else {
		// ต้องทำการตรวจสอบ token ก่อน ถึงจะทำคำสั่ง handler func อื่นต่อ
		tokenStr := c.Request.Header.Get("Authorization")
		if len(tokenStr) == 0 {
			c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
			// ใช้คำสั่ง c.Abort() จะหยุดการเรียก HandlerFunc อื่นต่อจากนี้
			c.Abort()
			return
		}

		if strings.HasPrefix(tokenStr, "Bearer ") {
			tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")
			isToken, err := middleware.NewJWTMaker(secretKey).VerifyAccessToken(tokenStr)
			if err != nil {
				c.IndentedJSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
				c.Abort()
				return
			}
			// เมื่อ token ถูกต้องจะเรียก Handler อื่นต่อจากนี้
			if isToken {
				c.Next()
			}
		}
	}
}

func (jwtHandler) JWTRefresh(c *gin.Context) {
	// กำหนด type ของ key value
	mapToken := map[string]string{}
	if err := c.ShouldBindJSON(&mapToken); err != nil {
		c.IndentedJSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}
	// กำหนดค่า key เป็น refresh_token
	refreshToken := mapToken["refresh_token"]

	// ตรวจสอบ refresh token
	isRt, err := middleware.NewJWTMaker(secretKey).VerifyRefreshToken(refreshToken)
	if err != nil {
		c.IndentedJSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}
	// กรณีไม่ใช่ refresh token
	if !isRt {
		c.IndentedJSON(http.StatusUnprocessableEntity, gin.H{"message": "refresh token not found"})
		return
	}

	claimsDetail, err := middleware.NewJWTMaker(secretKey).GetClaimsToken(refreshToken)
	if err != nil {
		c.IndentedJSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}
	user := repository.User{
		UserId:   int(claimsDetail.Id),
		Username: claimsDetail.Subject,
		UserRole: claimsDetail.Roles,
	}

	td, err := middleware.NewJWTMaker(secretKey).GenerateToken(user)
	if err != nil {
		c.IndentedJSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}
	resToken := middleware.TokenResponse{
		AccessToken:  td.AccessToken,
		RefreshToken: td.RefreshToken,
	}

	c.IndentedJSON(http.StatusCreated, resToken)

}
