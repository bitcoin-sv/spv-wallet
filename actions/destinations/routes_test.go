package destinations

import (
	"testing"

	"github.com/bitcoin-sv/spv-wallet/config"
	"github.com/stretchr/testify/assert"
)

// TestDestinationRegisterRoutes will test routes
func (ts *TestSuite) TestDestinationRegisterRoutes() {
	ts.T().Run("test routes", func(t *testing.T) {
		testCases := []struct {
			method string
			url    string
		}{
			{"GET", "/" + config.APIVersion + "/destination"},
			{"POST", "/" + config.APIVersion + "/destination"},
			{"PATCH", "/" + config.APIVersion + "/destination"},
			{"GET", "/" + config.APIVersion + "/destination/search"},
			{"POST", "/" + config.APIVersion + "/destination/search"},
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
