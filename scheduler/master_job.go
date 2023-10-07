package scheduler

import (
	"fmt"
	"math/rand"

	"github.com/robfig/cron/v3"
	"kerdo.dev/taavi/data"
	"kerdo.dev/taavi/logger"
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
	randomJobs, err := data.GetRandomJobs()
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

		startMeta, err := data.GetJobMeta(&job, "begin")
		if err != nil {
			panic(err)
		}
		endMeta, err := data.GetJobMeta(&job, "end")
		if err != nil {
			panic(err)
		}

		timeParts := parseStartEnd(startMeta.Value, endMeta.Value)
		rndHour := rand.Intn(timeParts[2]-timeParts[0]) + timeParts[0]
		rndMin := (rand.Intn(6) + 1) * 10
		cronTime := fmt.Sprintf("%d %d * * *", rndMin, rndHour)

		_, err = Scheduler.AddJob(cronTime, cronJob)
		if err != nil {
			logger.Errorw("error scheduling random job", logger.M{
				"err": err.Error(),
			})
		}
	}
}
