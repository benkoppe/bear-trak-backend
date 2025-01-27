package static

import (
	"net/url"
)

type Alert struct {
	ID int `json:"id"`

	Title   string `json:"title"`
	Message string `json:"message"`

	Enabled  bool `json:"enabled"`
	ShowOnce bool `json:"showOnce"`

	MaxBuild *int `json:"maxBuild,omitempty"` // alerts won't be displayed on builds after this optional property

	Button *AlertButton `json:"button,omitempty"`
}

type AlertButton struct {
	Title string  `json:"title"`
	URL   URLText `json:"url"`
}

type URLText struct {
	*url.URL
}

func (u *URLText) UnmarshalJSON(data []byte) error {
	s := string(data)
	s = s[1 : len(s)-1] // trim quotes
	parsedURL, err := url.Parse(s)
	if err != nil {
		return err
	}

	u.URL = parsedURL
	return nil
}
