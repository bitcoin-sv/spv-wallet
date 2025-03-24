package base

import (
	"testing"

	"github.com/bitcoin-sv/spv-wallet/config"
	"github.com/stretchr/testify/assert"
)

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
