package errorcases

import (
	"context"
	"testing"

	"github.com/bitcoin-sv/spv-wallet/api/manualtests"
	"github.com/bitcoin-sv/spv-wallet/api/manualtests/client"
	"github.com/stretchr/testify/require"
)

func TestAddPaymailBadRequest(t *testing.T) {
	badRequests := map[string]struct {
		makeRequest func(manualtests.User, *manualtests.State, testing.TB) client.RequestsAddPaymail
	}{
		"bad request: invalid paymail address": {
			makeRequest: func(user manualtests.User, _ *manualtests.State, _ testing.TB) client.RequestsAddPaymail {
				return client.RequestsAddPaymail{
					Address: "invalid",
				}
			},
		},
		"bad request: inconsistent paymail address": {
			makeRequest: func(user manualtests.User, _ *manualtests.State, _ testing.TB) client.RequestsAddPaymail {
				return client.RequestsAddPaymail{
					Address: user.Address(),
					Domain:  user.Domain,
					Alias:   "inconsistent",
				}
			},
		},
		"bad request: invalid domain": {
			makeRequest: func(user manualtests.User, _ *manualtests.State, t testing.TB) client.RequestsAddPaymail {
				t.Skip("Ensure paymail domain validation is enabled then comment out this")
				return client.RequestsAddPaymail{
					Domain: "unknown",
					Alias:  user.Alias,
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

			user := manualtests.User{
				Alias:     "badrequest",
				Domain:    state.Domain,
				Xpriv:     "xprv9s21ZrQH143K3TzwskUB1NW5iHx4EH7cxquXfwpFR5HrWe6HQrYvECJsj3sg1DJhWhwjtw5WdXwje8pkyuuvzJUingwo4f5BkD5dNubfNUn",
				Xpub:      "xpub661MyMwAqRbcFx5Qyn1BNWSpGKnYdjqUL4q8ULDryQpqPSRRxPsAmzdMaLnoyUzLAQ5ukXgMZjYs5LfNfsPFwBoSwxChePB1DxKvyFz6F67",
				PublicKey: "035dcb59eb7b5c5982ba6fbbccffbb2460f8daa07c5f9a21f2c2cf0845dcc6dda9",
				ID:        "1QCibNnc8CK7bzTdupM91f4PKixhMqqJQw",
			}

			req := test.makeRequest(user, state, t)
			res, err := c.AddPaymailToUserWithResponse(
				context.Background(),
				state.CurrentUser().ID,
				req,
			)
			require.NoError(t, err)

			manualtests.Print(res)

			require.Equal(t, 400, res.StatusCode())

			// IMPORTANT: DO NOT SAVE THE STATE HERE
		})
	}
}
