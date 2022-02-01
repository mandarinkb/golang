package service

type Web struct {
	WebId     int    `json:"webId"`
	WebName   string `json:"webName"`
	WebUrl    string `json:"webUrl"`
	WebStatus string `json:"webStatus"`
	IconUrl   string `json:"iconUrl"`
	Category  string `json:"category"`
}
type Product struct {
	Timestamp     string  `json:"timestamp"`
	WebName       string  `json:"webName"`
	ProductName   string  `json:"productName"`
	Category      string  `json:"category"`
	Price         float64 `json:"price"`
	OriginalPrice float64 `json:"originalPrice"`
	Discount      float64 `json:"discount"`
	ProductUrl    string  `json:"productUrl"`
	Image         string  `json:"image"`
	Icon          string  `json:"icon"`
}
type SwitchDatabaseService interface {
	ProdudtDetail(web Web) error
}
