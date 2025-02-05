package jsonrequire

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

// Getter allows to get values from the JSON data.
type Getter struct {
	t    testing.TB
	data map[string]any
}

// NewGetterWithJSON creates a new Getter based on JSON string
func NewGetterWithJSON(t testing.TB, jsonString string) *Getter {
	var data map[string]any
	if err := json.Unmarshal([]byte(jsonString), &data); err != nil {
		require.Fail(t, fmt.Sprintf("Provided value ('%s') is not valid json.\nJSON parsing error: '%s'", jsonString, err.Error()))
	}
	return &Getter{t: t, data: data}
}

// GetString returns a string value from the data.
func (g *Getter) GetString(xpath string) string {
	g.t.Helper()

	value := getByXPath(g.t, g.data, xpath)

	strValue, ok := value.(string)
	if !ok {
		require.Fail(g.t, "Value on xpath %s is not a string, it is %T (%v)", xpath, value, value)
	}

	return strValue
}
