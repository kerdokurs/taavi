package main

func main() {
	taavi := NewTaavi()
	taavi.Start()

	// koos seediga teeb iga kord uued randomid

	//taavi.NewStreams(getTeamsFromFile("tudengid.csv"))
	//users, _ := taavi.Bot.GetUsers()
	//for _, user := range users {
	//	fmt.Println(user)
	//}
	//fmt.Println(len(users))
	//teams := getTeamsFromFile("tudengid.csv")
	//team := Team{
	//	id:     1,2
	//	name:   RandomName(5),
	//	desc:   "Suvalise nimega striim veebilehe rühmatööks. Oma rühmaliikmete nägemiseks vajuta rühma nime kõrval olevat ikooni. Kui vajad abi, siis küsi oma seminarirühma juhendajalt.",
	//	emails: []string{"priidik.vastrik@ut.ee", "epp.haavasalu@ut.ee", "mirjam.paales@ut.ee"},
	//}

	//taavi.Bot.CreateStream(team.name, team.desc, true, team.emails)

}
