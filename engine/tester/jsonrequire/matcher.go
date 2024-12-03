package jsonrequire

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"testing"
	"text/template"
)

// Match helps to make assertions on JSON strings when some values are not known in advance.
// For example, when we do the assertion on JSON serialized models, we can't predict the values of fields like IDs or timestamps.
// In such cases, we can use a "template" with placeholders for these values.
// See matcher_test.go for examples
func Match(t testing.TB, expectedTemplateFormat string, params map[string]any, actual string) {
	t.Helper()
	expected := compileTemplate(t, expectedTemplateFormat, params)
	assertJSONWithPlaceholders(t, expected, actual)
}

func compileTemplate(t testing.TB, templateFormat string, params map[string]any) string {
	t.Helper()
	tmpl, err := template.
		New("").
		Funcs(funcsMap).
		Parse(templateFormat)

	if err != nil {
		require.Fail(t, "Failed to parse template", err)
	}
	var expected bytes.Buffer
	if err = tmpl.Execute(&expected, params); err != nil {
		require.Fail(t, "Failed to execute template", err)
	}
	return expected.String()
}
