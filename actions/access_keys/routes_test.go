package accesskeys

import (
	"net/http"
	"testing"

	"github.com/BuxOrg/spv-wallet/config"
	"github.com/stretchr/testify/assert"
)

// TestBaseRegisterRoutes will test routes
func (ts *TestSuite) TestRegisterRoutes() {
	ts.T().Run("test routes", func(t *testing.T) {

		// gey key
		handle, _, _ := ts.Router.HTTPRouter.Lookup(http.MethodGet, "/"+config.APIVersion+"/access-key")
		assert.NotNil(t, handle)

		// search key
		handle, _, _ = ts.Router.HTTPRouter.Lookup(http.MethodGet, "/"+config.APIVersion+"/access-key/search")
		assert.NotNil(t, handle)

		// search key
		handle, _, _ = ts.Router.HTTPRouter.Lookup(http.MethodPost, "/"+config.APIVersion+"/access-key/search")
		assert.NotNil(t, handle)

		// create key
		handle, _, _ = ts.Router.HTTPRouter.Lookup(http.MethodPost, "/"+config.APIVersion+"/access-key")
		assert.NotNil(t, handle)

		// delete key
		handle, _, _ = ts.Router.HTTPRouter.Lookup(http.MethodDelete, "/"+config.APIVersion+"/access-key")
		assert.NotNil(t, handle)
	})
}
