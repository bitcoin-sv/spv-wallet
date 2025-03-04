package errorcases

import (
	"context"
	"testing"

	"github.com/bitcoin-sv/spv-wallet/api/manualtests"
	"github.com/bitcoin-sv/spv-wallet/api/manualtests/client"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
)

func TestAddPaymailBadRequest(t *testing.T) {
	badRequests := map[string]struct {
		makeRequest func(manualtests.User, *manualtests.State, testing.TB) client.RequestsAddPaymail
		status      int
	}{
		"bad request: invalid paymail address": {
			makeRequest: func(user manualtests.User, _ *manualtests.State, _ testing.TB) client.RequestsAddPaymail {
				return client.RequestsAddPaymail{
					Address: "invalid",
				}
			},
			status: 400,
		},
		"bad request: inconsistent paymail address": {
			makeRequest: func(user manualtests.User, _ *manualtests.State, _ testing.TB) client.RequestsAddPaymail {
				return client.RequestsAddPaymail{
					Address: user.Address(),
					Domain:  user.Domain,
					Alias:   "inconsistent",
				}
			},
			status: 400,
		},
		"bad request: invalid domain": {
			makeRequest: func(user manualtests.User, _ *manualtests.State, t testing.TB) client.RequestsAddPaymail {
				t.Skip("Ensure paymail domain validation is enabled then comment out this")
				return client.RequestsAddPaymail{
					Domain: "unknown",
					Alias:  user.Alias,
				}
			},
			status: 400,
		},
		"bad request: invalid avatar url": {
			makeRequest: func(user manualtests.User, _ *manualtests.State, t testing.TB) client.RequestsAddPaymail {
				return client.RequestsAddPaymail{
					Address:   user.PaymailAddress(),
					AvatarURL: lo.ToPtr("https://[/]"),
				}
			},
			status: 422,
		},
	}
	for name, test := range badRequests {
		t.Run(name, func(t *testing.T) {
			state := manualtests.NewState()
			err := state.Load()
			require.NoError(t, err)

			c, err := state.AdminClient()
			require.NoError(t, err)

			user := UserDefinitionForMakingBadRequests(state.Domain)

			req := test.makeRequest(user, state, t)
			res, err := c.AddPaymailToUserWithResponse(
				context.Background(),
				state.CurrentUser().ID,
				req,
			)
			require.NoError(t, err)

			manualtests.Print(res)

			require.Equal(t, test.status, res.StatusCode())

			// IMPORTANT: DO NOT SAVE THE STATE HERE
		})
	}
}
