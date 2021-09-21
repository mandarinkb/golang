package service

import (
	"fmt"
	"strconv"

	"github.com/mandarinkb/go-rest-api-jwt-mariadb/repository"
	"github.com/mandarinkb/go-rest-api-jwt-mariadb/utils"
)

type productService struct {
	productRepo repository.ProductRepository
}

func NewProductServ(productRepo repository.ProductRepository) ProductService {
	return productService{productRepo: productRepo}
}

func (s productService) Search(name string) ([]ProductRespose, error) {
	productRepo, rows, err := s.productRepo.Search(name)
	if err != nil {
		return nil, err
	}
	fmt.Println("rows: ", rows)
	return mapDataProduct(productRepo), nil
}

func (s productService) Pagination(pageStr string, limitStr string) (*NewResponse, error) {

	// ตรวจสอบว่ามีค่า limit หรือไม่
	// กรณีไม่มีให้เซ็ตค่า default ไป
	if limitStr == "" {
		limitStr = utils.NewQueryParam().DefaultLimit()
	}

	pageDb, limit, err := utils.NewQueryParam().SetPage(pageStr, limitStr)
	if err != nil {
		return nil, err
	}

	productRepo, rows, err := s.productRepo.Pagination(pageDb, limit)
	if err != nil {
		return nil, err
	}

	totalPage, err := utils.NewQueryParam().GetTotalPage(rows, limit)
	if err != nil {
		return nil, err
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		return nil, err
	}

	products := NewResponse{
		Content:    mapDataProduct(productRepo),
		Page:       page,
		PageSize:   limit,
		TotalPage:  totalPage,
		IsLastPage: utils.NewQueryParam().IsLastPage(page, totalPage),
	}

	fmt.Println("rows: ", rows)
	return &products, nil
}

// ฟังก็ชัน map data จาก repository ไปยัง service
func mapDataProduct(productRepo []repository.Product) []ProductRespose {
	products := []ProductRespose{}
	for _, row := range productRepo {
		dataRepo := ProductRespose{
			Timestamp:     row.Timestamp,
			WebName:       row.WebName,
			ProductName:   row.ProductName,
			Category:      row.Category,
			Price:         row.Price,
			OriginalPrice: row.OriginalPrice,
			ProductUrl:    row.ProductUrl,
			Image:         row.Image,
			Icon:          row.Icon,
		}
		products = append(products, dataRepo)
	}

	return products
}
