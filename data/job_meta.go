package data

import "gorm.io/gorm"

type JobMeta struct {
	gorm.Model `json:"-"`

	Key   string `gorm:"key" json:"key"`
	Value string `gorm:"value" json:"value"`

	JobID int `json:"-"`
}
