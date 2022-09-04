package zlp

import (
	"encoding/json"
	"fmt"
)

type Group struct {
	Id                int    `json:"id"`
	Name              string `json:"name"`
	Description       string `json:"description"`
	Members           []int  `json:"members"`
	DirectSubgroupIds []int  `json:"direct_subgroup_ids"`
	IsSystemGroup     bool   `json:"is_system_group"`
}

func (b *Bot) GetUserGroups() ([]Group, error) {
	body, err := b.getResponseData("GET", "user_groups", nil)
	if err != nil {
		return nil, err
	}

	type response struct {
		Response
		UserGroups []Group `json:"user_groups"`
	}

	var res response
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, err
	}

	if res.IsError() {
		return nil, fmt.Errorf(res.Msg)
	}

	return res.UserGroups, nil
}
