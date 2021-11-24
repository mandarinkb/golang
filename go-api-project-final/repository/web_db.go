package repository

import (
	"database/sql"
	"errors"
)

type webRepo struct {
	db *sql.DB
}

func NewWebRepo(db *sql.DB) WebRepository {
	return webRepo{db: db}
}

func (w webRepo) Read() ([]Web, error) {
	err := w.db.Ping()
	if err != nil {
		return nil, err
	}

	query := "SELECT * FROM WEB ORDER BY WEB_ID DESC"
	rows, err := w.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	web := []Web{}
	for rows.Next() {
		dataWeb := Web{}
		err = rows.Scan(&dataWeb.WebId, &dataWeb.WebName, &dataWeb.WebUrl, &dataWeb.WebStatus, &dataWeb.IconUrl)
		if err != nil {
			return nil, err
		}
		web = append(web, dataWeb)
	}

	return web, nil
}
func (w webRepo) ReadById(id int) (*Web, error) {
	err := w.db.Ping()
	if err != nil {
		return nil, err
	}

	query := "SELECT * FROM WEB WHERE WEB_ID=?"
	row := w.db.QueryRow(query, id)

	web := Web{}
	err = row.Scan(&web.WebId, &web.WebName, &web.WebUrl, &web.WebStatus, &web.IconUrl)
	if err != nil {
		return nil, err
	}

	return &web, nil
}
func (w webRepo) Create(web Web) error {
	err := w.db.Ping()
	if err != nil {
		return err
	}

	query := "INSERT INTO WEB (WEB_NAME,WEB_URL,WEB_STATUS,ICON_URL) VALUES (?,?,?,?)"
	result, err := w.db.Exec(query, web.WebName, web.WebUrl, web.WebStatus, web.IconUrl)
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
func (w webRepo) Update(web Web) error {
	err := w.db.Ping()
	if err != nil {
		return err
	}

	query := "UPDATE WEB SET WEB_NAME=?, WEB_URL=?, WEB_STATUS=?, ICON_URL=? WHERE WEB_ID=?"
	result, err := w.db.Exec(query, web.WebName, web.WebUrl, web.WebStatus, web.IconUrl, web.WebId)
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
func (w webRepo) UpdateStatus(web Web) error {
	err := w.db.Ping()
	if err != nil {
		return err
	}

	query := "UPDATE WEB SET WEB_STATUS=? WHERE WEB_ID=?"
	result, err := w.db.Exec(query, web.WebStatus, web.WebId)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected <= 0 {
		return errors.New("cannot update status")
	}

	return nil
}
func (w webRepo) Delete(id int) error {
	err := w.db.Ping()
	if err != nil {
		return err
	}

	query := "DELETE FROM WEB WHERE WEB_ID=?"
	result, err := w.db.Exec(query, id)
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
