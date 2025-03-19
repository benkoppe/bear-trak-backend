package utils

import "strings"

func CleanString(s string) string {
	// Replace newline, carriage return, and tab with a space.
	s = strings.ReplaceAll(s, "\n", " ")
	s = strings.ReplaceAll(s, "\r", " ")
	s = strings.ReplaceAll(s, "\t", " ")

	// Split the string into fields (tokens) to remove any extra spaces,
	// then join with a single space.
	return strings.Join(strings.Fields(s), " ")
}
