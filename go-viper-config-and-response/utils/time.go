package utils

import "time"

const (
	DateTimeFormat string = "2006-01-02T15:04:05.000+07:00"
	DateFormat     string = "20060102"
	TimeFormat     string = "15:04:05"
)

var TimeZone, _ = time.LoadLocation("Asia/Bangkok")

func CurrentLocalDate() string {
	return time.Now().In(TimeZone).Format(DateTimeFormat)
}

func TimeFormated(t time.Time) string {
	return t.In(TimeZone).Format(DateTimeFormat)
}

func DateFormated(t time.Time) string {
	return t.In(TimeZone).Format(DateFormat)
}
