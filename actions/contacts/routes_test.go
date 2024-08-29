package contacts

import (
	"testing"

	"github.com/bitcoin-sv/spv-wallet/config"
	"github.com/stretchr/testify/assert"
)

// TestContactsRegisterRoutes will test routes
func (ts *TestSuite) TestContactsRegisterRoutes() {
	ts.T().Run("test routes", func(t *testing.T) {
		testCases := []struct {
			method string
			url    string
		}{
			{"PUT", "/" + config.APIVersion + "/contact/:paymail"},
			{"PATCH", "/" + config.APIVersion + "/contact/accepted/:paymail"},
			{"PATCH", "/" + config.APIVersion + "/contact/rejected/:paymail"},
			{"PATCH", "/" + config.APIVersion + "/contact/confirmed/:paymail"},
			{"PATCH", "/" + config.APIVersion + "/contact/unconfirmed/:paymail"},
			{"POST", "/" + config.APIVersion + "/contact/search"},

			{"PUT", "/api/" + config.APIVersion + "/contacts/:paymail"},
			{"DELETE", "/api/" + config.APIVersion + "/contacts/:paymail"},
			{"POST", "/api/" + config.APIVersion + "/contacts/:paymail/confirmation"},
			{"DELETE", "/api/" + config.APIVersion + "/contacts/:paymail/confirmation"},
			{"GET", "/api/" + config.APIVersion + "/contacts"},
			{"GET", "/api/" + config.APIVersion + "/contacts/:paymail"},

			{"POST", "/api/" + config.APIVersion + "/invitations/:paymail/contacts"},
			{"DELETE", "/api/" + config.APIVersion + "/invitations/:paymail"},
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
