package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/kerdokurs/zlp"
	"github.com/robfig/cron/v3"
)

type Job struct {
	Bot   *zlp.Bot
	Taavi *Taavi
}

type UpdateUrlJob struct {
	Job

	StreamId string
	TopicId  string
	Emails   []string
	Content  string
	Url      string
}

func (uj UpdateUrlJob) Run() {
	dayNr := time.Now().Day()
	// current format with advent of code url (https://adventofcode.com/2021/day/<current day>)
	mesContent := uj.Content + uj.Url + strconv.Itoa(dayNr)
	message := zlp.Message{
		Stream:  uj.StreamId,
		Topic:   uj.TopicId,
		Content: mesContent,
		Emails:  uj.Emails,
	}
	err := uj.Bot.Message(&message)
	if err != nil {
		log.Printf("Error sending message: #{err}\n")
		return
	}
	to := ""
	if len(uj.Emails) > 0 {
		to = strings.Join(uj.Emails, ", ")
	} else {
		to = fmt.Sprintf("%s:%s", uj.StreamId, uj.TopicId)
	}
	log.Printf("Ran message job to %s\n", to)
}

type MessageJob struct {
	Job

	StreamId string
	TopicId  string
	Emails   []string
	Content  string
}

func (m MessageJob) Run() {
	message := zlp.Message{
		Stream:  m.StreamId,
		Topic:   m.TopicId,
		Content: m.Content,
		Emails:  m.Emails,
	}
	err := m.Bot.Message(&message)
	if err != nil {
		log.Printf("Error sending message: %v\n", err)
		return
	}
	to := ""
	if len(m.Emails) > 0 {
		to = strings.Join(m.Emails, ", ")
	} else {
		to = fmt.Sprintf("%s:%s", m.StreamId, m.TopicId)
	}
	log.Printf("Ran message job to %s\n", to)
}

var cronInfos = NewObservable(make([]CronInfo, 0))
var scheduledJobs = make(map[string]*CronInfo)

type CronService struct {
	Taavi     *Taavi
	scheduler *cron.Cron
}

func (cs *CronService) Init() {
	cronInfos.Subscribe(cs)

	loc, err := time.LoadLocation("Europe/Tallinn")
	if err != nil {
		log.Fatalf("Error loading timezone: %v\n", err)
	}

	cs.scheduler = cron.New(cron.WithLocation(loc))
	cs.scheduler.Start()
}

func (cs *CronService) Stop() {
	ctx := cs.scheduler.Stop()
	<-ctx.Done()
	cronInfos.Close()
}

func (cs *CronService) OnChange(old *[]CronInfo, val *[]CronInfo) {
	// fmt.Println("CronInfos OnChange")
	// fmt.Printf("%+v\n", val)

	for _, info := range *val {
		if scheduledJob, ok := scheduledJobs[info.Id]; ok {
			// Job is scheduled. But did the configuration change?
			if scheduledJob.Changed(&info) {
				// It did, indeed.
				log.Printf("Cron configuration for job with id %s changed. Rescheduling.\n", info.Id)

				// Reschedule the job
				cs.scheduler.Remove(scheduledJob.EntryId)
				if err := cs.ScheduleJob(&info); err != nil {
					log.Printf("Could not schedule job with id %s and cron time %s: %v\n", info.Id, info.CronTime, err)
				}
			}
		} else {
			// Job is not scheduled. Schedule it!
			cs.ScheduleJob(&info)
		}
	}
}

func (cs *CronService) ScheduleJob(info *CronInfo) error {
	var job cron.Job
	switch info.Type {
	case Message:
		job = MessageJob{
			Job: Job{
				Bot: cs.Taavi.Bot,
			},
			StreamId: info.StreamId,
			TopicId:  info.TopicId,
			Emails:   info.Emails,
			Content:  info.Content,
		}
	case Update:
		job = UpdateUrlJob{
			Job: Job{
				Bot: cs.Taavi.Bot,
			},
			StreamId: info.StreamId,
			TopicId:  info.TopicId,
			Emails:   info.Emails,
			Content:  info.Content,
			Url:      info.Url,
		}
	default:
		return fmt.Errorf("unsupported job type")
	}
	id, err := cs.scheduler.AddJob(info.CronTime, job)
	if err != nil {
		return err
	}

	scheduledJob := *info
	scheduledJob.EntryId = id
	scheduledJobs[info.Id] = &scheduledJob
	log.Printf("Scheduled job for %s\n", info.CronTime)

	return nil
}
