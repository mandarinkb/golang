package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mandarinkb/go-api-project-final/middleware"
	"github.com/mandarinkb/go-api-project-final/service"
	"github.com/mandarinkb/go-api-project-final/utils"
)

type webHandler struct {
	webServ service.WebService
}

func NewWebHandler(webServ service.WebService) webHandler {
	return webHandler{webServ}
}

func (w webHandler) Read(c *gin.Context) {
	// config zap log
	logger, err := utils.LogConf()
	if err != nil {
		logger.Error(err.Error(), utils.Url(c.Request.URL.Path),
			utils.User(middleware.Username), utils.Type(utils.TypeWeb))
	}
	// close log
	defer logger.Sync()

	web, err := w.webServ.Read()
	if err != nil {
		logger.Error(err.Error(), utils.Url(c.Request.URL.Path),
			utils.User(middleware.Username), utils.Type(utils.TypeWeb))
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	logger.Info("read all web", utils.Url(c.Request.URL.Path),
		utils.User(middleware.Username), utils.Type(utils.TypeWeb))
	c.IndentedJSON(http.StatusOK, web)
}

func (w webHandler) ReadById(c *gin.Context) {
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
	web, err := w.webServ.ReadById(id)
	if err != nil {
		logger.Error(err.Error(), utils.Url(c.Request.URL.Path),
			utils.User(middleware.Username), utils.Type(utils.TypeWeb))
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	logger.Info("read web "+web.WebName, utils.Url(c.Request.URL.Path),
		utils.User(middleware.Username), utils.Type(utils.TypeWeb))
	c.IndentedJSON(http.StatusOK, web)
}

func (w webHandler) Create(c *gin.Context) {
	// config zap log
	logger, err := utils.LogConf()
	if err != nil {
		logger.Error(err.Error(), utils.Url(c.Request.URL.Path),
			utils.User(middleware.Username), utils.Type(utils.TypeWeb))
	}
	// close log
	defer logger.Sync()

	var reqBody service.Web
	// แปลงค่าจาก body payload เป็น struct
	err = c.BindJSON(&reqBody)
	if err != nil {
		logger.Error(err.Error(), utils.Url(c.Request.URL.Path),
			utils.User(middleware.Username), utils.Type(utils.TypeWeb))
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	err = w.webServ.Create(reqBody)
	if err != nil {
		logger.Error(err.Error(), utils.Url(c.Request.URL.Path),
			utils.User(middleware.Username), utils.Type(utils.TypeWeb))
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	logger.Info("create web "+reqBody.WebName, utils.Url(c.Request.URL.Path),
		utils.User(middleware.Username), utils.Type(utils.TypeWeb))
	c.IndentedJSON(http.StatusOK, gin.H{"message": "create web success"})
}

func (w webHandler) Update(c *gin.Context) {
	// config zap log
	logger, err := utils.LogConf()
	if err != nil {
		logger.Error(err.Error(), utils.Url(c.Request.URL.Path),
			utils.User(middleware.Username), utils.Type(utils.TypeWeb))
	}
	// close log
	defer logger.Sync()

	var reqBody service.Web
	// แปลงค่าจาก body payload เป็น struct
	err = c.BindJSON(&reqBody)
	if err != nil {
		logger.Error(err.Error(), utils.Url(c.Request.URL.Path),
			utils.User(middleware.Username), utils.Type(utils.TypeWeb))
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	err = w.webServ.Update(reqBody)
	if err != nil {
		logger.Error(err.Error(), utils.Url(c.Request.URL.Path),
			utils.User(middleware.Username), utils.Type(utils.TypeWeb))
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	logger.Info("update web "+reqBody.WebName, utils.Url(c.Request.URL.Path),
		utils.User(middleware.Username), utils.Type(utils.TypeWeb))
	c.IndentedJSON(http.StatusOK, gin.H{"message": "update web success"})
}

func (w webHandler) UpdateStatus(c *gin.Context) {
	// config zap log
	logger, err := utils.LogConf()
	if err != nil {
		logger.Error(err.Error(), utils.Url(c.Request.URL.Path),
			utils.User(middleware.Username), utils.Type(utils.TypeWeb))
	}
	// close log
	defer logger.Sync()

	var reqBody service.Web
	// แปลงค่าจาก body payload เป็น struct
	err = c.BindJSON(&reqBody)
	if err != nil {
		logger.Error(err.Error(), utils.Url(c.Request.URL.Path),
			utils.User(middleware.Username), utils.Type(utils.TypeWeb))
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	err = w.webServ.UpdateStatus(reqBody)
	if err != nil {
		logger.Error(err.Error(), utils.Url(c.Request.URL.Path),
			utils.User(middleware.Username), utils.Type(utils.TypeWeb))
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	logger.Info("update web status "+reqBody.WebName, utils.Url(c.Request.URL.Path),
		utils.User(middleware.Username), utils.Type(utils.TypeWeb))
	c.IndentedJSON(http.StatusOK, gin.H{"message": "update web success"})
}

func (w webHandler) Delete(c *gin.Context) {
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
	webName, err := w.webServ.Delete(id)
	if err != nil {
		logger.Error(err.Error(), utils.Url(c.Request.URL.Path),
			utils.User(middleware.Username), utils.Type(utils.TypeWeb))
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	logger.Info("delete web "+webName, utils.Url(c.Request.URL.Path),
		utils.User(middleware.Username), utils.Type(utils.TypeWeb))
	c.IndentedJSON(http.StatusOK, gin.H{"message": "delete web success"})
}
