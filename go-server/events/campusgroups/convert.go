package campusgroups

import (
	"fmt"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/benkoppe/bear-trak-backend/go-server/api"
	"github.com/benkoppe/bear-trak-backend/go-server/utils/timeutils"
)

func convertAndSort(events []ProcessedEvent) ([]api.Event, error) {
	converted := make([]api.Event, len(events))
	for i, event := range events {
		convEvent, err := convertEvent(event)
		if err != nil {
			return nil, fmt.Errorf("failed to convert event %d: %w", event.Event.ID, err)
		}
		converted[i] = *convEvent
	}

	sort.Slice(converted, func(i, j int) bool {
		return converted[i].Hours.Start.Before(converted[j].Hours.Start)
	})

	return converted, nil
}

func convertEvent(e ProcessedEvent) (*api.Event, error) {
	groupURL, err := url.Parse(e.Event.GroupURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse group URL %s: %w", e.Event.GroupURL, err)
	}

	eventHours, err := convertEventTimes(e)
	if err != nil {
		return nil, fmt.Errorf("failed to convert event times: %w", err)
	}

	pe := api.Event{
		ID:           e.Event.ID,
		Title:        e.Event.Title,
		Description:  e.Event.EventDescription,
		Hours:        *eventHours,
		ImageURL:     api.NilString{Null: true},
		LocationName: api.NewNilString(e.Event.EventLocation),
		Latitude:     api.NilFloat64{Null: true},
		Longitude:    api.NilFloat64{Null: true},
		Group: api.EventGroup{
			ID:   e.Event.ClubID,
			Name: e.Event.GroupName,
			URL:  *groupURL,
		},
	}

	if e.ImageURL != nil {
		pe.ImageURL = api.NewNilString(*e.ImageURL)
	}

	if e.Location != nil {
		pe.Latitude = api.NewNilFloat64(e.Location.Lat)
		pe.Longitude = api.NewNilFloat64(e.Location.Lng)
	}

	return &pe, nil
}

func convertEventTimes(e ProcessedEvent) (*api.Hours, error) {
	// Parse start time
	startTime, err := parseEventDateTime(e.Event.EventDateStr, e.Event.StartTime)
	if err != nil {
		return nil, fmt.Errorf("failed to parse start time (%s %s): %w",
			e.Event.EventDateStr, e.Event.StartTime, err)
	}

	// Parse end time
	endTime, err := parseEventDateTime(e.Event.EventEndDateStr, e.Event.EndTime)
	if err != nil {
		return nil, fmt.Errorf("failed to parse end time (%s %s): %w",
			e.Event.EventEndDateStr, e.Event.EndTime, err)
	}

	return &api.Hours{
		Start: startTime,
		End:   endTime,
	}, nil
}

func parseEventDateTime(dateStr, timeStr string) (time.Time, error) {
	cleanTimeStr := cleanTimeString(timeStr)

	dateTimeStr := fmt.Sprintf("%s %s", dateStr, cleanTimeStr)

	// Event-specific layouts
	eventLayouts := []string{
		"2006-01-02 3:04pm", // 2025-09-22 8:00am
		"2006-01-02 3:04PM", // 2025-09-22 8:00AM
		"2006-01-02 15:04",  // 2025-09-22 08:00
	}

	return timeutils.ParseDateTime(dateTimeStr, eventLayouts)
}

// cleanTimeString removes timezone information from time strings
func cleanTimeString(timeStr string) string {
	// Remove timezone info like "EDT (GMT-4)"
	if idx := strings.Index(timeStr, " EDT"); idx != -1 {
		return timeStr[:idx]
	}
	if idx := strings.Index(timeStr, " EST"); idx != -1 {
		return timeStr[:idx]
	}
	if idx := strings.Index(timeStr, " PDT"); idx != -1 {
		return timeStr[:idx]
	}
	if idx := strings.Index(timeStr, " PST"); idx != -1 {
		return timeStr[:idx]
	}
	// Add other timezone abbreviations as needed
	return timeStr
}
