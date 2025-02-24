//nolint:nolintlint,revive
package adminapi

import (
	"context"
	"testing"

	"github.com/bitcoin-sv/spv-wallet/api/manualtests"
	"github.com/bitcoin-sv/spv-wallet/api/manualtests/client"
	"github.com/stretchr/testify/require"
)

func TestAdminAPIRead(t *testing.T) {
	state := manualtests.NewState()
	err := state.Load()
	require.NoError(t, err)

	var calls = map[string]struct {
		call func(*client.ClientWithResponses) (manualtests.Result, error)
	}{
		"adminStatus": {
			call: func(c *client.ClientWithResponses) (manualtests.Result, error) {
				return c.AdminStatusWithResponse(context.Background())
			},
		},
		"userById": {
			call: func(c *client.ClientWithResponses) (manualtests.Result, error) {
				return c.UserByIdWithResponse(context.Background(), state.CurrentUser().ID)
			},
		},
	}
	for name, endpoint := range calls {
		t.Run(name, func(t *testing.T) {
			manualtests.APICallForAdmin(t).Call(endpoint.call).RequireSuccess()
		})
	}
}
