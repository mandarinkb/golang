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

func (s switchDBRepo) Read() ([]SwitchDatabase, error) {
	err := s.db.Ping()
	if err != nil {
		return nil, err
	}

	query := "SELECT * FROM SWITCH_DATABASE ORDER BY DATABASE_ID DESC"
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	swDB := []SwitchDatabase{}
	for rows.Next() {
		dataSwDB := SwitchDatabase{}
		err = rows.Scan(&dataSwDB.DatabaseId, &dataSwDB.DatabaseName, &dataSwDB.DatabaseStatus)
		if err != nil {
			return nil, err
		}
		swDB = append(swDB, dataSwDB)
	}

	return swDB, nil
}
func (s switchDBRepo) ReadById(id int) (*SwitchDatabase, error) {
	err := s.db.Ping()
	if err != nil {
		return nil, err
	}

	query := "SELECT * FROM SWITCH_DATABASE WHERE DATABASE_ID =?"
	row := s.db.QueryRow(query, id)

	swDB := SwitchDatabase{}
	err = row.Scan(&swDB.DatabaseId, &swDB.DatabaseName, &swDB.DatabaseStatus)
	if err != nil {
		return nil, err
	}

	return &swDB, nil
}
func (s switchDBRepo) ReadActivateSwitchDatabase() (*SwitchDatabase, error) {
	err := s.db.Ping()
	if err != nil {
		return nil, err
	}

	query := "SELECT * FROM SWITCH_DATABASE WHERE DATABASE_STATUS ='1'"
	row := s.db.QueryRow(query)

	swDB := SwitchDatabase{}
	err = row.Scan(&swDB.DatabaseId, &swDB.DatabaseName, &swDB.DatabaseStatus)
	if err != nil {
		return nil, err
	}

	return &swDB, nil
}
func (s switchDBRepo) Create(swDb SwitchDatabase) error {
	err := s.db.Ping()
	if err != nil {
		return err
	}

	query := "INSERT INTO SWITCH_DATABASE (DATABASE_NAME,DATABASE_STATUS) VALUES (?,?)"
	result, err := s.db.Exec(query, swDb.DatabaseName, swDb.DatabaseStatus)

	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected <= 0 {
		return errors.New("cannot insert")
	}
	return nil
}
func (s switchDBRepo) Update(swDb SwitchDatabase) error {
	err := s.db.Ping()
	if err != nil {
		return err
	}

	query := "UPDATE SWITCH_DATABASE SET DATABASE_NAME=?, DATABASE_STATUS=? WHERE DATABASE_ID=?"
	result, err := s.db.Exec(query, swDb.DatabaseName, swDb.DatabaseStatus, swDb.DatabaseId)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected <= 0 {
		return errors.New("cannot update")
	}

	return nil
}
func (s switchDBRepo) UpdateStatus(swDb SwitchDatabase) error {
	err := s.db.Ping()
	if err != nil {
		return err
	}

	query := "UPDATE SWITCH_DATABASE SET DATABASE_STATUS=? WHERE DATABASE_ID=?"
	result, err := s.db.Exec(query, swDb.DatabaseStatus, swDb.DatabaseId)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected <= 0 {
		return errors.New("cannot update switch")
	}
	// update สถานอื่นๆ ให้เป็น 0
	query2 := "UPDATE SWITCH_DATABASE SET DATABASE_STATUS = '0' WHERE DATABASE_ID <>?"
	result2, err2 := s.db.Exec(query2, swDb.DatabaseId)
	if err2 != nil {
		return err2
	}

	affected2, err2 := result2.RowsAffected()
	if err2 != nil {
		return err2
	}

	if affected2 <= 0 {
		return errors.New("cannot update")
	}

	return nil

}
func (s switchDBRepo) Delete(id int) error {
	err := s.db.Ping()
	if err != nil {
		return err
	}

	query := "DELETE FROM SWITCH_DATABASE WHERE DATABASE_ID=?"
	result, err := s.db.Exec(query, id)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected <= 0 {
		return errors.New("cannot delete")
	}

	return nil
}
