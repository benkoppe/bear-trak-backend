package predictions

import "testing"

func TestParsePredictionUsesNewYorkLocalTime(t *testing.T) {
	record := []string{
		"Helen Newman Fitness Center",
		"2026-03-01 10:30:00",
		"48",
		"2026-03-01 10:33:27.333743",
	}

	prediction, err := parsePrediction(record, 2)
	if err != nil {
		t.Fatalf("parsePrediction returned error: %v", err)
	}

	if got, want := prediction.Timestamp.Format("2006-01-02 15:04:05 -0700"), "2026-03-01 10:30:00 -0500"; got != want {
		t.Fatalf("unexpected timestamp: got %q want %q", got, want)
	}

	if got, want := prediction.PredictionMade.Format("2006-01-02 15:04:05.999999 -0700"), "2026-03-01 10:33:27.333743 -0500"; got != want {
		t.Fatalf("unexpected prediction made time: got %q want %q", got, want)
	}
}

func TestParsePredictionUsesDSTOffsetWhenApplicable(t *testing.T) {
	record := []string{
		"Helen Newman Fitness Center",
		"2026-07-01 10:30:00",
		"48",
		"2026-07-01 10:33:27.333743",
	}

	prediction, err := parsePrediction(record, 2)
	if err != nil {
		t.Fatalf("parsePrediction returned error: %v", err)
	}

	if got, want := prediction.Timestamp.Format("2006-01-02 15:04:05 -0700"), "2026-07-01 10:30:00 -0400"; got != want {
		t.Fatalf("unexpected timestamp: got %q want %q", got, want)
	}

	if got, want := prediction.PredictionMade.Format("2006-01-02 15:04:05.999999 -0700"), "2026-07-01 10:33:27.333743 -0400"; got != want {
		t.Fatalf("unexpected prediction made time: got %q want %q", got, want)
	}
}
