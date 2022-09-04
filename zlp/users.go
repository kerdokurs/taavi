package zlp

import (
	"encoding/json"
	"fmt"
)

type User struct {
	Id            int    `json:"user_id"`
	FullName      string `json:"full_name"`
	Email         string `json:"email"`
	DeliveryEmail string `json:"delivery_email"`
}

func (b *Bot) GetUsers() ([]User, error) {
	type response struct {
		Response
		Members []User `json:"members"`
	}

	var res response
	bytes, err := b.getResponseData("GET", "users", nil)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(bytes, &res); err != nil {
		return nil, err
	}

	return res.Members, nil
}

func (b *Bot) GetSelf() (*User, error) {
	bytes, err := b.getResponseData("GET", "users/me", nil)
	if err != nil {
		return nil, err
	}

	res := new(User)
	if err := json.Unmarshal(bytes, res); err != nil {
		return nil, err
	}

	return res, nil
}

func (b *Bot) GetUser(id int) (*User, error) {
	type response struct {
		User User `json:"user"`
	}

	path := fmt.Sprintf("users/%d", id)
	bytes, err := b.getResponseData("GET", path, nil)
	if err != nil {
		return nil, err
	}

	var res response
	if err := json.Unmarshal(bytes, &res); err != nil {
		return nil, err
	}

	return &res.User, nil
}

func (b *Bot) GetUserByEmail(email string) (*User, error) {
	type response struct {
		Response
		User User `json:"user"`
	}

	path := fmt.Sprintf("users/%s", email)
	bytes, err := b.getResponseData("GET", path, nil)
	if err != nil {
		return nil, err
	}

	var res response
	if err := json.Unmarshal(bytes, &res); err != nil {
		return nil, err
	}

	if res.IsError() {
		return nil, fmt.Errorf(res.Msg)
	}

	return &res.User, nil
}
