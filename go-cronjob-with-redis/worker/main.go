package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"log/slog"

	"github.com/google/uuid"
	"github.com/mandarinkb/go-cronjob-with-redis/util"
	"github.com/redis/go-redis/v9"
)

const (
	pollInterval   = 2 * time.Second
	defaultTimeout = 30 * time.Second
	stuckThreshold = 5 * time.Minute
	historyMaxLen  = 100

	// maxConcurrentWorkflows จำกัดจำนวน workflow run ที่รันพร้อมกันได้
	// ป้องกัน goroutine leak เมื่อมี workflow จำนวนมาก
	maxConcurrentWorkflows = 10

	// maxConcurrentStandaloneJobs จำกัดจำนวน standalone job ที่รันพร้อมกัน
	maxConcurrentStandaloneJobs = 20

	// shutdownTimeout เวลาสูงสุดที่รอ graceful shutdown
	shutdownTimeout = 60 * time.Second
)

// ── Helpers ───────────────────────────────────────────────────────────────────

// generateRunID สร้าง unique run ID ที่ปลอดภัยจาก panic (ไม่ใช้ slice index)
func generateRunID(workflowID string) string {
	prefix := workflowID
	if len(prefix) > 8 {
		prefix = prefix[:8]
	}
	suffix := strings.ReplaceAll(uuid.New().String(), "-", "")[:8]
	return fmt.Sprintf("%s-%s", prefix, suffix)
}

// ── Shared execution ──────────────────────────────────────────────────────────

// JobHistoryRecord เก็บ log ของ standalone job execution
type JobHistoryRecord struct {
	ExecutedAt string `json:"executed_at"`
	Status     string `json:"status"` // success | failed | timeout | recovered
	ExitCode   int    `json:"exit_code"`
	Output     string `json:"output,omitempty"`
	Error      string `json:"error,omitempty"`
	DurationMs int64  `json:"duration_ms"`
}

// executeCommand รัน shell command พร้อม context timeout
func executeCommand(ctx context.Context, command string) (output string, exitCode int, err error) {
	cmd := exec.CommandContext(ctx, "bash", "-c", command)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf

	err = cmd.Run()
	output = buf.String()
	if len(output) > 2000 {
		output = output[:2000] + "... [truncated]"
	}
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			exitCode = -1
		}
	}
	return
}

// timeoutForStep คืน timeout ของ step (ถ้าไม่ได้กำหนดใช้ default)
func timeoutForStep(step util.WorkflowStep) time.Duration {
	if step.TimeoutSec > 0 {
		return time.Duration(step.TimeoutSec) * time.Second
	}
	return defaultTimeout
}

// ── Standalone Job processing ─────────────────────────────────────────────────

func saveJobHistory(ctx context.Context, jobID string, record JobHistoryRecord) {
	data, err := json.Marshal(record)
	if err != nil {
		util.Logger.Error("marshal job history failed", slog.String("id", jobID), slog.Any("error", err))
		return
	}
	pipe := util.RDB.Pipeline()
	pipe.LPush(ctx, util.JobHistoryKey(jobID), string(data))
	pipe.LTrim(ctx, util.JobHistoryKey(jobID), 0, int64(historyMaxLen-1))
	if _, err := pipe.Exec(ctx); err != nil {
		util.Logger.Error("save job history failed", slog.String("id", jobID), slog.Any("error", err))
	}
}

func recoverStuckJobs(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			cutoff := fmt.Sprintf("%f", float64(time.Now().Add(-stuckThreshold).Unix()))
			stuck, err := util.RDB.ZRangeByScoreWithScores(ctx, util.ProcessingKey(), &redis.ZRangeBy{
				Min: "-inf", Max: cutoff,
			}).Result()
			if err != nil {
				util.Logger.Error("stuck job check failed", slog.Any("error", err))
				continue
			}
			for _, z := range stuck {
				jobID := z.Member.(string)
				util.Logger.Warn("recovering stuck job", slog.String("id", jobID))
				schedule, err := util.RDB.HGet(ctx, util.JobKey(jobID), "schedule").Result()
				if err != nil {
					util.RDB.ZRem(ctx, util.ProcessingKey(), jobID)
					continue
				}
				nextRun, err := util.GetNextRun(schedule)
				if err != nil {
					util.RDB.ZRem(ctx, util.ProcessingKey(), jobID)
					continue
				}
				pipe := util.RDB.Pipeline()
				pipe.ZRem(ctx, util.ProcessingKey(), jobID)
				pipe.ZAdd(ctx, util.ZsetKey(), redis.Z{Score: float64(nextRun), Member: jobID})
				if _, err := pipe.Exec(ctx); err != nil {
					util.Logger.Error("requeue stuck job failed", slog.String("id", jobID), slog.Any("error", err))
					continue
				}
				util.Logger.Info("stuck job requeued", slog.String("id", jobID))
				saveJobHistory(ctx, jobID, JobHistoryRecord{
					ExecutedAt: time.Now().Format(time.RFC3339),
					Status:     "recovered",
					ExitCode:   -1,
					Error:      "job was stuck in processing and has been recovered",
				})
			}
		}
	}
}

// processStandaloneJobs poll และรัน standalone jobs
// wg ใช้สำหรับ graceful shutdown — รอให้ job ที่กำลังรันเสร็จก่อน
func processStandaloneJobs(ctx context.Context, wg *sync.WaitGroup) {
	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	// semaphore จำกัด concurrent standalone jobs
	sem := make(chan struct{}, maxConcurrentStandaloneJobs)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			now := time.Now().Unix()
			res, err := util.RDB.ZPopMin(ctx, util.ZsetKey(), 1).Result()
			if err != nil && err != redis.Nil {
				util.Logger.Error("zpopmin (jobs) failed", slog.Any("error", err))
				continue
			}
			if len(res) == 0 {
				continue
			}
			z := res[0]
			if int64(z.Score) > now {
				util.RDB.ZAdd(ctx, util.ZsetKey(), redis.Z{Score: z.Score, Member: z.Member})
				continue
			}

			jobID := z.Member.(string)
			job, err := util.RDB.HGetAll(ctx, util.JobKey(jobID)).Result()
			if err != nil || len(job) == 0 {
				util.Logger.Warn("standalone job data missing", slog.String("id", jobID))
				continue
			}
			if job["paused"] == "true" {
				nextRun, err := util.GetNextRun(job["schedule"])
				if err != nil {
					continue
				}
				util.RDB.ZAdd(ctx, util.ZsetKey(), redis.Z{Score: float64(nextRun), Member: jobID})
				continue
			}

			// ตรวจ semaphore — ถ้าเต็มให้ push กลับ queue แล้วรอรอบถัดไป
			select {
			case sem <- struct{}{}:
			default:
				util.Logger.Warn("standalone job semaphore full, requeueing",
					slog.String("id", jobID))
				util.RDB.ZAdd(ctx, util.ZsetKey(), redis.Z{
					Score:  float64(time.Now().Add(pollInterval).Unix()),
					Member: jobID,
				})
				continue
			}

			// mark processing
			if err := util.RDB.ZAdd(ctx, util.ProcessingKey(), redis.Z{
				Score: float64(time.Now().Unix()), Member: jobID,
			}).Err(); err != nil {
				<-sem
				util.RDB.ZAdd(ctx, util.ZsetKey(), redis.Z{Score: z.Score, Member: jobID})
				continue
			}

			wg.Add(1)
			capturedJob := job
			capturedID := jobID
			capturedScore := z.Score
			go func() {
				defer wg.Done()
				defer func() { <-sem }()

				util.Logger.Info("executing standalone job",
					slog.String("id", capturedID), slog.String("name", capturedJob["name"]))

				startTime := time.Now()
				execCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
				output, exitCode, execErr := executeCommand(execCtx, capturedJob["command"])
				cancel()
				duration := time.Since(startTime)

				record := JobHistoryRecord{
					ExecutedAt: startTime.Format(time.RFC3339),
					ExitCode:   exitCode,
					Output:     output,
					DurationMs: duration.Milliseconds(),
				}
				switch {
				case execErr == nil:
					record.Status = "success"
					util.Logger.Info("standalone job success",
						slog.String("id", capturedID),
						slog.Int64("duration_ms", duration.Milliseconds()))
				case execCtx.Err() == context.DeadlineExceeded:
					record.Status = "timeout"
					record.Error = fmt.Sprintf("exceeded timeout of %s", defaultTimeout)
					util.Logger.Warn("standalone job timeout", slog.String("id", capturedID))
				default:
					record.Status = "failed"
					record.Error = execErr.Error()
					util.Logger.Warn("standalone job failed",
						slog.String("id", capturedID), slog.Int("exit_code", exitCode))
				}
				// ใช้ context.Background() สำหรับ save เพราะ ctx อาจถูก cancel แล้ว
				// แต่เราต้องการ save result ก่อนออก
				saveCtx := context.Background()
				saveJobHistory(saveCtx, capturedID, record)

				nextRun, err := util.GetNextRun(capturedJob["schedule"])
				if err != nil {
					util.RDB.ZRem(saveCtx, util.ProcessingKey(), capturedID)
					return
				}
				pipe := util.RDB.Pipeline()
				pipe.ZRem(saveCtx, util.ProcessingKey(), capturedID)
				pipe.ZAdd(saveCtx, util.ZsetKey(), redis.Z{Score: float64(nextRun), Member: capturedID})
				if _, err := pipe.Exec(saveCtx); err != nil {
					util.Logger.Error("reschedule standalone job failed",
						slog.String("id", capturedID), slog.Any("error", err))
				}
				_ = capturedScore
			}()
		}
	}
}

// ── Workflow processing ───────────────────────────────────────────────────────

// runStep execute หนึ่ง step และ return StepRunState
func runStep(ctx context.Context, step util.WorkflowStep, job map[string]string) util.StepRunState {
	timeout := timeoutForStep(step)
	execCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	startTime := time.Now()
	output, exitCode, execErr := executeCommand(execCtx, job["command"])
	duration := time.Since(startTime)

	state := util.StepRunState{
		StartedAt:  startTime.Format(time.RFC3339),
		FinishedAt: time.Now().Format(time.RFC3339),
		ExitCode:   exitCode,
		Output:     output,
		DurationMs: duration.Milliseconds(),
	}
	switch {
	case execErr == nil:
		state.Status = util.StepStatusSuccess
	case execCtx.Err() == context.DeadlineExceeded:
		state.Status = util.StepStatusTimeout
		state.Error = fmt.Sprintf("exceeded timeout of %s", timeout)
	default:
		state.Status = util.StepStatusFailed
		state.Error = execErr.Error()
	}
	return state
}

// stepResultMsg ใช้ส่งผลจาก goroutine กลับมาที่ orchestrator loop ผ่าน channel
type stepResultMsg struct {
	stepID string
	result util.StepRunState
}

// executeWorkflowRun รัน workflow run ตาม DAG พร้อม concurrent fan-out
//
// แนวคิด:
//   - แต่ละรอบ loop ดู readySteps ทั้งหมดแล้วยิง goroutine พร้อมกัน (fan-out)
//   - goroutine ส่งผลกลับผ่าน resultCh (fan-in)
//   - orchestrator loop อ่าน result ทีละตัว update state แล้ววน loop ใหม่
//   - state ทั้งหมดอยู่ใน goroutine เดียว (orchestrator) → ไม่ต้องใช้ mutex
//   - inFlight นับจำนวน goroutine ที่ยังรันอยู่
func executeWorkflowRun(ctx context.Context, wf *util.Workflow, runID string) {
	log := util.Logger.With(slog.String("workflow_id", wf.ID), slog.String("run_id", runID))

	// สร้าง initial run state
	steps := make(map[string]util.StepRunState, len(wf.Steps))
	for _, s := range wf.Steps {
		steps[s.ID] = util.StepRunState{Status: util.StepStatusPending}
	}
	state := &util.WorkflowRunState{
		RunID:      runID,
		WorkflowID: wf.ID,
		Status:     util.WorkflowStatusRunning,
		Steps:      steps,
		StartedAt:  time.Now().Format(time.RFC3339),
	}
	if err := util.SaveWorkflowRun(ctx, runID, wf.ID, state); err != nil {
		log.Error("save initial run state failed", slog.Any("error", err))
		return
	}
	log.Info("workflow run started", slog.String("name", wf.Name))

	// build step lookup map
	stepMap := make(map[string]util.WorkflowStep, len(wf.Steps))
	for _, s := range wf.Steps {
		stepMap[s.ID] = s
	}

	// channel รับผลจาก goroutine — buffer = จำนวน steps ทั้งหมด (ไม่บล็อก goroutine)
	resultCh := make(chan stepResultMsg, len(wf.Steps))
	inFlight := 0
	stopFlag := false

	// saveState บันทึก state กลับ Redis โดยใช้ context.Background()
	// เพื่อให้ save ได้แม้ว่า ctx ถูก cancel (worker กำลัง shutdown)
	saveState := func() {
		saveCtx := context.Background()
		_ = util.SaveWorkflowRun(saveCtx, runID, wf.ID, state)
	}

	// launchStep ยิง goroutine สำหรับ step เดียว
	launchStep := func(stepID string) {
		step := stepMap[stepID]

		// โหลด job data (ทำใน orchestrator เพื่อ handle error ง่าย)
		job, err := util.RDB.HGetAll(ctx, util.JobKey(step.JobID)).Result()
		if err != nil || len(job) == 0 {
			log.Error("job data missing for step",
				slog.String("step_id", stepID), slog.String("job_id", step.JobID))
			resultCh <- stepResultMsg{
				stepID: stepID,
				result: util.StepRunState{
					Status:     util.StepStatusFailed,
					StartedAt:  time.Now().Format(time.RFC3339),
					FinishedAt: time.Now().Format(time.RFC3339),
					Error:      fmt.Sprintf("job_id %q not found in Redis", step.JobID),
				},
			}
			inFlight++
			return
		}

		// mark running ก่อน launch (orchestrator thread → thread-safe)
		s := state.Steps[stepID]
		s.Status = util.StepStatusRunning
		s.StartedAt = time.Now().Format(time.RFC3339)
		state.Steps[stepID] = s

		inFlight++
		log.Info("launching step",
			slog.String("step_id", stepID),
			slog.String("job_id", step.JobID),
			slog.String("command", job["command"]),
		)

		capturedStep := step
		capturedJob := job
		go func() {
			result := runStep(ctx, capturedStep, capturedJob)
			resultCh <- stepResultMsg{stepID: capturedStep.ID, result: result}
		}()
	}

	// ── Orchestrator loop ──────────────────────────────────────────────────
	for {
		if ctx.Err() != nil {
			log.Warn("workflow run interrupted by context cancellation")
			// drain goroutines ที่ยังรันอยู่ก่อน finalize
			for inFlight > 0 {
				msg := <-resultCh
				inFlight--
				state.Steps[msg.stepID] = msg.result
			}
			break
		}

		// 1. Launch ready steps (fan-out wave)
		if !stopFlag {
			readySteps := util.GetReadySteps(wf, state)
			for _, stepID := range readySteps {
				launchStep(stepID)
			}
			if len(readySteps) > 0 {
				saveState()
			}
		}

		// 2. ถ้าไม่มี goroutine รันอยู่ → ตัดสินใจ
		if inFlight == 0 {
			if util.IsWorkflowDone(state) {
				break
			}
			// ทางตัน: pending อยู่แต่ไม่มี ready step
			log.Warn("workflow blocked, skipping remaining pending steps")
			for stepID, s := range state.Steps {
				if s.Status == util.StepStatusPending {
					s.Status = util.StepStatusSkipped
					s.Error = "skipped: all dependencies failed or workflow stopped"
					state.Steps[stepID] = s
				}
			}
			break
		}

		// 3. รอรับผลจาก goroutine (fan-in)
		msg := <-resultCh
		inFlight--
		state.Steps[msg.stepID] = msg.result

		log.Info("step finished",
			slog.String("step_id", msg.stepID),
			slog.String("status", string(msg.result.Status)),
			slog.Int64("duration_ms", msg.result.DurationMs),
		)
		saveState()

		// 4. Handle failure policy
		stepFailed := msg.result.Status == util.StepStatusFailed ||
			msg.result.Status == util.StepStatusTimeout

		if stepFailed {
			if wf.OnFailure == util.OnFailureStop {
				stopFlag = true
				markDownstreamSkipped(wf, state, msg.stepID)
				// drain goroutines ที่รันอยู่ก่อน break (เก็บ result จริงไว้)
				for inFlight > 0 {
					drained := <-resultCh
					inFlight--
					state.Steps[drained.stepID] = drained.result
					log.Info("step finished (after stop)",
						slog.String("step_id", drained.stepID),
						slog.String("status", string(drained.result.Status)),
					)
				}
				for stepID, s := range state.Steps {
					if s.Status == util.StepStatusPending {
						s.Status = util.StepStatusSkipped
						s.Error = "skipped: upstream step failed with on_failure=stop"
						state.Steps[stepID] = s
					}
				}
				saveState()
				break
			}
			// on_failure=continue: skip เฉพาะ downstream ของ step นี้
			markDownstreamSkipped(wf, state, msg.stepID)
		}
	}

	// ── Finalize ───────────────────────────────────────────────────────────
	state.Status = util.ComputeWorkflowStatus(state)
	state.FinishedAt = time.Now().Format(time.RFC3339)
	saveState()

	failed, skipped := 0, 0
	for _, s := range state.Steps {
		switch s.Status {
		case util.StepStatusFailed, util.StepStatusTimeout:
			failed++
		case util.StepStatusSkipped:
			skipped++
		}
	}

	util.AppendWorkflowHistory(context.Background(), wf.ID, util.WorkflowHistoryRecord{
		RunID:      runID,
		Status:     state.Status,
		StartedAt:  state.StartedAt,
		FinishedAt: state.FinishedAt,
		StepCount:  len(wf.Steps),
		Failed:     failed,
		Skipped:    skipped,
	})

	log.Info("workflow run finished",
		slog.String("status", string(state.Status)),
		slog.Int("failed_steps", failed),
		slog.Int("skipped_steps", skipped),
	)
}

// markDownstreamSkipped mark step ที่ downstream ของ failedStepID ว่า skipped (transitive)
func markDownstreamSkipped(wf *util.Workflow, state *util.WorkflowRunState, failedStepID string) {
	if failedStepID == "" {
		return
	}
	affected := map[string]bool{failedStepID: true}
	changed := true
	for changed {
		changed = false
		for _, s := range wf.Steps {
			if affected[s.ID] {
				continue
			}
			for _, dep := range s.DependsOn {
				if affected[dep] {
					affected[s.ID] = true
					changed = true
					break
				}
			}
		}
	}
	for stepID, s := range state.Steps {
		if affected[stepID] && stepID != failedStepID && s.Status == util.StepStatusPending {
			s.Status = util.StepStatusSkipped
			s.Error = fmt.Sprintf("skipped: upstream step %q failed", failedStepID)
			state.Steps[stepID] = s
		}
	}
}

// processWorkflows poll และรัน workflow runs
// wg ใช้สำหรับ graceful shutdown
func processWorkflows(ctx context.Context, wg *sync.WaitGroup) {
	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	// semaphore จำกัด concurrent workflow runs
	sem := make(chan struct{}, maxConcurrentWorkflows)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			now := time.Now().Unix()
			res, err := util.RDB.ZPopMin(ctx, util.WorkflowZsetKey(), 1).Result()
			if err != nil && err != redis.Nil {
				util.Logger.Error("zpopmin (workflows) failed", slog.Any("error", err))
				continue
			}
			if len(res) == 0 {
				continue
			}
			z := res[0]
			if int64(z.Score) > now {
				util.RDB.ZAdd(ctx, util.WorkflowZsetKey(), redis.Z{Score: z.Score, Member: z.Member})
				continue
			}

			wfID := z.Member.(string)
			data, err := util.RDB.HGetAll(ctx, util.WorkflowKey(wfID)).Result()
			if err != nil || len(data) == 0 {
				util.Logger.Warn("workflow data missing", slog.String("id", wfID))
				continue
			}
			wf, err := util.UnmarshalWorkflow(data)
			if err != nil {
				util.Logger.Error("unmarshal workflow failed",
					slog.String("id", wfID), slog.Any("error", err))
				continue
			}

			// paused → requeue ด้วย nextRun ปกติ
			if wf.Paused {
				nextRun, err := util.GetNextRun(wf.Schedule)
				if err != nil {
					util.Logger.Error("invalid workflow schedule",
						slog.String("id", wfID), slog.Any("error", err))
					continue
				}
				util.RDB.ZAdd(ctx, util.WorkflowZsetKey(),
					redis.Z{Score: float64(nextRun), Member: wfID})
				continue
			}

			// ตรวจ semaphore — ถ้าเต็มให้ push กลับ queue แล้วรอรอบถัดไป
			select {
			case sem <- struct{}{}:
			default:
				util.Logger.Warn("workflow semaphore full, requeueing",
					slog.String("id", wfID))
				util.RDB.ZAdd(ctx, util.WorkflowZsetKey(), redis.Z{
					Score:  float64(time.Now().Add(pollInterval).Unix()),
					Member: wfID,
				})
				continue
			}

			runID := generateRunID(wfID)

			// mark workflow เข้า processing
			util.RDB.ZAdd(ctx, util.WorkflowRunZsetKey(), redis.Z{
				Score:  float64(time.Now().Unix()),
				Member: runID,
			})

			wg.Add(1)
			capturedWF := wf
			capturedRunID := runID
			go func() {
				defer wg.Done()
				defer func() { <-sem }()
				defer util.RDB.ZRem(context.Background(), util.WorkflowRunZsetKey(), capturedRunID)

				executeWorkflowRun(ctx, capturedWF, capturedRunID)

				// reschedule ด้วย context.Background() เพราะ ctx อาจ cancel แล้ว
				nextRun, err := util.GetNextRun(capturedWF.Schedule)
				if err != nil {
					util.Logger.Error("cannot reschedule workflow",
						slog.String("id", capturedWF.ID), slog.Any("error", err))
					return
				}
				if err := util.RDB.ZAdd(context.Background(), util.WorkflowZsetKey(),
					redis.Z{Score: float64(nextRun), Member: capturedWF.ID}).Err(); err != nil {
					util.Logger.Error("reschedule workflow failed",
						slog.String("id", capturedWF.ID), slog.Any("error", err))
				}
			}()
		}
	}
}

// ── main ──────────────────────────────────────────────────────────────────────

func main() {
	util.Init()

	ctx, cancel := context.WithCancel(context.Background())

	// wg ติดตาม goroutines ที่กำลังรัน job/workflow อยู่
	// graceful shutdown จะรอ wg.Wait() ก่อนออก
	var wg sync.WaitGroup

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	// Graceful shutdown goroutine
	go func() {
		sig := <-signals
		util.Logger.Info("received stop signal, stopping new jobs...",
			slog.String("signal", sig.String()))

		// บอกว่าไม่รับ job ใหม่แล้ว
		cancel()

		// รอ job ที่รันอยู่เสร็จ ภายใน shutdownTimeout
		done := make(chan struct{})
		go func() {
			wg.Wait()
			close(done)
		}()

		select {
		case <-done:
			util.Logger.Info("all jobs completed, worker stopped cleanly")
		case <-time.After(shutdownTimeout):
			util.Logger.Warn("shutdown timeout reached, forcing exit",
				slog.String("timeout", shutdownTimeout.String()))
		}
	}()

	util.Logger.Info("worker started",
		slog.String("poll_interval", pollInterval.String()),
		slog.String("default_timeout", defaultTimeout.String()),
		slog.String("stuck_threshold", stuckThreshold.String()),
		slog.Int("max_concurrent_workflows", maxConcurrentWorkflows),
		slog.Int("max_concurrent_standalone_jobs", maxConcurrentStandaloneJobs),
		slog.String("shutdown_timeout", shutdownTimeout.String()),
	)

	// รัน background goroutines
	go recoverStuckJobs(ctx)
	go processStandaloneJobs(ctx, &wg)

	// รัน workflow processor (blocking จนกว่า ctx จะ cancel)
	processWorkflows(ctx, &wg)

	// รอ goroutines ทั้งหมดเสร็จ (กรณีที่ shutdown goroutine ยังไม่ถึง timeout)
	wg.Wait()
	util.Logger.Info("worker stopped")
}
