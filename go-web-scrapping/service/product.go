package service

import (
	"webscrapping/repository"
)

// กำหนดค่าที่ต้องการเองได้
// from web scrapping
type ProductReponse struct {
	Timestamp     string `json:"timestamp"`
	WebName       string `json:"webName"`
	ProductName   string `json:"productName"`
	Category      string `json:"category"`
	Price         string `json:"price"`         // เปลี่ยนใหม่เป็น string
	OriginalPrice string `json:"originalPrice"` // เปลี่ยนใหม่เป็น string
	ProductUrl    string `json:"productUrl"`
	Image         string `json:"image"`
	Icon          string `json:"icon"`
}

// เอาไว้ insert to database
// from database
type ProductDB struct {
	ProductId     int     `db:"productId"`
	Timestamp     string  `db:"timestamp"`
	WebName       string  `db:"webName"`
	ProductName   string  `db:"productName"`
	Category      string  `db:"category"`
	Price         float32 `db:"price"`         // เปลี่ยนใหม่เป็น string
	OriginalPrice float32 `db:"originalPrice"` // เปลี่ยนใหม่เป็น string
	ProductUrl    string  `db:"productUrl"`
	Image         string  `db:"image"`
	Icon          string  `db:"icon"`
}

type ProductService interface {
	CsvData(repo []repository.Product) ([]ProductReponse, error)
}
