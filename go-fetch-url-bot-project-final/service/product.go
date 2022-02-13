package service

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
type Web struct {
	WebId     int    `json:"webId"`
	WebName   string `json:"webName"`
	WebUrl    string `json:"webUrl"`
	WebStatus string `json:"webStatus"`
	IconUrl   string `json:"iconUrl"`
	Category  string `json:"category"`
	MenuId    string `json:"menuId"`    // for makro
	MakroPage int    `json:"makroPage"` // for makro
	EntityId  string `json:"entityId"`  // for bigc
	BigcPage  int    `json:"bigcPage"`  // for bigc
}
