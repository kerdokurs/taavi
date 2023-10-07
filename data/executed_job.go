package data

import (
	"time"

	"gorm.io/gorm"
)

type ExecutedJob struct {
	gorm.Model

	RanAt time.Time `gorm:"ran_at" json:"ran_at"`

	Error string `gorm:"error" json:"error"`

	JobID int
	Job   Job
}
