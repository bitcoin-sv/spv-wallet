package admin

import (
	"github.com/bitcoin-sv/spv-wallet/config"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestXPubRegisterRoutes will test routes
func (ts *TestSuite) TestXPubRegisterRoutes() {
	ts.T().Run("test routes", func(t *testing.T) {

		testCases := []struct {
			method string
			url    string
		}{
			{"GET", "/" + config.APIVersion + "/admin/stats"},
			{"GET", "/" + config.APIVersion + "/admin/status"},
			{"POST", "/" + config.APIVersion + "/admin/access-keys/search"},
			{"POST", "/" + config.APIVersion + "/admin/access-keys/count"},
			{"POST", "/" + config.APIVersion + "/admin/destinations/search"},
			{"POST", "/" + config.APIVersion + "/admin/destinations/count"},
			{"POST", "/" + config.APIVersion + "/admin/paymail/get"},
			{"POST", "/" + config.APIVersion + "/admin/paymails/search"},
			{"POST", "/" + config.APIVersion + "/admin/paymails/count"},
			{"POST", "/" + config.APIVersion + "/admin/paymail/create"},
			{"DELETE", "/" + config.APIVersion + "/admin/paymail/delete"},
			{"POST", "/" + config.APIVersion + "/admin/transactions/search"},
			{"POST", "/" + config.APIVersion + "/admin/transactions/count"},
			{"POST", "/" + config.APIVersion + "/admin/transactions/record"},
			{"POST", "/" + config.APIVersion + "/admin/utxos/search"},
			{"POST", "/" + config.APIVersion + "/admin/utxos/count"},
			{"POST", "/" + config.APIVersion + "/admin/xpub"},
			{"POST", "/" + config.APIVersion + "/admin/xpubs/search"},
			{"POST", "/" + config.APIVersion + "/admin/xpubs/count"},
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
