package scheduler

import (
	"context"
	"time"

	"github.com/kerdokurs/zlp"
	"kerdo.dev/taavi/data"
	"kerdo.dev/taavi/logger"
	"kerdo.dev/taavi/zulip"
)

type MessageJob struct {
	JobID    int
	StreamID string
	TopicID  string
	Content  string
}

func (m *MessageJob) Run() {
	content := replaceVariables(context.Background(), m.Content)

	msg := zlp.Message{
		Stream:  m.StreamID,
		Topic:   m.TopicID,
		Content: content,
	}
	err := zulip.Client.Message(&msg)

	executedJob := data.ExecutedJob{
		RanAt: time.Now(),
		JobID: m.JobID,
	}
	defer func() {
		data.DB.Create(&executedJob)
	}()

	if err != nil {
		logger.Errorw("error sending message", logger.M{
			"err": err.Error(),
		})
		executedJob.Error = err.Error()
		return
	}
	logger.Infow("message sent", logger.M{
		"topic":  msg.Topic,
		"stream": msg.Stream,
	})
}
