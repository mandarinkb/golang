package service

import "webscrapping/repository"

// กำหนดค่าที่ต้องการเองได้
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

type ProductService interface {
	CsvData(repo []repository.Product) ([]ProductReponse, error)
}
