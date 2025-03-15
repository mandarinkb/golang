package workflows

import (
	"fmt"
	"mandarinkb/go-temporal/activities"

	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

// 📌 Workflow: เพิ่ม Retry Policy ใน ActivityOptions
func JobWorkflowWithRetry(ctx workflow.Context, jobID string) (string, error) {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Second * 10,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second * 1,
			BackoffCoefficient: 2.0,
			MaximumInterval:    time.Second * 5,
			MaximumAttempts:    3, // retry 3 ครั้ง
		},
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	var result string
	err := workflow.ExecuteActivity(ctx, activities.ProcessJob, jobID).Get(ctx, &result)
	if err != nil {
		return "", err
	}
	return result, nil
}

// 📌 Workflow: ทำ Parallel Execution ของหลาย Activity
func ParallelJobWorkflow(ctx workflow.Context, jobIDs []string) ([]string, error) {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	results := make([]string, len(jobIDs))
	futures := make([]workflow.Future, len(jobIDs))

	// 🔹 เริ่ม Activity หลายตัวพร้อมกัน
	for i, jobID := range jobIDs {
		futures[i] = workflow.ExecuteActivity(ctx, activities.ProcessJob, jobID)
	}

	// 🔹 รอให้ Activity ทั้งหมดเสร็จ
	for i, future := range futures {
		err := future.Get(ctx, &results[i])
		if err != nil {
			return nil, err
		}
	}

	return results, nil
}

// 📌 Dynamic Workflow: สร้างจำนวนงานที่ต้องทำตอน runtime
func DynamicJobWorkflow(ctx workflow.Context, numJobs int) ([]string, error) {
	jobIDs := make([]string, numJobs)
	for i := 0; i < numJobs; i++ {
		jobIDs[i] = fmt.Sprintf("job-%d", i+1)
	}
	return ParallelJobWorkflow(ctx, jobIDs)
}
