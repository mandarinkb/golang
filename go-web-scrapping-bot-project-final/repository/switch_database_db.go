package repository

import (
	"database/sql"
)

type switchDatabaseDB struct {
	db *sql.DB
}

func NewSwitchDatabaseDB(db *sql.DB) SwitchDatabaseRepository {
	return switchDatabaseDB{db}
}

func (s switchDatabaseDB) GetDatabaseName() (*SwitchDatabase, error) {
	err := s.db.Ping()
	if err != nil {
		return nil, err
	}
	query := "SELECT DATABASE_NAME FROM SWITCH_DATABASE WHERE DATABASE_STATUS='1'"
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	swData := SwitchDatabase{}
	for rows.Next() {
		err = rows.Scan(&swData.DatabaseName)
		if err != nil {
			return nil, err
		}
	}
	return &swData, nil
}
