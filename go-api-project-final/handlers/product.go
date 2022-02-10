package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mandarinkb/go-api-project-final/service"
	"github.com/mandarinkb/go-api-project-final/utils"
)

var typeApp string = "application"

type productHandler struct {
	productSrv service.ProductService
}

func NewProductHandler(productSrv service.ProductService) productHandler {
	return productHandler{productSrv}
}

func (p productHandler) GetSearch(c *gin.Context) {
	// config zap log
	logger, err := utils.LogConf()
	if err != nil {
		logger.Error(err.Error(), utils.Url(c.Request.URL.Path),
			utils.User("-"), utils.Type(typeApp))
	}
	// close log
	defer logger.Sync()

	from := c.Query("from")
	// request body จะถูกเรียกเมื่อ ใช้ ฟังก์ชัน BindJSON โดยค่าใน body จะตรงกับ reqBody
	// reqBody ชื่อ filed ในนี้จะต้องตรงกับ ค่า key ใน body ด้วย
	var reqBody service.ProductRequest
	// แปลงค่าจาก body payload เป็น struct
	err = c.BindJSON(&reqBody)
	if err != nil {
		logger.Error(err.Error(), utils.Url(c.Request.URL.Path),
			utils.User(reqBody.UserId), utils.Type(typeApp))
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	searchData, err := p.productSrv.GetSearch(from, reqBody.Name)
	if err != nil {
		logger.Error(err.Error(), utils.Url(c.Request.URL.Path),
			utils.User(reqBody.UserId), utils.Type(typeApp))
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	logger.Info("[search] "+reqBody.Name, utils.Url(c.Request.URL.Path),
		utils.User(reqBody.UserId), utils.Type(typeApp))
	c.IndentedJSON(http.StatusOK, searchData)
}
func (p productHandler) GetFilterSearch(c *gin.Context) {
	// config zap log
	logger, err := utils.LogConf()
	if err != nil {
		logger.Error(err.Error(), utils.Url(c.Request.URL.Path),
			utils.User("-"), utils.Type(typeApp))
	}
	// close log
	defer logger.Sync()

	from := c.Query("from")
	// request body จะถูกเรียกเมื่อ ใช้ ฟังก์ชัน BindJSON โดยค่าใน body จะตรงกับ reqBody
	// reqBody ชื่อ filed ในนี้จะต้องตรงกับ ค่า key ใน body ด้วย
	var reqBody service.ProductRequest
	// แปลงค่าจาก body payload เป็น struct
	err = c.BindJSON(&reqBody)
	if err != nil {
		logger.Error(err.Error(), utils.Url(c.Request.URL.Path),
			utils.User(reqBody.UserId), utils.Type(typeApp))
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	filterSearchData, err := p.productSrv.
		GetFilterSearch(from, reqBody.Name, reqBody.WebName, reqBody.MinPrice, reqBody.MaxPrice)
	if err != nil {
		logger.Error(err.Error(), utils.Url(c.Request.URL.Path),
			utils.User(reqBody.UserId), utils.Type(typeApp))
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	logger.Info("[filter search] "+reqBody.Name, utils.Url(c.Request.URL.Path),
		utils.User(reqBody.UserId), utils.Type(typeApp))
	c.IndentedJSON(http.StatusOK, filterSearchData)
}
func (p productHandler) GetCategory(c *gin.Context) {
	// config zap log
	logger, err := utils.LogConf()
	if err != nil {
		logger.Error(err.Error(), utils.Url(c.Request.URL.Path),
			utils.User("-"), utils.Type(typeApp))
	}
	// close log
	defer logger.Sync()

	from := c.Query("from")
	// request body จะถูกเรียกเมื่อ ใช้ ฟังก์ชัน BindJSON โดยค่าใน body จะตรงกับ reqBody
	// reqBody ชื่อ filed ในนี้จะต้องตรงกับ ค่า key ใน body ด้วย
	var reqBody service.ProductRequest
	// แปลงค่าจาก body payload เป็น struct
	err = c.BindJSON(&reqBody)
	if err != nil {
		logger.Error(err.Error(), utils.Url(c.Request.URL.Path),
			utils.User(reqBody.UserId), utils.Type(typeApp))
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	categoryData, err := p.productSrv.GetCategory(from, reqBody.Category, reqBody.WebName)
	if err != nil {
		logger.Error(err.Error(), utils.Url(c.Request.URL.Path),
			utils.User(reqBody.UserId), utils.Type(typeApp))
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	logger.Info("[category] "+reqBody.Category, utils.Url(c.Request.URL.Path),
		utils.User(reqBody.UserId), utils.Type(typeApp))
	c.IndentedJSON(http.StatusOK, categoryData)
}
func (p productHandler) GetHistory(c *gin.Context) {
	// config zap log
	logger, err := utils.LogConf()
	if err != nil {
		logger.Error(err.Error(), utils.Url(c.Request.URL.Path),
			utils.User("-"), utils.Type(typeApp))
	}
	// close log
	defer logger.Sync()
	from := c.Query("from")
	// request body จะถูกเรียกเมื่อ ใช้ ฟังก์ชัน BindJSON โดยค่าใน body จะตรงกับ reqBody
	// reqBody ชื่อ filed ในนี้จะต้องตรงกับ ค่า key ใน body ด้วย
	var reqBody service.ProductRequest
	// แปลงค่าจาก body payload เป็น struct
	err = c.BindJSON(&reqBody)
	if err != nil {
		logger.Error(err.Error(), utils.Url(c.Request.URL.Path),
			utils.User(reqBody.UserId), utils.Type(typeApp))
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	historyData, err := p.productSrv.GetHistory(from, reqBody.History)
	if err != nil {
		logger.Error(err.Error(), utils.Url(c.Request.URL.Path),
			utils.User(reqBody.UserId), utils.Type(typeApp))
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, historyData)
}
func (p productHandler) GetWebName(c *gin.Context) {
	// config zap log
	logger, err := utils.LogConf()
	if err != nil {
		logger.Error(err.Error(), utils.Url(c.Request.URL.Path),
			utils.User("-"), utils.Type(typeApp))
	}
	// close log
	defer logger.Sync()
	webNameData, err := p.productSrv.GetWebName()
	if err != nil {
		logger.Error(err.Error(), utils.Url(c.Request.URL.Path),
			utils.User("-"), utils.Type(typeApp))
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, webNameData)
}
