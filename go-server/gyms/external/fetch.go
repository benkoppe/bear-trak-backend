package external

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func FetchData(url string) ([]Gym, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to make the request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var gyms []Gym
	if err := json.Unmarshal(body, &gyms); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return gyms, nil
}
