package scheduler

import (
	"github.com/gofiber/fiber/v2"
)

func HandleScheduledGet(ctx *fiber.Ctx) error {
	entries := Scheduler.Entries()

	return ctx.Status(fiber.StatusOK).JSON(entries)
}
