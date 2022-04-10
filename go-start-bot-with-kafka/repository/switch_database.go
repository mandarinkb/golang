package repository

type SwitchDatabase struct {
	DatabaseId     int    `db:"DATABASE_ID"`
	DatabaseName   string `db:"DATABASE_NAME"`
	DatabaseStatus string `db:"DATABASE_STATUS"`
}

type SwitchDatabaseRepository interface {
	SwapDatabase() error
	ReadInActivateSwitchDatabase() (*SwitchDatabase, error)
}
