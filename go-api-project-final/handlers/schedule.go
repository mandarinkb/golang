package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mandarinkb/go-api-project-final/middleware"
	"github.com/mandarinkb/go-api-project-final/service"
	"github.com/mandarinkb/go-api-project-final/utils"
)

type scheduleHandler struct {
	scheduleServ service.ScheduleService
}

func NewScheduleHandler(scheduleServ service.ScheduleService) scheduleHandler {
	return scheduleHandler{scheduleServ}
}

func (s scheduleHandler) Read(c *gin.Context) {
	// config zap log
	logger, err := utils.LogConf()
	if err != nil {
		logger.Error(err.Error(), utils.Url(c.Request.URL.Path),
			utils.User(middleware.Username), utils.Type(utils.TypeWeb))
	}
	// close log
	defer logger.Sync()

	schedule, err := s.scheduleServ.Read()
	if err != nil {
		logger.Error(err.Error(), utils.Url(c.Request.URL.Path),
			utils.User(middleware.Username), utils.Type(utils.TypeWeb))
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	logger.Info("read all schedule", utils.Url(c.Request.URL.Path),
		utils.User(middleware.Username), utils.Type(utils.TypeWeb))
	c.IndentedJSON(http.StatusOK, schedule)
}

func (s scheduleHandler) ReadById(c *gin.Context) {
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
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	schedule, err := s.scheduleServ.ReadById(id)
	if err != nil {
		logger.Error(err.Error(), utils.Url(c.Request.URL.Path),
			utils.User(middleware.Username), utils.Type(utils.TypeWeb))
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	logger.Info("read schedule "+schedule.ScheduleName, utils.Url(c.Request.URL.Path),
		utils.User(middleware.Username), utils.Type(utils.TypeWeb))
	c.IndentedJSON(http.StatusOK, schedule)
}

func (s scheduleHandler) Create(c *gin.Context) {
	// config zap log
	logger, err := utils.LogConf()
	if err != nil {
		logger.Error(err.Error(), utils.Url(c.Request.URL.Path),
			utils.User(middleware.Username), utils.Type(utils.TypeWeb))
	}
	// close log
	defer logger.Sync()

	var reqBody service.Schedule
	// แปลงค่าจาก body payload เป็น struct
	err = c.BindJSON(&reqBody)
	if err != nil {
		logger.Error(err.Error(), utils.Url(c.Request.URL.Path),
			utils.User(middleware.Username), utils.Type(utils.TypeWeb))
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	err = s.scheduleServ.Create(reqBody)
	if err != nil {
		logger.Error(err.Error(), utils.Url(c.Request.URL.Path),
			utils.User(middleware.Username), utils.Type(utils.TypeWeb))
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	logger.Info("create schedule "+reqBody.ScheduleName, utils.Url(c.Request.URL.Path),
		utils.User(middleware.Username), utils.Type(utils.TypeWeb))
	c.IndentedJSON(http.StatusOK, gin.H{"message": "create schedule success"})
}

func (s scheduleHandler) Update(c *gin.Context) {
	// config zap log
	logger, err := utils.LogConf()
	if err != nil {
		logger.Error(err.Error(), utils.Url(c.Request.URL.Path),
			utils.User(middleware.Username), utils.Type(utils.TypeWeb))
	}
	// close log
	defer logger.Sync()

	var reqBody service.Schedule
	// แปลงค่าจาก body payload เป็น struct
	err = c.BindJSON(&reqBody)
	if err != nil {
		logger.Error(err.Error(), utils.Url(c.Request.URL.Path),
			utils.User(middleware.Username), utils.Type(utils.TypeWeb))
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	err = s.scheduleServ.Update(reqBody)
	if err != nil {
		logger.Error(err.Error(), utils.Url(c.Request.URL.Path),
			utils.User(middleware.Username), utils.Type(utils.TypeWeb))
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	logger.Info("update schedule "+reqBody.ScheduleName, utils.Url(c.Request.URL.Path),
		utils.User(middleware.Username), utils.Type(utils.TypeWeb))
	c.IndentedJSON(http.StatusOK, gin.H{"message": "update schedule success"})
}

func (s scheduleHandler) Delete(c *gin.Context) {
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
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	scheduleName, err := s.scheduleServ.Delete(id)
	if err != nil {
		logger.Error(err.Error(), utils.Url(c.Request.URL.Path),
			utils.User(middleware.Username), utils.Type(utils.TypeWeb))
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	logger.Info("delete schedule "+scheduleName, utils.Url(c.Request.URL.Path),
		utils.User(middleware.Username), utils.Type(utils.TypeWeb))
	c.IndentedJSON(http.StatusOK, gin.H{"message": "delete schedule success"})
}
