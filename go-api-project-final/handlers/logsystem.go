package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mandarinkb/go-api-project-final/middleware"
	"github.com/mandarinkb/go-api-project-final/service"
	"github.com/mandarinkb/go-api-project-final/utils"
)

type logHandler struct {
	logSrv service.LogSystemService
}

func NewLogHanler(logSrv service.LogSystemService) logHandler {
	return logHandler{logSrv}
}

func (l logHandler) Read(c *gin.Context) {
	// config zap log
	logger, err := utils.LogConf()
	if err != nil {
		logger.Error(err.Error(), utils.Url(c.Request.URL.Path),
			utils.User(middleware.Username), utils.Type(utils.TypeWeb))
	}
	// close log
	defer logger.Sync()

	var reqBody service.LogReq
	// แปลงค่าจาก body payload เป็น struct
	err = c.BindJSON(&reqBody)
	if err != nil {
		logger.Error(err.Error(), utils.Url(c.Request.URL.Path),
			utils.User(middleware.Username), utils.Type(utils.TypeWeb))
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	logs, err := l.logSrv.GetLogs(reqBody)
	if err != nil {
		logger.Error(err.Error(), utils.Url(c.Request.URL.Path),
			utils.User(middleware.Username), utils.Type(utils.TypeWeb))
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	logger.Info("read logs", utils.Url(c.Request.URL.Path),
		utils.User(middleware.Username), utils.Type(utils.TypeWeb))
	c.IndentedJSON(http.StatusOK, logs)
}
