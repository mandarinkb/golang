package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mandarinkb/go-cronjob-with-sqlite3/service"
)

type cronJobHandler struct {
	cronJobServ service.CronJobService
}

func NewCronJobHandler(cronJobServ service.CronJobService) *cronJobHandler {
	return &cronJobHandler{cronJobServ}
}

func (h *cronJobHandler) GetCronJob(c *fiber.Ctx) error {
	cronJobData, err := h.cronJobServ.GetCronJob()
	if err != nil {
		return c.Status(400).SendString(err.Error())
	}
	return c.Status(200).JSON(cronJobData)
}
