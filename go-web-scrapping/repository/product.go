package repository

type Product struct {
	Timestamp     string  `json:"timestamp"`
	WebName       string  `json:"webName"`
	ProductName   string  `json:"productName"`
	Category      string  `json:"category"`
	Price         float32 `json:"price"`
	OriginalPrice float32 `json:"originalPrice"`
	ProductUrl    string  `json:"productUrl"`
	Image         string  `json:"image"`
	Icon          string  `json:"icon"`
}

type ProductRepositoy interface {
	Makro() ([]Product, error)
	Bigc() ([]Product, error)
	Lotus() ([]Product, error)
}
