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
		handle, _, _ := ts.Router.HTTPRouter.Lookup(http.MethodGet, "/"+config.CurrentMajorVersion+"/utxo")
		assert.NotNil(t, handle)

		// count utxo
		handle, _, _ = ts.Router.HTTPRouter.Lookup(http.MethodPost, "/"+config.CurrentMajorVersion+"/utxo/count")
		assert.NotNil(t, handle)

		// search utxo
		handle, _, _ = ts.Router.HTTPRouter.Lookup(http.MethodPost, "/"+config.CurrentMajorVersion+"/utxo/search")
		assert.NotNil(t, handle)

		// unreserve utxo
		handle, _, _ = ts.Router.HTTPRouter.Lookup(http.MethodPatch, "/"+config.CurrentMajorVersion+"/utxo/unreserve")
		assert.NotNil(t, handle)
	})
}
