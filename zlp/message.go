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

func (b *Bot) Message(m *Message) (*http.Response, error) {
	if m.Content == "" {
		return nil, fmt.Errorf("message content cannot be empty")
	}

	if len(m.Emails) > 0 {
		return b.PrivateMessage(m)
	}

	if m.Stream == "" {
		return nil, fmt.Errorf("message stream cannot be empty")
	}

	if m.Topic == "" {
		return nil, fmt.Errorf("message topic cannot be empty")
	}

	req, err := b.constructMessageRequest(m)
	if err != nil {
		return nil, err
	}
	return b.Client.Do(req)
}

func (b *Bot) PrivateMessage(m *Message) (*http.Response, error) {
	if len(m.Emails) == 0 {
		return nil, fmt.Errorf("private message must contain atleast one recipient")
	}
	req, err := b.constructMessageRequest(m)
	if err != nil {
		return nil, err
	}
	return b.Client.Do(req)
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
