package jsonrequire

import "testing"

func TestJSONPlaceholders(t *testing.T) {
	tests := map[string]struct {
		template string
		actual   string
	}{
		"flat structure": {
			template: `{"a": 1, "b": "/^[a-zA-Z]+$/", "c": "exact-match"}`,
			actual:   `{"a": 1, "b": "asd", "c": "exact-match"}`,
		},
		"no placeholders": {
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
			Match(t, test.template, nil, test.actual)
		})
	}
}

func TestJSONTemplate(t *testing.T) {
	tests := map[string]struct {
		template string
		actual   string
		params   map[string]any
	}{
		"match timestamp": {
			template: `{"timestamp": "{{ matchTimestamp }}", "withTimezone": "{{ matchTimestamp }}"}`,
			actual:   `{"timestamp": "2024-12-03T09:39:42.5515364Z", "withTimezone": "2024-12-03T15:21:27.1080542+01:00"}`,
		},
		"or empty with empty value": {
			template: `{"timestamp": "{{ matchTimestamp | orEmpty }}"}`,
			actual:   `{"timestamp": ""}`,
		},
		"or empty with actual value": {
			template: `{"timestamp": "{{ matchTimestamp | orEmpty }}"}`,
			actual:   `{"timestamp": "2024-12-03T09:39:42.5515364Z"}`,
		},
		"match URL http": {
			template: `{ "url": "{{ matchURL }}" }`,
			actual:   `{ "url": "http://example.com" }`,
		},
		"match URL https": {
			template: `{ "url": "{{ matchURL }}" }`,
			actual:   `{ "url": "https://example.com" }`,
		},
		"match URL ftp": {
			template: `{ "url": "{{ matchURL }}" }`,
			actual:   `{ "url": "ftp://example.com" }`,
		},
		"match URL localhost": {
			template: `{ "url": "{{ matchURL }}" }`,
			actual:   `{ "url": "http://localhost" }`,
		},
		"match URL with path and search params": {
			template: `{ "url": "{{ matchURL }}" }`,
			actual:   `{ "url": "https://example.com/path?hello=123" }`,
		},
		"match URL with orEmpty on empty value": {
			template: `{ "url": "{{ matchURL | orEmpty }}" }`,
			actual:   `{ "url": "" }`,
		},
		"match URL with orEmpty on actual value": {
			template: `{ "url": "{{ matchURL | orEmpty }}" }`,
			actual:   `{ "url": "https://example.com" }`,
		},
		"match string ID with 64 characters": {
			template: `{ "id": "{{ matchID64 }}" }`,
			actual:   `{ "id": "d425432e0d10a46af1ec6d00f380e9581ebf7907f3486572b3cd561a4c326e14" }`,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			Match(t, test.template, test.params, test.actual)
		})
	}
}
