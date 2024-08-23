package query

import (
	"errors"
	"fmt"
	"strings"
)

// GetMap returns a map, which satisfies conditions.
func GetMap(query map[string][]string, key string) (dicts map[string]interface{}, err error) {
	result := make(map[string]interface{})
	getAll := key == ""
	var allErrors = make([]error, 0)
	for qk, value := range query {
		kType := getType(qk, key, getAll)
		switch kType {
		case "filtered_unsupported":
			allErrors = append(allErrors, fmt.Errorf("invalid access to map %s", qk))
			continue
		case "filtered_map":
			fallthrough
		case "map":
			path, err := parsePath(qk)
			if err != nil {
				allErrors = append(allErrors, err)
				continue
			}
			if !getAll {
				path = path[1:]
			}
			setValueOnPath(result, path, value)
		case "array":
			result[keyWithoutArraySymbol(qk)] = value
		case "filtered_rejected":
			continue
		default:
			result[qk] = value[0]
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
	return j-i > 1
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
	first, rawPath := k[:firstKeyEnd], k[firstKeyEnd:]

	split := strings.Split(rawPath, "]")
	if split[len(split)-1] != "" {
		return nil, fmt.Errorf("invalid access to map key %s", k)
	}

	// -2 because after split the last element should be empty string.
	last := len(split) - 2

	paths := []string{first}
	for i := 0; i <= last; i++ {
		p := split[i]
		if strings.HasPrefix(p, "[") {
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
func setValueOnPath(dicts map[string]interface{}, paths []string, value []string) {
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
			currentLevel = currentLevel[p].(map[string]interface{})
		}
	}
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
