package main

import (
	"log"

	"github.com/Erb3/switchcraftgo"
)

func main() {
	players, err := switchcraftgo.GetPlayerCounts()

	if err != nil {
		log.Fatalf("Error while fetching player counts! %s", err.Error())
	}

	log.Printf("There are %d players online, %d of which are active.", players.Total, players.Active)
}
