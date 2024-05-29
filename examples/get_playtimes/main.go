package main

import (
	"log"

	"github.com/Erb3/switchcraftgo"
)

func main() {
	leaderboard, err := switchcraftgo.GetPlaytimeLeaderboard()

	if err != nil {
		log.Fatalf("Error while fetching leaderboard! %s", err.Error())
	}

	log.Printf("Last Updated: %s", leaderboard.UpdatedAt)
	for idx, player := range leaderboard.Entries {
		log.Printf("#%d: %s with %.0f hours.", idx+1, player.Username, float64(player.Seconds/60/60))
	}
}
