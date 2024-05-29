package switchcraftgo

import "time"

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
		Username string `json:"name"`
		Seconds  int    `json:"time"`
	}
}

// Fetches the leaderboard of who has been the most active player
func GetPlaytimeLeaderboard() (*PlayTimeLeaderboard, error) {
	var leaderboardData struct {
		UpdatedAt string `json:"lastUpdated"`
		Entries   []struct {
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
