package utils

import (
	"fmt"
	"time"
)

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

func MaintenancePeriod(start, end string) bool {
	tn := time.Now().In(TimeZone)

	h := addPrefixTime(tn.Hour())
	m := addPrefixTime(tn.Minute())
	s := addPrefixTime(tn.Second())
	current := fmt.Sprintf("%v:%v:%v", h, m, s)

	cuerrentTime, _ := time.Parse(TimeFormat, current)
	startTime, _ := time.Parse(TimeFormat, start)
	endTime, _ := time.Parse(TimeFormat, end)

	return IsInTimeRange(startTime, endTime, cuerrentTime)
}

func IsInTimeRange(startDate, endDate, current time.Time) bool {

	if current.Before(startDate) {
		return false
	}

	if current.Before(endDate.Add(1 * time.Second)) {
		return true
	}

	return false
}

func addPrefixTime(t int) string {
	var newT string
	switch t {
	case 0:
		newT = fmt.Sprintf("0%v", t)
	case 1:
		newT = fmt.Sprintf("0%v", t)
	case 2:
		newT = fmt.Sprintf("0%v", t)
	case 3:
		newT = fmt.Sprintf("0%v", t)
	case 4:
		newT = fmt.Sprintf("0%v", t)
	case 5:
		newT = fmt.Sprintf("0%v", t)
	case 6:
		newT = fmt.Sprintf("0%v", t)
	case 7:
		newT = fmt.Sprintf("0%v", t)
	case 8:
		newT = fmt.Sprintf("0%v", t)
	case 9:
		newT = fmt.Sprintf("0%v", t)
	default:
		newT = fmt.Sprintf("%v", t)
	}
	return newT

}

func IsDateEqual(date1, date2 time.Time) bool {
	y1, m1, d1 := date1.Date()
	y2, m2, d2 := date2.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}
