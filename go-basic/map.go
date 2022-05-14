package main

import (
	"encoding/json"
	"fmt"
)

func main() {
	// map[Key_Type]Value_Type{key1: value1, ..., keyN: valueN}
	// ใช้ map กรณีมีหลายค่า ,ค่า Key ต้องไม่ซ้ำกัน
	mapValue := map[string]string{"user": "joke", "user2": "joke"}
	fmt.Println(mapValue)
	js, _ := json.Marshal(mapValue)
	fmt.Println(string(js))
}
