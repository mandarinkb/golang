package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mandarinkb/go-rest-api-jwt-mariadb/service"
	"github.com/mandarinkb/go-rest-api-jwt-mariadb/utils"
)

type productHandler struct {
	prodSrv service.ProductService
}

func NewProductHandler(prodSrv service.ProductService) productHandler {
	return productHandler{prodSrv: prodSrv}
}

func (h productHandler) SearchProduct(c *gin.Context) {
	param := utils.Param{}
	param.ProductName = c.Query("product_name")
	products, err := h.prodSrv.Search(param.ProductName)
	if err != nil {
		fmt.Println(err)
		// gin.H{} สามารถใส่ข้อความแล้วจะ return json ให้
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	// return ค่า http status และค่า json obj หรือ json array ออกให้
	c.IndentedJSON(http.StatusOK, products)
}

func (h productHandler) PaginationProduct(c *gin.Context) {
	param := utils.Param{}
	param.Page = c.Query("page")
	param.Limit = c.Query("limit")
	products, err := h.prodSrv.Pagination(param.Page, param.Limit)
	if err != nil {
		fmt.Println(err)
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, products)
}
