package scrape

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/benkoppe/bear-trak-backend/go-server/gyms/static"
)

// tries a few common date layouts
func parseDate(s string) (time.Time, error) {
	// Remove any weekday prefix, for example "Tues, Jan 21" -> "Jan 21"
	if strings.Contains(s, ",") {
		parts := strings.SplitN(s, ",", 2)
		s = strings.TrimSpace(parts[1])
	} else {
		s = strings.TrimSpace(s)
	}

	// Try several layouts
	layouts := []string{
		"Jan 2",
		"1/2/06",
		"01/02/06",
	}

	var err error
	for _, layout := range layouts {
		if t, err2 := time.Parse(layout, s); err2 == nil {
			return t, nil
		} else {
			err = err2
		}
	}
	return time.Time{}, fmt.Errorf("unable to parse date %q: %v", s, err)
}

// parseCaption will try multiple patterns to extract title, start and/or end dates.
func parseCaption(caption string) captionData {
	// Pattern 1: Parenthesized date range.
	// E.g.: "Finals Period Fall 2024 Hours (12/16/24 - 12/23/24)"
	reParen := regexp.MustCompile(`^(.*?)\s*\(\s*([^()]+?)\s*-\s*([^()]+?)\s*\)\s*$`)
	if matches := reParen.FindStringSubmatch(caption); len(matches) == 4 {
		title := strings.TrimSpace(matches[1])
		startStr := matches[2]
		endStr := matches[3]

		startDate, err1 := parseDate(startStr)
		endDate, err2 := parseDate(endStr)

		var startPtr, endPtr *time.Time
		if err1 == nil && err2 == nil {
			startPtr = &startDate
			endPtr = &endDate
		}
		return captionData{
			Title:     title,
			StartDate: startPtr,
			EndDate:   endPtr,
		}
	}

	// Pattern 2: Colon-based date range.
	// E.g.: "February Break: Feb 17 - Feb 23"
	reColon := regexp.MustCompile(`^(.*):\s*([A-Za-z]+\s+\d{1,2})\s*-\s*([A-Za-z]+\s+\d{1,2})\s*$`)
	if matches := reColon.FindStringSubmatch(caption); len(matches) == 4 {
		title := strings.TrimSpace(matches[1])
		startStr := matches[2]
		endStr := matches[3]

		startDate, err1 := parseDate(startStr)
		endDate, err2 := parseDate(endStr)

		var startPtr, endPtr *time.Time
		if err1 == nil && err2 == nil {
			startPtr = &startDate
			endPtr = &endDate
		}
		return captionData{
			Title:     title,
			StartDate: startPtr,
			EndDate:   endPtr,
		}
	}

	// Pattern 3: "starting" a single date.
	// E.g.: "Regular Semester Hours starting Tues, Jan 21"
	reStarting := regexp.MustCompile(`^(.*?)\s+starting\s+(.+)$`)
	if matches := reStarting.FindStringSubmatch(caption); len(matches) == 3 {
		title := strings.TrimSpace(matches[1])
		dateStr := matches[2]
		startDate, err := parseDate(dateStr)
		var startPtr *time.Time
		if err == nil {
			startPtr = &startDate
		}
		return captionData{
			Title:     title,
			StartDate: startPtr,
		}
	}

	// If nothing matched, treat the whole caption as the title.
	return captionData{Title: caption}
}

var daysOfWeek = []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}

func extractDayName(s string) []string {
	s = strings.ToLower(strings.TrimSpace(s))

	if strings.Contains(s, "weekend") {
		return []string{"Saturday", "Sunday"}
	}

	for _, day := range daysOfWeek {
		if strings.Contains(s, strings.ToLower(day)) {
			return []string{day}
		}
	}

	return nil
}

// parseHeaderDays maps headers like "Monday - Thursday", "Weekend", or "Friday" to actual day names.
func parseHeaderDays(header string) []string {
	if strings.Contains(header, "&") {
		parts := strings.Split(header, "&")
		var days []string
		for _, p := range parts {
			days = append(days, extractDayName(p)...)
		}
		return days
	}
	if strings.Contains(header, "-") {
		parts := strings.Split(header, "-")
		if len(parts) == 2 {
			startDay := extractDayName(parts[0])
			endDay := extractDayName(parts[1])
			if len(startDay) != 1 || len(endDay) != 1 {
				fmt.Printf("Couldn't parse days from header: %s", header)
				return nil
			}
			daysOfWeek := []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}
			startIndex, endIndex := -1, -1
			for i, d := range daysOfWeek {
				if d == startDay[0] {
					startIndex = i
				}
				if d == endDay[0] {
					endIndex = i
				}
			}
			if startIndex != -1 && endIndex != -1 && startIndex <= endIndex {
				return daysOfWeek[startIndex : endIndex+1]
			}
		}
	}
	return extractDayName(header)
}

// parseCellHours splits a cell value into one or more Hours.
// It handles values like "6am - 9pm" or "7am - 8:30am / 10am - 10:45pm".
func parseCellHours(cell string) []static.Hours {
	var hrs []static.Hours

	if strings.Contains(cell, "Closed") {
		return hrs
	}

	ranges := strings.Split(cell, "/")
	for _, r := range ranges {
		r = strings.TrimSpace(r)
		// split by dash for open/close
		parts := strings.Split(r, "-")
		if len(parts) != 2 {
			log.Printf("Couldn't format a static gym cell: %s", r)
			continue // doesn't match expected pattern
		}
		open := strings.TrimSpace(parts[0])
		close := strings.TrimSpace(parts[1])
		hrs = append(hrs, static.Hours{
			Open:  static.TimeString(open),
			Close: static.TimeString(close),
		})
	}
	return hrs
}

// parseGymSchedule takes a headers slice and corresponding values slice and builds a GymSchedule.
// Assumes the first header/value is the gym name.
func parseGymSchedule(headers, values []string) GymSchedule {
	schedule := GymSchedule{
		GymName:   values[0],
		WeekHours: static.WeekHours{},
	}
	// process remaining headers
	for i := 1; i < len(headers) && i < len(values); i++ {
		header := headers[i]
		cell := values[i]

		days := parseHeaderDays(header)
		hrs := parseCellHours(cell)
		for _, day := range days {
			err := schedule.WeekHours.AddHours(day, hrs)
			if err != nil {
				log.Printf("Error adding hours for %s: %v", day, err)
			}
		}
	}
	return schedule
}

func parseSchedule(table tableData) ParsedSchedule {
	schedule := ParsedSchedule{
		Caption:      parseCaption(table.Caption),
		GymSchedules: make([]GymSchedule, 0),
	}

	for _, row := range table.Rows {
		gymSchedule := parseGymSchedule(table.Headers, row.Columns)
		schedule.GymSchedules = append(schedule.GymSchedules, gymSchedule)
	}

	return schedule
}
