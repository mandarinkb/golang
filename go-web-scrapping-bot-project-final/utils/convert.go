package utils

import (
	"fmt"
	"strconv"
)

func FloatToString(f float64) string {
	return fmt.Sprintf("%v", f)
}
func StrToFloat64(str string) float64 {
	pTemp, _ := strconv.ParseFloat(str, 32)
	return float64(pTemp)
}
