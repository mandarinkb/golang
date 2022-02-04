package service

import "github.com/mandarinkb/go-api-project-final/repository"

type LogSystem struct {
	Level     string `json:"level"`
	Timestamp string `json:"timestamp"`
	Caller    string `json:"caller"`
	User      string `json:"user"`
	Massage   string `json:"msg"`
	Url       string `json:"url"`
	TypeLog   string `json:"typeLog"`
}
type LogReq struct {
	Date string `json:"date"`
}

type logSystemService struct {
	logRepo repository.LogRepository
}

type LogSystemService interface {
	GetLogs(date LogReq) ([]LogSystem, error)
}

func NewLogSystem(log repository.LogRepository) LogSystemService {
	return logSystemService{log}
}

func (l logSystemService) GetLogs(date LogReq) (logArr []LogSystem, err error) {
	logRepo, err := l.logRepo.GetLogs(date.Date)
	if err != nil {
		return nil, err
	}
	for _, row := range logRepo {
		logArr = append(logArr, mapDataLogRespose(row))
	}
	return logArr, nil
}

func mapDataLogRespose(log repository.LogSystem) LogSystem {
	return LogSystem{Level: log.Level,
		Timestamp: log.Timestamp,
		Caller:    log.Caller,
		User:      log.User,
		Massage:   log.Massage,
		Url:       log.Url,
		TypeLog:   log.TypeLog}
}
