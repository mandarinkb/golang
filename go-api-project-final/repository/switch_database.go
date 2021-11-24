package repository

type SwitchDatabase struct {
	DatabaseId     int    `db:"DATABASE_ID"`
	DatabaseName   string `db:"DATABASE_NAME"`
	DatabaseStatus string `db:"DATABASE_STATUS"`
}

type SwitchDatabaseRepository interface {
	Read() ([]SwitchDatabase, error)
	ReadById(id int) (*SwitchDatabase, error)
	Create(swDb SwitchDatabase) error
	Update(swDb SwitchDatabase) error
	UpdateStatus(swDb SwitchDatabase) error
	Delete(id int) error
}
