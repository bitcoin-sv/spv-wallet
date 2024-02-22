package utils

import "net/http"

// HTTPInterface is the HTTP client interface
type HTTPInterface interface {
	Do(req *http.Request) (*http.Response, error)
}
