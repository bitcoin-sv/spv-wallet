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
		handle, _, _ := ts.Router.HTTPRouter.Lookup(http.MethodGet, "/"+config.APIVersion+"/utxo")
		assert.NotNil(t, handle)

		// count utxo
		handle, _, _ = ts.Router.HTTPRouter.Lookup(http.MethodPost, "/"+config.APIVersion+"/utxo/count")
		assert.NotNil(t, handle)

		// search utxo
		handle, _, _ = ts.Router.HTTPRouter.Lookup(http.MethodPost, "/"+config.APIVersion+"/utxo/search")
		assert.NotNil(t, handle)
	})
}
