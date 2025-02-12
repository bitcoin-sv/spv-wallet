package paymailmock

import (
	"bytes"
	"io"
	"net/http"
	"regexp"
	"sync"

	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
)

// CallDetails holds the details of a call made to the mocked server
type CallDetails struct {
	URL string

	RequestMethod string
	RequestBody   []byte

	ResponseBody []byte
	ResponseCode int
}

type httpSniffer struct {
	next   http.RoundTripper
	called map[string]CallDetails
	lock   sync.Mutex
}

func newHTTPSniffer() *httpSniffer {
	return &httpSniffer{
		called: make(map[string]CallDetails),
	}
}

func (s *httpSniffer) setTransport(next http.RoundTripper) {
	s.next = next
}

func (s *httpSniffer) getCallByRegex(r string) *CallDetails {
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

func (s *httpSniffer) RoundTrip(req *http.Request) (*http.Response, error) {
	var details CallDetails
	details.URL = req.URL.String()
	details.RequestMethod = req.Method

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
