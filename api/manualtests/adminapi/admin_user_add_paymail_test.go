package adminapi

import (
	"context"
	"testing"

	"github.com/bitcoin-sv/spv-wallet/api/manualtests"
	"github.com/bitcoin-sv/spv-wallet/api/manualtests/client"
)

func TestAddPaymail(t *testing.T) {
	t.Skip("Don't run it yet")

	manualtests.APICallForAdmin(t).
		CallWithUpdateState(func(state manualtests.StateForCall, c *client.ClientWithResponses) (manualtests.Result, error) {
			user := state.CurrentUser()

			additionalAlias := user.CreateAdditionalAlias()

			return c.AddPaymailToUserWithResponse(
				context.Background(),
				state.CurrentUser().ID,
				client.RequestsAddPaymail{
					Alias:      additionalAlias.String(),
					Domain:     user.Domain,
					AvatarURL:  user.AvatarURL(),
					PublicName: additionalAlias.PublicName(),
				},
			)
		}).
		RequireSuccess()
}
