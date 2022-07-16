package handlers

import (
	"encoding/json"
	"strconv"

	"github.com/mandarinkb/test-git/service"
	routing "github.com/qiangxue/fasthttp-routing"
)

type dataHandler struct {
	dataServ service.DataService
}

func NewDataService(dataServ service.DataService) dataHandler {
	return dataHandler{dataServ}
}

func (h dataHandler) GetAll(c *routing.Context) error {
	data, err := h.dataServ.GetAll()
	if err != nil {
		return err
	}
	resp, err := json.Marshal(data)
	if err != nil {
		return err
	}
	c.Response.Header.Set("Content-Type", "application/json")
	c.SetStatusCode(200)
	c.Write(resp)
	return nil
}

func (h dataHandler) GetById(c *routing.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return err
	}
	data, err := h.dataServ.GetById(id)
	if err != nil {
		return err
	}
	resp, err := json.Marshal(data)
	if err != nil {
		return err
	}
	c.Response.Header.Set("Content-Type", "application/json")
	c.SetStatusCode(200)
	c.Write(resp)
	return nil
}
