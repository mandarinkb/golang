package service

import "context"

type CronJobService interface {
	RunJobService(ctx context.Context) error
}
