package service

import "github.com/mandarinkb/go-start-bot-with-kafka/repository"

type scheduleService struct {
	scheduleRepo repository.ScheduleRepository
}

func NewScheduleService(scheRepo repository.ScheduleRepository) ScheduleService {
	return scheduleService{scheRepo}
}

func (s scheduleService) Read() (schedule []Schedule, err error) {
	scheduleRepo, err := s.scheduleRepo.Read()
	if err != nil {
		return nil, err
	}
	for _, row := range scheduleRepo {
		schedule = append(schedule, mapDataSchedule(row))
	}
	return schedule, nil
}

func mapDataSchedule(scheduleRepo repository.Schedule) Schedule {
	return Schedule{
		ScheduleId:     scheduleRepo.ScheduleId,
		ScheduleName:   scheduleRepo.ScheduleName,
		CronExpression: scheduleRepo.CronExpression,
		MethodName:     scheduleRepo.MethodName,
		ProjectName:    scheduleRepo.ProjectName,
	}
}
