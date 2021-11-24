package repository

type Schedule struct {
	ScheduleId     int    `db:"SCHEDULE_ID"`
	ScheduleName   string `db:"SCHEDULE_NAME"`
	CronExpression string `db:"CRON_EXPRESSION"`
	MethodName     string `db:"METHOD_NAME"`
	ProjectName    string `db:"PROJECT_NAME"`
}

type ScheduleRepository interface {
	Read() ([]Schedule, error)
	ReadById(id int) (*Schedule, error)
	Create(schedule Schedule) error
	Update(schedule Schedule) error
	Delete(id int) error
}
