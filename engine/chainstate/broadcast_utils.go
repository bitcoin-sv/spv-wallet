package chainstate

import (
	"strings"
)

// containsAny checks if the given string contains any of the provided substrings
func containsAny(s string, substr []string) bool {
	lower := strings.ToLower(s)
	for _, str := range substr {
		if strings.Contains(lower, str) {
			return true
		}
	}
	return false
}

