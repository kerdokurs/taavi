package zulip

import (
	"os"

	"github.com/kerdokurs/zlp"
)

var Client *zlp.Bot

func Init() {
    cfgs := []zlp.ConfigFunction{nil}

    if os.Getenv("TAAVI_ENV") == "dev" {
    }

	if os.Getenv("TAAVI_ENV") == "dev" {
        cfgs[0] = zlp.WithRCFile(".zuliprc")
	} else {
        cfgs[0] = zlp.WithRCEnv()
	}
	Client = zlp.NewBot(cfgs...)
}
