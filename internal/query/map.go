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
			exists = true
			path, err := parsePath(qk, key)
			if err != nil {
				exists = false
				continue
			}
			setValueOnPath(result, path, value)
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
	if rawPath == "" {
		return nil, fmt.Errorf("expect %s to be a map but got value", key)
	}
	splitted := strings.Split(rawPath, "]")
	paths := make([]string, 0)
	for _, p := range splitted {
		if p == "" {
			continue
		}
		if strings.HasPrefix(p, "[") {
			p = p[1:]
		} else {
			return nil, fmt.Errorf("invalid access to map key %s", p)
		}
		if p == "" {
			return nil, fmt.Errorf("expect %s to be a map but got array", key)
		}
		paths = append(paths, p)
	}
	return paths, nil
}

// setValueOnPath is an internal function to set value a path on dicts.
func setValueOnPath(dicts map[string]interface{}, paths []string, value []string) {
	nesting := len(paths)
	currentLevel := dicts
	for i, p := range paths {
		if isLast(i, nesting) {
			currentLevel[p] = value[0]
		} else {
			initNestingIfNotExists(currentLevel, p)
			currentLevel = currentLevel[p].(map[string]interface{})
		}
	}
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
