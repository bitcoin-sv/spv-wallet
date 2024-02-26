package base

import (
	"testing"

	"github.com/bitcoin-sv/spv-wallet/config"
	"github.com/stretchr/testify/assert"
)

// TestBaseRegisterRoutes will test routes
func (ts *TestSuite) TestBaseRegisterRoutes() {
	ts.T().Run("test routes", func(t *testing.T) {
		testCases := []struct {
			method string
			url    string
		}{
			{"GET", "/"},
			{"OPTIONS", "/"},
			{"HEAD", "/"},
			{"GET", "/" + config.HealthRequestPath},
			{"OPTIONS", "/" + config.HealthRequestPath},
			{"HEAD", "/" + config.HealthRequestPath},
			{"GET", "/debug/pprof/"},
			{"GET", "/debug/pprof/cmdline"},
			{"GET", "/debug/pprof/profile"},
			{"POST", "/debug/pprof/symbol"},
			{"GET", "/debug/pprof/trace"},
			{"GET", "/debug/pprof/goroutine"},
			{"GET", "/debug/pprof/heap"},
			{"GET", "/debug/pprof/threadcreate"},
			{"GET", "/debug/pprof/block"},
		}

		ts.Router.Routes()

		for _, testCase := range testCases {
			found := false
			for _, routeInfo := range ts.Router.Routes() {
				if testCase.url == routeInfo.Path && testCase.method == routeInfo.Method {
					assert.NotNil(t, routeInfo.HandlerFunc)
					found = true
					break
				}
			}
			assert.True(t, found)
		}
	})
}
