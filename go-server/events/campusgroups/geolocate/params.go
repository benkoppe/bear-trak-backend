package geolocate

import (
	"sync"

	"googlemaps.github.io/maps"
)

type GeoLocator struct {
	APIKey       string
	PreferBounds *maps.LatLngBounds

	mu     sync.Mutex
	client *maps.Client
}

// clientOrInit lazily initializes the maps.Client if needed and reuses it afterwards.
func (g *GeoLocator) clientOrInit() (*maps.Client, error) {
	g.mu.Lock()
	defer g.mu.Unlock()

	if g.client != nil {
		return g.client, nil
	}

	c, err := maps.NewClient(maps.WithAPIKey(g.APIKey))
	if err != nil {
		return nil, err
	}

	g.client = c
	return g.client, nil
}
