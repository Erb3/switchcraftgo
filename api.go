package switchcraftgo

import (
	"strconv"
	"time"
)

type PlayerCountsResponse struct {
	Total  int `json:"total"`
	Active int `json:"active"`
}

// Fetches the current player counts from the API.
// If successful, it returns the response, otherwise it returns the error.
func GetPlayerCounts() (*PlayerCountsResponse, error) {
	var counts PlayerCountsResponse

	err := makeGetJsonRequest("https://api.sc3.io/v3/players", &counts)
	if err != nil {
		return nil, err
	}

	return &counts, nil
}

type TpsResponse struct {
	AverageTps                float32 `json:"tps"`
	AverageMillisecondPerTick float32 `json:"avgMsPerTick"`
	MillisecondsLastTick      float32 `json:"lastMsPerTick"`
}

// Fetches information regarding ticks per seconds.
//
// If successful, it returns a variety of statistics:
//   - The average ticks per second, averaged over the last 100 ticks
//   - The time per tick, in milliseconds, averaged over the last 100 ticks
//   - The time per tick, in milliseconds, for the last tick
//
// If the request failed, it will return nil, and the the error.
func GetTps() (*TpsResponse, error) {
	var tps TpsResponse

	err := makeGetJsonRequest("https://api.sc3.io/v3/tps", &tps)
	if err != nil {
		return nil, err
	}

	return &tps, nil
}

type PlayTimeLeaderboard struct {
	UpdatedAt *time.Time
	Entries   []struct {
		Uuid     string `json:"uuid"`
		Username string `json:"name"`
		Seconds  int    `json:"time"`
	}
}

// Fetches the leaderboard of who has been the most active player.
// Players are able to opt-out of this leaderboard.
// Sorted in descending order.
func GetPlaytimeLeaderboard() (*PlayTimeLeaderboard, error) {
	var leaderboardData struct {
		UpdatedAt string `json:"lastUpdated"`
		Entries   []struct {
			Uuid     string `json:"uuid"`
			Username string `json:"name"`
			Seconds  int    `json:"time"`
		} `json:"entries"`
	}

	err := makeGetJsonRequest("https://api.sc3.io/v3/activetime", &leaderboardData)
	if err != nil {
		return nil, err
	}

	timestamp, err := parseScTimestamp(leaderboardData.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &PlayTimeLeaderboard{
		UpdatedAt: &timestamp,
		Entries:   leaderboardData.Entries,
	}, nil
}

type SupporterGoal struct {
	SupporterUrl string  `json:"supporterUrl"`
	Current      float64 `json:"current"`
	Goal         float64 `json:"goal"`
	GoalMet      bool    `json:"goalMet"`
}

// Fetches information about the supporter goal
//
// Assuming success, the following information is returned:
//   - SupporterUrl is a string containing the URL to support SwitchCraft
//   - Current is a float of the current donated money this month, in dollars
//   - Goal is a float of the supporter goal, in dollars
//   - GoalMet reports whether the goal has been met
//
// In the case that something fails, the first return value is nil and the second is the error
func GetSupporterGoal() (*SupporterGoal, error) {
	var supporterData struct {
		SupporterUrl string `json:"supporterUrl"`
		Current      string `json:"current"`
		Goal         string `json:"goal"`
		GoalMet      bool   `json:"goalMet"`
	}

	err := makeGetJsonRequest("https://api.sc3.io/v3/supporter", &supporterData)
	if err != nil {
		return nil, err
	}

	current, err := strconv.ParseFloat(supporterData.Current, 64)
	if err != nil {
		return nil, err
	}

	goal, err := strconv.ParseFloat(supporterData.Goal, 64)
	if err != nil {
		return nil, err
	}

	return &SupporterGoal{
		SupporterUrl: supporterData.SupporterUrl,
		Current:      current,
		Goal:         goal,
		GoalMet:      supporterData.GoalMet,
	}, nil
}

type DeathsLeaderboard []struct {
	Uuid  string `json:"uuid"`
	Name  string `json:"name"`
	Count int    `json:"count"`
}

// Fetches the leaderboard of most deaths.
//
// Returns an array of players along with the death count.
// The array is sorted in descending order by most deaths.
func GetDeathsLeaderboard() (DeathsLeaderboard, error) {
	var leaderboard DeathsLeaderboard

	err := makeGetJsonRequest("https://api.sc3.io/v3/deaths", &leaderboard)
	if err != nil {
		return nil, err
	}

	return leaderboard, nil
}

// Fetches the current ComputerCraft proxies
//
// Returns an array of IP ranges in CIDR notation that HTTP requests
// originate from when sent from ComputerCraft.
// This information can be used to block any traffic not originating from inside SwitchCraft.
// For an IPv4 address, a CIDR range of /32 means that the IP address will match exactly.
func GetProxyRanges() ([]string, error) {
	var proxyInfo struct {
		Proxies []string `json:"httpProxies"`
	}

	err := makeGetJsonRequest("https://api.sc3.io/v3/proxies", &proxyInfo)
	if err != nil {
		return nil, err
	}

	return proxyInfo.Proxies, nil
}
