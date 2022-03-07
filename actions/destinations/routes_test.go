package destinations

import (
	"net/http"
	"testing"

	"github.com/BuxOrg/bux-server/config"
	"github.com/stretchr/testify/assert"
)

// TestDestinationRegisterRoutes will test routes
func (ts *TestSuite) TestDestinationRegisterRoutes() {
	ts.T().Run("test routes", func(t *testing.T) {

		handle, _, _ := ts.Router.HTTPRouter.Lookup(http.MethodGet, "/"+config.CurrentMajorVersion+"/destination")
		assert.NotNil(t, handle)

		handle, _, _ = ts.Router.HTTPRouter.Lookup(http.MethodPost, "/"+config.CurrentMajorVersion+"/destination")
		assert.NotNil(t, handle)
	})
}
