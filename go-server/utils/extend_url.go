package utils

import (
	"fmt"
	"net/url"
	"path"
)

// ExtendURL extends urls, such as ExtendURL("https://example.com/api", "path/to/users")
func ExtendURL(baseURL string, appendPath string) (*string, error) {
	fullURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse base URL: %w", err)
	}

	fullURL.Path = path.Join(fullURL.Path, appendPath)
	output := fullURL.String()

	return &output, nil
}
