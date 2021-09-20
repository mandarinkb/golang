package repository

type Product struct {
	ProductId     int     `db:"PRODUCT_ID"`
	Timestamp     string  `db:"TIMESTAMP"`
	WebName       string  `db:"WEB_NAME"`
	ProductName   string  `db:"PRODUCT_NAME"`
	Category      string  `db:"CATEGORY"`
	Price         float32 `db:"PRICE"`
	OriginalPrice float32 `db:"ORIGINAL_PRICE"`
	ProductUrl    string  `db:"PRODUCT_URL"`
	Image         string  `db:"IMAGE"`
	Icon          string  `db:"ICON"`
}
type ProductRepository interface {
	Pagination(page int, limit int) ([]Product, int, error)
	Search(name string) ([]Product, int, error)
}
