package zulip

import (
	"os"

	"github.com/kerdokurs/zlp"
)

var Client *zlp.Bot

func Init() {
    cfgs := []zlp.ConfigFunction{
        zlp.WithUserAgent("Taavi/2.0"),
        nil,
    }

	if os.Getenv("TAAVI_ENV") == "dev" {
        cfgs[1] = zlp.WithRCFile(".zuliprc")
	} else {
        cfgs[1] = zlp.WithRCEnv()
	}
	Client = zlp.NewBot(cfgs...)
}
