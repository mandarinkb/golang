package repository

type Web struct {
	WebId     int    `db:"WEB_ID"`
	WebName   string `db:"WEB_NAME"`
	WebUrl    string `db:"WEB_URL"`
	WebStatus string `db:"WEB_STATUS"`
	IconUrl   string `db:"ICON_URL"`
}

type WebRepository interface {
	Read() ([]Web, error)
	ReadById(id int) (*Web, error)
	ReadActivateWeb() ([]Web, error)
	Create(web Web) error
	Update(web Web) error
	UpdateStatus(web Web) error
	Delete(id int) error
}
