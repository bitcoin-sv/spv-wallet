package xpubs

import (
	"net/http"
	"testing"

	"github.com/BuxOrg/spv-wallet/config"
	"github.com/stretchr/testify/assert"
)

// TestXPubRegisterRoutes will test routes
func (ts *TestSuite) TestXPubRegisterRoutes() {
	ts.T().Run("test routes", func(t *testing.T) {

		handle, _, _ := ts.Router.HTTPRouter.Lookup(http.MethodPost, "/"+config.APIVersion+"/xpub")
		assert.NotNil(t, handle)

		handle, _, _ = ts.Router.HTTPRouter.Lookup(http.MethodGet, "/"+config.APIVersion+"/xpub")
		assert.NotNil(t, handle)

		handle, _, _ = ts.Router.HTTPRouter.Lookup(http.MethodPatch, "/"+config.APIVersion+"/xpub")
		assert.NotNil(t, handle)
	})
}
