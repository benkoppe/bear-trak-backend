// Package predictions loads predictions data from github csv
package predictions

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

func fetchData(url string) ([]Prediction, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch URL %q: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected HTTP status %d when fetching %q", resp.StatusCode, url)
	}

	reader := csv.NewReader(resp.Body)
	reader.FieldsPerRecord = -1

	// skip header row
	if _, err := reader.Read(); err != nil {
		return nil, fmt.Errorf("failed to read CSV header: %w", err)
	}

	var predictions []Prediction
	line := 2 // header was line 1
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("line %d: failed to read CSV record: %w", line, err)
		}

		p, err := parsePrediction(record, line)
		if err != nil {
			return nil, err
		}

		predictions = append(predictions, p)
		line++
	}

	return predictions, nil
}

func parsePrediction(record []string, line int) (Prediction, error) {
	if len(record) < 4 {
		return Prediction{}, fmt.Errorf("line %d: invalid record length (expected 4 fields, got %d): %v", line, len(record), record)
	}

	ts, err := time.Parse("2006-01-02 15:04:05-07:00", record[1])
	if err != nil {
		return Prediction{}, fmt.Errorf("line %d: failed to parse Timestamp %q: %w", line, record[1], err)
	}

	pred, err := strconv.Atoi(record[2])
	if err != nil {
		return Prediction{}, fmt.Errorf("line %d: failed to parse Predicted value %q: %w", line, record[2], err)
	}

	pm, err := time.Parse("2006-01-02 15:04:05.999999999-07:00", record[3])
	if err != nil {
		return Prediction{}, fmt.Errorf("line %d: failed to parse PredictionMadeAt %q: %w", line, record[3], err)
	}

	return Prediction{
		GymName:        record[0],
		Timestamp:      ts,
		Predicted:      pred,
		PredictionMade: pm,
	}, nil
}
