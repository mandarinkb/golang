package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mandarinkb/go-api-project-final/service"
)

type scheduleHandler struct {
	scheduleServ service.ScheduleService
}

func NewScheduleHandler(scheduleServ service.ScheduleService) scheduleHandler {
	return scheduleHandler{scheduleServ}
}

func (s scheduleHandler) Read(c *gin.Context) {
	schedule, err := s.scheduleServ.Read()
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, schedule)
}

func (s scheduleHandler) ReadById(c *gin.Context) {
	idStr := c.Param("id")
	// แปลงเป็น int
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	schedule, err := s.scheduleServ.ReadById(id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, schedule)
}

func (s scheduleHandler) Create(c *gin.Context) {
	var reqBody service.Schedule
	// แปลงค่าจาก body payload เป็น struct
	err := c.BindJSON(&reqBody)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	err = s.scheduleServ.Create(reqBody)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"message": "create schedule success"})
}

func (s scheduleHandler) Update(c *gin.Context) {
	var reqBody service.Schedule
	// แปลงค่าจาก body payload เป็น struct
	err := c.BindJSON(&reqBody)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	err = s.scheduleServ.Update(reqBody)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"message": "update schedule success"})
}

func (s scheduleHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	// แปลงเป็น int
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	err = s.scheduleServ.Delete(id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"message": "delete schedule success"})
}
