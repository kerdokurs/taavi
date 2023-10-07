package handlers

import (
	"github.com/gofiber/fiber/v2"
	"kerdo.dev/taavi/logger"
	"kerdo.dev/taavi/scheduler"
)

func HandleScheduledGet(c *fiber.Ctx) error {
	entries := scheduler.Scheduler.Entries()

	return c.Status(fiber.StatusOK).JSON(entries)
}

func HandleReschedulePost(c *fiber.Ctx) error {
	type request struct {
		IgnoreMaster bool `json:"ignore_master" form:"ignore_master"`
	}

	var req request
	if err := c.BodyParser(&req); err != nil {
		logger.Errorw("error decoding request", logger.M{
			"err": err.Error(),
		})
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	scheduler.RescheduleAll(req.IgnoreMaster)
	return c.SendStatus(200)
}
