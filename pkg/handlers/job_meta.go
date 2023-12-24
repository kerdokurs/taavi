package handlers

import (
	"github.com/gofiber/fiber/v2"
	"kerdo.dev/taavi/pkg/data"
	"kerdo.dev/taavi/pkg/logger"
)

func HandleJobMetaDelete(c *fiber.Ctx) error {
	jobID, err := c.ParamsInt("id", -1)
	if err != nil {
		logger.Errorw("error parsing int param", logger.M{
			"err": err.Error(),
		})
		return c.SendStatus(fiber.StatusBadRequest)
	}
	metaID, err := c.ParamsInt("metaId", -1)
	if err != nil {
		logger.Errorw("error parsing int param", logger.M{
			"err": err.Error(),
		})
		return c.SendStatus(fiber.StatusBadRequest)
	}

	ctx := c.Context()
	if err := data.DeleteJobMeta(ctx, uint(jobID), uint(metaID)); err != nil {
		logger.Errorw("error deleting job meta", logger.M{
			"err": err.Error(),
		})
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	c.Response().Header.Set("Content-Type", "text/html; charset=utf-8")
	c.Status(fiber.StatusOK)
	_, err = c.Response().BodyWriter().Write([]byte(""))
	return err
}
