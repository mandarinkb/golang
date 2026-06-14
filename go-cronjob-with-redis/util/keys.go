package util

// --- Standalone Job keys ---

// JobKey เก็บ metadata ของ job (Hash)
func JobKey(id string) string {
	return "job:" + id
}

// ZsetKey เป็น priority queue ของ standalone jobs (ZSet, score=nextRun)
func ZsetKey() string {
	return "cron_jobs"
}

// ProcessingKey tracking standalone job ที่กำลัง execute (ZSet, score=startedAt)
func ProcessingKey() string {
	return "cron_processing"
}

// JobHistoryKey เก็บ execution history ของ standalone job (List, max 100)
func JobHistoryKey(id string) string {
	return "job_history:" + id
}

// --- Workflow keys ---

// WorkflowKey เก็บ metadata ของ workflow (Hash)
func WorkflowKey(id string) string {
	return "workflow:" + id
}

// WorkflowZsetKey เป็น priority queue ของ workflows (ZSet, score=nextRun)
func WorkflowZsetKey() string {
	return "workflow_runs"
}

// WorkflowRunKey เก็บ runtime state ของ workflow run หนึ่งครั้ง (Hash)
// runID = unique id ต่อ execution round
func WorkflowRunKey(workflowID, runID string) string {
	return "workflow_run:" + workflowID + ":" + runID
}

// WorkflowRunZsetKey tracking workflow run ที่กำลัง active (ZSet, score=startedAt)
func WorkflowRunZsetKey() string {
	return "workflow_processing"
}

// WorkflowHistoryKey เก็บ execution history ของ workflow (List, max 50)
func WorkflowHistoryKey(id string) string {
	return "workflow_history:" + id
}
