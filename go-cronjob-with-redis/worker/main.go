package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"log/slog"

	"github.com/mandarinkb/go-cronjob-with-redis/util"
	"github.com/redis/go-redis/v9"
)

func processCronJobs(ctx context.Context) {
	pollInterval := 2 * time.Second
	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			util.Logger.Info("worker stopping")
			return
		case <-ticker.C:
			now := time.Now().Unix()

			// ZPopMin to claim one job (atomic)
			res, err := util.RDB.ZPopMin(util.CTX, util.ZsetKey(), 1).Result()
			if err != nil && err != redis.Nil {
				util.Logger.Error("zpopmin failed", slog.Any("error", err))
				continue
			}
			if len(res) == 0 {
				continue
			}

			z := res[0]
			// If score is in the future (possible if another worker pushed back), re-add
			if int64(z.Score) > now {
				// push it back
				if err := util.RDB.ZAdd(util.CTX, util.ZsetKey(), redis.Z{Score: z.Score, Member: z.Member}).Err(); err != nil {
					util.Logger.Error("requeue failed", slog.Any("error", err))
				}
				continue
			}

			jobID := z.Member.(string)
			job, err := util.RDB.HGetAll(util.CTX, util.JobKey(jobID)).Result()
			if err != nil || len(job) == 0 {
				util.Logger.Warn("job data missing", slog.String("id", jobID))
				continue
			}

			if job["paused"] == "true" {
				util.Logger.Info("job paused - skipping", slog.String("id", jobID))
				// skip but do not schedule next run -> we might want to schedule in future; for now re-add to zset with same score+60s
				// choose to requeue with +60s to avoid busy-loop
				retryAt := time.Now().Add(60 * time.Second).Unix()
				util.RDB.ZAdd(util.CTX, util.ZsetKey(), redis.Z{Score: float64(retryAt), Member: jobID})
				continue
			}

			util.Logger.Info("executing job", slog.String("id", jobID), slog.String("cmd", job["command"]))

			// TODO: run job command properly (exec.Command with timeout/context)
			// Example placeholder: simulate execution
			time.Sleep(100 * time.Millisecond)

			// compute next run and schedule it
			nextRun := util.GetNextRun(job["schedule"])
			if err := util.RDB.ZAdd(util.CTX, util.ZsetKey(), redis.Z{Score: float64(nextRun), Member: jobID}).Err(); err != nil {
				util.Logger.Error("schedule next failed", slog.Any("error", err), slog.String("id", jobID))
			} else {
				util.Logger.Info("scheduled next run", slog.String("id", jobID),
					slog.String("next", time.Unix(nextRun, 0).Format(time.RFC3339)))
			}
		}
	}
}

func main() {
	util.Init()

	ctx, cancel := context.WithCancel(context.Background())
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-signals
		util.Logger.Info("received stop signal")
		cancel()
	}()

	util.Logger.Info("worker started")
	processCronJobs(ctx)
}
