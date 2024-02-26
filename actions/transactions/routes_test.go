package transactions

import (
	"testing"

	"github.com/bitcoin-sv/spv-wallet/config"
	"github.com/stretchr/testify/assert"
)

// TestTransactionRegisterRoutes will test routes
func (ts *TestSuite) TestTransactionRegisterRoutes() {
	ts.T().Run("test routes", func(t *testing.T) {
		testCases := []struct {
			method string
			url    string
		}{
			{"GET", "/" + config.APIVersion + "/transaction"},
			{"POST", "/" + config.APIVersion + "/transaction"},
			{"PATCH", "/" + config.APIVersion + "/transaction"},
			{"GET", "/" + config.APIVersion + "/transaction/search"},
			{"POST", "/" + config.APIVersion + "/transaction/search"},
			{"POST", "/" + config.APIVersion + "/transaction/record"},
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
