package service

type NewResponse struct {
	Content    []ProductRespose `json:"content"`
	Page       int              `json:"page"`
	PageSize   int              `json:"pageSize"`
	TotalPage  int              `json:"totalPage"`
	IsLastPage bool             `json:"isLastPage"`
}

type ProductRespose struct {
	Timestamp     string  `json:"timestamp"`
	WebName       string  `json:"web_name"`
	ProductName   string  `json:"product_name"`
	Category      string  `json:"category"`
	Price         float32 `json:"price"`
	OriginalPrice float32 `json:"original_price"`
	ProductUrl    string  `json:"product_url"`
	Image         string  `json:"image"`
	Icon          string  `json:"icon"`
}
type ProductService interface {
	Search(name string) ([]ProductRespose, error)
	Pagination(page string, limit string) (*NewResponse, error)
}
