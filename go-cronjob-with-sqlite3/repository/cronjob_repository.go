package repository

import (
	"database/sql"
	"errors"

	"github.com/robfig/cron/v3"
)

type cronJobRepository struct {
	db *sql.DB
}

func NewCronJobRepository(db *sql.DB) CronJobRepository {
	return &cronJobRepository{db}
}

func (c *cronJobRepository) GetCronJob() (cronjob []CronJob, err error) {
	err = c.db.Ping()
	if err != nil {
		return nil, err
	}

	query := "SELECT * FROM CRONJOB"

	rows, err := c.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		c := CronJob{}
		err = rows.Scan(&c.CronID, &c.CronName, &c.CronExpression, &c.CronFunctionRef, &c.CronStatus)
		if err != nil {
			return nil, err
		}
		cronjob = append(cronjob, c)
	}
	return cronjob, nil
}

func (s *cronJobRepository) GetCronJobByID(id int) (*CronJob, error) {
	err := s.db.Ping()
	if err != nil {
		return nil, err
	}

	query := "SELECT * FROM CRONJOB WHERE CRONID =?"
	row := s.db.QueryRow(query, id)
	cronjob := CronJob{}

	err = row.Scan(&cronjob.CronID, &cronjob.CronName, &cronjob.CronExpression, &cronjob.CronFunctionRef, &cronjob.CronStatus)
	if err != nil {
		return nil, err
	}

	return &cronjob, nil
}

func (s *cronJobRepository) CreateCronJob(cronjob CronJob) error {
	err := s.db.Ping()
	if err != nil {
		return err
	}

	query := "INSERT INTO CRONJOB(CRONNAME,CRONEXPRESSION,CRONFUNCTIONREF,CRONSTATUS) VALUES (?,?,?,?)"
	result, err := s.db.Exec(query, cronjob.CronName, cronjob.CronExpression, cronjob.CronFunctionRef, cronjob.CronStatus)
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
func (s *cronJobRepository) UpdateCronJob(cronjob CronJob) error {
	err := s.db.Ping()
	if err != nil {
		return err
	}

	query := "UPDATE CRONJOB SET CRONNAME=?, CRONEXPRESSION=?, CRONFUNCTIONREF=?, CRONSTATUS=?  WHERE CRONID=?"
	result, err := s.db.Exec(query, cronjob.CronName, cronjob.CronExpression, cronjob.CronFunctionRef, cronjob.CronStatus, cronjob.CronID)
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
func (s *cronJobRepository) DeleteCronJob(id int) error {
	err := s.db.Ping()
	if err != nil {
		return err
	}

	query := "DELETE FROM CRONJOB WHERE CRONID=?"
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

func (c *cronJobRepository) RunJob(cronExpression string, cmd func()) (cID cron.EntryID, err error) {
	return Cron.AddFunc(cronExpression, cmd)
}

func (c *cronJobRepository) RemoveJob(id cron.EntryID) {
	Cron.Remove(id)
}
