package switchcraftgo

import (
	"encoding/json"
	"io"
	"net/http"
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
