package service

import (
	"encoding/json"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/mandarinkb/go-web-scraping-bot-with-kafka/repository"
	"github.com/mandarinkb/go-web-scraping-bot-with-kafka/utils"
)

type productService struct {
	swDbRepo repository.SwitchDatabaseRepository
}

func NewProductService(swDbRepo repository.SwitchDatabaseRepository) ProductService {
	return productService{swDbRepo}
}
func (s productService) Tescolotus(web Web) error {
	logger, err := utils.LogConf()
	if err != nil {
		logger.Error(err.Error(), utils.Url("-"),
			utils.User("-"), utils.Type(utils.TypeBot))
	}
	// close log
	defer logger.Sync()

	var products Product
	var img string
	var priceAll string
	var name string
	var price float64
	var originalPrice float64
	var discount float64
	var slicePrice []string
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.131 Safari/537.36"))
	// ดึงจาก class has-trolley.main
	c.OnHTML(".product-details-page", func(e *colly.HTMLElement) {
		name = e.ChildText(".product-details-tile__title")
		promotion := e.ChildText(".offer-text")
		if promotion != "" {
			e.ForEach(".offer-text", func(_ int, el *colly.HTMLElement) {
				slicePrice = append(slicePrice, el.Text)
			})
			priceAll = slicePrice[0] //offer-text .first
			parts := strings.Split(priceAll, "บาท")
			p := strings.ReplaceAll(parts[0], "ราคาพิเศษ ", "")
			p = strings.ReplaceAll(p, ".00 ", "")
			o := strings.ReplaceAll(parts[1], " จากราคาปกติ  ", "")
			o = strings.ReplaceAll(o, ".00 ", "")
			price, originalPrice = utils.StrToFloat64(p), utils.StrToFloat64(o)
			discount = (((originalPrice - price) / originalPrice) * 100)
			// ปัดเศษ
			discount = math.Round(discount)
		}

		e.ForEach(".product-image__container", func(_ int, el *colly.HTMLElement) {
			img = el.ChildAttr("img", "src")
		})
		products = Product{
			Timestamp:     time.Now().Format(time.RFC3339),
			WebName:       web.WebName,
			ProductName:   name,
			Category:      web.Category,
			Price:         price,
			OriginalPrice: originalPrice,
			Discount:      discount,
			ProductUrl:    web.WebUrl,
			Image:         img,
			Icon:          web.IconUrl,
		}
		jsonProducts, err := json.Marshal(products)
		if err != nil {
			logger.Error(err.Error(), utils.Url("-"),
				utils.User("-"), utils.Type(utils.TypeBot))
		}
		dbRepo, err := s.swDbRepo.GetInActivateDatabaseName()
		if err != nil {
			logger.Error(err.Error(), utils.Url("-"),
				utils.User("-"), utils.Type(utils.TypeBot))
		}
		dbname := dbRepo.DatabaseName

		err = insertToElasticsearch(dbname, string(jsonProducts))
		if err != nil {
			logger.Error(err.Error(), utils.Url("-"),
				utils.User("-"), utils.Type(utils.TypeBot))
		}
		fmt.Println(time.Now().Format(time.RFC3339), " : ", "lotus : ", name)
	})

	// start scraping (ไว้ล่างสุด)
	err = c.Visit(web.WebUrl)
	if err != nil {
		logger.Error(err.Error(), utils.Url("-"),
			utils.User("-"), utils.Type(utils.TypeBot))
	}
	return nil
}
