package jsonrequire

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

var arrayKeyRegex = regexp.MustCompile(`^([a-zA-Z0-9_-]+)\[(\d+)]$`)

func getByXPath(t testing.TB, data map[string]any, path string) any {
	t.Helper()
	keys := strings.Split(path, "/")
	current := any(data)

	for _, key := range keys {
		if key == "" {
			continue
		}
		// Check if the key is accessing an array element (e.g., "key[0]")
		if matches := arrayKeyRegex.FindStringSubmatch(key); matches != nil {
			mapKey := matches[1]
			index, err := strconv.Atoi(matches[2])
			if err != nil {
				failOnGettingXpath(t, fmt.Sprintf("invalid array index '%s'", matches[2]), data, path)
			}

			m, ok := current.(map[string]any)
			if !ok {
				failOnGettingXpath(t, "path does not exist or is not a map", data, path)
			}

			array, exists := m[mapKey]
			if !exists {
				failOnGettingXpath(t, fmt.Sprintf("key '%s' not found", mapKey), data, path)
			}

			slice, ok := array.([]any)
			if !ok {
				failOnGettingXpath(t, fmt.Sprintf("key '%s' is not an array", mapKey), data, path)
			}

			if index < 0 || index >= len(slice) {
				failOnGettingXpath(t, fmt.Sprintf("index '%d' out of bounds for key '%s'", index, mapKey), data, path)
			}

			// Move to the array element
			current = slice[index]
		} else {
			// Normal map traversal
			m, ok := current.(map[string]any)
			if !ok {
				failOnGettingXpath(t, "path does not exist or is not a map", data, path)
			}

			val, exists := m[key]
			if !exists {
				failOnGettingXpath(t, fmt.Sprintf("key '%s' not found", key), data, path)
			}
			current = val
		}
	}

	return current
}

func failOnGettingXpath(t testing.TB, failureMessage string, data map[string]any, path string) {
	t.Helper()
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		require.Failf(t, failureMessage, "Failed to get path %s from response %+v", path, data)
	}
	require.Fail(t, failureMessage, "Failed to get path %s from response %+v", path, string(jsonData))
}
