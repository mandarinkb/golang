package controller

import (
	"encoding/json"

	"github.com/Shopify/sarama"
	"github.com/mandarinkb/go-fetch-url-bot-with-kafka/service"
	"github.com/mandarinkb/go-fetch-url-bot-with-kafka/utils"
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

	webData := service.Web{}
	for msg := range claim.Messages() {
		switch string(msg.Key) {
		// จัดการหมวดหมู่สินค้าใหม่
		case "start-url":
			err := json.Unmarshal((msg.Value), &webData)
			if err != nil {
				logger.Error(err.Error(), utils.Url("-"),
					utils.User("-"), utils.Type(utils.TypeBot))
			}
			switch webData.WebName {
			case "tescolotus":
				service.TescolotusMainPage(webData)
			case "makroclick":
				service.MakroMainPage(webData)
			case "bigc":
				service.BigcMainPage(webData)
			}
			// หา url ของสินค้าในแต่ละหมวดหมู่ โดยจะหาทุกๆหน้า
		case "fetch-url":
			err := json.Unmarshal((msg.Value), &webData)
			if err != nil {
				logger.Error(err.Error(), utils.Url("-"),
					utils.User("-"), utils.Type(utils.TypeBot))
			}
			switch webData.WebName {
			case "tescolotus":
				service.TescolotusFindUrlInPage(webData)
			case "makroclick":
				service.MakroFindUrlInPage(webData)
			case "bigc":
				service.BigcFindUrlInPage(webData)
			}
		}
		session.MarkMessage(msg, "")
	}
	return nil
}
