package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/mandarinkb/go-start-bot-project-final/database"
	"github.com/mandarinkb/go-start-bot-project-final/repository"
)

type webService struct {
	webServ repository.WebRepository
}

func NewWebService(webServ repository.WebRepository) WebService {
	return webService{webServ}
}

func (w webService) Read() error {
	webRepo, err := w.webServ.Read()
	if err != nil {
		return err
	}

	redis := database.RedisConn()
	defer redis.Close()

	for _, row := range webRepo {
		webStr, err := json.Marshal(mapDataWeb(row))
		if err != nil {
			return err
		}
		redis.RPush(context.Background(), "startUrl", string(webStr))
	}
	fmt.Println(time.Now(), "start bot stop")
	return nil
}

// แปลงค่า เพื่อส่งไปยัง handler
func mapDataWeb(webRepo repository.Web) Web {
	return Web{
		WebId:     webRepo.WebId,
		WebName:   webRepo.WebName,
		WebUrl:    webRepo.WebUrl,
		WebStatus: webRepo.WebStatus,
		IconUrl:   webRepo.IconUrl,
	}
}
