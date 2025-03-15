package main

import (
	"context"
	"fmt"
	"log"
	"mandarinkb/go-temporal/workflows"

	"go.temporal.io/sdk/client"
)

func main() {
	// เชื่อมต่อ Temporal Client
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalf("❌ Failed to create Temporal client: %v", err)
	}
	defer c.Close()

	// กำหนด Workflow Options
	options := client.StartWorkflowOptions{
		ID:        "dynamic-job-workflow",
		TaskQueue: "job-task-queue",
	}

	// เรียกใช้ Dynamic Workflow พร้อมจำนวนงานที่ต้องการ
	we, err := c.ExecuteWorkflow(context.Background(), options, workflows.DynamicJobWorkflow, 5)
	if err != nil {
		log.Fatalf("❌ Failed to start workflow: %v", err)
	}

	fmt.Printf("✅ Started Dynamic Workflow: %s\n", we.GetID())

	// รับผลลัพธ์ของ Workflow
	var results []string
	err = we.Get(context.Background(), &results)
	if err != nil {
		log.Fatalf("❌ Workflow failed: %v", err)
	}

	// แสดงผลลัพธ์
	fmt.Println("🎉 Workflow results:", results)
}
