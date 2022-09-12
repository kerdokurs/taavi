package main

import (
	"fmt"
	"log"
	"time"
)

type logWriter struct {
}

func (writer logWriter) Write(bytes []byte) (int, error) {
	loc, err := time.LoadLocation("Europe/Tallinn")
	if err != nil {
		log.Fatalf("Error loading timezone: %v\n", err)
	}
	return fmt.Print(time.Now().In(loc).Format("2006/01/02 15:04:00") + " " + string(bytes))
}

func init() {
	log.SetFlags(0)
	log.SetOutput(new(logWriter))
}
