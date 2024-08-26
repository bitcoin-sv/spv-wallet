package query

import (
	"errors"
	"fmt"
	"strings"
)

const MaxNestedMapDepth = 100

// GetMap returns a map, which satisfies conditions.
func GetMap(query map[string][]string, filteredKey string) (dicts map[string]interface{}, err error) {
	result := make(map[string]interface{})
	getAll := filteredKey == ""
	var allErrors = make([]error, 0)
	for key, value := range query {
		kType := getType(key, filteredKey, getAll)
		switch kType {
		case "filtered_unsupported":
			allErrors = append(allErrors, fmt.Errorf("invalid access to map %s", key))
			continue
		case "filtered_map":
			fallthrough
		case "map":
			path, mapErr := parsePath(key)
			if mapErr != nil {
				allErrors = append(allErrors, mapErr)
				continue
			}
			if !getAll {
				path = path[1:]
			}
			mapErr = setValueOnPath(result, path, value)
			if mapErr != nil {
				allErrors = append(allErrors, mapErr)
				continue
			}
		case "array":
			result[keyWithoutArraySymbol(key)] = value
		case "filtered_rejected":
			continue
		default:
			result[key] = value[0]
		}
	}
	if len(allErrors) > 0 {
		return nil, errors.Join(allErrors...)
	} else if len(result) == 0 {
		return nil, nil
	}
	return result, nil
}

// getType is an internal function to get the type of query key.
func getType(qk string, key string, getAll bool) string {
	if getAll {
		if isMap(qk) {
			return "map"
		} else if isArray(qk) {
			return "array"
		}
		return "value"
	}
	if isFilteredKey(qk, key) {
		if isMap(qk) {
			return "filtered_map"
		}
		return "filtered_unsupported"
	}
	return "filtered_rejected"
}

// isFilteredKey is an internal function to check if k is accepted when searching for map with given key.
func isFilteredKey(k string, key string) bool {
	return k == key || strings.HasPrefix(k, key+"[")
}

// isMap is an internal function to check if k is a map query key.
func isMap(k string) bool {
	i := strings.IndexByte(k, '[')
	j := strings.IndexByte(k, ']')
	return j-i > 1 || (i >= 0 && j == -1)
}

// isArray is an internal function to check if k is an array query key.
func isArray(k string) bool {
	i := strings.IndexByte(k, '[')
	j := strings.IndexByte(k, ']')
	return j-i == 1
}

// keyWithoutArraySymbol is an internal function to remove array symbol from query key.
func keyWithoutArraySymbol(qk string) string {
	return qk[:len(qk)-2]
}

// parsePath is an internal function to parse key access path.
// For example, key[foo][bar] will be parsed to ["foo", "bar"].
func parsePath(k string) ([]string, error) {
	firstKeyEnd := strings.IndexByte(k, '[')
	if firstKeyEnd == -1 {
		return nil, fmt.Errorf("invalid access to map key %s", k)
	}
	first, rawPath := k[:firstKeyEnd], k[firstKeyEnd:]

	split := strings.Split(rawPath, "]")

	// Bear in mind that split of the valid map will always have "" as the last element.
	if split[len(split)-1] != "" {
		return nil, fmt.Errorf("invalid access to map key %s", k)
	} else if len(split)-1 > MaxNestedMapDepth {
		return nil, fmt.Errorf("maximum depth [%d] of nesting in map exceeded [%d]", MaxNestedMapDepth, len(split)-1)
	}

	// -2 because after split the last element should be empty string.
	last := len(split) - 2

	paths := []string{first}
	for i := 0; i <= last; i++ {
		p := split[i]

		// this way we can handle both error cases: foo] and [foo[bar
		if strings.LastIndex(p, "[") == 0 {
			p = p[1:]
		} else {
			return nil, fmt.Errorf("invalid access to map key %s", p)
		}
		if p == "" && i != last {
			return nil, fmt.Errorf("unsupported array-like access to map key %s", k)
		}
		paths = append(paths, p)
	}
	return paths, nil
}

// setValueOnPath is an internal function to set value a path on dicts.
func setValueOnPath(dicts map[string]interface{}, paths []string, value []string) error {
	nesting := len(paths)
	previousLevel := dicts
	currentLevel := dicts
	for i, p := range paths {
		if isLast(i, nesting) {
			if isArrayOnPath(p) {
				previousLevel[paths[i-1]] = value
			} else {
				currentLevel[p] = value[0]
			}
		} else {
			initNestingIfNotExists(currentLevel, p)
			previousLevel = currentLevel
			switch currentLevel[p].(type) {
			case map[string]any:
				currentLevel = currentLevel[p].(map[string]any)
			case []string:
				return fmt.Errorf("trying to set array and nested map at the same key [%s]", p)
			case string:
				return fmt.Errorf("trying to set value and nested map at the same key [%s]", p)
			}
		}
	}
	return nil
}

// isArrayOnPath is an internal function to check if the current parsed map path is an array.
func isArrayOnPath(p string) bool {
	return p == ""
}

// initNestingIfNotExists is an internal function to initialize a nested map if not exists.
func initNestingIfNotExists(currentLevel map[string]interface{}, p string) {
	if _, ok := currentLevel[p]; !ok {
		currentLevel[p] = make(map[string]interface{})
	}
}

// isLast is an internal function to check if the current level is the last level.
func isLast(i int, nesting int) bool {
	return i == nesting-1
}
