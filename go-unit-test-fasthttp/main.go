package main

import (
	"fmt"

	"github.com/mandarinkb/test-git/handlers"
	"github.com/mandarinkb/test-git/repository"
	"github.com/mandarinkb/test-git/service"
	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
)

func main() {
	dataMock := repository.NewMock()
	dataServ := service.NewDataService(dataMock)
	dataHandler := handlers.NewDataService(dataServ)
	router := routing.New()
	router.Get("/mock", dataHandler.GetAll)
	router.Get("/mock/<id>", dataHandler.GetById)
	fmt.Println("service run port 8080")
	fasthttp.ListenAndServe(":8080", router.HandleRequest)
}
