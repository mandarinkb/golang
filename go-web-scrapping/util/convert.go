package util

import (
	"fmt"
	"strconv"
)

func FloatToString(f float32) string {
	return fmt.Sprintf("%v", f)
}
func StrToFloat32(str string) float32 {
	pTemp, _ := strconv.ParseFloat(str, 32)
	return float32(pTemp)
}
