package service

type Web struct {
	WebId     int    `json:"webId"`
	WebName   string `json:"webName"`
	WebUrl    string `json:"webUrl"`
	WebStatus string `json:"webStatus"`
	IconUrl   string `json:"iconUrl"`
}

type WebService interface {
	Read() ([]Web, error)
	ReadById(id int) (*Web, error)
	Create(web Web) error
	Update(web Web) error
	UpdateStatus(web Web) error
	Delete(id int) error
}
