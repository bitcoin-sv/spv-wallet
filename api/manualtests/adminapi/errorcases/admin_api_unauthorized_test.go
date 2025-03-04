//nolint:nolintlint,revive
package errorcases

import (
	"context"
	"testing"

	"github.com/bitcoin-sv/spv-wallet/api/manualtests"
	"github.com/bitcoin-sv/spv-wallet/api/manualtests/client"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
)

func TestUnauthorized(t *testing.T) {
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
		"createUser": {
			call: func(c *client.ClientWithResponses) (manualtests.Result, error) {
				return c.CreateUserWithResponse(context.Background(), client.CreateUserJSONRequestBody{
					PublicKey: "1234567890",
				})
			},
		},
		"addPaymailToUser": {
			call: func(c *client.ClientWithResponses) (manualtests.Result, error) {
				return c.AddPaymailToUserWithResponse(
					context.Background(),
					state.CurrentUser().ID,
					client.RequestsAddPaymail{
						Alias:      "unauthorized",
						Domain:     state.Domain,
						PublicName: lo.ToPtr("Unauthorized"),
					},
				)
			},
		},
	}
	for name, endpoint := range calls {
		t.Run(name+"_user", func(t *testing.T) {
			manualtests.APICallForCurrentUser(t).
				Call(endpoint.call).
				RequireUnauthorizedForUserOnAdminAPI()
		})
	}

	for name, endpoint := range calls {
		t.Run(name+"_anonymous", func(t *testing.T) {
			manualtests.APICallForAnonymous(t).
				Call(endpoint.call).
				RequireUnauthorizedForAnonymous()
		})
	}

	for name, endpoint := range calls {
		t.Run(name+"_unknown_user", func(t *testing.T) {
			manualtests.APICallForUnknownUser(t).
				Call(endpoint.call).
				RequireUnauthorizedForUnknownUser()
		})
	}
}
