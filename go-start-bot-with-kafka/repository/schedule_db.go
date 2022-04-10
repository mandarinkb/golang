package repository

import "database/sql"

type scheduleRepo struct {
	db *sql.DB
}

func NewScheduleRepo(db *sql.DB) ScheduleRepository {
	return scheduleRepo{db: db}
}

func (s scheduleRepo) Read() (schedules []Schedule, err error) {
	err = s.db.Ping()
	if err != nil {
		return nil, err
	}

	query := "SELECT * FROM SCHEDULE ORDER BY SCHEDULE_ID DESC"

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

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
