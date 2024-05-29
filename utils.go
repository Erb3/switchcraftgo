package switchcraftgo

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"
)

func makeGetRequest(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "erb3/switchcraftgo/v0")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func makeGetJsonRequest(url string, v any) error {
	resp, err := makeGetRequest(url)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	json.Unmarshal(body, &v)
	return nil
}

// parseScTimestamp parses a timestamp string with the format `2022-10-22T20:50:45.6221607+01:00`
// or `2024-05-29T15:16:28.042866052Z`, ignoring any bracketed timezone information.
// Returns a time.Time object, or an error.
func parseScTimestamp(timestamp string) (time.Time, error) {
	// Remove any bracketed timezone information
	if bracketPos := strings.Index(timestamp, "["); bracketPos != -1 {
		timestamp = timestamp[:bracketPos]
	}

	// Define the layout that matches the timestamp format
	const layout = "2006-01-02T15:04:05.999999999Z07:00"

	// Parse the timestamp string using the specified layout
	t, err := time.Parse(layout, timestamp)
	if err != nil {
		return time.Time{}, err
	}

	return t, nil
}
