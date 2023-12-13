package transactions

import (
	"net/http"
	"testing"

	"github.com/BuxOrg/bux-server/config"
	"github.com/stretchr/testify/assert"
)

// TestTransactionRegisterRoutes will test routes
func (ts *TestSuite) TestTransactionRegisterRoutes() {
	ts.T().Run("test routes", func(t *testing.T) {

		// new transaction
		handle, _, _ := ts.Router.HTTPRouter.Lookup(http.MethodPost, "/"+config.ApiVersion+"/transaction")
		assert.NotNil(t, handle)

		// record transaction
		handle, _, _ = ts.Router.HTTPRouter.Lookup(http.MethodPost, "/"+config.ApiVersion+"/transaction/record")
		assert.NotNil(t, handle)

		// get transaction
		handle, _, _ = ts.Router.HTTPRouter.Lookup(http.MethodGet, "/"+config.ApiVersion+"/transaction")
		assert.NotNil(t, handle)

		// search transaction
		handle, _, _ = ts.Router.HTTPRouter.Lookup(http.MethodGet, "/"+config.ApiVersion+"/transaction/search")
		assert.NotNil(t, handle)

		// search transaction
		handle, _, _ = ts.Router.HTTPRouter.Lookup(http.MethodPost, "/"+config.ApiVersion+"/transaction/search")
		assert.NotNil(t, handle)

		// update transaction
		handle, _, _ = ts.Router.HTTPRouter.Lookup(http.MethodPatch, "/"+config.ApiVersion+"/transaction")
		assert.NotNil(t, handle)
	})
}
