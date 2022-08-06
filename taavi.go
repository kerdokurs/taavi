package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

type Taavi struct {
	bot *discordgo.Session

	keepAlive chan os.Signal
}

func NewTaavi() *Taavi {
	token := os.Getenv("DISCORD_TOKEN")
	if token == "" {
		log.Fatal("Discord token is not defined")
	}

	bot, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalf("Could not create new bot: %v\n", err)
	}

	return &Taavi{
		bot:       bot,
		keepAlive: make(chan os.Signal, 1),
	}
}

func (t *Taavi) Start() {
	err := t.bot.Open()
	if err != nil {
		log.Fatalf("Could not open bot: %s\n", err)
	}

	log.Print("Tiiger Taavi is now up!")

	signal.Notify(t.keepAlive, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-t.keepAlive

	log.Print("Tiiger Taavi is now shutting down.")
}
