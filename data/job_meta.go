package data

import "gorm.io/gorm"

type JobMeta struct {
	gorm.Model `json:"-"`

	Key   string `gorm:"key" json:"key"`
	Value string `gorm:"value" json:"value"`

	JobID int `json:"-"`
}

func GetMetas(jobID int) ([]JobMeta, error) {
	var meta []JobMeta

	return meta, DB.Model(JobMeta{}).Find(&meta, "job_id = ?", jobID).Error
}
