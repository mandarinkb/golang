package main

// func main() {
// 	//###test goroutine###
// 	// go func() {
// 	// 	fmt.Println("func 1")
// 	// }()

// 	// go func(msg string) {
// 	// 	fmt.Println(msg)
// 	// }("going")

// 	// str := []string{"sub1 func2", "sub2 func2"}
// 	// for _, s := range str {
// 	// 	go func(ss string) {
// 	// 		fmt.Println("ss : ", ss)
// 	// 	}(s) // s จะส่งค่าไปยัง ss

// 	// }

// 	// func3 := func() {
// 	// 	fmt.Println("func 3")
// 	// }
// 	// go func3()

// 	// // ต้องมีคำสั่ง time.Sleep เพื่อรอ บรรทัดที่มีคำสั่ง go ทำงานเสร็จ
// 	// // มิฉะนั้น ใน function main() จะปิดตัวทันที ไม่ทันให้ function อื่นแสดงผล
// 	// time.Sleep(time.Second)

// 	//###test channal###
// 	// fmt.Println("main")
// 	// name := "World"
// 	// result := make(chan string)
// 	// go hello(name, result)

// 	// msg := <-result //แปลงค่าที่อยู่ใน channel เป็น string
// 	// fmt.Println(msg)

// 	//###test channal ex2###
// 	resultA := make(chan int)
// 	resultB := make(chan int)
// 	go A(9, resultA)
// 	go B(7, resultB)
// 	intA := <-resultA
// 	intB := <-resultB
// 	fmt.Println(intA - intB)

// 	// มีงานบางประเภทที่เราอยากให้ Gotoutine นั้นหยุดรอการประมวลผลให้เสร็จทุกชิ้นก่อน ที่จะไปรันคำสั่งถัดไป
// 	// var wg sync.WaitGroup
// 	// บอก Go ว่ามีจำนวน Gotoutine ที่ต้องหยุดรอ ให้ใช้คำสั่ง
// 	// wg.Add(1)
// 	// เมื่อ Goroutine ใดๆก็ตามทำงานเสร็จแล้ว เราจำเป็นต้องบอก Go เพื่อลดจำนวนที่ต้องหยุดรอลงไป ด้วยคำสั่ง
// 	// wg.Done()
// 	// และการหยุดรอ Goroutine ให้เสร็จทั้งหมด ก่อนทำคำสั่งถัดไปได้ ต้องใช้คำสั่ง
// 	// wg.Wait()
// }

// //###test channal###
// // func hello(name string, result chan string) {
// // 	msg := "Hello " + name
// // 	result <- msg // ส่งค่าเข้า channel
// // }

// //###test channal ex2###
// func A(a int, result chan int) {
// 	result <- a
// }

// func B(b int, result chan int) {
// 	result <- b
// }
