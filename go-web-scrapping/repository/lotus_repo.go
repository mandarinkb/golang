package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"
	"webscrapping/util"

	"github.com/gocolly/colly"
)

var lotusDomain = "https://shoponline.tescolotus.com"

type category struct {
	Category string `json:"category"`
	Url      string `json:"url"`
}

// ต้องมี NewLotus
// เพราะเป็นตัวเชื่อมใช้ implement interface
// โดยจะเชื่อมเฉพาะ ที่เป็น receiver function เท่านั้น
func NewLotus() ProductRepositoy {
	return Product{}
}
func (Product) Lotus() ([]Product, error) {
	var products []Product
	url := "https://shoponline.tescolotus.com/groceries/th-TH/"
	cateMain, _ := mainLotus(url)
	for _, row := range cateMain {
		categories, _ := categoryLink(row.Url, row.Category)
		for _, newCate := range categories {
			allPages, _ := categoryAllPage(newCate.Url, newCate.Category)
			for _, page := range allPages {
				allLink, _ := linkAllProductInPage(page.Url, page.Category)
				for _, link := range allLink {
					pd, _ := productDetail(link.Url, link.Category)
					products = append(products, pd)
				}
			}
		}
	}
	return products, nil
}

// step 1
func mainLotus(url string) ([]category, error) {
	var categories []category
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.131 Safari/537.36"))
	// ดึงจาก class menu-superdepartment
	c.OnHTML(".menu-superdepartment", func(e *colly.HTMLElement) {
		// ดึงจาก li tag
		e.ForEach("li", func(_ int, el *colly.HTMLElement) {
			categoryStr := strings.ReplaceAll(el.ChildText("span"), "แผนกซื้อแผนก", "")
			newCategory := strings.ReplaceAll(categoryStr, "ซื้อ", "")
			urlCategory := lotusDomain + el.ChildAttr("a", "href")
			// สร้าง struct
			category := category{
				Category: newCategory,
				Url:      urlCategory,
			}
			categories = append(categories, category)
			fmt.Println(time.Now().Format(time.RFC3339), " : ", "lotus : ", urlCategory)
		})
	})
	// start scraping (ไว้ล่างสุด)
	err := c.Visit(url)
	if err != nil {
		fmt.Println(err)
		//return nil, err
	}
	return categories, nil
}

// step 2 find link category
func categoryLink(url string, cate string) ([]category, error) {
	var categories []category
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.131 Safari/537.36"))
	// ดึงจาก class list-item.list-subheader
	c.OnHTML(".list-item.list-subheader", func(e *colly.HTMLElement) {
		urlCategory := lotusDomain + e.ChildAttr("a", "href")
		// สร้าง struct
		category := category{
			Category: cate,
			Url:      urlCategory,
		}
		categories = append(categories, category)
		fmt.Println(time.Now().Format(time.RFC3339), " : ", "lotus : ", urlCategory)

	})
	// start scraping (ไว้ล่างสุด)
	err := c.Visit(url)
	if err != nil {
		fmt.Println(err)
		//return nil, err
	}
	return categories, nil
}

// step 3 find category all page
func categoryAllPage(url string, cate string) ([]category, error) {
	var sliceUrl []string
	var categories []category
	checkLast := true
	for checkLast {
		c := colly.NewCollector(
			colly.UserAgent("Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.131 Safari/537.36"))
		c.OnHTML(".pagination--page-selector-wrapper", func(e *colly.HTMLElement) {
			e.ForEach(".pagination-btn-holder", func(_ int, el *colly.HTMLElement) {
				el.ChildAttr("a", "href")
				// เก็บค่าใน pagination เพื่อหาค่าหน้าสุดท้าย
				sliceUrl = append(sliceUrl, el.ChildAttr("a", "href"))
			})
			lasturl := sliceUrl[len(sliceUrl)-1]
			if lasturl != "" {
				url = lotusDomain + lasturl
				category := category{
					Category: cate,
					Url:      url,
				}
				categories = append(categories, category)
				fmt.Println(time.Now().Format(time.RFC3339), " : ", "lotus : ", url)
			} else {
				fmt.Println("this last page")
				checkLast = false
			}
		})
		// start scraping (ไว้ล่างสุด)
		err := c.Visit(url)
		if err != nil {
			fmt.Println(err)
			//return nil, err
		}
	}
	return categories, nil
}

// step 4  get all product in page
func linkAllProductInPage(url string, cate string) ([]category, error) {
	var categories []category
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.131 Safari/537.36"))
	// ดึงจาก class has-trolley.main
	c.OnHTML(".has-trolley.main", func(e *colly.HTMLElement) {
		e.ForEach(".tile-content", func(_ int, el *colly.HTMLElement) {
			urlDetail := lotusDomain + el.ChildAttr("a", "href")
			category := category{
				Category: cate,
				Url:      urlDetail,
			}
			categories = append(categories, category)
			fmt.Println(time.Now().Format(time.RFC3339), " : ", "lotus : ", urlDetail)
		})
	})
	// start scraping (ไว้ล่างสุด)
	err := c.Visit(url)
	if err != nil {
		fmt.Println(err)
		//return nil, err
	}
	return categories, nil
}

// step 5 web scrapping in page
func productDetail(url string, cate string) (Product, error) {
	var products Product
	var img string
	var priceAll string
	var name string
	var price float32
	var orignalPrice float32
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
			price, orignalPrice = util.StrToFloat32(p), util.StrToFloat32(o)
		} else {
			// fmt.Println("==== not promotion ====")
			pStr := e.ChildText(".price-per-sellable-unit.price-per-sellable-unit--price.price-per-sellable-unit--price-per-item")
			pStr = strings.ReplaceAll(pStr, "฿ ", "")
			pStr = strings.ReplaceAll(pStr, ".00 ", "")
			price = util.StrToFloat32(pStr)
			orignalPrice = price
		}

		e.ForEach(".product-image__container", func(_ int, el *colly.HTMLElement) {
			img = el.ChildAttr("img", "src")
		})
		products = Product{
			Timestamp:     time.Now().Format(time.RFC3339),
			WebName:       "lotus",
			ProductName:   name,
			Category:      cate,
			Price:         price,
			OriginalPrice: orignalPrice,
			ProductUrl:    url,
			Image:         img,
			Icon:          "https://www.tescolotus.com/assets/theme2018/tl-theme/img/logo.png",
		}
		fmt.Println(time.Now().Format(time.RFC3339), " : ", "lotus : ", name)

		InsertLotusToDB(products)

	})
	// start scraping (ไว้ล่างสุด)
	err := c.Visit(url)
	if err != nil {
		fmt.Println(err)
		//return Product{}, err
	}

	return products, nil
}

func InsertLotusToDB(repo Product) error {
	db, err := sql.Open("mysql", "root:mandarinkb@tcp(127.0.0.1)/TEST_DB?charset=utf8")
	if err != nil {
		fmt.Print(err)
	}
	err = db.Ping()
	if err != nil {
		return err
	}

	query := "INSERT INTO PRODUCT (TIMESTAMP,WEB_NAME,PRODUCT_NAME,CATEGORY,PRICE,ORIGINAL_PRICE,PRODUCT_URL,IMAGE,ICON) VALUES (?,?,?,?,?,?,?,?,?)"
	result, err := db.Exec(query, repo.Timestamp,
		repo.WebName,
		repo.ProductName,
		repo.Category,
		repo.Price,
		repo.OriginalPrice,
		repo.ProductUrl,
		repo.Image,
		repo.Icon)
	if err != nil {
		return err
	}

	// รับค่ามาเพื่อตรวสสอบว่า insert สำเร็จหรือไม่
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	// กรณี insert ไม่สำเร็จ
	if affected <= 0 {
		return errors.New("cannot insert")
	}

	db.Close()
	return nil
}
