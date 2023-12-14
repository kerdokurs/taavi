package handlers

import (
	"github.com/gofiber/fiber/v2"
)

func Init(app *fiber.App) {
	app.Get("/api/jobs", HandleJobsGet)
	app.Post("/api/jobs", HandleJobsPost)
	app.Delete("/api/jobs/:id", HandleJobsDelete)
	app.Post("/api/jobs/:id/toggle", HandleJobEnabledToggle)

	app.Post("/api/jobs/:id/run", HandleJobRun)

	app.Get("/api/scheduled", HandleScheduledGet)
	app.Post("/api/scheduled/:id/run", HandleScheduledRun)
	app.Post("/api/reschedule", HandleReschedulePost)

	app.Get("/", HandleIndex)
}
