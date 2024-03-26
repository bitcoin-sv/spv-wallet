package utils

import "regexp"

func SanitizeInput(input string) string {
	// Remove non-alphanumeric characters using regular expression
	sanitized := regexp.MustCompile("[^a-zA-Z0-9]").ReplaceAllString(input, "")
	return sanitized
}
