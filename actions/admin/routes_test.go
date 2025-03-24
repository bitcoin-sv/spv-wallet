package admin

import (
	"testing"

	"github.com/bitcoin-sv/spv-wallet/config"
	"github.com/stretchr/testify/assert"
)

func (ts *TestSuite) TestXPubRegisterRoutes() {
	ts.T().Run("test routes", func(t *testing.T) {

		testCases := []struct {
			method string
			url    string
		}{
			// tx
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
