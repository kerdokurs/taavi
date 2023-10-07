package zulip

import (
	"log"
	"os"

	"github.com/kerdokurs/zlp"
)

var Client *zlp.Bot

func Init() {
	var rc *zlp.ZulipRC
	var err error

	if os.Getenv("TAAVI_ENV") == "dev" {
		rc, err = zlp.LoadRC(".zuliprc")
		if err != nil {
			log.Fatalf("could not parse Zulip configuration")
		}
	} else {
		if rc, err = zlp.RCFromEnv(); err != nil {
			log.Fatalf("could not parse Zulip configuration from env")
		}
	}
	Client = zlp.NewBot(rc)
	Client.Init()
}
