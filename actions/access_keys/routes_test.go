package accesskeys

import (
	"testing"

	"github.com/bitcoin-sv/spv-wallet/config"
	"github.com/stretchr/testify/assert"
)

func (ts *TestSuite) TestRegisterRoutes() {
	ts.T().Run("test routes", func(t *testing.T) {
		testCases := []struct {
			method string
			url    string
		}{
			{"GET", "/api/" + config.APIVersion + "/users/current/keys/:id"},
			{"POST", "/api/" + config.APIVersion + "/users/current/keys"},
			{"DELETE", "/api/" + config.APIVersion + "/users/current/keys/:id"},
			{"GET", "/api/" + config.APIVersion + "/users/current/keys"},
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
