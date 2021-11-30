package repository

import "database/sql"

type webRepo struct {
	db *sql.DB
}

func NewWeb(db *sql.DB) WebRepository {
	return webRepo{db}
}

func (w webRepo) Read() (web []Web, err error) {
	err = w.db.Ping()
	if err != nil {
		return nil, err
	}

	query := "SELECT * FROM WEB ORDER BY WEB_ID DESC"
	rows, err := w.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

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
