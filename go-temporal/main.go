package main

import (
	"log"
	"mandarinkb/go-temporal/activities"
	"mandarinkb/go-temporal/workflows"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {
	// เชื่อมต่อ Temporal Client
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalf("❌ Failed to create Temporal client: %v", err)
	}
	defer c.Close()

	// สร้าง Worker
	w := worker.New(c, "job-task-queue", worker.Options{})
	w.RegisterWorkflow(workflows.DynamicJobWorkflow)
	w.RegisterActivity(activities.ProcessJob)

	// รัน Worker
	log.Println("🚀 Temporal Worker Started!")
	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalf("❌ Worker failed: %v", err)
	}
}
