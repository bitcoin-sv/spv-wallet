package errorcases

import (
	"context"
	"testing"

	"github.com/bitcoin-sv/spv-wallet/api/manualtests"
	"github.com/bitcoin-sv/spv-wallet/api/manualtests/client"
	"github.com/stretchr/testify/require"
)

func TestCreateUserBadRequest(t *testing.T) {
	badRequests := map[string]struct {
		makeRequest func(manualtests.User, *manualtests.State, testing.TB) client.CreateUserJSONRequestBody
	}{
		"bad request: invalid public key": {
			makeRequest: func(_ manualtests.User, _ *manualtests.State, _ testing.TB) client.CreateUserJSONRequestBody {
				return client.CreateUserJSONRequestBody{
					PublicKey: "123",
				}
			},
		},
		"bad request: invalid paymail address": {
			makeRequest: func(user manualtests.User, _ *manualtests.State, _ testing.TB) client.CreateUserJSONRequestBody {
				return client.CreateUserJSONRequestBody{
					PublicKey: user.PublicKey,
					Paymail: &client.RequestsAddPaymail{
						Address: "invalid",
					},
				}
			},
		},
		"bad request: inconsistent paymail address": {
			makeRequest: func(user manualtests.User, _ *manualtests.State, _ testing.TB) client.CreateUserJSONRequestBody {
				return client.CreateUserJSONRequestBody{
					PublicKey: user.PublicKey,
					Paymail: &client.RequestsAddPaymail{
						Address: user.Address(),
						Domain:  user.Domain,
						Alias:   "inconsistent",
					},
				}
			},
		},
		"bad request: invalid domain": {
			makeRequest: func(user manualtests.User, _ *manualtests.State, t testing.TB) client.CreateUserJSONRequestBody {
				t.Skip("Ensure paymail domain validation is enabled then comment out this")
				return client.CreateUserJSONRequestBody{
					PublicKey: user.PublicKey,
					Paymail: &client.RequestsAddPaymail{
						Domain: "unknown",
						Alias:  user.Alias,
					},
				}
			},
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
			res, err := c.CreateUserWithResponse(context.Background(), req)
			require.NoError(t, err)

			manualtests.Print(res)

			require.Equal(t, 400, res.StatusCode())

			// IMPORTANT: DO NOT SAVE THE STATE HERE
		})
	}
}
