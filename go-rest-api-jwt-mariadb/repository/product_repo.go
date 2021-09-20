package repository

import (
	"database/sql"
)

type productRepository struct {
	db *sql.DB
}

func NewProductRepo(db *sql.DB) ProductRepository {
	return productRepository{db: db}
}

// ฟังก์ชันทำ pagination ผ่าน api
// จะ return ค่า ([]Product, totalRows, error) ไปให้
func (p productRepository) Pagination(rowCount int, limit int) ([]Product, int, error) {
	err := p.db.Ping()
	if err != nil {
		return nil, 0, err
	}

	qurey := "SELECT PRODUCT_ID, TIMESTAMP, WEB_NAME, PRODUCT_NAME, CATEGORY, PRICE, ORIGINAL_PRICE, PRODUCT_URL, IMAGE, ICON FROM PRODUCTS LIMIT ?,?"
	rows, err := p.db.Query(qurey, rowCount, limit)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	// หาจำนวน rows ทั้งหมดเพื่อทำ pagination
	q := "SELECT COUNT(*) FROM PRODUCTS"
	totalRows, err := totalRows(p.db, q, "")
	if err != nil {
		return nil, 0, err
	}

	products := []Product{}
	for rows.Next() {
		product := Product{}
		err := rows.Scan(&product.ProductId,
			&product.Timestamp,
			&product.WebName,
			&product.ProductName,
			&product.Category,
			&product.Price,
			&product.OriginalPrice,
			&product.ProductUrl,
			&product.Image,
			&product.Icon)
		if err != nil {
			return nil, 0, err
		}
		products = append(products, product)
	}
	return products, totalRows, nil
}

// ฟังก์ชันค้นหาข้อความ โดยจะค้นหาข้อความที่อยู่ใน filed database มาให้
// จะ return ค่า ([]Product, totalRows, error) ไปให้
func (p productRepository) Search(name string) ([]Product, int, error) {
	err := p.db.Ping()
	if err != nil {
		return nil, 0, err
	}

	name = "%" + name + "%"
	// ต้องแปลงค่าแบบนี้ถึงจะใส่ใน LIKE ? ได้
	// '%name%'
	qurey := "SELECT PRODUCT_ID, TIMESTAMP, WEB_NAME, PRODUCT_NAME, CATEGORY, PRICE, ORIGINAL_PRICE, PRODUCT_URL, IMAGE, ICON FROM PRODUCTS WHERE PRODUCT_NAME LIKE ?"
	rows, err := p.db.Query(qurey, name)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	// หาจำนวน rows ทั้งหมดเพื่อทำ pagination
	q := "SELECT COUNT(*) FROM PRODUCTS WHERE PRODUCT_NAME LIKE ?"
	totalRows, err := totalRows(p.db, q, name)
	if err != nil {
		return nil, 0, err
	}

	products := []Product{}
	for rows.Next() {
		product := Product{}
		err := rows.Scan(&product.ProductId,
			&product.Timestamp,
			&product.WebName,
			&product.ProductName,
			&product.Category,
			&product.Price,
			&product.OriginalPrice,
			&product.ProductUrl,
			&product.Image,
			&product.Icon)
		if err != nil {
			return nil, 0, err
		}
		products = append(products, product)
	}
	return products, totalRows, nil
}

// ฟังก็ชั่นหาจำนวน rows ที่ query ทั้งหมด เพื่องทำ pagination
func totalRows(db *sql.DB, query string, param string) (int, error) {
	var rows *sql.Rows
	var err error
	if param == "" {
		rows, err = db.Query(query)
		if err != nil {
			return 0, err
		}
	} else {
		rows, err = db.Query(query, param)
		if err != nil {
			return 0, err
		}
	}
	defer rows.Close()

	var count int
	for rows.Next() {
		err := rows.Scan(&count)
		if err != nil {
			return 0, err
		}
	}
	return count, nil
}
