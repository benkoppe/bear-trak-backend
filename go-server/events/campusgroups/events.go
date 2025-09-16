package campusgroups

import (
	"fmt"
	"net/url"
	"path"
	"time"

	"github.com/benkoppe/bear-trak-backend/go-server/events/campusgroups/login"
	"github.com/benkoppe/bear-trak-backend/go-server/utils"
)

func buildEventsURL(base *url.URL, start, end string) string {
	eventsURL := *base
	eventsURL.Path = path.Join(eventsURL.Path, "mobile_ws", "v17", "mobile_calendar.aspx")

	q := eventsURL.Query()
	q.Set("view", "group")
	q.Set("start_date", start)
	q.Set("end_date", end)
	q.Set("calendarView", "detail")

	q.Set("etl", "82908") // i don't get this parameter, but its necessary for proper filtering "by open event"

	eventsURL.RawQuery = q.Encode()

	return eventsURL.String()
}

func buildLoginURL(base *url.URL) string {
	loginURL := *base
	loginURL.Path = path.Join(loginURL.Path, "login_only")
	return loginURL.String()
}

func fetchEvents(base *url.URL, start time.Time, days int, loginParams login.LoginParams) (*EventsResponse, error) {
	// inclusive, so subtract 1 before adding
	end := start.AddDate(0, 0, days-1)

	startStr := start.Format("2006-01-02")
	endStr := end.Format("2006-01-02")

	eventsURL := buildEventsURL(base, startStr, endStr)
	loginURL := buildLoginURL(base)

	loginCookie, err := login.GetLoginCookie(
		loginURL,
		loginParams,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get campusgroups login cookie: %w", err)
	}

	headers := map[string]string{
		"User-Agent":       "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/127.0 Safari/537.36",
		"Accept":           "application/json,text/javascript,*/*;q=0.01",
		"X-Requested-With": "XMLHttpRequest",
		"Cookie":           fmt.Sprintf("%s=%s", login.SessionCookieName, loginCookie),
	}

	return utils.DoGetRequest[EventsResponse](eventsURL, headers)
}
