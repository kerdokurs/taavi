package main

import (
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/robfig/cron/v3"
)

type Job struct {
	Id    uuid.UUID
	Taavi *Taavi

	ServerId  string
	ChannelId string
}

type MessageJob struct {
	Job
	Message string
}

func (m MessageJob) Run() {
	log.Printf("Running job with message: %s\n", m.Message)
	// _, err := m.Taavi.bot.ChannelMessageSend(m.ChannelId, m.Message)
	// if err != nil {
	// 	log.Printf("Error sending message for job %s: %v\n", m.Id.String(), err)
	// }
}

func (t *Taavi) CronSync(initial bool) {
	if initial && t.mainSchedulerTask == -1 {
		// The first initial setup
		entries := t.scheduler.Entries()

		for _, entry := range entries {
			if entry.ID == t.mainSchedulerTask {
				continue
			}

			t.scheduler.Remove(entry.ID)
		}
	}

	var jobInfos []CronInfo
	if initial {
		t.db.Find(&jobInfos)
	} else {
		t.db.Where("scheduled_at < updated_at").Find(&jobInfos)
	}

	for _, info := range jobInfos {
		sid, err := uuid.NewUUID()
		if err != nil {
			log.Printf("Error creating new UUID for job %d: %v\n", info.ID, err)
			continue
		}

		var job cron.Job
		switch info.Type {
		case Message:
			job = MessageJob{
				Job: Job{
					Id:        sid,
					Taavi:     t,
					ServerId:  info.ServerId,
					ChannelId: info.ChannelId,
				},
				Message: info.Message,
			}
		default:
			log.Printf("Unsupported job type: %d\n", info.Type)
		}

		id, err := t.scheduler.AddJob(info.TimeString, job)
		if err != nil {
			log.Printf("Could not start job with id %d: %v\n", info.ID, err)
			continue
		}

		info.SchedulerId = int(id)
		info.ScheduleId = sid
		info.ScheduledAt = time.Now()
		t.db.Save(&info)
	}
}

func (t *Taavi) ScheduledJobs() []cron.Entry {
	return t.scheduler.Entries()
}
