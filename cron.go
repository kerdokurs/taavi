package main

import (
	"fmt"
	"log"
	"taavi/zlp"
	"time"

	"github.com/google/uuid"
	"github.com/robfig/cron/v3"
)

type Job struct {
	Id    uuid.UUID
	Taavi *Taavi

	StreamId string
	TopicId  string
}

type MessageJob struct {
	Job
	Message string
}

func (m MessageJob) Run() {
	message := zlp.Message{
		Stream:  m.StreamId,
		Topic:   m.TopicId,
		Content: m.Message,
	}
	_, err := m.Taavi.bot.Message(&message)
	if err != nil {
		log.Printf("Error sending message: %v\n", err)
	}
}

func (t *Taavi) CronSync(initial bool) {
	// If the bot has just started (meaning initial == true),
	// we shall cancel ALL currently queued tasks except the main scheduler.
	if initial {
		t.deleteAllExceptMain()
	}

	jobInfos, err := t.getTodaysJobs()
	if err != nil {
		log.Printf("Error loading today's jobs: %v\n", err)
		return
	}

	// Starting all jobs based on their infos
	for _, info := range jobInfos {
		if err := t.ScheduleJob(&info); err != nil {
			log.Printf("Error scheduling job: %v\n", err)
		}
	}
}

func (t *Taavi) ScheduleJob(info *CronInfo) error {
	// UUID used to cancel the job after it's done
	taskUid, _ := uuid.NewUUID()

	var job cron.Job
	switch info.Type {
	case Message:
		job = MessageJob{
			Job: Job{
				Id:       taskUid,
				Taavi:    t,
				StreamId: info.StreamId,
				TopicId:  info.TopicId,
			},
			Message: info.Message,
		}
	default:
		return fmt.Errorf("unsupported job type")
	}

	id, err := t.scheduler.AddJob(info.TimeString, job)
	if err != nil {
		return err
	}

	info.SchedulerId = int(id)
	info.ScheduleId = taskUid
	info.ScheduledAt = time.Now()
	t.db.Save(info)

	return nil
}

func (t *Taavi) ScheduledJobs() []cron.Entry {
	return t.scheduler.Entries()
}

func (t *Taavi) deleteAllExceptMain() {
	entries := t.scheduler.Entries()

	for _, entry := range entries {
		if entry.ID == t.mainSchedulerTask {
			continue
		}

		t.scheduler.Remove(entry.ID)
	}
}

func (t *Taavi) cancelByUUID(id uuid.UUID) error {
	info := CronInfo{
		ScheduleId: id,
	}
	tx := t.db.Find(&info)
	if tx.Error != nil {
		return tx.Error
	}

	t.scheduler.Remove(cron.EntryID(info.SchedulerId))
	info.SchedulerId = -1
	t.db.Save(&info)

	return nil
}

func (t *Taavi) getTodaysJobs() ([]CronInfo, error) {
	var allInfos []CronInfo
	if tx := t.db.Find(&allInfos); tx.Error != nil {
		return nil, tx.Error
	}

	// TODO: Pick out today's jobs
	jobs := make([]CronInfo, len(allInfos))
	copy(jobs, allInfos)

	return jobs, nil
}
