package main

import (
	"fmt"
	"log"
	"testing"
	"time"
)

// --- Unit Test เปรียบเทียบ ---
func TestCompareBreakerVsNormalClient(t *testing.T) {
	// ------- 1. ทดสอบ HTTP Client ปกติ -------
	t.Run("HTTP Client ธรรมดา ล้ม 3 ครั้ง ยังยิงต่อ (ใช้เวลานานทุกครั้ง)", func(t *testing.T) {
		for i := 0; i < 5; i++ {
			start := time.Now()
			err := callAPIWithoutBreaker()
			duration := time.Since(start)

			log.Printf("[HTTP Client %d] Error: %v | Duration: %v\n", i+1, err, duration)
		}
		log.Printf("=> HTTP Client ธรรมดายิงทุกครั้ง แม้จะ error และ delay ทุกครั้ง")
		fmt.Println()
	})

	// ------- 2. ทดสอบ Circuit Breaker -------
	t.Run("Circuit Breaker ล้ม 3 ครั้ง จากนั้น fail fast (ใช้เวลาแค่บางรอบแรก)", func(t *testing.T) {
		for i := 0; i < 40; i++ {
			start := time.Now()
			err := callAPIWithBreaker()
			duration := time.Since(start)

			// ตัดสินว่า fail fast หรือไม่จากเวลา
			if duration < 50*time.Millisecond {
				log.Printf("[CB %d] Fail Fast  | Error: %v | Duration: %v\n", i+1, err, duration)
			} else {
				log.Printf("[CB %d] CALL TRIED | Error: %v | Duration: %v\n", i+1, err, duration)
			}
			time.Sleep(1 * time.Second)
		}
		log.Printf("=> Circuit Breaker หยุดยิงหลังล้ม 3 ครั้ง และตอบกลับเร็ว (fail fast)")
	})
}
