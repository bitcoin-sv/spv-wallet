package jsonrequire

import (
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func assertJSONWithPlaceholders(t testing.TB, expected, actual string) {
	t.Helper()

	var expectedJSONValue, actualJSONValue map[string]any
	if err := json.Unmarshal([]byte(expected), &expectedJSONValue); err != nil {
		require.Fail(t, fmt.Sprintf("Expected value ('%s') is not valid json.\nJSON parsing error: '%s'", expected, err.Error()))
	}

	if err := json.Unmarshal([]byte(actual), &actualJSONValue); err != nil {
		require.Fail(t, fmt.Sprintf("Input value ('%s') is not valid json.\nJSON parsing error: '%s'", actual, err.Error()))
	}

	processJSONPlaceholders(t, expectedJSONValue, actualJSONValue, "")

	if assert.ObjectsAreEqual(expectedJSONValue, actualJSONValue) {
		return
	}

	// We want a message that shows the diff between the two JSON strings not the decoded objects.
	expectedJsonString := expected
	actualJSONString := actual

	// Try to unify the JSON strings to make the diff more readable.
	marshaled, err := json.MarshalIndent(expectedJSONValue, "", "  ")
	if err == nil {
		expectedJsonString = string(marshaled)
	}

	marshaled, err = json.MarshalIndent(actualJSONValue, "", "  ")
	if err == nil {
		actualJSONString = string(marshaled)
	}

	assert.Equal(t, expectedJsonString, actualJSONString)
}

func processJSONPlaceholders(t testing.TB, template any, base map[string]any, xpath string) {
	t.Helper()

	rval, kind, _ := typeof(template)

	if kind == reflect.Map {
		for _, k := range rval.MapKeys() {
			currentXPath := addToXPath(xpath, k.String())
			_, templateKind, templateVal := typeof(rval.MapIndex(k).Interface())

			if templateKind == reflect.String {
				str := templateVal.(string)
				toCopy := processTemplateCandidate(t, str, base, currentXPath)
				if toCopy == nil {
					continue
				}

				// copy the value from the base to the template
				rval.SetMapIndex(k, *toCopy)
			} else {
				processJSONPlaceholders(t, templateVal, base, currentXPath)
			}
		}
	} else if kind == reflect.Slice {
		for i := 0; i < rval.Len(); i++ {
			currentXPath := xpath + fmt.Sprintf("[%d]", i)
			value := rval.Index(i)

			_, itemKind, itemVal := typeof(value.Interface())
			if itemKind == reflect.String {
				str := itemVal.(string)
				toCopy := processTemplateCandidate(t, str, base, currentXPath)
				if toCopy != nil {
					rval.Index(i).Set(*toCopy)
				}
			} else {
				processJSONPlaceholders(t, value.Interface(), base, currentXPath)
			}
		}
	}
}

func processTemplateCandidate(t testing.TB, templateVal string, base map[string]any, xpath string) *reflect.Value {
	t.Helper()
	isRegex := containsRegex(templateVal)

	if !isRegex && templateVal != "*" {
		return nil
	}

	valueFromBase := getByXPath(t, base, xpath)

	if isRegex {
		reg := regexp.MustCompile(extractRegex(templateVal))
		valueFromBaseAsStr := stringValue(t, valueFromBase)
		if !reg.MatchString(stringValue(t, valueFromBase)) {
			require.Fail(t, fmt.Sprintf("value at '%s' (%s) does not match '%s'", xpath, valueFromBaseAsStr, templateVal))
		}
	}

	valueToCopy := reflect.ValueOf(valueFromBase)
	return &valueToCopy
}

func isNumberLike(kind reflect.Kind) bool {
	return kind >= reflect.Int && kind <= reflect.Uint64 || kind >= reflect.Float32 && kind <= reflect.Float64
}

func stringValue(t testing.TB, v any) string {
	_, kind, _ := typeof(v)
	if isNumberLike(kind) {
		return fmt.Sprintf("%v", v)
	}
	if kind == reflect.String {
		return v.(string)
	}
	require.Fail(t, "value is not a string or number-like")
	return ""
}

func containsRegex(str string) bool {
	return len(str) > 2 && strings.HasPrefix(str, "/") && strings.HasSuffix(str, "/")
}

func extractRegex(str string) string {
	return str[1 : len(str)-1]
}

func typeof(v any) (reflect.Value, reflect.Kind, any) {
	if v == nil {
		return reflect.Value{}, reflect.Invalid, nil
	}
	rval := reflect.ValueOf(v)
	kind := reflect.TypeOf(v).Kind()

	if kind == reflect.Ptr {
		rval = rval.Elem()
		kind = rval.Kind()
	}

	return rval, kind, rval.Interface()
}

func addToXPath(xpath, key string) string {
	if xpath == "" {
		return key
	}
	return xpath + "/" + key
}
