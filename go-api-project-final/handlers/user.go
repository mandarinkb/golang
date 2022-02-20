package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mandarinkb/go-api-project-final/database"
	"github.com/mandarinkb/go-api-project-final/middleware"
	"github.com/mandarinkb/go-api-project-final/service"
	"github.com/mandarinkb/go-api-project-final/utils"
)

type userHandler struct {
	userSrv service.UserService
}

func NewUserHandler(userSrv service.UserService) userHandler {
	return userHandler{userSrv: userSrv}
}

func (h userHandler) Authenticate(c *gin.Context) {
	// config zap log
	logger, err := utils.LogConf()
	if err != nil {
		logger.Error(err.Error(), utils.Url(c.Request.URL.Path),
			utils.User("-"), utils.Type(utils.TypeWeb))
	}
	// close log
	defer logger.Sync()

	// request body จะถูกเรียกเมื่อ ใช้ ฟังก์ชัน BindJSON โดยค่าใน body จะตรงกับ reqBody
	// reqBody ชื่อ filed ในนี้จะต้องตรงกับ ค่า key ใน body ด้วย
	var reqBody service.UserRequest
	// แปลงค่าจาก body payload เป็น struct
	err = c.BindJSON(&reqBody)
	if err != nil {
		logger.Error(err.Error(), utils.Url(c.Request.URL.Path),
			utils.User("-"), utils.Type(utils.TypeWeb))
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	token, err := h.userSrv.Authenticate(reqBody.Username, reqBody.Password)
	if err != nil {
		logger.Error(err.Error(), utils.Url(c.Request.URL.Path),
			utils.User(reqBody.Username), utils.Type(utils.TypeWeb))
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	logger.Info("login", utils.Url(c.Request.URL.Path),
		utils.User(reqBody.Username), utils.Type(utils.TypeWeb))
	c.IndentedJSON(http.StatusCreated, token)
}

func (h userHandler) Logout(c *gin.Context) {
	// config zap log
	logger, err := utils.LogConf()
	if err != nil {
		logger.Error(err.Error(), utils.Url(c.Request.URL.Path),
			utils.User(middleware.Username), utils.Type(utils.TypeWeb))
	}
	// close log
	defer logger.Sync()

	// รับค่า token from header
	token, err := middleware.GetToken(c.Request)
	if err != nil {
		logger.Error(err.Error(), utils.Url(c.Request.URL.Path),
			utils.User(middleware.Username), utils.Type(utils.TypeWeb))
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	// ดึงค่า claims จาก token
	claimsDetail, err := middleware.GetClaimsToken(token)
	if err != nil {
		logger.Error(err.Error(), utils.Url(c.Request.URL.Path),
			utils.User(claimsDetail.Subject), utils.Type(utils.TypeWeb))
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	// connect redis
	rdb := database.RedisConn()
	defer rdb.Close()

	// ลบ access token และ refresh token
	ctx := context.Background()
	rdb.Del(ctx, claimsDetail.AccessUuid, claimsDetail.RefRefreshUuid)

	logger.Info("logout", utils.Url(c.Request.URL.Path),
		utils.User(claimsDetail.Subject), utils.Type(utils.TypeWeb))
	c.IndentedJSON(http.StatusOK, gin.H{"message": "logout successful"})
}

func (h userHandler) ReadUsers(c *gin.Context) {
	// config zap log
	logger, err := utils.LogConf()
	if err != nil {
		logger.Error(err.Error(), utils.Url(c.Request.URL.Path),
			utils.User(middleware.Username), utils.Type(utils.TypeWeb))
	}
	// close log
	defer logger.Sync()

	users, err := h.userSrv.Read()
	if err != nil {
		logger.Error(err.Error(), utils.Url(c.Request.URL.Path),
			utils.User(middleware.Username), utils.Type(utils.TypeWeb))
		// gin.H{} สามารถใส่ข้อความแล้วจะ return json ให้
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	logger.Info("read all user", utils.Url(c.Request.URL.Path),
		utils.User(middleware.Username), utils.Type(utils.TypeWeb))
	// return ค่า http status และค่า json obj หรือ json array ออกให้
	c.IndentedJSON(http.StatusOK, users)
}

func (h userHandler) ReadUserByID(c *gin.Context) {
	// config zap log
	logger, err := utils.LogConf()
	if err != nil {
		logger.Error(err.Error(), utils.Url(c.Request.URL.Path),
			utils.User(middleware.Username), utils.Type(utils.TypeWeb))
	}
	// close log
	defer logger.Sync()

	idStr := c.Param("id")
	// แปลงเป็น int
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logger.Error(err.Error(), utils.Url(c.Request.URL.Path),
			utils.User(middleware.Username), utils.Type(utils.TypeWeb))
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	user, err := h.userSrv.ReadById(id)
	if err != nil {
		logger.Error(err.Error(), utils.Url(c.Request.URL.Path),
			utils.User(middleware.Username), utils.Type(utils.TypeWeb))
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	logger.Info("read user "+user.Username, utils.Url(c.Request.URL.Path),
		utils.User(middleware.Username), utils.Type(utils.TypeWeb))
	c.IndentedJSON(http.StatusOK, user)
}

func (h userHandler) CreateUsers(c *gin.Context) {
	// config zap log
	logger, err := utils.LogConf()
	if err != nil {
		logger.Error(err.Error(), utils.Url(c.Request.URL.Path),
			utils.User(middleware.Username), utils.Type(utils.TypeWeb))
	}
	// close log
	defer logger.Sync()

	// request body จะถูกเรียกเมื่อ ใช้ ฟังก์ชัน BindJSON โดยค่าใน body จะตรงกับ reqBody
	// reqBody ชื่อ filed ในนี้จะต้องตรงกับ ค่า key ใน body ด้วย
	var reqBody service.UserRequest
	// แปลงค่าจาก body payload เป็น struct
	err = c.BindJSON(&reqBody)
	if err != nil {
		logger.Error(err.Error(), utils.Url(c.Request.URL.Path),
			utils.User(middleware.Username), utils.Type(utils.TypeWeb))
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	user, err := h.userSrv.Create(reqBody)
	if err != nil {
		logger.Error(err.Error(), utils.Url(c.Request.URL.Path),
			utils.User(middleware.Username), utils.Type(utils.TypeWeb))
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	logger.Info("create user "+reqBody.Username, utils.Url(c.Request.URL.Path),
		utils.User(middleware.Username), utils.Type(utils.TypeWeb))
	c.IndentedJSON(http.StatusCreated, user)
}

func (h userHandler) UpdateUsers(c *gin.Context) {
	// config zap log
	logger, err := utils.LogConf()
	if err != nil {
		logger.Error(err.Error(), utils.Url(c.Request.URL.Path),
			utils.User(middleware.Username), utils.Type(utils.TypeWeb))
	}
	// close log
	defer logger.Sync()

	var reqBody service.UserRequest
	err = c.BindJSON(&reqBody)
	if err != nil {
		logger.Error(err.Error(), utils.Url(c.Request.URL.Path),
			utils.User(middleware.Username), utils.Type(utils.TypeWeb))
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	user, err := h.userSrv.Update(reqBody)
	if err != nil {
		logger.Error(err.Error(), utils.Url(c.Request.URL.Path),
			utils.User(middleware.Username), utils.Type(utils.TypeWeb))
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	logger.Info("update user "+reqBody.Username, utils.Url(c.Request.URL.Path),
		utils.User(middleware.Username), utils.Type(utils.TypeWeb))
	c.IndentedJSON(http.StatusOK, user)
}

func (h userHandler) DeleteUsers(c *gin.Context) {
	// config zap log
	logger, err := utils.LogConf()
	if err != nil {
		logger.Error(err.Error(), utils.Url(c.Request.URL.Path),
			utils.User(middleware.Username), utils.Type(utils.TypeWeb))
	}
	// close log
	defer logger.Sync()

	idStr := c.Param("id")
	// แปลงเป็น int
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logger.Error(err.Error(), utils.Url(c.Request.URL.Path),
			utils.User(middleware.Username), utils.Type(utils.TypeWeb))
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	username, err := h.userSrv.Delete(id)
	if err != nil {
		logger.Error(err.Error(), utils.Url(c.Request.URL.Path),
			utils.User(middleware.Username), utils.Type(utils.TypeWeb))
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	logger.Info("delete user "+username, utils.Url(c.Request.URL.Path),
		utils.User(middleware.Username), utils.Type(utils.TypeWeb))
	c.IndentedJSON(http.StatusOK, gin.H{"message": "delete user success"})
}
