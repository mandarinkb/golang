package utils

import "github.com/mandarinkb/go-cronjob-with-sqlite3/repository"

func RemoveElement(s []repository.CronJobData, i int) []repository.CronJobData {
	s[i] = s[len(s)-1]
	s = s[:len(s)-1]
	return s
}
