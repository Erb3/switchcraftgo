package switchcraftgo

import (
	"fmt"
	"log"
	"strings"
	"testing"
)

func TestGetPlayerCounts(t *testing.T) {
	res, err := GetPlayerCounts()

	if err != nil {
		t.Fatalf("GetPlayerCounts() returned error %s", err.Error())
	}

	if res.Active > res.Total {
		t.Fatalf("GetPlayerCounts() returned more active players than total players. Received %d active out of %d total", res.Active, res.Total)
	}

	if res.Active < 0 {
		t.Fatalf("GetPlayerCounts() retured active below zero (received %d)", res.Active)
	}

	if res.Total < 0 {
		t.Fatalf("GetPlayerCounts() retured total below zero (received %d)", res.Total)
	}
}

func ExampleGetPlayerCounts() {
	counts, err := GetPlayerCounts()
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("There are currently %d players online.", counts.Total)
}

func TestGetTps(t *testing.T) {
	res, err := GetTps()

	if err != nil {
		t.Fatalf("GetTps() returned error %s", err.Error())
	}

	if res.AverageTps == 0 {
		t.Fatalf("GetTps() returned TPS of zero, less than zero, or nil (got %f)", res.AverageTps)
	}

	if res.AverageMillisecondPerTick == 0 {
		t.Fatalf("GetTps() returned average milliseconds of zero, less than zero, or nil (got %f)", res.AverageMillisecondPerTick)
	}

	if res.MillisecondsLastTick == 0 {
		t.Fatalf("GetTps() returned last milliseconds of zero, less than zero, or nil (got %f)", res.MillisecondsLastTick)
	}
}

func ExampleGetTps() {
	tpsData, err := GetTps()
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("Average TPS: %f", tpsData.AverageTps)
}

func TestGetPlaytimeLeaderboard(t *testing.T) {
	res, err := GetPlaytimeLeaderboard()

	if err != nil {
		t.Fatalf("GetPlaytimeLeaderboard() returned error %s", err.Error())
	}

	if res.UpdatedAt == nil {
		t.Fatalf("GetPlaytimeLeaderboard() returned no updated at value")
	}

	entry := res.Entries[0]
	if entry.Username == "" {
		t.Fatalf("GetPlaytimeLeaderboard() returned empty username on first entry")
	}

	if entry.Seconds <= 0 {
		t.Fatalf("GetPlaytimeLeaderboard() returned zero seconds, less than zero, or nil on first entry (got %d)", entry.Seconds)
	}
}

func ExampleGetPlaytimeLeaderboard() {
	leaderboard, err := GetPlaytimeLeaderboard()

	if err != nil {
		log.Fatalln(err)
	}

	for idx, player := range leaderboard.Entries {
		log.Printf("#%d: %s with %.0f hours.", idx+1, player.Username, float64(player.Seconds/60/60))
	}
}

func TestGetSupporterGoal(t *testing.T) {
	res, err := GetSupporterGoal()

	if err != nil {
		t.Fatalf("GetSupporterGoal() returned error %s", err.Error())
	}

	if res.SupporterUrl == "" {
		t.Fatalf("GetSupporterGoal() returned no supporter url")
	}

	if res.Goal <= 0 {
		t.Fatalf("GetSupporterGoal() returned goal less than 1$ (got %f)", res.Goal)
	}

	if res.Current >= res.Goal != res.GoalMet {
		t.Fatalf("GetSupporterGoal() returned wrong GoalMet field. Got %f current, and %f as a goal.", res.Current, res.Goal)
	}
}

func ExampleGetSupporterGoal() {
	goal, err := GetSupporterGoal()
	if err != nil {
		log.Fatalln(err)
	}

	if !goal.GoalMet {
		fmt.Printf("The supporter goal has not been reached this month. Consider supporting! %f has been donated so far this month, out of the %f required to operate.", goal.Current, goal.Goal)
	}
}

func TestGetDeathsLeaderboard(t *testing.T) {
	deaths, err := GetDeathsLeaderboard()

	if err != nil {
		t.Fatalf("GetDeathsLeaderboard() returned error %s", err.Error())
	}

	if deaths[0].Name == "" {
		t.Fatalf("GetDeathsLeaderboard() returned no name for first entry")
	}

	if deaths[0].Uuid == "" {
		t.Fatalf("GetDeathsLeaderboard() returned no uuid for first entry")
	}

	if deaths[0].Count == 0 {
		t.Fatalf("GetDeathsLeaderboard() returned 0 deaths for first entry")
	}

	last := -1
	for _, player := range deaths {
		if last < player.Count && last != -1 {
			t.Fatalf("GetDeathsLeaderboard() did not return leaderboard in descending order")
		}

		last = player.Count
	}
}

func ExampleGetDeathsLeaderboard() {
	deaths, err := GetDeathsLeaderboard()
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("%s has most deaths, with %d deaths", deaths[0].Name, deaths[0].Count)
}

func TestGetProxyRanges(t *testing.T) {
	proxies, err := GetProxyRanges()

	if err != nil {
		t.Fatalf("GetProxyRanges() returned error %s", err.Error())
	}

	if len(proxies) == 0 {
		t.Fatalf("GetProxyRanges() received no proxies")
	}
}

func ExampleGetProxyRanges() {
	proxies, err := GetProxyRanges()
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("Current proxies: %s", strings.Join(proxies, ", "))
}
