package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mandarinkb/go-api-project-final/service"
)

type webHandler struct {
	webServ service.WebService
}

func NewWebHandler(webServ service.WebService) webHandler {
	return webHandler{webServ}
}

func (w webHandler) Read(c *gin.Context) {
	web, err := w.webServ.Read()
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, web)
}

func (w webHandler) ReadById(c *gin.Context) {
	idStr := c.Param("id")
	// แปลงเป็น int
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	web, err := w.webServ.ReadById(id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, web)
}

func (w webHandler) Create(c *gin.Context) {
	var reqBody service.Web
	// แปลงค่าจาก body payload เป็น struct
	err := c.BindJSON(&reqBody)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	err = w.webServ.Create(reqBody)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"message": "create web success"})
}

func (w webHandler) Update(c *gin.Context) {
	var reqBody service.Web
	// แปลงค่าจาก body payload เป็น struct
	err := c.BindJSON(&reqBody)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	err = w.webServ.Update(reqBody)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"message": "update web success"})
}

func (w webHandler) UpdateStatus(c *gin.Context) {
	var reqBody service.Web
	// แปลงค่าจาก body payload เป็น struct
	err := c.BindJSON(&reqBody)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	err = w.webServ.UpdateStatus(reqBody)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"message": "update web success"})
}

func (w webHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	// แปลงเป็น int
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	err = w.webServ.Delete(id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"message": "delete web success"})
}
