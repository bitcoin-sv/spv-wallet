package testabilities

import "testing"

func TestJSONTemplate(t *testing.T) {
	tests := map[string]struct {
		template string
		actual   string
	}{
		"flat structure": {
			template: `{"a": 1, "b": "/^[a-zA-Z]+$/", "c": "exact-match"}`,
			actual:   `{"a": 1, "b": "asd", "c": "exact-match"}`,
		},
		"no template": {
			template: `{"a": 1, "b": "b", "c": "c"}`,
			actual:   `{"a": 1, "b": "b", "c": "c"}`,
		},
		"regex for number": {
			template: `{"a": 1, "b": "/^\\d{1,3}$/", "c": 3}`,
			actual:   `{"a": 1, "b": 122, "c": 3}`,
		},
		"regex in nested obj": {
			template: `{"a": 1, "b": { "c": "/^[a-zA-Z]+$/", "d": "exact-match" }}`,
			actual:   `{"a": 1, "b": { "c": "asd", "d": "exact-match" }}`,
		},
		"any string regex": {
			template: `{"a": 1, "b": "/.*/", "c": "exact-match"}`,
			actual:   `{"a": 1, "b": "asd", "c": "exact-match"}`,
		},
		"regex in obj in array": {
			template: `{"a": 1, "b": [{"b": "/^[a-zA-Z]+$/", "c": "exact-match"}]}`,
			actual:   `{"a": 1, "b": [{"b": "asd", "c": "exact-match"}]}`,
		},
		"regex in obj in array (two elements)": {
			template: `{"a": 1, "b": [{"b": "/^[a-zA-Z]+$/", "c": "exact-match"}, {"d":"/^\\d{1,3}$/"}]}`,
			actual:   `{"a": 1, "b": [{"b": "asd", "c": "exact-match"}, {"d": 122}]}`,
		},
		"regex in value in array": {
			template: `{"a": 1, "b": ["/^[a-zA-Z]+$/"]}`,
			actual:   `{"a": 1, "b": ["asd"]}`,
		},
		"accept any (but must exist)": {
			template: `{"a": 1, "metadata": "*"}`,
			actual:   `{"a": 1, "metadata": {"b": "asd", "c": "exact-match"}}`,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assertJSONTemplate(t, test.template, test.actual)
		})
	}
}
