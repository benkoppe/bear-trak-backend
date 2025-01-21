package external_availtec

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
)

func FetchRoutes(baseUrl string) ([]Route, error) {
	fullUrl, err := url.Parse(baseUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to parse base URL: %w", err)
	}

	fullUrl.Path = path.Join(fullUrl.Path, "Routes/GetVisibleRoutes")

	resp, err := http.Get(fullUrl.String())
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

	var routes []Route
	if err := json.Unmarshal(body, &routes); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return routes, nil
}
