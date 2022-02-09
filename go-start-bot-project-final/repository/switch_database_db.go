package repository

import (
	"database/sql"
	"errors"
)

type switchDBRepo struct {
	db *sql.DB
}

func NewSwitchDBRepo(db *sql.DB) SwitchDatabaseRepository {
	return switchDBRepo{db: db}
}
func (s switchDBRepo) SwapDatabase() error {
	err := s.db.Ping()
	if err != nil {
		return err
	}

	query := `update SWITCH_DATABASE
				set DATABASE_STATUS =
				case when DATABASE_STATUS = '1' then '0'
		 			when DATABASE_STATUS = '0' then '1'
				end
				where DATABASE_NAME in ('web-scrapping-db-1','web-scrapping-db-2')`
	result, err := s.db.Exec(query)

	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected <= 0 {
		return errors.New("cannot swap database")
	}
	return nil
}
func (s switchDBRepo) ReadInActivateSwitchDatabase() (*SwitchDatabase, error) {
	err := s.db.Ping()
	if err != nil {
		return nil, err
	}

	query := "SELECT * FROM SWITCH_DATABASE WHERE DATABASE_STATUS ='0'"
	row := s.db.QueryRow(query)

	swDB := SwitchDatabase{}
	err = row.Scan(&swDB.DatabaseId, &swDB.DatabaseName, &swDB.DatabaseStatus)
	if err != nil {
		return nil, err
	}

	return &swDB, nil
}
