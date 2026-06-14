package util

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"
)

// OnFailurePolicy กำหนดพฤติกรรมเมื่อ step fail
type OnFailurePolicy string

const (
	OnFailureStop     OnFailurePolicy = "stop"     // หยุด workflow ทันที (default)
	OnFailureContinue OnFailurePolicy = "continue" // ข้าม step ที่ fail แล้วรันต่อ
)

// StepStatus state machine ของแต่ละ step
type StepStatus string

const (
	StepStatusPending   StepStatus = "pending"
	StepStatusRunning   StepStatus = "running"
	StepStatusSuccess   StepStatus = "success"
	StepStatusFailed    StepStatus = "failed"
	StepStatusSkipped   StepStatus = "skipped" // เกิดเมื่อ on_failure=continue และ dependency fail
	StepStatusTimeout   StepStatus = "timeout"
)

// WorkflowStatus state machine ของทั้ง workflow run
type WorkflowStatus string

const (
	WorkflowStatusPending   WorkflowStatus = "pending"
	WorkflowStatusRunning   WorkflowStatus = "running"
	WorkflowStatusSuccess   WorkflowStatus = "success"
	WorkflowStatusFailed    WorkflowStatus = "failed"
	WorkflowStatusPartial   WorkflowStatus = "partial" // บาง step skipped แต่ไม่ stop
)

// WorkflowStep คือ node หนึ่งใน DAG
type WorkflowStep struct {
	ID         string   `json:"id"`          // unique ภายใน workflow เช่น "extract"
	JobID      string   `json:"job_id"`      // อ้างถึง standalone job
	DependsOn  []string `json:"depends_on"`  // step IDs ที่ต้องเสร็จก่อน (ว่างได้ = entry point)
	TimeoutSec int      `json:"timeout_sec"` // 0 = ใช้ default jobTimeout
}

// Workflow คือ definition ของ workflow (immutable หลังสร้าง)
type Workflow struct {
	ID        string          `json:"id"`
	Name      string          `json:"name"`
	Schedule  string          `json:"schedule"`           // cron expression
	Steps     []WorkflowStep  `json:"steps"`
	OnFailure OnFailurePolicy `json:"on_failure"`         // default: stop
	Paused    bool            `json:"paused"`
	CreatedAt string          `json:"created_at"`
	UpdatedAt string          `json:"updated_at"`
}

// WorkflowRunState เก็บ runtime state ของ workflow run หนึ่งครั้ง
type WorkflowRunState struct {
	RunID      string                    `json:"run_id"`
	WorkflowID string                    `json:"workflow_id"`
	Status     WorkflowStatus            `json:"status"`
	Steps      map[string]StepRunState   `json:"steps"` // key = step.ID
	StartedAt  string                    `json:"started_at"`
	FinishedAt string                    `json:"finished_at,omitempty"`
}

// StepRunState เก็บ runtime state ของแต่ละ step ใน run นี้
type StepRunState struct {
	Status     StepStatus `json:"status"`
	StartedAt  string     `json:"started_at,omitempty"`
	FinishedAt string     `json:"finished_at,omitempty"`
	ExitCode   int        `json:"exit_code"`
	Output     string     `json:"output,omitempty"`
	Error      string     `json:"error,omitempty"`
	DurationMs int64      `json:"duration_ms,omitempty"`
}

// WorkflowHistoryRecord เก็บ summary ของ workflow run สำหรับ history list
type WorkflowHistoryRecord struct {
	RunID      string         `json:"run_id"`
	Status     WorkflowStatus `json:"status"`
	StartedAt  string         `json:"started_at"`
	FinishedAt string         `json:"finished_at,omitempty"`
	StepCount  int            `json:"step_count"`
	Failed     int            `json:"failed_steps"`
	Skipped    int            `json:"skipped_steps"`
}

// --- Serialization helpers ---

// MarshalWorkflow แปลง Workflow เป็น map[string]interface{} สำหรับ HSET
func MarshalWorkflow(wf *Workflow) (map[string]interface{}, error) {
	stepsJSON, err := json.Marshal(wf.Steps)
	if err != nil {
		return nil, fmt.Errorf("marshal steps: %w", err)
	}
	return map[string]interface{}{
		"id":         wf.ID,
		"name":       wf.Name,
		"schedule":   wf.Schedule,
		"steps":      string(stepsJSON),
		"on_failure": string(wf.OnFailure),
		"paused":     fmt.Sprintf("%v", wf.Paused),
		"created_at": wf.CreatedAt,
		"updated_at": wf.UpdatedAt,
	}, nil
}

// UnmarshalWorkflow แปลง map[string]string จาก HGETALL เป็น Workflow
func UnmarshalWorkflow(data map[string]string) (*Workflow, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("workflow not found")
	}
	var steps []WorkflowStep
	if err := json.Unmarshal([]byte(data["steps"]), &steps); err != nil {
		return nil, fmt.Errorf("unmarshal steps: %w", err)
	}
	policy := OnFailurePolicy(data["on_failure"])
	if policy == "" {
		policy = OnFailureStop
	}
	return &Workflow{
		ID:        data["id"],
		Name:      data["name"],
		Schedule:  data["schedule"],
		Steps:     steps,
		OnFailure: policy,
		Paused:    data["paused"] == "true",
		CreatedAt: data["created_at"],
		UpdatedAt: data["updated_at"],
	}, nil
}

// MarshalRunState แปลง WorkflowRunState เป็น map สำหรับ HSET
func MarshalRunState(state *WorkflowRunState) (map[string]interface{}, error) {
	stepsJSON, err := json.Marshal(state.Steps)
	if err != nil {
		return nil, fmt.Errorf("marshal run steps: %w", err)
	}
	return map[string]interface{}{
		"run_id":      state.RunID,
		"workflow_id": state.WorkflowID,
		"status":      string(state.Status),
		"steps":       string(stepsJSON),
		"started_at":  state.StartedAt,
		"finished_at": state.FinishedAt,
	}, nil
}

// UnmarshalRunState แปลง map[string]string จาก HGETALL เป็น WorkflowRunState
func UnmarshalRunState(data map[string]string) (*WorkflowRunState, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("run state not found")
	}
	var steps map[string]StepRunState
	if err := json.Unmarshal([]byte(data["steps"]), &steps); err != nil {
		return nil, fmt.Errorf("unmarshal run steps: %w", err)
	}
	return &WorkflowRunState{
		RunID:      data["run_id"],
		WorkflowID: data["workflow_id"],
		Status:     WorkflowStatus(data["status"]),
		Steps:      steps,
		StartedAt:  data["started_at"],
		FinishedAt: data["finished_at"],
	}, nil
}

// ValidateWorkflow ตรวจสอบ workflow definition
// - step ID ต้องไม่ซ้ำ
// - depends_on ต้องอ้างถึง step ที่มีอยู่จริง
// - ต้องไม่มี cycle (topological sort)
// - job_id ทุกตัวต้องมีใน Redis
func ValidateWorkflow(wf *Workflow) error {
	if wf.Name == "" {
		return fmt.Errorf("workflow name is required")
	}
	if wf.Schedule == "" {
		return fmt.Errorf("workflow schedule is required")
	}
	if len(wf.Steps) == 0 {
		return fmt.Errorf("workflow must have at least one step")
	}

	stepIDs := make(map[string]bool, len(wf.Steps))
	for _, s := range wf.Steps {
		if s.ID == "" {
			return fmt.Errorf("all steps must have an id")
		}
		if s.JobID == "" {
			return fmt.Errorf("step %q must have a job_id", s.ID)
		}
		if stepIDs[s.ID] {
			return fmt.Errorf("duplicate step id %q", s.ID)
		}
		stepIDs[s.ID] = true
	}

	// ตรวจ depends_on references
	for _, s := range wf.Steps {
		for _, dep := range s.DependsOn {
			if !stepIDs[dep] {
				return fmt.Errorf("step %q depends_on %q which does not exist", s.ID, dep)
			}
		}
	}

	// cycle detection: Kahn's algorithm
	if err := detectCycle(wf.Steps); err != nil {
		return err
	}

	return nil
}

// detectCycle ใช้ topological sort (Kahn's algorithm) ตรวจ cycle ใน DAG
func detectCycle(steps []WorkflowStep) error {
	inDegree := make(map[string]int, len(steps))
	for _, s := range steps {
		if _, ok := inDegree[s.ID]; !ok {
			inDegree[s.ID] = 0
		}
		for _, dep := range s.DependsOn {
			inDegree[dep] = inDegree[dep] // ensure key exists
			inDegree[s.ID]++
		}
	}

	queue := []string{}
	for id, deg := range inDegree {
		if deg == 0 {
			queue = append(queue, id)
		}
	}

	visited := 0
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		visited++
		for _, s := range steps {
			for _, dep := range s.DependsOn {
				if dep == cur {
					inDegree[s.ID]--
					if inDegree[s.ID] == 0 {
						queue = append(queue, s.ID)
					}
				}
			}
		}
	}

	if visited != len(steps) {
		return fmt.Errorf("workflow steps contain a cycle")
	}
	return nil
}

// GetReadySteps คืน step IDs ที่พร้อมรัน:
// - ยังไม่ถูก execute (pending)
// - dependency ทั้งหมด success หรือ skipped แล้ว
func GetReadySteps(wf *Workflow, state *WorkflowRunState) []string {
	ready := []string{}
	for _, step := range wf.Steps {
		s, ok := state.Steps[step.ID]
		if !ok || s.Status != StepStatusPending {
			continue
		}
		allDepsDone := true
		for _, dep := range step.DependsOn {
			depState, exists := state.Steps[dep]
			if !exists {
				allDepsDone = false
				break
			}
			switch depState.Status {
			case StepStatusSuccess, StepStatusSkipped:
				// dependency สำเร็จหรือถูก skip → ถือว่าผ่าน
			default:
				allDepsDone = false
			}
		}
		if allDepsDone {
			ready = append(ready, step.ID)
		}
	}
	return ready
}

// IsWorkflowDone ตรวจว่า workflow run เสร็จแล้วหรือยัง
// เสร็จ = ทุก step ไม่ใช่ pending หรือ running
func IsWorkflowDone(state *WorkflowRunState) bool {
	for _, s := range state.Steps {
		if s.Status == StepStatusPending || s.Status == StepStatusRunning {
			return false
		}
	}
	return true
}

// ComputeWorkflowStatus คำนวณ final status จาก step states
func ComputeWorkflowStatus(state *WorkflowRunState) WorkflowStatus {
	hasFailed := false
	hasSkipped := false
	for _, s := range state.Steps {
		switch s.Status {
		case StepStatusFailed, StepStatusTimeout:
			hasFailed = true
		case StepStatusSkipped:
			hasSkipped = true
		}
	}
	if hasFailed && !hasSkipped {
		return WorkflowStatusFailed
	}
	if hasFailed && hasSkipped {
		return WorkflowStatusPartial
	}
	if hasSkipped {
		return WorkflowStatusPartial
	}
	return WorkflowStatusSuccess
}

// SaveWorkflowRun บันทึก run state กลับ Redis
func SaveWorkflowRun(ctx context.Context, runID, workflowID string, state *WorkflowRunState) error {
	fields, err := MarshalRunState(state)
	if err != nil {
		return err
	}
	key := WorkflowRunKey(workflowID, runID)
	if err := RDB.HSet(ctx, key, fields).Err(); err != nil {
		return fmt.Errorf("save run state: %w", err)
	}
	// TTL 7 วัน ไม่ให้ run state เก่าค้างใน Redis ตลอดกาล
	RDB.Expire(ctx, key, 7*24*time.Hour)
	return nil
}

// AppendWorkflowHistory บันทึก summary ลง history list
func AppendWorkflowHistory(ctx context.Context, workflowID string, rec WorkflowHistoryRecord) {
	data, err := json.Marshal(rec)
	if err != nil {
		Logger.Error("marshal workflow history failed", slog.String("wf_id", workflowID), slog.Any("error", err))
		return
	}
	pipe := RDB.Pipeline()
	hk := WorkflowHistoryKey(workflowID)
	pipe.LPush(ctx, hk, string(data))
	pipe.LTrim(ctx, hk, 0, 49) // เก็บ 50 records
	if _, err := pipe.Exec(ctx); err != nil {
		Logger.Error("append workflow history failed", slog.String("wf_id", workflowID), slog.Any("error", err))
	}
}
