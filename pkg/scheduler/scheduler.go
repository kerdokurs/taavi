package scheduler

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
	"kerdo.dev/taavi/pkg/data"
	"kerdo.dev/taavi/pkg/logger"
)

var Scheduler *cron.Cron

func Init() {
	loc, err := time.LoadLocation("Europe/Tallinn")
	if err != nil {
		logger.Fatalw("error loading location", logger.M{
			"err": err.Error(),
		})
	}

	Scheduler = cron.New(cron.WithLocation(loc))
	Scheduler.Start()
}

func ScheduleAll() {
	ctx := context.TODO()
	jobs, err := data.GetJobs(ctx)
	if err != nil {
		logger.Errorw("error getting jobs for scheduling", logger.M{
			"err": err.Error(),
		})
	}

	for _, job := range jobs {
		id, err := ScheduleJob(&job)

		if err != nil {
			logger.Errorw("error scheduling job", logger.M{
				"err": err.Error(),
				"id":  job.ID,
			})
			break
		}
		logger.Infow("scheduled job", logger.M{
			"id":        id,
			"stream_id": job.StreamID,
			"topic_id":  job.TopicID,
			"cron_time": job.CronTime,
		})
	}
}

func ScheduleMaster() {
	var masterJob cron.Job = &MasterJob{}
	_, err := Scheduler.AddJob(MasterJobCronTime, masterJob)
	if err != nil {
		logger.Errorw("error scheduling master job", logger.M{
			"err": err.Error(),
		})
		return
	}
	logger.Infow("scheduled master job", nil)
}

func ScheduleJob(job *data.Job) (cron.EntryID, error) {
	if cronJob, cronTime, err := CreateJob(job); err != nil {
		return -1, err
	} else {
		return Scheduler.AddJob(cronTime, cronJob)
	}
}

func CreateJob(job *data.Job) (cron.Job, string, error) {
	var cronJob cron.Job
	var cronTime string

	switch job.Type {
	case data.Simple:
		cronTime = job.CronTime
		cronJob = &MessageJob{
			JobID:    int(job.ID),
			StreamID: job.StreamID,
			TopicID:  job.TopicID,
			Content:  job.Content,
		}

	case data.Random:
		return nil, "", errors.New("random jobs must be scheduled by the master scheduler")
	default:
		return nil, "", fmt.Errorf("unsupported error type %s", job.Type)
	}

	return cronJob, cronTime, nil
}

func UnscheduleAll() {
	entryIDs := make([]cron.EntryID, 0)
	for _, entry := range Scheduler.Entries() {
		entryIDs = append(entryIDs, entry.ID)
	}

	for _, entryID := range entryIDs {
		Scheduler.Remove(entryID)
	}
}

func RescheduleAll() {
	logger.Infow("rescheduling jobs", nil)

	UnscheduleAll()
	ScheduleMaster()
	ScheduleAll()
}
