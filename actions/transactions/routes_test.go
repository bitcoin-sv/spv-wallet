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

		handle, _, _ := ts.Router.HTTPRouter.Lookup(http.MethodPost, "/"+config.CurrentMajorVersion+"/transaction")
		assert.NotNil(t, handle)

		handle, _, _ = ts.Router.HTTPRouter.Lookup(http.MethodPost, "/"+config.CurrentMajorVersion+"/transaction/record")
		assert.NotNil(t, handle)

		handle, _, _ = ts.Router.HTTPRouter.Lookup(http.MethodGet, "/"+config.CurrentMajorVersion+"/transaction")
		assert.NotNil(t, handle)
	})
}
