# Go Cron Job with Redis — Architecture Overview

## ภาพรวมระบบ

ระบบเป็น **Distributed Cron Job Manager** ที่รองรับทั้ง standalone jobs และ workflow (DAG)
ใช้ Redis เป็น job store, scheduler, และ state store กลาง แบ่งเป็น 2 process ที่ deploy แยกกัน

```
┌──────────────────────────────┐      ┌──────────────────────────────────────┐
│         API Server           │      │               Worker                 │
│       (main.go :3000)        │      │          (worker/main.go)            │
│                              │      │                                      │
│  /cron      — standalone     │      │  processStandaloneJobs  (goroutine)  │
│  /workflow  — workflow DAG   │      │    semaphore: 20 concurrent          │
│  /workflow/:id/runs/:run_id  │      │  processWorkflows       (main loop)  │
│                              │      │    semaphore: 10 concurrent          │
│                              │      │  recoverStuckJobs       (background) │
│                              │      │  graceful shutdown via sync.WaitGroup│
└──────────────────────────────┘      └──────────────────────────────────────┘
              │                                         │
              └──────────────────┬──────────────────────┘
                                 ▼
                         ┌──────────────┐
                         │    Redis     │
                         │  pool: 20   │
                         │  retry: 3x  │
                         │              │
                         │  job:{id}    │ Hash — job metadata
                         │  cron_jobs   │ ZSet — score=nextRun
                         │  cron_proc.  │ ZSet — score=startedAt
                         │  job_history │ List — max 100
                         │  :{id}       │
                         │              │
                         │  workflow:   │ Hash — workflow definition
                         │  {id}        │
                         │  workflow_   │ ZSet — score=nextRun
                         │  runs        │
                         │  workflow_   │ ZSet — score=startedAt
                         │  processing  │
                         │  workflow_   │ Hash — step runtime state, TTL 7d
                         │  run:{wf}:   │
                         │  {runID}     │
                         │  workflow_   │ List — max 50
                         │  history:{id}│
                         └──────────────┘
```

---

## โครงสร้าง Project

```
go-cronjob-with-redis/
├── main.go                  # API Server (Fiber)
├── worker/
│   └── main.go              # Worker: standalone + workflow processor
├── util/
│   ├── redis.go             # Redis client, connection pool, Logger, GetEnv
│   ├── cron.go              # GetNextRun() → (int64, error)
│   ├── keys.go              # Redis key naming conventions
│   └── workflow.go          # Types, DAG helpers, serialization
├── script/
│   └── job_script.sh
├── docker-compose.yml       # Redis + API + Worker services
├── Dockerfile               # Multi-stage build golang:1.23
└── go.mod
```

---

## Standalone Job

### Data Model

| Key | Type | Description |
|-----|------|-------------|
| `job:{id}` | Hash | id, name, command, schedule, paused |
| `cron_jobs` | ZSet | priority queue, score = nextRun (Unix ts) |
| `cron_processing` | ZSet | in-flight tracker, score = startedAt |
| `job_history:{id}` | List | execution records, max 100 |

### API Endpoints

| Method | Path | Description |
|--------|------|-------------|
| `POST` | `/cron` | สร้าง job |
| `GET` | `/cron` | list jobs พร้อม next_run |
| `GET` | `/cron/:id` | job detail |
| `PUT` | `/cron/:id` | update job |
| `DELETE` | `/cron/:id` | ลบ job (history ยังอยู่) |
| `PATCH` | `/cron/:id/now` | รัน job ทันที |
| `POST` | `/cron/:id/pause` | หยุด job ชั่วคราว |
| `POST` | `/cron/:id/resume` | เปิด job อีกครั้ง |
| `GET` | `/cron/:id/history` | execution history 100 รายการ |

---

## Workflow (DAG)

### Concept

Workflow คือกลุ่มของ steps ที่มี dependency กัน สร้างเป็น DAG (Directed Acyclic Graph)
แต่ละ step อ้างถึง standalone job และกำหนดว่าต้องรอ step ใดก่อน

```json
{
  "name": "ETL Pipeline",
  "schedule": "0 2 * * *",
  "on_failure": "stop",
  "steps": [
    { "id": "extract",    "job_id": "abc123", "depends_on": [] },
    { "id": "transform",  "job_id": "def456", "depends_on": ["extract"] },
    { "id": "load",       "job_id": "ghi789", "depends_on": ["transform"] }
  ]
}
```

### Fan-out / Fan-in (Concurrent Execution)

Step ที่ dependency ครบพร้อมกันหลายตัวจะถูก **launch พร้อมกันเป็น goroutine** (fan-out)
แล้วรวมผลกลับผ่าน channel (fan-in) ก่อนตัดสินใจ step ถัดไป

```
extract → transform_a ─┐
        ↘ transform_b ─┴→ load
```

```json
"steps": [
  { "id": "extract",     "depends_on": [] },
  { "id": "transform_a", "depends_on": ["extract"] },
  { "id": "transform_b", "depends_on": ["extract"] },
  { "id": "load",        "depends_on": ["transform_a", "transform_b"] }
]
```

`transform_a` และ `transform_b` รันพร้อมกัน, `load` รอทั้งคู่เสร็จ

### Step State Machine

```
pending → running → success
                 ↘ failed
                 ↘ timeout
                 ↘ skipped  ← dependency fail + on_failure=continue
                              หรือ on_failure=stop (all downstream)
```

### on_failure Policy

| Policy | พฤติกรรม |
|--------|----------|
| `stop` (default) | หยุดทันที drain goroutines ที่รันอยู่ แล้ว mark downstream ทั้งหมดว่า skipped |
| `continue` | ข้าม step ที่ fail mark เฉพาะ downstream ของ step นั้นว่า skipped แต่รัน branch อื่นต่อ |

### Workflow Status

| Status | เงื่อนไข |
|--------|---------|
| `success` | ทุก step success |
| `failed` | มี step fail และไม่มี skipped |
| `partial` | มีทั้ง failed และ skipped |

### Data Model

| Key | Type | Description |
|-----|------|-------------|
| `workflow:{id}` | Hash | name, schedule, steps(JSON), on_failure, paused |
| `workflow_runs` | ZSet | priority queue, score = nextRun |
| `workflow_processing` | ZSet | active runs tracker, score = startedAt |
| `workflow_run:{wf_id}:{run_id}` | Hash | runtime step states, TTL 7 วัน |
| `workflow_history:{id}` | List | run summary records, max 50 |

### API Endpoints

| Method | Path | Description |
|--------|------|-------------|
| `POST` | `/workflow` | สร้าง workflow (validate DAG + job existence) |
| `GET` | `/workflow` | list workflows (pipeline HGetAll) |
| `GET` | `/workflow/:id` | workflow detail + next_run |
| `DELETE` | `/workflow/:id` | ลบ workflow (history ยังอยู่) |
| `PATCH` | `/workflow/:id/now` | trigger ทันที |
| `POST` | `/workflow/:id/pause` | หยุด workflow |
| `POST` | `/workflow/:id/resume` | เปิด workflow |
| `GET` | `/workflow/:id/history` | run history 50 รายการ |
| `GET` | `/workflow/:id/runs/:run_id` | step-level detail ของ run |

### ตัวอย่าง Request

```bash
# 1. สร้าง jobs ก่อน
curl -X POST http://localhost:3000/cron \
  -H "Content-Type: application/json" \
  -d '{"name":"extract","schedule":"0 2 * * *","command":"echo extract"}'
# → {"id":"abc123"}

# 2. สร้าง workflow
curl -X POST http://localhost:3000/workflow \
  -H "Content-Type: application/json" \
  -d '{
    "name": "ETL Pipeline",
    "schedule": "0 2 * * *",
    "on_failure": "stop",
    "steps": [
      {"id":"extract",   "job_id":"abc123","depends_on":[]},
      {"id":"transform", "job_id":"def456","depends_on":["extract"]},
      {"id":"load",      "job_id":"ghi789","depends_on":["transform"]}
    ]
  }'

# 3. trigger ทันที
curl -X PATCH http://localhost:3000/workflow/{wf_id}/now

# 4. ดู step-level detail ของ run
curl http://localhost:3000/workflow/{wf_id}/runs/{run_id}
```

---

## Worker Architecture

### Concurrency Model

```
main()
  ├── go recoverStuckJobs(ctx)           background, every 1m
  ├── go processStandaloneJobs(ctx, wg)  goroutine, poll 2s
  │       semaphore: 20 concurrent jobs
  │       each job → go func() { wg.Add/Done }
  │
  └── processWorkflows(ctx, wg)          main loop, poll 2s
          semaphore: 10 concurrent workflows
          each workflow → go func() { wg.Add/Done }
                            executeWorkflowRun()
                              fan-out: launch ready steps concurrently
                              fan-in:  collect via resultCh channel
```

### Standalone Job Flow

```
poll 2s:
  ZPopMin cron_jobs → check score (future? push back)
  → load job metadata → check paused
  → semaphore check (full? requeue)
  → ZAdd cron_processing (mark in-flight)
  → go func:
      exec.CommandContext(timeout=30s)
      saveJobHistory (context.Background)
      ZRem processing + ZAdd nextRun

recoverStuckJobs (every 1m):
  ZRangeByScore cron_processing score < (now - 5m)
  → requeue + saveHistory status="recovered"
```

### Workflow Execution Flow

```
poll 2s:
  ZPopMin workflow_runs → check score → load workflow → check paused
  → semaphore check (full? requeue)
  → ZAdd workflow_processing
  → go func (wg tracked):
      executeWorkflowRun(ctx, wf, runID)
        loop until done:
          GetReadySteps() → launch all as goroutines (fan-out)
          saveState() → รอรับ result จาก resultCh (fan-in)
          handle failure policy:
            stop     → drain goroutines + markDownstreamSkipped + break
            continue → markDownstreamSkipped (branch only)
        ComputeWorkflowStatus() → AppendWorkflowHistory()
      reschedule nextRun (context.Background)
```

### Graceful Shutdown

```
SIGTERM/SIGINT
  → cancel(ctx)           หยุดรับ job/workflow ใหม่
  → wg.Wait()             รอ job/workflow ที่รันอยู่เสร็จ
  → timeout 60s           force exit ถ้ารอนานเกินไป

save operations ใช้ context.Background()
เพื่อให้ save result ได้แม้ ctx ถูก cancel แล้ว
```

---

## Redis Configuration

```go
DialTimeout:     5s
ReadTimeout:     3s
WriteTimeout:    3s
PoolSize:        20    // รองรับ concurrent workflows + jobs
MinIdleConns:    5     // ป้องกัน cold start latency
ConnMaxIdleTime: 5m
MaxRetries:      3     // retry อัตโนมัติเมื่อ network hiccup
MinRetryBackoff: 100ms
MaxRetryBackoff: 1s
```

---

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `REDIS_ADDR` | `localhost:6379` | Redis address |
| `REDIS_PASSWORD` | `mandarinkb` | Redis password |
| `API_PORT` | `3000` | HTTP port |

---

## การรันด้วย Docker Compose

```bash
# รัน Redis + API + Worker
docker compose up --build

# scale worker (ZPopMin atomic รองรับ multi-worker)
docker compose up --scale worker=3
```

---

## Performance Characteristics

| Operation | Redis Calls | หมายเหตุ |
|-----------|------------|---------|
| `GET /cron` | 1 pipeline | ZRangeWithScores + N HGetAll ใน 1 round trip |
| `GET /workflow` | 1 pipeline | ZRangeWithScores + N HGetAll ใน 1 round trip |
| `POST /workflow` | 1 pipeline | N EXISTS ใน 1 round trip สำหรับ validate job IDs |
| Worker poll | 1 ZPopMin | atomic claim |
| Step save progress | 1 HSet | บันทึกหลังแต่ละ step เสร็จ |

---

## Dependencies

| Package | Version | ใช้สำหรับ |
|---------|---------|---------|
| `gofiber/fiber/v2` | v2.52.13 | HTTP API server |
| `redis/go-redis/v9` | v9.20.1 | Redis client |
| `robfig/cron/v3` | v3.0.1 | Parse cron expressions |
| `google/uuid` | v1.6.0 | Generate unique IDs |
