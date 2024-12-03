package jsonrequire

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

var arrayKeyRegex = regexp.MustCompile(`^([a-zA-Z0-9_-]+)\[(\d+)]$`)

func getByXPath(t testing.TB, data map[string]any, path string) any {
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
				require.Fail(t, fmt.Sprintf("invalid array index '%s'", matches[2]))
			}

			m, ok := current.(map[string]any)
			if !ok {
				require.Fail(t, "path does not exist or is not a map")
			}

			array, exists := m[mapKey]
			if !exists {
				require.Fail(t, fmt.Sprintf("key '%s' not found", mapKey))
			}

			slice, ok := array.([]any)
			if !ok {
				require.Fail(t, fmt.Sprintf("key '%s' is not an array", mapKey))
			}

			if index < 0 || index >= len(slice) {
				require.Fail(t, fmt.Sprintf("index '%d' out of bounds for key '%s'", index, mapKey))
			}

			// Move to the array element
			current = slice[index]
		} else {
			// Normal map traversal
			m, ok := current.(map[string]any)
			if !ok {
				require.Fail(t, "path does not exist or is not a map")
			}

			val, exists := m[key]
			if !exists {
				require.Fail(t, fmt.Sprintf("key '%s' not found", key))
			}
			current = val
		}
	}

	return current
}
