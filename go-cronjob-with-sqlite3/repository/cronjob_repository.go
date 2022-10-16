package repository

import "database/sql"

type cronJobRepository struct {
	db *sql.DB
}

func NewCronJobRepository(db *sql.DB) CronJobRepository {
	return &cronJobRepository{db}
}

func (c *cronJobRepository) Read() (cronjob []CronJob, err error) {
	err = c.db.Ping()
	if err != nil {
		return nil, err
	}

	query := "SELECT * FROM cronjob"

	rows, err := c.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		c := CronJob{}
		err = rows.Scan(&c.CronID, &c.CronName, &c.CronExpression, &c.CronUseInFunc, &c.CronStatus)
		if err != nil {
			return nil, err
		}
		cronjob = append(cronjob, c)
	}
	return cronjob, nil
}
