package campusgroups

import "googlemaps.github.io/maps"

type EventsResponse struct {
	Events []Event `json:"events"`
}

type Event struct {
	EventDateStr         string `json:"eventDateStr"`
	EventEndDateStr      string `json:"eventEndDateStr"`
	EventDate            string `json:"eventDate"`
	Title                string `json:"title"`
	StartTime            string `json:"startTime"`
	EndTime              string `json:"endTime"`
	EventAddress         string `json:"event_address"`
	EventLocation        string `json:"event_location"`
	EventDescription     string `json:"eventDescription"`
	GroupName            string `json:"groupName"`
	GroupURL             string `json:"groupURL"`
	ClubLogo             string `json:"clubLogo"`
	EventFlyer           string `json:"eventFlyer"`
	AuthorID             string `json:"author_id"`
	ID                   int    `json:"id"`
	ClubID               int    `json:"club_id"`
	RegistrationRequired string `json:"registrationRequired"`
	// HasTickets           bool        `json:"hasTickets"`
	// TicketsSold          interface{} `json:"ticketsSold"` // Can be null or a string
	// SoldTickets          string      `json:"soldTickets"`
	Attendees           []Attendee `json:"attendees"`
	IsAlreadyRegistered bool       `json:"isAlreadyRegistered"`
}

type Attendee struct {
	AttendeesPhotosSrc string `json:"attendeesPhotosSrc"`
	AttendeeUID        string `json:"attendee_uid"`
}

type ProcessedEvent struct {
	Event    Event
	ImageURL *string
	Location *maps.LatLng
}
