// Package campusgroups loads event data from campusgroups
package campusgroups

import (
	"fmt"
	"log"
	"net/url"
	"sync"
	"time"

	"github.com/benkoppe/bear-trak-backend/go-server/events/campusgroups/geolocate"
	"github.com/benkoppe/bear-trak-backend/go-server/events/campusgroups/login"
)

func fetchAllData(baseURL string, loginParams login.LoginParams, locator *geolocate.GeoLocator) ([]ProcessedEvent, error) {
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse base URL: %w", err)
	}

	events, err := fetchEvents(
		parsedURL,
		time.Now(),
		7,
		loginParams,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch campusgroups events: %w", err)
	}

	const maxWorkers = 20
	sem := make(chan struct{}, maxWorkers)
	results := make(chan ProcessedEvent, len(events.Events))
	var wg sync.WaitGroup

	for _, e := range events.Events {
		wg.Add(1)
		sem <- struct{}{}
		go func(e Event) {
			defer wg.Done()
			processed := processEvent(parsedURL, e, locator)
			results <- processed
			<-sem
		}(e)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	var processedEvents []ProcessedEvent
	for pe := range results {
		processedEvents = append(processedEvents, pe)
	}

	return processedEvents, nil
}

func processEvent(base *url.URL, e Event, locator *geolocate.GeoLocator) ProcessedEvent {
	imageURL, err := fetchEventImage(base, e.ID)
	if err != nil {
		log.Printf("failed to fetch campusgroups image for event %d: %v\n", e.ID, err)
	}

	loc, err := locator.Geolocate(e.EventLocation)
	if err != nil {
		log.Printf("failed to fetch location lat/lng: %v", err)
	}

	return ProcessedEvent{
		Event:    e,
		ImageURL: imageURL,
		Location: loc,
	}
}
