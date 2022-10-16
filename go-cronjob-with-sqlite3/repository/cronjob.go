package repository

type CronJob struct {
	CronID         int    `json:"cronID"`
	CronName       string `json:"cronName"`
	CronExpression string `json:"cronExpression"`
	CronUseInFunc  string `json:"cronUseInFunc"`
	CronStatus     string `json:"cronStatus"`
}
type CronJobRepository interface {
	Read() (cronjob []CronJob, err error)
}
