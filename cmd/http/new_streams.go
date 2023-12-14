package main

import (
	"encoding/csv"
	"golang.org/x/exp/maps"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

type Team struct {
	id     int64
	name   string
	desc   string
	emails []string
}

func getTeamsFromFile(fileName string) []Team {
	f, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}

	defer func(f *os.File) {
		err := f.Close()
		if err != nil {

		}
	}(f)

	csvReader := csv.NewReader(f)
	data, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	teamMap := make(map[int64]Team)

	for i, line := range data {
		if i > 0 {
			teamField := line[5]
			email := line[7]
			idStr := strings.Split(teamField, " ")[1]
			idInt, err := strconv.ParseInt(string(idStr), 0, 8)
			if err != nil {
				log.Fatal(err)
			}

			currentTeam, present := teamMap[idInt]
			if present == true {
				currentTeam.emails = append(currentTeam.emails, email)
				teamMap[idInt] = currentTeam
			} else {
				teamMap[idInt] = Team{
					id:     idInt,
					name:   RandomName(5),
					desc:   "Suvalise nimega striim veebilehe rühmatööks. Oma rühmaliikmete nägemiseks vajuta rühma nime kõrval olevat ikooni. Kui vajad abi, siis küsi oma seminarirühma juhendajalt.",
					emails: []string{email},
				}

			}
		}

	}

	return maps.Values(teamMap)
}

func RandomName(n int) string {
	var letterRunes = []rune("abdefghijklmnoprstuvöäüõABDEFGHIJKLMNOPRSTUVÖÄÜÕ0123456789")
	name := make([]rune, n)
	for i := range name {
		name[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(name)
}

func (t *Taavi) NewStreams(teams []Team) {
	for _, team := range teams {
		t.Bot.CreateStream(team.name, team.desc, true, team.emails)
		time.Sleep(time.Millisecond * 10)
	}
}
