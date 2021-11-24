package service

type Schedule struct {
	ScheduleId     int    `json:"scheduleId"`
	ScheduleName   string `json:"scheduleName"`
	CronExpression string `json:"cronExpression"`
	MethodName     string `json:"methodName"`
	ProjectName    string `json:"projectName"`
}
type ScheduleService interface {
	Read() ([]Schedule, error)
	ReadById(id int) (*Schedule, error)
	Create(schedule Schedule) error
	Update(schedule Schedule) error
	Delete(id int) error
}
