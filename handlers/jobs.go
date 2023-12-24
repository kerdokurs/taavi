package handlers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"kerdo.dev/taavi/logger"
	"kerdo.dev/taavi/pkg/data"
	"kerdo.dev/taavi/scheduler"
	"kerdo.dev/taavi/views/partials"
)

func HandleJobsGet(c *fiber.Ctx) error {
	isHtmx := c.Locals("htmx").(bool)

	ctx := c.Context()

	jobs, err := data.GetAllJobs(ctx)
	if err != nil {
		logger.Errorw("error getting jobs", logger.M{
			"err": err.Error(),
		})
		jobs = []data.Job{}
	}
	var jobIDs []uint
	for _, job := range jobs {
		jobIDs = append(jobIDs, job.ID)
	}
	metas, err := data.GetAllJobMetas(ctx, jobIDs)
	if err != nil {
		logger.Errorw("error getting job metas", logger.M{
			"err": err.Error(),
		})
		metas = map[uint][]data.JobMeta{}
	}
	for i, job := range jobs {
		jobs[i].Meta = metas[job.ID]
	}

	if !isHtmx {
		return c.Status(fiber.StatusOK).JSON(jobs)
	}

	return c.SendStatus(fiber.StatusUnsupportedMediaType)
}

type NewJobDto struct {
	Type     data.JobType `json:"type" form:"type"`
	StreamID string       `json:"stream_id" form:"stream_id"`
	TopicID  string       `json:"topic_id" form:"topic_id"`
	Content  string       `json:"content" form:"content"`
	CronTime string       `json:"cron_time" form:"cron_time"`
}

func HandleJobsPost(c *fiber.Ctx) error {
	type request struct {
		NewJobDto
		Metas []data.JobMeta
	}

	var req request

	if err := c.BodyParser(&req); err != nil {
		logger.Errorw("error decoding JSON payload", logger.M{
			"err": err.Error(),
		})
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	job := data.Job{
		Type:     req.Type,
		StreamID: req.StreamID,
		TopicID:  req.TopicID,
		Content:  req.Content,
		CronTime: req.CronTime,
		Enabled:  true,
	}
	fmt.Printf("%+v\n", job)
	for _, meta := range req.Metas {
		fmt.Printf("%s -> %s\n", meta.Key, meta.Value)
	}

	// if err := data.DB.Transaction(func(tx *gorm.DB) error {
	// 	if err := tx.Create(&job).Error; err != nil {
	// 		return err
	// 	}

	// 	for _, meta := range req.Metas {
	// 		jobMeta := data.JobMeta{
	// 			Key:   meta.Key,
	// 			Value: meta.Value,
	// 			JobID: int(job.ID),
	// 		}
	// 		if err := tx.Create(&jobMeta).Error; err != nil {
	// 			return err
	// 		}
	// 	}

	// 	return nil
	// }); err != nil {
	// 	logger.Errorw("error creating job", logger.M{
	// 		"err": err.Error(),
	// 	})
	// 	return c.SendStatus(fiber.StatusInternalServerError)
	// }

	if err := data.InsertJobWithMeta(c.Context(), job, req.Metas); err != nil {
		logger.Errorw("error creating job", logger.M{
			"err": err.Error(),
		})
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	c.Response().Header.Set("Content-Type", "text/html; charset=utf-8")
	return partials.JobRow(&job).Render(c.Context(), c.Response().BodyWriter())
}

func HandleJobsDelete(c *fiber.Ctx) error {
	jobID, err := c.ParamsInt("id")
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	ctx := c.Context()

	if err := data.DeleteJob(ctx, uint(jobID)); err != nil {
		logger.Errorw("error deleting job", logger.M{
			"err": err.Error(),
		})
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	isHtmx := c.Locals("htmx").(bool)

	if !isHtmx {
		return c.SendStatus(fiber.StatusOK)
	}

	return c.Status(fiber.StatusOK).SendString("")
}

func HandleJobRun(c *fiber.Ctx) error {
	jobID, err := c.ParamsInt("id", -1)
	if err != nil {
		logger.Errorw("error parsing id parameter", logger.M{
			"err": err.Error(),
		})
		return c.SendStatus(fiber.StatusBadRequest)
	}

	ctx := c.Context()

	job, err := data.GetJob(ctx, uint(jobID))
	if err != nil {
		logger.Errorw("error getting job", logger.M{
			"err": err.Error(),
			"id":  jobID,
		})
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	cronJob, _, err := scheduler.CreateJob(job)
	if err != nil {
		logger.Errorw("error creating job", logger.M{
			"err": err.Error(),
		})
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	go cronJob.Run()

	return c.Status(fiber.StatusOK).SendString("set to run")
}

func HandleJobEnabledToggle(c *fiber.Ctx) error {
	jobID, err := c.ParamsInt("id", -1)
	if err != nil {
		logger.Errorw("error parsing int param", logger.M{
			"err": err.Error(),
		})
		return c.SendStatus(fiber.StatusBadRequest)
	}

	ctx := c.Context()
	newState, err := data.ToggleJobEnabled(ctx, uint(jobID))
	if err != nil {
		logger.Errorw("error toggling job enabled state", logger.M{
			"err": err.Error(),
		})
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	isHtmx := c.Locals("htmx").(bool)

	if !isHtmx {
		return c.Status(fiber.StatusOK).JSON(newState)
	}

	job := data.Job{
		Model: data.Model{
			ID: uint(jobID),
		},
		Enabled: newState,
	}
	return partials.JobCheckbox(&job).Render(c.Context(), c.Response().BodyWriter())
}

func HandleJobMetaGet(c *fiber.Ctx) error {
	jobID, err := c.ParamsInt("id", -1)
	if err != nil {
		logger.Errorw("error parsing int param", logger.M{
			"err": err.Error(),
		})
		return c.SendStatus(fiber.StatusBadRequest)
	}

	ctx := c.Context()

	metas, err := data.GetJobMetas(ctx, uint(jobID))
	if err != nil {
		logger.Errorw("error getting job metas", logger.M{
			"err": err.Error(),
		})
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	isHtmx := c.Locals("htmx").(bool)

	if !isHtmx {
		return c.Status(fiber.StatusOK).JSON(metas)
	}

	return partials.JobMetaList(uint(jobID), metas).Render(c.Context(), c.Response().BodyWriter())
}
