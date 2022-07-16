package handlers_test

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"path"
	"runtime"
	"testing"

	"github.com/mandarinkb/test-git/handlers"
	"github.com/mandarinkb/test-git/helper"
	"github.com/mandarinkb/test-git/repository"
	"github.com/mandarinkb/test-git/service"
	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttputil"
)

// [fix bug test] no such file or directory
func init() {
	_, filename, _, _ := runtime.Caller(0)
	//The ".." may change depending on you folder structure
	dir := path.Join(path.Dir(filename), "..")
	err := os.Chdir(dir)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Current test filename: %s", dir)

}
func serve(handler fasthttp.RequestHandler, req *http.Request) (*http.Response, error) {
	ln := fasthttputil.NewInmemoryListener()
	defer ln.Close()

	go func() {
		err := fasthttp.Serve(ln, handler)
		if err != nil {
			panic(fmt.Errorf("failed to serve: %v", err))
		}
	}()

	client := http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				return ln.Dial()
			},
		},
	}

	return client.Do(req)
}
func TestGetAll(t *testing.T) {
	dataRepo := repository.NewMock()
	dataServ := service.NewDataService(dataRepo)
	dataHdlr := handlers.NewDataService(dataServ)
	router := routing.New()
	router.Get("/mock", dataHdlr.GetAll)
	// send request
	req, err := http.NewRequest("GET", "http://localhost:8080/mock", nil)
	if err != nil {
		t.Error(err)
	}
	// get response
	res, err := serve(router.HandleRequest, req)
	if err != nil {
		t.Error(err)
	}
	helper.Equals(t, http.StatusOK, res.StatusCode)
}

func TestGetById(t *testing.T) {
	dataRepo := repository.NewMock()
	dataServ := service.NewDataService(dataRepo)
	dataHdlr := handlers.NewDataService(dataServ)
	router := routing.New()
	router.Get("/mock/<id>", dataHdlr.GetById)
	// send request
	req, err := http.NewRequest("GET", "http://localhost:8080/mock/1", nil)
	if err != nil {
		t.Error(err)
	}
	// get response
	res, err := serve(router.HandleRequest, req)
	if err != nil {
		t.Error(err)
	}
	helper.Equals(t, http.StatusOK, res.StatusCode)
}
