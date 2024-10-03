package policy_test

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
)

const (
	wrongButReachable = "/wrong/url"
	arcURL            = "https://arc.taal.com"
	arcToken          = "mainnet_06770f425eb00298839a24a49cbdc02c"
)

func arcMockActivate(applyTimeout bool) *resty.Client {
	transport := httpmock.NewMockTransport()
	client := resty.New()
	client.GetClient().Transport = transport

	responder := func(status int, content string) func(req *http.Request) (*http.Response, error) {
		return func(req *http.Request) (*http.Response, error) {
			if applyTimeout {
				time.Sleep(100 * time.Millisecond)
			}
			if req.Header.Get("Authorization") != arcToken {
				return httpmock.NewStringResponse(http.StatusUnauthorized, ""), nil
			}
			res := httpmock.NewStringResponse(status, content)
			res.Header.Set("Content-Type", "application/json")
			return res, nil
		}
	}

	transport.RegisterResponder("GET", fmt.Sprintf("%s/v1/policy", arcURL), responder(http.StatusOK, `{
			"policy": {
				"maxscriptsizepolicy": 100000000,
				"maxtxsigopscountspolicy": 4294967295,
				"maxtxsizepolicy": 100000000,
				"miningFee": {
					"bytes": 1000,
					"satoshis": 1
				}
			},
			"timestamp": "2024-10-02T07:36:33.589144918Z"
		}`),
	)

	transport.RegisterResponder("GET", arcURL+wrongButReachable, responder(http.StatusNotFound, `{
			"message": "no matching operation was found"
		}`),
	)

	return client
}
