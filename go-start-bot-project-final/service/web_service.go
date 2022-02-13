package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/mandarinkb/go-start-bot-project-final/database"
	"github.com/mandarinkb/go-start-bot-project-final/repository"
	"github.com/mandarinkb/go-start-bot-project-final/utils"
)

type webService struct {
	webServ repository.WebRepository
}

func NewWebService(webServ repository.WebRepository) WebService {
	return webService{webServ}
}

func (w webService) Read() error {
	logger, err := utils.LogConf()
	if err != nil {
		logger.Error(err.Error(), utils.Url("-"),
			utils.User("-"), utils.Type(utils.TypeBot))
	}
	// close log
	defer logger.Sync()

	logger.Info("[start bot] start", utils.Url("-"),
		utils.User("-"), utils.Type(utils.TypeBot))
	fmt.Println(time.Now(), " : start bot start")

	err = clearOldData()
	if err != nil {
		logger.Error(err.Error(), utils.Url("-"),
			utils.User("-"), utils.Type(utils.TypeBot))
		return err
	}

	redis := database.RedisConn()
	defer redis.Close()

	webRepo, err := w.webServ.Read()
	if err != nil {
		logger.Error(err.Error(), utils.Url("-"),
			utils.User("-"), utils.Type(utils.TypeBot))
		return err
	}

	for _, row := range webRepo {
		webStr, err := json.Marshal(mapDataWeb(row))
		if err != nil {
			logger.Error(err.Error(), utils.Url("-"),
				utils.User("-"), utils.Type(utils.TypeBot))
			return err
		}
		redis.RPush(context.Background(), "startUrl", string(webStr))
	}
	fmt.Println(time.Now(), " : start bot stop")

	logger.Info("[start bot] stop", utils.Url("-"),
		utils.User("-"), utils.Type(utils.TypeBot))
	return nil
}
func clearOldData() error {
	logger, err := utils.LogConf()
	if err != nil {
		logger.Error(err.Error(), utils.Url("-"),
			utils.User("-"), utils.Type(utils.TypeBot))
	}
	// close log
	defer logger.Sync()

	redis := database.RedisConn()
	defer redis.Close()

	// delete data in redis
	redis.Del(context.Background(), "startUrl")
	redis.Del(context.Background(), "fetchUrl")
	redis.Del(context.Background(), "detailUrl")
	logger.Info("clear data in redis", utils.Url("-"),
		utils.User("-"), utils.Type(utils.TypeBot))

	db, err := database.Conn()
	if err != nil {
		return err
	}
	defer db.Close()

	// swap database
	err = repository.NewSwitchDBRepo(db).SwapDatabase()
	if err != nil {
		return err
	}
	logger.Info("swap database in elasticsearch", utils.Url("-"),
		utils.User("-"), utils.Type(utils.TypeBot))

	// เลือก database ในสถานะที่เป็น 0
	InActavateDb, err := repository.NewSwitchDBRepo(db).ReadInActivateSwitchDatabase()
	if err != nil {
		return err
	}

	// ลบ index ใน elasticsearch
	err = repository.DeleteIndex(InActavateDb.DatabaseName)
	if err != nil {
		return err
	}
	logger.Info("clear data in "+InActavateDb.DatabaseName, utils.Url("-"),
		utils.User("-"), utils.Type(utils.TypeBot))

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
