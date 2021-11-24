package service

type SwitchDatabase struct {
	DatabaseId     int    `json:"databaseId"`
	DatabaseName   string `json:"databaseName"`
	DatabaseStatus string `json:"databaseStatus"`
}

type SwitchDatabaseService interface {
	Read() ([]SwitchDatabase, error)
	ReadById(id int) (*SwitchDatabase, error)
	Create(swDb SwitchDatabase) error
	Update(swDb SwitchDatabase) error
	UpdateStatus(swDb SwitchDatabase) error
	Delete(id int) error
}
