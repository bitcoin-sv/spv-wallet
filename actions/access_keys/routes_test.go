package accesskeys

import (
	"testing"

	"github.com/bitcoin-sv/spv-wallet/config"
	"github.com/stretchr/testify/assert"
)

// TestBaseRegisterRoutes will test routes
func (ts *TestSuite) TestRegisterRoutes() {
	ts.T().Run("test routes", func(t *testing.T) {
		testCases := []struct {
			method string
			url    string
		}{
			{"GET", "/" + config.APIVersion + "/access-key"},
			{"POST", "/" + config.APIVersion + "/access-key"},
			{"DELETE", "/" + config.APIVersion + "/access-key"},
			{"POST", "/" + config.APIVersion + "/access-key/search"},
			{"GET", "/" + config.APIVersion + "/access-key/search"},
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
