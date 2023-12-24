package app

import (
	"github.com/go-ini/ini"
	"kerdo.dev/taavi/pkg/logger"
)

func loadAuthUsers() map[string]string {
	users := make(map[string]string)

	file, err := ini.Load("users.ini")
	if err != nil {
		logger.Fatalw("error loading users file", logger.M{
			"err": err.Error(),
		})
		return users
	}
	keys := file.Section("").Keys()
	for _, key := range keys {
		users[key.Name()] = key.Value()
	}

	return users
}
