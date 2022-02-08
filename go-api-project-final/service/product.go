package service

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/mandarinkb/go-api-project-final/database"
	"github.com/mandarinkb/go-api-project-final/repository"
)

type Product struct {
	Image         string  `json:"image"`
	OriginalPrice float64 `json:"originalPrice"`
	Price         float64 `json:"price"`
	Name          string  `json:"name"`
	Icon          string  `json:"icon"`
	Discount      float64 `json:"discount"`
	Category      string  `json:"category"`
	ProductUrl    string  `json:"productUrl"`
}

type ProductRequest struct {
	UserId   string   `json:"userId"`
	Name     string   `json:"name"`
	History  []string `json:"history"`
	Category string   `json:"category"`
	WebName  []string `json:"webName"`
	MinPrice int      `json:"minPrice"`
	MaxPrice int      `json:"maxPrice"`
}
type productService struct {
	productRepo repository.ProductRepository
}

type ProductService interface {
	GetSearch(from string, name string) ([]Product, error)
	GetFilterSearch(from string, name string, web []string, min int, max int) ([]Product, error)
	GetCategory(from string, category string, web []string) ([]Product, error)
	GetHistory(from string, history []string) ([]Product, error)
	GetWebName() ([]Web, error)
}

func NewProductServ(productRepo repository.ProductRepository) ProductService {
	return productService{productRepo}
}

func (p productService) GetSearch(from string, name string) (product []Product, err error) {
	db, err := database.Conn()
	if err != nil {
		return nil, err
	}
	defer db.Close()
	// หา index โดยค้นหาจาก db
	index, err := repository.NewSwitchDBRepo(db).ReadActivateSwitchDatabase()
	if err != nil {
		return nil, err
	}
	fmt.Println("index : ", index.DatabaseName)

	productRepo, err := p.productRepo.GetSearch(index.DatabaseName, from, jsonName(name))
	if err != nil {
		return nil, err
	}
	for _, row := range productRepo {
		product = append(product, mapDataProduct(row))
	}

	return product, nil
}
func (p productService) GetFilterSearch(from string, name string, web []string,
	min int, max int) (product []Product, err error) {
	db, err := database.Conn()
	if err != nil {
		return nil, err
	}
	defer db.Close()
	// หา index โดยค้นหาจาก db
	index, err := repository.NewSwitchDBRepo(db).ReadActivateSwitchDatabase()
	if err != nil {
		return nil, err
	}
	fmt.Println("index : ", index.DatabaseName)
	// แปลงเห็น string
	minStr := strconv.Itoa(min)
	maxStr := strconv.Itoa(max)
	productRepo, err := p.productRepo.GetFilterSearch(index.DatabaseName, from, jsonName(name), jsonWeb(web), minStr, maxStr)
	if err != nil {
		return nil, err
	}
	for _, row := range productRepo {
		product = append(product, mapDataProduct(row))
	}
	return product, nil
}
func (p productService) GetCategory(from string, category string, web []string) (product []Product, err error) {
	db, err := database.Conn()
	if err != nil {
		return nil, err
	}
	defer db.Close()
	// หา index โดยค้นหาจาก db
	index, err := repository.NewSwitchDBRepo(db).ReadActivateSwitchDatabase()
	if err != nil {
		return nil, err
	}
	fmt.Println("index : ", index.DatabaseName)

	productRepo, err := p.productRepo.GetCategory(index.DatabaseName, from, category, jsonWeb(web))
	if err != nil {
		return nil, err
	}
	for _, row := range productRepo {
		product = append(product, mapDataProduct(row))
	}
	return product, nil
}
func (p productService) GetHistory(from string, history []string) (product []Product, err error) {
	db, err := database.Conn()
	if err != nil {
		return nil, err
	}
	defer db.Close()
	// หา index โดยค้นหาจาก db
	index, err := repository.NewSwitchDBRepo(db).ReadActivateSwitchDatabase()
	if err != nil {
		return nil, err
	}
	fmt.Println("index : ", index.DatabaseName)

	if len(history) == 0 {
		fmt.Println("history")
		productRepo, err := p.productRepo.GetHistory(index.DatabaseName, from)
		if err != nil {
			return nil, err
		}
		for _, row := range productRepo {
			product = append(product, mapDataProduct(row))
		}
		return product, nil
	} else {
		fmt.Println("history with data")
		productRepo, err := p.productRepo.GetHistoryWithHistoryData(index.DatabaseName, from, jsonNameSlice(history))
		if err != nil {
			return nil, err
		}
		for _, row := range productRepo {
			product = append(product, mapDataProduct(row))
		}
		return product, nil
	}
}
func (p productService) GetWebName() (web []Web, err error) {
	db, err := database.Conn()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	webRepo, err := repository.NewWebRepo(db).ReadActivateWeb()
	if err != nil {
		return nil, err
	}
	for _, row := range webRepo {
		web = append(web, mapDataWebResponse(row))
	}
	return web, nil
}

// แปลงค่า name ที่รับเข้ามาเพื่อจัดข้อมูล json เพื่อ query ใน elasticsearch
func jsonName(name string) string {
	nameSlice := strings.Split(name, " ")
	var jsonName string
	// กรณีมีคำเดียว
	if len(nameSlice) == 1 {
		for _, row := range nameSlice {
			jsonName = `{"regexp": {"productName": {"value": "(.*)` + row + `(.*)"}}}`
		}
	}
	// กรณีมีหลายคำ
	if len(nameSlice) > 1 {
		for _, row := range nameSlice {
			// ดักไว้ ป้องกัน space
			if row != "" {
				jsonName = jsonName + `{"regexp": {"productName": {"value": "(.*)` + row + `(.*)"}}},`
			}
		}
		// ตัด , ตัวสุดท้ายออก
		jsonName = strings.TrimSuffix(jsonName, ",")
	}
	return jsonName
}

// แปลงค่า slice name ที่รับเข้ามาเพื่อจัดข้อมูล json เพื่อ query ใน elasticsearch
func jsonNameSlice(name []string) string {
	var jsonName string
	// กรณีมีคำเดียว
	if len(name) == 1 {
		for _, row := range name {
			jsonName = `{"regexp": {"productName": {"value": "(.*)` + row + `(.*)"}}}`
		}
	}
	// กรณีมีหลายคำ
	if len(name) > 1 {
		for _, row := range name {
			// ดักไว้ ป้องกัน space
			if row != "" {
				jsonName = jsonName + `{"regexp": {"productName": {"value": "(.*)` + row + `(.*)"}}},`
			}
		}
		// ตัด , ตัวสุดท้ายออก
		jsonName = strings.TrimSuffix(jsonName, ",")
	}
	return jsonName
}

func jsonWeb(web []string) string {
	var jsonWeb string
	// กรณีมีค่าเดียว
	if len(web) == 1 {
		for _, row := range web {
			jsonWeb = `{"match_phrase": {"webName": "` + row + `"}}`
		}
	}
	// กรณีมีหลายค่า
	if len(web) > 1 {
		for _, row := range web {
			// ดักไว้ ป้องกัน space
			if row != "" {
				jsonWeb = jsonWeb + `{"match_phrase": {"webName": "` + row + `"}},`
			}
		}
		// ตัด , ตัวสุดท้ายออก
		jsonWeb = strings.TrimSuffix(jsonWeb, ",")
	}
	return jsonWeb
}

func mapDataProduct(product repository.Product) Product {
	return Product{
		Image:         product.Image,
		OriginalPrice: product.OriginalPrice,
		Price:         product.Price,
		Name:          product.ProductName,
		Icon:          product.Icon,
		Discount:      product.Discount,
		Category:      product.Category,
		ProductUrl:    product.ProductUrl,
	}
}
