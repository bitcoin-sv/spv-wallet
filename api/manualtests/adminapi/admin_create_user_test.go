//nolint:nolintlint,revive
package adminapi_test // This should make it easier to run mutations separately from the queries

import (
	"context"
	"testing"

	bip32 "github.com/bitcoin-sv/go-sdk/compat/bip32"
	"github.com/bitcoin-sv/spv-wallet/api/manualtests"
	"github.com/bitcoin-sv/spv-wallet/api/manualtests/client"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {
	t.Skip("Don't run it yet")

	manualtests.APICallForAdmin(t).
		CallWithUpdateState(func(state manualtests.StateForCall, c *client.ClientWithResponses) (manualtests.Result, error) {
			xpriv, xpub, err := bip32.GenerateHDKeyPair(bip32.RecommendedSeedLength)
			require.NoError(t, err)

			user, err := state.NewUser(xpriv, xpub)
			require.NoError(t, err)

			return c.CreateUserWithResponse(context.Background(), client.CreateUserJSONRequestBody{
				Paymail: &client.RequestsAddPaymail{
					Address:    user.PaymailAddress(),
					AvatarURL:  lo.ToPtr(user.AvatarURL()),
					PublicName: lo.ToPtr(user.PublicName()),
				},
				PublicKey: user.PublicKey,
			})
		}).
		RequireSuccess()
}

func TestCreateUserAgain(t *testing.T) {
	t.Skip("Don't run it yet")

	manualtests.APICallForAdmin(t).
		CallWithUpdateState(func(state manualtests.StateForCall, c *client.ClientWithResponses) (manualtests.Result, error) {
			user := state.CurrentUser()

			user.RemoveTag("deleted")

			return c.CreateUserWithResponse(context.Background(), client.CreateUserJSONRequestBody{
				Paymail: &client.RequestsAddPaymail{
					Address:    user.PaymailAddress(),
					AvatarURL:  lo.ToPtr(user.AvatarURL()),
					PublicName: lo.ToPtr(user.PublicName()),
				},
				PublicKey: user.PublicKey,
			})
		})
}
