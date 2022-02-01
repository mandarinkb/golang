package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gocolly/colly"
	"github.com/mandarinkb/go-fetch-url-bot-project-final/database"
	"github.com/mandarinkb/go-fetch-url-bot-project-final/repository"
)

type Web struct {
	WebId     int    `json:"webId"`
	WebName   string `json:"webName"`
	WebUrl    string `json:"webUrl"`
	WebStatus string `json:"webStatus"`
	IconUrl   string `json:"iconUrl"`
	Category  string `json:"category"`
}

var baseUrlTescolotus string = "https://shoponline.tescolotus.com"

// หน้าแรกของเว็บใซต์
func MainPage(web Web) {
	c := colly.NewCollector(
		// ต้องใส่ UserAgent ถ้าไม่ใส่อาจจะขึ้น Forbidden
		colly.UserAgent("Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.131 Safari/537.36"))
	// ดึงจาก class
	c.OnHTML(".promotions-by-department--facets.promotions--department.by-department", func(e *colly.HTMLElement) {
		// ดึงจาก li tag
		e.ForEach("li", func(_ int, el *colly.HTMLElement) {

			// ดึงจาก class name
			category := el.ChildText(".name")
			detectNotTakeCategory, err := repository.IsNotTakeCategory(web.WebName, category)
			if err != nil {
				fmt.Println(err)
			}
			// ไม่เอาหมวดหมู่สินค้าที่ได้ตั้งค่าไว้
			if !detectNotTakeCategory {
				// จัดหมวดหมู่ใหม่
				newCategory, err := repository.NewCategory(category)
				if err != nil {
					fmt.Println(err)
				}
				web.Category = newCategory
				fmt.Println(category)
				fmt.Println("new => ", newCategory)
				// ดึง url หน้ารายละเอียด
				web.WebUrl = baseUrlTescolotus + el.ChildAttr("a", "href")
				categoryAllPage(web)
			}
		})
	})
	err := c.Visit(web.WebUrl)
	if err != nil {
		fmt.Println("c.Visit(web.WebUrl) : ", err)
	}

}

// ดึงข้อมูลทุกหมวดหมู่สินค้า และทุกหน้าเว็บไซต์ พร้อมทั้งจัดหมวดหมู่สินค้าใหม่
func categoryAllPage(web Web) {
	redis := database.RedisConn()
	defer redis.Close()

	// ดึงข้อมูลหน้าแรก
	webStr, err := json.Marshal(web)
	if err != nil {
		fmt.Println(err)
	}
	redis.RPush(context.Background(), "fetchUrl", string(webStr))

	// วนหาข้อมูลหน้าถัดไปจนถึงหน้าสุดท้าย
	var sliceUrl []string
	checkLast := true
	for checkLast {
		c := colly.NewCollector(
			colly.UserAgent("Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.131 Safari/537.36"))
		c.OnHTML(".pagination--page-selector-wrapper", func(e *colly.HTMLElement) {
			e.ForEach(".pagination-btn-holder", func(_ int, el *colly.HTMLElement) {
				// เก็บค่าใน pagination เพื่อหาค่าหน้าสุดท้าย
				sliceUrl = append(sliceUrl, el.ChildAttr("a", "href"))
			})
			postfixUrl := sliceUrl[len(sliceUrl)-1]
			// ignoreCategory
			if postfixUrl != "" {
				web.WebUrl = baseUrlTescolotus + postfixUrl
				webStr, err := json.Marshal(web)
				if err != nil {
					fmt.Println(err)
				}
				redis.RPush(context.Background(), "fetchUrl", string(webStr))
				fmt.Println(time.Now().Format(time.RFC3339), " : ", web.Category, " : ", web.WebUrl)
			} else {
				fmt.Println("this last page")
				checkLast = false
			}
		})
		// start scraping (ไว้ล่างสุด)
		err := c.Visit(web.WebUrl)
		if err != nil {
			fmt.Println(err)
		}
	}
}

//
func FindUrlInPage(web Web) {
	redis := database.RedisConn()
	defer redis.Close()

	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.131 Safari/537.36"))
	// ดึงจาก class has-trolley.main
	c.OnHTML(".has-trolley.main", func(e *colly.HTMLElement) {
		e.ForEach(".tile-content", func(_ int, el *colly.HTMLElement) {
			web.WebUrl = baseUrlTescolotus + el.ChildAttr("a", "href")
			webStr, err := json.Marshal(web)
			if err != nil {
				fmt.Println(err)
			}
			redis.RPush(context.Background(), "detailUrl", string(webStr))
			fmt.Println(time.Now().Format(time.RFC3339), " : ", "lotus : ", web.WebUrl)
		})
	})
	// start scraping (ไว้ล่างสุด)
	err := c.Visit(web.WebUrl)
	if err != nil {
		fmt.Println(err)
	}
}
