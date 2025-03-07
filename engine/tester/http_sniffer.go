package tester

import (
	"bytes"
	"io"
	"net/http"
	"regexp"
	"strings"
	"sync"

	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
)

type Headers map[string]string

func (h Headers) Get(key string) (string, bool) {
	val, ok := h[strings.ToLower(key)]
	return val, ok
}

func (h Headers) Set(key, value string) {
	h[strings.ToLower(key)] = value
}

// CallDetails holds the details of a call made to the mocked server
type CallDetails struct {
	URL string

	RequestMethod  string
	RequestBody    []byte
	RequestHeaders Headers

	ResponseBody []byte
	ResponseCode int
}

type HTTPSniffer struct {
	next   http.RoundTripper
	called map[string]CallDetails
	lock   sync.Mutex
}

func NewHTTPSniffer(next http.RoundTripper) *HTTPSniffer {
	return &HTTPSniffer{
		called: make(map[string]CallDetails),
		next:   next,
	}
}

func (s *HTTPSniffer) GetCallByRegex(r string) *CallDetails {
	reg := regexp.MustCompile(r)
	s.lock.Lock()
	defer s.lock.Unlock()
	for url, details := range s.called {
		if reg.MatchString(url) {
			return &details
		}
	}
	return nil
}

func (s *HTTPSniffer) RoundTrip(req *http.Request) (*http.Response, error) {
	var details CallDetails
	details.URL = req.URL.String()
	details.RequestMethod = req.Method

	details.RequestHeaders = make(map[string]string)
	for k, v := range req.Header {
		details.RequestHeaders.Set(k, v[0])
	}

	var err error
	if req.Body != nil {
		details.RequestBody, err = io.ReadAll(req.Body)
		if err != nil {
			panic(spverrors.Wrapf(err, "cannot read request body"))
		}
		req.Body = io.NopCloser(bytes.NewReader(details.RequestBody)) // Restore body after reading
	}

	resp, err := s.next.RoundTrip(req)
	if err != nil {
		return nil, spverrors.Wrapf(err, "error in round trip")
	}

	details.ResponseCode = resp.StatusCode
	if resp.Body != nil {
		details.ResponseBody, err = io.ReadAll(resp.Body)
		if err != nil {
			panic(spverrors.Wrapf(err, "cannot read response body"))
		}
		resp.Body = io.NopCloser(bytes.NewReader(details.ResponseBody)) // Restore body after reading
	}

	s.lock.Lock()
	defer s.lock.Unlock()
	s.called[details.URL] = details

	return resp, nil
}
