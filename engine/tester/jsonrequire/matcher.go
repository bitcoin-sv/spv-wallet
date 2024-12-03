package jsonrequire

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"testing"
	"text/template"
)

func Match(t testing.TB, expectedTemplateFormat string, params map[string]any, actual string) {
	t.Helper()
	expected := compileTemplate(t, expectedTemplateFormat, params)
	assertJSONWithPlaceholders(t, expected, actual)
}

func compileTemplate(t testing.TB, templateFormat string, params map[string]any) string {
	t.Helper()
	tmpl, err := template.New("").Parse(templateFormat)
	if err != nil {
		require.Fail(t, "Failed to parse template", err)
	}
	var expected bytes.Buffer
	if err = tmpl.Execute(&expected, params); err != nil {
		require.Fail(t, "Failed to execute template", err)
	}
	return expected.String()
}
