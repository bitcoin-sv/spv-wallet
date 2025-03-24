package userapi

import (
	"context"
	"testing"

	"github.com/bitcoin-sv/spv-wallet/api/manualtests"
	"github.com/bitcoin-sv/spv-wallet/api/manualtests/client"
	"github.com/joomcode/errorx"
)

func TestDeleteCurrentUser(t *testing.T) {
	t.Skip("don't run yet")

	logger := manualtests.Logger()

	manualtests.APICallForCurrentUser(t).CallWithUpdateState(func(state manualtests.StateForCall, c *client.ClientWithResponses) (manualtests.Result, error) {
		err := state.User.MarkAsDeleted()
		if errorx.IsOfType(err, manualtests.UserDeleted) {
			logger.Warn().Msg("You are trying to delete a user that potentially is already deleted, it is tagged with deleted in state.yaml")
		}

		return c.DeleteCurrentUserWithResponse(context.Background())
	}).RequireSuccess()
}
