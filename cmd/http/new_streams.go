package main

import (
	"encoding/csv"
	"golang.org/x/exp/maps"
	"log"
	"math/rand"
	"os"
	"time"
)

type Team struct {
	id     int64
	name   string
	desc   string
	emails []string
}

/*
csv file
first row is column names
columns for teamname and email can be mapped in function
iterates over each line and groups together emails with same team name
(emails should be Zulip emails, (hopefully UT emails), can be fetched with bot.GetUsers())
returns list of teams that is suitable for bot.NewStreams()
*/

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

	teamMap := make(map[string]Team)

	for i, line := range data {
		if i > 0 {
			// set appropriate column numbers
			teamName := line[2]
			email := line[3]
			if err != nil {
				log.Fatal(err)
			}

			currentTeam, present := teamMap[teamName]
			if present == true {
				currentTeam.emails = append(currentTeam.emails, email)
				teamMap[teamName] = currentTeam
			} else {
				teamMap[teamName] = Team{
					id:     int64(i),
					name:   teamName,
					desc:   "Suvalise nimega striim veebilehe rühmatööks. Oma rühmaliikmete nägemiseks vajuta rühma nime kõrval olevat ikooni. Kui vajad abi, siis küsi oma seminarirühma juhendajalt.",
					emails: []string{email},
				}

			}
		}

	}

	return maps.Values(teamMap)
}

/*
generates random string with given length
*/
func RandomName(n int) string {
	var letterRunes = []rune("abdefghijklmnoprstuvöäüõABDEFGHIJKLMNOPRSTUVÖÄÜÕ0123456789")
	name := make([]rune, n)
	for i := range name {
		name[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(name)
}

/*
creates streams from list
!! if one email is invalid, then entire stream is EMPTY !!
*/
func (t *Taavi) NewStreams(teams []Team) {
	for _, team := range teams {
		t.Bot.CreateStream(team.name, team.desc, true, team.emails)
		time.Sleep(time.Millisecond * 10)
	}
}
