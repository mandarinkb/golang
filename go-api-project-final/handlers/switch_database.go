package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mandarinkb/go-api-project-final/middleware"
	"github.com/mandarinkb/go-api-project-final/service"
	"github.com/mandarinkb/go-api-project-final/utils"
)

type switchDatabaseHandler struct {
	swDbServ service.SwitchDatabaseService
}

func NewSwitchDabaseHandler(swDbServ service.SwitchDatabaseService) switchDatabaseHandler {
	return switchDatabaseHandler{swDbServ}
}

func (s switchDatabaseHandler) Read(c *gin.Context) {
	// config zap log
	logger, err := utils.LogConf()
	if err != nil {
		logger.Error(err.Error(), utils.Url(c.Request.URL.Path),
			utils.User(middleware.Username), utils.Type(utils.TypeWeb))
	}
	// close log
	defer logger.Sync()

	swDb, err := s.swDbServ.Read()
	if err != nil {
		logger.Error(err.Error(), utils.Url(c.Request.URL.Path),
			utils.User(middleware.Username), utils.Type(utils.TypeWeb))
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	logger.Info("read all switch database", utils.Url(c.Request.URL.Path),
		utils.User(middleware.Username), utils.Type(utils.TypeWeb))
	c.IndentedJSON(http.StatusOK, swDb)
}

func (s switchDatabaseHandler) ReadById(c *gin.Context) {
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
	swDb, err := s.swDbServ.ReadById(id)
	if err != nil {
		logger.Error(err.Error(), utils.Url(c.Request.URL.Path),
			utils.User(middleware.Username), utils.Type(utils.TypeWeb))
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	logger.Info("read switch database "+swDb.DatabaseName, utils.Url(c.Request.URL.Path),
		utils.User(middleware.Username), utils.Type(utils.TypeWeb))
	c.IndentedJSON(http.StatusOK, swDb)
}

func (s switchDatabaseHandler) Create(c *gin.Context) {
	// config zap log
	logger, err := utils.LogConf()
	if err != nil {
		logger.Error(err.Error(), utils.Url(c.Request.URL.Path),
			utils.User(middleware.Username), utils.Type(utils.TypeWeb))
	}
	// close log
	defer logger.Sync()

	var reqBody service.SwitchDatabase
	// แปลงค่าจาก body payload เป็น struct
	err = c.BindJSON(&reqBody)
	if err != nil {
		logger.Error(err.Error(), utils.Url(c.Request.URL.Path),
			utils.User(middleware.Username), utils.Type(utils.TypeWeb))
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	err = s.swDbServ.Create(reqBody)
	if err != nil {
		logger.Error(err.Error(), utils.Url(c.Request.URL.Path),
			utils.User(middleware.Username), utils.Type(utils.TypeWeb))
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	logger.Info("create switch database "+reqBody.DatabaseName, utils.Url(c.Request.URL.Path),
		utils.User(middleware.Username), utils.Type(utils.TypeWeb))
	c.IndentedJSON(http.StatusOK, gin.H{"message": "create schedule success"})
}

func (s switchDatabaseHandler) Update(c *gin.Context) {
	// config zap log
	logger, err := utils.LogConf()
	if err != nil {
		logger.Error(err.Error(), utils.Url(c.Request.URL.Path),
			utils.User(middleware.Username), utils.Type(utils.TypeWeb))
	}
	// close log
	defer logger.Sync()

	var reqBody service.SwitchDatabase
	// แปลงค่าจาก body payload เป็น struct
	err = c.BindJSON(&reqBody)
	if err != nil {
		logger.Error(err.Error(), utils.Url(c.Request.URL.Path),
			utils.User(middleware.Username), utils.Type(utils.TypeWeb))
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	err = s.swDbServ.Update(reqBody)
	if err != nil {
		logger.Error(err.Error(), utils.Url(c.Request.URL.Path),
			utils.User(middleware.Username), utils.Type(utils.TypeWeb))
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	logger.Info("update switch database "+reqBody.DatabaseName, utils.Url(c.Request.URL.Path),
		utils.User(middleware.Username), utils.Type(utils.TypeWeb))
	c.IndentedJSON(http.StatusOK, gin.H{"message": "update switch database success"})
}

func (s switchDatabaseHandler) UpdateStatus(c *gin.Context) {
	// config zap log
	logger, err := utils.LogConf()
	if err != nil {
		logger.Error(err.Error(), utils.Url(c.Request.URL.Path),
			utils.User(middleware.Username), utils.Type(utils.TypeWeb))
	}
	// close log
	defer logger.Sync()

	var reqBody service.SwitchDatabase
	// แปลงค่าจาก body payload เป็น struct
	err = c.BindJSON(&reqBody)
	if err != nil {
		logger.Error(err.Error(), utils.Url(c.Request.URL.Path),
			utils.User(middleware.Username), utils.Type(utils.TypeWeb))
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	err = s.swDbServ.UpdateStatus(reqBody)
	if err != nil {
		logger.Error(err.Error(), utils.Url(c.Request.URL.Path),
			utils.User(middleware.Username), utils.Type(utils.TypeWeb))
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	logger.Info("update switch database status "+reqBody.DatabaseName, utils.Url(c.Request.URL.Path),
		utils.User(middleware.Username), utils.Type(utils.TypeWeb))
	c.IndentedJSON(http.StatusOK, gin.H{"message": "update status switch database success"})
}

func (s switchDatabaseHandler) Delete(c *gin.Context) {
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
	dbName, err := s.swDbServ.Delete(id)
	if err != nil {
		logger.Error(err.Error(), utils.Url(c.Request.URL.Path),
			utils.User(middleware.Username), utils.Type(utils.TypeWeb))
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	logger.Info("delete switch database "+dbName, utils.Url(c.Request.URL.Path),
		utils.User(middleware.Username), utils.Type(utils.TypeWeb))
	c.IndentedJSON(http.StatusOK, gin.H{"message": "delete switch database success"})
}
