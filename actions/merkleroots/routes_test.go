package merkleroots

import (
	"testing"

	"github.com/bitcoin-sv/spv-wallet/config"
	"github.com/stretchr/testify/require"
)

// TestMerklerootsRoutes will test routes
func (ts *TestSuite) TestMerklerootsRoutes() {

	ts.T().Run("test routes", func(t *testing.T) {
		testCases := []struct {
			method string
			url    string
		}{
			{"GET", "/api/" + config.APIVersion + "/merkleroots"},
		}

		for _, tc := range testCases {
			found := false
			for _, routeInfo := range ts.Router.Routes() {
				if tc.url == routeInfo.Path && tc.method == routeInfo.Method {
					require.NotNil(t, routeInfo.HandlerFunc)
					found = true
					break
				}
			}
			require.True(t, found)
		}
	})
}
