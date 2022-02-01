package repository

type switchDatabaseMock struct {
	DatabaseId     int
	DatabaseName   string
	DatabaseStatus string
}

func NewSwitchDatabaseMock() SwitchDatabaseRepository {
	return switchDatabaseMock{
		DatabaseId:     1,
		DatabaseName:   "web-scrapping-db-1",
		DatabaseStatus: "1"}
}

func (s switchDatabaseMock) GetDatabaseName() (*SwitchDatabase, error) {
	swData := SwitchDatabase{
		DatabaseId:     s.DatabaseId,
		DatabaseName:   s.DatabaseName,
		DatabaseStatus: s.DatabaseStatus,
	}
	return &swData, nil
}
