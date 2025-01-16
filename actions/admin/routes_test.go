package admin

import (
	"testing"

	"github.com/bitcoin-sv/spv-wallet/config"
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

			// tx
			{"GET", "/" + config.APIVersion + "/admin/transactions/:id"},     // get tx by id old
			{"GET", "/" + config.APIVersion + "/admin/transactions"},         // search old
			{"GET", "/api/" + config.APIVersion + "/admin/transactions/:id"}, // get tx by id
			{"GET", "/api/" + config.APIVersion + "/admin/transactions"},     // search

			// contacts
			{"POST", "/api/" + config.APIVersion + "/admin/invitations/:id"},   // accept
			{"DELETE", "/api/" + config.APIVersion + "/admin/invitations/:id"}, // reject
			{"GET", "/api/" + config.APIVersion + "/admin/contacts"},           // search
			{"DELETE", "/api/" + config.APIVersion + "/admin/contacts/:id"},    // delete
			{"PUT", "/api/" + config.APIVersion + "/admin/contacts/:id"},       // update

			// access keys
			{"GET", "/api/" + config.APIVersion + "/admin/users/keys"}, // search

			// paymails
			{"GET", "/api/" + config.APIVersion + "/admin/paymails/:id"},    // get paymail by id
			{"GET", "/api/" + config.APIVersion + "/admin/paymails"},        // search
			{"POST", "/api/" + config.APIVersion + "/admin/paymails"},       // create
			{"DELETE", "/api/" + config.APIVersion + "/admin/paymails/:id"}, // delete

			// utxos
			{"GET", "/api/" + config.APIVersion + "/admin/utxos"}, // get utxo

			// webhooks
			{"POST", "/api/" + config.APIVersion + "/admin/webhooks/subscriptions"},   // subscribe
			{"DELETE", "/api/" + config.APIVersion + "/admin/webhooks/subscriptions"}, // unsubscribe

			// xpubs
			{"POST", "/api/" + config.APIVersion + "/admin/users"}, // create
			{"GET", "/api/" + config.APIVersion + "/admin/users"},  // search
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
