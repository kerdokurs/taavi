package app

import (
	"fmt"
	"os"
)

func GetListenAddr() string {
	host := os.Getenv("HOST")
	if host == "" {
		host = "127.0.0.1"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8088"
	}

	return fmt.Sprintf("%s:%s", host, port)
}
