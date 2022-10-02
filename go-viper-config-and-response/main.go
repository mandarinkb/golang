package main

import (
	"fmt"
	"time"

	"github.com/mandarinkb/go-viper-config-and-response/assets"
	"github.com/mandarinkb/go-viper-config-and-response/config"
	"github.com/mandarinkb/go-viper-config-and-response/utils"
)

const (
	startDate string = "00:41:30"
	endDate   string = "00:42:00"
)

func CurrentLocalTime() string {
	tn := time.Now().In(utils.TimeZone)

	h := convertTime(tn.Hour())
	m := convertTime(tn.Minute())
	s := convertTime(tn.Second())
	return fmt.Sprintf("%v:%v:%v", h, m, s)
}

func convertTime(t int) string {
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

func isInTimeRange(startDate, endDate string) bool {

	check, _ := time.Parse(utils.TimeFormat, CurrentLocalTime())
	start, _ := time.Parse(utils.TimeFormat, startDate)
	end, _ := time.Parse(utils.TimeFormat, endDate)

	//
	if check.Before(start) {
		return false
	}

	if check.Before(end.Add(1 * time.Second)) {
		return true
	}

	return false
}

func main() {
	config.LoadConfig("config", "config")
	assets.LoadAssets("assets", "error")

	for {
		fmt.Println(isInTimeRange(startDate, endDate))
		if isInTimeRange(startDate, endDate) {
			fmt.Println("return error")
		} else {
			fmt.Println("working..")
		}
		fmt.Println(CurrentLocalTime())
		fmt.Println("===========================")
		time.Sleep(1 * time.Second)
	}
}
