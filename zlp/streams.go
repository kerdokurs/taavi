package zlp

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Stream struct {
	Id          int    `json:"stream_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	InviteOnly  bool   `json:"invite_only"`
	// CreatedAt   time.Time `json:"date_created"`
}

func (b *Bot) GetSubscribedStreams() ([]Stream, error) {
	req, err := b.constructRequest("GET", "users/me/subscriptions", nil)
	if err != nil {
		return nil, err
	}

	type response struct {
		Result        string   `json:"result"`
		Msg           string   `json:"msg"`
		Subscriptions []Stream `json:"subscriptions"`
	}

	resp, err := b.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var res response
	body, err := ioutil.ReadAll(req.Body)
	fmt.Println(body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return nil, err
	}

	return res.Subscriptions, nil
}

func (b *Bot) GetStreams() ([]Stream, error) {
	req, err := b.constructRequest("GET", "streams", nil)
	if err != nil {
		return nil, err
	}
	resp, err := b.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	type response struct {
		Msg     string   `json:"msg"`
		Streams []Stream `json:"streams"`
		Result  string   `json:"result"`
	}

	var res response
	err = json.Unmarshal(body, &res)
	if err != nil {
		return nil, err
	}

	return res.Streams, nil
}

func (b *Bot) GetStream(id int) (*Stream, error) {
	req, err := b.constructRequest("GET", fmt.Sprintf("streams/%d", id), nil)
	if err != nil {
		return nil, err
	}
	resp, err := b.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	type response struct {
		Response
		Msg    string `json:"msg"`
		Stream Stream `json:"stream"`
		Result string `json:"result"`
	}

	var res response
	err = json.Unmarshal(body, &res)
	if err != nil {
		return nil, err
	}
	stream := res.Stream

	return &stream, nil
}
