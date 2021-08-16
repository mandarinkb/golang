package service

import (
	"webscrapping/repository"
	"webscrapping/util"
)

// ต้องมี NewProductService
// เพราะเป็นตัวเชื่อมใช้ implement interface
// โดยจะเชื่อมเฉพาะ ที่เป็น receiver function เท่านั้น
func NewProductService() ProductService {
	return ProductReponse{}
}
func (ProductReponse) CsvData(repo []repository.Product) ([]ProductReponse, error) {
	var newProducts []ProductReponse
	for i := range repo {
		product := ProductReponse{
			Timestamp:     repo[i].Timestamp,
			WebName:       repo[i].WebName,
			ProductName:   repo[i].ProductName,
			Category:      repo[i].Category,
			Price:         util.FloatToString(repo[i].Price),
			OriginalPrice: util.FloatToString(repo[i].OriginalPrice),
			ProductUrl:    repo[i].ProductUrl,
			Image:         repo[i].Image,
			Icon:          repo[i].Icon,
		}
		newProducts = append(newProducts, product)
	}
	return newProducts, nil
}
