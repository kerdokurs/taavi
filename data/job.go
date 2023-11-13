package data

import (
	"errors"

	"gorm.io/gorm"
)

type JobType string

const (
	Simple JobType = "SIMPLE"
	Random JobType = "RANDOM"
)

type Job struct {
	gorm.Model

	Type JobType `gorm:"job_type" json:"type"`

	StreamID string `gorm:"stream_id" json:"stream_id"`
	TopicID  string `gorm:"topic_id" json:"topic_id"`
	Content  string `gorm:"content" json:"content"`

	CronTime string `gorm:"cron_time" json:"cron_time"`

	Enabled bool `gorm:"enabled" json:"enabled"`

	Meta []JobMeta `gorm:"meta" json:"meta"`
}

func GetAllJobs() ([]Job, error) {
	var jobs []Job
	return jobs, DB.Model(Job{}).Preload("Meta").Find(&jobs).Error
}

func GetJobs() ([]Job, error) {
	var jobs []Job

	if tx := DB.Model(Job{}).Preload("Meta").Find(&jobs, "enabled = true AND type = 'SIMPLE'"); tx.Error != nil {
		return jobs, tx.Error
	}

	return jobs, nil
}

func GetRandomJobs() ([]Job, error) {
	var jobs []Job
	return jobs, DB.Model(Job{}).Preload("Meta").Find(&jobs, "type = 'RANDOM' AND enabled = true").Error
}

func GetJobMeta(job *Job, key string) (JobMeta, error) {
	for _, meta := range job.Meta {
		if meta.Key == key {
			return meta, nil
		}
	}

	return JobMeta{}, errors.New("not found")
}

func GetJob(id int) (Job, error) {
	var job Job
	return job, DB.Model(Job{}).Preload("Meta").Find(&job, "ID = ?", id).Error
}

func ToggleJobEnabled(id int) (bool, error) {
	var job Job
	if tx := DB.Model(Job{}).Find(&job, "ID = ?", id); tx.Error != nil {
		return false, tx.Error
	}

	return job.Enabled, DB.Model(&job).Update("enabled", !job.Enabled).Error
}

func DeleteJob(id int) error {
	return DB.Delete(&Job{}, "ID = ?", id).Error
}
