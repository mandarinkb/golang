package service

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/Shopify/sarama"
	"github.com/mandarinkb/go-start-bot-with-kafka/database"
	"github.com/mandarinkb/go-start-bot-with-kafka/repository"
	"github.com/mandarinkb/go-start-bot-with-kafka/utils"
)

var (
	topic string = "start-url"
	key   string = "start-url"
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

	producer := database.KafkaConn()
	defer producer.Close()

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
		message := &sarama.ProducerMessage{
			Topic: topic,
			Key:   sarama.StringEncoder(key),
			Value: sarama.StringEncoder(string(webStr)),
		}

		pid, offset, err := producer.SendMessage(message)
		if err != nil {
			return err
		}
		fmt.Println("partition : ", pid)
		fmt.Println("offset : ", offset)
		fmt.Println(string(webStr))
		fmt.Println()
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

	clusterAdmin := database.KafkaAdminConn()
	defer clusterAdmin.Close()
	clusterAdmin.DeleteTopic("start-url")
	clusterAdmin.DeleteTopic("fetch-url")
	clusterAdmin.DeleteTopic("detail-url")

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
