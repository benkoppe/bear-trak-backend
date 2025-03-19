package scrape

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/benkoppe/bear-trak-backend/go-server/gyms/cornell/static"
	"github.com/benkoppe/bear-trak-backend/go-server/utils"
	"github.com/benkoppe/bear-trak-backend/go-server/utils/time_utils"
)

// parseCaption will try multiple patterns to extract title, start and/or end dates.
func parseCaption(caption string) captionData {
	// Pattern 1: Parenthesized date range.
	// E.g.: "Finals Period Fall 2024 Hours (12/16/24 - 12/23/24)"
	reParen := regexp.MustCompile(`^(.*?)\s*\(\s*([^()]+?)\s*-\s*([^()]+?)\s*\)\s*$`)
	if matches := reParen.FindStringSubmatch(caption); len(matches) == 4 {
		title := strings.TrimSpace(matches[1])
		startStr := matches[2]
		endStr := matches[3]

		startDate, err1 := time_utils.ParseCommonDateTimeYearOptional(startStr)
		endDate, err2 := time_utils.ParseCommonDateTimeYearOptional(endStr)

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

		startDate, err1 := time_utils.ParseCommonDateTimeYearOptional(startStr)
		endDate, err2 := time_utils.ParseCommonDateTimeYearOptional(endStr)

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
		startDate, err := time_utils.ParseCommonDateTimeYearOptional(dateStr)
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
		hrs := time_utils.ParseHours(cell, nil)
		for _, day := range days {
			err := schedule.WeekHours.AddHours(day, hrs)
			if err != nil {
				log.Printf("Error adding hours for %s: %v", day, err)
			}
		}
	}
	return schedule
}

func parseSchedule(table utils.TableData) ParsedSchedule {
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
