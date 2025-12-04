package campusgroups

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// scraping is needed only for event images.
func fetchEventImage(base *url.URL, eventID int) (*string, error) {
	detailURL := buildDetailURL(base, eventID)

	resp, err := http.Get(detailURL)
	if err != nil {
		return nil, fmt.Errorf("error fetching page: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP error: %d", resp.StatusCode)
	}

	return scrapeEventImageData(base, resp.Body)
}

func buildDetailURL(base *url.URL, eventID int) string {
	detailURL := *base
	detailURL.Path = path.Join(detailURL.Path, "einhorn", "rsvp_boot")

	q := detailURL.Query()
	q.Set("id", strconv.Itoa(eventID))
	detailURL.RawQuery = q.Encode()

	return detailURL.String()
}

func scrapeEventImageData(base *url.URL, htmlReader io.Reader) (*string, error) {
	doc, err := goquery.NewDocumentFromReader(htmlReader)
	if err != nil {
		return nil, fmt.Errorf("error parsing HTML: %w", err)
	}

	imageURL := doc.Find("#event_main_card .card-block .row .col-md-8 img.img-responsive").First().AttrOr("src", "")
	if imageURL != "" {
		// Convert relative URLs to absolute if needed
		if strings.HasPrefix(imageURL, "/") {
			builtImageURL := *base
			builtImageURL.Path = path.Join(builtImageURL.Path, imageURL)
			imageURL = builtImageURL.String()
		}
		return &imageURL, nil
	}
	// no image found
	return nil, nil
}
