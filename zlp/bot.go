package zlp

import (
	"net/http"
)

const DefaultApiVersion = "v1"

type Bot struct {
	Email  string
	Key    string
	ApiUrl string

	ApiVersion string

	Client Doer
}

func NewBot(rc *ZulipRC) *Bot {
	return &Bot{
		Email:  rc.Email,
		Key:    rc.APIKey,
		ApiUrl: rc.APIUrl,
	}
}

type Doer interface {
	Do(*http.Request) (*http.Response, error)
}

func (b *Bot) Init() {
	b.Email = "taavi-test-bot@se-22.zulip.cs.ut.ee"
	b.Key = "V1FY7ESgEIOpndh04IPNbr6kCgbKLzdC"
	b.ApiUrl = "https://se-22.zulip.cs.ut.ee"

	b.Client = &http.Client{}
	if b.ApiVersion == "" {
		b.ApiVersion = DefaultApiVersion
	}
}
