package scheduler

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/robfig/cron/v3"
	"kerdo.dev/taavi/logger"
	"kerdo.dev/taavi/pkg/data"
)

const MasterJobCronTime = "0 8 * * *"

// Goals of MasterJob
// 1. Schedule all random-timed jobs for current day

type MasterJob struct {
}

func (m *MasterJob) Run() {
	m.scheduleRandomJobs()
}

func (m *MasterJob) scheduleRandomJobs() {
	ctx := context.TODO()
	randomJobs, err := data.GetRandomJobs(ctx)
	if err != nil {
		logger.Errorw("error getting random jobs", logger.M{
			"err": err.Error(),
		})
		return
	}

	for _, job := range randomJobs {
		var cronJob cron.Job = &MessageJob{
			JobID:    int(job.ID),
			StreamID: job.StreamID,
			TopicID:  job.TopicID,
			Content:  job.Content,
		}

		startMeta, err := data.GetJobMeta(ctx, &job, "begin")
		if err != nil {
			logger.Errorw("error getting begin meta for random job", logger.M{
				"err": err.Error(),
				"job": job,
			})
			continue
		}
		endMeta, err := data.GetJobMeta(ctx, &job, "end")
		if err != nil {
			logger.Errorw("error getting end meta for random job", logger.M{
				"err": err.Error(),
				"job": job,
			})
			continue
		}

		timeParts := parseStartEnd(startMeta.Value, endMeta.Value)
		rndHour := rand.Intn(timeParts[2]-timeParts[0]) + timeParts[0]
		rndMin := rand.Intn(6) * 10
		cronTime := fmt.Sprintf("%d %d * * *", rndMin, rndHour)

		id, err := Scheduler.AddJob(cronTime, cronJob)
		if err != nil {
			logger.Errorw("error scheduling random job", logger.M{
				"err": err.Error(),
			})
			continue
		}
		logger.Infow("scheduled random job", logger.M{
			"id":        id,
			"stream_id": job.StreamID,
			"topic_id":  job.TopicID,
			"cron_time": cronTime,
		})
	}
}
