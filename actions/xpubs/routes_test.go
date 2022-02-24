package xpubs

import (
	"net/http"
	"testing"

	"github.com/BuxOrg/bux-server/config"
	"github.com/stretchr/testify/assert"
)

// TestXPubRegisterRoutes will test routes
func (ts *TestSuite) TestXPubRegisterRoutes() {
	ts.T().Run("test routes", func(t *testing.T) {

		handle, _, _ := ts.Router.HTTPRouter.Lookup(http.MethodPost, "/"+config.CurrentMajorVersion+"/xpubs")
		assert.NotNil(t, handle)

		handle, _, _ = ts.Router.HTTPRouter.Lookup(http.MethodGet, "/"+config.CurrentMajorVersion+"/xpub")
		assert.NotNil(t, handle)
	})
}
