package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kerdokurs/zlp"
	"kerdo.dev/taavi/data"
	"kerdo.dev/taavi/logger"
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

	return c.Render("index", map[string]any{
		"jobs":    jobs,
		"streams": streams,
	}, "layouts/main")
}
