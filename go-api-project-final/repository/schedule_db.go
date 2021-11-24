package repository

import (
	"database/sql"
	"errors"
)

type scheduleRepo struct {
	db *sql.DB
}

func NewScheduleRepo(db *sql.DB) ScheduleRepository {
	return scheduleRepo{db: db}
}

func (s scheduleRepo) Read() ([]Schedule, error) {
	err := s.db.Ping()
	if err != nil {
		return nil, err
	}

	query := "SELECT * FROM SCHEDULE ORDER BY SCHEDULE_ID DESC"

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	schedules := []Schedule{}
	for rows.Next() {
		schedule := Schedule{}
		err = rows.Scan(&schedule.ScheduleId, &schedule.ScheduleName, &schedule.CronExpression, &schedule.MethodName, &schedule.ProjectName)
		if err != nil {
			return nil, err
		}
		schedules = append(schedules, schedule)
	}
	return schedules, nil
}
func (s scheduleRepo) ReadById(id int) (*Schedule, error) {
	err := s.db.Ping()
	if err != nil {
		return nil, err
	}

	query := "SELECT * FROM SCHEDULE WHERE SCHEDULE_ID =?"
	row := s.db.QueryRow(query, id)
	schedule := Schedule{}

	err = row.Scan(&schedule.ScheduleId, &schedule.ScheduleName, &schedule.CronExpression, &schedule.MethodName, &schedule.ProjectName)
	if err != nil {
		return nil, err
	}

	return &schedule, nil
}
func (s scheduleRepo) Create(schedule Schedule) error {
	err := s.db.Ping()
	if err != nil {
		return err
	}

	query := "INSERT INTO SCHEDULE(SCHEDULE_NAME,CRON_EXPRESSION,METHOD_NAME,PROJECT_NAME) VALUES (?,?,?,?)"
	result, err := s.db.Exec(query, schedule.ScheduleName, schedule.CronExpression, schedule.MethodName, schedule.ProjectName)
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
func (s scheduleRepo) Update(schedule Schedule) error {
	err := s.db.Ping()
	if err != nil {
		return err
	}

	query := "UPDATE SCHEDULE SET SCHEDULE_NAME=?, CRON_EXPRESSION=?, METHOD_NAME=?, PROJECT_NAME=?  WHERE SCHEDULE_ID=?"
	result, err := s.db.Exec(query, schedule.ScheduleName, schedule.CronExpression, schedule.MethodName, schedule.ProjectName, schedule.ScheduleId)
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
func (s scheduleRepo) Delete(id int) error {
	err := s.db.Ping()
	if err != nil {
		return err
	}

	query := "DELETE FROM SCHEDULE WHERE SCHEDULE_ID=?"
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
