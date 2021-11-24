package service

type Web struct {
	WebId     int    `db:"webId"`
	WebName   string `db:"webName"`
	WebUrl    string `db:"webUrl"`
	WebStatus string `db:"webStatus"`
	IconUrl   string `db:"iconUrl"`
}

type WebService interface {
	Read() ([]Web, error)
	ReadById(id int) (*Web, error)
	Create(web Web) error
	Update(web Web) error
	UpdateStatus(web Web) error
	Delete(id int) error
}
