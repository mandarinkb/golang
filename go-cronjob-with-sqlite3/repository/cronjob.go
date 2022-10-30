package repository

import "github.com/robfig/cron/v3"

var (
	Cron *cron.Cron = cron.New()
	//store cron job detail
	CronJobDetail []CronJobData
)

type CronJobData struct {
	CronEntryID     cron.EntryID
	CronFunctionRef string
}

type CronJob struct {
	CronID          int    `json:"cronID" db:"CRONID"`
	CronName        string `json:"cronName" db:"CRONNAME"`
	CronExpression  string `json:"cronExpression" db:"CRONEXPRESSION"`
	CronFunctionRef string `json:"cronFunctionRef" db:"CRONFUNCTIONREF"`
	CronStatus      int    `json:"cronStatus" db:"CRONSTATUS"`
}
type CronJobRepository interface {
	GetCronJob() (cronjob []CronJob, err error)
	GetCronJobByID(id int) (*CronJob, error)
	CreateCronJob(cronjob CronJob) error
	UpdateCronJob(cronjob CronJob) error
	DeleteCronJob(id int) error
	RunJob(cronExpression string, cmd func()) (cron.EntryID, error)
	RemoveJob(id cron.EntryID)
}
