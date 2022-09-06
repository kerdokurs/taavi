package main

import (
	"log"
	"time"

	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const databasePath = "test.db"

func (t *Taavi) setupDb() {
	var err error
	t.db, err = gorm.Open(sqlite.Open(databasePath), &gorm.Config{})
	if err != nil {
		log.Fatalf("Could not open database %s: %v\n", databasePath, err)
	}

	t.db.AutoMigrate(&CronInfo{})
}

type CronType int

const (
	Message CronType = iota
)

type CronInfo struct {
	gorm.Model

	TimeString string
	Type       CronType

	ScheduledAt time.Time
	SchedulerId int
	ScheduleId  uuid.UUID

	StreamId string
	TopicId  string

	Message string
}
