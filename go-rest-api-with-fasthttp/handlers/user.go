package handlers

import (
	"encoding/json"
	"strconv"

	"github.com/mandarinkb/go-rest-api-fasthttp/middleware"
	"github.com/mandarinkb/go-rest-api-fasthttp/repository"
	"github.com/mandarinkb/go-rest-api-fasthttp/service"
	"github.com/valyala/fasthttp"
)

type response struct {
	Message string
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

func (h userHandler) Authenticate(ctx *fasthttp.RequestCtx) {
	var reqBody repository.User
	err := json.Unmarshal(ctx.PostBody(), &reqBody)
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}
	user, err := h.userServ.Authenticate(reqBody.Username, reqBody.Password)
	if err != nil {
		// ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		ctx.Response.Header.Set("Content-Type", "application/json")
		ctx.SetStatusCode(500)
		ctx.Write(message("username or password incorrect"))
		return
	}

	ctx.Response.Header.Set("Content-Type", "application/json")
	ctx.SetStatusCode(200)
	ctx.Write(message("auth success"))
	middleware.SetCookie(ctx, user.Username)
}
func (h userHandler) LogOut(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set("Content-Type", "application/json")
	ctx.SetStatusCode(200)
	ctx.Write(message("logout success"))
}

func (h userHandler) Create(ctx *fasthttp.RequestCtx) {
	var reqBody repository.User
	err := json.Unmarshal(ctx.PostBody(), &reqBody)
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}
	err = h.userServ.Create(reqBody)
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}
	ctx.Response.Header.Set("Content-Type", "application/json")
	ctx.SetStatusCode(201)
	ctx.Write(message("creted success"))
}

func (h userHandler) Read(ctx *fasthttp.RequestCtx) {
	users, err := h.userServ.Read()
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}
	resp, err := json.Marshal(users)
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}
	ctx.Response.Header.Set("Content-Type", "application/json")
	ctx.SetStatusCode(200)
	ctx.Write(resp)
}

func (h userHandler) ReadById(ctx *fasthttp.RequestCtx) {
	id := ctx.UserValue("id")
	idStr := id.(string)              // แปลงค่าจาก interfacd{} to string
	idInt, err := strconv.Atoi(idStr) // แปลงค่าจาก string to int
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}
	user, err := h.userServ.ReadById(idInt)
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}
	resp, err := json.Marshal(user)
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}
	ctx.Response.Header.Set("Content-Type", "application/json")
	ctx.SetStatusCode(200)
	ctx.Write(resp)
}

func (h userHandler) Update(ctx *fasthttp.RequestCtx) {
	var reqBody repository.User
	err := json.Unmarshal(ctx.PostBody(), &reqBody)
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}
	err = h.userServ.Update(reqBody)
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}
	ctx.Response.Header.Set("Content-Type", "application/json")
	ctx.SetStatusCode(200)
	ctx.Write(message("update success"))
}
func (h userHandler) Delete(ctx *fasthttp.RequestCtx) {
	id := ctx.UserValue("id")
	idStr := id.(string)              // แปลงค่าจาก interfacd{} to string
	idInt, err := strconv.Atoi(idStr) // แปลงค่าจาก string to int
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}
	err = h.userServ.Delete(idInt)
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}
	ctx.Response.Header.Set("Content-Type", "application/json")
	ctx.SetStatusCode(200)
	ctx.Write(message("deleted success"))
}
