package zlp

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type Message struct {
	Stream  string
	Topic   string
	Emails  []string
	Content string
}

func (b *Bot) Message(m *Message) error {
	if m.Content == "" {
		return fmt.Errorf("message content cannot be empty")
	}

	if len(m.Emails) > 0 {
		return b.PrivateMessage(m)
	}

	if m.Stream == "" {
		return fmt.Errorf("message stream cannot be empty")
	}

	if m.Topic == "" {
		return fmt.Errorf("message topic cannot be empty")
	}

	req, err := b.constructMessageRequest(m)
	if err != nil {
		return err
	}
	resp, err := b.Client.Do(req)
	if err != nil {
		return err
	}
	return b.respToError(resp)
}

func (b *Bot) PrivateMessage(m *Message) error {
	if len(m.Emails) == 0 {
		return fmt.Errorf("private message must contain atleast one recipient")
	}
	req, err := b.constructMessageRequest(m)
	if err != nil {
		return err
	}
	resp, err := b.Client.Do(req)
	if err != nil {
		return err
	}
	return b.respToError(resp)
}

func (b *Bot) constructMessageRequest(m *Message) (*http.Request, error) {
	to := m.Stream
	messageType := "stream"

	if len(m.Emails) > 0 {
		messageType = "private"
		to = strings.Join(m.Emails, ",")
	}

	values := url.Values{}
	values.Set("type", messageType)
	values.Set("to", to)
	values.Set("content", m.Content)
	if messageType == "stream" {
		values.Set("subject", m.Topic)
	}

	return b.constructRequest("POST", "messages", &values)
}
