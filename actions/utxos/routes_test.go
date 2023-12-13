package utxos

import (
	"net/http"
	"testing"

	"github.com/BuxOrg/bux-server/config"
	"github.com/stretchr/testify/assert"
)

// TestUtxoRegisterRoutes will test routes
func (ts *TestSuite) TestUtxoRegisterRoutes() {
	ts.T().Run("test routes", func(t *testing.T) {
		// get utxo
		handle, _, _ := ts.Router.HTTPRouter.Lookup(http.MethodGet, "/"+config.ApiVersion+"/utxo")
		assert.NotNil(t, handle)

		// count utxo
		handle, _, _ = ts.Router.HTTPRouter.Lookup(http.MethodPost, "/"+config.ApiVersion+"/utxo/count")
		assert.NotNil(t, handle)

		// search utxo
		handle, _, _ = ts.Router.HTTPRouter.Lookup(http.MethodPost, "/"+config.ApiVersion+"/utxo/search")
		assert.NotNil(t, handle)

		// unreserve utxo
		handle, _, _ = ts.Router.HTTPRouter.Lookup(http.MethodPatch, "/"+config.ApiVersion+"/utxo/unreserve")
		assert.NotNil(t, handle)
	})
}
