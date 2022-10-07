package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/kerdokurs/zlp"
)

type Taavi struct {
	Bot  *zlp.Bot
	Db   DatabaseService
	Cron CronService

	keepAlive chan os.Signal
}

func NewTaavi() *Taavi {
	// Connection with Zulip
	var rc *zlp.ZulipRC
	var err error
	if os.Getenv("ENV") == "prod" {
		rc, err = zlp.RCFromEnv()
	} else {
		rc, err = zlp.LoadRC(".zuliprc")
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not load ZulipRC file: %v\n", err)
		os.Exit(1)
	}
	bot := zlp.NewBot(rc)
	bot.Init()

	return &Taavi{
		Bot:       bot,
		keepAlive: make(chan os.Signal, 1),
	}
}

func (t *Taavi) Start() {
	t.Cron = CronService{
		Taavi: t,
	}
	t.Cron.Init()
	defer t.Cron.Stop()

	t.Db = DatabaseService{}
	t.Db.Init()
	defer t.Db.Stop()

	log.Print("Tiiger Taavi is now up!")

	signal.Notify(t.keepAlive, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-t.keepAlive

	log.Print("Tiiger Taavi is now shutting down.")
}
