package handlers

import (
	"github.com/gofiber/fiber/v2"
	"kerdo.dev/taavi/scheduler"
)

func HandleScheduledGet(c *fiber.Ctx) error {
	entries := scheduler.Scheduler.Entries()
	return c.Status(fiber.StatusOK).JSON(entries)
}

func HandleReschedulePost(c *fiber.Ctx) error {
	scheduler.RescheduleAll()
	return c.SendStatus(fiber.StatusOK)
}
