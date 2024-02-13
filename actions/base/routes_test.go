package base

import (
	"net/http"
	"testing"

	"github.com/BuxOrg/spv-wallet/config"
	"github.com/stretchr/testify/assert"
)

// TestBaseRegisterRoutes will test routes
func (ts *TestSuite) TestBaseRegisterRoutes() {
	ts.T().Run("test routes", func(t *testing.T) {
		// index
		handle, _, _ := ts.Router.HTTPRouter.Lookup(http.MethodGet, "/")
		assert.NotNil(t, handle)

		// options
		handle, _, _ = ts.Router.HTTPRouter.Lookup(http.MethodOptions, "/")
		assert.NotNil(t, handle)

		// head
		handle, _, _ = ts.Router.HTTPRouter.Lookup(http.MethodHead, "/")
		assert.NotNil(t, handle)

		// health
		handle, _, _ = ts.Router.HTTPRouter.Lookup(http.MethodGet, "/"+config.HealthRequestPath)
		assert.NotNil(t, handle)

		// health options
		handle, _, _ = ts.Router.HTTPRouter.Lookup(http.MethodOptions, "/"+config.HealthRequestPath)
		assert.NotNil(t, handle)

		// health head
		handle, _, _ = ts.Router.HTTPRouter.Lookup(http.MethodHead, "/"+config.HealthRequestPath)
		assert.NotNil(t, handle)
	})

	ts.T().Run("test debug profile routes", func(t *testing.T) {
		handle, _, _ := ts.Router.HTTPRouter.Lookup(http.MethodGet, "/debug/pprof/")
		assert.NotNil(t, handle)

		handle, _, _ = ts.Router.HTTPRouter.Lookup(http.MethodGet, "/debug/pprof/cmdline")
		assert.NotNil(t, handle)

		handle, _, _ = ts.Router.HTTPRouter.Lookup(http.MethodGet, "/debug/pprof/profile")
		assert.NotNil(t, handle)

		handle, _, _ = ts.Router.HTTPRouter.Lookup(http.MethodGet, "/debug/pprof/symbol")
		assert.NotNil(t, handle)

		handle, _, _ = ts.Router.HTTPRouter.Lookup(http.MethodGet, "/debug/pprof/trace")
		assert.NotNil(t, handle)

		handle, _, _ = ts.Router.HTTPRouter.Lookup(http.MethodGet, "/debug/pprof/goroutine")
		assert.NotNil(t, handle)

		handle, _, _ = ts.Router.HTTPRouter.Lookup(http.MethodGet, "/debug/pprof/heap")
		assert.NotNil(t, handle)

		handle, _, _ = ts.Router.HTTPRouter.Lookup(http.MethodGet, "/debug/pprof/threadcreate")
		assert.NotNil(t, handle)

		handle, _, _ = ts.Router.HTTPRouter.Lookup(http.MethodGet, "/debug/pprof/block")
		assert.NotNil(t, handle)
	})
}
