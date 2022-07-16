package handlers

import (
	"encoding/json"
	"strconv"

	"github.com/mandarinkb/go-rest-api-with-fasthttp-and-fasthttp-routing/middleware"
	"github.com/mandarinkb/go-rest-api-with-fasthttp-and-fasthttp-routing/repository"
	"github.com/mandarinkb/go-rest-api-with-fasthttp-and-fasthttp-routing/service"
	routing "github.com/qiangxue/fasthttp-routing"
)

type response struct {
	Message string `json:"message"`
}

func message(msg string) []byte {
	message := response{
		Message: msg,
	}
	resp, _ := json.Marshal(message)
	return resp
}

type userHandler struct {
	userServ service.UserService
}

func NewUserHandler(userServ service.UserService) userHandler {

	return userHandler{userServ}
}

func (h userHandler) Authenticate(c *routing.Context) error {
	var reqBody repository.User
	err := json.Unmarshal(c.PostBody(), &reqBody)
	if err != nil {
		return err
	}
	user, err := h.userServ.Authenticate(reqBody.Username, reqBody.Password)
	if err != nil {
		return err
	}

	c.Response.Header.Set("Content-Type", "application/json")
	c.SetStatusCode(200)
	c.Write(message("auth success"))
	middleware.SetCookie(c, user.Username)
	return nil
}

func (h userHandler) Logout(c *routing.Context) error {
	c.Response.Header.Set("Content-Type", "application/json")
	c.SetStatusCode(200)
	c.Write(message("logout success"))
	return nil
}

func (h userHandler) Read(c *routing.Context) error {
	users, err := h.userServ.Read()
	if err != nil {
		return err
	}

	resp, err := json.Marshal(users)
	if err != nil {
		return err
	}

	c.Response.Header.Set("Content-Type", "application/json")
	c.SetStatusCode(200)
	c.Write(resp)
	return nil
}

func (h userHandler) ReadById(c *routing.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return err
	}
	user, err := h.userServ.ReadById(id)
	if err != nil {
		return err
	}
	resp, err := json.Marshal(user)
	if err != nil {
		return err
	}
	c.Response.Header.Set("Content-Type", "application/json")
	c.SetStatusCode(200)
	c.Write(resp)
	return nil
}
func (h userHandler) Create(c *routing.Context) error {
	var reqBody repository.User
	err := json.Unmarshal(c.PostBody(), &reqBody)
	if err != nil {
		return err
	}

	err = h.userServ.Create(reqBody)
	if err != nil {
		return err
	}
	c.Response.Header.Set("Content-Type", "application/json")
	c.SetStatusCode(201)
	c.Write(message("creted success"))
	return nil

}

func (h userHandler) Update(c *routing.Context) error {
	var reqBody repository.User
	err := json.Unmarshal(c.PostBody(), &reqBody)
	if err != nil {
		return err
	}

	err = h.userServ.Update(reqBody)
	if err != nil {
		return err
	}
	c.Response.Header.Set("Content-Type", "application/json")
	c.SetStatusCode(201)
	c.Write(message("update success"))
	return nil
}

func (h userHandler) Delete(c *routing.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return err
	}
	err = h.userServ.Delete(id)
	if err != nil {
		return err
	}
	c.Response.Header.Set("Content-Type", "application/json")
	c.SetStatusCode(200)
	c.Write(message("deleted success"))
	return nil
}
