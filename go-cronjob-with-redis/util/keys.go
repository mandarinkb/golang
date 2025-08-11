package util

func JobKey(id string) string {
	return "job_name-" + id
}

func ZsetKey() string {
	return "cron_jobs"
}
