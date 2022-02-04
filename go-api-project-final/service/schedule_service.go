package service

import "github.com/mandarinkb/go-api-project-final/repository"

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
		schedule = append(schedule, mapDataScheduleResponse(row))
	}
	return schedule, nil
}

func (s scheduleService) ReadById(id int) (*Schedule, error) {
	scheduleRepo, err := s.scheduleRepo.ReadById(id)
	if err != nil {
		return nil, err
	}
	scheduleRes := mapDataScheduleResponse(*scheduleRepo)
	return &scheduleRes, nil
}

func (s scheduleService) Create(schedule Schedule) error {
	return s.scheduleRepo.Create(mapDataScheduleRequest(schedule))
}

func (s scheduleService) Update(schedule Schedule) error {
	return s.scheduleRepo.Update(mapDataScheduleRequest(schedule))
}

func (s scheduleService) Delete(id int) (string, error) {
	scheduleRepo, err := s.scheduleRepo.ReadById(id)
	if err != nil {
		return "", err
	}
	scheduleRes := mapDataScheduleResponse(*scheduleRepo)
	return scheduleRes.ScheduleName, s.scheduleRepo.Delete(id)
}

// แปลงค่า เพื่อส่งไปยัง repository
func mapDataScheduleRequest(schedule Schedule) repository.Schedule {
	return repository.Schedule{
		ScheduleId:     schedule.ScheduleId,
		ScheduleName:   schedule.ScheduleName,
		CronExpression: schedule.CronExpression,
		MethodName:     schedule.MethodName,
		ProjectName:    schedule.ProjectName,
	}
}

// แปลงค่า เพื่อส่งไปยัง handler
func mapDataScheduleResponse(scheduleRepo repository.Schedule) Schedule {
	return Schedule{
		ScheduleId:     scheduleRepo.ScheduleId,
		ScheduleName:   scheduleRepo.ScheduleName,
		CronExpression: scheduleRepo.CronExpression,
		MethodName:     scheduleRepo.MethodName,
		ProjectName:    scheduleRepo.ProjectName,
	}
}
