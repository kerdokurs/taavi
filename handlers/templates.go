package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kerdokurs/zlp"
	"kerdo.dev/taavi/data"
	"kerdo.dev/taavi/logger"
	"kerdo.dev/taavi/views"
	"kerdo.dev/taavi/zulip"
)

func HandleIndex(c *fiber.Ctx) error {
	jobs, err := data.GetAllJobs()
	if err != nil {
		logger.Errorw("error getting jobs", logger.M{
			"err": err.Error(),
		})
		jobs = []data.Job{}
	}
	streams, err := zulip.Client.GetStreams()
	if err != nil {
		logger.Errorw("error getting Zulip streams", logger.M{
			"err": err.Error(),
		})
		streams = []zlp.Stream{}
	}

	comp := views.Index(jobs, streams)
	c.Response().Header.Set("Content-Type", "text/html; charset=utf-8")
	return comp.Render(c.Context(), c.Response().BodyWriter())
}
