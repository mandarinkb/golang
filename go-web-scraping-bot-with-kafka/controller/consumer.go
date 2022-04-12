package controller

import (
	"encoding/json"

	"github.com/Shopify/sarama"
	"github.com/mandarinkb/go-web-scraping-bot-with-kafka/database"
	"github.com/mandarinkb/go-web-scraping-bot-with-kafka/repository"
	"github.com/mandarinkb/go-web-scraping-bot-with-kafka/service"
	"github.com/mandarinkb/go-web-scraping-bot-with-kafka/utils"
)

type ConsumerGroupHandler struct{}

func (ConsumerGroupHandler) Setup(sarama.ConsumerGroupSession) error   { return nil }
func (ConsumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error { return nil }
func (ConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	logger, err := utils.LogConf()
	if err != nil {
		logger.Error(err.Error(), utils.Url("-"),
			utils.User("-"), utils.Type(utils.TypeBot))
	}
	// close log
	defer logger.Sync()

	db, err := database.Conn()
	if err != nil {
		logger.Error(err.Error(), utils.Url("-"),
			utils.User("-"), utils.Type(utils.TypeBot))
	}
	defer db.Close()

	webData := service.Web{}
	for msg := range claim.Messages() {
		switch string(msg.Key) {
		case "detail-url":
			err := json.Unmarshal((msg.Value), &webData)
			if err != nil {
				logger.Error(err.Error(), utils.Url("-"),
					utils.User("-"), utils.Type(utils.TypeBot))
			}
			switch webData.WebName {
			case "tescolotus":
				swRepo := repository.NewSwitchDatabaseDB(db)
				service.NewProductService(swRepo).Tescolotus(webData)
			case "makroclick":
				service.Makroclick(webData)
			case "bigc":
				service.Bigc(webData)
			}
		}
		session.MarkMessage(msg, "")
	}
	return nil
}
