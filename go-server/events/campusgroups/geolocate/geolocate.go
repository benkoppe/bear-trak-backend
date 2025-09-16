// Package geolocate takes in a string address and returns its location
package geolocate

import (
	"context"
	"fmt"

	"googlemaps.github.io/maps"
)

// specialLocations are strings we don’t need to geocode.
var specialLocations = map[string]*maps.LatLng{
	"TBD":                                    nil,
	"Private Location (register to display)": nil,
	"Online Event":                           nil,
	"Willard Straight Theatre, 104 Willard Straight Hall": {Lat: 42.44648867509554, Lng: -76.48567704394083},
}

func (g *GeoLocator) Geolocate(location string) (*maps.LatLng, error) {
	// Handle special cases first
	if loc, ok := handleSpecialLocation(location); ok {
		return loc, nil
	}

	// Obtain maps client
	mapsC, err := g.clientOrInit()
	if err != nil {
		return nil, fmt.Errorf("failed to init maps client: %w", err)
	}

	// Try geocoding
	loc, err := geocodeLocation(mapsC, location, g.PreferBounds)
	if err != nil {
		return nil, fmt.Errorf("geocoding failed for %q: %w", location, err)
	}
	if loc != nil {
		return loc, nil
	}

	// Fallback to Places API
	// loc, err = placesLocation(mapsC, location)
	// if err != nil {
	// 	log.Printf("places search failed for %q: %v\n", location, err)
	// 	return nil, err
	// }
	// if loc == nil {
	// 	log.Printf("\tno location found for %s in geocode or places.\n", location)
	// }

	return loc, nil
}

func handleSpecialLocation(location string) (*maps.LatLng, bool) {
	if loc, ok := specialLocations[location]; ok {
		return loc, true
	}
	return nil, false
}

// geocodeLocation calls the Google Maps Geocoding API to resolve a location string.
func geocodeLocation(mapsC *maps.Client, location string, preferBounds *maps.LatLngBounds) (*maps.LatLng, error) {
	r := &maps.GeocodingRequest{
		Address: location,
		Bounds:  preferBounds,
	}

	response, err := mapsC.Geocode(context.Background(), r)
	if err != nil {
		return nil, fmt.Errorf("failed to geocode location %s: %w", location, err)
	}

	if len(response) == 0 {
		return nil, nil
	}

	loc := response[0].Geometry.Location
	return &loc, nil
}

// placesLocation calls the Google Places API to resolve a string into a Lat/Lng.
// CURRENTLY DISABLED: this requires enabling billing.

// func placesLocation(mapsC *maps.Client, location string) (*maps.LatLng, error) {
// 	r := &maps.FindPlaceFromTextRequest{
// 		Input:                 location,
// 		InputType:             maps.FindPlaceFromTextInputTypeTextQuery,
// 		Fields:                []maps.PlaceSearchFieldMask{maps.PlaceSearchFieldMaskGeometry},
// 		LocationBiasNorthEast: &maps.LatLng{Lat: 42.470, Lng: -76.458},
// 		LocationBiasSouthWest: &maps.LatLng{Lat: 42.416, Lng: -76.541},
// 	}
//
// 	resp, err := mapsC.FindPlaceFromText(context.Background(), r)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to call Places API for %s: %w", location, err)
// 	}
//
// 	if len(resp.Candidates) == 0 {
// 		return nil, nil
// 	}
//
// 	// Take the first candidate
// 	loc := resp.Candidates[0].Geometry.Location
// 	return &loc, nil
// }
