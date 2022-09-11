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
	b.Client = &http.Client{}
	if b.ApiVersion == "" {
		b.ApiVersion = DefaultApiVersion
	}
}
