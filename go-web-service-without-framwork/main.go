package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type User struct {
	USER_ID  int    `json:"userId"`
	USERNAME string `json:"username"`
	PASSWORD string `json:"password"`
	ROLE     string `json:"role"`
}

var users = []User{
	{USER_ID: 1, USERNAME: "test", PASSWORD: "test-password", ROLE: "admin"},
	{USER_ID: 2, USERNAME: "test1", PASSWORD: "test1-password", ROLE: "assistant"},
	{USER_ID: 3, USERNAME: "test2", PASSWORD: "test2-password", ROLE: "admin"},
	{USER_ID: 4, USERNAME: "test3", PASSWORD: "test3-password", ROLE: "assistant"},
}

type Respose struct {
	Message string `json:"message"`
}

// หน้าแรก
func homePage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json") // set เป็นค่า json
	res := Respose{
		Message: "Hello rest api",
	}
	w.WriteHeader(http.StatusOK)   // set http status
	json.NewEncoder(w).Encode(res) // return ค่า json ออกให้เลย
}

func getAllUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
}

func getUsersById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	params := mux.Vars(r)
	idStr := params["id"] //รับ id มา

	id, err := strconv.Atoi(idStr) // แปลงเป็น int
	if err != nil {
		fmt.Println(err)
	}
	// ตรวจเช็ค id ที่รับเข้ามาก่อน (คล้ายกับ select id ใน database)
	haveId := false
	for _, user := range users {
		if user.USER_ID == id {
			haveId = true
		}
	}
	if haveId { // กรณีพบ id
		for i, user := range users {
			if user.USER_ID == id {
				newUser := User{ //User strucrt ต้องขึ้นต้นด้วยตัวใหญ่ไม่งั้น error
					USER_ID:  users[i].USER_ID,
					USERNAME: users[i].USERNAME,
					PASSWORD: users[i].PASSWORD,
					ROLE:     users[i].ROLE,
				}
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(newUser)
				return
			}
		}
	} else { //กรณีไม่พบ
		res := Respose{
			Message: "not found user",
		}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(res)
	}
}

func createUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	var newUser User
	body, err := ioutil.ReadAll(r.Body) // รับค่าจาก body payload
	if err != nil {
		fmt.Println(err)
	}
	json.Unmarshal(body, &newUser) // แปลงค่าจาก payload เป็น struct
	users = append(users, newUser) // เก็บลง slice struct
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newUser)
}
func updateUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	params := mux.Vars(r)
	idStr := params["id"] //รับ id มา

	id, err := strconv.Atoi(idStr) // แปลงเป็น int
	if err != nil {
		fmt.Println(err)
	}
	var bodyValue User
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
	}
	json.Unmarshal(body, &bodyValue) // แปลงค่าจาก payload เป็น struct

	var newUser User     // ใช้เก็บ user ใหม่ และ user เก่าที่ id ไม่ตรงกัน
	var userUpdate User  // ใช้แสดงค่าที่ อัพเดท
	var sliceUser []User // ไว้ อัพเดท slice User

	// ตรวจเช็ค id ที่รับเข้ามาก่อน
	haveId := false
	for _, user := range users {
		if user.USER_ID == id {
			haveId = true
		}
	}
	if haveId { // กรณีพบ id
		for i, user := range users {
			if user.USER_ID == id {
				newUser = User{
					USER_ID:  users[i].USER_ID,
					USERNAME: bodyValue.USERNAME,
					PASSWORD: bodyValue.PASSWORD,
					ROLE:     bodyValue.ROLE,
				}
				userUpdate = newUser
				sliceUser = append(sliceUser, newUser)
			} else {
				newUser = User{
					USER_ID:  users[i].USER_ID,
					USERNAME: users[i].USERNAME,
					PASSWORD: users[i].PASSWORD,
					ROLE:     users[i].ROLE,
				}
				sliceUser = append(sliceUser, newUser)
			}
		}
		users = sliceUser
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(userUpdate)
	} else { //กรณีไม่พบ
		res := Respose{
			Message: "not found user",
		}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(res)
	}
}

func deleteUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	params := mux.Vars(r)
	idStr := params["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Println(err)
	}

	var newUser User     // ใช้เก็บ user ใหม่ และ user เก่าที่ id ไม่ตรงกัน
	var sliceUser []User // ไว้ อัพเดท slice User

	// ตรวจเช็ค id ที่รับเข้ามาก่อน
	haveId := false
	for _, user := range users {
		if user.USER_ID == id {
			haveId = true
		}
	}
	if haveId { // กรณีพบ id
		for i, user := range users {
			if user.USER_ID != id { // id ที่เหมือนกันลบทิ้ง
				newUser = User{
					USER_ID:  users[i].USER_ID,
					USERNAME: users[i].USERNAME,
					PASSWORD: users[i].PASSWORD,
					ROLE:     users[i].ROLE,
				}
				sliceUser = append(sliceUser, newUser)
			}
		}
		users = sliceUser //อัพเดทค่าใน slice
		res := Respose{
			Message: "delete user success",
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(res)
	} else { //กรณีไม่พบ
		res := Respose{
			Message: "not found user",
		}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(res)
	}
}

func handleRequest() {
	// ใช้ routing ของ  gorilla/mux มาช่วยจัดการ
	r := mux.NewRouter()
	r.HandleFunc("/", homePage)
	r.HandleFunc("/users", getAllUsers).Methods("GET")
	r.HandleFunc("/users/{id}", getUsersById).Methods("GET")
	r.HandleFunc("/users", createUsers).Methods("POST")
	r.HandleFunc("/users/{id}", updateUsers).Methods("PUT")
	r.HandleFunc("/users/{id}", deleteUsers).Methods("DELETE")
	http.ListenAndServe(":8080", r)
}
func main() {
	fmt.Println("service running at port 8080")
	handleRequest()
}
