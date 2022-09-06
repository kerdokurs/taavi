package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
	"kerdo.dev/taavi/zlp"
)

type Taavi struct {
	bot *zlp.Bot
	db  *gorm.DB

	scheduler         *cron.Cron
	mainSchedulerTask cron.EntryID

	keepAlive chan os.Signal
}

func NewTaavi() *Taavi {
	// Connection with Zulip
	rc, err := zlp.LoadRC(".zuliprc")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not load ZulipRC file: %v\n", err)
		os.Exit(1)
	}
	bot := zlp.NewBot(rc)
	bot.Init()

	return &Taavi{
		bot:               bot,
		scheduler:         cron.New(),
		mainSchedulerTask: -1,
		keepAlive:         make(chan os.Signal, 1),
	}
}

func (t *Taavi) Start() {
	t.setupDb()

	// Sync up CRON
	t.CronSync(true)
	go t.scheduler.Run()

	// Main task scheduler is the task that starts up each day's tasks at midnight
	// t.mainSchedulerTask, err = t.scheduler.AddFunc("@every 5s", func() {
	// 	t.CronSync(true)
	// })
	// if err != nil {
	// 	log.Fatalf("Could not start main job scheduler: %v\n", err)
	// }
	go t.runGRPC()

	log.Print("Tiiger Taavi is now up!")

	signal.Notify(t.keepAlive, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-t.keepAlive

	ForEach(t.scheduler.Entries(), func(entry cron.Entry) {
		t.scheduler.Remove(entry.ID)
	})
	ctx := t.scheduler.Stop()
	<-ctx.Done()

	// TODO: Stop gRPC HTTP server

	log.Print("Tiiger Taavi is now shutting down.")
}
