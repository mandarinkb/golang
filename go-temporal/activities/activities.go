package activities

import (
	"context"
	"fmt"
	"time"
)

// 📌 Activity: จำลองการประมวลผลงาน
func ProcessJob(ctx context.Context, jobID string) (string, error) {
	// จำลองการทำงาน
	fmt.Printf("🔄 Processing job ID: %s\n", jobID)
	time.Sleep(2 * time.Second)

	// สมมุติให้เกิดข้อผิดพลาดบางอย่าง (ลองใช้ Retry)
	if jobID == "job-3" {
		return "", fmt.Errorf("❌ Job %s failed", jobID)
	}

	return fmt.Sprintf("✅ Job %s processed successfully", jobID), nil
}
