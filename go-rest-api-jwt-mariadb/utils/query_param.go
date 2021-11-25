package utils

import (
	"math"
	"strconv"
)

type Param struct {
	ProductName string
	Page        string
	Limit       string
}

func DefaultLimit() string {
	return "10"
}

func SetPage(pageStr string, limitStr string) (int, int, error) {
	//strconv.Atoi() convert string to int
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		return 0, 0, err
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		return 0, 0, err
	}
	// ต้อง - 1 เพราะใน database เริ่มจาก 0
	rowCount := (page - 1) * limit
	return rowCount, limit, nil
}

// หาจำนวนหน้า page ทั้งหมด กรณีทำ pagination
func GetTotalPage(totalRows int, limit int) (int, error) {
	var total int
	// แปลงเป็น float64 ก่อน เพื่อใช้ฟังก์ชัน  math.Trunc()
	totalPage := float64(totalRows) / float64(limit)
	// กรณีมีทศนิยม ให้เพิ่มค่าไปอีก 1
	if totalPage != math.Trunc(totalPage) {
		total = (totalRows / limit) + 1
	} else { // กรณีไม่มีทศนิยม
		total = (totalRows / limit)
	}
	return total, nil
}

// เช็คว่าคือ lasty page หรือไม่
func IsLastPage(page int, totalPage int) bool {
	return page == totalPage
}
