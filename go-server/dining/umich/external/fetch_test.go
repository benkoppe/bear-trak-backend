package external

import (
	"net/url"
	"testing"
	"time"
)

func TestBuildRequestURL(t *testing.T) {
	day := time.Date(2026, time.February, 11, 0, 0, 0, 0, time.UTC)

	params := url.Values{
		"key":      {"test-key"},
		"location": {"Berts Cafe"},
		"date":     {day.Format("02-01-2006")},
		"meal":     {"LUNCH"},
	}

	requestURL, err := buildRequestURL("https://example.com/dining", "menu", params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	parsed, err := url.Parse(requestURL)
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}

	if parsed.Path != "/dining/menu" {
		t.Fatalf("unexpected path: %s", parsed.Path)
	}
	if parsed.Query().Get("date") != "11-02-2026" {
		t.Fatalf("unexpected date param: %s", parsed.Query().Get("date"))
	}
	if parsed.Query().Get("meal") != "LUNCH" {
		t.Fatalf("unexpected meal param: %s", parsed.Query().Get("meal"))
	}
}
