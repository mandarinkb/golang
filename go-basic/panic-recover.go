package main

// func main() {

// 	fmt.Println("Hello World")
// 	// ใช้คู่กับกรณีเกิด panic
// 	// โดยจะให้ทำอะไรต่อหลังเกิด panic
// 	defer func() {
// 		r := recover()
// 		if r == "close db" {
// 			fmt.Println("call recover")
// 		}
// 	}()

// 	//  กรณีเกิด error ที่ส่งผลกระทบต่อระบบ ให้ใช้ panic จะหยุดการทำงานทันที
// 	// เช่น เชื่อมต่อ database ไม่ได้ , เปิดไฟล์ไม่ได้
// 	i := 1
// 	if i == 1 {
// 		panic("close db")
// 	}

// 	func() {
// 		fmt.Println("call function")
// 	}()
// }
