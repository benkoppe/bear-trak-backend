package campusgroups

import (
	"os"
	"testing"
	"time"

	"github.com/benkoppe/bear-trak-backend/go-server/events/campusgroups/geolocate"
	"github.com/benkoppe/bear-trak-backend/go-server/events/campusgroups/login"
	"googlemaps.github.io/maps"
)

func TestAll(t *testing.T) {
	loginParams := login.LoginParams{
		LoginEmail:       os.Getenv("LOGIN_EMAIL"),
		OtpEmail:         os.Getenv("OTP_EMAIL"),
		OtpEmailPassword: os.Getenv("OTP_EMAIL_PASSWORD"),
	}
	locator := &geolocate.GeoLocator{
		APIKey: os.Getenv("API_KEY"),
		PreferBounds: &maps.LatLngBounds{
			NorthEast: maps.LatLng{Lat: 42.470, Lng: -76.458},
			SouthWest: maps.LatLng{Lat: 42.416, Lng: -76.541},
		},
	}

	events, err := fetchAllData("https://cornell.campusgroups.com", loginParams, locator)
	if err != nil {
		t.Fatalf("fetchAllData failed: %v", err)
	}

	for _, e := range events {
		t.Logf(
			"Event: %s (%s@%s - %s@%s) at %s | Image: %v | Location: %v",
			e.Event.Title,
			e.Event.EventDateStr,
			e.Event.StartTime,
			e.Event.EventEndDateStr,
			e.Event.EndTime,
			e.Event.EventLocation,
			e.ImageURL,
			e.Location,
		)
		t.Logf("\tgroup #%d - %s (%s)", e.Event.ClubID, e.Event.GroupName, e.Event.GroupURL)
	}

	t.Logf("\nPOST CONVERSION: -------------------------------------------\n")

	converted, err := convertAndSort(events)
	if err != nil {
		t.Fatalf("convertAndSort failed: %v", err)
	}

	for _, e := range converted {
		t.Logf(
			"Converted: #%d %s | %s → %s | Group: %s | Lat: %v, Lng: %v | Location: %s | Image: %s",
			e.ID,
			e.Title,
			e.Hours.Start.Format(time.RFC3339),
			e.Hours.End.Format(time.RFC3339),
			e.Group.Name,
			e.Latitude,
			e.Longitude,
			e.LocationName.Value,
			e.ImageURL.Value,
		)
	}
}
