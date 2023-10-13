package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/robfig/cron/v3"
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

func HandleScheduledRun(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
	go scheduler.Scheduler.Entry(cron.EntryID(id)).Job.Run()
	return c.SendStatus(fiber.StatusOK)
}
