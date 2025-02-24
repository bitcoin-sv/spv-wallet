//nolint:nolintlint,revive
package adminapi

import (
	"context"
	"testing"

	"github.com/bitcoin-sv/spv-wallet/api/manualtests"
	"github.com/bitcoin-sv/spv-wallet/api/manualtests/client"
)

func TestAdminStatus(t *testing.T) {
	manualtests.APICallForAdmin(t).Call(
		func(c *client.ClientWithResponses) (manualtests.Result, error) {
			return c.AdminStatusWithResponse(context.Background())
		},
	).RequireSuccess()
}

func TestAdmin_UserById(t *testing.T) {
	manualtests.APICallForAdmin(t).CallWithState(
		func(state manualtests.StateForCall, c *client.ClientWithResponses) (manualtests.Result, error) {
			return c.UserByIdWithResponse(context.Background(), state.CurrentUser().ID)
		},
	).RequireSuccess()
}
