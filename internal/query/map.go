package query

import (
	"fmt"
	"strings"
)

// GetMap returns a map, which satisfies conditions.
func GetMap(query map[string][]string, key string) (dicts map[string]interface{}, exists bool) {
	result := make(map[string]interface{})
	for qk, value := range query {
		if isKey(qk, key) {
			path, err := parsePath(qk, key)
			if err != nil {
				continue
			}
			setValueOnPath(result, path, value)
			exists = true
		}
	}
	if !exists {
		return nil, exists
	}
	return result, exists

}

// isKey is an internal function to check if a k is a map key.
func isKey(k string, key string) bool {
	i := strings.IndexByte(k, '[')
	return i >= 1 && k[0:i] == key
}

// parsePath is an internal function to parse key access path.
// For example, key[foo][bar] will be parsed to ["foo", "bar"].
func parsePath(k string, key string) ([]string, error) {
	rawPath := strings.TrimPrefix(k, key)
	splitted := strings.Split(rawPath, "]")
	paths := make([]string, 0)
	for i, p := range splitted {
		if p == "" {
			continue
		}
		if strings.HasPrefix(p, "[") {
			p = p[1:]
		} else {
			return nil, fmt.Errorf("invalid access to map key %s", p)
		}
		if i == 0 && p == "" {
			return nil, fmt.Errorf("expect %s to be a map but got array", key)
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
			if isArray(p) {
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

func isArray(p string) bool {
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
