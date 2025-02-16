package utils

import (
	"fmt"
	"net/url"
	"path"
)

// extends urls, such as ExtendURL("https://example.com/api", "path/to/users")
func ExtendUrl(baseUrl string, appendPath string) (*string, error) {
	fullUrl, err := url.Parse(baseUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to parse base URL: %w", err)
	}

	fullUrl.Path = path.Join(fullUrl.Path, appendPath)
	output := fullUrl.String()

	return &output, nil
}
