package main

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/sony/gobreaker"
)

// Mock ที่ล้มตลอด
func alwaysFailingAPI() (*http.Response, error) {
	time.Sleep(500 * time.Millisecond) // จำลองว่า API ช้า
	return nil, errors.New("mock error")
}

// --- Case 1: HTTP client ธรรมดา ---
func callAPIWithoutBreaker() error {
	_, err := alwaysFailingAPI()
	return err
}

// --- Case 2: HTTP client ห่อด้วย Circuit Breaker ---
var cb *gobreaker.CircuitBreaker

func init() {
	cb = gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:        "QRCodeBreaker",
		Interval:    30 * time.Second, // ล้าง count ทุก 30 วิ (ลด false open)
		Timeout:     5 * time.Second,  // breaker เปิด 5 วิ
		MaxRequests: 1,                // ใน HALF-OPEN อนุญาต 1 request เท่านั้น
		ReadyToTrip: func(c gobreaker.Counts) bool {
			return c.ConsecutiveFailures >= 3 // ถ้า error 3 ครั้งติดกัน → เปิดวงจร
		},
		OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
			log.Printf("[CircuitBreaker] %s: %s → %s\n", name, from.String(), to.String())
		},
	})
}

func callAPIWithBreaker() error {
	_, err := cb.Execute(func() (interface{}, error) {
		return alwaysFailingAPI()
	})
	return err
}
