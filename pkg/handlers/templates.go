package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kerdokurs/zlp"
	"kerdo.dev/taavi/pkg/data"
	"kerdo.dev/taavi/pkg/logger"
	"kerdo.dev/taavi/pkg/zulip"
	"kerdo.dev/taavi/views"
)

func HandleIndex(c *fiber.Ctx) error {
	ctx := c.Context()
	jobs, err := data.GetAllJobs(ctx)
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
