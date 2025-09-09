package shared

import "time"

func FilterByTimeRange[T any](
	items []T,
	getTime func(T) time.Time,
	start, end *time.Time,
) []T {
	if start == nil || end == nil {
		return items
	}

	var filtered []T
	for _, item := range items {
		t := getTime(item)
		if !t.Before(*start) && !t.After(*end) {
			filtered = append(filtered, item)
		}
	}
	return filtered
}
